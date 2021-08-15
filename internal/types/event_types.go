package types

import (
	"sort"
	"time"
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

func (e Events) SortEventsByStartDate() {
	sort.Sort(SortByStartDate(e.Events))
}

// -----------------------

type SortByStartDate []Event

func (a SortByStartDate) Len() int           { return len(a) }
func (a SortByStartDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByStartDate) Less(i, j int) bool { return a[i].StartDate.Before(a[j].StartDate) }

// -----------------------
