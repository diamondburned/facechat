package oauth

import (
	"net/http"

	"github.com/diamondburned/facechat/backend/db"
	"github.com/diamondburned/facechat/backend/http/routes/oauth/github"
	"github.com/diamondburned/facechat/backend/http/routes/oauth/twitter"
	"github.com/go-chi/chi"
)

func Mount(db *db.DB) http.Handler {
	r := chi.NewMux()
	r.Mount("/", github.Mount(db))
	r.Mount("/", twitter.Mount(db))
	return r
}
