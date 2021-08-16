package handlers

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	service "github.com/frankbo/trash-app-api/internal/services"
)

func HandleLambdaRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	locationId := request.QueryStringParameters["locationId"]
	streetId := request.QueryStringParameters["streetId"]

	calEvents, err := service.FetchEvents(locationId, streetId)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	calEventsMarshaled, err := json.Marshal(calEvents)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode:        200,
		IsBase64Encoded:   false,
		Headers:           map[string]string{"Content-Type": "application/json"},
		Body:              string(calEventsMarshaled),
		MultiValueHeaders: nil}, nil

}
