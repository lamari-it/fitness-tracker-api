package controllers

import (
	"lamari-fit-api/database"
	"lamari-fit-api/models"
	"lamari-fit-api/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserSearchResponse includes distance when radius search is used
type UserSearchResponse struct {
	ID                  uuid.UUID                `json:"id"`
	FirstName           string                   `json:"first_name"`
	LastName            string                   `json:"last_name"`
	Bio                 string                   `json:"bio,omitempty"`
	IsLookingForTrainer bool                     `json:"is_looking_for_trainer"`
	Location            *models.LocationResponse `json:"location,omitempty"`
	DistanceKm          *float64                 `json:"distance_km,omitempty"`
	CreatedAt           time.Time                `json:"created_at"`
}

// TrainerSearchResponse includes distance when radius search is used
type TrainerSearchResponse struct {
	ID                  uuid.UUID                  `json:"id"`
	UserID              uuid.UUID                  `json:"user_id"`
	Bio                 string                     `json:"bio"`
	Specialties         []models.SpecialtyResponse `json:"specialties"`
	HourlyRate          *float64                   `json:"hourly_rate,omitempty"`
	Location            *models.LocationResponse   `json:"location,omitempty"`
	Visibility          string                     `json:"visibility"`
	IsLookingForClients bool                       `json:"is_looking_for_clients"`
	User                *models.UserPublicResponse `json:"user"`
	ReviewCount         int                        `json:"review_count"`
	AverageRating       float64                    `json:"average_rating"`
	DistanceKm          *float64                   `json:"distance_km,omitempty"`
	CreatedAt           time.Time                  `json:"created_at"`
}

// SearchUsers searches for public users with location filtering
// GET /api/v1/search/users
func SearchUsers(c *gin.Context) {
	var queryParams UserSearchQuery
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	SetDefaultPagination(&queryParams.PaginationQuery)

	query := database.DB.Model(&models.User{}).Where("profile_visibility = ?", "public")

	// Apply is_looking_for_trainer filter
	if queryParams.IsLookingForTrainer == "true" {
		query = query.Where("is_looking_for_trainer = ?", true)
	} else if queryParams.IsLookingForTrainer == "false" {
		query = query.Where("is_looking_for_trainer = ?", false)
	}

	// Apply location-based filtering
	strategy := queryParams.SearchStrategy()
	var centerLat, centerLng float64

	switch strategy {
	case "radius":
		centerLat = *queryParams.Latitude
		centerLng = *queryParams.Longitude
		radiusKm := utils.DefaultRadiusKm
		if queryParams.RadiusKm != nil {
			radiusKm = *queryParams.RadiusKm
		}

		// Apply bounding box pre-filter for performance
		bbox := utils.CalculateBoundingBox(centerLat, centerLng, radiusKm)
		query = query.Where("location_latitude IS NOT NULL AND location_longitude IS NOT NULL")
		query = query.Where("location_latitude BETWEEN ? AND ?", bbox.MinLat, bbox.MaxLat)
		query = query.Where("location_longitude BETWEEN ? AND ?", bbox.MinLng, bbox.MaxLng)

		// Apply Haversine distance filter
		distanceExpr := utils.HaversineSQL("location_latitude", "location_longitude", centerLat, centerLng)
		query = query.Where(distanceExpr+" <= ?", radiusKm)

	case "structured":
		if queryParams.CountryCode != "" {
			query = query.Where("location_country_code = ?", queryParams.CountryCode)
		}
		if queryParams.Region != "" {
			query = query.Where("location_region ILIKE ?", "%"+queryParams.Region+"%")
		}
		if queryParams.City != "" {
			query = query.Where("location_city ILIKE ?", "%"+queryParams.City+"%")
		}
		if queryParams.District != "" {
			query = query.Where("location_district ILIKE ?", "%"+queryParams.District+"%")
		}
		if queryParams.PostalCode != "" {
			query = query.Where("location_postal_code = ?", queryParams.PostalCode)
		}

	case "free_text":
		searchPattern := "%" + queryParams.Q + "%"
		query = query.Where(`
			location_city ILIKE ? OR
			location_region ILIKE ? OR
			location_district ILIKE ? OR
			location_raw_address ILIKE ? OR
			first_name ILIKE ? OR
			last_name ILIKE ?`,
			searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Get total count before pagination
	var total int64
	countQuery := query.Session(&gorm.Session{})
	countQuery.Count(&total)

	// Apply ordering
	if strategy == "radius" {
		distanceExpr := utils.HaversineSQL("location_latitude", "location_longitude", centerLat, centerLng)
		query = query.Order(distanceExpr + " ASC")
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	offset := queryParams.GetOffset()

	var users []models.User
	if err := query.Offset(offset).Limit(queryParams.Limit).Find(&users).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve users")
		return
	}

	// Build responses
	responses := make([]UserSearchResponse, len(users))
	for i, user := range users {
		resp := UserSearchResponse{
			ID:                  user.ID,
			FirstName:           user.FirstName,
			LastName:            user.LastName,
			Bio:                 user.Bio,
			IsLookingForTrainer: user.IsLookingForTrainer,
			Location:            user.Location.ToResponse(),
			CreatedAt:           user.CreatedAt,
		}

		// Calculate distance if radius search and user has coordinates
		if strategy == "radius" && user.Location.HasCoordinates() {
			distance := utils.HaversineDistance(centerLat, centerLng, *user.Location.Latitude, *user.Location.Longitude)
			resp.DistanceKm = &distance
		}

		responses[i] = resp
	}

	utils.PaginatedResponse(c, "Users retrieved successfully", responses, queryParams.Page, queryParams.Limit, int(total))
}

// SearchTrainers searches for public trainers with location filtering
// GET /api/v1/search/trainers
func SearchTrainers(c *gin.Context) {
	var queryParams TrainerSearchQuery
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	SetDefaultPagination(&queryParams.PaginationQuery)

	query := database.DB.Model(&models.TrainerProfile{}).
		Preload("User").
		Preload("Specialties").
		Where("visibility = ?", "public")

	// Apply is_looking_for_clients filter
	if queryParams.IsLookingForClients == "true" {
		query = query.Where("is_looking_for_clients = ?", true)
	} else if queryParams.IsLookingForClients == "false" {
		query = query.Where("is_looking_for_clients = ?", false)
	}

	// Filter by specialty
	if queryParams.Specialty != "" {
		query = query.Joins("JOIN trainer_specialties ON trainer_specialties.trainer_profile_id = trainer_profiles.id").
			Joins("JOIN specialties ON specialties.id = trainer_specialties.specialty_id").
			Where("specialties.name ILIKE ?", "%"+queryParams.Specialty+"%")
	}

	// Apply location-based filtering
	strategy := queryParams.SearchStrategy()
	var centerLat, centerLng float64

	switch strategy {
	case "radius":
		centerLat = *queryParams.Latitude
		centerLng = *queryParams.Longitude
		radiusKm := utils.DefaultRadiusKm
		if queryParams.RadiusKm != nil {
			radiusKm = *queryParams.RadiusKm
		}

		// Apply bounding box pre-filter for performance
		bbox := utils.CalculateBoundingBox(centerLat, centerLng, radiusKm)
		query = query.Where("location_latitude IS NOT NULL AND location_longitude IS NOT NULL")
		query = query.Where("location_latitude BETWEEN ? AND ?", bbox.MinLat, bbox.MaxLat)
		query = query.Where("location_longitude BETWEEN ? AND ?", bbox.MinLng, bbox.MaxLng)

		// Apply Haversine distance filter
		distanceExpr := utils.HaversineSQL("location_latitude", "location_longitude", centerLat, centerLng)
		query = query.Where(distanceExpr+" <= ?", radiusKm)

	case "structured":
		if queryParams.CountryCode != "" {
			query = query.Where("location_country_code = ?", queryParams.CountryCode)
		}
		if queryParams.Region != "" {
			query = query.Where("location_region ILIKE ?", "%"+queryParams.Region+"%")
		}
		if queryParams.City != "" {
			query = query.Where("location_city ILIKE ?", "%"+queryParams.City+"%")
		}
		if queryParams.District != "" {
			query = query.Where("location_district ILIKE ?", "%"+queryParams.District+"%")
		}
		if queryParams.PostalCode != "" {
			query = query.Where("location_postal_code = ?", queryParams.PostalCode)
		}

	case "free_text":
		searchPattern := "%" + queryParams.Q + "%"
		query = query.Joins("JOIN users ON users.id = trainer_profiles.user_id").
			Where(`
				trainer_profiles.location_city ILIKE ? OR
				trainer_profiles.location_region ILIKE ? OR
				trainer_profiles.location_district ILIKE ? OR
				trainer_profiles.location_raw_address ILIKE ? OR
				users.first_name ILIKE ? OR
				users.last_name ILIKE ?`,
				searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Get total count before pagination
	var total int64
	countQuery := query.Session(&gorm.Session{})
	countQuery.Count(&total)

	// Apply ordering
	switch queryParams.SortBy {
	case "distance":
		if strategy == "radius" {
			distanceExpr := utils.HaversineSQL("location_latitude", "location_longitude", centerLat, centerLng)
			query = query.Order(distanceExpr + " ASC")
		} else {
			query = query.Order("created_at DESC")
		}
	case "rate":
		query = query.Order("hourly_rate ASC")
	case "recent":
		query = query.Order("created_at DESC")
	default:
		if strategy == "radius" {
			distanceExpr := utils.HaversineSQL("location_latitude", "location_longitude", centerLat, centerLng)
			query = query.Order(distanceExpr + " ASC")
		} else {
			query = query.Order("created_at DESC")
		}
	}

	// Apply pagination
	offset := queryParams.GetOffset()

	var trainers []models.TrainerProfile
	if err := query.Offset(offset).Limit(queryParams.Limit).Find(&trainers).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve trainers")
		return
	}

	// Build responses with review stats
	responses := make([]TrainerSearchResponse, len(trainers))
	for i, trainer := range trainers {
		// Calculate review stats
		var reviewCount int64
		var avgRating float64

		database.DB.Model(&models.TrainerReview{}).Where("trainer_id = ?", trainer.ID).Count(&reviewCount)
		if reviewCount > 0 {
			database.DB.Model(&models.TrainerReview{}).Where("trainer_id = ?", trainer.ID).Select("COALESCE(AVG(rating), 0)").Scan(&avgRating)
		}

		// Convert specialties
		specialties := make([]models.SpecialtyResponse, 0, len(trainer.Specialties))
		for _, s := range trainer.Specialties {
			specialties = append(specialties, s.ToResponse())
		}

		resp := TrainerSearchResponse{
			ID:                  trainer.ID,
			UserID:              trainer.UserID,
			Bio:                 trainer.Bio,
			Specialties:         specialties,
			HourlyRate:          trainer.HourlyRate,
			Location:            trainer.Location.ToResponse(),
			Visibility:          trainer.Visibility,
			IsLookingForClients: trainer.IsLookingForClients,
			ReviewCount:         int(reviewCount),
			AverageRating:       avgRating,
			CreatedAt:           trainer.CreatedAt,
		}

		if trainer.User.ID != uuid.Nil {
			resp.User = &models.UserPublicResponse{
				ID:        trainer.User.ID,
				FirstName: trainer.User.FirstName,
				LastName:  trainer.User.LastName,
			}
		}

		// Calculate distance if radius search and trainer has coordinates
		if strategy == "radius" && trainer.Location.HasCoordinates() {
			distance := utils.HaversineDistance(centerLat, centerLng, *trainer.Location.Latitude, *trainer.Location.Longitude)
			resp.DistanceKm = &distance
		}

		responses[i] = resp
	}

	// Filter by min_rating if specified (post-query filter since we need to calculate ratings)
	if queryParams.MinRating > 0 {
		filteredResponses := make([]TrainerSearchResponse, 0)
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
