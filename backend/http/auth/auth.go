package auth

import (
	"context"
	"net/http"

	"github.com/diamondburned/facechat/backend/db"
	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/diamondburned/facechat/backend/http/tx"
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

func require(verifyAccounts bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("token")
			if err != nil {
				httperr.WriteErr(w, ErrTokenNotFound)
				return
			}

			var s *facechat.Session
			err = tx.RAcquire(r, func(tx *db.ReadTx) (err error) {
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
		SameSite: http.SameSiteStrictMode,
	})
}

func Session(r *http.Request) *facechat.Account {
	return SessionFromCtx(r.Context())
}

func SessionFromCtx(ctx context.Context) *facechat.Account {
	sm, ok := ctx.Value(sessionKey).(*facechat.Account)
	if !ok {
		return nil
	}

	return sm
}
