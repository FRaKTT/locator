package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fraktt/locator/internal/app"
	"github.com/fraktt/locator/internal/cache"
	"github.com/fraktt/locator/internal/opensky"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	addr := ":" + port

	a := app.New(
		opensky.New(),
		cache.New(),
	)

	http.HandleFunc("/", a.GetHandler())
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
