package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ExerciseType represents a type/category of exercise (e.g., compound, isolation, cardio)
type ExerciseType struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Slug        string         `gorm:"type:varchar(50);not null;unique" json:"slug"`
	Name        string         `gorm:"type:varchar(100);not null;unique" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	ExerciseLinks []ExerciseExerciseType `gorm:"foreignKey:ExerciseTypeID" json:"exercise_links,omitempty"`
}

// ExerciseExerciseType is the junction table for many-to-many relationship between Exercise and ExerciseType
type ExerciseExerciseType struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ExerciseID     uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:unique_exercise_type_combo" json:"exercise_id"`
	ExerciseTypeID uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:unique_exercise_type_combo" json:"exercise_type_id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Exercise     Exercise     `gorm:"foreignKey:ExerciseID;constraint:OnDelete:CASCADE" json:"exercise,omitempty"`
	ExerciseType ExerciseType `gorm:"foreignKey:ExerciseTypeID;constraint:OnDelete:CASCADE" json:"exercise_type,omitempty"`
}

// BeforeCreate hooks
func (et *ExerciseType) BeforeCreate(tx *gorm.DB) (err error) {
	if et.ID == uuid.Nil {
		et.ID = uuid.New()
	}
	return
}

func (eet *ExerciseExerciseType) BeforeCreate(tx *gorm.DB) (err error) {
	if eet.ID == uuid.Nil {
		eet.ID = uuid.New()
	}
	return
}

// Response DTOs
type ExerciseTypeResponse struct {
	ID          uuid.UUID `json:"id"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ExerciseExerciseTypeResponse struct {
	ID             uuid.UUID            `json:"id"`
	ExerciseID     uuid.UUID            `json:"exercise_id"`
	ExerciseTypeID uuid.UUID            `json:"exercise_type_id"`
	ExerciseType   ExerciseTypeResponse `json:"exercise_type"`
}

type ExerciseTypeWithExercises struct {
	ExerciseTypeResponse
	ExerciseCount int                            `json:"exercise_count"`
	Exercises     []ExerciseExerciseTypeResponse `json:"exercises,omitempty"`
}

// ToResponse converts ExerciseType to ExerciseTypeResponse
func (et *ExerciseType) ToResponse() ExerciseTypeResponse {
	return ExerciseTypeResponse{
		ID:          et.ID,
		Slug:        et.Slug,
		Name:        et.Name,
		Description: et.Description,
		CreatedAt:   et.CreatedAt,
		UpdatedAt:   et.UpdatedAt,
	}
}

// ToResponse converts ExerciseExerciseType to ExerciseExerciseTypeResponse
func (eet *ExerciseExerciseType) ToResponse() ExerciseExerciseTypeResponse {
	return ExerciseExerciseTypeResponse{
		ID:             eet.ID,
		ExerciseID:     eet.ExerciseID,
		ExerciseTypeID: eet.ExerciseTypeID,
		ExerciseType:   eet.ExerciseType.ToResponse(),
	}
}

// Request DTOs
type CreateExerciseTypeRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
}

type UpdateExerciseTypeRequest struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
}

type AssignExerciseTypeRequest struct {
	ExerciseTypeID uuid.UUID `json:"exercise_type_id" binding:"required"`
}
