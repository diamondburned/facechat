package main

import (
	"log"

	nethttp "net/http"

	"github.com/diamondburned/facechat/backend/db"
	"github.com/diamondburned/facechat/backend/http"
	"github.com/diamondburned/facechat/backend/http/addr"
	"github.com/diamondburned/facechat/backend/http/routes/gateway/pubsub"
	"github.com/diamondburned/facechat/backend/internal/dotenv"
	"github.com/go-chi/chi"
)

func main() {
	d, err := db.Open(dotenv.Getenv("SQL_ADDRESS"))
	if err != nil {
		log.Fatalln("Failed to connect to PostgreSQL:", err)
	}

	c := pubsub.NewCollection(d)

	r := chi.NewMux()
	r.Mount("/api", http.Mount(d, c))

	h := addr.HTTP()

	if err := nethttp.ListenAndServe(h.Host, r); err != nil {
		log.Fatalln("Failed to listen:", err)
	}
}
