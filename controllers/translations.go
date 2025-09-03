package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateTranslation creates a new translation
func CreateTranslation(c *gin.Context) {
	var req models.CreateTranslationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	translation := models.Translation{
		ResourceType: req.ResourceType,
		ResourceID:   req.ResourceID,
		FieldName:    req.FieldName,
		Language:     req.Language,
		Content:      req.Content,
	}

	if err := translation.Validate(); err != nil {
		utils.BadRequestResponse(c, "Validation failed.", err.Error())
		return
	}

	if err := database.DB.Create(&translation).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create translation.")
		return
	}

	utils.CreatedResponse(c, "Translation created successfully.", translation.ToResponse())
}

// GetTranslations retrieves translations for a specific resource
func GetTranslations(c *gin.Context) {
	resourceType := c.Query("resource_type")
	resourceIDStr := c.Query("resource_id")
	language := c.Query("language")

	// Pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}
	
	offset := (page - 1) * limit

	query := database.DB.Model(&models.Translation{})

	if resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}

	if resourceIDStr != "" {
		resourceID, err := uuid.Parse(resourceIDStr)
		if err != nil {
			utils.BadRequestResponse(c, "Invalid resource ID format.", nil)
			return
		}
		query = query.Where("resource_id = ?", resourceID)
	}

	if language != "" {
		query = query.Where("language = ?", language)
	}

	// Get total count
	var total int64
	query.Count(&total)

	var translations []models.Translation
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&translations).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve translations.")
		return
	}

	var responses []models.TranslationResponse
	for _, translation := range translations {
		responses = append(responses, translation.ToResponse())
	}

	utils.PaginatedResponse(c, "Translations retrieved successfully.", responses, page, limit, int(total))
}

// GetTranslation retrieves a specific translation
func GetTranslation(c *gin.Context) {
	id := c.Param("id")
	translationID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID format.", nil)
		return
	}

	var translation models.Translation
	if err := database.DB.First(&translation, translationID).Error; err != nil {
		utils.NotFoundResponse(c, "Translation not found.")
		return
	}

	utils.SuccessResponse(c, "Translation retrieved successfully.", translation.ToResponse())
}

// UpdateTranslation updates an existing translation
func UpdateTranslation(c *gin.Context) {
	id := c.Param("id")
	translationID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID format.", nil)
		return
	}

	var req models.UpdateTranslationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	var translation models.Translation
	if err := database.DB.First(&translation, translationID).Error; err != nil {
		utils.NotFoundResponse(c, "Translation not found.")
		return
	}

	translation.Content = req.Content

	if err := translation.Validate(); err != nil {
		utils.BadRequestResponse(c, "Validation failed.", err.Error())
		return
	}

	if err := database.DB.Save(&translation).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update translation.")
		return
	}

	utils.SuccessResponse(c, "Translation updated successfully.", translation.ToResponse())
}

// DeleteTranslation deletes a translation
func DeleteTranslation(c *gin.Context) {
	id := c.Param("id")
	translationID, err := uuid.Parse(id)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid ID format.", nil)
		return
	}

	var translation models.Translation
	if err := database.DB.First(&translation, translationID).Error; err != nil {
		utils.NotFoundResponse(c, "Translation not found.")
		return
	}

	if err := database.DB.Delete(&translation).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete translation.")
		return
	}

	utils.DeletedResponse(c, "Translation deleted successfully.")
}

// GetResourceTranslations retrieves all translations for a specific resource
func GetResourceTranslations(c *gin.Context) {
	resourceType := c.Param("resource_type")
	resourceIDStr := c.Param("resource_id")

	resourceID, err := uuid.Parse(resourceIDStr)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid resource ID format.", nil)
		return
	}

	var translations []models.Translation
	if err := database.DB.Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).Find(&translations).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve resource translations.")
		return
	}

	// Group translations by field and language
	result := make(map[string]map[string]string)
	for _, translation := range translations {
		if result[translation.FieldName] == nil {
			result[translation.FieldName] = make(map[string]string)
		}
		result[translation.FieldName][translation.Language] = translation.Content
	}

	utils.SuccessResponse(c, "Resource translations retrieved successfully.", result)
}

// CreateOrUpdateTranslation creates or updates a translation
func CreateOrUpdateTranslation(c *gin.Context) {
	var req models.CreateTranslationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	var translation models.Translation
	err := database.DB.Where("resource_type = ? AND resource_id = ? AND field_name = ? AND language = ?",
		req.ResourceType, req.ResourceID, req.FieldName, req.Language).First(&translation).Error

	if err != nil {
		// Create new translation
		translation = models.Translation{
			ResourceType: req.ResourceType,
			ResourceID:   req.ResourceID,
			FieldName:    req.FieldName,
			Language:     req.Language,
			Content:      req.Content,
		}

		if err := translation.Validate(); err != nil {
			utils.BadRequestResponse(c, "Validation failed.", err.Error())
			return
		}

		if err := database.DB.Create(&translation).Error; err != nil {
			utils.InternalServerErrorResponse(c, "Failed to create translation.")
			return
		}

		utils.CreatedResponse(c, "Translation created successfully.", translation.ToResponse())
	} else {
		// Update existing translation
		translation.Content = req.Content

		if err := translation.Validate(); err != nil {
			utils.BadRequestResponse(c, "Validation failed.", err.Error())
			return
		}

		if err := database.DB.Save(&translation).Error; err != nil {
			utils.InternalServerErrorResponse(c, "Failed to update translation.")
			return
		}

		utils.SuccessResponse(c, "Translation updated successfully.", translation.ToResponse())
	}
}