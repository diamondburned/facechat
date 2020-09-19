package room

import (
	"encoding/json"
	"net/http"

	"github.com/diamondburned/facechat/backend/db"
	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/diamondburned/facechat/backend/http/auth"
	"github.com/diamondburned/facechat/backend/http/routes/gateway/pubsub"
	"github.com/diamondburned/facechat/backend/http/routes/room/message"
	"github.com/diamondburned/facechat/backend/http/routes/room/roomid"
	"github.com/diamondburned/facechat/backend/http/tx"
	"github.com/diamondburned/facechat/backend/internal/form"
	"github.com/diamondburned/facechat/backend/internal/httperr"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

func Mount() http.Handler {
	r := chi.NewMux()
	r.Use(auth.Require())
	r.Get("/", searchRoom) // TODO: aggregate room w/ the most people
	r.Post("/", createPublicLobby)
	r.Route(roomid.Route, func(r chi.Router) {
		r.Mount("/messages", message.Mount())

		r.Get("/", queryRoom)
		r.Post("/", joinRoom)
		r.Delete("/", leaveRoom)
	})

	return r
}

type SearchRoomQuery struct {
	Query string `schema:"q,required"`
}

func searchRoom(w http.ResponseWriter, r *http.Request) {
	var q SearchRoomQuery
	if err := form.Unmarshal(r, &q); err != nil {
		httperr.WriteErr(w, httperr.Wrap(err, 400, "failed to unmarshal query"))
		return
	}

	var rooms []facechat.Room

	err := tx.RAcquire(r, func(tx *db.ReadTx) (err error) {
		rooms, err = tx.SearchRoom(q.Query)
		if err != nil {
			err = errors.Wrap(err, "failed to get room")
		}
		return
	})

	if err != nil {
		httperr.WriteErr(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(rooms); err != nil {
		httperr.WriteErr(w, errors.Wrap(err, "failed to encode room"))
	}
}

type CreatePublicLobbyBody struct {
	Name  string               `json:"name"`
	Level facechat.SecretLevel `json:"level"`
}

func createPublicLobby(w http.ResponseWriter, r *http.Request) {
	var cj CreatePublicLobbyBody
	if err := json.NewDecoder(r.Body).Decode(&cj); err != nil {
		httperr.WriteErr(w, httperr.Wrap(err, 400, "failed to decode body"))
		return
	}

	s := auth.Session(r)

	var room *facechat.Room
	err := tx.Acquire(r, func(tx *db.Tx) (err error) {
		room, err = tx.CreatePublicLobby(cj.Name, cj.Level)
		if err != nil {
			return errors.Wrap(err, "failed to get room")
		}

		if err := tx.JoinRoom(room.ID, s.UserID); err != nil {
			return errors.Wrap(err, "failed to join newly created room")
		}

		return nil
	})

	if err != nil {
		httperr.WriteErr(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(room); err != nil {
		httperr.WriteErr(w, errors.Wrap(err, "failed to encode room"))
	}
}

func queryRoom(w http.ResponseWriter, r *http.Request) {
	roomID := roomid.Get(r)

	var room *facechat.Room
	err := tx.RAcquire(r, func(tx *db.ReadTx) (err error) {
		room, err = tx.Room(roomID)
		if err != nil {
			err = errors.Wrap(err, "failed to get room")
		}
		return
	})

	if err != nil {
		httperr.WriteErr(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(room); err != nil {
		httperr.WriteErr(w, errors.Wrap(err, "failed to encode room"))
	}
}

func joinRoom(w http.ResponseWriter, r *http.Request) {
	roomID := roomid.Get(r)
	s := auth.Session(r)

	err := tx.Acquire(r, func(tx *db.Tx) (err error) {
		err = tx.JoinRoom(roomID, s.UserID)
		err = errors.Wrap(err, "failed to join room")
		return
	})

	if err != nil {
		httperr.WriteErr(w, err)
		return
	}

	coll := pubsub.RequestCollection(r)
	coll.SubscribeRoom(s.UserID, roomID) // attempt to subscribe
}

func leaveRoom(w http.ResponseWriter, r *http.Request) {
	roomID := roomid.Get(r)
	s := auth.Session(r)

	err := tx.Acquire(r, func(tx *db.Tx) (err error) {
		err = tx.LeaveRoom(roomID, s.UserID)
		err = errors.Wrap(err, "failed to leave room")
		return
	})

	if err != nil {
		httperr.WriteErr(w, err)
		return
	}

	coll := pubsub.RequestCollection(r)
	coll.UnsubscribeRoom(s.UserID, roomID) // attempt to subscribe
}
