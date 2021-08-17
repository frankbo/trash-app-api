package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/frankbo/trash-app-api/internal/types"
)

func HandleHttpRequest(fetchEvents func(string, string) (types.Events, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		locationId := r.FormValue("locationId")
		streetId := r.FormValue("streetId")

		events, err := fetchEvents(locationId, streetId)

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
