package models

import (
	"time"
)

type Service struct {
	ID              uint64    `json:"id"`
	UUID            string    `json:"uuid"`
	Title           string    `json:"title"`
	Description     *string   `json:"description,omitempty"`
	Price           float64   `json:"price"`
	Currency        string    `json:"currency"`
	DurationMinutes *uint     `json:"duration_minutes,omitempty"`
	Active          bool      `json:"active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
