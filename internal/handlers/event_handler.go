package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/frankbo/trash-app-api/internal/types"
)

func formatStartDate(date string) (time.Time, error) {
	value, err := time.Parse("20060102", date)

	if err != nil {
		return time.Time{}, fmt.Errorf("error while parsing string to time %g", err)
	}
	return value, nil
}

func marshalIcalEvents(icalEvents []*ics.VEvent) (types.Events, error) {
	events := []types.Event{}
	for _, value := range icalEvents {
		startDate, err := formatStartDate(value.GetProperty(ics.ComponentPropertyDtStart).Value)
		if err != nil {
			return types.Events{}, err
		}
		location := strings.TrimSpace(value.GetProperty(ics.ComponentPropertyLocation).Value)
		event := types.Event{
			Summary:     value.GetProperty(ics.ComponentPropertySummary).Value,
			Location:    location,
			StartDate:   startDate,
			Description: value.GetProperty(ics.ComponentPropertyDescription).Value,
		}

		events = append(events, event)
	}
	return types.Events{Events: events}, nil
}

func parseResponse(responseBody []byte) ([]*ics.VEvent, error) {
	cal, err := ics.ParseCalendar(strings.NewReader(string(responseBody)))
	if err != nil {
		return nil, fmt.Errorf("error while parsing response body %g", err)
	}
	return cal.Events(), nil
}

func createRequstUrl(locationId string, streetId string) string {
	baseUrl := "https://www.bad-berleburg.de/"
	location := "&ort=" + locationId
	street := "&strasse=" + streetId
	fullUrl := baseUrl + "output/abfall_export.php?csv_export=1&mode=vcal" + location + street + "&vtyp=4&vMo=1&vJ=2021&bMo=12"
	return fullUrl
}

func EventHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		locationId := r.FormValue("locationId")
		streetId := r.FormValue("streetId")

		if locationId == "" {
			http.Error(w, "Missing query parameter locationId", http.StatusBadRequest)
			return
		}

		if streetId == "" {
			streetId = locationId
		}

		url := createRequstUrl(locationId, streetId)
		resp, err := http.Get(url)

		if err != nil {
			http.Error(w, fmt.Sprintf("Error while fetching from url: %s", err), http.StatusBadRequest)
			return
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		icalEvents, err := parseResponse(body)

		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		events, err := marshalIcalEvents(icalEvents)

		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		events.SortEventsByStartDate()

		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(events); err != nil {
			println("encoding of the events failed")
		}
	})
}
