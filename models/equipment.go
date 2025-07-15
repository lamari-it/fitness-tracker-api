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
	Slug        string    `gorm:"type:varchar(100);unique" json:"slug"` // e.g. 'dumbbells', 'barbell'
	Description string    `gorm:"type:text" json:"description"`
	Category    string    `gorm:"type:varchar(50)" json:"category"` // machine, free_weight, cable, cardio, other
	ImageURL    string    `gorm:"type:text" json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	ExerciseLinks []ExerciseEquipment `gorm:"foreignKey:EquipmentID" json:"exercise_links,omitempty"`
	UserEquipment []UserEquipment     `gorm:"foreignKey:EquipmentID" json:"user_equipment,omitempty"`
}

// UserEquipment represents what equipment a user has access to at different locations
type UserEquipment struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	EquipmentID  uuid.UUID `gorm:"type:uuid;not null;index" json:"equipment_id"`
	LocationType string    `gorm:"type:varchar(10);not null" json:"location_type"` // 'home' or 'gym'
	GymLocation  string    `gorm:"type:text" json:"gym_location,omitempty"`        // Optional: specific gym location
	Notes        string    `gorm:"type:text" json:"notes,omitempty"`               // Optional: notes (e.g., "20lb dumbbells")
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relationships
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Equipment Equipment `gorm:"foreignKey:EquipmentID;constraint:OnDelete:CASCADE" json:"equipment,omitempty"`
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
	Slug        string    `json:"slug"`
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

type UserEquipmentResponse struct {
	ID           uuid.UUID         `json:"id"`
	UserID       uuid.UUID         `json:"user_id"`
	EquipmentID  uuid.UUID         `json:"equipment_id"`
	LocationType string            `json:"location_type"`
	GymLocation  string            `json:"gym_location,omitempty"`
	Notes        string            `json:"notes,omitempty"`
	Equipment    EquipmentResponse `json:"equipment"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// Helper methods
func (e *Equipment) ToResponse() EquipmentResponse {
	return EquipmentResponse{
		ID:          e.ID,
		Name:        e.Name,
		Slug:        e.Slug,
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

func (ue *UserEquipment) ToResponse() UserEquipmentResponse {
	return UserEquipmentResponse{
		ID:           ue.ID,
		UserID:       ue.UserID,
		EquipmentID:  ue.EquipmentID,
		LocationType: ue.LocationType,
		GymLocation:  ue.GymLocation,
		Notes:        ue.Notes,
		Equipment:    ue.Equipment.ToResponse(),
		CreatedAt:    ue.CreatedAt,
		UpdatedAt:    ue.UpdatedAt,
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

func (ue *UserEquipment) Validate() error {
	if ue.UserID == uuid.Nil || ue.EquipmentID == uuid.Nil {
		return gorm.ErrInvalidValue
	}
	
	validLocationTypes := []string{"home", "gym"}
	valid := false
	for _, lt := range validLocationTypes {
		if ue.LocationType == lt {
			valid = true
			break
		}
	}
	if !valid {
		return gorm.ErrInvalidValue
	}
	
	return nil
}

// BeforeCreate hook for UserEquipment
func (ue *UserEquipment) BeforeCreate(tx *gorm.DB) (err error) {
	if ue.ID == uuid.Nil {
		ue.ID = uuid.New()
	}
	return
}

// Request DTOs
type CreateEquipmentRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category" binding:"omitempty,oneof=machine free_weight cable cardio other"`
	ImageURL    string `json:"image_url"`
}

type UpdateEquipmentRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Category    string `json:"category" binding:"omitempty,oneof=machine free_weight cable cardio other"`
	ImageURL    string `json:"image_url"`
}

type AssignEquipmentToExerciseRequest struct {
	EquipmentID uuid.UUID `json:"equipment_id" binding:"required"`
	Optional    bool      `json:"optional"`
	Notes       string    `json:"notes"`
}

// User Equipment Request DTOs
type AddUserEquipmentRequest struct {
	EquipmentID  uuid.UUID `json:"equipment_id" binding:"required"`
	LocationType string    `json:"location_type" binding:"required,oneof=home gym"`
	GymLocation  string    `json:"gym_location,omitempty"`
	Notes        string    `json:"notes,omitempty"`
}

type UpdateUserEquipmentRequest struct {
	LocationType string `json:"location_type" binding:"omitempty,oneof=home gym"`
	GymLocation  string `json:"gym_location,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

type UserEquipmentFilter struct {
	LocationType string `form:"location_type" binding:"omitempty,oneof=home gym"`
}