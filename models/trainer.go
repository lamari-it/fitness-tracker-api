package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TrainerProfile struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null;unique" json:"user_id"`
	Bio        string         `gorm:"type:text" json:"bio"`
	HourlyRate float64        `gorm:"type:numeric(10,2)" json:"hourly_rate"`
	Location   string         `gorm:"type:text" json:"location"`
	Visibility string         `gorm:"type:varchar(20);default:'public'" json:"visibility"` // public, link_only, private
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	User        User        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Specialties []Specialty `gorm:"many2many:trainer_specialties;" json:"specialties,omitempty"`
}

type TrainerReview struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TrainerID  uuid.UUID      `gorm:"type:uuid;not null" json:"trainer_id"`
	ReviewerID uuid.UUID      `gorm:"type:uuid;not null" json:"reviewer_id"`
	Rating     int            `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Comment    string         `gorm:"type:text" json:"comment"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Trainer  User `gorm:"foreignKey:TrainerID;constraint:OnDelete:CASCADE" json:"trainer,omitempty"`
	Reviewer User `gorm:"foreignKey:ReviewerID;constraint:OnDelete:CASCADE" json:"reviewer,omitempty"`
}

type TrainerClientLink struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TrainerID uuid.UUID      `gorm:"type:uuid;not null" json:"trainer_id"`
	ClientID  uuid.UUID      `gorm:"type:uuid;not null" json:"client_id"`
	Status    string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, active, inactive
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

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

// Request DTOs
type CreateTrainerProfileRequest struct {
	Bio          string      `json:"bio" binding:"omitempty,max=1000"`
	SpecialtyIDs []uuid.UUID `json:"specialty_ids" binding:"omitempty,max=20"`
	HourlyRate   float64     `json:"hourly_rate" binding:"omitempty,gte=0,lte=9999.99"`
	Location     string      `json:"location" binding:"omitempty,max=500"`
	Visibility   string      `json:"visibility" binding:"omitempty,oneof=public link_only private"`
}

type UpdateTrainerProfileRequest struct {
	Bio          string      `json:"bio" binding:"omitempty,max=1000"`
	SpecialtyIDs []uuid.UUID `json:"specialty_ids" binding:"omitempty,max=20"`
	HourlyRate   float64     `json:"hourly_rate" binding:"omitempty,gte=0,lte=9999.99"`
	Location     string      `json:"location" binding:"omitempty,max=500"`
	Visibility   string      `json:"visibility" binding:"omitempty,oneof=public link_only private"`
}

// Response DTOs
type TrainerProfileResponse struct {
	ID          uuid.UUID           `json:"id"`
	UserID      uuid.UUID           `json:"user_id"`
	Bio         string              `json:"bio"`
	Specialties []SpecialtyResponse `json:"specialties"`
	HourlyRate  float64             `json:"hourly_rate"`
	Location    string              `json:"location"`
	Visibility  string              `json:"visibility"`
	User        *UserResponse       `json:"user,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type TrainerPublicResponse struct {
	ID            uuid.UUID           `json:"id"`
	UserID        uuid.UUID           `json:"user_id"`
	Bio           string              `json:"bio"`
	Specialties   []SpecialtyResponse `json:"specialties"`
	HourlyRate    float64             `json:"hourly_rate"`
	Location      string              `json:"location"`
	Visibility    string              `json:"visibility"`
	User          *UserPublicResponse `json:"user"`
	ReviewCount   int                 `json:"review_count"`
	AverageRating float64             `json:"average_rating"`
	CreatedAt     time.Time           `json:"created_at"`
}

type UserPublicResponse struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

// TrainerClientLink Request DTOs
type InviteClientRequest struct {
	ClientID uuid.UUID `json:"client_id" binding:"required"`
}

type RespondToInvitationRequest struct {
	Action string `json:"action" binding:"required,oneof=accept reject"`
}

// TrainerClientLink Response DTOs
type TrainerClientLinkResponse struct {
	ID        uuid.UUID           `json:"id"`
	TrainerID uuid.UUID           `json:"trainer_id"`
	ClientID  uuid.UUID           `json:"client_id"`
	Status    string              `json:"status"`
	Trainer   *UserPublicResponse `json:"trainer,omitempty"`
	Client    *UserPublicResponse `json:"client,omitempty"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

// ToResponse converts TrainerClientLink to response format
func (tcl *TrainerClientLink) ToResponse() TrainerClientLinkResponse {
	resp := TrainerClientLinkResponse{
		ID:        tcl.ID,
		TrainerID: tcl.TrainerID,
		ClientID:  tcl.ClientID,
		Status:    tcl.Status,
		CreatedAt: tcl.CreatedAt,
		UpdatedAt: tcl.UpdatedAt,
	}

	if tcl.Trainer.ID != uuid.Nil {
		resp.Trainer = &UserPublicResponse{
			ID:        tcl.Trainer.ID,
			FirstName: tcl.Trainer.FirstName,
			LastName:  tcl.Trainer.LastName,
		}
	}

	if tcl.Client.ID != uuid.Nil {
		resp.Client = &UserPublicResponse{
			ID:        tcl.Client.ID,
			FirstName: tcl.Client.FirstName,
			LastName:  tcl.Client.LastName,
		}
	}

	return resp
}

// Helper methods
func (tp *TrainerProfile) ToResponse() TrainerProfileResponse {
	// Convert specialties to response format
	specialties := make([]SpecialtyResponse, 0, len(tp.Specialties))
	for _, s := range tp.Specialties {
		specialties = append(specialties, s.ToResponse())
	}

	resp := TrainerProfileResponse{
		ID:          tp.ID,
		UserID:      tp.UserID,
		Bio:         tp.Bio,
		Specialties: specialties,
		HourlyRate:  tp.HourlyRate,
		Location:    tp.Location,
		Visibility:  tp.Visibility,
		CreatedAt:   tp.CreatedAt,
		UpdatedAt:   tp.UpdatedAt,
	}
	if tp.User.ID != uuid.Nil {
		userResp := tp.User.ToResponse()
		resp.User = &userResp
	}
	return resp
}

func (tp *TrainerProfile) ToPublicResponse(reviewCount int, avgRating float64) TrainerPublicResponse {
	// Convert specialties to response format
	specialties := make([]SpecialtyResponse, 0, len(tp.Specialties))
	for _, s := range tp.Specialties {
		specialties = append(specialties, s.ToResponse())
	}

	resp := TrainerPublicResponse{
		ID:            tp.ID,
		UserID:        tp.UserID,
		Bio:           tp.Bio,
		Specialties:   specialties,
		HourlyRate:    tp.HourlyRate,
		Location:      tp.Location,
		Visibility:    tp.Visibility,
		ReviewCount:   reviewCount,
		AverageRating: avgRating,
		CreatedAt:     tp.CreatedAt,
	}
	if tp.User.ID != uuid.Nil {
		resp.User = &UserPublicResponse{
			ID:        tp.User.ID,
			FirstName: tp.User.FirstName,
			LastName:  tp.User.LastName,
		}
	}
	return resp
}

// Validation method
func (tp *TrainerProfile) Validate() error {
	// Bio: allow empty or 1-1000 characters
	if tp.Bio != "" && len(tp.Bio) > 1000 {
		return fmt.Errorf("bio must be at most 1000 characters")
	}
	// Specialties: allow empty or up to 20 items
	if len(tp.Specialties) > 20 {
		return fmt.Errorf("specialties must have at most 20 items")
	}
	// HourlyRate: allow 0-9999.99
	if tp.HourlyRate < 0 || tp.HourlyRate > 9999.99 {
		return fmt.Errorf("hourly rate must be between 0 and 9999.99")
	}
	// Location: allow empty or 1-500 characters
	if tp.Location != "" && len(tp.Location) > 500 {
		return fmt.Errorf("location must be at most 500 characters")
	}
	// Visibility: must be valid enum or empty
	if tp.Visibility != "" && tp.Visibility != "public" && tp.Visibility != "link_only" && tp.Visibility != "private" {
		return fmt.Errorf("visibility must be one of: public, link_only, private")
	}
	return nil
}
