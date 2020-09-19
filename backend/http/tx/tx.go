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

func Middleware(db *db.Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			next.ServeHTTP(w, r)
		})
	}
}

func getDB(ctx context.Context) *db.Database {
	db, ok := ctx.Value(dbKey).(*db.Database)
	if !ok {
		return nil
	}
	return db
}

func Acquire(r *http.Request, fn func(*db.Transaction) error) error {
	var db = getDB(r.Context())
	if db == nil {
		return ErrMissingDB
	}
}

func RAcquire(r *http.Request, fn func(*db.ReadTransaction) error) error {
	var db = getDB(r.Context())
	if db == nil {
		return ErrMissingDB
	}
}
