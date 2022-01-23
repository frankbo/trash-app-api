package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/frankbo/trash-app-api/internal/handlers"
	service "github.com/frankbo/trash-app-api/internal/services"
)

func main() {
	localFlag := flag.Bool("local", false, "Run service local or as lambda")
	flag.Parse()

	if *localFlag {
		address := ":8080"
		http.Handle("/events", handlers.HandleHttpRequest(service.FetchEvents))

		s := &http.Server{
			Addr: address,
		}
		log.Printf("Listen on %s...", address)
		log.Fatal(s.ListenAndServe())
	} else {
		lambda.Start(handlers.HandleLambdaRequest(service.FetchEvents))
		log.Printf("Run as Lambda")
	}
}
