package controllers

import (
	"lamari-fit-api/database"
	"lamari-fit-api/models"
	"lamari-fit-api/utils"

	"github.com/gin-gonic/gin"
)

// ListSpecialties retrieves all available trainer specialties
func ListSpecialties(c *gin.Context) {
	var specialties []models.Specialty

	if err := database.DB.Order("name ASC").Find(&specialties).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve specialties")
		return
	}

	// Convert to response format
	responses := make([]models.SpecialtyResponse, len(specialties))
	for i, specialty := range specialties {
		responses[i] = specialty.ToResponse()
	}

	utils.SuccessResponse(c, "Specialties retrieved successfully", responses)
}
