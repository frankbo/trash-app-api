package main

import (
	"log"
	"net/http"

	"github.com/frankbo/trash-app-api/internal/handlers"
)

const address = ":8080"

func main() {

	http.Handle("/events", handlers.EventHandler())

	s := &http.Server{
		Addr: address,
	}
	log.Printf("Listen on %s...", address)
	log.Fatal(s.ListenAndServe())
}
