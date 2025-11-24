package models

import (
	"time"
)

type Note struct {
	ID             string    `json:"id"`
	EncryptedData  string    `json:"encrypted_data"`
	ViewsRemaining int       `json:"views_remaining"`
	ExpiresAt      time.Time `json:"expires_at"`
	CreatedAt      time.Time `json:"created_at"`
}
