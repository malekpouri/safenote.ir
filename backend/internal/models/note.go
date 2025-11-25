package models

import (
	"time"
)

type Note struct {
	ID                  string    `json:"id" gorm:"primaryKey"`
	EncryptedData       string    `json:"encrypted_data"`
	PasswordHash        string    `json:"password_hash"`
	IsPasswordProtected bool      `json:"is_password_protected"`
	ViewsRemaining      int       `json:"views_remaining"`
	ExpiresAt           time.Time `json:"expires_at"`
	CreatedAt           time.Time `json:"created_at"`
}
