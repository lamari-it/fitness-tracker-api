package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TrainerInvitation stores email-based invitations from trainers to potential clients
type TrainerInvitation struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TrainerID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"trainer_id"`
	InviteeEmail    string         `gorm:"type:varchar(255);not null;index" json:"invitee_email"`
	InvitationToken string         `gorm:"type:varchar(64);not null;uniqueIndex" json:"-"`
	Status          string         `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	ExpiresAt       time.Time      `gorm:"not null" json:"expires_at"`
	AcceptedAt      *time.Time     `json:"accepted_at,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Trainer User `gorm:"foreignKey:TrainerID;constraint:OnDelete:CASCADE" json:"trainer,omitempty"`
}

// Invitation statuses
const (
	InvitationStatusPending  = "pending"
	InvitationStatusAccepted = "accepted"
	InvitationStatusRejected = "rejected"
	InvitationStatusExpired  = "expired"
)

// Request DTOs

type CreateEmailInvitationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// Response DTOs

type TrainerInvitationResponse struct {
	ID           uuid.UUID           `json:"id"`
	TrainerID    uuid.UUID           `json:"trainer_id"`
	InviteeEmail string              `json:"invitee_email"`
	Status       string              `json:"status"`
	ExpiresAt    time.Time           `json:"expires_at"`
	AcceptedAt   *time.Time          `json:"accepted_at,omitempty"`
	CreatedAt    time.Time           `json:"created_at"`
	Trainer      *TrainerInfoResponse `json:"trainer,omitempty"`
}

type TrainerInfoResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

type VerifyInvitationResponse struct {
	Valid        bool                 `json:"valid"`
	InviteeEmail string               `json:"invitee_email,omitempty"`
	Trainer      *TrainerInfoResponse `json:"trainer,omitempty"`
	ExpiresAt    time.Time            `json:"expires_at,omitempty"`
	Message      string               `json:"message,omitempty"`
}

// ToResponse converts the model to a response DTO
func (i *TrainerInvitation) ToResponse() TrainerInvitationResponse {
	resp := TrainerInvitationResponse{
		ID:           i.ID,
		TrainerID:    i.TrainerID,
		InviteeEmail: i.InviteeEmail,
		Status:       i.Status,
		ExpiresAt:    i.ExpiresAt,
		AcceptedAt:   i.AcceptedAt,
		CreatedAt:    i.CreatedAt,
	}

	// Include trainer info if loaded
	if i.Trainer.ID != uuid.Nil {
		resp.Trainer = &TrainerInfoResponse{
			ID:        i.Trainer.ID,
			Email:     i.Trainer.Email,
			FirstName: i.Trainer.FirstName,
			LastName:  i.Trainer.LastName,
		}
	}

	return resp
}

// IsExpired checks if the invitation has expired
func (i *TrainerInvitation) IsExpired() bool {
	return time.Now().After(i.ExpiresAt)
}
