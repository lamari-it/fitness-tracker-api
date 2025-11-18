package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FitnessLevel represents different fitness experience levels
type FitnessLevel struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(50);not null;unique" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"` // For ordering levels (beginner->advanced)
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Users []User `gorm:"foreignKey:FitnessLevelID" json:"users,omitempty"`
}

// FitnessGoal represents different fitness objectives
type FitnessGoal struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null;unique" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Category    string         `gorm:"type:varchar(50)" json:"category"`  // strength, cardio, flexibility, etc.
	IconName    string         `gorm:"type:varchar(50)" json:"icon_name"` // For UI icon display
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	UserGoals []UserFitnessGoal `gorm:"foreignKey:FitnessGoalID" json:"user_goals,omitempty"`
}

// UserFitnessGoal represents the many-to-many relationship between users and fitness goals
type UserFitnessGoal struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:unique_user_fitness_goal_combo" json:"user_id"`
	FitnessGoalID uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:unique_user_fitness_goal_combo" json:"fitness_goal_id"`
	Priority      int            `gorm:"default:0" json:"priority"` // 1=primary, 2=secondary, etc.
	TargetDate    *time.Time     `json:"target_date,omitempty"`     // Optional target completion date
	Notes         string         `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	User        User        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	FitnessGoal FitnessGoal `gorm:"foreignKey:FitnessGoalID;constraint:OnDelete:CASCADE" json:"fitness_goal,omitempty"`
}

// BeforeCreate hooks
func (fl *FitnessLevel) BeforeCreate(tx *gorm.DB) (err error) {
	if fl.ID == uuid.Nil {
		fl.ID = uuid.New()
	}
	return
}

func (fg *FitnessGoal) BeforeCreate(tx *gorm.DB) (err error) {
	if fg.ID == uuid.Nil {
		fg.ID = uuid.New()
	}
	return
}

func (ufg *UserFitnessGoal) BeforeCreate(tx *gorm.DB) (err error) {
	if ufg.ID == uuid.Nil {
		ufg.ID = uuid.New()
	}
	return
}

// Response DTOs
type FitnessLevelResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FitnessGoalResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	IconName    string    `json:"icon_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserFitnessGoalResponse struct {
	ID            uuid.UUID           `json:"id"`
	UserID        uuid.UUID           `json:"user_id"`
	FitnessGoalID uuid.UUID           `json:"fitness_goal_id"`
	Priority      int                 `json:"priority"`
	TargetDate    *time.Time          `json:"target_date,omitempty"`
	Notes         string              `json:"notes"`
	FitnessGoal   FitnessGoalResponse `json:"fitness_goal"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
}

// Helper methods
func (fl *FitnessLevel) ToResponse() FitnessLevelResponse {
	return FitnessLevelResponse{
		ID:          fl.ID,
		Name:        fl.Name,
		Description: fl.Description,
		SortOrder:   fl.SortOrder,
		CreatedAt:   fl.CreatedAt,
		UpdatedAt:   fl.UpdatedAt,
	}
}

func (fg *FitnessGoal) ToResponse() FitnessGoalResponse {
	return FitnessGoalResponse{
		ID:          fg.ID,
		Name:        fg.Name,
		Description: fg.Description,
		Category:    fg.Category,
		IconName:    fg.IconName,
		CreatedAt:   fg.CreatedAt,
		UpdatedAt:   fg.UpdatedAt,
	}
}

func (ufg *UserFitnessGoal) ToResponse() UserFitnessGoalResponse {
	return UserFitnessGoalResponse{
		ID:            ufg.ID,
		UserID:        ufg.UserID,
		FitnessGoalID: ufg.FitnessGoalID,
		Priority:      ufg.Priority,
		TargetDate:    ufg.TargetDate,
		Notes:         ufg.Notes,
		FitnessGoal:   ufg.FitnessGoal.ToResponse(),
		CreatedAt:     ufg.CreatedAt,
		UpdatedAt:     ufg.UpdatedAt,
	}
}

// Request DTOs
type CreateFitnessLevelRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

type UpdateFitnessLevelRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

type CreateFitnessGoalRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category"`
	IconName    string `json:"icon_name"`
}

type UpdateFitnessGoalRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	IconName    string `json:"icon_name"`
}

type SetUserFitnessGoalsRequest struct {
	Goals []UserFitnessGoalInput `json:"goals" binding:"required,dive"`
}

type UserFitnessGoalInput struct {
	FitnessGoalID uuid.UUID  `json:"fitness_goal_id" binding:"required"`
	Priority      int        `json:"priority"`
	TargetDate    *time.Time `json:"target_date,omitempty"`
	Notes         string     `json:"notes"`
}

type UpdateUserFitnessLevelRequest struct {
	FitnessLevelID *uuid.UUID `json:"fitness_level_id"`
}
