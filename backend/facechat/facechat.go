package facechat

import (
	"encoding/json"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/diamondburned/facechat/backend/internal/httperr"
	"github.com/pkg/errors"
)

var node *snowflake.Node

func init() {
	snowflake.Epoch = 1288834974657
	var err error
	node, err = snowflake.NewNode(0)
	if err != nil {
		panic(errors.Wrap(err, "error creating snowflake node"))
	}
}

type ID uint64

func GenerateID() ID {
	return ID(node.Generate())
}

type User struct {
	ID    ID     `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
}

const MinAccounts = 1

var ErrNoAccountsLinked = httperr.New(400, "no accounts linked")

type Account struct {
	Service string          `json:"service" db:"service"`
	Name    string          `json:"name"    db:"name"`
	URL     string          `json:"url"     db:"url"`
	Data    json.RawMessage `json:"data"    db:"data"`
	UserID  ID              `json:"userID"  db:"user_id"`
}

const SessionTimeout = 7 * 24 * time.Hour

type Session struct {
	UserID ID        `json:"-" db:"user_id"`
	Token  string    `json:"-" db:"token"`
	Expiry time.Time `json:"-" db:"expiry"`
}
