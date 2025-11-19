package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateTrainerProfile creates a new trainer profile for the authenticated user
func CreateTrainerProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type")
		return
	}

	// Check if user already has a trainer profile
	var existingProfile models.TrainerProfile
	if err := database.DB.Where("user_id = ?", userID).First(&existingProfile).Error; err == nil {
		utils.ConflictResponse(c, "Trainer profile already exists for this user")
		return
	}

	var req models.CreateTrainerProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Validate that all specialty IDs exist
	var specialties []models.Specialty
	if err := database.DB.Where("id IN ?", req.SpecialtyIDs).Find(&specialties).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to validate specialties")
		return
	}
	if len(specialties) != len(req.SpecialtyIDs) {
		utils.BadRequestResponse(c, "One or more specialty IDs are invalid", nil)
		return
	}

	// Set default visibility if not provided
	visibility := req.Visibility
	if visibility == "" {
		visibility = "public"
	}

	trainerProfile := models.TrainerProfile{
		UserID:      userID,
		Bio:         req.Bio,
		Specialties: specialties,
		HourlyRate:  req.HourlyRate,
		Location:    req.Location,
		Visibility:  visibility,
	}

	if err := trainerProfile.Validate(); err != nil {
		utils.BadRequestResponse(c, "Validation failed", err.Error())
		return
	}

	if err := database.DB.Create(&trainerProfile).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create trainer profile")
		return
	}

	// Associate specialties with the trainer profile
	if err := database.DB.Model(&trainerProfile).Association("Specialties").Replace(&specialties); err != nil {
		utils.InternalServerErrorResponse(c, "Failed to associate specialties")
		return
	}

	// Preload user and specialties for response
	database.DB.Preload("User").Preload("Specialties").First(&trainerProfile, "id = ?", trainerProfile.ID)

	utils.CreatedResponse(c, "Trainer profile created successfully", trainerProfile.ToResponse())
}

// GetTrainerProfile retrieves the authenticated user's trainer profile
func GetTrainerProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type")
		return
	}

	var trainerProfile models.TrainerProfile
	if err := database.DB.Preload("User").Preload("Specialties").Where("user_id = ?", userID).First(&trainerProfile).Error; err != nil {
		utils.NotFoundResponse(c, "Trainer profile not found")
		return
	}

	utils.SuccessResponse(c, "Trainer profile retrieved successfully", trainerProfile.ToResponse())
}

// UpdateTrainerProfile updates the authenticated user's trainer profile
func UpdateTrainerProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type")
		return
	}

	// Check if user is admin
	var currentUser models.User
	if err := database.DB.First(&currentUser, "id = ?", userID).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve user")
		return
	}

	var trainerProfile models.TrainerProfile
	if err := database.DB.Where("user_id = ?", userID).First(&trainerProfile).Error; err != nil {
		// If not found and user is admin, they might be trying to update someone else's profile
		if !currentUser.IsAdmin {
			utils.NotFoundResponse(c, "Trainer profile not found")
			return
		}
		utils.NotFoundResponse(c, "Trainer profile not found")
		return
	}

	var req models.UpdateTrainerProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Update only provided fields
	if req.Bio != "" {
		trainerProfile.Bio = req.Bio
	}
	if len(req.SpecialtyIDs) > 0 {
		// Validate that all specialty IDs exist
		var specialties []models.Specialty
		if err := database.DB.Where("id IN ?", req.SpecialtyIDs).Find(&specialties).Error; err != nil {
			utils.InternalServerErrorResponse(c, "Failed to validate specialties")
			return
		}
		if len(specialties) != len(req.SpecialtyIDs) {
			utils.BadRequestResponse(c, "One or more specialty IDs are invalid", nil)
			return
		}

		// Replace the specialties association
		if err := database.DB.Model(&trainerProfile).Association("Specialties").Replace(&specialties); err != nil {
			utils.InternalServerErrorResponse(c, "Failed to update specialties")
			return
		}

		trainerProfile.Specialties = specialties
	}
	if req.HourlyRate > 0 {
		trainerProfile.HourlyRate = req.HourlyRate
	}
	if req.Location != "" {
		trainerProfile.Location = req.Location
	}
	if req.Visibility != "" {
		trainerProfile.Visibility = req.Visibility
	}

	if err := trainerProfile.Validate(); err != nil {
		utils.BadRequestResponse(c, "Validation failed", err.Error())
		return
	}

	if err := database.DB.Save(&trainerProfile).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update trainer profile")
		return
	}

	// Preload user and specialties for response
	database.DB.Preload("User").Preload("Specialties").First(&trainerProfile, "id = ?", trainerProfile.ID)

	utils.SuccessResponse(c, "Trainer profile updated successfully", trainerProfile.ToResponse())
}

// DeleteTrainerProfile deletes the authenticated user's trainer profile
func DeleteTrainerProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type")
		return
	}

	var trainerProfile models.TrainerProfile
	if err := database.DB.Where("user_id = ?", userID).First(&trainerProfile).Error; err != nil {
		utils.NotFoundResponse(c, "Trainer profile not found")
		return
	}

	// Check for active trainer-client relationships
	var activeCount int64
	database.DB.Model(&models.TrainerClientLink{}).Where("trainer_id = ? AND status = ?", trainerProfile.ID, "active").Count(&activeCount)
	if activeCount > 0 {
		utils.ConflictResponse(c, "Cannot delete trainer profile with active client relationships")
		return
	}

	if err := database.DB.Delete(&trainerProfile).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete trainer profile")
		return
	}

	utils.NoContentResponse(c)
}

// GetTrainerPublicProfile retrieves a trainer's public profile by ID
func GetTrainerPublicProfile(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	trainerID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid UUID format", nil)
		return
	}

	// Get current user info for access control
	userIDVal, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type")
		return
	}

	var currentUser models.User
	if err := database.DB.First(&currentUser, "id = ?", userID).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve user")
		return
	}

	var trainerProfile models.TrainerProfile
	if err := database.DB.Preload("User").Preload("Specialties").First(&trainerProfile, "id = ?", trainerID).Error; err != nil {
		utils.NotFoundResponse(c, "Trainer not found")
		return
	}

	// Check visibility access
	// - public: anyone can view
	// - link_only: anyone with the ID can view
	// - private: only owner or admin can view
	if trainerProfile.Visibility == "private" {
		if trainerProfile.UserID != userID && !currentUser.IsAdmin {
			utils.NotFoundResponse(c, "Trainer not found")
			return
		}
	}

	// Calculate review statistics
	var reviewCount int64
	var avgRating float64

	database.DB.Model(&models.TrainerReview{}).Where("trainer_id = ?", trainerID).Count(&reviewCount)
	if reviewCount > 0 {
		database.DB.Model(&models.TrainerReview{}).Where("trainer_id = ?", trainerID).Select("COALESCE(AVG(rating), 0)").Scan(&avgRating)
	}

	utils.SuccessResponse(c, "Trainer retrieved successfully", trainerProfile.ToPublicResponse(int(reviewCount), avgRating))
}

// ListTrainers retrieves a paginated list of trainers with optional filtering
func ListTrainers(c *gin.Context) {
	var queryParams TrainerQuery
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	SetDefaultPagination(&queryParams.PaginationQuery)

	query := database.DB.Model(&models.TrainerProfile{}).Preload("User").Preload("Specialties")

	// Only show public trainers in list
	query = query.Where("visibility = ?", "public")

	// Search by user first_name, last_name, or location
	if queryParams.Search != "" {
		searchPattern := "%" + queryParams.Search + "%"
		query = query.Joins("JOIN users ON users.id = trainer_profiles.user_id").
			Where("users.first_name ILIKE ? OR users.last_name ILIKE ? OR trainer_profiles.location ILIKE ?",
				searchPattern, searchPattern, searchPattern)
	}

	// Filter by specialty (join with trainer_specialties and specialties)
	if queryParams.Specialty != "" {
		query = query.Joins("JOIN trainer_specialties ON trainer_specialties.trainer_profile_id = trainer_profiles.id").
			Joins("JOIN specialties ON specialties.id = trainer_specialties.specialty_id").
			Where("specialties.name ILIKE ?", "%"+queryParams.Specialty+"%")
	}

	// Filter by location
	if queryParams.Location != "" {
		query = query.Where("location ILIKE ?", "%"+queryParams.Location+"%")
	}

	// Get total count before pagination
	var total int64
	countQuery := query.Session(&gorm.Session{})
	countQuery.Count(&total)

	// Apply sorting
	switch queryParams.SortBy {
	case "rate":
		query = query.Order("hourly_rate ASC")
	case "recent":
		query = query.Order("created_at DESC")
	default:
		query = query.Order("created_at DESC")
	}

	offset := (queryParams.Page - 1) * queryParams.Limit

	var trainers []models.TrainerProfile
	if err := query.Offset(offset).Limit(queryParams.Limit).Find(&trainers).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve trainers")
		return
	}

	// Build public responses with review stats
	responses := make([]models.TrainerPublicResponse, len(trainers))
	for i, trainer := range trainers {
		var reviewCount int64
		var avgRating float64

		database.DB.Model(&models.TrainerReview{}).Where("trainer_id = ?", trainer.ID).Count(&reviewCount)
		if reviewCount > 0 {
			database.DB.Model(&models.TrainerReview{}).Where("trainer_id = ?", trainer.ID).Select("COALESCE(AVG(rating), 0)").Scan(&avgRating)
		}

		responses[i] = trainer.ToPublicResponse(int(reviewCount), avgRating)
	}

	// Filter by min_rating if specified (post-query filter since we need to calculate ratings)
	if queryParams.MinRating > 0 {
		filteredResponses := make([]models.TrainerPublicResponse, 0)
		for _, resp := range responses {
			if resp.AverageRating >= queryParams.MinRating {
				filteredResponses = append(filteredResponses, resp)
			}
		}
		responses = filteredResponses
		total = int64(len(responses))
	}

	utils.PaginatedResponse(c, "Trainers retrieved successfully", responses, queryParams.Page, queryParams.Limit, int(total))
}
