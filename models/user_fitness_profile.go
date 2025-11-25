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
	DateOfBirth                time.Time `gorm:"type:date;not null" json:"date_of_birth"`
	Gender                     string    `gorm:"type:varchar(20);not null" json:"gender"`
	HeightCm                   float64   `gorm:"type:decimal(5,2);not null" json:"height_cm"`
	CurrentWeightKg            *float64  `gorm:"type:decimal(6,2)" json:"-"`
	OriginalCurrentWeightValue *float64  `gorm:"type:decimal(6,2)" json:"-"`
	OriginalCurrentWeightUnit  *string   `gorm:"type:varchar(2)" json:"-"`

	// Fitness goals
	TargetWeightKg            *float64 `gorm:"type:decimal(6,2)" json:"-"`
	OriginalTargetWeightValue *float64 `gorm:"type:decimal(6,2)" json:"-"`
	OriginalTargetWeightUnit  *string  `gorm:"type:varchar(2)" json:"-"`
	TargetWeeklyWorkouts      int      `gorm:"not null;default:3" json:"target_weekly_workouts"`

	// Activity level
	ActivityLevel string `gorm:"type:varchar(20);not null;default:'moderate'" json:"activity_level"`

	// Training preferences
	TrainingLocations            pq.StringArray `gorm:"type:text[];default:ARRAY['gym']::text[]" json:"training_locations"`
	PreferredWorkoutDurationMins int            `gorm:"not null;default:45" json:"preferred_workout_duration_mins"`
	AvailableDays                pq.StringArray `gorm:"type:text[];default:ARRAY['monday','wednesday','friday']::text[]" json:"available_days"`

	// Additional info (optional)
	HealthConditions string `gorm:"type:text" json:"health_conditions,omitempty"`
	InjuriesNotes    string `gorm:"type:text" json:"injuries_notes,omitempty"`

	// Fitness level (moved from User table)
	FitnessLevelID *uuid.UUID `gorm:"type:uuid" json:"fitness_level_id,omitempty"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User         User              `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	FitnessGoals []UserFitnessGoal `gorm:"foreignKey:UserFitnessProfileID" json:"fitness_goals,omitempty"`
	FitnessLevel *FitnessLevel     `gorm:"foreignKey:FitnessLevelID;constraint:OnDelete:SET NULL" json:"fitness_level,omitempty"`
}

// Request DTOs

type CreateUserFitnessProfileRequest struct {
	// Required fields
	DateOfBirth    string      `json:"date_of_birth" binding:"required"`
	Gender         string      `json:"gender" binding:"required,oneof=male female other prefer_not_to_say"`
	HeightCm       float64     `json:"height_cm" binding:"required,gt=50,lt=300"`
	CurrentWeight  WeightInput `json:"current_weight" binding:"required"`
	FitnessGoalIDs []uuid.UUID `json:"fitness_goal_ids" binding:"required,min=1,max=5"`

	// Optional fields with defaults
	TargetWeight                 *WeightInput `json:"target_weight,omitempty"`
	TargetWeeklyWorkouts         int          `json:"target_weekly_workouts" binding:"omitempty,min=1,max=7"`
	ActivityLevel                string       `json:"activity_level" binding:"omitempty,oneof=sedentary lightly_active moderate active very_active"`
	TrainingLocations            []string     `json:"training_locations" binding:"omitempty,dive,oneof=home gym outdoors"`
	PreferredWorkoutDurationMins int          `json:"preferred_workout_duration_mins" binding:"omitempty,min=10,max=180"`
	AvailableDays                []string     `json:"available_days" binding:"omitempty,dive,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	HealthConditions             string       `json:"health_conditions" binding:"omitempty,max=1000"`
	InjuriesNotes                string       `json:"injuries_notes" binding:"omitempty,max=1000"`
	FitnessLevelID               *uuid.UUID   `json:"fitness_level_id,omitempty"`
}

type UpdateUserFitnessProfileRequest struct {
	DateOfBirth                  string       `json:"date_of_birth" binding:"omitempty"`
	Gender                       string       `json:"gender" binding:"omitempty,oneof=male female other prefer_not_to_say"`
	HeightCm                     float64      `json:"height_cm" binding:"omitempty,gt=50,lt=300"`
	CurrentWeight                *WeightInput `json:"current_weight,omitempty"`
	FitnessGoalIDs               []uuid.UUID  `json:"fitness_goal_ids" binding:"omitempty,max=5"`
	TargetWeight                 *WeightInput `json:"target_weight,omitempty"`
	TargetWeeklyWorkouts         int          `json:"target_weekly_workouts" binding:"omitempty,min=1,max=7"`
	ActivityLevel                string       `json:"activity_level" binding:"omitempty,oneof=sedentary lightly_active moderate active very_active"`
	TrainingLocations            []string     `json:"training_locations" binding:"omitempty,dive,oneof=home gym outdoors"`
	PreferredWorkoutDurationMins int          `json:"preferred_workout_duration_mins" binding:"omitempty,min=10,max=180"`
	AvailableDays                []string     `json:"available_days" binding:"omitempty,dive,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	HealthConditions             string       `json:"health_conditions" binding:"omitempty,max=1000"`
	InjuriesNotes                string       `json:"injuries_notes" binding:"omitempty,max=1000"`
	FitnessLevelID               *uuid.UUID   `json:"fitness_level_id,omitempty"`
}

// Response DTOs

type UserFitnessProfileResponse struct {
	ID                           uuid.UUID                 `json:"id"`
	UserID                       uuid.UUID                 `json:"user_id"`
	DateOfBirth                  string                    `json:"date_of_birth"`
	Age                          int                       `json:"age"`
	Gender                       string                    `json:"gender"`
	HeightCm                     float64                   `json:"height_cm"`
	HeightFtIn                   string                    `json:"height_ft_in"`
	CurrentWeight                WeightOutput              `json:"current_weight"`
	FitnessGoals                 []UserFitnessGoalResponse `json:"fitness_goals"`
	TargetWeight                 *WeightOutput             `json:"target_weight,omitempty"`
	TargetWeeklyWorkouts         int                       `json:"target_weekly_workouts"`
	ActivityLevel                string                    `json:"activity_level"`
	TrainingLocations            []string                  `json:"training_locations"`
	PreferredWorkoutDurationMins int                       `json:"preferred_workout_duration_mins"`
	AvailableDays                []string                  `json:"available_days"`
	HealthConditions             string                    `json:"health_conditions,omitempty"`
	InjuriesNotes                string                    `json:"injuries_notes,omitempty"`
	FitnessLevelID               *uuid.UUID                `json:"fitness_level_id,omitempty"`
	FitnessLevel                 *FitnessLevelResponse     `json:"fitness_level,omitempty"`
	CreatedAt                    time.Time                 `json:"created_at"`
	UpdatedAt                    time.Time                 `json:"updated_at"`
}

// ToResponse converts the model to a response DTO with unit conversions
// NOTE: This method requires the user's preferred_weight_unit to be passed in
// It will be called from the controller with the authenticated user's preferences
func (p *UserFitnessProfile) ToResponse(preferredWeightUnit string) UserFitnessProfileResponse {
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

	// Convert current weight to user's preferred unit
	var currentWeightValue *float64
	if p.CurrentWeightKg != nil {
		if preferredWeightUnit == "lb" {
			lbs := *p.CurrentWeightKg * 2.20462
			currentWeightValue = &lbs
		} else {
			currentWeightValue = p.CurrentWeightKg
		}
	}
	unit := preferredWeightUnit
	currentWeight := WeightOutput{
		WeightValue: currentWeightValue,
		WeightUnit:  &unit,
	}

	// Convert target weight to user's preferred unit
	var targetWeight *WeightOutput
	if p.TargetWeightKg != nil {
		var targetWeightValue *float64
		if preferredWeightUnit == "lb" {
			lbs := *p.TargetWeightKg * 2.20462
			targetWeightValue = &lbs
		} else {
			targetWeightValue = p.TargetWeightKg
		}
		targetWeight = &WeightOutput{
			WeightValue: targetWeightValue,
			WeightUnit:  &unit,
		}
	}

	// Convert fitness goals to response
	fitnessGoalResponses := make([]UserFitnessGoalResponse, len(p.FitnessGoals))
	for i, goal := range p.FitnessGoals {
		fitnessGoalResponses[i] = goal.ToResponse()
	}

	response := UserFitnessProfileResponse{
		ID:                           p.ID,
		UserID:                       p.UserID,
		DateOfBirth:                  p.DateOfBirth.Format("2006-01-02"),
		Age:                          age,
		Gender:                       p.Gender,
		HeightCm:                     p.HeightCm,
		HeightFtIn:                   heightFtIn,
		CurrentWeight:                currentWeight,
		FitnessGoals:                 fitnessGoalResponses,
		TargetWeight:                 targetWeight,
		TargetWeeklyWorkouts:         p.TargetWeeklyWorkouts,
		ActivityLevel:                p.ActivityLevel,
		TrainingLocations:            p.TrainingLocations,
		PreferredWorkoutDurationMins: p.PreferredWorkoutDurationMins,
		AvailableDays:                p.AvailableDays,
		HealthConditions:             p.HealthConditions,
		InjuriesNotes:                p.InjuriesNotes,
		FitnessLevelID:               p.FitnessLevelID,
		CreatedAt:                    p.CreatedAt,
		UpdatedAt:                    p.UpdatedAt,
	}

	if p.FitnessLevel != nil {
		levelResponse := p.FitnessLevel.ToResponse()
		response.FitnessLevel = &levelResponse
	}

	return response
}

// Helper function to format feet and inches
func formatFeetInches(feet, inches int) string {
	return fmt.Sprintf("%d'%d\"", feet, inches)
}
