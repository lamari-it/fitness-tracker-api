package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Specialty represents a trainer specialty category
type Specialty struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null;unique" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TrainerSpecialty represents the many-to-many relationship between trainers and specialties
type TrainerSpecialty struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TrainerProfileID uuid.UUID `gorm:"type:uuid;not null"`
	SpecialtyID      uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time

	// Relationships
	TrainerProfile TrainerProfile `gorm:"foreignKey:TrainerProfileID;constraint:OnDelete:CASCADE"`
	Specialty      Specialty      `gorm:"foreignKey:SpecialtyID;constraint:OnDelete:CASCADE"`
}

// TableName overrides the table name for TrainerSpecialty
func (TrainerSpecialty) TableName() string {
	return "trainer_specialties"
}

// SpecialtyResponse is the public response for a specialty
type SpecialtyResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse converts a Specialty model to SpecialtyResponse
func (s *Specialty) ToResponse() SpecialtyResponse {
	return SpecialtyResponse{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}
