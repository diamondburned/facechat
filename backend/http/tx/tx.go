package tx

import (
	"net/http"

	"github.com/diamondburned/facechat/backend/db"
	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/diamondburned/facechat/backend/http/auth"
	"github.com/diamondburned/facechat/backend/http/tx/dbctx"
	"github.com/diamondburned/facechat/backend/internal/httperr"
)

var ErrMissingDB = httperr.New(500, "missing database middleware")

func Acquire(r *http.Request, fn func(*db.Tx) error) error {
	var db = dbctx.Database(r)
	if db == nil {
		return ErrMissingDB
	}

	var selfID facechat.ID
	if s := auth.Session(r); s != nil {
		selfID = s.UserID
	}

	return db.Acquire(r.Context(), selfID, fn)
}

func RAcquire(r *http.Request, fn func(*db.ReadTx) error) error {
	var db = dbctx.Database(r)
	if db == nil {
		return ErrMissingDB
	}

	var selfID facechat.ID
	if s := auth.Session(r); s != nil {
		selfID = s.UserID
	}

	return db.RAcquire(r.Context(), selfID, fn)
}
