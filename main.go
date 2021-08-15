package main

import (
	"log"
	"net/http"

	"github.com/frankbo/trash-app-api/internal/handlers"
)

func main() {

	http.Handle("/events", handlers.EventHandler())

	s := &http.Server{
		Addr: ":8080",
	}
	log.Fatal(s.ListenAndServe())
}
