package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateWorkoutPlanRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Visibility  string `json:"visibility"`
}

type UpdateWorkoutPlanRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Visibility  string `json:"visibility"`
}

func CreateWorkoutPlan(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req CreateWorkoutPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan := models.WorkoutPlan{
		UserID:      userID.(uuid.UUID),
		Title:       req.Title,
		Description: req.Description,
		Visibility:  req.Visibility,
	}

	if plan.Visibility == "" {
		plan.Visibility = "private"
	}

	if err := database.DB.Create(&plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workout plan"})
		return
	}

	c.JSON(http.StatusCreated, plan)
}

func GetWorkoutPlans(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var plans []models.WorkoutPlan
	var total int64

	database.DB.Model(&models.WorkoutPlan{}).Where("user_id = ?", userID).Count(&total)
	
	if err := database.DB.Where("user_id = ?", userID).
		Preload("Workouts").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&plans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch workout plans"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"plans": plans,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func GetWorkoutPlan(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	var plan models.WorkoutPlan
	if err := database.DB.Where("id = ? AND user_id = ?", planID, userID).
		Preload("Workouts.WorkoutExercises.Exercise").
		First(&plan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout plan not found"})
		return
	}

	c.JSON(http.StatusOK, plan)
}

func UpdateWorkoutPlan(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	var req UpdateWorkoutPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var plan models.WorkoutPlan
	if err := database.DB.Where("id = ? AND user_id = ?", planID, userID).First(&plan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout plan not found"})
		return
	}

	if req.Title != "" {
		plan.Title = req.Title
	}
	if req.Description != "" {
		plan.Description = req.Description
	}
	if req.Visibility != "" {
		plan.Visibility = req.Visibility
	}

	if err := database.DB.Save(&plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update workout plan"})
		return
	}

	c.JSON(http.StatusOK, plan)
}

func DeleteWorkoutPlan(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	result := database.DB.Where("id = ? AND user_id = ?", planID, userID).Delete(&models.WorkoutPlan{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete workout plan"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout plan not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workout plan deleted successfully"})
}