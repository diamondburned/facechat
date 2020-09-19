package oauth

import (
	"net/http"

	"github.com/diamondburned/facechat/backend/http/auth"
	"github.com/diamondburned/facechat/backend/http/routes/oauth/github"
	"github.com/diamondburned/facechat/backend/http/routes/oauth/twitter"
	"github.com/go-chi/chi"
)

func Mount() http.Handler {
	r := chi.NewMux()
	r.Use(auth.Require())
	r.Mount("/github", github.Mount())
	r.Mount("/twitter", twitter.Mount())
	return r
}
