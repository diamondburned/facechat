package db

import (
	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/jmoiron/sqlx"
)

type Tx struct {
	tx *sqlx.Tx
}

type ReadTx struct {
	tx *sqlx.Tx
}

func (tx *ReadTransaction) User(id facechat.ID) (*facechat.User, error) {}

func (tx *ReadTransaction) UserAccounts(id facechat.ID) ([]facechat.Account, error) {}

func (tx *ReadTransaction) UserVerifyPassword(id facechat.ID) ([]byte, error) {}
