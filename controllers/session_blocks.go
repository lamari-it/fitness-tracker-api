package controllers

import (
	"time"

	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetSessionBlock retrieves a single session block
func GetSessionBlock(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	blockID, ok := utils.ParseUUID(c, params.ID, "session block")
	if !ok {
		return
	}

	authUserID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var block models.SessionBlock
	if err := database.DB.
		Preload("SessionExercises.Exercise").
		Preload("SessionExercises.SessionSets.RPEValue").
		Preload("Session").
		First(&block, "id = ?", blockID).Error; err != nil {
		utils.NotFoundResponse(c, "Session block not found")
		return
	}

	// Authorization: check session ownership
	if !isAuthorizedForSession(block.Session, authUserID) {
		utils.NotFoundResponse(c, "Session block not found")
		return
	}

	// Get user's preferred weight unit for response conversion
	preferredWeightUnit := getUserPreferredWeightUnit(c, authUserID)

	utils.SuccessResponse(c, "Session block retrieved successfully", block.ToResponse(preferredWeightUnit))
}

// CompleteSessionBlock marks a session block as complete
func CompleteSessionBlock(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	blockID, ok := utils.ParseUUID(c, params.ID, "session block")
	if !ok {
		return
	}

	authUserID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var block models.SessionBlock
	if err := database.DB.Preload("Session").First(&block, "id = ?", blockID).Error; err != nil {
		utils.NotFoundResponse(c, "Session block not found")
		return
	}

	// Authorization
	if !isAuthorizedForSession(block.Session, authUserID) {
		utils.ForbiddenResponse(c, "Not authorized to complete this block")
		return
	}

	// Set completion time
	now := time.Now()
	block.CompletedAt = &now
	block.Skipped = false

	if err := database.DB.Save(&block).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to complete session block")
		return
	}

	// Reload with relationships
	database.DB.
		Preload("SessionExercises.Exercise").
		Preload("SessionExercises.SessionSets.RPEValue").
		First(&block, "id = ?", block.ID)

	// Get user's preferred weight unit for response conversion
	preferredWeightUnit := getUserPreferredWeightUnit(c, authUserID)

	utils.SuccessResponse(c, "Session block completed successfully", block.ToResponse(preferredWeightUnit))
}

// SkipSessionBlock marks a session block as skipped
func SkipSessionBlock(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	blockID, ok := utils.ParseUUID(c, params.ID, "session block")
	if !ok {
		return
	}

	authUserID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var block models.SessionBlock
	if err := database.DB.Preload("Session").First(&block, "id = ?", blockID).Error; err != nil {
		utils.NotFoundResponse(c, "Session block not found")
		return
	}

	// Authorization
	if !isAuthorizedForSession(block.Session, authUserID) {
		utils.ForbiddenResponse(c, "Not authorized to skip this block")
		return
	}

	// Mark as skipped
	block.Skipped = true
	block.CompletedAt = nil

	if err := database.DB.Save(&block).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to skip session block")
		return
	}

	// Reload with relationships
	database.DB.
		Preload("SessionExercises.Exercise").
		Preload("SessionExercises.SessionSets.RPEValue").
		First(&block, "id = ?", block.ID)

	// Get user's preferred weight unit for response conversion
	preferredWeightUnit := getUserPreferredWeightUnit(c, authUserID)

	utils.SuccessResponse(c, "Session block skipped successfully", block.ToResponse(preferredWeightUnit))
}

// UpdateSessionBlockRPE updates the perceived exertion for a block
func UpdateSessionBlockRPE(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	blockID, ok := utils.ParseUUID(c, params.ID, "session block")
	if !ok {
		return
	}

	authUserID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var block models.SessionBlock
	if err := database.DB.Preload("Session").First(&block, "id = ?", blockID).Error; err != nil {
		utils.NotFoundResponse(c, "Session block not found")
		return
	}

	// Authorization
	if !isAuthorizedForSession(block.Session, authUserID) {
		utils.ForbiddenResponse(c, "Not authorized to update this block")
		return
	}

	var req struct {
		PerceivedExertion *int `json:"perceived_exertion" binding:"omitempty,min=1,max=10"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	block.PerceivedExertion = req.PerceivedExertion

	if err := database.DB.Save(&block).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update session block")
		return
	}

	// Reload with relationships
	database.DB.
		Preload("SessionExercises.Exercise").
		Preload("SessionExercises.SessionSets.RPEValue").
		First(&block, "id = ?", block.ID)

	// Get user's preferred weight unit for response conversion
	preferredWeightUnit := getUserPreferredWeightUnit(c, authUserID)

	utils.SuccessResponse(c, "Session block updated successfully", block.ToResponse(preferredWeightUnit))
}

// Helper functions

// isAuthorizedForSession checks if user owns or created the session.
func isAuthorizedForSession(session models.WorkoutSession, authUserID uuid.UUID) bool {
	isOwner := session.UserID == authUserID
	isCreator := session.CreatedByID != nil && *session.CreatedByID == authUserID
	return isOwner || isCreator
}

// getUserPreferredWeightUnit retrieves user's preferred weight unit.
func getUserPreferredWeightUnit(c *gin.Context, userID uuid.UUID) string {
	var user models.User
	if err := database.DB.Select("preferred_weight_unit").First(&user, "id = ?", userID).Error; err != nil {
		return "kg" // default fallback
	}
	if user.PreferredWeightUnit == "" {
		return "kg"
	}
	return user.PreferredWeightUnit
}
