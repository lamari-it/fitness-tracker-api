package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MuscleGroup struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null;unique" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Category    string    `gorm:"type:varchar(50)" json:"category"` // upper, lower, core, cardio
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	ExerciseLinks []ExerciseMuscleGroup `gorm:"foreignKey:MuscleGroupID" json:"exercise_links,omitempty"`
}

type ExerciseMuscleGroup struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ExerciseID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:unique_exercise_muscle_combo" json:"exercise_id"`
	MuscleGroupID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:unique_exercise_muscle_combo" json:"muscle_group_id"`
	Primary       bool      `gorm:"default:false" json:"primary"`
	Intensity     string    `gorm:"type:varchar(20);default:'moderate'" json:"intensity"` // high, moderate, low
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relationships
	Exercise    Exercise    `gorm:"foreignKey:ExerciseID;constraint:OnDelete:CASCADE" json:"exercise,omitempty"`
	MuscleGroup MuscleGroup `gorm:"foreignKey:MuscleGroupID;constraint:OnDelete:CASCADE" json:"muscle_group,omitempty"`
}

// BeforeCreate hooks
func (mg *MuscleGroup) BeforeCreate(tx *gorm.DB) (err error) {
	if mg.ID == uuid.Nil {
		mg.ID = uuid.New()
	}
	return
}

func (emg *ExerciseMuscleGroup) BeforeCreate(tx *gorm.DB) (err error) {
	if emg.ID == uuid.Nil {
		emg.ID = uuid.New()
	}
	return
}

// Response DTOs
type MuscleGroupResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ExerciseMuscleGroupResponse struct {
	ID            uuid.UUID           `json:"id"`
	ExerciseID    uuid.UUID           `json:"exercise_id"`
	MuscleGroupID uuid.UUID           `json:"muscle_group_id"`
	Primary       bool                `json:"primary"`
	Intensity     string              `json:"intensity"`
	MuscleGroup   MuscleGroupResponse `json:"muscle_group"`
}

type MuscleGroupWithExercises struct {
	MuscleGroupResponse
	ExerciseCount int                           `json:"exercise_count"`
	Exercises     []ExerciseMuscleGroupResponse `json:"exercises,omitempty"`
}

// Helper methods
func (mg *MuscleGroup) ToResponse() MuscleGroupResponse {
	return MuscleGroupResponse{
		ID:          mg.ID,
		Name:        mg.Name,
		Description: mg.Description,
		Category:    mg.Category,
		CreatedAt:   mg.CreatedAt,
		UpdatedAt:   mg.UpdatedAt,
	}
}

func (emg *ExerciseMuscleGroup) ToResponse() ExerciseMuscleGroupResponse {
	return ExerciseMuscleGroupResponse{
		ID:            emg.ID,
		ExerciseID:    emg.ExerciseID,
		MuscleGroupID: emg.MuscleGroupID,
		Primary:       emg.Primary,
		Intensity:     emg.Intensity,
		MuscleGroup:   emg.MuscleGroup.ToResponse(),
	}
}

// Validation methods
func (mg *MuscleGroup) Validate() error {
	if mg.Name == "" {
		return gorm.ErrInvalidValue
	}
	return nil
}

func (emg *ExerciseMuscleGroup) Validate() error {
	if emg.ExerciseID == uuid.Nil || emg.MuscleGroupID == uuid.Nil {
		return gorm.ErrInvalidValue
	}

	validIntensities := []string{"high", "moderate", "low"}
	if emg.Intensity != "" {
		valid := false
		for _, intensity := range validIntensities {
			if emg.Intensity == intensity {
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
