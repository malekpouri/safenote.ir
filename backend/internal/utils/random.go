package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateShortID returns a cryptographically random base62 string.
func GenerateShortID(length int) string {
	// 248 is the largest multiple of 62 that fits in a byte; rejecting bytes
	// above it keeps the distribution uniform.
	const maxByte = 255 - (256 % len(charset))
	out := make([]byte, 0, length)
	buf := make([]byte, length*2)
	for len(out) < length {
		if _, err := rand.Read(buf); err != nil {
			panic("crypto/rand unavailable: " + err.Error())
		}
		for _, b := range buf {
			if int(b) > maxByte {
				continue
			}
			out = append(out, charset[int(b)%len(charset)])
			if len(out) == length {
				break
			}
		}
	}
	return string(out)
}

// HashToken returns the hex SHA-256 of a token.
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// ConstantTimeEqual compares two strings without leaking timing information.
func ConstantTimeEqual(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
