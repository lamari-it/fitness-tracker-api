package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListRPEScales retrieves all RPE scales (global + user's custom scales)
func ListRPEScales(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	// Get trainer IDs where user is an active client
	var trainerIDs []uuid.UUID
	database.DB.Model(&models.TrainerClientLink{}).
		Where("client_id = ? AND status = ?", userID, "active").
		Pluck("trainer_id", &trainerIDs)

	var scales []models.RPEScale
	// Get global scales, user's custom scales, and scales from active trainers
	query := database.DB.Preload("Values").
		Where("is_global = ? OR trainer_id = ?", true, userID)

	if len(trainerIDs) > 0 {
		query = query.Or("trainer_id IN ?", trainerIDs)
	}

	if err := query.Order("is_global DESC, name ASC").
		Find(&scales).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve RPE scales")
		return
	}

	responses := make([]models.RPEScaleResponse, len(scales))
	for i, scale := range scales {
		responses[i] = scale.ToResponse()
	}

	utils.SuccessResponse(c, "RPE scales retrieved successfully", responses)
}

// GetRPEScale retrieves a single RPE scale by ID
func GetRPEScale(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	scaleID, ok := utils.ParseUUID(c, params.ID, "RPE scale")
	if !ok {
		return
	}

	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var scale models.RPEScale
	if err := database.DB.Preload("Values").First(&scale, "id = ?", scaleID).Error; err != nil {
		utils.NotFoundResponse(c, "RPE scale not found")
		return
	}

	// Check access: global scales are accessible to all, custom scales to owner or active clients
	if !scale.IsGlobal {
		isOwner := scale.TrainerID != nil && *scale.TrainerID == userID
		if !isOwner {
			// Check if user is an active client of the trainer who owns this scale
			var link models.TrainerClientLink
			if scale.TrainerID != nil {
				err := database.DB.Where(
					"trainer_id = ? AND client_id = ? AND status = ?",
					*scale.TrainerID, userID, "active",
				).First(&link).Error
				if err != nil {
					utils.NotFoundResponse(c, "RPE scale not found")
					return
				}
			} else {
				utils.NotFoundResponse(c, "RPE scale not found")
				return
			}
		}
	}

	utils.SuccessResponse(c, "RPE scale retrieved successfully", scale.ToResponse())
}

// CreateRPEScale creates a new custom RPE scale for the trainer
func CreateRPEScale(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var req models.CreateRPEScaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Set defaults if not provided
	minValue := req.MinValue
	maxValue := req.MaxValue
	if maxValue == 0 {
		maxValue = 10
	}
	if minValue == 0 {
		minValue = 1
	}

	scale := models.RPEScale{
		Name:        req.Name,
		Description: req.Description,
		MinValue:    minValue,
		MaxValue:    maxValue,
		IsGlobal:    false,
		TrainerID:   &userID,
	}

	if err := scale.Validate(); err != nil {
		utils.BadRequestResponse(c, "Validation failed", err.Error())
		return
	}

	// Start transaction
	tx := database.DB.Begin()

	if err := tx.Create(&scale).Error; err != nil {
		tx.Rollback()
		utils.InternalServerErrorResponse(c, "Failed to create RPE scale")
		return
	}

	// Create scale values if provided
	if len(req.Values) > 0 {
		for _, valueReq := range req.Values {
			value := models.RPEScaleValue{
				ScaleID:     scale.ID,
				Value:       valueReq.Value,
				Label:       valueReq.Label,
				Description: valueReq.Description,
			}

			if err := value.Validate(&scale); err != nil {
				tx.Rollback()
				utils.BadRequestResponse(c, "Validation failed", err.Error())
				return
			}

			if err := tx.Create(&value).Error; err != nil {
				tx.Rollback()
				utils.InternalServerErrorResponse(c, "Failed to create RPE scale value")
				return
			}
		}
	}

	tx.Commit()

	// Reload with values
	database.DB.Preload("Values").First(&scale, "id = ?", scale.ID)

	utils.CreatedResponse(c, "RPE scale created successfully", scale.ToResponse())
}

// UpdateRPEScale updates an existing custom RPE scale
func UpdateRPEScale(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	scaleID, ok := utils.ParseUUID(c, params.ID, "RPE scale")
	if !ok {
		return
	}

	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var scale models.RPEScale
	if err := database.DB.First(&scale, "id = ?", scaleID).Error; err != nil {
		utils.NotFoundResponse(c, "RPE scale not found")
		return
	}

	// Cannot update global scales
	if scale.IsGlobal {
		utils.ForbiddenResponse(c, "Cannot modify global RPE scales")
		return
	}

	// Check ownership
	if scale.TrainerID == nil || *scale.TrainerID != userID {
		utils.ForbiddenResponse(c, "You can only update your own RPE scales")
		return
	}

	var req models.UpdateRPEScaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Update fields
	if req.Name != "" {
		scale.Name = req.Name
	}
	if req.Description != "" {
		scale.Description = req.Description
	}

	if err := scale.Validate(); err != nil {
		utils.BadRequestResponse(c, "Validation failed", err.Error())
		return
	}

	if err := database.DB.Save(&scale).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update RPE scale")
		return
	}

	// Reload with values
	database.DB.Preload("Values").First(&scale, "id = ?", scale.ID)

	utils.SuccessResponse(c, "RPE scale updated successfully", scale.ToResponse())
}

// DeleteRPEScale deletes a custom RPE scale
func DeleteRPEScale(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	scaleID, ok := utils.ParseUUID(c, params.ID, "RPE scale")
	if !ok {
		return
	}

	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var scale models.RPEScale
	if err := database.DB.First(&scale, "id = ?", scaleID).Error; err != nil {
		utils.NotFoundResponse(c, "RPE scale not found")
		return
	}

	// Cannot delete global scales
	if scale.IsGlobal {
		utils.ForbiddenResponse(c, "Cannot delete global RPE scales")
		return
	}

	// Check ownership
	if scale.TrainerID == nil || *scale.TrainerID != userID {
		utils.ForbiddenResponse(c, "You can only delete your own RPE scales")
		return
	}

	if err := database.DB.Delete(&scale).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete RPE scale")
		return
	}

	utils.NoContentResponse(c)
}

// AddRPEScaleValue adds a new value to an existing custom RPE scale
func AddRPEScaleValue(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	scaleID, ok := utils.ParseUUID(c, params.ID, "RPE scale")
	if !ok {
		return
	}

	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var scale models.RPEScale
	if err := database.DB.First(&scale, "id = ?", scaleID).Error; err != nil {
		utils.NotFoundResponse(c, "RPE scale not found")
		return
	}

	// Cannot modify global scales
	if scale.IsGlobal {
		utils.ForbiddenResponse(c, "Cannot modify global RPE scales")
		return
	}

	// Check ownership
	if scale.TrainerID == nil || *scale.TrainerID != userID {
		utils.ForbiddenResponse(c, "You can only modify your own RPE scales")
		return
	}

	var req models.AddRPEScaleValueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	value := models.RPEScaleValue{
		ScaleID:     scaleID,
		Value:       req.Value,
		Label:       req.Label,
		Description: req.Description,
	}

	if err := value.Validate(&scale); err != nil {
		utils.BadRequestResponse(c, "Validation failed", err.Error())
		return
	}

	if err := database.DB.Create(&value).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to add RPE scale value")
		return
	}

	utils.CreatedResponse(c, "RPE scale value added successfully", value.ToResponse())
}

// GetGlobalRPEScale retrieves the global RPE scale (convenience endpoint)
func GetGlobalRPEScale(c *gin.Context) {
	var scale models.RPEScale
	if err := database.DB.Preload("Values").First(&scale, "is_global = ?", true).Error; err != nil {
		utils.NotFoundResponse(c, "Global RPE scale not found")
		return
	}

	utils.SuccessResponse(c, "Global RPE scale retrieved successfully", scale.ToResponse())
}
