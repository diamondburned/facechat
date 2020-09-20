package addr

import (
	"log"
	"net/url"

	"github.com/diamondburned/facechat/backend/internal/dotenv"
)

var http url.URL

func HTTP() url.URL {
	if http.String() == "" {
		loadURL()
	}

	return http
}

func init() {
	loadURL()
}

func loadURL() {
	addr := dotenv.Getenv("HTTP_ADDRESS")
	if addr == "" {
		log.Fatalln("Missing $HTTP_ADDRESS")
	}

	u, err := url.Parse(addr)
	if err != nil {
		log.Fatalln("Invalid $HTTP_ADDRESS:", err)
	}

	http = *u
	http.Path = "/api"
}
