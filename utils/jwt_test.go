package utils

import (
	"fit-flow-api/config"
	"testing"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	config.AppConfig = &config.Config{
		JWTSecret:  "test-secret",
		JWTExpires: "1h",
	}

	userID := uuid.New()
	email := "test@example.com"

	token, err := GenerateJWT(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	claims, err := ValidateJWT(token)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %v, got %v", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("Expected email %s, got %s", email, claims.Email)
	}
}

func TestInvalidJWT(t *testing.T) {
	config.AppConfig = &config.Config{
		JWTSecret:  "test-secret",
		JWTExpires: "1h",
	}

	_, err := ValidateJWT("invalid-token")
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}
