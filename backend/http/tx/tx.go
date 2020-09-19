package tx

import (
	"context"
	"net/http"

	"github.com/diamondburned/facechat/backend/db"
	"github.com/diamondburned/facechat/backend/internal/httperr"
)

type ctxKey int

const (
	dbKey ctxKey = iota
)

var ErrMissingDB = httperr.New(500, "missing database middleware")

func Middleware(db *db.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(
				context.WithValue(r.Context(), dbKey, db),
			))
		})
	}
}

func GetDB(ctx context.Context) *db.DB {
	db, ok := ctx.Value(dbKey).(*db.DB)
	if !ok {
		return nil
	}
	return db
}

func Acquire(r *http.Request, fn func(*db.Tx) error) error {
	var db = GetDB(r.Context())
	if db == nil {
		return ErrMissingDB
	}

	return db.Acquire(r.Context(), fn)
}

func RAcquire(r *http.Request, fn func(*db.ReadTx) error) error {
	var db = GetDB(r.Context())
	if db == nil {
		return ErrMissingDB
	}

	return db.RAcquire(r.Context(), fn)
}
