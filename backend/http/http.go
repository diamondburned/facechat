package http

import (
	"net/http"

	"github.com/diamondburned/facechat/backend/db"
	"github.com/diamondburned/facechat/backend/http/routes/gateway"
	"github.com/diamondburned/facechat/backend/http/routes/gateway/pubsub"
	"github.com/diamondburned/facechat/backend/http/routes/login"
	"github.com/diamondburned/facechat/backend/http/routes/oauth"
	"github.com/diamondburned/facechat/backend/http/routes/register"
	"github.com/diamondburned/facechat/backend/http/routes/room"
	"github.com/diamondburned/facechat/backend/http/routes/user"
	"github.com/diamondburned/facechat/backend/http/tx"
	"github.com/go-chi/chi"
)

func Mount(db *db.DB, coll *pubsub.Collection) http.Handler {
	mux := chi.NewMux()
	mux.Use(tx.Middleware(db))
	mux.Use(pubsub.UseCollection(coll))

	mux.Mount("/register", register.Mount())
	mux.Mount("/login", login.Mount())
	mux.Mount("/gateway", gateway.Mount())
	mux.Mount("/user", user.Mount())
	mux.Mount("/room", room.Mount())
	mux.Mount("/oauth", oauth.Mount())

	return mux
}
