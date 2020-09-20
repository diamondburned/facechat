package main

import (
	"log"
	"os"
	"os/exec"

	nethttp "net/http"

	"github.com/diamondburned/facechat/backend/db"
	"github.com/diamondburned/facechat/backend/dotenv"
	"github.com/diamondburned/facechat/backend/http"
	"github.com/diamondburned/facechat/backend/http/addr"
	"github.com/diamondburned/facechat/backend/http/routes/gateway/pubsub"
	"github.com/go-chi/chi"
)

func main() {
	sh("cd frontend && npm run build")

	d, err := db.Open(dotenv.Getenv("SQL_ADDRESS"))
	if err != nil {
		log.Fatalln("Failed to connect to PostgreSQL:", err)
	}

	c := pubsub.NewCollection(d)

	r := chi.NewMux()
	r.Mount("/api", http.Mount(d, c))
	r.Mount("/", nethttp.FileServer(nethttp.Dir("./frontend/dist/")))

	h := addr.HTTP()

	if err := nethttp.ListenAndServe(h.Host, r); err != nil {
		log.Fatalln("Failed to listen:", err)
	}
}

func sh(eval string) {
	cmd := exec.Command("sh", "-c", eval)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalln("sh failed:", err)
	}
}
