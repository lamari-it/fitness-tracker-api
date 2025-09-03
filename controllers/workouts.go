package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateWorkoutRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Visibility  string `json:"visibility"`
}

type UpdateWorkoutRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Visibility  string `json:"visibility"`
}

type AddWorkoutExerciseRequest struct {
	ExerciseID        uuid.UUID `json:"exercise_id" binding:"required"`
	SetGroupID        uuid.UUID `json:"set_group_id" binding:"required"`
	OrderNumber       int       `json:"order_number" binding:"required"`
	TargetSets        int       `json:"target_sets"`
	TargetReps        int       `json:"target_reps"`
	TargetWeight      float64   `json:"target_weight"`
	TargetRestSec     int       `json:"target_rest_sec"`
	Prescription      string    `json:"prescription"` // reps or time
	TargetDurationSec int       `json:"target_duration_sec"`
}

func CreateWorkout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	// Ensure userID is properly typed as UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type.")
		return
	}

	var req CreateWorkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	workout := models.Workout{
		UserID:      userUUID,
		Title:       req.Title,
		Description: req.Description,
		Visibility:  req.Visibility,
	}

	if workout.Visibility == "" {
		workout.Visibility = "private"
	}

	if err := database.DB.Create(&workout).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create workout.")
		return
	}

	utils.CreatedResponse(c, "Workout created successfully.", workout)
}

func GetUserWorkouts(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	// Ensure userID is properly typed as UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type.")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	offset := (page - 1) * limit

	var workouts []models.Workout
	var total int64

	if err := database.DB.Model(&models.Workout{}).Where("user_id = ?", userUUID).Count(&total).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to count workouts.")
		return
	}

	if err := database.DB.Where("user_id = ?", userUUID).
		Preload("SetGroups").
		Preload("WorkoutExercises.Exercise").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&workouts).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch workouts.")
		return
	}

	utils.PaginatedResponse(c, "Workouts fetched successfully.", workouts, page, limit, int(total))
}

func GetWorkout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	// Ensure userID is properly typed as UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type.")
		return
	}

	workoutID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		validationErrors := utils.ValidationErrors{
			"id": []string{"Invalid workout ID format."},
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	var workout models.Workout
	if err := database.DB.Where("id = ? AND user_id = ?", workoutID, userUUID).
		Preload("SetGroups").
		Preload("WorkoutExercises.Exercise").
		Preload("WorkoutExercises.SetGroup").
		First(&workout).Error; err != nil {
		utils.NotFoundResponse(c, "Workout not found.")
		return
	}

	utils.SuccessResponse(c, "Workout fetched successfully.", workout)
}

func UpdateWorkout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	// Ensure userID is properly typed as UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type.")
		return
	}

	workoutID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		validationErrors := utils.ValidationErrors{
			"id": []string{"Invalid workout ID format."},
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	var req UpdateWorkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	var workout models.Workout
	if err := database.DB.Where("id = ? AND user_id = ?", workoutID, userUUID).First(&workout).Error; err != nil {
		utils.NotFoundResponse(c, "Workout not found.")
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Visibility != "" {
		updates["visibility"] = req.Visibility
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&workout).Updates(updates).Error; err != nil {
			utils.InternalServerErrorResponse(c, "Failed to update workout.")
			return
		}
	}

	database.DB.Where("id = ?", workoutID).First(&workout)

	utils.SuccessResponse(c, "Workout updated successfully.", workout)
}

func DeleteWorkout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	// Ensure userID is properly typed as UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type.")
		return
	}

	workoutID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		validationErrors := utils.ValidationErrors{
			"id": []string{"Invalid workout ID format."},
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	var workout models.Workout
	if err := database.DB.Where("id = ? AND user_id = ?", workoutID, userUUID).First(&workout).Error; err != nil {
		utils.NotFoundResponse(c, "Workout not found.")
		return
	}

	if err := database.DB.Delete(&workout).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete workout.")
		return
	}

	utils.DeletedResponse(c, "Workout deleted successfully.")
}

func AddExerciseToWorkout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	workoutID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		validationErrors := utils.ValidationErrors{
			"id": []string{"Invalid workout ID format."},
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	var req AddWorkoutExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Ensure userID is properly typed as UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type.")
		return
	}

	// Check if workout exists and belongs to user
	var workout models.Workout
	if err := database.DB.Where("id = ? AND user_id = ?", workoutID, userUUID).First(&workout).Error; err != nil {
		utils.NotFoundResponse(c, "Workout not found.")
		return
	}

	// Check if exercise exists
	var exercise models.Exercise
	if err := database.DB.Where("id = ?", req.ExerciseID).First(&exercise).Error; err != nil {
		utils.NotFoundResponse(c, "Exercise not found.")
		return
	}

	// Verify set group exists and belongs to this workout
	var setGroup models.SetGroup
	if err := database.DB.Where("id = ? AND workout_id = ?", req.SetGroupID, workoutID).First(&setGroup).Error; err != nil {
		utils.NotFoundResponse(c, "Set group not found in this workout.")
		return
	}

	// Set default prescription if not provided
	prescription := req.Prescription
	if prescription == "" {
		prescription = "reps"
	}

	// Create workout exercise
	workoutExercise := models.WorkoutExercise{
		WorkoutID:         workoutID,
		ExerciseID:        req.ExerciseID,
		SetGroupID:        req.SetGroupID,
		OrderNumber:       req.OrderNumber,
		TargetSets:        req.TargetSets,
		TargetReps:        req.TargetReps,
		TargetWeight:      req.TargetWeight,
		TargetRestSec:     req.TargetRestSec,
		Prescription:      prescription,
		TargetDurationSec: req.TargetDurationSec,
	}

	if err := database.DB.Create(&workoutExercise).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to add exercise to workout.")
		return
	}

	// Load the exercise details
	database.DB.Preload("Exercise").First(&workoutExercise, workoutExercise.ID)

	utils.CreatedResponse(c, "Exercise added to workout successfully.", workoutExercise)
}

func RemoveExerciseFromWorkout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	workoutID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		validationErrors := utils.ValidationErrors{
			"id": []string{"Invalid workout ID format."},
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	exerciseID, err := uuid.Parse(c.Param("exercise_id"))
	if err != nil {
		validationErrors := utils.ValidationErrors{
			"exercise_id": []string{"Invalid exercise ID format."},
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	// Ensure userID is properly typed as UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type.")
		return
	}

	// Check if workout exists and belongs to user
	var workout models.Workout
	if err := database.DB.Where("id = ? AND user_id = ?", workoutID, userUUID).First(&workout).Error; err != nil {
		utils.NotFoundResponse(c, "Workout not found.")
		return
	}

	// Delete the workout exercise
	result := database.DB.Where("workout_id = ? AND id = ?", workoutID, exerciseID).Delete(&models.WorkoutExercise{})
	if result.Error != nil {
		utils.InternalServerErrorResponse(c, "Failed to remove exercise from workout.")
		return
	}

	if result.RowsAffected == 0 {
		utils.NotFoundResponse(c, "Exercise not found in this workout.")
		return
	}

	utils.DeletedResponse(c, "Exercise removed from workout successfully.")
}

func GetWorkoutExercises(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	workoutID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		validationErrors := utils.ValidationErrors{
			"id": []string{"Invalid workout ID format."},
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	// Ensure userID is properly typed as UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type.")
		return
	}

	// Check if workout exists and belongs to user
	var workout models.Workout
	if err := database.DB.Where("id = ? AND user_id = ?", workoutID, userUUID).First(&workout).Error; err != nil {
		utils.NotFoundResponse(c, "Workout not found.")
		return
	}

	var exercises []models.WorkoutExercise
	if err := database.DB.Where("workout_id = ?", workoutID).
		Preload("Exercise").
		Preload("SetGroup").
		Order("order_number").
		Find(&exercises).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch workout exercises.")
		return
	}

	utils.SuccessResponse(c, "Workout exercises fetched successfully.", exercises)
}

func DuplicateWorkout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	workoutID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		validationErrors := utils.ValidationErrors{
			"id": []string{"Invalid workout ID format."},
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	// Ensure userID is properly typed as UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type.")
		return
	}

	// Fetch the original workout with all related data
	var originalWorkout models.Workout
	if err := database.DB.Where("id = ? AND user_id = ?", workoutID, userUUID).
		Preload("SetGroups").
		Preload("WorkoutExercises").
		First(&originalWorkout).Error; err != nil {
		utils.NotFoundResponse(c, "Workout not found.")
		return
	}

	// Start a transaction
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// Create a new workout
		newWorkout := models.Workout{
			UserID:      userUUID,
			Title:       originalWorkout.Title + " (Copy)",
			Description: originalWorkout.Description,
			Visibility:  originalWorkout.Visibility,
		}

		if err := tx.Create(&newWorkout).Error; err != nil {
			return err
		}

		// Duplicate set groups if any
		setGroupMapping := make(map[uuid.UUID]uuid.UUID) // Map old set group IDs to new ones
		for _, setGroup := range originalWorkout.SetGroups {
			newSetGroup := models.SetGroup{
				WorkoutID:       newWorkout.ID,
				GroupType:       setGroup.GroupType,
				Name:            setGroup.Name,
				Notes:           setGroup.Notes,
				OrderNumber:     setGroup.OrderNumber,
				RestBetweenSets: setGroup.RestBetweenSets,
				Rounds:          setGroup.Rounds,
			}
			if err := tx.Create(&newSetGroup).Error; err != nil {
				return err
			}
			setGroupMapping[setGroup.ID] = newSetGroup.ID
		}

		// Duplicate workout exercises
		for _, exercise := range originalWorkout.WorkoutExercises {
			newExercise := models.WorkoutExercise{
				WorkoutID:         newWorkout.ID,
				ExerciseID:        exercise.ExerciseID,
				OrderNumber:       exercise.OrderNumber,
				TargetSets:        exercise.TargetSets,
				TargetReps:        exercise.TargetReps,
				TargetWeight:      exercise.TargetWeight,
				TargetRestSec:     exercise.TargetRestSec,
				Prescription:      exercise.Prescription,
				TargetDurationSec: exercise.TargetDurationSec,
			}

			// Map the set group ID
			if newSetGroupID, ok := setGroupMapping[exercise.SetGroupID]; ok {
				newExercise.SetGroupID = newSetGroupID
			} else {
				// If mapping doesn't exist, use the original (this shouldn't happen normally)
				newExercise.SetGroupID = exercise.SetGroupID
			}

			if err := tx.Create(&newExercise).Error; err != nil {
				return err
			}
		}

		// Load the new workout with all relations
		if err := tx.Where("id = ?", newWorkout.ID).
			Preload("SetGroups").
			Preload("WorkoutExercises.Exercise").
			First(&newWorkout).Error; err != nil {
			return err
		}

		// Store the new workout for response
		originalWorkout = newWorkout
		return nil
	})

	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to duplicate workout.")
		return
	}

	utils.CreatedResponse(c, "Workout duplicated successfully.", originalWorkout)
}
