package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/middleware"
	"fit-flow-api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateMuscleGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

type UpdateMuscleGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

type AssignMuscleGroupRequest struct {
	MuscleGroupID uuid.UUID `json:"muscle_group_id" binding:"required"`
	Primary       bool      `json:"primary"`
	Intensity     string    `json:"intensity"`
}

func CreateMuscleGroup(c *gin.Context) {
	var req CreateMuscleGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
		return
	}

	muscleGroup := models.MuscleGroup{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
	}

	if err := muscleGroup.Validate(); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
		return
	}

	if err := database.DB.Create(&muscleGroup).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create muscle group"})
		return
	}

	c.JSON(http.StatusCreated, muscleGroup.ToResponse())
}

func GetMuscleGroups(c *gin.Context) {
	search := c.Query("search")
	category := c.Query("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	query := database.DB.Model(&models.MuscleGroup{})

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	var muscleGroups []models.MuscleGroup
	var total int64

	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).Order("name ASC").Find(&muscleGroups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch muscle groups"})
		return
	}

	var responses []models.MuscleGroupResponse
	for _, mg := range muscleGroups {
		responses = append(responses, mg.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"muscle_groups": responses,
		"total":         total,
		"page":          page,
		"limit":         limit,
	})
}

func GetMuscleGroup(c *gin.Context) {
	muscleGroupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid muscle group ID"})
		return
	}

	var muscleGroup models.MuscleGroup
	if err := database.DB.Where("id = ?", muscleGroupID).
		Preload("ExerciseLinks.Exercise").
		First(&muscleGroup).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Muscle group not found"})
		return
	}

	response := models.MuscleGroupWithExercises{
		MuscleGroupResponse: muscleGroup.ToResponse(),
		ExerciseCount:       len(muscleGroup.ExerciseLinks),
	}

	for _, link := range muscleGroup.ExerciseLinks {
		exerciseResponse := models.ExerciseMuscleGroupResponse{
			ID:            link.ID,
			ExerciseID:    link.ExerciseID,
			MuscleGroupID: link.MuscleGroupID,
			Primary:       link.Primary,
			Intensity:     link.Intensity,
			MuscleGroup:   muscleGroup.ToResponse(),
		}
		response.Exercises = append(response.Exercises, exerciseResponse)
	}

	c.JSON(http.StatusOK, response)
}

func UpdateMuscleGroup(c *gin.Context) {
	muscleGroupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid muscle group ID"})
		return
	}

	var req UpdateMuscleGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
		return
	}

	var muscleGroup models.MuscleGroup
	if err := database.DB.Where("id = ?", muscleGroupID).First(&muscleGroup).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Muscle group not found"})
		return
	}

	if req.Name != "" {
		muscleGroup.Name = req.Name
	}
	if req.Description != "" {
		muscleGroup.Description = req.Description
	}
	if req.Category != "" {
		muscleGroup.Category = req.Category
	}

	if err := muscleGroup.Validate(); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
		return
	}

	if err := database.DB.Save(&muscleGroup).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update muscle group"})
		return
	}

	c.JSON(http.StatusOK, muscleGroup.ToResponse())
}

func DeleteMuscleGroup(c *gin.Context) {
	muscleGroupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid muscle group ID"})
		return
	}

	// Check if muscle group is being used by any exercises
	var count int64
	database.DB.Model(&models.ExerciseMuscleGroup{}).Where("muscle_group_id = ?", muscleGroupID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete muscle group that is assigned to exercises"})
		return
	}

	result := database.DB.Where("id = ?", muscleGroupID).Delete(&models.MuscleGroup{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete muscle group"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Muscle group not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Muscle group deleted successfully"})
}

func AssignMuscleGroupToExercise(c *gin.Context) {
	exerciseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	var req AssignMuscleGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
		return
	}

	// Check if exercise exists
	var exercise models.Exercise
	if err := database.DB.Where("id = ?", exerciseID).First(&exercise).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}

	// Check if muscle group exists
	var muscleGroup models.MuscleGroup
	if err := database.DB.Where("id = ?", req.MuscleGroupID).First(&muscleGroup).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Muscle group not found"})
		return
	}

	// Check if assignment already exists
	var existing models.ExerciseMuscleGroup
	if err := database.DB.Where("exercise_id = ? AND muscle_group_id = ?", exerciseID, req.MuscleGroupID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Muscle group already assigned to this exercise"})
		return
	}

	// If this is being set as primary, unset any existing primary muscle groups
	if req.Primary {
		database.DB.Model(&models.ExerciseMuscleGroup{}).
			Where("exercise_id = ? AND primary = true", exerciseID).
			Update("primary", false)
	}

	assignment := models.ExerciseMuscleGroup{
		ExerciseID:    exerciseID,
		MuscleGroupID: req.MuscleGroupID,
		Primary:       req.Primary,
		Intensity:     req.Intensity,
	}

	if assignment.Intensity == "" {
		assignment.Intensity = "moderate"
	}

	if err := assignment.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment data"})
		return
	}

	if err := database.DB.Create(&assignment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign muscle group to exercise"})
		return
	}

	// Load the muscle group for response
	database.DB.Where("id = ?", req.MuscleGroupID).First(&assignment.MuscleGroup)

	c.JSON(http.StatusCreated, assignment.ToResponse())
}

func RemoveMuscleGroupFromExercise(c *gin.Context) {
	exerciseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	muscleGroupID, err := uuid.Parse(c.Param("muscle_group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid muscle group ID"})
		return
	}

	result := database.DB.Where("exercise_id = ? AND muscle_group_id = ?", exerciseID, muscleGroupID).Delete(&models.ExerciseMuscleGroup{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove muscle group from exercise"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Muscle group assignment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Muscle group removed from exercise successfully"})
}

func GetExerciseMuscleGroups(c *gin.Context) {
	exerciseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	var assignments []models.ExerciseMuscleGroup
	if err := database.DB.Where("exercise_id = ?", exerciseID).
		Preload("MuscleGroup").
		Find(&assignments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exercise muscle groups"})
		return
	}

	var responses []models.ExerciseMuscleGroupResponse
	for _, assignment := range assignments {
		responses = append(responses, assignment.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{"muscle_groups": responses})
}