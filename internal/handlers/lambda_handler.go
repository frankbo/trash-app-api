package handlers

import (
	"context"

	service "github.com/frankbo/trash-app-api/internal/services"
	"github.com/frankbo/trash-app-api/internal/types"
)

type TrashEvent struct {
	LocationId string `json:"locationId"`
	StreetId   string `json:"streetId"`
}

type LambdaResponse struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       types.Events      `json:"body"`
}

func HandleRequest(ctx context.Context, trashEvent TrashEvent) (LambdaResponse, error) {

	locationId := trashEvent.LocationId
	streetId := trashEvent.StreetId

	events, err := service.FetchEvents(locationId, streetId)

	if err != nil {
		return LambdaResponse{}, err
	}

	return LambdaResponse{StatusCode: 200, Headers: map[string]string{"Content-Type": "application/json"}, Body: events}, nil
}
