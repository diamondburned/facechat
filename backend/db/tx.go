package db

import (
	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/jmoiron/sqlx"
)

type Tx struct {
	*ReadTx
}

func (tx *Tx) Register(username, password, email string) (*facechat.User, error) {}

type ReadTx struct {
	tx *sqlx.Tx
}

func (tx *ReadTx) User(id facechat.ID) (*facechat.User, error) {}

func (tx *ReadTx) UserAccounts(id facechat.ID) ([]facechat.Account, error) {}

func (tx *ReadTx) UserVerifyPassword(id facechat.ID) ([]byte, error) {}
