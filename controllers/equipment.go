package controllers

import (
	"lamari-fit-api/database"
	"lamari-fit-api/models"
	"lamari-fit-api/utils"

	"github.com/gin-gonic/gin"
)

// CreateEquipment creates a new equipment item
func CreateEquipment(c *gin.Context) {
	var req models.CreateEquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Check if equipment with same name or slug already exists
	var existingEquipment models.Equipment
	if err := database.DB.Where("name = ? OR slug = ?", req.Name, req.Slug).First(&existingEquipment).Error; err == nil {
		utils.ConflictResponse(c, "Equipment with this name or slug already exists")
		return
	}

	equipment := models.Equipment{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
	}

	if err := equipment.Validate(); err != nil {
		utils.BadRequestResponse(c, "Validation failed", err.Error())
		return
	}

	if err := database.DB.Create(&equipment).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create equipment")
		return
	}

	utils.CreatedResponse(c, "Equipment created successfully", equipment.ToResponse())
}

// GetAllEquipment retrieves all equipment with optional filtering
func GetAllEquipment(c *gin.Context) {
	var queryParams EquipmentQuery
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Set default pagination values
	SetDefaultPagination(&queryParams.PaginationQuery)

	var equipment []models.Equipment
	query := database.DB.Model(&models.Equipment{})

	// Search by name
	if queryParams.Search != "" {
		query = query.Where("name ILIKE ?", "%"+queryParams.Search+"%")
	}

	// Filter by category
	if queryParams.Category != "" {
		query = query.Where("category = ?", queryParams.Category)
	}

	offset := (queryParams.Page - 1) * queryParams.Limit

	// Get total count
	var total int64
	query.Count(&total)

	// Get equipment with pagination
	if err := query.Offset(offset).Limit(queryParams.Limit).Order("name").Find(&equipment).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve equipment")
		return
	}

	// Convert to response format
	responses := make([]models.EquipmentResponse, len(equipment))
	for i, eq := range equipment {
		responses[i] = eq.ToResponse()
	}

	utils.PaginatedResponse(c, "  list retrieved successfully", responses, queryParams.Page, queryParams.Limit, int(total))
}

// GetEquipmentByID retrieves a specific equipment by ID
func GetEquipmentByID(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	id, ok := utils.ParseUUID(c, params.ID, "equipment")
	if !ok {
		return
	}

	var equipment models.Equipment
	if err := database.DB.Preload("ExerciseLinks.Exercise").First(&equipment, "id = ?", id).Error; err != nil {
		utils.NotFoundResponse(c, "Equipment not found")
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

	utils.SuccessResponse(c, "Equipment retrieved successfully", response)
}

// UpdateEquipment updates an existing equipment
func UpdateEquipment(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	id, ok := utils.ParseUUID(c, params.ID, "equipment")
	if !ok {
		return
	}

	var equipment models.Equipment
	if err := database.DB.First(&equipment, "id = ?", id).Error; err != nil {
		utils.NotFoundResponse(c, "Equipment not found")
		return
	}

	var req models.UpdateEquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Update fields if provided
	if req.Name != "" {
		// Check for duplicate name
		var existingEquipment models.Equipment
		if err := database.DB.Where("name = ? AND id != ?", req.Name, id).First(&existingEquipment).Error; err == nil {
			utils.ConflictResponse(c, "Equipment with this name already exists")
			return
		}
		equipment.Name = req.Name
	}
	if req.Slug != "" {
		// Check for duplicate slug
		var existingEquipment models.Equipment
		if err := database.DB.Where("slug = ? AND id != ?", req.Slug, id).First(&existingEquipment).Error; err == nil {
			utils.ConflictResponse(c, "Equipment with this slug already exists")
			return
		}
		equipment.Slug = req.Slug
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
		utils.BadRequestResponse(c, "Validation failed", err.Error())
		return
	}

	if err := database.DB.Save(&equipment).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update equipment")
		return
	}

	utils.SuccessResponse(c, "Equipment updated successfully", equipment.ToResponse())
}

// DeleteEquipment deletes an equipment
func DeleteEquipment(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	id, ok := utils.ParseUUID(c, params.ID, "equipment")
	if !ok {
		return
	}

	var equipment models.Equipment
	if err := database.DB.First(&equipment, "id = ?", id).Error; err != nil {
		utils.NotFoundResponse(c, "Equipment not found")
		return
	}

	// Check if equipment is being used by any exercises
	var count int64
	database.DB.Model(&models.ExerciseEquipment{}).Where("equipment_id = ?", id).Count(&count)
	if count > 0 {
		utils.ConflictResponse(c, "Equipment is currently in use and cannot be deleted")
		return
	}

	if err := database.DB.Delete(&equipment).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete equipment")
		return
	}

	utils.NoContentResponse(c)
}

// AssignEquipmentToExercise assigns equipment to an exercise
func AssignEquipmentToExercise(c *gin.Context) {
	exerciseID, ok := utils.ParseUUIDParam(c, "exercise_id", "exercise")
	if !ok {
		return
	}

	var req models.AssignEquipmentToExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Check if exercise exists
	var exercise models.Exercise
	if err := database.DB.First(&exercise, "id = ?", exerciseID).Error; err != nil {
		utils.NotFoundResponse(c, "Exercise not found")
		return
	}

	// Check if equipment exists
	var equipment models.Equipment
	if err := database.DB.First(&equipment, "id = ?", req.EquipmentID).Error; err != nil {
		utils.NotFoundResponse(c, "Equipment not found")
		return
	}

	// Check if relationship already exists
	var existing models.ExerciseEquipment
	if err := database.DB.Where("exercise_id = ? AND equipment_id = ?", exerciseID, req.EquipmentID).First(&existing).Error; err == nil {
		utils.ConflictResponse(c, "Equipment is already assigned to this exercise")
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
		utils.InternalServerErrorResponse(c, "Failed to assign equipment to exercise")
		return
	}

	// Load the equipment for response
	database.DB.Preload("Equipment").First(&exerciseEquipment, "id = ?", exerciseEquipment.ID)

	utils.CreatedResponse(c, "Equipment assigned to exercise successfully", exerciseEquipment.ToResponse())
}

// RemoveEquipmentFromExercise removes equipment from an exercise
func RemoveEquipmentFromExercise(c *gin.Context) {
	exerciseID, ok := utils.ParseUUIDParam(c, "exercise_id", "exercise")
	if !ok {
		return
	}

	equipmentID, ok := utils.ParseUUIDParam(c, "equipment_id", "equipment")
	if !ok {
		return
	}

	// Find and delete the relationship
	result := database.DB.Where("exercise_id = ? AND equipment_id = ?", exerciseID, equipmentID).Delete(&models.ExerciseEquipment{})
	if result.Error != nil {
		utils.InternalServerErrorResponse(c, "Failed to remove equipment from exercise")
		return
	}

	if result.RowsAffected == 0 {
		utils.NotFoundResponse(c, "Equipment is not assigned to this exercise")
		return
	}

	utils.NoContentResponse(c)
}

// GetExerciseEquipment gets all equipment for an exercise
func GetExerciseEquipment(c *gin.Context) {
	exerciseID, ok := utils.ParseUUIDParam(c, "exercise_id", "exercise")
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

	// Check if exercise exists
	var exercise models.Exercise
	if err := database.DB.First(&exercise, "id = ?", exerciseID).Error; err != nil {
		utils.NotFoundResponse(c, "Exercise not found")
		return
	}

	// Get total count
	var total int64
	if err := database.DB.Model(&models.ExerciseEquipment{}).
		Where("exercise_id = ?", exerciseID).
		Count(&total).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to count exercise equipment")
		return
	}

	// Get all equipment for the exercise with pagination
	var exerciseEquipment []models.ExerciseEquipment
	if err := database.DB.Preload("Equipment").
		Where("exercise_id = ?", exerciseID).
		Offset(offset).
		Limit(queryParams.Limit).
		Order("created_at DESC").
		Find(&exerciseEquipment).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve exercise equipment")
		return
	}

	// Convert to response format
	responses := make([]models.ExerciseEquipmentResponse, len(exerciseEquipment))
	for i, ee := range exerciseEquipment {
		responses[i] = ee.ToResponse()
	}

	utils.PaginatedResponse(c, "Exercise equipment list retrieved successfully", responses, queryParams.Page, queryParams.Limit, int(total))
}
