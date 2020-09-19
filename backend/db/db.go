package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type DB struct {
	db *sqlx.DB
}

const schema = `
	CREATE TABLE IF NOT EXISTS users (
		id    BIGINT     PRIMARY KEY, -- const always
		pass  BINARY(40) NOT NULL,
		name  TEXT       NOT NULL,
		email TEXT       NOT NULL UNIQUE
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
		type  SMALLINT NOT NULL,
		name  TEXT     NOT NULL,
		topic TEXT     NOT NULL,
		level SMALLINT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS room_participants (
		room_id BIGINT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
		user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

		UNIQUE (room_id, user_id)
	);

	CREATE TABLE IF NOT EXISTS messages (
		id        BIGINT   PRIMARY KEY,
		type      SMALLINT NOT NULL,
		room_id   BIGINT   NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
		author_id BIGINT   NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
		markdown  TEXT     NOT NULL
	);
`

func Open(source string) (*DB, error) {
	sqldb, err := sql.Open("pgx", source)
	if err != nil {
		return nil, err
	}
	db := sqlx.NewDb(sqldb, "pgx")
	return &DB{db}, nil
}

func (db *DB) Acquire(ctx context.Context, fn func(tx *Tx) error) error {
	sqlxtx, err := db.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "error starting transaction")
	}
	tx := &Tx{&ReadTx{sqlxtx}}
	err = fn(tx)
	return err
}

func (db *DB) RAcquire(ctx context.Context, fn func(tx *ReadTx) error) error {
	sqlxtx, err := db.db.BeginTxx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return errors.Wrap(err, "error starting transaction")
	}
	tx := &ReadTx{sqlxtx}
	err = fn(tx)
	return err
}
