package user

import (
	"net/http"

	"github.com/diamondburned/facechat/backend/db"
	"github.com/go-chi/chi"
)

func Mount(db *db.DB) http.Handler {
	r := chi.NewMux()
	r.Route("/{user}", func(r chi.Router) {

	})

	return r
}

func user(w http.ResponseWriter, r *http.Request) {
	var username = chi.URLParam(r, "user")

	switch username {
	case "@me":
	default:
	}
}
