package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/middleware"
	"fit-flow-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateTranslation creates a new translation
func CreateTranslation(c *gin.Context) {
	var req models.CreateTranslationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
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
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
		return
	}

	if err := database.DB.Create(&translation).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "general.internal_error", nil)
		return
	}

	middleware.TranslateResponse(c, http.StatusCreated, "general.success", translation.ToResponse())
}

// GetTranslations retrieves translations for a specific resource
func GetTranslations(c *gin.Context) {
	resourceType := c.Query("resource_type")
	resourceIDStr := c.Query("resource_id")
	language := c.Query("language")

	var translations []models.Translation
	query := database.DB.Model(&models.Translation{})

	if resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}

	if resourceIDStr != "" {
		resourceID, err := uuid.Parse(resourceIDStr)
		if err != nil {
			middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_uuid", nil)
			return
		}
		query = query.Where("resource_id = ?", resourceID)
	}

	if language != "" {
		query = query.Where("language = ?", language)
	}

	if err := query.Find(&translations).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "general.internal_error", nil)
		return
	}

	var responses []models.TranslationResponse
	for _, translation := range translations {
		responses = append(responses, translation.ToResponse())
	}

	middleware.TranslateResponse(c, http.StatusOK, "general.success", responses)
}

// GetTranslation retrieves a specific translation
func GetTranslation(c *gin.Context) {
	id := c.Param("id")
	translationID, err := uuid.Parse(id)
	if err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_uuid", nil)
		return
	}

	var translation models.Translation
	if err := database.DB.First(&translation, translationID).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusNotFound, "general.not_found", nil)
		return
	}

	middleware.TranslateResponse(c, http.StatusOK, "general.success", translation.ToResponse())
}

// UpdateTranslation updates an existing translation
func UpdateTranslation(c *gin.Context) {
	id := c.Param("id")
	translationID, err := uuid.Parse(id)
	if err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_uuid", nil)
		return
	}

	var req models.UpdateTranslationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
		return
	}

	var translation models.Translation
	if err := database.DB.First(&translation, translationID).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusNotFound, "general.not_found", nil)
		return
	}

	translation.Content = req.Content

	if err := translation.Validate(); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
		return
	}

	if err := database.DB.Save(&translation).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "general.internal_error", nil)
		return
	}

	middleware.TranslateResponse(c, http.StatusOK, "general.success", translation.ToResponse())
}

// DeleteTranslation deletes a translation
func DeleteTranslation(c *gin.Context) {
	id := c.Param("id")
	translationID, err := uuid.Parse(id)
	if err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_uuid", nil)
		return
	}

	var translation models.Translation
	if err := database.DB.First(&translation, translationID).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusNotFound, "general.not_found", nil)
		return
	}

	if err := database.DB.Delete(&translation).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "general.internal_error", nil)
		return
	}

	middleware.TranslateResponse(c, http.StatusOK, "general.success", nil)
}

// GetResourceTranslations retrieves all translations for a specific resource
func GetResourceTranslations(c *gin.Context) {
	resourceType := c.Param("resource_type")
	resourceIDStr := c.Param("resource_id")

	resourceID, err := uuid.Parse(resourceIDStr)
	if err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_uuid", nil)
		return
	}

	var translations []models.Translation
	if err := database.DB.Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).Find(&translations).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "general.internal_error", nil)
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

	middleware.TranslateResponse(c, http.StatusOK, "general.success", result)
}

// CreateOrUpdateTranslation creates or updates a translation
func CreateOrUpdateTranslation(c *gin.Context) {
	var req models.CreateTranslationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
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
			middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
			return
		}

		if err := database.DB.Create(&translation).Error; err != nil {
			middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "general.internal_error", nil)
			return
		}

		middleware.TranslateResponse(c, http.StatusCreated, "general.success", translation.ToResponse())
	} else {
		// Update existing translation
		translation.Content = req.Content

		if err := translation.Validate(); err != nil {
			middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
			return
		}

		if err := database.DB.Save(&translation).Error; err != nil {
			middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "general.internal_error", nil)
			return
		}

		middleware.TranslateResponse(c, http.StatusOK, "general.success", translation.ToResponse())
	}
}