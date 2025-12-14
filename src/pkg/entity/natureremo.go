package entity

import (
	"time"
)

type Device struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	NewestEvents map[string]NewestEvent `json:"newest_events"`
}

type NewestEvent struct {
	Val       float64   `json:"val"`
	CreatedAt time.Time `json:"created_at"`
}
