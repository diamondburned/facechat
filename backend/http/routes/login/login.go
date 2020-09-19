package login

import (
	"encoding/json"
	"net/http"

	"github.com/diamondburned/facechat/backend/db"
	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/diamondburned/facechat/backend/http/auth"
	"github.com/diamondburned/facechat/backend/http/tx"
	"github.com/diamondburned/facechat/backend/internal/httperr"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

func Mount() http.Handler {
	r := chi.NewMux()
	r.Post("/", login)

	return r
}

type LoginJSON struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func login(w http.ResponseWriter, r *http.Request) {
	var l LoginJSON
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		httperr.WriteErr(w, httperr.Wrap(err, 400, "failed to decode login body"))
		return
	}

	var s *facechat.Session

	err := tx.Acquire(r, func(tx *db.Tx) (err error) {
		s, err = tx.Login(l.Email, l.Password)
		if err != nil {
			err = errors.Wrap(err, "failed to login")
		}
		return
	})

	if err != nil {
		httperr.WriteErr(w, err)
		return
	}

	auth.WriteSession(w, s)
}
