package main

import (
	"log"
	"net/http"
)

var version string = "0.0.0-local"

func main() {
	Logger.Info("Starting the server")

	app := newApp()

	server := http.Server{
		Addr:    ":8080",
		Handler: app,
	}

	log.Printf("HTTP server listening on %v", server.Addr)
	server.ListenAndServe()
}
