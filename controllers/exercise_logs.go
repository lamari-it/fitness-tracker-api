package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Request DTOs
type CreateExerciseLogRequest struct {
	SessionID        uuid.UUID  `json:"session_id" binding:"required"`
	PrescriptionID   *uuid.UUID `json:"prescription_id"`
	ExerciseID       uuid.UUID  `json:"exercise_id" binding:"required"`
	OrderNumber      int        `json:"order_number" binding:"required"`
	Notes            string     `json:"notes"`
	DifficultyRating int        `json:"difficulty_rating"`
	DifficultyType   string     `json:"difficulty_type"`
}

type UpdateExerciseLogRequest struct {
	Notes            string `json:"notes"`
	DifficultyRating int    `json:"difficulty_rating"`
	DifficultyType   string `json:"difficulty_type"`
}

// ExerciseLogResponse is the response DTO for exercise logs
type ExerciseLogResponse struct {
	ID               uuid.UUID               `json:"id"`
	SessionID        uuid.UUID               `json:"session_id"`
	PrescriptionID   *uuid.UUID              `json:"prescription_id,omitempty"`
	GroupID          *uuid.UUID              `json:"group_id,omitempty"`
	GroupName        *string                 `json:"group_name,omitempty"`
	GroupType        models.PrescriptionType `json:"group_type,omitempty"`
	ExerciseID       uuid.UUID               `json:"exercise_id"`
	ExerciseName     string                  `json:"exercise_name"`
	OrderNumber      int                     `json:"order_number"`
	Notes            string                  `json:"notes"`
	DifficultyRating int                     `json:"difficulty_rating"`
	DifficultyType   string                  `json:"difficulty_type"`
	SetLogs          []SetLogResponse        `json:"set_logs,omitempty"`
}

// SetLogResponse is the response DTO for set logs
type SetLogResponse struct {
	ID                uuid.UUID  `json:"id"`
	SetNumber         int        `json:"set_number"`
	WeightKg          float64    `json:"weight_kg"`           // Canonical weight in kg
	WeightDisplay     float64    `json:"weight_display"`      // Weight in user's preferred unit
	WeightDisplayUnit string     `json:"weight_display_unit"` // User's preferred unit (kg/lb)
	InputWeight       float64    `json:"input_weight"`        // Original input value
	InputWeightUnit   string     `json:"input_weight_unit"`   // Original input unit
	Reps              int        `json:"reps"`
	RestAfterSec      int        `json:"rest_after_sec"`
	Tempo             string     `json:"tempo"`
	RPE               float64    `json:"rpe"`
	RPEValueID        *uuid.UUID `json:"rpe_value_id,omitempty"`
}

// buildSetLogResponse creates a SetLogResponse with weight conversion
func buildSetLogResponse(sl models.SetLog, preferredUnit string) SetLogResponse {
	return SetLogResponse{
		ID:                sl.ID,
		SetNumber:         sl.SetNumber,
		WeightKg:          sl.Weight,
		WeightDisplay:     utils.ConvertFromKg(sl.Weight, preferredUnit),
		WeightDisplayUnit: preferredUnit,
		InputWeight:       sl.InputWeight,
		InputWeightUnit:   sl.InputWeightUnit,
		Reps:              sl.Reps,
		RestAfterSec:      sl.RestAfterSec,
		Tempo:             sl.Tempo,
		RPE:               sl.RPE,
		RPEValueID:        sl.RPEValueID,
	}
}

// Helper to check if user can access a session
func canAccessSession(authUserID uuid.UUID, session *models.WorkoutSession) bool {
	isOwner := session.UserID == authUserID
	isCreator := session.CreatedByID != nil && *session.CreatedByID == authUserID
	return isOwner || isCreator
}

// getUserPreferredWeightUnit fetches the user's preferred weight unit from their fitness profile
// Returns "kg" as default if profile not found
func getUserPreferredWeightUnit(userID uuid.UUID) string {
	var profile models.UserFitnessProfile
	if err := database.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return "kg" // Default to kg if no profile
	}
	return utils.GetUserPreferredWeightUnit(profile.PreferredWeightUnit)
}

// CreateExerciseLog creates a new exercise log in a session
func CreateExerciseLog(c *gin.Context) {
	authUserIDVal, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	authUserID, ok := authUserIDVal.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type")
		return
	}

	var req CreateExerciseLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Verify session exists and user can access it
	var session models.WorkoutSession
	if err := database.DB.First(&session, "id = ?", req.SessionID).Error; err != nil {
		utils.NotFoundResponse(c, "Workout session not found")
		return
	}

	if !canAccessSession(authUserID, &session) {
		utils.ForbiddenResponse(c, "Not authorized to log exercises in this session")
		return
	}

	// Verify exercise exists
	var exercise models.Exercise
	if err := database.DB.First(&exercise, "id = ?", req.ExerciseID).Error; err != nil {
		utils.NotFoundResponse(c, "Exercise not found")
		return
	}

	exerciseLog := models.ExerciseLog{
		SessionID:        req.SessionID,
		PrescriptionID:   req.PrescriptionID,
		ExerciseID:       req.ExerciseID,
		OrderNumber:      req.OrderNumber,
		Notes:            req.Notes,
		DifficultyRating: req.DifficultyRating,
		DifficultyType:   req.DifficultyType,
	}

	if err := database.DB.Create(&exerciseLog).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create exercise log")
		return
	}

	// Reload with relationships
	database.DB.Preload("Exercise").Preload("Prescription").Preload("SetLogs").First(&exerciseLog, "id = ?", exerciseLog.ID)

	response := ExerciseLogResponse{
		ID:               exerciseLog.ID,
		SessionID:        exerciseLog.SessionID,
		PrescriptionID:   exerciseLog.PrescriptionID,
		ExerciseID:       exerciseLog.ExerciseID,
		ExerciseName:     exerciseLog.Exercise.Name,
		OrderNumber:      exerciseLog.OrderNumber,
		Notes:            exerciseLog.Notes,
		DifficultyRating: exerciseLog.DifficultyRating,
		DifficultyType:   exerciseLog.DifficultyType,
		SetLogs:          []SetLogResponse{},
	}

	// Add prescription group info if available
	if exerciseLog.Prescription != nil {
		response.GroupID = &exerciseLog.Prescription.GroupID
		response.GroupName = exerciseLog.Prescription.GroupName
		response.GroupType = exerciseLog.Prescription.Type
	}

	utils.CreatedResponse(c, "Exercise log created successfully", response)
}

// GetExerciseLog retrieves a single exercise log
func GetExerciseLog(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	exerciseLogID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid UUID format", nil)
		return
	}

	authUserIDVal, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	authUserID, ok := authUserIDVal.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type")
		return
	}

	var exerciseLog models.ExerciseLog
	if err := database.DB.
		Preload("Session").
		Preload("Exercise").
		Preload("Prescription").
		Preload("SetLogs").
		First(&exerciseLog, "id = ?", exerciseLogID).Error; err != nil {
		utils.NotFoundResponse(c, "Exercise log not found")
		return
	}

	// Check authorization via session
	if !canAccessSession(authUserID, &exerciseLog.Session) {
		utils.NotFoundResponse(c, "Exercise log not found")
		return
	}

	// Get user's preferred weight unit
	preferredUnit := getUserPreferredWeightUnit(authUserID)

	setLogs := make([]SetLogResponse, len(exerciseLog.SetLogs))
	for i, sl := range exerciseLog.SetLogs {
		setLogs[i] = buildSetLogResponse(sl, preferredUnit)
	}

	response := ExerciseLogResponse{
		ID:               exerciseLog.ID,
		SessionID:        exerciseLog.SessionID,
		PrescriptionID:   exerciseLog.PrescriptionID,
		ExerciseID:       exerciseLog.ExerciseID,
		ExerciseName:     exerciseLog.Exercise.Name,
		OrderNumber:      exerciseLog.OrderNumber,
		Notes:            exerciseLog.Notes,
		DifficultyRating: exerciseLog.DifficultyRating,
		DifficultyType:   exerciseLog.DifficultyType,
		SetLogs:          setLogs,
	}

	// Add prescription group info if available
	if exerciseLog.Prescription != nil {
		response.GroupID = &exerciseLog.Prescription.GroupID
		response.GroupName = exerciseLog.Prescription.GroupName
		response.GroupType = exerciseLog.Prescription.Type
	}

	utils.SuccessResponse(c, "Exercise log retrieved successfully", response)
}

// UpdateExerciseLog updates an exercise log
func UpdateExerciseLog(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	exerciseLogID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid UUID format", nil)
		return
	}

	authUserIDVal, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	authUserID, ok := authUserIDVal.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type")
		return
	}

	var exerciseLog models.ExerciseLog
	if err := database.DB.Preload("Session").First(&exerciseLog, "id = ?", exerciseLogID).Error; err != nil {
		utils.NotFoundResponse(c, "Exercise log not found")
		return
	}

	// Check authorization via session
	if !canAccessSession(authUserID, &exerciseLog.Session) {
		utils.ForbiddenResponse(c, "Not authorized to update this exercise log")
		return
	}

	var req UpdateExerciseLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	exerciseLog.Notes = req.Notes
	exerciseLog.DifficultyRating = req.DifficultyRating
	exerciseLog.DifficultyType = req.DifficultyType

	if err := database.DB.Save(&exerciseLog).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update exercise log")
		return
	}

	// Reload with relationships
	database.DB.Preload("Exercise").Preload("Prescription").Preload("SetLogs").First(&exerciseLog, "id = ?", exerciseLog.ID)

	// Get user's preferred weight unit
	preferredUnit := getUserPreferredWeightUnit(authUserID)

	setLogs := make([]SetLogResponse, len(exerciseLog.SetLogs))
	for i, sl := range exerciseLog.SetLogs {
		setLogs[i] = buildSetLogResponse(sl, preferredUnit)
	}

	response := ExerciseLogResponse{
		ID:               exerciseLog.ID,
		SessionID:        exerciseLog.SessionID,
		PrescriptionID:   exerciseLog.PrescriptionID,
		ExerciseID:       exerciseLog.ExerciseID,
		ExerciseName:     exerciseLog.Exercise.Name,
		OrderNumber:      exerciseLog.OrderNumber,
		Notes:            exerciseLog.Notes,
		DifficultyRating: exerciseLog.DifficultyRating,
		DifficultyType:   exerciseLog.DifficultyType,
		SetLogs:          setLogs,
	}

	// Add prescription group info if available
	if exerciseLog.Prescription != nil {
		response.GroupID = &exerciseLog.Prescription.GroupID
		response.GroupName = exerciseLog.Prescription.GroupName
		response.GroupType = exerciseLog.Prescription.Type
	}

	utils.SuccessResponse(c, "Exercise log updated successfully", response)
}

// DeleteExerciseLog deletes an exercise log
func DeleteExerciseLog(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	exerciseLogID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid UUID format", nil)
		return
	}

	authUserIDVal, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	authUserID, ok := authUserIDVal.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type")
		return
	}

	var exerciseLog models.ExerciseLog
	if err := database.DB.Preload("Session").First(&exerciseLog, "id = ?", exerciseLogID).Error; err != nil {
		utils.NotFoundResponse(c, "Exercise log not found")
		return
	}

	// Check authorization via session
	if !canAccessSession(authUserID, &exerciseLog.Session) {
		utils.ForbiddenResponse(c, "Not authorized to delete this exercise log")
		return
	}

	if err := database.DB.Delete(&exerciseLog).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete exercise log")
		return
	}

	utils.NoContentResponse(c)
}
