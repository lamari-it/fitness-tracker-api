package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"lamari-fit-api/config"
	"time"
)

const (
	// RefreshTokenBytes is the number of random bytes for the token
	RefreshTokenBytes = 32
)

// GenerateRefreshToken creates a secure random token and its hash
// Returns: raw token (to send to client), hash (to store in DB), error
func GenerateRefreshToken() (token string, hash string, err error) {
	// Generate random bytes
	bytes := make([]byte, RefreshTokenBytes)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}

	// Encode as base64 for the raw token
	token = base64.URLEncoding.EncodeToString(bytes)

	// Create SHA-256 hash for storage
	hash = HashRefreshToken(token)

	return token, hash, nil
}

// HashRefreshToken creates a SHA-256 hash of a refresh token
func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// GetRefreshTokenExpiration returns the expiration time for refresh tokens
func GetRefreshTokenExpiration() time.Time {
	duration, err := time.ParseDuration(config.AppConfig.RefreshTokenExpires)
	if err != nil {
		// Default to 7 days if parsing fails
		duration = 7 * 24 * time.Hour
	}
	return time.Now().Add(duration)
}

// GetAccessTokenExpiresIn returns the access token expiry in seconds
func GetAccessTokenExpiresIn() int64 {
	duration, err := time.ParseDuration(config.AppConfig.JWTExpires)
	if err != nil {
		// Default to 1 hour if parsing fails
		duration = time.Hour
	}
	return int64(duration.Seconds())
}
