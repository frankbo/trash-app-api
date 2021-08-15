package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

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

	

	return LambdaResponse{StatusCode: 200, Headers: map[string]string{"Content-Type": "application/json"}, Body: events}, nil
}
