package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func init() {
	d, err := filepath.Glob("env*")
	if err != nil {
		log.Fatalln("Failed to get env* files:", err)
	}

	if len(d) == 0 {
		log.Fatalln("No env files found.")
	}

	for _, f := range d {
		if err := godotenv.Load(f); err != nil {
			log.Fatalf("Failed to load %q: %v\n", f, err)
		}
	}
}

func main() {
	d, err := sql.Open("pgx", os.Getenv("SQL_ADDRESS"))
	if err != nil {
		log.Fatalln("Failed to connect to PostgreSQL:", err)
	}
}
