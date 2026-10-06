package models

import (
	"time"
)

type Note struct {
	ID                  string    `json:"id"`
	EncryptedData       string    `json:"encrypted_data"`
	IsPasswordProtected bool      `json:"is_password_protected"`
	ViewsRemaining      int       `json:"views_remaining"`
	ExpiresAt           time.Time `json:"expires_at"`
	CreatedAt           time.Time `json:"created_at"`

	// AccessTokenHash is SHA-256 (hex) of the client-derived access token.
	// Empty for legacy (v1) notes created before tokens existed.
	AccessTokenHash string `json:"-"`
	// KdfSalt is the random per-note PBKDF2 salt (base64url). It is not secret;
	// the viewer needs it to derive the keys. Empty for legacy notes.
	KdfSalt string `json:"-"`
	// PasswordHash is only set on legacy (v1) notes and is used solely to
	// authorize their deletion.
	PasswordHash string `json:"-"`
}

// IsLegacy reports whether the note predates access tokens (v1 format).
func (n *Note) IsLegacy() bool {
	return n.AccessTokenHash == ""
}
