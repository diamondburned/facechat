package db

import (
	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Tx struct {
	*ReadTx
}

func (tx *Tx) Register(username, password, email string) (*facechat.User, error) {
	id := facechat.GenerateID()
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "error hashing password")
	}
	_, err = tx.tx.Exec("INSERT INTO users(id, name, pass, email) VALUES (?, ?)", id, username, hashed, email)
	if err != nil {
		return nil, errors.Wrap(err, "error inserting user into db")
	}
	user := &facechat.User{
		ID:    id,
		Name:  username,
		Email: email,
	}
	return user, nil
}

type ReadTx struct {
	tx *sqlx.Tx
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

func (tx *ReadTx) UserVerifyPassword(email, pass string) error {
	var hashed []byte
	row := tx.tx.QueryRow("SELECT pass FROM users WHERE email = ?", email)
	err := row.Scan(&hashed)
	if err != nil {
		return errors.Wrap(err, "error scanning row into password")
	}
	return bcrypt.CompareHashAndPassword(hashed, []byte(pass))
}
