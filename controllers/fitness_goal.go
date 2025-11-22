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
	var queryParams FitnessGoalQuery
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Set default pagination values
	SetDefaultPagination(&queryParams.PaginationQuery)

	offset := (queryParams.Page - 1) * queryParams.Limit

	query := database.DB.Model(&models.FitnessGoal{}).Order("category, name")

	// Optional category filter
	if queryParams.Category != "" {
		query = query.Where("category = ?", queryParams.Category)
	}

	// Get total count
	var total int64
	query.Count(&total)

	var goals []models.FitnessGoal
	if err := query.Offset(offset).Limit(queryParams.Limit).Find(&goals).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve fitness goals.")
		return
	}

	response := make([]models.FitnessGoalResponse, len(goals))
	for i, goal := range goals {
		response[i] = goal.ToResponse()
	}

	utils.PaginatedResponse(c, "Fitness goals retrieved successfully.", response, queryParams.Page, queryParams.Limit, int(total))
}

// GetFitnessGoal retrieves a single fitness goal by ID
func GetFitnessGoal(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	id, err := uuid.Parse(params.ID)
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

	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	id, err := uuid.Parse(params.ID)
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

	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	id, err := uuid.Parse(params.ID)
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
