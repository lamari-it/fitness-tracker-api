package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// CreateUserFitnessProfile creates a fitness profile for the authenticated user
func CreateUserFitnessProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type")
		return
	}

	// Check if user already has a fitness profile
	var existingProfile models.UserFitnessProfile
	if err := database.DB.Where("user_id = ?", userID).First(&existingProfile).Error; err == nil {
		utils.ConflictResponse(c, "Fitness profile already exists for this user")
		return
	}

	var req models.CreateUserFitnessProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Parse date of birth
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid date format. Use YYYY-MM-DD", nil)
		return
	}

	// Validate that all fitness goal IDs exist
	var fitnessGoals []models.FitnessGoal
	if err := database.DB.Where("id IN ?", req.FitnessGoalIDs).Find(&fitnessGoals).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to validate fitness goals")
		return
	}
	if len(fitnessGoals) != len(req.FitnessGoalIDs) {
		utils.BadRequestResponse(c, "One or more fitness goal IDs are invalid", nil)
		return
	}

	// Create profile with defaults for optional fields
	profile := models.UserFitnessProfile{
		UserID:          userID,
		DateOfBirth:     dob,
		Gender:          req.Gender,
		HeightCm:        req.HeightCm,
		CurrentWeightKg: req.CurrentWeightKg,
	}

	// Set optional fields with defaults
	if req.PreferredUnitSystem != "" {
		profile.PreferredUnitSystem = req.PreferredUnitSystem
	} else {
		profile.PreferredUnitSystem = "metric"
	}

	if req.TargetWeightKg != nil {
		profile.TargetWeightKg = req.TargetWeightKg
	}

	if req.TargetWeeklyWorkouts > 0 {
		profile.TargetWeeklyWorkouts = req.TargetWeeklyWorkouts
	} else {
		profile.TargetWeeklyWorkouts = 3
	}

	if req.ActivityLevel != "" {
		profile.ActivityLevel = req.ActivityLevel
	} else {
		profile.ActivityLevel = "moderate"
	}

	if len(req.TrainingLocations) > 0 {
		profile.TrainingLocations = pq.StringArray(req.TrainingLocations)
	} else {
		profile.TrainingLocations = pq.StringArray{"gym"}
	}

	if req.PreferredWorkoutDurationMins > 0 {
		profile.PreferredWorkoutDurationMins = req.PreferredWorkoutDurationMins
	} else {
		profile.PreferredWorkoutDurationMins = 45
	}

	if len(req.AvailableDays) > 0 {
		profile.AvailableDays = pq.StringArray(req.AvailableDays)
	} else {
		profile.AvailableDays = pq.StringArray{"monday", "wednesday", "friday"}
	}

	profile.HealthConditions = req.HealthConditions
	profile.InjuriesNotes = req.InjuriesNotes

	// Use transaction for profile creation and goal associations
	tx := database.DB.Begin()

	if err := tx.Create(&profile).Error; err != nil {
		tx.Rollback()
		utils.InternalServerErrorResponse(c, "Failed to create fitness profile")
		return
	}

	// Create fitness goal associations
	for i, goalID := range req.FitnessGoalIDs {
		userGoal := models.UserFitnessGoal{
			UserFitnessProfileID: profile.ID,
			FitnessGoalID:        goalID,
			Priority:             i + 1, // First goal is primary (priority 1)
		}
		if err := tx.Create(&userGoal).Error; err != nil {
			tx.Rollback()
			utils.InternalServerErrorResponse(c, "Failed to associate fitness goals")
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create fitness profile")
		return
	}

	// Preload fitness goals for response
	database.DB.Preload("FitnessGoals.FitnessGoal").First(&profile, "id = ?", profile.ID)

	utils.CreatedResponse(c, "Fitness profile created successfully", profile.ToResponse())
}

// GetUserFitnessProfile retrieves the fitness profile for the authenticated user
func GetUserFitnessProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type")
		return
	}

	var profile models.UserFitnessProfile
	if err := database.DB.Preload("FitnessGoals.FitnessGoal").Where("user_id = ?", userID).First(&profile).Error; err != nil {
		utils.NotFoundResponse(c, "Fitness profile not found")
		return
	}

	utils.SuccessResponse(c, "Fitness profile retrieved successfully", profile.ToResponse())
}

// UpdateUserFitnessProfile updates the fitness profile for the authenticated user
func UpdateUserFitnessProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type")
		return
	}

	var profile models.UserFitnessProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		utils.NotFoundResponse(c, "Fitness profile not found")
		return
	}

	var req models.UpdateUserFitnessProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Update fields only if provided
	if req.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			utils.BadRequestResponse(c, "Invalid date format. Use YYYY-MM-DD", nil)
			return
		}
		profile.DateOfBirth = dob
	}

	if req.Gender != "" {
		profile.Gender = req.Gender
	}

	if req.HeightCm > 0 {
		profile.HeightCm = req.HeightCm
	}

	if req.CurrentWeightKg > 0 {
		profile.CurrentWeightKg = req.CurrentWeightKg
	}

	if req.PreferredUnitSystem != "" {
		profile.PreferredUnitSystem = req.PreferredUnitSystem
	}

	if req.TargetWeightKg != nil {
		profile.TargetWeightKg = req.TargetWeightKg
	}

	if req.TargetWeeklyWorkouts > 0 {
		profile.TargetWeeklyWorkouts = req.TargetWeeklyWorkouts
	}

	if req.ActivityLevel != "" {
		profile.ActivityLevel = req.ActivityLevel
	}

	if len(req.TrainingLocations) > 0 {
		profile.TrainingLocations = pq.StringArray(req.TrainingLocations)
	}

	if req.PreferredWorkoutDurationMins > 0 {
		profile.PreferredWorkoutDurationMins = req.PreferredWorkoutDurationMins
	}

	if len(req.AvailableDays) > 0 {
		profile.AvailableDays = pq.StringArray(req.AvailableDays)
	}

	if req.HealthConditions != "" {
		profile.HealthConditions = req.HealthConditions
	}

	if req.InjuriesNotes != "" {
		profile.InjuriesNotes = req.InjuriesNotes
	}

	// Use transaction if updating fitness goals
	if len(req.FitnessGoalIDs) > 0 {
		// Validate that all fitness goal IDs exist
		var fitnessGoals []models.FitnessGoal
		if err := database.DB.Where("id IN ?", req.FitnessGoalIDs).Find(&fitnessGoals).Error; err != nil {
			utils.InternalServerErrorResponse(c, "Failed to validate fitness goals")
			return
		}
		if len(fitnessGoals) != len(req.FitnessGoalIDs) {
			utils.BadRequestResponse(c, "One or more fitness goal IDs are invalid", nil)
			return
		}

		tx := database.DB.Begin()

		// Save profile updates
		if err := tx.Save(&profile).Error; err != nil {
			tx.Rollback()
			utils.InternalServerErrorResponse(c, "Failed to update fitness profile")
			return
		}

		// Delete existing goal associations
		if err := tx.Where("user_fitness_profile_id = ?", profile.ID).Delete(&models.UserFitnessGoal{}).Error; err != nil {
			tx.Rollback()
			utils.InternalServerErrorResponse(c, "Failed to update fitness goals")
			return
		}

		// Create new goal associations
		for i, goalID := range req.FitnessGoalIDs {
			userGoal := models.UserFitnessGoal{
				UserFitnessProfileID: profile.ID,
				FitnessGoalID:        goalID,
				Priority:             i + 1,
			}
			if err := tx.Create(&userGoal).Error; err != nil {
				tx.Rollback()
				utils.InternalServerErrorResponse(c, "Failed to associate fitness goals")
				return
			}
		}

		if err := tx.Commit().Error; err != nil {
			utils.InternalServerErrorResponse(c, "Failed to update fitness profile")
			return
		}
	} else {
		if err := database.DB.Save(&profile).Error; err != nil {
			utils.InternalServerErrorResponse(c, "Failed to update fitness profile")
			return
		}
	}

	// Preload fitness goals for response
	database.DB.Preload("FitnessGoals.FitnessGoal").First(&profile, "id = ?", profile.ID)

	utils.SuccessResponse(c, "Fitness profile updated successfully", profile.ToResponse())
}

// DeleteUserFitnessProfile deletes the fitness profile for the authenticated user
func DeleteUserFitnessProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type")
		return
	}

	var profile models.UserFitnessProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		utils.NotFoundResponse(c, "Fitness profile not found")
		return
	}

	if err := database.DB.Delete(&profile).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete fitness profile")
		return
	}

	utils.NoContentResponse(c)
}

// LogWeight logs a weight update for the user (updates current weight)
func LogWeight(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type")
		return
	}

	var profile models.UserFitnessProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		utils.NotFoundResponse(c, "Fitness profile not found")
		return
	}

	var req struct {
		WeightKg float64 `json:"weight_kg" binding:"required,gt=20,lt=500"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	profile.CurrentWeightKg = req.WeightKg

	if err := database.DB.Save(&profile).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to log weight")
		return
	}

	// Also create a weight log entry for historical tracking
	weightLog := models.WeightLog{
		UserID:   userID,
		WeightKg: req.WeightKg,
	}
	database.DB.Create(&weightLog)

	utils.SuccessResponse(c, "Weight logged successfully", profile.ToResponse())
}
