package github

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danilopolani/gocialite"
	"github.com/diamondburned/facechat/backend/db"
	"github.com/diamondburned/facechat/backend/dotenv"
	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/diamondburned/facechat/backend/http/addr"
	"github.com/diamondburned/facechat/backend/http/auth"
	"github.com/diamondburned/facechat/backend/http/tx"
	"github.com/diamondburned/facechat/backend/internal/httperr"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

var (
	clientID     string
	clientSecret string
	redirectURL  string

	gocial = gocialite.NewDispatcher()
)

func init() {
	clientID = dotenv.Getenv("GITHUB_CLIENT_ID")
	clientSecret = dotenv.Getenv("GITHUB_CLIENT_SECRET")

	// TODO: improve
	url := addr.HTTP()
	url.Path += "/oauth/github/callback"
	redirectURL = url.String()
}

func Mount() http.Handler {
	r := chi.NewMux()
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/callback", callbackHandler)
	return r
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	gocialnew := gocial.New()
	// spew.Dump(gocialnew)

	a, err := gocialnew.Driver("github").Redirect(
		clientID, clientSecret, redirectURL,
	)
	if err != nil {
		httperr.WriteErr(w, errors.Wrap(err, "failed to create new oauth session"))
		return
	}

	http.Redirect(w, r, a, http.StatusFound)
}

type Data struct {
	Username  string `json:"username"`   // .Login
	AvatarURL string `json:"avatar_url"` // .AvatarURL
}

var (
	ErrUnknown         = httperr.New(500, "unknown error")
	ErrMissingID       = httperr.New(500, "missing ID")
	ErrMissingUsername = httperr.New(500, "missing username")
)

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	s := auth.Session(r)

	u, _, err := gocial.Handle(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		httperr.WriteErr(w, httperr.Wrap(err, 400, "failed to handle callback"))
		return
	}

	var data = Data{
		Username:  u.Username,
		AvatarURL: u.Avatar,
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		httperr.WriteErr(w, errors.Wrap(err, "Failed to marshal github data"))
		return
	}

	var name = u.Username
	if u.FullName != "<nil>" { // gocial is a trash library
		name = u.FullName
	}

	var account = facechat.Account{
		Service: "GitHub",
		Name:    name,
		URL:     fmt.Sprintf("https://github.com/%s", data.Username),
		Data:    dataJSON,
		UserID:  s.UserID,
	}

	err = tx.Acquire(r, func(tx *db.Tx) error {
		return errors.Wrap(tx.AddAccount(account), "failed to save")
	})

	if err != nil {
		httperr.WriteErr(w, err)
		return
	}

	// TODO: redirect back to profile?
	http.Redirect(w, r, "/user", http.StatusFound)
}
