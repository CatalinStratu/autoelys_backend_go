package models

import "time"

type Automobile struct {
	ID           uint64     `json:"id"`
	URLHash      string     `json:"url_hash"`
	URL          string     `json:"url"`
	BrandID      uint64     `json:"brand_id"`
	Name         string     `json:"name"`
	Description  *string    `json:"description,omitempty"`
	PressRelease *string    `json:"press_release,omitempty"`
	Photos       *string    `json:"photos,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`

	// Relationship
	Brand *Brand `json:"brand,omitempty"`
}
