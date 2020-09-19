package form

import (
	"net/http"

	"github.com/diamondburned/facechat/backend/internal/httperr"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func AlwaysParse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			httperr.WriteErr(w, httperr.Wrap(err, 400, "failed to parse form"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Unmarshal decodes the form in the given request into the interface.
func Unmarshal(r *http.Request, v interface{}) error {
	switch r.Method {
	case http.MethodPatch, http.MethodPost, http.MethodPut:
		return decoder.Decode(v, r.PostForm)
	default:
		return decoder.Decode(v, r.Form)
	}
}
