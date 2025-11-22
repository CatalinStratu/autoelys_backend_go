package models

import "time"

type PasswordResetToken struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}
