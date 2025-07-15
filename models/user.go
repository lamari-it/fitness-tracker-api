package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email          string     `gorm:"unique;not null" json:"email"`
	Password       string     `gorm:"not null" json:"-"`
	FirstName      string     `gorm:"not null" json:"first_name"`
	LastName       string     `gorm:"not null" json:"last_name"`
	Provider       string     `gorm:"default:'local'" json:"provider"`
	GoogleID       string     `gorm:"unique" json:"google_id,omitempty"`
	AppleID        string     `gorm:"unique" json:"apple_id,omitempty"`
	FitnessLevelID *uuid.UUID `gorm:"type:uuid" json:"fitness_level_id,omitempty"`
	IsActive       bool       `gorm:"default:true" json:"is_active"`
	IsAdmin        bool       `gorm:"default:false" json:"is_admin"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Relationships
	FitnessLevel   *FitnessLevel      `gorm:"foreignKey:FitnessLevelID;constraint:OnDelete:SET NULL" json:"fitness_level,omitempty"`
	FitnessGoals   []UserFitnessGoal  `gorm:"foreignKey:UserID" json:"fitness_goals,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

type UserResponse struct {
	ID             uuid.UUID                  `json:"id"`
	Email          string                     `json:"email"`
	FirstName      string                     `json:"first_name"`
	LastName       string                     `json:"last_name"`
	Provider       string                     `json:"provider"`
	FitnessLevelID *uuid.UUID                 `json:"fitness_level_id,omitempty"`
	FitnessLevel   *FitnessLevelResponse      `json:"fitness_level,omitempty"`
	FitnessGoals   []UserFitnessGoalResponse  `json:"fitness_goals,omitempty"`
	IsActive       bool                       `json:"is_active"`
	IsAdmin        bool                       `json:"is_admin"`
	CreatedAt      time.Time                  `json:"created_at"`
	UpdatedAt      time.Time                  `json:"updated_at"`
}

func (u *User) ToResponse() UserResponse {
	response := UserResponse{
		ID:             u.ID,
		Email:          u.Email,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Provider:       u.Provider,
		FitnessLevelID: u.FitnessLevelID,
		IsActive:       u.IsActive,
		IsAdmin:        u.IsAdmin,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}

	if u.FitnessLevel != nil {
		levelResponse := u.FitnessLevel.ToResponse()
		response.FitnessLevel = &levelResponse
	}

	if len(u.FitnessGoals) > 0 {
		response.FitnessGoals = make([]UserFitnessGoalResponse, len(u.FitnessGoals))
		for i, goal := range u.FitnessGoals {
			response.FitnessGoals[i] = goal.ToResponse()
		}
	}

	return response
}