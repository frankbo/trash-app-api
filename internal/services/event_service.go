package service

import (
	"errors"
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
	fullUrl := baseUrl + "output/abfall_export.php?csv_export=1&mode=vcal" + location + street + "&vtyp=4&vMo=1&vJ=2022&bMo=12"
	return fullUrl
}

func FetchEvents(locationId string, streetId string) (types.Events, error) {
	if locationId == "" {
		return types.Events{}, errors.New("Missing query parameter locationId")
	}

	// Add special case for location Muesse. It is the only location besides Bad Berleburg that needs a separate streetId
	if locationId == "1746.24" {
		streetId = "1746.30.1"
	}

	if streetId == "" {
		streetId = locationId
	}

	url := createRequstUrl(locationId, streetId)
	resp, err := http.Get(url)

	if err != nil {
		return types.Events{}, fmt.Errorf("Error while fetching from url: %s", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.Events{}, fmt.Errorf("%s", err)
	}

	icalEvents, err := parseResponse(body)

	if err != nil {
		return types.Events{}, fmt.Errorf("%s", err)
	}

	events, err := marshalIcalEvents(icalEvents)

	if err != nil {
		return types.Events{}, fmt.Errorf("Error while fetching from url: %s", err)
	}

	events.SortEventsByStartDate()
	return events, nil
}
