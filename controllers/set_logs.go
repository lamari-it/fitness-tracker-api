package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Request DTOs
type CreateSetLogRequest struct {
	ExerciseLogID uuid.UUID  `json:"exercise_log_id" binding:"required"`
	SetNumber     int        `json:"set_number" binding:"required"`
	Weight        float64    `json:"weight"`
	WeightUnit    string     `json:"weight_unit"`
	Reps          int        `json:"reps"`
	RestAfterSec  int        `json:"rest_after_sec"`
	Tempo         string     `json:"tempo"`
	RPE           float64    `json:"rpe"`
	RPEValueID    *uuid.UUID `json:"rpe_value_id"`
}

type UpdateSetLogRequest struct {
	Weight       float64    `json:"weight"`
	WeightUnit   string     `json:"weight_unit"`
	Reps         int        `json:"reps"`
	RestAfterSec int        `json:"rest_after_sec"`
	Tempo        string     `json:"tempo"`
	RPE          float64    `json:"rpe"`
	RPEValueID   *uuid.UUID `json:"rpe_value_id"`
}

// CreateSetLog creates a new set log for an exercise
func CreateSetLog(c *gin.Context) {
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

	var req CreateSetLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Verify exercise log exists and get its session
	var exerciseLog models.ExerciseLog
	if err := database.DB.Preload("Session").First(&exerciseLog, "id = ?", req.ExerciseLogID).Error; err != nil {
		utils.NotFoundResponse(c, "Exercise log not found")
		return
	}

	// Check authorization via session
	if !canAccessSession(authUserID, &exerciseLog.Session) {
		utils.ForbiddenResponse(c, "Not authorized to log sets in this exercise")
		return
	}

	// Determine the input weight unit
	inputUnit := req.WeightUnit
	if inputUnit == "" {
		// Default to user's preferred unit if not specified
		inputUnit = getUserPreferredWeightUnit(authUserID)
	}

	// Normalize the unit
	normalizedUnit := utils.NormalizeWeightUnit(inputUnit)

	// Convert to canonical kg for storage
	weightInKg := utils.ConvertToKg(req.Weight, normalizedUnit)

	setLog := models.SetLog{
		ExerciseLogID:   req.ExerciseLogID,
		SetNumber:       req.SetNumber,
		Weight:          weightInKg,      // Canonical storage in kg
		WeightUnit:      normalizedUnit,  // Keep for backward compatibility
		InputWeight:     req.Weight,      // Original input value
		InputWeightUnit: normalizedUnit,  // Original input unit
		Reps:            req.Reps,
		RestAfterSec:    req.RestAfterSec,
		Tempo:           req.Tempo,
		RPE:             req.RPE,
		RPEValueID:      req.RPEValueID,
	}

	if err := database.DB.Create(&setLog).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create set log")
		return
	}

	// Get user's preferred unit for response
	preferredUnit := getUserPreferredWeightUnit(authUserID)
	response := buildSetLogResponse(setLog, preferredUnit)

	utils.CreatedResponse(c, "Set log created successfully", response)
}

// GetSetLog retrieves a single set log
func GetSetLog(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	setLogID, err := uuid.Parse(params.ID)
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

	var setLog models.SetLog
	if err := database.DB.
		Preload("ExerciseLog.Session").
		First(&setLog, "id = ?", setLogID).Error; err != nil {
		utils.NotFoundResponse(c, "Set log not found")
		return
	}

	// Check authorization via session
	if !canAccessSession(authUserID, &setLog.ExerciseLog.Session) {
		utils.NotFoundResponse(c, "Set log not found")
		return
	}

	// Get user's preferred unit for response
	preferredUnit := getUserPreferredWeightUnit(authUserID)
	response := buildSetLogResponse(setLog, preferredUnit)

	utils.SuccessResponse(c, "Set log retrieved successfully", response)
}

// UpdateSetLog updates a set log
func UpdateSetLog(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	setLogID, err := uuid.Parse(params.ID)
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

	var setLog models.SetLog
	if err := database.DB.
		Preload("ExerciseLog.Session").
		First(&setLog, "id = ?", setLogID).Error; err != nil {
		utils.NotFoundResponse(c, "Set log not found")
		return
	}

	// Check authorization via session
	if !canAccessSession(authUserID, &setLog.ExerciseLog.Session) {
		utils.ForbiddenResponse(c, "Not authorized to update this set log")
		return
	}

	var req UpdateSetLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Determine the input weight unit
	inputUnit := req.WeightUnit
	if inputUnit == "" {
		// Default to user's preferred unit if not specified
		inputUnit = getUserPreferredWeightUnit(authUserID)
	}

	// Normalize the unit
	normalizedUnit := utils.NormalizeWeightUnit(inputUnit)

	// Convert to canonical kg for storage
	weightInKg := utils.ConvertToKg(req.Weight, normalizedUnit)

	setLog.Weight = weightInKg           // Canonical storage in kg
	setLog.WeightUnit = normalizedUnit   // Keep for backward compatibility
	setLog.InputWeight = req.Weight      // Original input value
	setLog.InputWeightUnit = normalizedUnit // Original input unit
	setLog.Reps = req.Reps
	setLog.RestAfterSec = req.RestAfterSec
	setLog.Tempo = req.Tempo
	setLog.RPE = req.RPE
	setLog.RPEValueID = req.RPEValueID

	if err := database.DB.Save(&setLog).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update set log")
		return
	}

	// Get user's preferred unit for response
	preferredUnit := getUserPreferredWeightUnit(authUserID)
	response := buildSetLogResponse(setLog, preferredUnit)

	utils.SuccessResponse(c, "Set log updated successfully", response)
}

// DeleteSetLog deletes a set log
func DeleteSetLog(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	setLogID, err := uuid.Parse(params.ID)
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

	var setLog models.SetLog
	if err := database.DB.
		Preload("ExerciseLog.Session").
		First(&setLog, "id = ?", setLogID).Error; err != nil {
		utils.NotFoundResponse(c, "Set log not found")
		return
	}

	// Check authorization via session
	if !canAccessSession(authUserID, &setLog.ExerciseLog.Session) {
		utils.ForbiddenResponse(c, "Not authorized to delete this set log")
		return
	}

	if err := database.DB.Delete(&setLog).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete set log")
		return
	}

	utils.NoContentResponse(c)
}
