package db

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/diamondburned/facechat/backend/internal/httperr"
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

// func (tx *Tx) UpdateSession(token string, expiry time.Time) error {
// 	return errors.New("unimplemented")
// }

// func (tx *Tx) DeleteSession(token string) error {
// 	return errors.New("unimplemented")
// }

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

func (tx *Tx) AddAccount(acc facechat.Account) error {
	panic("Implement me")
}

func (tx *Tx) SetRelationship(userID, targetID facechat.ID, rel facechat.RelationshipType) error {
	// TODO(diamond to sam): this requires further discussion, contact me.
	panic("Implement me")
}

func (tx *Tx) CreatePublicLobby(name string, lvl facechat.SecretLevel) (*facechat.Room, error) {
	room := facechat.Room{
		Type:  facechat.PublicLobby,
		Name:  name,
		Level: lvl,
	}

	if err := tx.createRoom(&room); err != nil {
		return nil, err
	}

	return &room, nil
}

func (tx *Tx) CreatePrivateRoom(targetUser facechat.ID) (*facechat.Room, error) {
	r, err := tx.Relationship(targetUser)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get relationship")
	}

	if r != facechat.Friend {
		return nil, err
	}

	// TODO: write a complex query that searches room_participants where both
	// current user and target user has the same room_id (aka shares the same
	// room), then query rooms for type == PrivateRoom

	room := facechat.Room{
		Type: facechat.PrivateRoom,
		// TODO: query the target user for the room name.
		Name: "",
		// TODO: if UpdateRoom method, DO NOT change level if the room type is a
		// Private one.
		Level: facechat.FullyOpen,
	}

	if err := tx.createRoom(&room); err != nil {
		return nil, err
	}

	return &room, nil
}

func (tx *Tx) createRoom(room *facechat.Room) error {
	(*room).ID = 0 // TODO
	panic("Implement me")
}

func (tx *Tx) CreateMessage(room, author facechat.ID, content string) (*facechat.Message, error) {
	msg := facechat.Message{
		Type:     facechat.NormalMessage,
		RoomID:   room,
		AuthorID: author,
		Markdown: content,
	}

	if err := tx.AddMessage(&msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

func (tx *Tx) AddMessage(msg *facechat.Message) error {
	// do a room_participants check on AuthorID
	// *msg.ID = setIDhere
	panic("Implement me")
}

func (tx *Tx) JoinRoom(room, user facechat.ID) error {
	// room_participants
	// TYPE == facechat.PublicLobby (MUST DO THIS)
	panic("Implement me")
}

func (tx *Tx) LeaveRoom(room, user facechat.ID) error {
	// room_participants
	// TYPE == facechat.PublicLobby (MUST DO THIS)
	panic("Implement me")
}

func randToken() (string, error) {
	var token = make([]byte, 32)

	if _, err := rand.Read(token); err != nil {
		return "", errors.Wrap(err, "failed to generate randomness")
	}

	return base64.RawURLEncoding.EncodeToString(token), nil
}
