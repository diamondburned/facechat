package facechat

import (
	"encoding/json"
	"time"

	"github.com/bwmarrin/snowflake"
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

	Accounts []Account `json:"accounts" db:"accounts"`
}

type Account struct {
	ID      string          `json:"id" db:"id"`
	Service string          `json:"service" db:"service"`
	Cookies json.RawMessage `json:"cookies" db:"cookies"`

	UserID ID `json:"userID" db:"user_id"`
}

type Session struct {
	UserID ID              `json:"-" db:"user_id"`
	Data   json.RawMessage `json:"-" db:"data"`
	Token  string          `json:"-" db:"token"`
	Expiry time.Time       `json:"-" db:"expiry"`

	User *User `json:"-" db:"-"`
}
