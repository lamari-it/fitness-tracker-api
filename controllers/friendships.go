package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SendFriendRequestRequest struct {
	FriendEmail string `json:"friend_email" binding:"required,email,max=255"`
}

func SendFriendRequest(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	var req SendFriendRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	var friend models.User
	if err := database.DB.Where("email = ?", req.FriendEmail).First(&friend).Error; err != nil {
		utils.NotFoundResponse(c, "User not found.")
		return
	}

	if friend.ID == userID.(uuid.UUID) {
		utils.BadRequestResponse(c, "Cannot send friend request to yourself.", nil)
		return
	}

	var existingFriendship models.Friendship
	if err := database.DB.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", 
		userID, friend.ID, friend.ID, userID).First(&existingFriendship).Error; err == nil {
		utils.ConflictResponse(c, "Friend request already exists or you are already friends.")
		return
	}

	friendship := models.Friendship{
		UserID:   userID.(uuid.UUID),
		FriendID: friend.ID,
		Status:   "pending",
	}

	if err := database.DB.Create(&friendship).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to send friend request.")
		return
	}

	database.DB.Preload("Friend").First(&friendship, friendship.ID)

	utils.CreatedResponse(c, "Friend request sent successfully.", friendship.ToResponse())
}

func GetFriendRequests(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	var query PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.HandleBindingError(c, err)
		return
	}
	
	SetDefaultPagination(&query)
	offset := (query.Page - 1) * query.Limit

	// Get total count
	var total int64
	if err := database.DB.Model(&models.Friendship{}).Where("friend_id = ? AND status = ?", userID, "pending").Count(&total).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to count friend requests.")
		return
	}

	var friendships []models.Friendship
	if err := database.DB.Where("friend_id = ? AND status = ?", userID, "pending").
		Preload("User").
		Offset(offset).
		Limit(query.Limit).
		Order("created_at DESC").
		Find(&friendships).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch friend requests.")
		return
	}

	var responses []models.FriendshipResponse
	for _, friendship := range friendships {
		response := models.FriendshipResponse{
			ID:        friendship.ID,
			UserID:    friendship.UserID,
			FriendID:  friendship.FriendID,
			Status:    friendship.Status,
			Friend:    friendship.User.ToResponse(),
			CreatedAt: friendship.CreatedAt,
			UpdatedAt: friendship.UpdatedAt,
		}
		responses = append(responses, response)
	}

	utils.PaginatedResponse(c, "Friend requests retrieved successfully.", responses, query.Page, query.Limit, int(total))
}

func RespondToFriendRequest(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid request ID.", nil)
		return
	}

	action := c.Param("action")
	if action != "accept" && action != "decline" {
		utils.BadRequestResponse(c, "Invalid action. Use 'accept' or 'decline'.", nil)
		return
	}

	var friendship models.Friendship
	if err := database.DB.Where("id = ? AND friend_id = ? AND status = ?", 
		requestID, userID, "pending").First(&friendship).Error; err != nil {
		utils.NotFoundResponse(c, "Friend request not found.")
		return
	}

	if action == "accept" {
		friendship.Status = "accepted"
	} else {
		friendship.Status = "declined"
	}

	if err := database.DB.Save(&friendship).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update friend request.")
		return
	}

	database.DB.Preload("User").First(&friendship, friendship.ID)

	utils.SuccessResponse(c, "Friend request updated successfully.", friendship.ToResponse())
}

func GetFriends(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	var query PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.HandleBindingError(c, err)
		return
	}
	
	SetDefaultPagination(&query)
	offset := (query.Page - 1) * query.Limit

	// Get total count
	var total int64
	if err := database.DB.Model(&models.Friendship{}).
		Where("(user_id = ? OR friend_id = ?) AND status = ?", userID, userID, "accepted").
		Count(&total).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to count friends.")
		return
	}

	var friendships []models.Friendship
	if err := database.DB.Where("(user_id = ? OR friend_id = ?) AND status = ?", 
		userID, userID, "accepted").
		Preload("User").
		Preload("Friend").
		Offset(offset).
		Limit(query.Limit).
		Order("created_at DESC").
		Find(&friendships).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch friends.")
		return
	}

	var responses []models.FriendshipResponse
	for _, friendship := range friendships {
		response := models.FriendshipResponse{
			ID:        friendship.ID,
			UserID:    friendship.UserID,
			FriendID:  friendship.FriendID,
			Status:    friendship.Status,
			CreatedAt: friendship.CreatedAt,
			UpdatedAt: friendship.UpdatedAt,
		}

		if friendship.UserID == userID.(uuid.UUID) {
			response.Friend = friendship.Friend.ToResponse()
		} else {
			response.Friend = friendship.User.ToResponse()
		}

		responses = append(responses, response)
	}

	utils.PaginatedResponse(c, "Friends retrieved successfully.", responses, query.Page, query.Limit, int(total))
}

func RemoveFriend(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	friendshipID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid friendship ID.", nil)
		return
	}

	result := database.DB.Where("id = ? AND (user_id = ? OR friend_id = ?)", 
		friendshipID, userID, userID).Delete(&models.Friendship{})
	
	if result.Error != nil {
		utils.InternalServerErrorResponse(c, "Failed to remove friend.")
		return
	}

	if result.RowsAffected == 0 {
		utils.NotFoundResponse(c, "Friendship not found.")
		return
	}

	utils.DeletedResponse(c, "Friend removed successfully.")
}