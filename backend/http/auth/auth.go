package auth

import (
	"context"
	"net/http"

	"github.com/diamondburned/facechat/backend/db"
	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/diamondburned/facechat/backend/http/tx/dbctx"
	"github.com/diamondburned/facechat/backend/internal/httperr"
)

type ctxKey int

const (
	sessionKey ctxKey = iota
)

var ErrTokenNotFound = httperr.New(403, "token not found")

func Require() func(http.Handler) http.Handler {
	return require(true)
}

func RequireUnverified() func(http.Handler) http.Handler {
	return require(false)
}

func require(verifyAccounts bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Prevent the same session from being obtained twice.
			if Session(r) != nil {
				return
			}

			c, err := r.Cookie("token")
			if err != nil {
				httperr.WriteErr(w, ErrTokenNotFound)
				return
			}

			// Get the database and manually acquire a transaction without
			// package tx to avoid cyclical dependencies.
			d := dbctx.Database(r)

			var s *facechat.Session

			err = d.RAcquire(r.Context(), 0, func(tx *db.ReadTx) (err error) {
				s, err = tx.Session(c.Value)
				if err != nil {
					return err
				}

				if verifyAccounts {
					n, err := tx.UserAccountsLen(s.UserID)
					if err != nil {
						return err
					}

					// TODO: change
					if n < facechat.MinAccounts {
						return facechat.ErrNoAccountsLinked
					}
				}

				return nil
			})

			if err != nil {
				// TODO: session not found.
				httperr.WriteErr(w, err)
				return
			}

			next.ServeHTTP(w, r.WithContext(
				context.WithValue(r.Context(), sessionKey, s),
			))
		})
	}
}

func WriteSession(w http.ResponseWriter, s *facechat.Session) {
	if s == nil {
		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: "",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    s.Token,
		Path:     "/",
		Expires:  s.Expiry,
		SameSite: http.SameSiteLaxMode, // TODO: MAJOR CSRF RISK!!!!!!! Required for OAuth.
	})
}

func Session(r *http.Request) *facechat.Session {
	return SessionFromCtx(r.Context())
}

func SessionFromCtx(ctx context.Context) *facechat.Session {
	sm, ok := ctx.Value(sessionKey).(*facechat.Session)
	if !ok {
		return nil
	}

	return sm
}
