package controllers

import (
	"errors"
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateWorkoutRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=200"`
	Description string `json:"description" binding:"omitempty,max=1000"`
	Visibility  string `json:"visibility" binding:"omitempty,oneof=public private friends"`
}

type UpdateWorkoutRequest struct {
	Title       string `json:"title" binding:"omitempty,min=1,max=200"`
	Description string `json:"description" binding:"omitempty,max=1000"`
	Visibility  string `json:"visibility" binding:"omitempty,oneof=public private friends"`
}

func CreateWorkout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

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
		Preload("Prescriptions", func(db *gorm.DB) *gorm.DB {
			return db.Order("group_order ASC, exercise_order ASC")
		}).
		Preload("Prescriptions.Exercise").
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
		Preload("Prescriptions", func(db *gorm.DB) *gorm.DB {
			return db.Order("group_order ASC, exercise_order ASC")
		}).
		Preload("Prescriptions.Exercise").
		Preload("Prescriptions.RPEValue").
		First(&workout).Error; err != nil {
		utils.NotFoundResponse(c, "Workout not found.")
		return
	}

	// Group prescriptions by group_id for response
	groupedPrescriptions := models.GroupPrescriptionsByGroupID(workout.Prescriptions)

	response := map[string]interface{}{
		"id":            workout.ID,
		"user_id":       workout.UserID,
		"title":         workout.Title,
		"description":   workout.Description,
		"visibility":    workout.Visibility,
		"created_at":    workout.CreatedAt,
		"updated_at":    workout.UpdatedAt,
		"prescriptions": groupedPrescriptions,
	}

	utils.SuccessResponse(c, "Workout fetched successfully.", response)
}

func UpdateWorkout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

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

// CreatePrescriptionGroup creates a new prescription group for a workout
func CreatePrescriptionGroup(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

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

	var req models.CreatePrescriptionGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Validate prescription type
	if !models.IsValidPrescriptionType(req.Type) {
		validationErrors := utils.ValidationErrors{
			"type": []string{"Invalid prescription type."},
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	// Check if workout exists and belongs to user
	var workout models.Workout
	if err := database.DB.Where("id = ? AND user_id = ?", workoutID, userUUID).First(&workout).Error; err != nil {
		utils.NotFoundResponse(c, "Workout not found.")
		return
	}

	// Generate group_id if not provided
	groupID := uuid.New()
	if req.GroupID != nil {
		groupID = *req.GroupID
	}

	// Start a transaction
	var createdPrescriptions []models.WorkoutPrescription
	var exerciseNotFound bool
	var validationError string
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// Create prescription rows for each exercise
		for _, exerciseReq := range req.Exercises {
			// Verify exercise exists
			var exercise models.Exercise
			if err := tx.Where("id = ?", exerciseReq.ExerciseID).First(&exercise).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					exerciseNotFound = true
				}
				return err
			}

			prescription := models.WorkoutPrescription{
				WorkoutID:       workoutID,
				ExerciseID:      exerciseReq.ExerciseID,
				GroupID:         groupID,
				Type:            req.Type,
				GroupOrder:      req.GroupOrder,
				GroupRounds:     req.GroupRounds,
				RestBetweenSets: req.RestBetweenSets,
				GroupName:       req.GroupName,
				GroupNotes:      req.GroupNotes,
				ExerciseOrder:   exerciseReq.ExerciseOrder,
				Sets:            exerciseReq.Sets,
				Reps:            exerciseReq.Reps,
				HoldSeconds:     exerciseReq.HoldSeconds,
				WeightKg:        exerciseReq.WeightKg,
				TargetWeightKg:  exerciseReq.TargetWeightKg,
				RPEValueID:      exerciseReq.RPEValueID,
				Notes:           exerciseReq.Notes,
			}

			if err := tx.Create(&prescription).Error; err != nil {
				// Check for validation errors from BeforeSave hook
				if err.Error() == "prescription cannot have both reps and hold_seconds" ||
					err.Error() == "prescription must have either reps or hold_seconds" ||
					err.Error() == "group_order must be at least 1" ||
					err.Error() == "exercise_order must be at least 1" ||
					err.Error() == "invalid prescription type" {
					validationError = err.Error()
				}
				return err
			}

			// Load the exercise details
			prescription.Exercise = exercise
			createdPrescriptions = append(createdPrescriptions, prescription)
		}

		return nil
	})

	if err != nil {
		if exerciseNotFound {
			validationErrors := utils.ValidationErrors{
				"exercise_id": []string{"Exercise not found."},
			}
			utils.ValidationErrorResponse(c, validationErrors)
			return
		}
		if validationError != "" {
			validationErrors := utils.ValidationErrors{
				"prescription": []string{validationError},
			}
			utils.ValidationErrorResponse(c, validationErrors)
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to create prescription group.")
		return
	}

	// Format response as grouped prescriptions
	groupedResponse := models.GroupPrescriptionsByGroupID(createdPrescriptions)
	if len(groupedResponse) > 0 {
		utils.CreatedResponse(c, "Prescription group created successfully.", groupedResponse[0])
	} else {
		utils.CreatedResponse(c, "Prescription group created successfully.", nil)
	}
}

// UpdatePrescriptionGroup updates an existing prescription group
func UpdatePrescriptionGroup(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

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

	groupID, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		validationErrors := utils.ValidationErrors{
			"group_id": []string{"Invalid group ID format."},
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	var req models.UpdatePrescriptionGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Check if workout exists and belongs to user
	var workout models.Workout
	if err := database.DB.Where("id = ? AND user_id = ?", workoutID, userUUID).First(&workout).Error; err != nil {
		utils.NotFoundResponse(c, "Workout not found.")
		return
	}

	// Check if group exists in this workout
	var existingPrescriptions []models.WorkoutPrescription
	if err := database.DB.Where("workout_id = ? AND group_id = ?", workoutID, groupID).
		Order("exercise_order ASC").
		Find(&existingPrescriptions).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch prescription group.")
		return
	}

	if len(existingPrescriptions) == 0 {
		utils.NotFoundResponse(c, "Prescription group not found.")
		return
	}

	// Start a transaction
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// Build updates for group-level fields
		groupUpdates := map[string]interface{}{}
		if req.Type != nil {
			if !models.IsValidPrescriptionType(*req.Type) {
				return gorm.ErrInvalidValue
			}
			groupUpdates["type"] = *req.Type
		}
		if req.GroupOrder != nil {
			groupUpdates["group_order"] = *req.GroupOrder
		}
		if req.GroupRounds != nil {
			groupUpdates["group_rounds"] = *req.GroupRounds
		}
		if req.RestBetweenSets != nil {
			groupUpdates["rest_between_sets"] = *req.RestBetweenSets
		}
		if req.GroupName != nil {
			groupUpdates["group_name"] = *req.GroupName
		}
		if req.GroupNotes != nil {
			groupUpdates["group_notes"] = *req.GroupNotes
		}

		// Update group-level fields on all rows
		if len(groupUpdates) > 0 {
			if err := tx.Model(&models.WorkoutPrescription{}).
				Where("workout_id = ? AND group_id = ?", workoutID, groupID).
				Updates(groupUpdates).Error; err != nil {
				return err
			}
		}

		// If exercises are provided, replace them
		if len(req.Exercises) > 0 {
			// Delete existing prescriptions
			if err := tx.Where("workout_id = ? AND group_id = ?", workoutID, groupID).
				Delete(&models.WorkoutPrescription{}).Error; err != nil {
				return err
			}

			// Get the current group-level values (use first existing or updated values)
			groupType := existingPrescriptions[0].Type
			groupOrder := existingPrescriptions[0].GroupOrder
			groupRounds := existingPrescriptions[0].GroupRounds
			restBetweenSets := existingPrescriptions[0].RestBetweenSets
			groupName := existingPrescriptions[0].GroupName
			groupNotes := existingPrescriptions[0].GroupNotes

			if req.Type != nil {
				groupType = *req.Type
			}
			if req.GroupOrder != nil {
				groupOrder = *req.GroupOrder
			}
			if req.GroupRounds != nil {
				groupRounds = req.GroupRounds
			}
			if req.RestBetweenSets != nil {
				restBetweenSets = req.RestBetweenSets
			}
			if req.GroupName != nil {
				groupName = req.GroupName
			}
			if req.GroupNotes != nil {
				groupNotes = req.GroupNotes
			}

			// Create new prescriptions
			for _, exerciseReq := range req.Exercises {
				// Verify exercise exists
				var exercise models.Exercise
				if err := tx.Where("id = ?", exerciseReq.ExerciseID).First(&exercise).Error; err != nil {
					return err
				}

				prescription := models.WorkoutPrescription{
					WorkoutID:       workoutID,
					ExerciseID:      exerciseReq.ExerciseID,
					GroupID:         groupID,
					Type:            groupType,
					GroupOrder:      groupOrder,
					GroupRounds:     groupRounds,
					RestBetweenSets: restBetweenSets,
					GroupName:       groupName,
					GroupNotes:      groupNotes,
					ExerciseOrder:   exerciseReq.ExerciseOrder,
					Sets:            exerciseReq.Sets,
					Reps:            exerciseReq.Reps,
					HoldSeconds:     exerciseReq.HoldSeconds,
					WeightKg:        exerciseReq.WeightKg,
					TargetWeightKg:  exerciseReq.TargetWeightKg,
					RPEValueID:      exerciseReq.RPEValueID,
					Notes:           exerciseReq.Notes,
				}

				if err := tx.Create(&prescription).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			validationErrors := utils.ValidationErrors{
				"exercise_id": []string{"Exercise not found."},
			}
			utils.ValidationErrorResponse(c, validationErrors)
			return
		}
		if errors.Is(err, gorm.ErrInvalidValue) {
			validationErrors := utils.ValidationErrors{
				"type": []string{"Invalid prescription type."},
			}
			utils.ValidationErrorResponse(c, validationErrors)
			return
		}
		// Check for validation errors from BeforeSave hook
		errMsg := err.Error()
		if errMsg == "prescription cannot have both reps and hold_seconds" ||
			errMsg == "prescription must have either reps or hold_seconds" ||
			errMsg == "group_order must be at least 1" ||
			errMsg == "exercise_order must be at least 1" ||
			errMsg == "invalid prescription type" {
			validationErrors := utils.ValidationErrors{
				"validation": []string{errMsg},
			}
			utils.ValidationErrorResponse(c, validationErrors)
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to update prescription group.")
		return
	}

	// Fetch updated prescriptions
	var updatedPrescriptions []models.WorkoutPrescription
	if err := database.DB.Where("workout_id = ? AND group_id = ?", workoutID, groupID).
		Preload("Exercise").
		Preload("RPEValue").
		Order("exercise_order ASC").
		Find(&updatedPrescriptions).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch updated prescription group.")
		return
	}

	groupedResponse := models.GroupPrescriptionsByGroupID(updatedPrescriptions)
	if len(groupedResponse) > 0 {
		utils.SuccessResponse(c, "Prescription group updated successfully.", groupedResponse[0])
	} else {
		utils.SuccessResponse(c, "Prescription group updated successfully.", nil)
	}
}

// DeletePrescriptionGroup deletes a prescription group from a workout
func DeletePrescriptionGroup(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

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

	groupID, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		validationErrors := utils.ValidationErrors{
			"group_id": []string{"Invalid group ID format."},
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	// Check if workout exists and belongs to user
	var workout models.Workout
	if err := database.DB.Where("id = ? AND user_id = ?", workoutID, userUUID).First(&workout).Error; err != nil {
		utils.NotFoundResponse(c, "Workout not found.")
		return
	}

	// Delete all prescriptions in the group
	result := database.DB.Where("workout_id = ? AND group_id = ?", workoutID, groupID).
		Delete(&models.WorkoutPrescription{})
	if result.Error != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete prescription group.")
		return
	}

	if result.RowsAffected == 0 {
		utils.NotFoundResponse(c, "Prescription group not found.")
		return
	}

	utils.DeletedResponse(c, "Prescription group deleted successfully.")
}

// ReorderPrescriptionGroups reorders the groups within a workout
func ReorderPrescriptionGroups(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

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

	var req models.ReorderPrescriptionGroupsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Check if workout exists and belongs to user
	var workout models.Workout
	if err := database.DB.Where("id = ? AND user_id = ?", workoutID, userUUID).First(&workout).Error; err != nil {
		utils.NotFoundResponse(c, "Workout not found.")
		return
	}

	// Start a transaction
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		for _, groupOrder := range req.GroupOrders {
			if err := tx.Model(&models.WorkoutPrescription{}).
				Where("workout_id = ? AND group_id = ?", workoutID, groupOrder.GroupID).
				Update("group_order", groupOrder.GroupOrder).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to reorder prescription groups.")
		return
	}

	utils.SuccessResponse(c, "Prescription groups reordered successfully.", nil)
}

// GetWorkoutPrescriptions returns all prescriptions for a workout grouped by group_id
func GetWorkoutPrescriptions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

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

	// Check if workout exists and belongs to user
	var workout models.Workout
	if err := database.DB.Where("id = ? AND user_id = ?", workoutID, userUUID).First(&workout).Error; err != nil {
		utils.NotFoundResponse(c, "Workout not found.")
		return
	}

	var prescriptions []models.WorkoutPrescription
	if err := database.DB.Where("workout_id = ?", workoutID).
		Preload("Exercise").
		Preload("RPEValue").
		Order("group_order ASC, exercise_order ASC").
		Find(&prescriptions).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch workout prescriptions.")
		return
	}

	groupedResponse := models.GroupPrescriptionsByGroupID(prescriptions)
	utils.SuccessResponse(c, "Workout prescriptions fetched successfully.", groupedResponse)
}

// AddExerciseToPrescriptionGroup adds an exercise to an existing prescription group
func AddExerciseToPrescriptionGroup(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

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

	groupID, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		validationErrors := utils.ValidationErrors{
			"group_id": []string{"Invalid group ID format."},
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	var req models.AddExerciseToPrescriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Check if workout exists and belongs to user
	var workout models.Workout
	if err := database.DB.Where("id = ? AND user_id = ?", workoutID, userUUID).First(&workout).Error; err != nil {
		utils.NotFoundResponse(c, "Workout not found.")
		return
	}

	// Get existing prescriptions in the group
	var existingPrescriptions []models.WorkoutPrescription
	if err := database.DB.Where("workout_id = ? AND group_id = ?", workoutID, groupID).
		Order("exercise_order ASC").
		Find(&existingPrescriptions).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch prescription group.")
		return
	}

	if len(existingPrescriptions) == 0 {
		utils.NotFoundResponse(c, "Prescription group not found.")
		return
	}

	// Verify exercise exists
	var exercise models.Exercise
	if err := database.DB.Where("id = ?", req.ExerciseID).First(&exercise).Error; err != nil {
		utils.NotFoundResponse(c, "Exercise not found.")
		return
	}

	// Get next exercise order
	nextOrder, err := models.GetNextExerciseOrder(database.DB, groupID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to determine exercise order.")
		return
	}

	// Use group-level values from existing prescription
	firstPrescription := existingPrescriptions[0]

	prescription := models.WorkoutPrescription{
		WorkoutID:       workoutID,
		ExerciseID:      req.ExerciseID,
		GroupID:         groupID,
		Type:            firstPrescription.Type,
		GroupOrder:      firstPrescription.GroupOrder,
		GroupRounds:     firstPrescription.GroupRounds,
		RestBetweenSets: firstPrescription.RestBetweenSets,
		GroupName:       firstPrescription.GroupName,
		GroupNotes:      firstPrescription.GroupNotes,
		ExerciseOrder:   nextOrder,
		Sets:            req.Sets,
		Reps:            req.Reps,
		HoldSeconds:     req.HoldSeconds,
		WeightKg:        req.WeightKg,
		TargetWeightKg:  req.TargetWeightKg,
		RPEValueID:      req.RPEValueID,
		Notes:           req.Notes,
	}

	if err := database.DB.Create(&prescription).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to add exercise to prescription group.")
		return
	}

	// Load exercise details
	prescription.Exercise = exercise

	utils.CreatedResponse(c, "Exercise added to prescription group successfully.", prescription)
}

// DuplicateWorkout creates a copy of an existing workout with all its prescriptions
func DuplicateWorkout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

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

	// Fetch the original workout with all prescriptions
	var originalWorkout models.Workout
	if err := database.DB.Where("id = ? AND user_id = ?", workoutID, userUUID).
		Preload("Prescriptions").
		First(&originalWorkout).Error; err != nil {
		utils.NotFoundResponse(c, "Workout not found.")
		return
	}

	var newWorkout models.Workout

	// Start a transaction
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// Create a new workout
		newWorkout = models.Workout{
			UserID:      userUUID,
			Title:       originalWorkout.Title + " (Copy)",
			Description: originalWorkout.Description,
			Visibility:  originalWorkout.Visibility,
		}

		if err := tx.Create(&newWorkout).Error; err != nil {
			return err
		}

		// Map old group IDs to new group IDs
		groupIDMapping := make(map[uuid.UUID]uuid.UUID)

		// Duplicate prescriptions
		for _, prescription := range originalWorkout.Prescriptions {
			// Get or create new group ID
			newGroupID, exists := groupIDMapping[prescription.GroupID]
			if !exists {
				newGroupID = uuid.New()
				groupIDMapping[prescription.GroupID] = newGroupID
			}

			newPrescription := models.WorkoutPrescription{
				WorkoutID:       newWorkout.ID,
				ExerciseID:      prescription.ExerciseID,
				RPEValueID:      prescription.RPEValueID,
				GroupID:         newGroupID,
				Type:            prescription.Type,
				GroupOrder:      prescription.GroupOrder,
				GroupRounds:     prescription.GroupRounds,
				RestBetweenSets: prescription.RestBetweenSets,
				GroupName:       prescription.GroupName,
				GroupNotes:      prescription.GroupNotes,
				ExerciseOrder:   prescription.ExerciseOrder,
				Sets:            prescription.Sets,
				Reps:            prescription.Reps,
				HoldSeconds:     prescription.HoldSeconds,
				WeightKg:        prescription.WeightKg,
				TargetWeightKg:  prescription.TargetWeightKg,
				Notes:           prescription.Notes,
			}

			if err := tx.Create(&newPrescription).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to duplicate workout.")
		return
	}

	// Load the new workout with all relations
	if err := database.DB.Where("id = ?", newWorkout.ID).
		Preload("Prescriptions", func(db *gorm.DB) *gorm.DB {
			return db.Order("group_order ASC, exercise_order ASC")
		}).
		Preload("Prescriptions.Exercise").
		First(&newWorkout).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to load duplicated workout.")
		return
	}

	// Group prescriptions for response
	groupedPrescriptions := models.GroupPrescriptionsByGroupID(newWorkout.Prescriptions)

	response := map[string]interface{}{
		"id":            newWorkout.ID,
		"user_id":       newWorkout.UserID,
		"title":         newWorkout.Title,
		"description":   newWorkout.Description,
		"visibility":    newWorkout.Visibility,
		"created_at":    newWorkout.CreatedAt,
		"updated_at":    newWorkout.UpdatedAt,
		"prescriptions": groupedPrescriptions,
	}

	utils.CreatedResponse(c, "Workout duplicated successfully.", response)
}
