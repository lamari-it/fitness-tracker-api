package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RPEScale represents a Rate of Perceived Exertion scale definition
type RPEScale struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	MinValue    int            `gorm:"not null;default:1" json:"min_value"`
	MaxValue    int            `gorm:"not null;default:10" json:"max_value"`
	IsGlobal    bool           `gorm:"not null;default:false" json:"is_global"`
	TrainerID   *uuid.UUID     `gorm:"type:uuid" json:"trainer_id,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Trainer *User           `gorm:"foreignKey:TrainerID;constraint:OnDelete:CASCADE" json:"trainer,omitempty"`
	Values  []RPEScaleValue `gorm:"foreignKey:ScaleID" json:"values,omitempty"`
}

// RPEScaleValue represents a single value within an RPE scale
type RPEScaleValue struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ScaleID     uuid.UUID `gorm:"type:uuid;not null" json:"scale_id"`
	Value       int       `gorm:"not null" json:"value"`
	Label       string    `gorm:"type:varchar(50);not null" json:"label"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`

	// Relationships
	Scale RPEScale `gorm:"foreignKey:ScaleID;constraint:OnDelete:CASCADE" json:"scale,omitempty"`
}

func (s *RPEScale) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}

func (v *RPEScaleValue) BeforeCreate(tx *gorm.DB) (err error) {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return
}

// Validate validates the RPE scale
func (s *RPEScale) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(s.Name) > 100 {
		return fmt.Errorf("name must be at most 100 characters")
	}
	if s.MinValue >= s.MaxValue {
		return fmt.Errorf("min_value must be less than max_value")
	}
	if s.MinValue < 0 {
		return fmt.Errorf("min_value must be non-negative")
	}
	if s.MaxValue > 100 {
		return fmt.Errorf("max_value must be at most 100")
	}
	// Global scales must not have a trainer ID
	if s.IsGlobal && s.TrainerID != nil {
		return fmt.Errorf("global scales cannot have a trainer_id")
	}
	// Non-global scales must have a trainer ID
	if !s.IsGlobal && s.TrainerID == nil {
		return fmt.Errorf("custom scales must have a trainer_id")
	}
	return nil
}

// Validate validates the RPE scale value
func (v *RPEScaleValue) Validate(scale *RPEScale) error {
	if v.Label == "" {
		return fmt.Errorf("label is required")
	}
	if len(v.Label) > 50 {
		return fmt.Errorf("label must be at most 50 characters")
	}
	if scale != nil {
		if v.Value < scale.MinValue || v.Value > scale.MaxValue {
			return fmt.Errorf("value must be between %d and %d", scale.MinValue, scale.MaxValue)
		}
	}
	return nil
}

// Request DTOs

type CreateRPEScaleRequest struct {
	Name        string                       `json:"name" binding:"required,max=100"`
	Description string                       `json:"description" binding:"omitempty,max=500"`
	MinValue    int                          `json:"min_value" binding:"omitempty,gte=0,lte=99"`
	MaxValue    int                          `json:"max_value" binding:"omitempty,gte=1,lte=100"`
	Values      []CreateRPEScaleValueRequest `json:"values" binding:"omitempty,dive"`
}

type CreateRPEScaleValueRequest struct {
	Value       int    `json:"value" binding:"required"`
	Label       string `json:"label" binding:"required,max=50"`
	Description string `json:"description" binding:"omitempty,max=500"`
}

type UpdateRPEScaleRequest struct {
	Name        string `json:"name" binding:"omitempty,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
}

type AddRPEScaleValueRequest struct {
	Value       int    `json:"value" binding:"required"`
	Label       string `json:"label" binding:"required,max=50"`
	Description string `json:"description" binding:"omitempty,max=500"`
}

// Response DTOs

type RPEScaleResponse struct {
	ID          uuid.UUID               `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	MinValue    int                     `json:"min_value"`
	MaxValue    int                     `json:"max_value"`
	IsGlobal    bool                    `json:"is_global"`
	TrainerID   *uuid.UUID              `json:"trainer_id,omitempty"`
	Values      []RPEScaleValueResponse `json:"values,omitempty"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
}

type RPEScaleValueResponse struct {
	ID          uuid.UUID `json:"id"`
	Value       int       `json:"value"`
	Label       string    `json:"label"`
	Description string    `json:"description"`
}

// ToResponse converts RPEScale to response format
func (s *RPEScale) ToResponse() RPEScaleResponse {
	values := make([]RPEScaleValueResponse, 0, len(s.Values))
	for _, v := range s.Values {
		values = append(values, v.ToResponse())
	}

	return RPEScaleResponse{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		MinValue:    s.MinValue,
		MaxValue:    s.MaxValue,
		IsGlobal:    s.IsGlobal,
		TrainerID:   s.TrainerID,
		Values:      values,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

// ToResponse converts RPEScaleValue to response format
func (v *RPEScaleValue) ToResponse() RPEScaleValueResponse {
	return RPEScaleValueResponse{
		ID:          v.ID,
		Value:       v.Value,
		Label:       v.Label,
		Description: v.Description,
	}
}
