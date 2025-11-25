package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetSessionExercise retrieves a single session exercise
func GetSessionExercise(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	exerciseID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid UUID format", nil)
		return
	}

	authUserID, err := getAuthUserID(c)
	if err != nil {
		return
	}

	var sessionExercise models.SessionExercise
	if err := database.DB.
		Preload("Exercise").
		Preload("SessionSets.RPEValue").
		Preload("SessionBlock.Session").
		First(&sessionExercise, "id = ?", exerciseID).Error; err != nil {
		utils.NotFoundResponse(c, "Session exercise not found")
		return
	}

	// Authorization: check session ownership
	if !isAuthorizedForSession(sessionExercise.SessionBlock.Session, authUserID) {
		utils.NotFoundResponse(c, "Session exercise not found")
		return
	}

	// Get user's preferred weight unit for response conversion
	preferredWeightUnit := getUserPreferredWeightUnit(c, authUserID)

	utils.SuccessResponse(c, "Session exercise retrieved successfully", sessionExercise.ToResponse(preferredWeightUnit))
}

// CompleteSessionExercise marks a session exercise as complete
func CompleteSessionExercise(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	exerciseID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid UUID format", nil)
		return
	}

	authUserID, err := getAuthUserID(c)
	if err != nil {
		return
	}

	var sessionExercise models.SessionExercise
	if err := database.DB.
		Preload("SessionBlock.Session").
		First(&sessionExercise, "id = ?", exerciseID).Error; err != nil {
		utils.NotFoundResponse(c, "Session exercise not found")
		return
	}

	// Authorization
	if !isAuthorizedForSession(sessionExercise.SessionBlock.Session, authUserID) {
		utils.ForbiddenResponse(c, "Not authorized to complete this exercise")
		return
	}

	// Set completion time
	now := time.Now()
	sessionExercise.CompletedAt = &now
	sessionExercise.Skipped = false

	if err := database.DB.Save(&sessionExercise).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to complete session exercise")
		return
	}

	// Reload with relationships
	database.DB.
		Preload("Exercise").
		Preload("SessionSets.RPEValue").
		First(&sessionExercise, "id = ?", sessionExercise.ID)

	// Get user's preferred weight unit for response conversion
	preferredWeightUnit := getUserPreferredWeightUnit(c, authUserID)

	utils.SuccessResponse(c, "Session exercise completed successfully", sessionExercise.ToResponse(preferredWeightUnit))
}

// SkipSessionExercise marks a session exercise as skipped
func SkipSessionExercise(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	exerciseID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid UUID format", nil)
		return
	}

	authUserID, err := getAuthUserID(c)
	if err != nil {
		return
	}

	var sessionExercise models.SessionExercise
	if err := database.DB.
		Preload("SessionBlock.Session").
		First(&sessionExercise, "id = ?", exerciseID).Error; err != nil {
		utils.NotFoundResponse(c, "Session exercise not found")
		return
	}

	// Authorization
	if !isAuthorizedForSession(sessionExercise.SessionBlock.Session, authUserID) {
		utils.ForbiddenResponse(c, "Not authorized to skip this exercise")
		return
	}

	// Mark as skipped
	sessionExercise.Skipped = true
	sessionExercise.CompletedAt = nil

	if err := database.DB.Save(&sessionExercise).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to skip session exercise")
		return
	}

	// Reload with relationships
	database.DB.
		Preload("Exercise").
		Preload("SessionSets.RPEValue").
		First(&sessionExercise, "id = ?", sessionExercise.ID)

	// Get user's preferred weight unit for response conversion
	preferredWeightUnit := getUserPreferredWeightUnit(c, authUserID)

	utils.SuccessResponse(c, "Session exercise skipped successfully", sessionExercise.ToResponse(preferredWeightUnit))
}

// UpdateSessionExerciseNotes updates the notes for a session exercise
func UpdateSessionExerciseNotes(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	exerciseID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid UUID format", nil)
		return
	}

	authUserID, err := getAuthUserID(c)
	if err != nil {
		return
	}

	var sessionExercise models.SessionExercise
	if err := database.DB.
		Preload("SessionBlock.Session").
		First(&sessionExercise, "id = ?", exerciseID).Error; err != nil {
		utils.NotFoundResponse(c, "Session exercise not found")
		return
	}

	// Authorization
	if !isAuthorizedForSession(sessionExercise.SessionBlock.Session, authUserID) {
		utils.ForbiddenResponse(c, "Not authorized to update this exercise")
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	sessionExercise.Notes = req.Notes

	if err := database.DB.Save(&sessionExercise).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update session exercise")
		return
	}

	// Reload with relationships
	database.DB.
		Preload("Exercise").
		Preload("SessionSets.RPEValue").
		First(&sessionExercise, "id = ?", sessionExercise.ID)

	// Get user's preferred weight unit for response conversion
	preferredWeightUnit := getUserPreferredWeightUnit(c, authUserID)

	utils.SuccessResponse(c, "Session exercise updated successfully", sessionExercise.ToResponse(preferredWeightUnit))
}

// AddSetToExercise adds a new set to a session exercise
func AddSetToExercise(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	exerciseID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid UUID format", nil)
		return
	}

	authUserID, err := getAuthUserID(c)
	if err != nil {
		return
	}

	var sessionExercise models.SessionExercise
	if err := database.DB.
		Preload("SessionBlock.Session").
		Preload("SessionSets").
		First(&sessionExercise, "id = ?", exerciseID).Error; err != nil {
		utils.NotFoundResponse(c, "Session exercise not found")
		return
	}

	// Authorization
	if !isAuthorizedForSession(sessionExercise.SessionBlock.Session, authUserID) {
		utils.ForbiddenResponse(c, "Not authorized to add sets to this exercise")
		return
	}

	var req models.CreateSessionSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Determine next set number
	nextSetNumber := len(sessionExercise.SessionSets) + 1

	// Process weight input using new unified weight system
	actualWeightKg, originalValue, originalUnit := utils.ProcessWeightInput(req.ActualWeight)

	set := models.SessionSet{
		SessionExerciseID:          sessionExercise.ID,
		SetNumber:                  nextSetNumber,
		Completed:                  false,
		ActualReps:                 req.ActualReps,
		ActualWeightKg:             actualWeightKg,
		OriginalActualWeightValue:  originalValue,
		OriginalActualWeightUnit:   originalUnit,
		ActualDurationSeconds:      req.ActualDurationSeconds,
		RPEValueID:                 req.RPEValueID,
		WasFailure:                 false,
		Notes:                      "",
	}

	if req.Notes != nil {
		set.Notes = *req.Notes
	}

	if err := database.DB.Create(&set).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to add set")
		return
	}

	// Reload with RPE value
	database.DB.Preload("RPEValue").First(&set, "id = ?", set.ID)

	// Get user's preferred weight unit for response conversion
	preferredWeightUnit := getUserPreferredWeightUnit(c, authUserID)

	utils.CreatedResponse(c, "Set added successfully", set.ToResponse(preferredWeightUnit))
}
