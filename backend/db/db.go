package db

import "github.com/jmoiron/sqlx"

type DB struct {
	db *sqlx.DB
}

const schema = `
	CREATE TABLE IF NOT EXISTS users (
		id    BIGINT PRIMARY KEY, -- const always
		pass  BINARY(40) NOT NULL,
		name  TEXT NOT NULL,
		email TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS accounts (
		id      TEXT PRIMARY KEY, -- slow
		service TEXT NOT NULL,
		cookies JSON NOT NULL,
		user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	);
`

func Open(source string) (*Database, error) {
	sqldb, err := sql.Open("pgx", source)
	if err != nil {
		return nil, err
	}
	db := sqlx.NewDB(sqldb, "pgx")
	return &DB{db}, nil
}

func (db *DB) Acquire(ctx context.Context, fn func(tx *TX) error) error {
	sqltx, err := db.db.BeginTX(ctx, nil)
	if err != nil {
		return errors.Wrap("error starting transaction", err)
	}
	tx := Tx{sqltx}
	err = fn(tx)
	return err
}

func (db *DB) RAcquire(fn func(tx *ReadTX) error) error {
	sqltx, err := db.db.BeginTX(ctx, &sql.TxOptions{ReadOnly:true})
	if err != nil {
		return errors.Wrap("error starting transaction", err)
	}
	tx := ReadTx{sqltx}
	err = fn(tx)
	return err
}
