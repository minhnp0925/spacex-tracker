package models

import "time"

type Launch struct {
	Id string `json:"id"`
	Name string `json:"name"`
	DateUTC time.Time `json:"date_utc"`
	Success *bool `json:"success,omitempty"` // nullable
	Upcoming bool `json:"upcoming"`
	Details string `json:"details,omitempty"`
}