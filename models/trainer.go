package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type TrainerProfile struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;unique" json:"user_id"`
	Bio         string         `gorm:"type:text" json:"bio"`
	Specialties pq.StringArray `gorm:"type:text[]" json:"specialties"`
	HourlyRate  float64        `gorm:"type:numeric(10,2)" json:"hourly_rate"`
	Location    string         `gorm:"type:text" json:"location"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

type TrainerReview struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TrainerID  uuid.UUID `gorm:"type:uuid;not null" json:"trainer_id"`
	ReviewerID uuid.UUID `gorm:"type:uuid;not null" json:"reviewer_id"`
	Rating     int       `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Comment    string    `gorm:"type:text" json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Trainer  User `gorm:"foreignKey:TrainerID;constraint:OnDelete:CASCADE" json:"trainer,omitempty"`
	Reviewer User `gorm:"foreignKey:ReviewerID;constraint:OnDelete:CASCADE" json:"reviewer,omitempty"`
}

type TrainerClientLink struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TrainerID uuid.UUID `gorm:"type:uuid;not null" json:"trainer_id"`
	ClientID  uuid.UUID `gorm:"type:uuid;not null" json:"client_id"`
	Status    string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, active, inactive
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Trainer User `gorm:"foreignKey:TrainerID;constraint:OnDelete:CASCADE" json:"trainer,omitempty"`
	Client  User `gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE" json:"client,omitempty"`
}

func (tp *TrainerProfile) BeforeCreate(tx *gorm.DB) (err error) {
	if tp.ID == uuid.Nil {
		tp.ID = uuid.New()
	}
	return
}

func (tr *TrainerReview) BeforeCreate(tx *gorm.DB) (err error) {
	if tr.ID == uuid.Nil {
		tr.ID = uuid.New()
	}
	return
}

func (tcl *TrainerClientLink) BeforeCreate(tx *gorm.DB) (err error) {
	if tcl.ID == uuid.Nil {
		tcl.ID = uuid.New()
	}
	return
}