package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/diamondburned/facechat/backend/db"
	"github.com/pkg/errors"
)

type ctxKey int

const (
	scsKey ctxKey = iota
)

func Middleware(db *db.DB) func(http.Handler) http.Handler {
	m := scs.New()
	m.Lifetime = 7 * 24 * time.Hour

	return m.LoadAndSave
}

func SessionManager(r *http.Request) *scs.SessionManager {
	return SessionManagerFromCtx(r.Context())
}

func SessionManagerFromCtx(ctx context.Context) *scs.SessionManager {
	sm, ok := ctx.Value(scsKey).(*scs.SessionManager)
	if !ok {
		return nil
	}

	return sm
}

type store struct {
	*db.DB
}

func newStore(db *db.DB) scs.Store {
	return &store{db}
}

func (s *store) Delete(token string) (err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	err := s.DB.Acquire(ctx, func(tx *db.Tx) error {
		return tx.DeleteSession(token)
	})

	return errors.Wrap(err, "Failed to delete token")
}

func (s *store) Find(token string) (b []byte, found bool, err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	err = s.DB.RAcquire(ctx, func(tx *db.ReadTx) error {
		s, err := tx.Session(token)
		if err != nil {
			return err
		}

		found = true
		b = []byte(s.Data)

		return nil
	})

	err = errors.Wrap(err, "Failed to find token")

	return
}

func (s *store) Commit(token string, b []byte, expiry time.Time) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	err := s.DB.Acquire(ctx, func(tx *db.Tx) error {
		tx.UpdateSession()
	})
}

// func store(db *db.DB) scs.Store {
// }
