package register

import (
	"net/http"

	"github.com/diamondburned/facechat/backend/db"
	"github.com/diamondburned/facechat/backend/http/tx"
	"github.com/diamondburned/facechat/backend/internal/httperr"
	"github.com/go-chi/chi"
)

func Mount(db *db.Database) http.Handler {
	mux := chi.NewMux()
	mux.Post("/", register)

	return mux
}

func register(w http.ResponseWriter, r *http.Request) {
	err := tx.Acquire(r, func(tx *db.Transaction) error {
		tx
	})

	if err != nil {
		httperr.WriteErr(w, err)
		return
	}
}
