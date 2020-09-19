package db

import (
	"database/sql"

	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type ReadTx struct {
	userState
	tx *sqlx.Tx
}

func (tx *ReadTx) Session(token string) (*facechat.Session, error) {
	// TODO: check session
	// TODO: check expiry
	// TODO: upgrade to writeable tx
	return nil, errors.New("unimplemented")
}

// IsBlocked returns true if target blocked self.
func (tx *ReadTx) IsBlocked(target facechat.ID) error {
	_, err := tx.Relationship(target, tx.UserID)
	return err
}

// IsMutual returns nil if self is a friend of target and target is also a
// friend of self.
func (tx *ReadTx) IsMutual(target facechat.ID) error {
	targetToSelf, err := tx.Relationship(target, tx.UserID)
	if err != nil {
		return err
	}
	if targetToSelf != facechat.Friend {
		return facechat.ErrNotMutual
	}

	selfToTarget, err := tx.Relationship(tx.UserID, target)
	if err != nil {
		return err
	}
	if selfToTarget != facechat.Friend {
		return facechat.ErrNotMutual
	}

	return nil
}

func (tx *ReadTx) Relationship(self, target facechat.ID) (facechat.RelationshipType, error) {
	r, err := tx.relationship(self, target)
	if err != nil {
		return 0, err
	}
	if r == facechat.Blocked {
		return 0, facechat.ErrUserNotFound
	}
	return r, nil
}

func (tx *ReadTx) relationship(self, target facechat.ID) (t facechat.RelationshipType, err error) {
	r := tx.tx.QueryRow(
		"SELECT type FROM relationships WHERE user_id = $1 AND target_id = $2",
		self, target,
	)

	if err := r.Scan(&t); err != nil {
		return 0, err
	}

	return
}

func (tx *ReadTx) User(id facechat.ID) (*facechat.User, error) {
	var user facechat.User

	err := tx.tx.
		QueryRowx("SELECT * FROM users WHERE id = ? LIMIT 1", id).
		StructScan(&user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, facechat.ErrUserNotFound
		}

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

func (tx *ReadTx) Messages(roomID, beforeID facechat.ID, limit int) ([]facechat.Message, error) {
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

// IsInRoom returns nil if the user is in the room.
func (tx *ReadTx) IsInRoom(roomID facechat.ID) error {
	panic("Implement me")
}

func (tx *ReadTx) RoomParticipants(roomID facechat.ID) ([]facechat.ID, error) {
	// TYPE == facechat.PublicLobby (MUST DO THIS)
	// room_participants
	// check if userID is in room_participants first
	panic("Implement me")
}
