package db

import (
	"context"
	"database/sql"

	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	_ "github.com/jackc/pgx/v4/stdlib"
)

const schema = `
	CREATE TABLE IF NOT EXISTS users (
		id    BIGINT      PRIMARY KEY, -- const always
		pass  VARCHAR(60) NOT NULL,
		name  TEXT        NOT NULL,
		email TEXT        NOT NULL UNIQUE
	);

	CREATE TABLE IF NOT EXISTS relationships (
		user_id   BIGINT   NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		target_id BIGINT   NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		type      SMALLINT NOT NULL,

		UNIQUE (user_id, target_id)
	);

	CREATE TABLE IF NOT EXISTS accounts (
		service TEXT   NOT NULL,
		name    TEXT   NOT NULL,
		url     TEXT   NOT NULL,
		data    JSON   NOT NULL,
		user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

		UNIQUE (service, user_id) -- 1 service per person
	);

	CREATE TABLE IF NOT EXISTS sessions (
		user_id BIGINT    NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token   TEXT      NOT NULL UNIQUE,
		expiry  TIMESTAMP NOT NULL,

		UNIQUE (user_id, token)
	);

	CREATE TABLE IF NOT EXISTS rooms (
		id    BIGINT   PRIMARY KEY,
		name  TEXT     NOT NULL,
		topic TEXT     NOT NULL,
		level SMALLINT NOT NULL,

		UNIQUE (name, level)
	);

	CREATE TABLE IF NOT EXISTS room_participants (
		room_id BIGINT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
		user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

		UNIQUE (room_id, user_id)
	);

	CREATE TABLE IF NOT EXISTS private_rooms (
		room_id    BIGINT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
		recipient1 BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		recipient2 BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

		UNIQUE (recipient1, recipient2)
	);

	CREATE TABLE IF NOT EXISTS messages (
		id        BIGINT   PRIMARY KEY,
		type      SMALLINT NOT NULL,
		room_id   BIGINT   NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
		author_id BIGINT   NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
		markdown  TEXT     NOT NULL
	);
`

type userState struct {
	UserID facechat.ID
}

type DB struct {
	db *sqlx.DB
}

func Open(source string) (*DB, error) {
	db, err := sqlx.Open("pgx", source)
	if err != nil {
		return nil, errors.Wrap(err, "error opening db")
	}
	_, err = db.Exec(schema)
	if err != nil {
		return nil, errors.Wrap(err, "error executing schema")
	}
	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) Acquire(ctx context.Context, user facechat.ID, fn func(tx *Tx) error) error {
	sqlxtx, err := db.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "error starting transaction")
	}
	defer sqlxtx.Rollback()

	if err = fn(&Tx{&ReadTx{userState{user}, sqlxtx}}); err != nil {
		return err
	}

	return sqlxtx.Commit()
}

var txRO = &sql.TxOptions{
	Isolation: sql.LevelReadCommitted,
	ReadOnly:  true,
}

func (db *DB) RAcquire(ctx context.Context, user facechat.ID, fn func(tx *ReadTx) error) error {
	sqlxtx, err := db.db.BeginTxx(ctx, txRO)
	if err != nil {
		return errors.Wrap(err, "error starting transaction")
	}

	if err = fn(&ReadTx{userState{user}, sqlxtx}); err != nil {
		return err
	}

	return sqlxtx.Commit()
}
