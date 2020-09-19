package db

import "github.com/jmoiron/sqlx"

type Database struct {
	*sqlx.DB
}

const schema = `
	CREATE TABLE users (
		id    BIGINT PRIMARY KEY, -- const always
		pass  BINARY(40) NOT NULL,
		name  TEXT NOT NULL,
		email TEXT NOT NULL
	);
	
	CREATE TABLE accounts (
		id      TEXT PRIMARY KEY, -- slow
		service TEXT NOT NULL,
		user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	);
`

func NewDatabase(db string) (*Database, error) {}

func (db *Database) Acquire(fn func(tx *Transaction) error) error {}

func (db *Database) RAcquire(fn func(tx *ReadTransaction) error) error {}
