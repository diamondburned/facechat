package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/diamondburned/facechat/backend/db"
	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/diamondburned/facechat/backend/http/auth"
	"github.com/diamondburned/facechat/backend/http/tx"
	"github.com/diamondburned/facechat/backend/internal/httperr"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

func Mount() http.Handler {
	r := chi.NewMux()
	r.Use(auth.Require())
	r.Route("/{user}", func(r chi.Router) {
		r.Get("/", user)
	})

	return r
}

type UserResponse struct {
	*facechat.User
	Accounts []facechat.Account `json:"accounts"`
}

func user(w http.ResponseWriter, r *http.Request) {
	var user = chi.URLParam(r, "user")

	var s = auth.Session(r)
	var u = UserResponse{}

	err := tx.RAcquire(r, func(tx *db.ReadTx) (err error) {
		switch user {
		case "@me":
			u.User, err = tx.User(s.UserID)
		default:
			i, err := strconv.ParseUint(user, 10, 64)
			if err != nil {
				return errors.Wrap(err, "failed to parse ID")
			}

			u.User, err = tx.User(facechat.ID(i))
		}

		if err != nil {
			err = errors.Wrap(err, "failed to get user")
			return
		}

		u.Accounts, err = tx.UserAccounts(u.User.ID)
		if err != nil {
			return errors.Wrap(err, "failed to get accounts")
		}
		return nil
	})

	if err != nil {
		return
	}

	if err := json.NewEncoder(w).Encode(u); err != nil {
		httperr.WriteErr(w, errors.Wrap(err, "failed to encode JSON"))
	}
}
