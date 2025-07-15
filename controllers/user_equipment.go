package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/middleware"
	"fit-flow-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetUserEquipment gets all equipment for the authenticated user
func GetUserEquipment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		middleware.TranslateErrorResponse(c, http.StatusUnauthorized, "auth.token_invalid", nil)
		return
	}

	// Get filter parameters
	var filter models.UserEquipmentFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
		return
	}

	query := database.DB.Where("user_id = ?", userID).Preload("Equipment")
	
	// Apply location filter if provided
	if filter.LocationType != "" {
		query = query.Where("location_type = ?", filter.LocationType)
	}

	var userEquipment []models.UserEquipment
	if err := query.Find(&userEquipment).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "user_equipment.fetch_failed", nil)
		return
	}

	// Convert to response format
	response := make([]models.UserEquipmentResponse, len(userEquipment))
	for i, ue := range userEquipment {
		response[i] = ue.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"equipment": response,
		"total":     len(response),
	})
}

// AddUserEquipment adds equipment to user's inventory
func AddUserEquipment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		middleware.TranslateErrorResponse(c, http.StatusUnauthorized, "auth.token_invalid", nil)
		return
	}

	var req models.AddUserEquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
		return
	}

	// Check if equipment exists
	var equipment models.Equipment
	if err := database.DB.First(&equipment, "id = ?", req.EquipmentID).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusNotFound, "equipment.not_found", nil)
		return
	}

	// Check if user already has this equipment at this location
	var existingUE models.UserEquipment
	err := database.DB.Where("user_id = ? AND equipment_id = ? AND location_type = ?", 
		userID, req.EquipmentID, req.LocationType).First(&existingUE).Error
	
	if err == nil {
		middleware.TranslateErrorResponse(c, http.StatusConflict, "user_equipment.already_exists", nil)
		return
	}

	// Create user equipment entry
	userEquipment := models.UserEquipment{
		UserID:       userID.(uuid.UUID),
		EquipmentID:  req.EquipmentID,
		LocationType: req.LocationType,
		GymLocation:  req.GymLocation,
		Notes:        req.Notes,
	}

	if err := userEquipment.Validate(); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
		return
	}

	if err := database.DB.Create(&userEquipment).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "user_equipment.create_failed", nil)
		return
	}

	// Load equipment for response
	database.DB.Preload("Equipment").First(&userEquipment, userEquipment.ID)

	middleware.TranslateResponse(c, http.StatusCreated, "user_equipment.created", userEquipment.ToResponse())
}

// UpdateUserEquipment updates user's equipment details
func UpdateUserEquipment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		middleware.TranslateErrorResponse(c, http.StatusUnauthorized, "auth.token_invalid", nil)
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_uuid", nil)
		return
	}

	// Find user equipment
	var userEquipment models.UserEquipment
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&userEquipment).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusNotFound, "user_equipment.not_found", nil)
		return
	}

	var req models.UpdateUserEquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
		return
	}

	// Update fields if provided
	if req.LocationType != "" {
		// Check if changing location would create a duplicate
		var existingUE models.UserEquipment
		err := database.DB.Where("user_id = ? AND equipment_id = ? AND location_type = ? AND id != ?", 
			userID, userEquipment.EquipmentID, req.LocationType, id).First(&existingUE).Error
		
		if err == nil {
			middleware.TranslateErrorResponse(c, http.StatusConflict, "user_equipment.already_exists", nil)
			return
		}
		userEquipment.LocationType = req.LocationType
	}
	
	if req.GymLocation != "" {
		userEquipment.GymLocation = req.GymLocation
	}
	
	if req.Notes != "" {
		userEquipment.Notes = req.Notes
	}

	if err := userEquipment.Validate(); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
		return
	}

	if err := database.DB.Save(&userEquipment).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "user_equipment.update_failed", nil)
		return
	}

	// Load equipment for response
	database.DB.Preload("Equipment").First(&userEquipment, userEquipment.ID)

	middleware.TranslateResponse(c, http.StatusOK, "user_equipment.updated", userEquipment.ToResponse())
}

// RemoveUserEquipment removes equipment from user's inventory
func RemoveUserEquipment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		middleware.TranslateErrorResponse(c, http.StatusUnauthorized, "auth.token_invalid", nil)
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_uuid", nil)
		return
	}

	// Delete user equipment
	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.UserEquipment{})
	if result.Error != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "user_equipment.delete_failed", nil)
		return
	}

	if result.RowsAffected == 0 {
		middleware.TranslateErrorResponse(c, http.StatusNotFound, "user_equipment.not_found", nil)
		return
	}

	middleware.TranslateResponse(c, http.StatusOK, "user_equipment.deleted", nil)
}

// GetUserEquipmentByLocation gets equipment filtered by location
func GetUserEquipmentByLocation(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		middleware.TranslateErrorResponse(c, http.StatusUnauthorized, "auth.token_invalid", nil)
		return
	}

	locationType := c.Param("location")
	if locationType != "home" && locationType != "gym" {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_location", nil)
		return
	}

	var userEquipment []models.UserEquipment
	if err := database.DB.Where("user_id = ? AND location_type = ?", userID, locationType).
		Preload("Equipment").Find(&userEquipment).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "user_equipment.fetch_failed", nil)
		return
	}

	// Group by category for better organization
	categoryMap := make(map[string][]models.UserEquipmentResponse)
	for _, ue := range userEquipment {
		response := ue.ToResponse()
		category := response.Equipment.Category
		if category == "" {
			category = "other"
		}
		categoryMap[category] = append(categoryMap[category], response)
	}

	c.JSON(http.StatusOK, gin.H{
		"location":           locationType,
		"equipment_by_category": categoryMap,
		"total":              len(userEquipment),
	})
}

// BulkAddUserEquipment allows adding multiple equipment items at once
func BulkAddUserEquipment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		middleware.TranslateErrorResponse(c, http.StatusUnauthorized, "auth.token_invalid", nil)
		return
	}

	var req struct {
		Equipment []models.AddUserEquipmentRequest `json:"equipment" binding:"required,dive"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
		return
	}

	// Validate all equipment exists
	var equipmentIDs []uuid.UUID
	for _, item := range req.Equipment {
		equipmentIDs = append(equipmentIDs, item.EquipmentID)
	}

	var count int64
	database.DB.Model(&models.Equipment{}).Where("id IN ?", equipmentIDs).Count(&count)
	if int(count) != len(equipmentIDs) {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "equipment.some_not_found", nil)
		return
	}

	// Start transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var created []models.UserEquipment
	for _, item := range req.Equipment {
		// Check if already exists
		var existing models.UserEquipment
		err := tx.Where("user_id = ? AND equipment_id = ? AND location_type = ?", 
			userID, item.EquipmentID, item.LocationType).First(&existing).Error
		
		if err == nil {
			continue // Skip if already exists
		}

		userEquipment := models.UserEquipment{
			UserID:       userID.(uuid.UUID),
			EquipmentID:  item.EquipmentID,
			LocationType: item.LocationType,
			GymLocation:  item.GymLocation,
			Notes:        item.Notes,
		}

		if err := tx.Create(&userEquipment).Error; err != nil {
			tx.Rollback()
			middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "user_equipment.bulk_create_failed", nil)
			return
		}
		created = append(created, userEquipment)
	}

	if err := tx.Commit().Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "user_equipment.bulk_create_failed", nil)
		return
	}

	// Load equipment for response
	database.DB.Where("id IN ?", created).Preload("Equipment").Find(&created)

	response := make([]models.UserEquipmentResponse, len(created))
	for i, ue := range created {
		response[i] = ue.ToResponse()
	}

	c.JSON(http.StatusCreated, gin.H{
		"created": response,
		"total":   len(response),
	})
}