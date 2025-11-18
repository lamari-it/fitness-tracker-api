package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"
)

// GetAllFitnessLevels retrieves all fitness levels
func GetAllFitnessLevels(c *gin.Context) {
	var queryParams PaginationQuery
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Set default pagination values
	SetDefaultPagination(&queryParams)

	offset := (queryParams.Page - 1) * queryParams.Limit

	// Get total count
	var total int64
	if err := database.DB.Model(&models.FitnessLevel{}).Count(&total).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to count fitness levels.")
		return
	}

	var levels []models.FitnessLevel
	if err := database.DB.Order("sort_order ASC").
		Offset(offset).
		Limit(queryParams.Limit).
		Find(&levels).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve fitness levels.")
		return
	}

	response := make([]models.FitnessLevelResponse, len(levels))
	for i, level := range levels {
		response[i] = level.ToResponse()
	}

	utils.PaginatedResponse(c, "Fitness levels retrieved successfully.", response, queryParams.Page, queryParams.Limit, int(total))
}

// GetFitnessLevel retrieves a single fitness level by ID
func GetFitnessLevel(c *gin.Context) {
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

	var level models.FitnessLevel
	if err := database.DB.First(&level, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "Fitness level not found.")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to retrieve fitness level.")
		return
	}

	utils.SuccessResponse(c, "Fitness level retrieved successfully.", level.ToResponse())
}

// CreateFitnessLevel creates a new fitness level (admin only)
func CreateFitnessLevel(c *gin.Context) {
	// Check if user is admin
	user, exists := c.Get("user")
	if !exists || !user.(models.User).IsAdmin {
		utils.ForbiddenResponse(c, "Admin access required.")
		return
	}

	var req models.CreateFitnessLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	level := models.FitnessLevel{
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}

	if err := database.DB.Create(&level).Error; err != nil {
		if err.Error() == `ERROR: duplicate key value violates unique constraint "idx_fitness_levels_name" (SQLSTATE 23505)` ||
			err.Error() == `ERROR: duplicate key value violates unique constraint "fitness_levels_name_key" (SQLSTATE 23505)` {
			utils.ConflictResponse(c, "A fitness level with this name already exists.")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to create fitness level.")
		return
	}

	utils.CreatedResponse(c, "Fitness level created successfully.", level.ToResponse())
}

// UpdateFitnessLevel updates an existing fitness level (admin only)
func UpdateFitnessLevel(c *gin.Context) {
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

	var req models.UpdateFitnessLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	var level models.FitnessLevel
	if err := database.DB.First(&level, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "Fitness level not found.")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to retrieve fitness level.")
		return
	}

	// Update fields
	if req.Name != "" {
		level.Name = req.Name
	}
	if req.Description != "" {
		level.Description = req.Description
	}
	if req.SortOrder != 0 {
		level.SortOrder = req.SortOrder
	}

	if err := database.DB.Save(&level).Error; err != nil {
		if err.Error() == `ERROR: duplicate key value violates unique constraint "idx_fitness_levels_name" (SQLSTATE 23505)` ||
			err.Error() == `ERROR: duplicate key value violates unique constraint "fitness_levels_name_key" (SQLSTATE 23505)` {
			utils.ConflictResponse(c, "A fitness level with this name already exists.")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to update fitness level.")
		return
	}

	utils.SuccessResponse(c, "Fitness level updated successfully.", level.ToResponse())
}

// DeleteFitnessLevel deletes a fitness level (admin only)
func DeleteFitnessLevel(c *gin.Context) {
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

	result := database.DB.Delete(&models.FitnessLevel{}, id)
	if result.Error != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete fitness level.")
		return
	}

	if result.RowsAffected == 0 {
		utils.NotFoundResponse(c, "Fitness level not found.")
		return
	}

	utils.DeletedResponse(c, "Fitness level deleted successfully.")
}
