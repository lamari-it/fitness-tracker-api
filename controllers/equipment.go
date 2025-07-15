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

// CreateEquipment creates a new equipment item
func CreateEquipment(c *gin.Context) {
	var req models.CreateEquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
		return
	}

	// Check if equipment with same name already exists
	var existingEquipment models.Equipment
	if err := database.DB.Where("name = ?", req.Name).First(&existingEquipment).Error; err == nil {
		middleware.TranslateErrorResponse(c, http.StatusConflict, "equipment.already_exists", nil)
		return
	}

	equipment := models.Equipment{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
	}

	if err := equipment.Validate(); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
		return
	}

	if err := database.DB.Create(&equipment).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "equipment.create_failed", nil)
		return
	}

	middleware.TranslateResponse(c, http.StatusCreated, "equipment.created", equipment.ToResponse())
}

// GetAllEquipment retrieves all equipment with optional filtering
func GetAllEquipment(c *gin.Context) {
	var equipment []models.Equipment
	query := database.DB.Model(&models.Equipment{})

	// Search by name
	if search := c.Query("search"); search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	// Filter by category
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	// Pagination
	page := 1
	limit := 10
	if p := c.Query("page"); p != "" {
		if pageNum, err := strconv.Atoi(p); err == nil && pageNum > 0 {
			page = pageNum
		}
	}
	if l := c.Query("limit"); l != "" {
		if limitNum, err := strconv.Atoi(l); err == nil && limitNum > 0 && limitNum <= 100 {
			limit = limitNum
		}
	}
	offset := (page - 1) * limit

	// Get total count
	var total int64
	query.Count(&total)

	// Get equipment with pagination
	if err := query.Offset(offset).Limit(limit).Order("name").Find(&equipment).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "general.internal_error", nil)
		return
	}

	// Convert to response format
	responses := make([]models.EquipmentResponse, len(equipment))
	for i, eq := range equipment {
		responses[i] = eq.ToResponse()
	}

	result := gin.H{
		"equipment": responses,
		"total":     total,
		"page":      page,
		"limit":     limit,
	}

	middleware.TranslateResponse(c, http.StatusOK, "equipment.list_retrieved", result)
}

// GetEquipmentByID retrieves a specific equipment by ID
func GetEquipmentByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_uuid", nil)
		return
	}

	var equipment models.Equipment
	if err := database.DB.Preload("ExerciseLinks.Exercise").First(&equipment, "id = ?", id).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusNotFound, "equipment.not_found", nil)
		return
	}

	// Count exercises
	var exerciseCount int64
	database.DB.Model(&models.ExerciseEquipment{}).Where("equipment_id = ?", id).Count(&exerciseCount)

	response := models.EquipmentWithExercises{
		EquipmentResponse: equipment.ToResponse(),
		ExerciseCount:     int(exerciseCount),
	}

	// If exercises are preloaded, include them
	if len(equipment.ExerciseLinks) > 0 {
		response.Exercises = make([]models.ExerciseEquipmentResponse, len(equipment.ExerciseLinks))
		for i, link := range equipment.ExerciseLinks {
			response.Exercises[i] = link.ToResponse()
		}
	}

	c.JSON(http.StatusOK, response)
}

// UpdateEquipment updates an existing equipment
func UpdateEquipment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_uuid", nil)
		return
	}

	var equipment models.Equipment
	if err := database.DB.First(&equipment, "id = ?", id).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusNotFound, "equipment.not_found", nil)
		return
	}

	var req models.UpdateEquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
		return
	}

	// Update fields if provided
	if req.Name != "" {
		// Check for duplicate name
		var existingEquipment models.Equipment
		if err := database.DB.Where("name = ? AND id != ?", req.Name, id).First(&existingEquipment).Error; err == nil {
			middleware.TranslateErrorResponse(c, http.StatusConflict, "equipment.already_exists", nil)
			return
		}
		equipment.Name = req.Name
	}
	if req.Description != "" {
		equipment.Description = req.Description
	}
	if req.Category != "" {
		equipment.Category = req.Category
	}
	if req.ImageURL != "" {
		equipment.ImageURL = req.ImageURL
	}

	if err := equipment.Validate(); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
		return
	}

	if err := database.DB.Save(&equipment).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "equipment.update_failed", nil)
		return
	}

	middleware.TranslateResponse(c, http.StatusOK, "equipment.updated", equipment.ToResponse())
}

// DeleteEquipment deletes an equipment
func DeleteEquipment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_uuid", nil)
		return
	}

	var equipment models.Equipment
	if err := database.DB.First(&equipment, "id = ?", id).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusNotFound, "equipment.not_found", nil)
		return
	}

	// Check if equipment is being used by any exercises
	var count int64
	database.DB.Model(&models.ExerciseEquipment{}).Where("equipment_id = ?", id).Count(&count)
	if count > 0 {
		middleware.TranslateErrorResponse(c, http.StatusConflict, "equipment.in_use", nil)
		return
	}

	if err := database.DB.Delete(&equipment).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "equipment.delete_failed", nil)
		return
	}

	c.Status(http.StatusNoContent)
}

// AssignEquipmentToExercise assigns equipment to an exercise
func AssignEquipmentToExercise(c *gin.Context) {
	exerciseIDStr := c.Param("exercise_id")
	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_uuid", nil)
		return
	}

	var req models.AssignEquipmentToExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
		return
	}

	// Check if exercise exists
	var exercise models.Exercise
	if err := database.DB.First(&exercise, "id = ?", exerciseID).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusNotFound, "exercise.not_found", nil)
		return
	}

	// Check if equipment exists
	var equipment models.Equipment
	if err := database.DB.First(&equipment, "id = ?", req.EquipmentID).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusNotFound, "equipment.not_found", nil)
		return
	}

	// Check if relationship already exists
	var existing models.ExerciseEquipment
	if err := database.DB.Where("exercise_id = ? AND equipment_id = ?", exerciseID, req.EquipmentID).First(&existing).Error; err == nil {
		middleware.TranslateErrorResponse(c, http.StatusConflict, "equipment.already_assigned", nil)
		return
	}

	// Create the relationship
	exerciseEquipment := models.ExerciseEquipment{
		ExerciseID:  exerciseID,
		EquipmentID: req.EquipmentID,
		Optional:    req.Optional,
		Notes:       req.Notes,
	}

	if err := database.DB.Create(&exerciseEquipment).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "equipment.assign_failed", nil)
		return
	}

	// Load the equipment for response
	database.DB.Preload("Equipment").First(&exerciseEquipment, "id = ?", exerciseEquipment.ID)

	middleware.TranslateResponse(c, http.StatusCreated, "equipment.assigned", exerciseEquipment.ToResponse())
}

// RemoveEquipmentFromExercise removes equipment from an exercise
func RemoveEquipmentFromExercise(c *gin.Context) {
	exerciseIDStr := c.Param("exercise_id")
	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_uuid", nil)
		return
	}

	equipmentIDStr := c.Param("equipment_id")
	equipmentID, err := uuid.Parse(equipmentIDStr)
	if err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_uuid", nil)
		return
	}

	// Find and delete the relationship
	result := database.DB.Where("exercise_id = ? AND equipment_id = ?", exerciseID, equipmentID).Delete(&models.ExerciseEquipment{})
	if result.Error != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "equipment.remove_failed", nil)
		return
	}

	if result.RowsAffected == 0 {
		middleware.TranslateErrorResponse(c, http.StatusNotFound, "equipment.not_assigned", nil)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetExerciseEquipment gets all equipment for an exercise
func GetExerciseEquipment(c *gin.Context) {
	exerciseIDStr := c.Param("exercise_id")
	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_uuid", nil)
		return
	}

	// Check if exercise exists
	var exercise models.Exercise
	if err := database.DB.First(&exercise, "id = ?", exerciseID).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusNotFound, "exercise.not_found", nil)
		return
	}

	// Get all equipment for the exercise
	var exerciseEquipment []models.ExerciseEquipment
	if err := database.DB.Preload("Equipment").Where("exercise_id = ?", exerciseID).Find(&exerciseEquipment).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "general.internal_error", nil)
		return
	}

	// Convert to response format
	responses := make([]models.ExerciseEquipmentResponse, len(exerciseEquipment))
	for i, ee := range exerciseEquipment {
		responses[i] = ee.ToResponse()
	}

	middleware.TranslateResponse(c, http.StatusOK, "equipment.list_retrieved", gin.H{
		"equipment": responses,
	})
}