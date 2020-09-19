package db

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/diamondburned/facechat/backend/internal/httperr"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Tx struct {
	*ReadTx
}

func (tx *Tx) Register(username, password, email string) (*facechat.User, *facechat.Session, error) {
	id := facechat.GenerateID()

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error hashing password")
	}

	_, err = tx.tx.Exec(
		"INSERT INTO users(id, name, pass, email) VALUES (?, ?)",
		id, username, hashed, email,
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error inserting user into db")
	}

	user := &facechat.User{
		ID:    id,
		Name:  username,
		Email: email,
	}
	ses, err := tx.insertSession(id)
	if err != nil {
		return nil, nil, err
	}

	return user, ses, nil
}

func (tx *Tx) Login(email, password string) (*facechat.Session, error) {
	var hashed []byte
	var id facechat.ID
	row := tx.tx.QueryRow("SELECT pass, id FROM users WHERE email = ?", email)
	err := row.Scan(&hashed, &id)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(hashed, []byte(password))
	if err != nil {
		return nil, httperr.Wrap(err, http.StatusUnauthorized, "invalid password")
	}

	ses, err := tx.insertSession(id)
	if err != nil {
		return nil, err
	}

	return ses, err
}

func (tx *Tx) UpdateSession(token string, expiry time.Time) error {
	return errors.New("unimplemented")
}

func (tx *Tx) DeleteSession(token string) error {
	return errors.New("unimplemented")
}

func (tx *Tx) insertSession(user facechat.ID) (*facechat.Session, error) {
	token, err := randToken()
	if err != nil {
		return nil, err
	}

	ses := facechat.Session{
		UserID: user,
		Token:  token,
		Expiry: time.Now().Add(facechat.SessionTimeout),
	}

	_, err = tx.tx.Exec(
		"INSERT INTO users(user_id, token, expiry) VALUES (?, ?, ?)",
		ses.UserID, ses.Token, ses.Expiry,
	)
	if err != nil {
		return nil, err
	}
	return &ses, nil
}

type ReadTx struct {
	tx *sqlx.Tx
}

func (tx *ReadTx) Session(token string) (*facechat.Session, error) {
	return nil, errors.New("unimplemented")
}

func (tx *ReadTx) User(id facechat.ID) (*facechat.User, error) {
	var user facechat.User
	row := tx.tx.QueryRowx("SELECT * FROM users WHERE id = ?", id)
	if err := row.StructScan(&user); err != nil {
		return nil, errors.Wrap(err, "error scanning row into user")
	}
	return &user, nil
}

func (tx *ReadTx) UserAccounts(id facechat.ID) ([]facechat.Account, error) {
	var accounts []facechat.Account
	rows, err := tx.tx.Queryx("SELECT * FROM accounts WHERE user_id = ?", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var account facechat.Account
		err := rows.StructScan(&account)
		if err != nil {
			return nil, errors.Wrap(err, "error scanning row into account")
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func randToken() (string, error) {
	var token = make([]byte, 32)

	if _, err := rand.Read(token); err != nil {
		return "", errors.Wrap(err, "failed to generate randomness")
	}

	return base64.RawURLEncoding.EncodeToString(token), nil
}
