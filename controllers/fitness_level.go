package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"fit-flow-api/database"
	"fit-flow-api/middleware"
	"fit-flow-api/models"
)

// GetAllFitnessLevels retrieves all fitness levels
func GetAllFitnessLevels(c *gin.Context) {
	var levels []models.FitnessLevel

	if err := database.DB.Order("sort_order ASC").Find(&levels).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "common.server_error", err)
		return
	}

	response := make([]models.FitnessLevelResponse, len(levels))
	for i, level := range levels {
		response[i] = level.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// GetFitnessLevel retrieves a single fitness level by ID
func GetFitnessLevel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "common.invalid_id", err)
		return
	}

	var level models.FitnessLevel
	if err := database.DB.First(&level, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.TranslateErrorResponse(c, http.StatusNotFound, "fitness.level_not_found", nil)
			return
		}
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "common.server_error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": level.ToResponse(),
	})
}

// CreateFitnessLevel creates a new fitness level (admin only)
func CreateFitnessLevel(c *gin.Context) {
	// Check if user is admin
	user, exists := c.Get("user")
	if !exists || !user.(models.User).IsAdmin {
		middleware.TranslateErrorResponse(c, http.StatusForbidden, "auth.forbidden", nil)
		return
	}

	var req models.CreateFitnessLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "common.invalid_request", err)
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
			middleware.TranslateErrorResponse(c, http.StatusConflict, "fitness.level_name_exists", nil)
			return
		}
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "common.server_error", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": level.ToResponse(),
	})
}

// UpdateFitnessLevel updates an existing fitness level (admin only)
func UpdateFitnessLevel(c *gin.Context) {
	// Check if user is admin
	user, exists := c.Get("user")
	if !exists || !user.(models.User).IsAdmin {
		middleware.TranslateErrorResponse(c, http.StatusForbidden, "auth.forbidden", nil)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "common.invalid_id", err)
		return
	}

	var req models.UpdateFitnessLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "common.invalid_request", err)
		return
	}

	var level models.FitnessLevel
	if err := database.DB.First(&level, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.TranslateErrorResponse(c, http.StatusNotFound, "fitness.level_not_found", nil)
			return
		}
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "common.server_error", err)
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
			middleware.TranslateErrorResponse(c, http.StatusConflict, "fitness.level_name_exists", nil)
			return
		}
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "common.server_error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": level.ToResponse(),
	})
}

// DeleteFitnessLevel deletes a fitness level (admin only)
func DeleteFitnessLevel(c *gin.Context) {
	// Check if user is admin
	user, exists := c.Get("user")
	if !exists || !user.(models.User).IsAdmin {
		middleware.TranslateErrorResponse(c, http.StatusForbidden, "auth.forbidden", nil)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "common.invalid_id", err)
		return
	}

	result := database.DB.Delete(&models.FitnessLevel{}, id)
	if result.Error != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "common.server_error", result.Error)
		return
	}

	if result.RowsAffected == 0 {
		middleware.TranslateErrorResponse(c, http.StatusNotFound, "fitness.level_not_found", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Fitness level deleted successfully",
	})
}