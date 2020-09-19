package register

import (
	"encoding/json"
	"net/http"

	"github.com/diamondburned/facechat/backend/db"
	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/diamondburned/facechat/backend/http/auth"
	"github.com/diamondburned/facechat/backend/http/tx"
	"github.com/diamondburned/facechat/backend/internal/httperr"
	"github.com/go-chi/chi"
)

func Mount(db *db.DB) http.Handler {
	mux := chi.NewMux()
	mux.Post("/", register)

	return mux
}

type RegisterBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func register(w http.ResponseWriter, r *http.Request) {
	var body RegisterBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httperr.WriteErr(w, httperr.Wrap(err, 400, "Failed to decode JSON"))
		return
	}

	var u *facechat.User
	var s *facechat.Session

	err := tx.Acquire(r, func(tx *db.Tx) (err error) {
		u, s, err = tx.Register(body.Username, body.Password, body.Email)
		return
	})

	if err != nil {
		httperr.WriteErr(w, err)
		return
	}

	auth.WriteSession(w, s)

	if err := json.NewEncoder(w).Encode(u); err != nil {
		httperr.WriteErr(w, err)
		return
	}
}
