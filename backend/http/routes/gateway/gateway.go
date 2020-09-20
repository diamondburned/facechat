package gateway

import (
	"net/http"

	"github.com/diamondburned/facechat/backend/db"
	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/diamondburned/facechat/backend/http/auth"
	"github.com/diamondburned/facechat/backend/http/routes/gateway/pubsub"
	"github.com/diamondburned/facechat/backend/http/tx"
	"github.com/diamondburned/facechat/backend/internal/httperr"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

func Mount() http.Handler {
	r := chi.NewMux()
	r.Use(auth.Require())
	r.Get("/", upgrade)

	return r
}

type Ready struct {
	Me           facechat.User          `json:"me"` // TODO: add pubsub mutate events
	PublicRooms  []facechat.Room        `json:"public_rooms"`
	PrivateRooms []facechat.PrivateRoom `json:"private_rooms"`
}

func upgrade(w http.ResponseWriter, r *http.Request) {
	s := auth.Session(r)
	c := pubsub.NewClient(s.UserID)

	if err := c.Connect(w, r); err != nil {
		httperr.WriteErr(w, err)
		return
	}

	defer c.Close()

	var ready Ready

	tx.RAcquire(r, func(tx *db.ReadTx) error {
		u, err := tx.User(tx.UserID)
		if err != nil {
			return errors.Wrap(err, "failed to get user")
		}
		ready.Me = *u

		r, err := tx.JoinedRooms()
		if err != nil {
			return errors.Wrap(err, "failed to get joined rooms")
		}
		ready.PublicRooms = r

		p, err := tx.PrivateRooms()
		if err != nil {
			return errors.Wrap(err, "failed to join private rooms")
		}
		ready.PrivateRooms = p

		return nil
	})

	coll := pubsub.RequestCollection(r)
	coll.Register(c)
	defer coll.Unregister(c.UserID)

	c.Start()
}
