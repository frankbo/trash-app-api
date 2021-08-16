package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/frankbo/trash-app-api/internal/handlers"
)

func main() {
	lambda.Start(handlers.HandleLambdaRequest)

	// TODO with params distinguish between lambda and service.
	// http.Handle("/events", handlers.EventHandler())

	// s := &http.Server{
	// 	Addr: address,
	// }
	// log.Printf("Listen on %s...", address)
	// log.Fatal(s.ListenAndServe())
}
