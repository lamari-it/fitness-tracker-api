package controllers

import (
	"lamari-fit-api/database"
	"lamari-fit-api/models"
	"lamari-fit-api/utils"

	"github.com/gin-gonic/gin"
)

// GetUserSettings retrieves the current user's settings
func GetUserSettings(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	utils.SuccessResponse(c, "User settings retrieved successfully", user.ToResponse())
}

// UpdateUserSettings updates the current user's settings/preferences
func UpdateUserSettings(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var req models.UpdateUserSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	// Update only provided fields
	if req.PreferredWeightUnit != "" {
		user.PreferredWeightUnit = req.PreferredWeightUnit
	}
	if req.PreferredHeightUnit != "" {
		user.PreferredHeightUnit = req.PreferredHeightUnit
	}
	if req.PreferredDistanceUnit != "" {
		user.PreferredDistanceUnit = req.PreferredDistanceUnit
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}

	if err := database.DB.Save(&user).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update user settings")
		return
	}

	utils.SuccessResponse(c, "User settings updated successfully", user.ToResponse())
}
