package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"
	"time"

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

	blockID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid UUID format", nil)
		return
	}

	authUserID, err := getAuthUserID(c)
	if err != nil {
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

	blockID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid UUID format", nil)
		return
	}

	authUserID, err := getAuthUserID(c)
	if err != nil {
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

	blockID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid UUID format", nil)
		return
	}

	authUserID, err := getAuthUserID(c)
	if err != nil {
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

	blockID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid UUID format", nil)
		return
	}

	authUserID, err := getAuthUserID(c)
	if err != nil {
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

func getAuthUserID(c *gin.Context) (uuid.UUID, error) {
	authUserIDVal, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return uuid.Nil, nil
	}

	authUserID, ok := authUserIDVal.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type")
		return uuid.Nil, nil
	}

	return authUserID, nil
}

func isAuthorizedForSession(session models.WorkoutSession, authUserID uuid.UUID) bool {
	isOwner := session.UserID == authUserID
	isCreator := session.CreatedByID != nil && *session.CreatedByID == authUserID
	return isOwner || isCreator
}

func getUserPreferredWeightUnit(c *gin.Context, userID uuid.UUID) string {
	var user models.User
	if err := database.DB.Select("preferred_weight_unit").First(&user, "id = ?", userID).Error; err != nil {
		return "kg" // default fallback
	}
	return user.PreferredWeightUnit
}
