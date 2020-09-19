package roomid

import (
	"net/http"
	"strconv"

	"github.com/diamondburned/facechat/backend/facechat"
	"github.com/go-chi/chi"
)

const Route = `/{roomid:\d+}`

func Get(r *http.Request) facechat.ID {
	v := chi.URLParam(r, "roomid")
	if v == "" {
		return 0
	}

	u, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0
	}

	return facechat.ID(u)
}
