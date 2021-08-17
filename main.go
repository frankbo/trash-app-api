package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/frankbo/trash-app-api/internal/handlers"
	service "github.com/frankbo/trash-app-api/internal/services"
)

func main() {
	lambda.Start(handlers.HandleLambdaRequest(service.FetchEvents))

	// TODO with params distinguish between lambda and service.
	// http.Handle("/events", handlers.HandleHttpRequest(service.FetchEvents))

	// s := &http.Server{
	// 	Addr: address,
	// }
	// log.Printf("Listen on %s...", address)
	// log.Fatal(s.ListenAndServe())
}
