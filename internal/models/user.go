package models

import (
	"time"
)

type User struct {
	ID              uint64     `json:"id"`
	UUID            string     `json:"uuid"`
	RoleID          uint64     `json:"role_id"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	Email           string     `json:"email"`
	Phone           *string    `json:"phone,omitempty"`
	PasswordHash    string     `json:"-"`
	Active          bool       `json:"active"`
	AcceptedTermsAt *time.Time `json:"accepted_terms_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type Role struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)
