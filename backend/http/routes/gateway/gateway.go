package gateway

import (
	"net/http"

	"github.com/diamondburned/facechat/backend/http/auth"
	"github.com/diamondburned/facechat/backend/http/routes/gateway/pubsub"
	"github.com/diamondburned/facechat/backend/internal/httperr"
	"github.com/go-chi/chi"
)

func Mount() http.Handler {
	r := chi.NewMux()
	r.Use(auth.Require())
	r.Get("/", upgrade)

	return r
}

func upgrade(w http.ResponseWriter, r *http.Request) {
	s := auth.Session(r)
	c := pubsub.NewClient(s.UserID)

	if err := c.Connect(w, r); err != nil {
		httperr.WriteErr(w, err)
		return
	}

	defer c.Close()

	coll := pubsub.RequestCollection(r)
	coll.Register(c)
	defer coll.Unregister(c.UserID)

	c.Start()
}
