package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// InviteClient allows a trainer to invite a client
func InviteClient(c *gin.Context) {
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

	// Check if user has a trainer profile
	var trainerProfile models.TrainerProfile
	if err := database.DB.Where("user_id = ?", userID).First(&trainerProfile).Error; err != nil {
		utils.ForbiddenResponse(c, "You must have a trainer profile to invite clients")
		return
	}

	var req models.InviteClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Check if client user exists
	var clientUser models.User
	if err := database.DB.First(&clientUser, "id = ?", req.ClientID).Error; err != nil {
		utils.NotFoundResponse(c, "Client user not found")
		return
	}

	// Cannot invite yourself
	if req.ClientID == userID {
		utils.BadRequestResponse(c, "You cannot invite yourself as a client", nil)
		return
	}

	// Check if relationship already exists
	var existingLink models.TrainerClientLink
	if err := database.DB.Where("trainer_id = ? AND client_id = ?", userID, req.ClientID).First(&existingLink).Error; err == nil {
		if existingLink.Status == "active" {
			utils.ConflictResponse(c, "This client is already linked to you")
			return
		} else if existingLink.Status == "pending" {
			utils.ConflictResponse(c, "An invitation is already pending for this client")
			return
		}
		// If inactive, we can reactivate by creating a new pending invitation
	}

	// Create trainer-client link
	link := models.TrainerClientLink{
		TrainerID: userID,
		ClientID:  req.ClientID,
		Status:    "pending",
	}

	if err := database.DB.Create(&link).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create client invitation")
		return
	}

	// Preload for response
	database.DB.Preload("Trainer").Preload("Client").First(&link, "id = ?", link.ID)

	utils.CreatedResponse(c, "Client invitation sent successfully", link.ToResponse())
}

// GetTrainerClients retrieves all clients for the authenticated trainer
func GetTrainerClients(c *gin.Context) {
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

	// Check if user has a trainer profile
	var trainerProfile models.TrainerProfile
	if err := database.DB.Where("user_id = ?", userID).First(&trainerProfile).Error; err != nil {
		utils.ForbiddenResponse(c, "You must have a trainer profile to view clients")
		return
	}

	// Get status filter from query params
	status := c.Query("status")

	query := database.DB.Preload("Client").Where("trainer_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var links []models.TrainerClientLink
	if err := query.Order("created_at DESC").Find(&links).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve clients")
		return
	}

	// Convert to response format
	responses := make([]models.TrainerClientLinkResponse, len(links))
	for i, link := range links {
		responses[i] = link.ToResponse()
	}

	utils.SuccessResponse(c, "Clients retrieved successfully", responses)
}

// RemoveClient removes or deactivates a client relationship
func RemoveClient(c *gin.Context) {
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

	// Parse link ID from URL
	linkID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid link ID", nil)
		return
	}

	// Find the link
	var link models.TrainerClientLink
	if err := database.DB.First(&link, "id = ?", linkID).Error; err != nil {
		utils.NotFoundResponse(c, "Client relationship not found")
		return
	}

	// Verify ownership
	if link.TrainerID != userID {
		utils.ForbiddenResponse(c, "You can only remove your own clients")
		return
	}

	// Set status to inactive
	link.Status = "inactive"
	if err := database.DB.Save(&link).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to remove client")
		return
	}

	utils.SuccessResponse(c, "Client removed successfully", nil)
}

// GetMyTrainers retrieves all trainers for the authenticated client
func GetMyTrainers(c *gin.Context) {
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

	// Get only active trainer relationships
	var links []models.TrainerClientLink
	if err := database.DB.Preload("Trainer").Where("client_id = ? AND status = ?", userID, "active").Order("created_at DESC").Find(&links).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve trainers")
		return
	}

	// Convert to response format
	responses := make([]models.TrainerClientLinkResponse, len(links))
	for i, link := range links {
		responses[i] = link.ToResponse()
	}

	utils.SuccessResponse(c, "Trainers retrieved successfully", responses)
}

// GetMyTrainerInvitations retrieves pending trainer invitations for the authenticated user
func GetMyTrainerInvitations(c *gin.Context) {
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

	// Get only pending invitations
	var links []models.TrainerClientLink
	if err := database.DB.Preload("Trainer").Where("client_id = ? AND status = ?", userID, "pending").Order("created_at DESC").Find(&links).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve invitations")
		return
	}

	// Convert to response format
	responses := make([]models.TrainerClientLinkResponse, len(links))
	for i, link := range links {
		responses[i] = link.ToResponse()
	}

	utils.SuccessResponse(c, "Trainer invitations retrieved successfully", responses)
}

// RespondToInvitation allows a client to accept or reject a trainer invitation
func RespondToInvitation(c *gin.Context) {
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

	// Parse link ID from URL
	linkID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid invitation ID", nil)
		return
	}

	var req models.RespondToInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Find the invitation
	var link models.TrainerClientLink
	if err := database.DB.First(&link, "id = ?", linkID).Error; err != nil {
		utils.NotFoundResponse(c, "Invitation not found")
		return
	}

	// Verify this invitation is for the current user
	if link.ClientID != userID {
		utils.ForbiddenResponse(c, "This invitation is not for you")
		return
	}

	// Check if invitation is still pending
	if link.Status != "pending" {
		utils.BadRequestResponse(c, "This invitation has already been processed", nil)
		return
	}

	// Process the response
	if req.Action == "accept" {
		link.Status = "active"
		if err := database.DB.Save(&link).Error; err != nil {
			utils.InternalServerErrorResponse(c, "Failed to accept invitation")
			return
		}

		// Preload for response
		database.DB.Preload("Trainer").Preload("Client").First(&link, "id = ?", link.ID)
		utils.SuccessResponse(c, "Invitation accepted successfully", link.ToResponse())
	} else {
		// Reject - delete the link
		if err := database.DB.Delete(&link).Error; err != nil {
			utils.InternalServerErrorResponse(c, "Failed to reject invitation")
			return
		}
		utils.SuccessResponse(c, "Invitation rejected successfully", nil)
	}
}
