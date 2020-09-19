package db

import (
	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type ReadTx struct {
	tx *sqlx.Tx
}

func (tx *ReadTx) Session(token string) (*facechat.Session, error) {
	// TODO: check session
	// TODO: check expiry
	// TODO: upgrade to writeable tx
	return nil, errors.New("unimplemented")
}

func (tx *ReadTx) Relationship(targetUser facechat.ID) (facechat.RelationshipType, error) {
	panic("Implement me")
}

func (tx *ReadTx) User(id facechat.ID) (*facechat.User, error) {
	var user facechat.User

	err := tx.tx.
		QueryRowx("SELECT * FROM users WHERE id = ? LIMIT 1", id).
		StructScan(&user)
	if err != nil {
		return nil, errors.Wrap(err, "error getting user")
	}

	return &user, nil
}

func (tx *ReadTx) UserAccountsLen(id facechat.ID) (n int, err error) {
	err = tx.tx.
		QueryRowx("SELECT COUNT(*) FROM accounts WHERE user_id = ?", id).
		Scan(&n)

	if err != nil {
		return 0, errors.Wrap(err, "error getting accounts number")
	}

	return
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

func (tx *ReadTx) Messages(userID, roomID, beforeID facechat.ID, limit int) ([]facechat.Message, error) {
	// https://www.the-art-of-web.com/sql/select-before-after/
	// ORDER BY id
	// check if userID is in room_participants first
	panic("Implement me")
}

func (tx *ReadTx) SearchRoom(query string) ([]facechat.Room, error) {
	// LIKE topic = ? || $1 || ? (not sure if this is the right syntax)
	// LIKE name  = ? || $1 || ?
	// TYPE == facechat.PublicLobby (MUST DO THIS)
	panic("Implement me")
}

func (tx *ReadTx) Room(roomID facechat.ID) (*facechat.Room, error) {
	// TYPE == facechat.PublicLobby (MUST DO THIS)
	panic("Implement me")
}

func (tx *ReadTx) RoomParticipants(userID, roomID facechat.ID) ([]facechat.ID, error) {
	// TYPE == facechat.PublicLobby (MUST DO THIS)
	// room_participants
	// check if userID is in room_participants first
	panic("Implement me")
}
