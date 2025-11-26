package controllers

import (
	"lamari-fit-api/database"
	"lamari-fit-api/models"
	"lamari-fit-api/utils"

	"github.com/gin-gonic/gin"
)

// GetSessionSet retrieves a single session set
func GetSessionSet(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	setID, ok := utils.ParseUUID(c, params.ID, "session set")
	if !ok {
		return
	}

	authUserID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var set models.SessionSet
	if err := database.DB.
		Preload("RPEValue").
		Preload("SessionExercise.SessionBlock.Session").
		First(&set, "id = ?", setID).Error; err != nil {
		utils.NotFoundResponse(c, "Session set not found")
		return
	}

	// Authorization: check session ownership
	if !isAuthorizedForSession(set.SessionExercise.SessionBlock.Session, authUserID) {
		utils.NotFoundResponse(c, "Session set not found")
		return
	}

	// Get user's preferred weight unit for response conversion
	preferredWeightUnit := getUserPreferredWeightUnit(c, authUserID)

	utils.SuccessResponse(c, "Session set retrieved successfully", set.ToResponse(preferredWeightUnit))
}

// UpdateSessionSet updates a session set
func UpdateSessionSet(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	setID, ok := utils.ParseUUID(c, params.ID, "session set")
	if !ok {
		return
	}

	authUserID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var set models.SessionSet
	if err := database.DB.
		Preload("SessionExercise.SessionBlock.Session").
		First(&set, "id = ?", setID).Error; err != nil {
		utils.NotFoundResponse(c, "Session set not found")
		return
	}

	// Authorization
	if !isAuthorizedForSession(set.SessionExercise.SessionBlock.Session, authUserID) {
		utils.ForbiddenResponse(c, "Not authorized to update this set")
		return
	}

	var req models.UpdateSessionSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Update fields if provided
	if req.ActualReps != nil {
		set.ActualReps = req.ActualReps
	}

	if req.ActualWeight != nil {
		// Process weight input using unified weight system
		actualWeightKg, originalValue, originalUnit := utils.ProcessWeightInput(req.ActualWeight)
		set.ActualWeightKg = actualWeightKg
		set.OriginalActualWeightValue = originalValue
		set.OriginalActualWeightUnit = originalUnit
	}

	if req.ActualDurationSeconds != nil {
		set.ActualDurationSeconds = req.ActualDurationSeconds
	}

	if req.RPEValueID != nil {
		set.RPEValueID = req.RPEValueID
	}

	if req.WasFailure != nil {
		set.WasFailure = *req.WasFailure
	}

	if req.Notes != nil {
		set.Notes = *req.Notes
	}

	if req.Completed != nil {
		set.Completed = *req.Completed
	}

	if err := database.DB.Save(&set).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update session set")
		return
	}

	// Reload with RPE value
	database.DB.Preload("RPEValue").First(&set, "id = ?", set.ID)

	// Get user's preferred weight unit for response conversion
	preferredWeightUnit := getUserPreferredWeightUnit(c, authUserID)

	utils.SuccessResponse(c, "Session set updated successfully", set.ToResponse(preferredWeightUnit))
}

// CompleteSessionSet marks a session set as complete
func CompleteSessionSet(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	setID, ok := utils.ParseUUID(c, params.ID, "session set")
	if !ok {
		return
	}

	authUserID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var set models.SessionSet
	if err := database.DB.
		Preload("SessionExercise.SessionBlock.Session").
		First(&set, "id = ?", setID).Error; err != nil {
		utils.NotFoundResponse(c, "Session set not found")
		return
	}

	// Authorization
	if !isAuthorizedForSession(set.SessionExercise.SessionBlock.Session, authUserID) {
		utils.ForbiddenResponse(c, "Not authorized to complete this set")
		return
	}

	// Optionally accept completion data
	var req models.UpdateSessionSetRequest
	if err := c.ShouldBindJSON(&req); err == nil {
		// Update fields if provided
		if req.ActualReps != nil {
			set.ActualReps = req.ActualReps
		}

		if req.ActualWeight != nil {
			// Process weight input using unified weight system
			actualWeightKg, originalValue, originalUnit := utils.ProcessWeightInput(req.ActualWeight)
			set.ActualWeightKg = actualWeightKg
			set.OriginalActualWeightValue = originalValue
			set.OriginalActualWeightUnit = originalUnit
		}

		if req.ActualDurationSeconds != nil {
			set.ActualDurationSeconds = req.ActualDurationSeconds
		}

		if req.RPEValueID != nil {
			set.RPEValueID = req.RPEValueID
		}

		if req.WasFailure != nil {
			set.WasFailure = *req.WasFailure
		}

		if req.Notes != nil {
			set.Notes = *req.Notes
		}
	}

	// Mark as completed
	set.Completed = true

	if err := database.DB.Save(&set).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to complete session set")
		return
	}

	// Reload with RPE value
	database.DB.Preload("RPEValue").First(&set, "id = ?", set.ID)

	// Get user's preferred weight unit for response conversion
	preferredWeightUnit := getUserPreferredWeightUnit(c, authUserID)

	utils.SuccessResponse(c, "Session set completed successfully", set.ToResponse(preferredWeightUnit))
}

// DeleteSessionSet deletes a session set
func DeleteSessionSet(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	setID, ok := utils.ParseUUID(c, params.ID, "session set")
	if !ok {
		return
	}

	authUserID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var set models.SessionSet
	if err := database.DB.
		Preload("SessionExercise.SessionBlock.Session").
		First(&set, "id = ?", setID).Error; err != nil {
		utils.NotFoundResponse(c, "Session set not found")
		return
	}

	// Authorization
	if !isAuthorizedForSession(set.SessionExercise.SessionBlock.Session, authUserID) {
		utils.ForbiddenResponse(c, "Not authorized to delete this set")
		return
	}

	if err := database.DB.Delete(&set).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete session set")
		return
	}

	// Renumber remaining sets
	var remainingSets []models.SessionSet
	database.DB.Where("session_exercise_id = ?", set.SessionExerciseID).Order("set_number ASC").Find(&remainingSets)

	for i, s := range remainingSets {
		s.SetNumber = i + 1
		database.DB.Save(&s)
	}

	utils.NoContentResponse(c)
}
