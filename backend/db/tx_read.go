package db

import (
	"database/sql"
	"strings"
	"time"

	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type ReadTx struct {
	userState
	tx *sqlx.Tx
}

func (tx *ReadTx) Session(token string) (*facechat.Session, error) {
	var s = facechat.Session{
		Token: token,
	}

	r := tx.tx.QueryRowx(
		"SELECT (user_id, expiry) FROM sessions WHERE token = $1",
		token,
	)
	if err := r.StructScan(&s); err != nil {
		// return nil, errors.Wrap(err, "failed to scan session")
		return nil, facechat.ErrUnknownSession
	}

	if s.Expiry.Before(time.Now()) {
		return nil, facechat.ErrUnknownSession
	}

	return &s, nil
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
	var t facechat.RelationshipType

	r := tx.tx.QueryRow(
		"SELECT type FROM relationships WHERE user_id = $1 AND target_id = $2",
		self, target,
	)

	if err := r.Scan(&t); err != nil {
		return 0, err
	}

	if t == facechat.Blocked {
		return 0, facechat.ErrUserNotFound
	}

	return t, nil
}

func (tx *ReadTx) User(id facechat.ID) (*facechat.User, error) {
	// TODO: to allow for anonymous peeking, this method should check if the
	// current user actually has a relationship with this user. For this to
	// work, there should be a RoomUser method, while this method should only
	// return a user if mutual.
	var user facechat.User

	err := tx.tx.
		QueryRowx("SELECT * FROM users WHERE id = $1 LIMIT 1", id).
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
		QueryRowx("SELECT COUNT(*) FROM accounts WHERE user_id = $1", id).
		Scan(&n)

	if err != nil {
		return 0, errors.Wrap(err, "error getting accounts number")
	}

	return
}

func (tx *ReadTx) UserAccounts(id facechat.ID) ([]facechat.Account, error) {
	var accounts []facechat.Account
	rows, err := tx.tx.Queryx("SELECT * FROM accounts WHERE user_id = $1", id)
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

// Messages returns a list of messages ordered by oldest first.
func (tx *ReadTx) Messages(roomID, beforeID facechat.ID, limit int) ([]facechat.Message, error) {
	if limit < 0 || limit > facechat.MaxMessagesQuery {
		return nil, facechat.ErrMessageLimitInvalid
	}

	if err := tx.IsInRoom(roomID); err != nil {
		return nil, err
	}

	q, err := tx.tx.Queryx(`
		SELECT * FROM messages WHERE room_id = $1 AND id > $2 LIMIT $3 ORDER BY id DESC`,
		roomID, beforeID, limit)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query")
	}

	defer q.Close()

	var msgs []facechat.Message

	for q.Next() {
		var m facechat.Message
		if err := q.StructScan(&m); err != nil {
			return nil, errors.Wrap(err, "failed to scan message")
		}

		msgs = append(msgs, m)
	}

	if err := q.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to read")
	}

	return msgs, nil
}

var percEscaper = strings.NewReplacer(
	"%", "\\%",
	"\\", "\\\\",
)

func (tx *ReadTx) SearchRoom(query string) ([]facechat.Room, error) {
	query = percEscaper.Replace(query)

	q, err := tx.tx.Queryx(`
		SELECT * FROM rooms
		WHERE  name = % || $1 || % OR TOPIC = % || $1 || % AND level < $2`,
		query, facechat.Private)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query")
	}

	defer q.Close()

	var rooms []facechat.Room

	for q.Next() {
		var r facechat.Room
		if err := q.StructScan(&r); err != nil {
			return nil, errors.Wrap(err, "failed to scan room")
		}

		rooms = append(rooms, r)
	}

	if err := q.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to read")
	}

	return rooms, nil
}

// Room returns a public room. It returns an ErrRoomNotFound if the room is
// private.
func (tx *ReadTx) Room(roomID facechat.ID) (*facechat.Room, error) {
	r, err := tx.room(roomID)
	if err != nil {
		return nil, err
	}

	if r.Level == facechat.Private {
		return nil, facechat.ErrRoomNotFound
	}

	return r, nil
}

func (tx *ReadTx) room(roomID facechat.ID) (*facechat.Room, error) {
	var room facechat.Room

	err := tx.tx.
		QueryRowx("SELECT * FROM rooms WHERE id = $1", roomID).
		StructScan(&room)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query room")
	}

	return &room, nil
}

// IsInRoom returns nil if the user is in the room.
func (tx *ReadTx) IsInRoom(roomID facechat.ID) error {
	row := tx.tx.QueryRowx(
		"SELECT COUNT(*) FROM room_participants WHERE room_id = $1 AND user_id = $2",
		roomID, tx.UserID,
	)

	var count int
	if err := row.Scan(&count); err != nil {
		return err
	}

	if count > 0 {
		return nil
	}
	return facechat.ErrNotInRoom
}

func (tx *ReadTx) JoinedRooms() ([]facechat.Room, error) {}

func (tx *ReadTx) PrivateRooms() ([]facechat.Room, error) {}

// func (tx *ReadTx) RoomParticipants(roomID facechat.ID) ([]facechat.ID, error) {
// 	if err := tx.IsInRoom(roomID); err != nil {
// 		return nil, err
// 	}

// 	panic("Implement me")
// }
