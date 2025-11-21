package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateSecureToken generates a cryptographically secure random token
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateInvitationToken generates a 32-byte (64 character hex) invitation token
func GenerateInvitationToken() (string, error) {
	return GenerateSecureToken(32)
}
