package db

import (
	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/jmoiron/sqlx"
)

type Transaction struct {
	*sqlx.Tx
}

type ReadTransaction struct {
	*sqlx.Tx
}

func (tx *ReadTransaction) User(id facechat.ID) (*facechat.User, error) {}

func (tx *ReadTransaction) UserAccounts(id facechat.ID) ([]facechat.Account, error) {}

func (tx *ReadTransaction) UserVerifyPassword(id facechat.ID) ([]byte, error) {}
