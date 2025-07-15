package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Equipment represents gym equipment
type Equipment struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null;unique" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Category    string    `gorm:"type:varchar(50)" json:"category"` // machine, free_weight, cable, cardio, other
	ImageURL    string    `gorm:"type:text" json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	ExerciseLinks []ExerciseEquipment `gorm:"foreignKey:EquipmentID" json:"exercise_links,omitempty"`
}

// ExerciseEquipment represents the many-to-many relationship between exercises and equipment
type ExerciseEquipment struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ExerciseID  uuid.UUID `gorm:"type:uuid;not null;index" json:"exercise_id"`
	EquipmentID uuid.UUID `gorm:"type:uuid;not null;index" json:"equipment_id"`
	Optional    bool      `gorm:"default:false" json:"optional"` // Whether this equipment is optional for the exercise
	Notes       string    `gorm:"type:text" json:"notes"`        // Additional notes about using this equipment
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Exercise  Exercise  `gorm:"foreignKey:ExerciseID;constraint:OnDelete:CASCADE" json:"exercise,omitempty"`
	Equipment Equipment `gorm:"foreignKey:EquipmentID;constraint:OnDelete:CASCADE" json:"equipment,omitempty"`
}

// BeforeCreate hooks
func (e *Equipment) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return
}

func (ee *ExerciseEquipment) BeforeCreate(tx *gorm.DB) (err error) {
	if ee.ID == uuid.Nil {
		ee.ID = uuid.New()
	}
	return
}

// Response DTOs
type EquipmentResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ExerciseEquipmentResponse struct {
	ID          uuid.UUID         `json:"id"`
	ExerciseID  uuid.UUID         `json:"exercise_id"`
	EquipmentID uuid.UUID         `json:"equipment_id"`
	Optional    bool              `json:"optional"`
	Notes       string            `json:"notes"`
	Equipment   EquipmentResponse `json:"equipment"`
}

type EquipmentWithExercises struct {
	EquipmentResponse
	ExerciseCount int                         `json:"exercise_count"`
	Exercises     []ExerciseEquipmentResponse `json:"exercises,omitempty"`
}

// Helper methods
func (e *Equipment) ToResponse() EquipmentResponse {
	return EquipmentResponse{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		Category:    e.Category,
		ImageURL:    e.ImageURL,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func (ee *ExerciseEquipment) ToResponse() ExerciseEquipmentResponse {
	return ExerciseEquipmentResponse{
		ID:          ee.ID,
		ExerciseID:  ee.ExerciseID,
		EquipmentID: ee.EquipmentID,
		Optional:    ee.Optional,
		Notes:       ee.Notes,
		Equipment:   ee.Equipment.ToResponse(),
	}
}

// Validation methods
func (e *Equipment) Validate() error {
	if e.Name == "" {
		return gorm.ErrInvalidValue
	}

	validCategories := []string{"machine", "free_weight", "cable", "cardio", "other"}
	if e.Category != "" {
		valid := false
		for _, cat := range validCategories {
			if e.Category == cat {
				valid = true
				break
			}
		}
		if !valid {
			return gorm.ErrInvalidValue
		}
	}

	return nil
}

func (ee *ExerciseEquipment) Validate() error {
	if ee.ExerciseID == uuid.Nil || ee.EquipmentID == uuid.Nil {
		return gorm.ErrInvalidValue
	}
	return nil
}

// Request DTOs
type CreateEquipmentRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category" binding:"omitempty,oneof=machine free_weight cable cardio other"`
	ImageURL    string `json:"image_url"`
}

type UpdateEquipmentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category" binding:"omitempty,oneof=machine free_weight cable cardio other"`
	ImageURL    string `json:"image_url"`
}

type AssignEquipmentToExerciseRequest struct {
	EquipmentID uuid.UUID `json:"equipment_id" binding:"required"`
	Optional    bool      `json:"optional"`
	Notes       string    `json:"notes"`
}