package message

import (
	"encoding/json"
	"net/http"

	"github.com/diamondburned/facechat/backend/db"
	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/diamondburned/facechat/backend/http/auth"
	"github.com/diamondburned/facechat/backend/http/routes/gateway/pubsub"
	"github.com/diamondburned/facechat/backend/http/routes/room/roomid"
	"github.com/diamondburned/facechat/backend/http/tx"
	"github.com/diamondburned/facechat/backend/internal/form"
	"github.com/diamondburned/facechat/backend/internal/httperr"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

func Mount() http.Handler {
	r := chi.NewMux()
	r.Use(form.AlwaysParse)
	r.Use(auth.Require())

	r.Get("/", listMessages)
	r.Post("/", createMessage)

	return r
}

type ListMessagesParam struct {
	Before facechat.ID `schema:"before"`
	Limit  int         `schema:"limit"`
}

var ErrLimitInvalid = httperr.New(400, "invalid 0 <= limit <= 100")

func listMessages(w http.ResponseWriter, r *http.Request) {
	var p ListMessagesParam
	if err := form.Unmarshal(r, &p); err != nil {
		httperr.WriteErr(w, httperr.Wrap(err, 400, "failed to unmarshal form"))
		return
	}

	if p.Limit < 0 || p.Limit > 100 {
		httperr.WriteErr(w, ErrLimitInvalid)
		return
	}

	roomID := roomid.Get(r)
	s := auth.Session(r)

	var msgs []facechat.Message

	err := tx.RAcquire(r, func(tx *db.ReadTx) (err error) {
		msgs, err = tx.Messages(s.UserID, roomID, p.Before, p.Limit)
		if err != nil {
			err = errors.Wrap(err, "failed to get messages")
		}
		return
	})

	if err != nil {
		httperr.WriteErr(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(msgs); err != nil {
		httperr.WriteErr(w, errors.Wrap(err, "failed to encode messages"))
	}
}

type CreateMessageJSON struct {
	Markdown string `json:"markdown"`
}

func createMessage(w http.ResponseWriter, r *http.Request) {
	var cj CreateMessageJSON
	if err := json.NewDecoder(r.Body).Decode(&cj); err != nil {
		httperr.WriteErr(w, httperr.Wrap(err, 400, "failed to decode create msg json"))
		return
	}

	roomID := roomid.Get(r)
	s := auth.Session(r)

	var msg *facechat.Message

	err := tx.Acquire(r, func(tx *db.Tx) (err error) {
		msg, err = tx.CreateMessage(roomID, s.UserID, cj.Markdown)
		if err != nil {
			err = errors.Wrap(err, "failed to create message")
		}
		return
	})

	if err != nil {
		httperr.WriteErr(w, err)
		return
	}

	// Attempt to broadcast.
	coll := pubsub.RequestCollection(r)
	coll.BroadcastMessage(*msg)

	if err := json.NewEncoder(w).Encode(msg); err != nil {
		httperr.WriteErr(w, errors.Wrap(err, "failed to encode messages"))
	}
}
