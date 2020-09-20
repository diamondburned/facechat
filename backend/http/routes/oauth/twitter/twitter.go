package twitter

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/twitter"
	"github.com/dghubble/oauth1"
	"github.com/diamondburned/facechat/backend/db"
	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/diamondburned/facechat/backend/http/addr"
	"github.com/diamondburned/facechat/backend/http/auth"
	"github.com/diamondburned/facechat/backend/http/tx"
	"github.com/diamondburned/facechat/backend/dotenv"
	"github.com/diamondburned/facechat/backend/internal/httperr"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	twitterauth "github.com/dghubble/oauth1/twitter"
)

var cfg = oauth1.Config{}

func init() {
	cfg.ConsumerKey = dotenv.Getenv("TWITTER_CONSUMER_KEY")
	if cfg.ConsumerKey == "" {
		log.Fatalln("Missing $TWITTER_CONSUMER_KEY")
	}

	cfg.ConsumerSecret = dotenv.Getenv("TWITTER_CONSUMER_SECRET")
	if cfg.ConsumerSecret == "" {
		log.Fatalln("Missing $TWITTER_CONSUMER_SECRET")
	}

	// TODO: improve
	url := addr.HTTP()
	url.Path += "/oauth/twitter/callback"
	cfg.CallbackURL = url.String()

	cfg.Endpoint = twitterauth.AuthorizeEndpoint
}

func Mount() http.Handler {
	r := chi.NewMux()
	r.Use(auth.Require())
	r.Mount("/login", twitter.LoginHandler(&cfg, http.HandlerFunc(onError)))
	r.Mount("/callback", twitter.CallbackHandler(
		&cfg,
		http.HandlerFunc(issueSession),
		http.HandlerFunc(onError),
	))
	return r
}

type Data struct {
	ID        int64  `json:"id"`
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

	u, err := twitter.UserFromContext(r.Context())
	if err != nil {
		httperr.WriteErr(w, errors.Wrap(err, "failed to get Twitter user"))
		return
	}

	var data = Data{
		ID:        u.ID,
		Username:  u.ScreenName,
		AvatarURL: u.ProfileImageURLHttps,
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		httperr.WriteErr(w, errors.Wrap(err, "failed to marshal Twitter data"))
		return
	}

	var name = data.Username
	if u.Name != "" {
		name = u.Name
	}

	var account = facechat.Account{
		Service: "Twitter",
		Name:    name,
		URL:     u.URL,
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
