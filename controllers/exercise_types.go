package controllers

import (
	"lamari-fit-api/database"
	"lamari-fit-api/models"
	"lamari-fit-api/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetExerciseTypes returns all exercise types
func GetExerciseTypes(c *gin.Context) {
	var exerciseTypes []models.ExerciseType
	if err := database.DB.Order("name ASC").Find(&exerciseTypes).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch exercise types.")
		return
	}

	responses := make([]models.ExerciseTypeResponse, len(exerciseTypes))
	for i, et := range exerciseTypes {
		responses[i] = et.ToResponse()
	}

	utils.SuccessResponse(c, "Exercise types fetched successfully.", responses)
}

// GetExerciseType returns a single exercise type by ID
func GetExerciseType(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	id, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid exercise type ID format", nil)
		return
	}

	var exerciseType models.ExerciseType
	if err := database.DB.First(&exerciseType, "id = ?", id).Error; err != nil {
		utils.NotFoundResponse(c, "Exercise type not found.")
		return
	}

	utils.SuccessResponse(c, "Exercise type fetched successfully.", exerciseType.ToResponse())
}

// CreateExerciseType creates a new exercise type (admin only)
func CreateExerciseType(c *gin.Context) {
	var req models.CreateExerciseTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Generate slug from name
	slug := strings.ToLower(strings.ReplaceAll(req.Name, " ", "_"))

	// Check if exercise type with same name or slug exists
	var existing models.ExerciseType
	if err := database.DB.Where("slug = ? OR name = ?", slug, req.Name).First(&existing).Error; err == nil {
		utils.BadRequestResponse(c, "Exercise type with this name already exists.", nil)
		return
	}

	exerciseType := models.ExerciseType{
		Slug:        slug,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := database.DB.Create(&exerciseType).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create exercise type.")
		return
	}

	utils.CreatedResponse(c, "Exercise type created successfully.", exerciseType.ToResponse())
}

// UpdateExerciseType updates an existing exercise type (admin only)
func UpdateExerciseType(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	id, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid exercise type ID format", nil)
		return
	}

	var exerciseType models.ExerciseType
	if err := database.DB.First(&exerciseType, "id = ?", id).Error; err != nil {
		utils.NotFoundResponse(c, "Exercise type not found.")
		return
	}

	var req models.UpdateExerciseTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
		updates["slug"] = strings.ToLower(strings.ReplaceAll(req.Name, " ", "_"))
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	if len(updates) == 0 {
		utils.BadRequestResponse(c, "No updates provided.", nil)
		return
	}

	if err := database.DB.Model(&exerciseType).Updates(updates).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update exercise type.")
		return
	}

	// Reload the exercise type
	database.DB.First(&exerciseType, "id = ?", id)

	utils.SuccessResponse(c, "Exercise type updated successfully.", exerciseType.ToResponse())
}

// DeleteExerciseType deletes an exercise type (admin only)
func DeleteExerciseType(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	id, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid exercise type ID format", nil)
		return
	}

	var exerciseType models.ExerciseType
	if err := database.DB.First(&exerciseType, "id = ?", id).Error; err != nil {
		utils.NotFoundResponse(c, "Exercise type not found.")
		return
	}

	if err := database.DB.Delete(&exerciseType).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete exercise type.")
		return
	}

	utils.SuccessResponse(c, "Exercise type deleted successfully.", nil)
}

// GetExerciseTypesByExercise returns all exercise types for an exercise
func GetExerciseTypesByExercise(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	exerciseID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid exercise ID format", nil)
		return
	}

	// Verify exercise exists
	var exercise models.Exercise
	if err := database.DB.First(&exercise, "id = ?", exerciseID).Error; err != nil {
		utils.NotFoundResponse(c, "Exercise not found.")
		return
	}

	var exerciseExerciseTypes []models.ExerciseExerciseType
	if err := database.DB.Preload("ExerciseType").Where("exercise_id = ?", exerciseID).Find(&exerciseExerciseTypes).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch exercise types.")
		return
	}

	responses := make([]models.ExerciseExerciseTypeResponse, len(exerciseExerciseTypes))
	for i, eet := range exerciseExerciseTypes {
		responses[i] = eet.ToResponse()
	}

	utils.SuccessResponse(c, "Exercise types fetched successfully.", responses)
}

// AssignExerciseType assigns an exercise type to an exercise
func AssignExerciseType(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	exerciseID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid exercise ID format", nil)
		return
	}

	// Verify exercise exists
	var exercise models.Exercise
	if err := database.DB.First(&exercise, "id = ?", exerciseID).Error; err != nil {
		utils.NotFoundResponse(c, "Exercise not found.")
		return
	}

	var req models.AssignExerciseTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Verify exercise type exists
	var exerciseType models.ExerciseType
	if err := database.DB.First(&exerciseType, "id = ?", req.ExerciseTypeID).Error; err != nil {
		utils.NotFoundResponse(c, "Exercise type not found.")
		return
	}

	// Check if already assigned
	var existing models.ExerciseExerciseType
	if err := database.DB.Where("exercise_id = ? AND exercise_type_id = ?", exerciseID, req.ExerciseTypeID).First(&existing).Error; err == nil {
		utils.BadRequestResponse(c, "Exercise type already assigned to this exercise.", nil)
		return
	}

	assignment := models.ExerciseExerciseType{
		ExerciseID:     exerciseID,
		ExerciseTypeID: req.ExerciseTypeID,
	}

	if err := database.DB.Create(&assignment).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to assign exercise type.")
		return
	}

	// Load the exercise type for response
	database.DB.Preload("ExerciseType").First(&assignment, "id = ?", assignment.ID)

	utils.CreatedResponse(c, "Exercise type assigned successfully.", assignment.ToResponse())
}

// RemoveExerciseType removes an exercise type from an exercise
func RemoveExerciseType(c *gin.Context) {
	exerciseIDStr := c.Param("id")
	typeIDStr := c.Param("type_id")

	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid exercise ID format", nil)
		return
	}

	typeID, err := uuid.Parse(typeIDStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid exercise type ID format", nil)
		return
	}

	var assignment models.ExerciseExerciseType
	if err := database.DB.Where("exercise_id = ? AND exercise_type_id = ?", exerciseID, typeID).First(&assignment).Error; err != nil {
		utils.NotFoundResponse(c, "Exercise type assignment not found.")
		return
	}

	if err := database.DB.Delete(&assignment).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to remove exercise type.")
		return
	}

	utils.SuccessResponse(c, "Exercise type removed successfully.", nil)
}
