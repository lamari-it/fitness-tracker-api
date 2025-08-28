package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"
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
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	var req CreateWorkoutPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
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
		utils.InternalServerErrorResponse(c, "Failed to create workout plan.")
		return
	}

	utils.CreatedResponse(c, "Workout plan created successfully.", plan)
}

func GetWorkoutPlans(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var plans []models.WorkoutPlan
	var total int64

	if err := database.DB.Model(&models.WorkoutPlan{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to count workout plans.")
		return
	}

	if err := database.DB.Where("user_id = ?", userID).
		Preload("Workouts").
		Preload("User").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&plans).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch workout plans.")
		return
	}

	utils.PaginatedResponse(c, "Workout plans fetched successfully.", plans, page, limit, int(total))
}

func GetWorkoutPlan(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		validationErrors := utils.ValidationErrors{
			"id": []string{"Invalid workout plan ID format."},
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	var plan models.WorkoutPlan
	if err := database.DB.Where("id = ? AND user_id = ?", planID, userID).
		Preload("Workouts.SetGroups.WorkoutExercises.Exercise").
		Preload("Workouts.SetGroups.WorkoutExercises.SetGroup").
		First(&plan).Error; err != nil {
		utils.NotFoundResponse(c, "Workout plan not found.")
		return
	}

	utils.SuccessResponse(c, "Workout plan fetched successfully.", plan)
}

func UpdateWorkoutPlan(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		validationErrors := utils.ValidationErrors{
			"id": []string{"Invalid workout plan ID format."},
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	var req UpdateWorkoutPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	var plan models.WorkoutPlan
	if err := database.DB.Where("id = ? AND user_id = ?", planID, userID).First(&plan).Error; err != nil {
		utils.NotFoundResponse(c, "Workout plan not found.")
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Visibility != "" {
		updates["visibility"] = req.Visibility
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&plan).Updates(updates).Error; err != nil {
			utils.InternalServerErrorResponse(c, "Failed to update workout plan.")
			return
		}
	}

	database.DB.Where("id = ?", planID).First(&plan)

	utils.SuccessResponse(c, "Workout plan updated successfully.", plan)
}

func DeleteWorkoutPlan(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		validationErrors := utils.ValidationErrors{
			"id": []string{"Invalid workout plan ID format."},
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	var plan models.WorkoutPlan
	if err := database.DB.Where("id = ? AND user_id = ?", planID, userID).First(&plan).Error; err != nil {
		utils.NotFoundResponse(c, "Workout plan not found.")
		return
	}

	if err := database.DB.Delete(&plan).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete workout plan.")
		return
	}

	utils.DeletedResponse(c, "Workout plan deleted successfully.")
}
