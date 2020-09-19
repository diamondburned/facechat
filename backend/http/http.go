package http

import (
	"net/http"

	"github.com/diamondburned/facechat/backend/db"
	"github.com/diamondburned/facechat/backend/http/routes/oauth"
	"github.com/diamondburned/facechat/backend/http/routes/register"
	"github.com/diamondburned/facechat/backend/http/routes/user"
	"github.com/diamondburned/facechat/backend/http/tx"
	"github.com/go-chi/chi"
)

func Mount(db *db.DB) http.Handler {
	mux := chi.NewMux()
	mux.Use(tx.Middleware(db))
	mux.Mount("/register", register.Mount(db))
	mux.Mount("/user", user.Mount(db))
	mux.Mount("/oauth", oauth.Mount(db))

	return mux
}
