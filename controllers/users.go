package controllers

import (
	"lamari-fit-api/database"
	"lamari-fit-api/models"
	"lamari-fit-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DiscoverUsers returns a paginated list of public users with optional search/filter
// GET /api/v1/users
func DiscoverUsers(c *gin.Context) {
	var queryParams UserDiscoveryQuery
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	SetDefaultPagination(&queryParams.PaginationQuery)

	query := database.DB.Model(&models.User{}).Where("profile_visibility = ?", "public")

	// Search by first_name, last_name, or email
	if queryParams.Search != "" {
		searchPattern := "%" + queryParams.Search + "%"
		query = query.Where(
			"first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	// Filter by is_looking_for_trainer
	if queryParams.IsLookingForTrainer == "true" {
		query = query.Where("is_looking_for_trainer = ?", true)
	} else if queryParams.IsLookingForTrainer == "false" {
		query = query.Where("is_looking_for_trainer = ?", false)
	}

	// Get total count before pagination
	var total int64
	countQuery := query.Session(&gorm.Session{})
	countQuery.Count(&total)

	// Apply pagination
	offset := queryParams.GetOffset()

	var users []models.User
	if err := query.Order("created_at DESC").Offset(offset).Limit(queryParams.Limit).Find(&users).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve users")
		return
	}

	// Convert to discovery responses (limited info)
	responses := make([]models.UserDiscoveryResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToDiscoveryResponse()
	}

	utils.PaginatedResponse(c, "Users retrieved successfully", responses, queryParams.Page, queryParams.Limit, int(total))
}

// GetUserPublicProfile returns a user's profile with privacy enforcement
// GET /api/v1/users/:id
func GetUserPublicProfile(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	targetID, ok := utils.ParseUUID(c, params.ID, "user")
	if !ok {
		return
	}

	// Get current user ID
	currentUserID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	// Fetch target user
	var targetUser models.User
	if err := database.DB.Preload("Roles").First(&targetUser, "id = ?", targetID).Error; err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	// Privacy check
	if canViewProfile(currentUserID, targetID, &targetUser) {
		// If viewing self, return full response
		if currentUserID == targetID {
			utils.SuccessResponse(c, "User profile retrieved successfully", targetUser.ToResponse())
			return
		}
		// Otherwise return discovery response (limited info)
		utils.SuccessResponse(c, "User profile retrieved successfully", targetUser.ToDiscoveryResponse())
		return
	}

	// Not allowed to view
	utils.NotFoundResponse(c, "User not found")
}

// canViewProfile checks if a viewer can view a target user's profile
func canViewProfile(viewerID, targetID uuid.UUID, target *models.User) bool {
	// Self can always view
	if viewerID == targetID {
		return true
	}

	switch target.ProfileVisibility {
	case "public":
		return true
	case "friends_only":
		return areFriends(viewerID, targetID)
	case "private":
		return false
	}
	return false
}

// areFriends checks if two users have an accepted friendship
func areFriends(userID1, userID2 uuid.UUID) bool {
	var count int64
	database.DB.Model(&models.Friendship{}).
		Where("((user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)) AND status = ?",
			userID1, userID2, userID2, userID1, "accepted").
		Count(&count)
	return count > 0
}
