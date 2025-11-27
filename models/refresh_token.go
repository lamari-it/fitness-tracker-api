package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RefreshToken represents a refresh token stored in the database
// for per-device session management
type RefreshToken struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash  string         `gorm:"type:varchar(64);not null;uniqueIndex" json:"-"`
	DeviceInfo string         `gorm:"type:varchar(255)" json:"device_info"`
	IPAddress  string         `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent  string         `gorm:"type:varchar(512)" json:"user_agent"`
	ExpiresAt  time.Time      `gorm:"not null;index" json:"expires_at"`
	LastUsedAt *time.Time     `json:"last_used_at,omitempty"`
	RevokedAt  *time.Time     `json:"revoked_at,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

// IsValid checks if the refresh token is valid (not expired and not revoked)
func (rt *RefreshToken) IsValid() bool {
	if rt.RevokedAt != nil {
		return false
	}
	if time.Now().After(rt.ExpiresAt) {
		return false
	}
	return true
}

// IsRevoked checks if the refresh token has been revoked
func (rt *RefreshToken) IsRevoked() bool {
	return rt.RevokedAt != nil
}

// IsExpired checks if the refresh token has expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// Revoke marks the refresh token as revoked
func (rt *RefreshToken) Revoke() {
	now := time.Now()
	rt.RevokedAt = &now
}

// UpdateLastUsed updates the last used timestamp
func (rt *RefreshToken) UpdateLastUsed() {
	now := time.Now()
	rt.LastUsedAt = &now
}

// SessionResponse represents a session for API responses
type SessionResponse struct {
	ID         uuid.UUID  `json:"id"`
	DeviceInfo string     `json:"device_info"`
	IPAddress  string     `json:"ip_address"`
	UserAgent  string     `json:"user_agent"`
	LastUsedAt *time.Time `json:"last_used_at"`
	CreatedAt  time.Time  `json:"created_at"`
	IsCurrent  bool       `json:"is_current"`
}

// ToSessionResponse converts a RefreshToken to a SessionResponse
func (rt *RefreshToken) ToSessionResponse(currentTokenHash string) SessionResponse {
	return SessionResponse{
		ID:         rt.ID,
		DeviceInfo: rt.DeviceInfo,
		IPAddress:  rt.IPAddress,
		UserAgent:  rt.UserAgent,
		LastUsedAt: rt.LastUsedAt,
		CreatedAt:  rt.CreatedAt,
		IsCurrent:  rt.TokenHash == currentTokenHash,
	}
}
