package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"
)

// GetAllFitnessGoals retrieves all fitness goals
func GetAllFitnessGoals(c *gin.Context) {
	var goals []models.FitnessGoal

	query := database.DB.Order("category, name")

	// Optional category filter
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Find(&goals).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve fitness goals.")
		return
	}

	response := make([]models.FitnessGoalResponse, len(goals))
	for i, goal := range goals {
		response[i] = goal.ToResponse()
	}

	utils.SuccessResponse(c, "Fitness goals retrieved successfully.", response)
}

// GetFitnessGoal retrieves a single fitness goal by ID
func GetFitnessGoal(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID format.", nil)
		return
	}

	var goal models.FitnessGoal
	if err := database.DB.First(&goal, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "Fitness goal not found.")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to retrieve fitness goal.")
		return
	}

	utils.SuccessResponse(c, "Fitness goal retrieved successfully.", goal.ToResponse())
}

// CreateFitnessGoal creates a new fitness goal (admin only)
func CreateFitnessGoal(c *gin.Context) {
	// Check if user is admin
	user, exists := c.Get("user")
	if !exists || !user.(models.User).IsAdmin {
		utils.ForbiddenResponse(c, "Admin access required.")
		return
	}

	var req models.CreateFitnessGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	goal := models.FitnessGoal{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		IconName:    req.IconName,
	}

	if err := database.DB.Create(&goal).Error; err != nil {
		if err.Error() == `ERROR: duplicate key value violates unique constraint "idx_fitness_goals_name" (SQLSTATE 23505)` ||
			err.Error() == `ERROR: duplicate key value violates unique constraint "fitness_goals_name_key" (SQLSTATE 23505)` {
			utils.ConflictResponse(c, "A fitness goal with this name already exists.")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to create fitness goal.")
		return
	}

	utils.CreatedResponse(c, "Fitness goal created successfully.", goal.ToResponse())
}

// UpdateFitnessGoal updates an existing fitness goal (admin only)
func UpdateFitnessGoal(c *gin.Context) {
	// Check if user is admin
	user, exists := c.Get("user")
	if !exists || !user.(models.User).IsAdmin {
		utils.ForbiddenResponse(c, "Admin access required.")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID format.", nil)
		return
	}

	var req models.UpdateFitnessGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	var goal models.FitnessGoal
	if err := database.DB.First(&goal, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "Fitness goal not found.")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to retrieve fitness goal.")
		return
	}

	// Update fields
	if req.Name != "" {
		goal.Name = req.Name
	}
	if req.Description != "" {
		goal.Description = req.Description
	}
	if req.Category != "" {
		goal.Category = req.Category
	}
	if req.IconName != "" {
		goal.IconName = req.IconName
	}

	if err := database.DB.Save(&goal).Error; err != nil {
		if err.Error() == `ERROR: duplicate key value violates unique constraint "idx_fitness_goals_name" (SQLSTATE 23505)` ||
			err.Error() == `ERROR: duplicate key value violates unique constraint "fitness_goals_name_key" (SQLSTATE 23505)` {
			utils.ConflictResponse(c, "A fitness goal with this name already exists.")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to update fitness goal.")
		return
	}

	utils.SuccessResponse(c, "Fitness goal updated successfully.", goal.ToResponse())
}

// DeleteFitnessGoal deletes a fitness goal (admin only)
func DeleteFitnessGoal(c *gin.Context) {
	// Check if user is admin
	user, exists := c.Get("user")
	if !exists || !user.(models.User).IsAdmin {
		utils.ForbiddenResponse(c, "Admin access required.")
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID format.", nil)
		return
	}

	result := database.DB.Delete(&models.FitnessGoal{}, id)
	if result.Error != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete fitness goal.")
		return
	}

	if result.RowsAffected == 0 {
		utils.NotFoundResponse(c, "Fitness goal not found.")
		return
	}

	utils.DeletedResponse(c, "Fitness goal deleted successfully.")
}

// GetUserFitnessGoals retrieves the authenticated user's fitness goals
func GetUserFitnessGoals(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "Authentication required.")
		return
	}

	currentUser := user.(models.User)

	var userGoals []models.UserFitnessGoal
	if err := database.DB.Preload("FitnessGoal").
		Where("user_id = ?", currentUser.ID).
		Order("priority ASC").
		Find(&userGoals).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve user fitness goals.")
		return
	}

	response := make([]models.UserFitnessGoalResponse, len(userGoals))
	for i, userGoal := range userGoals {
		response[i] = userGoal.ToResponse()
	}

	utils.SuccessResponse(c, "User fitness goals retrieved successfully.", response)
}

// SetUserFitnessGoals sets the authenticated user's fitness goals
func SetUserFitnessGoals(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "Authentication required.")
		return
	}

	currentUser := user.(models.User)

	var req models.SetUserFitnessGoalsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Start transaction
	tx := database.DB.Begin()

	// Delete existing user goals
	if err := tx.Where("user_id = ?", currentUser.ID).Delete(&models.UserFitnessGoal{}).Error; err != nil {
		tx.Rollback()
		utils.InternalServerErrorResponse(c, "Failed to update user fitness goals.")
		return
	}

	// Create new user goals
	userGoals := make([]models.UserFitnessGoal, len(req.Goals))
	for i, goalInput := range req.Goals {
		// Verify fitness goal exists
		var goal models.FitnessGoal
		if err := tx.First(&goal, goalInput.FitnessGoalID).Error; err != nil {
			tx.Rollback()
			if err == gorm.ErrRecordNotFound {
				utils.BadRequestResponse(c, "One or more fitness goals not found.", nil)
				return
			}
			utils.InternalServerErrorResponse(c, "Failed to validate fitness goals.")
			return
		}

		userGoals[i] = models.UserFitnessGoal{
			UserID:        currentUser.ID,
			FitnessGoalID: goalInput.FitnessGoalID,
			Priority:      goalInput.Priority,
			TargetDate:    goalInput.TargetDate,
			Notes:         goalInput.Notes,
		}
	}

	if len(userGoals) > 0 {
		if err := tx.Create(&userGoals).Error; err != nil {
			tx.Rollback()
			utils.InternalServerErrorResponse(c, "Failed to save user fitness goals.")
			return
		}
	}

	tx.Commit()

	// Reload with associations
	var updatedGoals []models.UserFitnessGoal
	if err := database.DB.Preload("FitnessGoal").
		Where("user_id = ?", currentUser.ID).
		Order("priority ASC").
		Find(&updatedGoals).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve updated user fitness goals.")
		return
	}

	response := make([]models.UserFitnessGoalResponse, len(updatedGoals))
	for i, userGoal := range updatedGoals {
		response[i] = userGoal.ToResponse()
	}

	utils.SuccessResponse(c, "User fitness goals updated successfully.", response)
}

// UpdateUserFitnessLevel updates the authenticated user's fitness level
func UpdateUserFitnessLevel(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.UnauthorizedResponse(c, "Authentication required.")
		return
	}

	currentUser := user.(models.User)

	var req models.UpdateUserFitnessLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// If fitness level ID is provided, verify it exists
	if req.FitnessLevelID != nil {
		var level models.FitnessLevel
		if err := database.DB.First(&level, *req.FitnessLevelID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				utils.BadRequestResponse(c, "Fitness level not found.", nil)
				return
			}
			utils.InternalServerErrorResponse(c, "Failed to validate fitness level.")
			return
		}
	}

	// Update user's fitness level
	if err := database.DB.Model(&currentUser).Update("fitness_level_id", req.FitnessLevelID).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update user fitness level.")
		return
	}

	// Reload user with fitness level
	if err := database.DB.Preload("FitnessLevel").First(&currentUser, currentUser.ID).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve updated user.")
		return
	}

	utils.SuccessResponse(c, "User fitness level updated successfully.", currentUser.ToResponse())
}
