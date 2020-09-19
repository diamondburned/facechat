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
	tx.UserID = facechat.GenerateID()

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error hashing password")
	}

	_, err = tx.tx.Exec(
		"INSERT INTO users(id, name, pass, email) VALUES (?, ?)",
		tx.UserID, username, hashed, email,
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error inserting user into db")
	}

	user := &facechat.User{
		ID:    tx.UserID,
		Name:  username,
		Email: email,
	}
	ses, err := tx.insertSession()
	if err != nil {
		return nil, nil, err
	}

	return user, ses, nil
}

func (tx *Tx) Login(email, password string) (*facechat.Session, error) {
	var hashed []byte

	row := tx.tx.QueryRow("SELECT pass, id FROM users WHERE email = ?", email)
	err := row.Scan(&hashed, &tx.UserID)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(hashed, []byte(password))
	if err != nil {
		return nil, httperr.Wrap(err, http.StatusUnauthorized, "invalid password")
	}

	ses, err := tx.insertSession()
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

func (tx *Tx) insertSession() (*facechat.Session, error) {
	token, err := randToken()
	if err != nil {
		return nil, err
	}

	ses := facechat.Session{
		UserID: tx.UserID,
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
	_, err := tx.tx.Exec(
		"INSERT INTO accounts VALUES($1, $2, $3, $4, $5)",
		acc.Service, acc.Name, acc.URL, acc.Data, acc.UserID,
	)
	return errors.Wrap(err, "failed to add account")
}

func (tx *Tx) SetRelationship(targetID facechat.ID, rel facechat.RelationshipType) error {
	if err := tx.IsBlocked(targetID); err != nil {
		return err
	}

	_, err := tx.tx.Exec(
		"INSERT INTO relationships VALUES ($1, $2, $3) ON CONFLICT DO UPDATE SET type = $3",
		tx.UserID, targetID, rel,
	)
	return errors.Wrap(err, "failed to update relationship")
}

func (tx *Tx) CreatePublicLobby(name string, lvl facechat.SecretLevel) (*facechat.Room, error) {
	room := facechat.Room{
		Name:  name,
		Level: lvl,
	}

	if err := tx.createRoom(&room); err != nil {
		return nil, err
	}

	return &room, nil
}

func (tx *Tx) CreatePrivateRoom(targetUser facechat.ID) (*facechat.Room, error) {
	if err := tx.IsMutual(targetUser); err != nil {
		return nil, err
	}

	var u1, u2 = tx.UserID, targetUser
	// Guarantee a specific ID order.
	if u1 > u2 {
		u1, u2 = u2, u1
	}

	var roomID facechat.ID
	row := tx.tx.QueryRow(
		"SELECT room_id FROM private_rooms WHERE recipient1 = $1 AND recipient2 = $2",
		u1, u2,
	)

	if err := row.Scan(&roomID); err == nil {
		return tx.Room(roomID)
	}

	room := facechat.Room{
		Level: facechat.FullyOpen,
	}

	if err := tx.createRoom(&room); err != nil {
		return nil, err
	}

	_, err := tx.tx.Exec(
		"INSERT INTO private_rooms VALUES ($1, $2, $3)",
		room.ID, u1, u2,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to register private room")
	}

	return &room, nil
}

func (tx *Tx) createRoom(room *facechat.Room) error {
	// Generate a new ID.
	(*room).ID = facechat.GenerateID()

	_, err := tx.tx.Exec(
		"INSERT INTO rooms VALUES ($1, $2, $3, $4)",
		room.ID, room.Name, room.Topic, room.Level,
	)
	return errors.Wrap(err, "failed to create room")
}

func (tx *Tx) CreateMessage(room facechat.ID, content string) (*facechat.Message, error) {
	msg := facechat.Message{
		Type:     facechat.NormalMessage,
		RoomID:   room,
		AuthorID: tx.UserID,
		Markdown: content,
	}

	if err := tx.AddMessage(&msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

func (tx *Tx) AddMessage(msg *facechat.Message) error {
	if err := tx.IsInRoom(msg.RoomID); err != nil {
		return err
	}

	// do a room_participants check on AuthorID
	(*msg).ID = facechat.GenerateID()

	_, err := tx.tx.Exec(
		"INSERT INTO messages VALUES ($1, $2, $3, $4, $5)",
		msg.ID, msg.Type, msg.RoomID, msg.AuthorID, msg.Markdown,
	)
	return err
}

func (tx *Tx) JoinRoom(room facechat.ID) error {
	if err := tx.IsInRoom(room); err == nil {
		// Exit if the user is already in the room.
		return nil
	}

	_, err := tx.tx.Exec(
		"INSERT INTO room_participants VALUES ($1, $2)",
		room, tx.UserID,
	)
	if err != nil {
		return errors.Wrap(err, "failed to exec join room")
	}

	return nil
}

func (tx *Tx) LeaveRoom(room facechat.ID) error {
	if err := tx.IsInRoom(room); err != nil {
		// Exit if the user is not in the room.
		return err
	}

	_, err := tx.tx.Exec(
		"DELETE FROM room_participants WHERE room_id = $1 AND user_id = $2",
		room, tx.UserID,
	)
	if err != nil {
		return errors.Wrap(err, "failed to delete from room")
	}

	return nil
}

func randToken() (string, error) {
	var token = make([]byte, 32)

	if _, err := rand.Read(token); err != nil {
		return "", errors.Wrap(err, "failed to generate randomness")
	}

	return base64.RawURLEncoding.EncodeToString(token), nil
}
