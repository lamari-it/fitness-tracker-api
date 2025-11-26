package controllers

import (
	"lamari-fit-api/database"
	"lamari-fit-api/models"
	"lamari-fit-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateWorkoutPlanRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=200"`
	Description string `json:"description" binding:"omitempty,max=1000"`
	Visibility  string `json:"visibility" binding:"omitempty,oneof=private public"`
}

type UpdateWorkoutPlanRequest struct {
	Title       string `json:"title" binding:"omitempty,min=1,max=200"`
	Description string `json:"description" binding:"omitempty,max=1000"`
	Visibility  string `json:"visibility" binding:"omitempty,oneof=private public"`
}

func CreateWorkoutPlan(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var req CreateWorkoutPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	plan := models.WorkoutPlan{
		UserID:      userID,
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
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var queryParams PaginationQuery
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Set default pagination values
	SetDefaultPagination(&queryParams)

	offset := (queryParams.Page - 1) * queryParams.Limit

	var plans []models.WorkoutPlan
	var total int64

	if err := database.DB.Model(&models.WorkoutPlan{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to count workout plans.")
		return
	}

	if err := database.DB.Where("user_id = ?", userID).
		Preload("Items").
		Preload("Items.Workout").
		Offset(offset).
		Limit(queryParams.Limit).
		Order("created_at DESC").
		Find(&plans).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch workout plans.")
		return
	}

	utils.PaginatedResponse(c, "Workout plans fetched successfully.", plans, queryParams.Page, queryParams.Limit, int(total))
}

func GetWorkoutPlan(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	planID, ok := utils.ParseUUID(c, params.ID, "workout plan")
	if !ok {
		return
	}

	var plan models.WorkoutPlan
	if err := database.DB.Where("id = ? AND user_id = ?", planID, userID).
		Preload("Items").
		Preload("Items.Workout").
		Preload("Items.Workout.SetGroups").
		Preload("Items.Workout.WorkoutExercises").
		Preload("Items.Workout.WorkoutExercises.Exercise").
		First(&plan).Error; err != nil {
		utils.NotFoundResponse(c, "Workout plan not found.")
		return
	}

	utils.SuccessResponse(c, "Workout plan fetched successfully.", plan)
}

func UpdateWorkoutPlan(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	planID, ok := utils.ParseUUID(c, params.ID, "workout plan")
	if !ok {
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
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	planID, ok := utils.ParseUUID(c, params.ID, "workout plan")
	if !ok {
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

type AddWorkoutToPlanRequest struct {
	WorkoutID uuid.UUID `json:"workout_id" binding:"required,uuid"`
	WeekIndex int       `json:"week_index" binding:"omitempty,min=0,max=52"`
}

func AddWorkoutToPlan(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	planID, ok := utils.ParseUUID(c, params.ID, "workout plan")
	if !ok {
		return
	}

	var req AddWorkoutToPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Check if plan exists and belongs to user
	var plan models.WorkoutPlan
	if err := database.DB.Where("id = ? AND user_id = ?", planID, userID).First(&plan).Error; err != nil {
		utils.NotFoundResponse(c, "Workout plan not found.")
		return
	}

	// Check if workout exists and belongs to user
	var workout models.Workout
	if err := database.DB.Where("id = ? AND user_id = ?", req.WorkoutID, userID).First(&workout).Error; err != nil {
		utils.NotFoundResponse(c, "Workout not found.")
		return
	}

	// Check if this workout is already in the plan
	var existingItem models.WorkoutPlanItem
	if err := database.DB.Where("plan_id = ? AND workout_id = ?", planID, req.WorkoutID).First(&existingItem).Error; err == nil {
		utils.BadRequestResponse(c, "This workout is already in the plan.", nil)
		return
	}

	// Create the workout plan item
	planItem := models.WorkoutPlanItem{
		PlanID:    planID,
		WorkoutID: req.WorkoutID,
		WeekIndex: req.WeekIndex,
	}

	if err := database.DB.Create(&planItem).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to add workout to plan.")
		return
	}

	// Load the workout details
	database.DB.Preload("Workout").First(&planItem, planItem.ID)

	utils.CreatedResponse(c, "Workout added to plan successfully.", planItem)
}

func RemoveWorkoutFromPlan(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var params struct {
		ID     string `uri:"id" binding:"required,uuid"`
		ItemID string `uri:"item_id" binding:"required,uuid"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	planID, ok := utils.ParseUUID(c, params.ID, "workout plan")
	if !ok {
		return
	}

	itemID, ok := utils.ParseUUID(c, params.ItemID, "plan item")
	if !ok {
		return
	}

	// Check if plan exists and belongs to user
	var plan models.WorkoutPlan
	if err := database.DB.Where("id = ? AND user_id = ?", planID, userID).First(&plan).Error; err != nil {
		utils.NotFoundResponse(c, "Workout plan not found.")
		return
	}

	// Delete the workout plan item
	result := database.DB.Where("id = ? AND plan_id = ?", itemID, planID).Delete(&models.WorkoutPlanItem{})
	if result.Error != nil {
		utils.InternalServerErrorResponse(c, "Failed to remove workout from plan.")
		return
	}

	if result.RowsAffected == 0 {
		utils.NotFoundResponse(c, "Workout item not found in this plan.")
		return
	}

	utils.DeletedResponse(c, "Workout removed from plan successfully.")
}

func GetPlanWorkouts(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	planID, ok := utils.ParseUUID(c, params.ID, "workout plan")
	if !ok {
		return
	}

	// Check if plan exists and belongs to user
	var plan models.WorkoutPlan
	if err := database.DB.Where("id = ? AND user_id = ?", planID, userID).First(&plan).Error; err != nil {
		utils.NotFoundResponse(c, "Workout plan not found.")
		return
	}

	var planItems []models.WorkoutPlanItem
	if err := database.DB.Where("plan_id = ?", planID).
		Preload("Workout.SetGroups").
		Preload("Workout.WorkoutExercises.Exercise").
		Order("week_index, created_at").
		Find(&planItems).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch plan workouts.")
		return
	}

	utils.SuccessResponse(c, "Plan workouts fetched successfully.", planItems)
}

type UpdatePlanItemRequest struct {
	WeekIndex int `json:"week_index" binding:"min=0,max=52"`
}

func UpdatePlanItem(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var params struct {
		ID     string `uri:"id" binding:"required,uuid"`
		ItemID string `uri:"item_id" binding:"required,uuid"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	planID, ok := utils.ParseUUID(c, params.ID, "workout plan")
	if !ok {
		return
	}

	itemID, ok := utils.ParseUUID(c, params.ItemID, "plan item")
	if !ok {
		return
	}

	var req UpdatePlanItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Check if plan exists and belongs to user
	var plan models.WorkoutPlan
	if err := database.DB.Where("id = ? AND user_id = ?", planID, userID).First(&plan).Error; err != nil {
		utils.NotFoundResponse(c, "Workout plan not found.")
		return
	}

	// Find and update the plan item
	var planItem models.WorkoutPlanItem
	if err := database.DB.Where("id = ? AND plan_id = ?", itemID, planID).First(&planItem).Error; err != nil {
		utils.NotFoundResponse(c, "Workout item not found in this plan.")
		return
	}

	planItem.WeekIndex = req.WeekIndex
	if err := database.DB.Save(&planItem).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update plan item.")
		return
	}

	// Load the workout details
	database.DB.Preload("Workout").First(&planItem, planItem.ID)

	utils.SuccessResponse(c, "Plan item updated successfully.", planItem)
}
