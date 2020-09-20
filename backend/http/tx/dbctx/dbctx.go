package dbctx

import (
	"context"
	"net/http"

	"github.com/diamondburned/facechat/backend/db"
)

type ctxKey int

const (
	dbKey ctxKey = iota
)

func Middleware(db *db.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(
				context.WithValue(r.Context(), dbKey, db),
			))
		})
	}
}

func Database(r *http.Request) *db.DB {
	return DatabaseFromCtx(r.Context())
}

func DatabaseFromCtx(ctx context.Context) *db.DB {
	db, ok := ctx.Value(dbKey).(*db.DB)
	if !ok {
		return nil
	}
	return db
}
