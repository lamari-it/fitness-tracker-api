package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	id, ok := utils.ParseUUID(c, params.ID, "fitness goal")
	if !ok {
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
	if !utils.RequireAdmin(c) {
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
	if !utils.RequireAdmin(c) {
		return
	}

	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	id, ok := utils.ParseUUID(c, params.ID, "fitness goal")
	if !ok {
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
	if !utils.RequireAdmin(c) {
		return
	}

	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	id, ok := utils.ParseUUID(c, params.ID, "fitness goal")
	if !ok {
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

	// Find user's fitness profile
	var profile models.UserFitnessProfile
	if err := database.DB.Where("user_id = ?", currentUser.ID).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "Fitness profile not found. Please create a fitness profile first.")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to find fitness profile.")
		return
	}

	// Update fitness profile's fitness level
	if err := database.DB.Model(&profile).Update("fitness_level_id", req.FitnessLevelID).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update fitness level.")
		return
	}

	// Reload profile with fitness level
	database.DB.Preload("FitnessLevel").Preload("FitnessGoals.FitnessGoal").First(&profile, "id = ?", profile.ID)

	utils.SuccessResponse(c, "Fitness level updated successfully.", profile.ToResponse(currentUser.PreferredWeightUnit))
}
