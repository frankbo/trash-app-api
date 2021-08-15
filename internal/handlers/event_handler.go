package handlers

import (
	"encoding/json"
	"net/http"

	service "github.com/frankbo/trash-app-api/internal/services"
)

func EventHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		locationId := r.FormValue("locationId")
		streetId := r.FormValue("streetId")

		events, err := service.FetchEvents(locationId, streetId)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(events); err != nil {
			println("encoding of the events failed")
		}
	})
}
