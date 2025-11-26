package controllers

import (
	"lamari-fit-api/database"
	"lamari-fit-api/models"
	"lamari-fit-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetUserEquipment gets all equipment for the authenticated user
func GetUserEquipment(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var queryParams UserEquipmentQuery
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Set default pagination values
	SetDefaultPagination(&queryParams.PaginationQuery)

	offset := (queryParams.Page - 1) * queryParams.Limit

	query := database.DB.Model(&models.UserEquipment{}).Where("user_id = ?", userID).Preload("Equipment")

	// Apply location filter if provided
	if queryParams.LocationType != "" {
		query = query.Where("location_type = ?", queryParams.LocationType)
	}

	// Get total count
	var total int64
	query.Count(&total)

	var userEquipment []models.UserEquipment
	if err := query.Offset(offset).Limit(queryParams.Limit).Order("created_at DESC").Find(&userEquipment).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve user equipment.")
		return
	}

	// Convert to response format
	response := make([]models.UserEquipmentResponse, len(userEquipment))
	for i, ue := range userEquipment {
		response[i] = ue.ToResponse()
	}

	utils.PaginatedResponse(c, "User equipment retrieved successfully.", response, queryParams.Page, queryParams.Limit, int(total))
}

// AddUserEquipment adds equipment to user's inventory
func AddUserEquipment(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var req models.AddUserEquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Check if equipment exists
	var equipment models.Equipment
	if err := database.DB.First(&equipment, "id = ?", req.EquipmentID).Error; err != nil {
		utils.NotFoundResponse(c, "Equipment not found.")
		return
	}

	// Check if user already has this equipment at this location
	var existingUE models.UserEquipment
	err := database.DB.Where("user_id = ? AND equipment_id = ? AND location_type = ?",
		userID, req.EquipmentID, req.LocationType).First(&existingUE).Error

	if err == nil {
		utils.ConflictResponse(c, "You already have this equipment at this location.")
		return
	}

	// Create user equipment entry
	userEquipment := models.UserEquipment{
		UserID:       userID,
		EquipmentID:  req.EquipmentID,
		LocationType: req.LocationType,
		GymLocation:  req.GymLocation,
		Notes:        req.Notes,
	}

	if err := userEquipment.Validate(); err != nil {
		utils.BadRequestResponse(c, "Validation failed.", err.Error())
		return
	}

	if err := database.DB.Create(&userEquipment).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to add equipment.")
		return
	}

	// Load equipment for response
	database.DB.Preload("Equipment").First(&userEquipment, userEquipment.ID)

	utils.CreatedResponse(c, "Equipment added successfully.", userEquipment.ToResponse())
}

// UpdateUserEquipment updates user's equipment details
func UpdateUserEquipment(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	id, ok := utils.ParseUUID(c, params.ID, "user equipment")
	if !ok {
		return
	}

	// Find user equipment
	var userEquipment models.UserEquipment
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&userEquipment).Error; err != nil {
		utils.NotFoundResponse(c, "User equipment not found.")
		return
	}

	var req models.UpdateUserEquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Update fields if provided
	if req.LocationType != "" {
		// Check if changing location would create a duplicate
		var existingUE models.UserEquipment
		err := database.DB.Where("user_id = ? AND equipment_id = ? AND location_type = ? AND id != ?",
			userID, userEquipment.EquipmentID, req.LocationType, id).First(&existingUE).Error

		if err == nil {
			utils.ConflictResponse(c, "You already have this equipment at this location.")
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
		utils.BadRequestResponse(c, "Validation failed.", err.Error())
		return
	}

	if err := database.DB.Save(&userEquipment).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update equipment.")
		return
	}

	// Load equipment for response
	database.DB.Preload("Equipment").First(&userEquipment, userEquipment.ID)

	utils.SuccessResponse(c, "Equipment updated successfully.", userEquipment.ToResponse())
}

// RemoveUserEquipment removes equipment from user's inventory
func RemoveUserEquipment(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	id, ok := utils.ParseUUID(c, params.ID, "user equipment")
	if !ok {
		return
	}

	// Delete user equipment
	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.UserEquipment{})
	if result.Error != nil {
		utils.InternalServerErrorResponse(c, "Failed to remove equipment.")
		return
	}

	if result.RowsAffected == 0 {
		utils.NotFoundResponse(c, "User equipment not found.")
		return
	}

	utils.DeletedResponse(c, "Equipment removed successfully.")
}

// GetUserEquipmentByLocation gets equipment filtered by location
func GetUserEquipmentByLocation(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var params struct {
		Location string `uri:"location" binding:"required,oneof=home gym"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
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

	// Get total count
	var total int64
	if err := database.DB.Model(&models.UserEquipment{}).
		Where("user_id = ? AND location_type = ?", userID, params.Location).
		Count(&total).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to count equipment.")
		return
	}

	var userEquipment []models.UserEquipment
	if err := database.DB.Where("user_id = ? AND location_type = ?", userID, params.Location).
		Preload("Equipment").
		Offset(offset).
		Limit(queryParams.Limit).
		Order("created_at DESC").
		Find(&userEquipment).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve equipment by location.")
		return
	}

	// Convert to response format
	response := make([]models.UserEquipmentResponse, len(userEquipment))
	for i, ue := range userEquipment {
		response[i] = ue.ToResponse()
	}

	utils.PaginatedResponse(c, "Equipment by location retrieved successfully.", response, queryParams.Page, queryParams.Limit, int(total))
}

// BulkAddUserEquipment allows adding multiple equipment items at once
func BulkAddUserEquipment(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var req struct {
		Equipment []models.AddUserEquipmentRequest `json:"equipment" binding:"required,dive"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
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
		utils.BadRequestResponse(c, "Some equipment items were not found.", nil)
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
			UserID:       userID,
			EquipmentID:  item.EquipmentID,
			LocationType: item.LocationType,
			GymLocation:  item.GymLocation,
			Notes:        item.Notes,
		}

		if err := tx.Create(&userEquipment).Error; err != nil {
			tx.Rollback()
			utils.InternalServerErrorResponse(c, "Failed to add equipment items.")
			return
		}
		created = append(created, userEquipment)
	}

	if err := tx.Commit().Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to add equipment items.")
		return
	}

	// Load equipment for response
	database.DB.Where("id IN ?", created).Preload("Equipment").Find(&created)

	response := make([]models.UserEquipmentResponse, len(created))
	for i, ue := range created {
		response[i] = ue.ToResponse()
	}

	utils.CreatedResponse(c, "Equipment items added successfully.", gin.H{
		"created": response,
		"total":   len(response),
	})
}
