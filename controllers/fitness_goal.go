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

// GetAllFitnessGoals retrieves all fitness goals
func GetAllFitnessGoals(c *gin.Context) {
	var goals []models.FitnessGoal

	query := database.DB.Order("category, name")
	
	// Optional category filter
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Find(&goals).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "common.server_error", err)
		return
	}

	response := make([]models.FitnessGoalResponse, len(goals))
	for i, goal := range goals {
		response[i] = goal.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// GetFitnessGoal retrieves a single fitness goal by ID
func GetFitnessGoal(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "common.invalid_id", err)
		return
	}

	var goal models.FitnessGoal
	if err := database.DB.First(&goal, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.TranslateErrorResponse(c, http.StatusNotFound, "fitness.goal_not_found", nil)
			return
		}
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "common.server_error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": goal.ToResponse(),
	})
}

// CreateFitnessGoal creates a new fitness goal (admin only)
func CreateFitnessGoal(c *gin.Context) {
	// Check if user is admin
	user, exists := c.Get("user")
	if !exists || !user.(models.User).IsAdmin {
		middleware.TranslateErrorResponse(c, http.StatusForbidden, "auth.forbidden", nil)
		return
	}

	var req models.CreateFitnessGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "common.invalid_request", err)
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
			middleware.TranslateErrorResponse(c, http.StatusConflict, "fitness.goal_name_exists", nil)
			return
		}
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "common.server_error", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": goal.ToResponse(),
	})
}

// UpdateFitnessGoal updates an existing fitness goal (admin only)
func UpdateFitnessGoal(c *gin.Context) {
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

	var req models.UpdateFitnessGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "common.invalid_request", err)
		return
	}

	var goal models.FitnessGoal
	if err := database.DB.First(&goal, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.TranslateErrorResponse(c, http.StatusNotFound, "fitness.goal_not_found", nil)
			return
		}
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "common.server_error", err)
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
			middleware.TranslateErrorResponse(c, http.StatusConflict, "fitness.goal_name_exists", nil)
			return
		}
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "common.server_error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": goal.ToResponse(),
	})
}

// DeleteFitnessGoal deletes a fitness goal (admin only)
func DeleteFitnessGoal(c *gin.Context) {
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

	result := database.DB.Delete(&models.FitnessGoal{}, id)
	if result.Error != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "common.server_error", result.Error)
		return
	}

	if result.RowsAffected == 0 {
		middleware.TranslateErrorResponse(c, http.StatusNotFound, "fitness.goal_not_found", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Fitness goal deleted successfully",
	})
}

// GetUserFitnessGoals retrieves the authenticated user's fitness goals
func GetUserFitnessGoals(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		middleware.TranslateErrorResponse(c, http.StatusUnauthorized, "auth.unauthorized", nil)
		return
	}

	currentUser := user.(models.User)

	var userGoals []models.UserFitnessGoal
	if err := database.DB.Preload("FitnessGoal").
		Where("user_id = ?", currentUser.ID).
		Order("priority ASC").
		Find(&userGoals).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "common.server_error", err)
		return
	}

	response := make([]models.UserFitnessGoalResponse, len(userGoals))
	for i, userGoal := range userGoals {
		response[i] = userGoal.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// SetUserFitnessGoals sets the authenticated user's fitness goals
func SetUserFitnessGoals(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		middleware.TranslateErrorResponse(c, http.StatusUnauthorized, "auth.unauthorized", nil)
		return
	}

	currentUser := user.(models.User)

	var req models.SetUserFitnessGoalsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "common.invalid_request", err)
		return
	}

	// Start transaction
	tx := database.DB.Begin()

	// Delete existing user goals
	if err := tx.Where("user_id = ?", currentUser.ID).Delete(&models.UserFitnessGoal{}).Error; err != nil {
		tx.Rollback()
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "common.server_error", err)
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
				middleware.TranslateErrorResponse(c, http.StatusBadRequest, "fitness.goal_not_found", nil)
				return
			}
			middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "common.server_error", err)
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
			middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "common.server_error", err)
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
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "common.server_error", err)
		return
	}

	response := make([]models.UserFitnessGoalResponse, len(updatedGoals))
	for i, userGoal := range updatedGoals {
		response[i] = userGoal.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// UpdateUserFitnessLevel updates the authenticated user's fitness level
func UpdateUserFitnessLevel(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		middleware.TranslateErrorResponse(c, http.StatusUnauthorized, "auth.unauthorized", nil)
		return
	}

	currentUser := user.(models.User)

	var req models.UpdateUserFitnessLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "common.invalid_request", err)
		return
	}

	// If fitness level ID is provided, verify it exists
	if req.FitnessLevelID != nil {
		var level models.FitnessLevel
		if err := database.DB.First(&level, *req.FitnessLevelID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				middleware.TranslateErrorResponse(c, http.StatusBadRequest, "fitness.level_not_found", nil)
				return
			}
			middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "common.server_error", err)
			return
		}
	}

	// Update user's fitness level
	if err := database.DB.Model(&currentUser).Update("fitness_level_id", req.FitnessLevelID).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "common.server_error", err)
		return
	}

	// Reload user with fitness level
	if err := database.DB.Preload("FitnessLevel").First(&currentUser, currentUser.ID).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "common.server_error", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": currentUser.ToResponse(),
	})
}