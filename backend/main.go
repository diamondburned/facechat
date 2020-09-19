package main

import (
	"log"
	"os"
	"path/filepath"

	nethttp "net/http"

	"github.com/diamondburned/facechat/backend/db"
	"github.com/diamondburned/facechat/backend/http"
	"github.com/diamondburned/facechat/backend/http/addr"
	"github.com/diamondburned/facechat/backend/http/routes/gateway/pubsub"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// Parse .env asap.
var _ = func() struct{} {
	d, err := filepath.Glob("env*")
	if err != nil {
		log.Fatalln("Failed to get env* files:", err)
	}

	for _, f := range d {
		if err := godotenv.Load(f); err != nil {
			log.Fatalf("Failed to load %q: %v\n", f, err)
		}
	}

	return struct{}{}
}()

func main() {
	d, err := db.Open(os.Getenv("SQL_ADDRESS"))
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
