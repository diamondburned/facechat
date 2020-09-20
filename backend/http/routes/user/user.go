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
		r.Get("/accounts", accounts)
	})

	return r
}

func userID(r *http.Request) (facechat.ID, error) {
	s := auth.Session(r)
	u := chi.URLParam(r, "user")

	if u == "@me" {
		return s.UserID, nil
	}

	i, err := strconv.ParseUint(u, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse ID")
	}

	return facechat.ID(i), nil
}

type UserResponse struct {
	*facechat.User
	Accounts []facechat.Account `json:"accounts"`
}

func user(w http.ResponseWriter, r *http.Request) {
	var u = UserResponse{}

	i, err := userID(r)
	if err != nil {
		httperr.WriteErr(w, err)
		return
	}

	err = tx.RAcquire(r, func(tx *db.ReadTx) (err error) {
		u.User, err = tx.User(i)
		if err != nil {
			return errors.Wrap(err, "failed to get user")
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

func accounts(w http.ResponseWriter, r *http.Request) {
	i, err := userID(r)
	if err != nil {
		httperr.WriteErr(w, err)
		return
	}

	var accounts []facechat.Account

	err = tx.RAcquire(r, func(tx *db.ReadTx) (err error) {
		accounts, err = tx.UserAccounts(i)
		if err != nil {
			return errors.Wrap(err, "failed to get user accounts")
		}
		return nil
	})

	if err != nil {
		httperr.WriteErr(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(accounts); err != nil {
		httperr.WriteErr(w, errors.Wrap(err, "failed to encode JSON"))
	}
}
