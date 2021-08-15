package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
)

type Event struct {
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"startDate"`
	Location    string    `json:"location"`
}

type Events struct {
	Events []Event `json:"events"`
}

func (e Events) sortEventsByStartDate() {
	sort.Sort(SortByStartDate(e.Events))
}

// -----------------------

type SortByStartDate []Event

func (a SortByStartDate) Len() int           { return len(a) }
func (a SortByStartDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByStartDate) Less(i, j int) bool { return a[i].StartDate.Before(a[j].StartDate) }

// -----------------------

func formatStartDate(date string) time.Time {
	value, err := time.Parse("20060102", date)

	if err != nil {
		println(err)
	}
	return value
}

func marshalIcalEvents(icalEvents []*ics.VEvent) Events {
	events := []Event{}
	for _, value := range icalEvents {
		startDate := formatStartDate(value.GetProperty(ics.ComponentPropertyDtStart).Value)
		location := strings.TrimSpace(value.GetProperty(ics.ComponentPropertyLocation).Value)
		event := Event{
			Summary:     value.GetProperty(ics.ComponentPropertySummary).Value,
			Location:    location,
			StartDate:   startDate,
			Description: value.GetProperty(ics.ComponentPropertyDescription).Value,
		}

		events = append(events, event)
	}
	return Events{Events: events}
}

func parseResponse(responseBody []byte) []*ics.VEvent {
	cal, err := ics.ParseCalendar(strings.NewReader(string(responseBody)))
	if err != nil {
		println("error while parsing response body")
	}
	return cal.Events()
}

func createRequstUrl(locationId string, streetId string) string {
	baseUrl := "https://www.bad-berleburg.de/"
	location := "&ort=" + locationId
	street := "&strasse=" + streetId
	if streetId == "" {
		street = "&strasse=" + locationId
	}
	fullUrl := baseUrl + "output/abfall_export.php?csv_export=1&mode=vcal" + location + street + "&vtyp=4&vMo=1&vJ=2021&bMo=12"
	return fullUrl
}

func EventHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		locationId := r.FormValue("locationId")
		streetId := r.FormValue("streetId")

		url := createRequstUrl(locationId, streetId)
		resp, err := http.Get(url)
		if err != nil {
			println("could fetch from url", err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		icalEvents := parseResponse(body)
		events := marshalIcalEvents(icalEvents)

		events.sortEventsByStartDate()

		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(events); err != nil {
			println("encoding of the events failed")
		}
	})
}
