package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserFavoriteExercise represents a user's favorited exercise
type UserFavoriteExercise struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:unique_user_favorite_exercise" json:"user_id"`
	ExerciseID uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:unique_user_favorite_exercise" json:"exercise_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	User     User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Exercise Exercise `gorm:"foreignKey:ExerciseID;constraint:OnDelete:CASCADE" json:"exercise,omitempty"`
}

// UserFavoriteWorkout represents a user's favorited workout
type UserFavoriteWorkout struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:unique_user_favorite_workout" json:"user_id"`
	WorkoutID uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:unique_user_favorite_workout" json:"workout_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	User    User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Workout Workout `gorm:"foreignKey:WorkoutID;constraint:OnDelete:CASCADE" json:"workout,omitempty"`
}

// BeforeCreate hooks
func (ufe *UserFavoriteExercise) BeforeCreate(tx *gorm.DB) (err error) {
	if ufe.ID == uuid.Nil {
		ufe.ID = uuid.New()
	}
	return
}

func (ufw *UserFavoriteWorkout) BeforeCreate(tx *gorm.DB) (err error) {
	if ufw.ID == uuid.Nil {
		ufw.ID = uuid.New()
	}
	return
}

// Response DTOs

// FavoriteExerciseResponse for exercise favorites
type FavoriteExerciseResponse struct {
	ID         uuid.UUID `json:"id"`
	ExerciseID uuid.UUID `json:"exercise_id"`
	Exercise   Exercise  `json:"exercise"`
	CreatedAt  time.Time `json:"created_at"`
}

// FavoriteWorkoutResponse for workout favorites
type FavoriteWorkoutResponse struct {
	ID        uuid.UUID `json:"id"`
	WorkoutID uuid.UUID `json:"workout_id"`
	Workout   Workout   `json:"workout"`
	CreatedAt time.Time `json:"created_at"`
}

// FavoriteResponse - unified response for the combined endpoint
type FavoriteResponse struct {
	ID        uuid.UUID   `json:"id"`
	Type      string      `json:"type"` // "exercise" or "workout"
	ItemID    uuid.UUID   `json:"item_id"`
	Item      interface{} `json:"item"`
	CreatedAt time.Time   `json:"created_at"`
}

// ToResponse methods
func (ufe *UserFavoriteExercise) ToResponse() FavoriteExerciseResponse {
	return FavoriteExerciseResponse{
		ID:         ufe.ID,
		ExerciseID: ufe.ExerciseID,
		Exercise:   ufe.Exercise,
		CreatedAt:  ufe.CreatedAt,
	}
}

func (ufw *UserFavoriteWorkout) ToResponse() FavoriteWorkoutResponse {
	return FavoriteWorkoutResponse{
		ID:        ufw.ID,
		WorkoutID: ufw.WorkoutID,
		Workout:   ufw.Workout,
		CreatedAt: ufw.CreatedAt,
	}
}

// ToGenericResponse converts to unified FavoriteResponse
func (ufe *UserFavoriteExercise) ToGenericResponse() FavoriteResponse {
	return FavoriteResponse{
		ID:        ufe.ID,
		Type:      "exercise",
		ItemID:    ufe.ExerciseID,
		Item:      ufe.Exercise,
		CreatedAt: ufe.CreatedAt,
	}
}

func (ufw *UserFavoriteWorkout) ToGenericResponse() FavoriteResponse {
	return FavoriteResponse{
		ID:        ufw.ID,
		Type:      "workout",
		ItemID:    ufw.WorkoutID,
		Item:      ufw.Workout,
		CreatedAt: ufw.CreatedAt,
	}
}

// Request DTOs
type AddFavoriteRequest struct {
	ItemID uuid.UUID `json:"item_id" binding:"required"`
}
