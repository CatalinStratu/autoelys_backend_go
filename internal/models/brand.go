package models

import "time"

type Brand struct {
	ID        uint64     `json:"id"`
	URLHash   string     `json:"url_hash"`
	URL       string     `json:"url"`
	Name      string     `json:"name"`
	Logo      *string    `json:"logo,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
