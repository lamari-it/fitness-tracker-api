package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SendFriendRequestRequest struct {
	FriendEmail string `json:"friend_email" binding:"required,email"`
}

func SendFriendRequest(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req SendFriendRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var friend models.User
	if err := database.DB.Where("email = ?", req.FriendEmail).First(&friend).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if friend.ID == userID.(uuid.UUID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot send friend request to yourself"})
		return
	}

	var existingFriendship models.Friendship
	if err := database.DB.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", 
		userID, friend.ID, friend.ID, userID).First(&existingFriendship).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Friend request already exists or you are already friends"})
		return
	}

	friendship := models.Friendship{
		UserID:   userID.(uuid.UUID),
		FriendID: friend.ID,
		Status:   "pending",
	}

	if err := database.DB.Create(&friendship).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send friend request"})
		return
	}

	database.DB.Preload("Friend").First(&friendship, friendship.ID)

	c.JSON(http.StatusCreated, friendship.ToResponse())
}

func GetFriendRequests(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var friendships []models.Friendship
	if err := database.DB.Where("friend_id = ? AND status = ?", userID, "pending").
		Preload("User").
		Find(&friendships).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch friend requests"})
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

	c.JSON(http.StatusOK, gin.H{"requests": responses})
}

func RespondToFriendRequest(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	action := c.Param("action")
	if action != "accept" && action != "decline" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action. Use 'accept' or 'decline'"})
		return
	}

	var friendship models.Friendship
	if err := database.DB.Where("id = ? AND friend_id = ? AND status = ?", 
		requestID, userID, "pending").First(&friendship).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Friend request not found"})
		return
	}

	if action == "accept" {
		friendship.Status = "accepted"
	} else {
		friendship.Status = "declined"
	}

	if err := database.DB.Save(&friendship).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update friend request"})
		return
	}

	database.DB.Preload("User").First(&friendship, friendship.ID)

	c.JSON(http.StatusOK, friendship.ToResponse())
}

func GetFriends(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var friendships []models.Friendship
	if err := database.DB.Where("(user_id = ? OR friend_id = ?) AND status = ?", 
		userID, userID, "accepted").
		Preload("User").
		Preload("Friend").
		Find(&friendships).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch friends"})
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

	c.JSON(http.StatusOK, gin.H{"friends": responses})
}

func RemoveFriend(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	friendshipID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friendship ID"})
		return
	}

	result := database.DB.Where("id = ? AND (user_id = ? OR friend_id = ?)", 
		friendshipID, userID, userID).Delete(&models.Friendship{})
	
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove friend"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Friendship not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend removed successfully"})
}