package github

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/github"
	"github.com/diamondburned/facechat/backend/db"
	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/diamondburned/facechat/backend/http/addr"
	"github.com/diamondburned/facechat/backend/http/auth"
	"github.com/diamondburned/facechat/backend/http/tx"
	"github.com/diamondburned/facechat/backend/internal/httperr"
	"github.com/drexedam/gravatar"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"

	githubauth "golang.org/x/oauth2/github"
)

var cookieCfg = gologin.CookieConfig{
	Name:   "gologin-github",
	Path:   "/",
	MaxAge: 5 * 60, // 5 minutes
}

var cfg = oauth2.Config{}

func init() {
	cfg.ClientID = os.Getenv("GITHUB_CLIENT_ID")
	if cfg.ClientID == "" {
		log.Fatalln("Missing $GITHUB_CLIENT_ID")
	}

	cfg.ClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
	if cfg.ClientSecret == "" {
		log.Fatalln("Missing $GITHUB_CLIENT_SECRET")
	}

	// TODO: improve
	url := addr.HTTP()
	url.Path += "/oauth/github/callback"
	cfg.RedirectURL = url.String()

	cfg.Endpoint = githubauth.Endpoint
}

func Mount() http.Handler {
	r := chi.NewMux()
	r.Use(auth.Require())
	r.Mount("/login", github.StateHandler(cookieCfg, http.HandlerFunc(onError)))
	r.Mount("/callback", github.StateHandler(
		cookieCfg,
		github.CallbackHandler(&cfg, http.HandlerFunc(issueSession), http.HandlerFunc(onError)),
	))
	return r
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

func issueSession(w http.ResponseWriter, r *http.Request) {
	s := auth.Session(r)

	u, err := github.UserFromContext(r.Context())
	if err != nil {
		httperr.WriteErr(w, errors.Wrap(err, "failed to get GitHub user"))
		return
	}

	if u.ID == nil {
		httperr.WriteErr(w, ErrMissingID)
		return
	}

	if u.Login == nil {
		httperr.WriteErr(w, ErrMissingUsername)
		return
	}

	var data = Data{
		Username: *u.Login,
	}

	switch {
	case u.AvatarURL != nil:
		data.AvatarURL = *u.AvatarURL
	case u.GravatarID != nil:
		data.AvatarURL = gravatar.New(*u.GravatarID).Size(200).AvatarURL()
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		httperr.WriteErr(w, errors.Wrap(err, "Failed to marshal github data"))
		return
	}

	var name = data.Username
	if u.Name != nil {
		name = *u.Name
	}

	var url = fmt.Sprintf("https://github.com/%s", data.Username)
	if u.HTMLURL != nil {
		url = *u.HTMLURL
	}

	var account = facechat.Account{
		Service: "GitHub",
		Name:    name,
		URL:     url,
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

	if err := json.NewEncoder(w).Encode(account); err != nil {
		httperr.WriteErr(w, errors.Wrap(err, "failed to encode"))
	}

	// TODO: redirect back to profile?
	http.Redirect(w, r, "/user", http.StatusFound)
}

func onError(w http.ResponseWriter, r *http.Request) {
	err := gologin.ErrorFromContext(r.Context())
	if err == nil {
		err = ErrUnknown
	}

	httperr.WriteErr(w, err)
}
