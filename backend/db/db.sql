CREATE TABLE users (
	id    BIGINT PRIMARY KEY, -- const always
	pass  BINARY(40) NOT NULL,
	name  TEXT NOT NULL,
	email TEXT NOT NULL
);

-- name: UserPassword :one
SELECT pass FROM users WHERE id = $1 LIMIT 1;

-- name: User :one
SELECT (id, name, email) FROM users WHERE id = $1 LIMIT 1;

CREATE TABLE accounts (
	id      TEXT PRIMARY KEY, -- slow
	service TEXT NOT NULL,
	user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
);

-- name: UserAccounts :many
SELECT (id, service) FROM accounts WHERE user_id = $1;
