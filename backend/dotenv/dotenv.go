package dotenv

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
)

var loadOnce sync.Once

var _ = func() struct{} {
	Load()
	return struct{}{}
}()

func Load() {
	loadOnce.Do(load)
}

func load() {
	d, err := filepath.Glob("env*")
	if err != nil {
		log.Fatalln("Failed to get env* files:", err)
	}

	for _, f := range d {
		log.Println("Loading", f)

		if err := godotenv.Load(f); err != nil {
			log.Fatalf("Failed to load %q: %v\n", f, err)
		}
	}
}

func Getenv(key string) string {
	Load()
	return os.Getenv(key)
}
