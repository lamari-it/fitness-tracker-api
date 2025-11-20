package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// UserFitnessProfile stores user's physical stats, goals, and training preferences
type UserFitnessProfile struct {
	ID     uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`

	// Basic physical stats (required)
	DateOfBirth  time.Time `gorm:"type:date;not null" json:"date_of_birth"`
	Gender       string    `gorm:"type:varchar(20);not null" json:"gender"`
	HeightCm     float64   `gorm:"type:decimal(5,2);not null" json:"height_cm"`
	CurrentWeightKg float64 `gorm:"type:decimal(5,2);not null" json:"current_weight_kg"`

	// Unit preference
	PreferredUnitSystem string `gorm:"type:varchar(10);not null;default:'metric'" json:"preferred_unit_system"`

	// Fitness goals
	PrimaryGoal         string   `gorm:"type:varchar(50);not null" json:"primary_goal"`
	TargetWeightKg      *float64 `gorm:"type:decimal(5,2)" json:"target_weight_kg,omitempty"`
	TargetWeeklyWorkouts int     `gorm:"not null;default:3" json:"target_weekly_workouts"`

	// Activity level
	ActivityLevel string `gorm:"type:varchar(20);not null;default:'moderate'" json:"activity_level"`

	// Training preferences
	TrainingLocations          pq.StringArray `gorm:"type:text[];default:ARRAY['gym']::text[]" json:"training_locations"`
	PreferredWorkoutDurationMins int          `gorm:"not null;default:45" json:"preferred_workout_duration_mins"`
	AvailableDays              pq.StringArray `gorm:"type:text[];default:ARRAY['monday','wednesday','friday']::text[]" json:"available_days"`

	// Additional info (optional)
	HealthConditions string `gorm:"type:text" json:"health_conditions,omitempty"`
	InjuriesNotes    string `gorm:"type:text" json:"injuries_notes,omitempty"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

// Request DTOs

type CreateUserFitnessProfileRequest struct {
	// Required fields
	DateOfBirth      string  `json:"date_of_birth" binding:"required"`
	Gender           string  `json:"gender" binding:"required,oneof=male female other prefer_not_to_say"`
	HeightCm         float64 `json:"height_cm" binding:"required,gt=50,lt=300"`
	CurrentWeightKg  float64 `json:"current_weight_kg" binding:"required,gt=20,lt=500"`
	PrimaryGoal      string  `json:"primary_goal" binding:"required,oneof=weight_loss muscle_gain strength general_fitness endurance flexibility"`

	// Optional fields with defaults
	PreferredUnitSystem       string   `json:"preferred_unit_system" binding:"omitempty,oneof=metric imperial"`
	TargetWeightKg            *float64 `json:"target_weight_kg" binding:"omitempty,gt=20,lt=500"`
	TargetWeeklyWorkouts      int      `json:"target_weekly_workouts" binding:"omitempty,min=1,max=7"`
	ActivityLevel             string   `json:"activity_level" binding:"omitempty,oneof=sedentary lightly_active moderate active very_active"`
	TrainingLocations         []string `json:"training_locations" binding:"omitempty,dive,oneof=home gym outdoors"`
	PreferredWorkoutDurationMins int   `json:"preferred_workout_duration_mins" binding:"omitempty,min=10,max=180"`
	AvailableDays             []string `json:"available_days" binding:"omitempty,dive,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	HealthConditions          string   `json:"health_conditions" binding:"omitempty,max=1000"`
	InjuriesNotes             string   `json:"injuries_notes" binding:"omitempty,max=1000"`
}

type UpdateUserFitnessProfileRequest struct {
	DateOfBirth               string   `json:"date_of_birth" binding:"omitempty"`
	Gender                    string   `json:"gender" binding:"omitempty,oneof=male female other prefer_not_to_say"`
	HeightCm                  float64  `json:"height_cm" binding:"omitempty,gt=50,lt=300"`
	CurrentWeightKg           float64  `json:"current_weight_kg" binding:"omitempty,gt=20,lt=500"`
	PreferredUnitSystem       string   `json:"preferred_unit_system" binding:"omitempty,oneof=metric imperial"`
	PrimaryGoal               string   `json:"primary_goal" binding:"omitempty,oneof=weight_loss muscle_gain strength general_fitness endurance flexibility"`
	TargetWeightKg            *float64 `json:"target_weight_kg" binding:"omitempty,gt=20,lt=500"`
	TargetWeeklyWorkouts      int      `json:"target_weekly_workouts" binding:"omitempty,min=1,max=7"`
	ActivityLevel             string   `json:"activity_level" binding:"omitempty,oneof=sedentary lightly_active moderate active very_active"`
	TrainingLocations         []string `json:"training_locations" binding:"omitempty,dive,oneof=home gym outdoors"`
	PreferredWorkoutDurationMins int   `json:"preferred_workout_duration_mins" binding:"omitempty,min=10,max=180"`
	AvailableDays             []string `json:"available_days" binding:"omitempty,dive,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	HealthConditions          string   `json:"health_conditions" binding:"omitempty,max=1000"`
	InjuriesNotes             string   `json:"injuries_notes" binding:"omitempty,max=1000"`
}

// Response DTOs

type UserFitnessProfileResponse struct {
	ID                        uuid.UUID      `json:"id"`
	UserID                    uuid.UUID      `json:"user_id"`
	DateOfBirth               string         `json:"date_of_birth"`
	Age                       int            `json:"age"`
	Gender                    string         `json:"gender"`
	HeightCm                  float64        `json:"height_cm"`
	HeightFtIn                string         `json:"height_ft_in"`
	CurrentWeightKg           float64        `json:"current_weight_kg"`
	CurrentWeightLbs          float64        `json:"current_weight_lbs"`
	PreferredUnitSystem       string         `json:"preferred_unit_system"`
	PrimaryGoal               string         `json:"primary_goal"`
	TargetWeightKg            *float64       `json:"target_weight_kg,omitempty"`
	TargetWeightLbs           *float64       `json:"target_weight_lbs,omitempty"`
	TargetWeeklyWorkouts      int            `json:"target_weekly_workouts"`
	ActivityLevel             string         `json:"activity_level"`
	TrainingLocations         []string       `json:"training_locations"`
	PreferredWorkoutDurationMins int         `json:"preferred_workout_duration_mins"`
	AvailableDays             []string       `json:"available_days"`
	HealthConditions          string         `json:"health_conditions,omitempty"`
	InjuriesNotes             string         `json:"injuries_notes,omitempty"`
	CreatedAt                 time.Time      `json:"created_at"`
	UpdatedAt                 time.Time      `json:"updated_at"`
}

// ToResponse converts the model to a response DTO with unit conversions
func (p *UserFitnessProfile) ToResponse() UserFitnessProfileResponse {
	// Calculate age
	now := time.Now()
	age := now.Year() - p.DateOfBirth.Year()
	if now.YearDay() < p.DateOfBirth.YearDay() {
		age--
	}

	// Convert height to feet and inches
	totalInches := p.HeightCm / 2.54
	feet := int(totalInches / 12)
	inches := int(totalInches) % 12
	heightFtIn := formatFeetInches(feet, inches)

	// Convert weights to lbs
	currentWeightLbs := p.CurrentWeightKg * 2.20462

	var targetWeightLbs *float64
	if p.TargetWeightKg != nil {
		lbs := *p.TargetWeightKg * 2.20462
		targetWeightLbs = &lbs
	}

	return UserFitnessProfileResponse{
		ID:                        p.ID,
		UserID:                    p.UserID,
		DateOfBirth:               p.DateOfBirth.Format("2006-01-02"),
		Age:                       age,
		Gender:                    p.Gender,
		HeightCm:                  p.HeightCm,
		HeightFtIn:                heightFtIn,
		CurrentWeightKg:           p.CurrentWeightKg,
		CurrentWeightLbs:          currentWeightLbs,
		PreferredUnitSystem:       p.PreferredUnitSystem,
		PrimaryGoal:               p.PrimaryGoal,
		TargetWeightKg:            p.TargetWeightKg,
		TargetWeightLbs:           targetWeightLbs,
		TargetWeeklyWorkouts:      p.TargetWeeklyWorkouts,
		ActivityLevel:             p.ActivityLevel,
		TrainingLocations:         p.TrainingLocations,
		PreferredWorkoutDurationMins: p.PreferredWorkoutDurationMins,
		AvailableDays:             p.AvailableDays,
		HealthConditions:          p.HealthConditions,
		InjuriesNotes:             p.InjuriesNotes,
		CreatedAt:                 p.CreatedAt,
		UpdatedAt:                 p.UpdatedAt,
	}
}

// Helper function to format feet and inches
func formatFeetInches(feet, inches int) string {
	return fmt.Sprintf("%d'%d\"", feet, inches)
}
