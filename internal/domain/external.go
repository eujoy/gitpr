package domain

import (
	"fmt"
	"strings"
	"time"
)

const (
	inputTimeLayout = "2006-01-02"
	outputTimeLayout = "2006-01-02"
)

var nilTime = (time.Time{}).UnixNano()

// CustomTime describes the time input and output formats that are used across the system.
type CustomTime struct {
	time.Time
}

// IsSet checks if the time is actually set.
func (ct *CustomTime) IsSet() bool {
	return ct.UnixNano() != nilTime
}

// UnmarshalJSON for the custom time type.
func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(inputTimeLayout, s)
	return
}

// MarshalJSON for the custom time type.
func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.UnixNano() == nilTime {
		return []byte("null"), nil
	}

	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(outputTimeLayout))), nil
}

// SprintSummary describe the sprint summary level details.
type SprintSummary struct {
	Number    int        `json:"number"`
	Name      string     `json:"name"`
	StartDate CustomTime `json:"start_date"`
	EndDate   CustomTime `json:"end_date"`
}
