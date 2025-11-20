package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WeightLog stores historical weight entries for a user
type WeightLog struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	WeightKg  float64        `gorm:"type:decimal(5,2);not null" json:"weight_kg"`
	Notes     string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

// Request DTOs

type CreateWeightLogRequest struct {
	WeightKg float64 `json:"weight_kg" binding:"required,gt=20,lt=500"`
	Notes    string  `json:"notes" binding:"omitempty,max=500"`
}

type UpdateWeightLogRequest struct {
	WeightKg float64 `json:"weight_kg" binding:"omitempty,gt=20,lt=500"`
	Notes    string  `json:"notes" binding:"omitempty,max=500"`
}

// Response DTOs

type WeightLogResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	WeightKg  float64   `json:"weight_kg"`
	WeightLbs float64   `json:"weight_lbs"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WeightStatsResponse struct {
	TotalEntries int       `json:"total_entries"`
	LatestWeight float64   `json:"latest_weight_kg"`
	MinWeight    float64   `json:"min_weight_kg"`
	MaxWeight    float64   `json:"max_weight_kg"`
	AvgWeight    float64   `json:"avg_weight_kg"`
	StartWeight  float64   `json:"start_weight_kg"`
	WeightChange float64   `json:"weight_change_kg"`
	PeriodDays   int       `json:"period_days"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
}

// ToResponse converts the model to a response DTO with unit conversions
func (w *WeightLog) ToResponse() WeightLogResponse {
	return WeightLogResponse{
		ID:        w.ID,
		UserID:    w.UserID,
		WeightKg:  w.WeightKg,
		WeightLbs: w.WeightKg * 2.20462,
		Notes:     w.Notes,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}
}
