package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateEmailInvitation creates an email-based invitation for a potential client
func CreateEmailInvitation(c *gin.Context) {
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

	var req models.CreateEmailInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Normalize email to lowercase
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Get trainer info for email
	var trainer models.User
	if err := database.DB.First(&trainer, userID).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get trainer info")
		return
	}

	// Check if trainer is trying to invite themselves
	if trainer.Email == email {
		utils.BadRequestResponse(c, "You cannot invite yourself", nil)
		return
	}

	// Check if user already exists and has an active/pending relationship
	var existingUser models.User
	if err := database.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		// User exists, check for existing relationship
		var existingLink models.TrainerClientLink
		if err := database.DB.Where("trainer_id = ? AND client_id = ?", userID, existingUser.ID).First(&existingLink).Error; err == nil {
			if existingLink.Status == "active" {
				utils.ConflictResponse(c, "This user is already your client")
				return
			}
			if existingLink.Status == "pending" {
				utils.ConflictResponse(c, "You already have a pending invitation for this user")
				return
			}
		}
	}

	// Check for existing pending email invitation
	var existingInvitation models.TrainerInvitation
	if err := database.DB.Where("trainer_id = ? AND invitee_email = ? AND status = ?", userID, email, models.InvitationStatusPending).First(&existingInvitation).Error; err == nil {
		utils.ConflictResponse(c, "You already have a pending invitation for this email")
		return
	}

	// Generate invitation token
	token, err := utils.GenerateInvitationToken()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to generate invitation token")
		return
	}

	// Create invitation (expires in 7 days)
	invitation := models.TrainerInvitation{
		TrainerID:       userID,
		InviteeEmail:    email,
		InvitationToken: token,
		Status:          models.InvitationStatusPending,
		ExpiresAt:       time.Now().AddDate(0, 0, 7),
	}

	if err := database.DB.Create(&invitation).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create invitation")
		return
	}

	// Send invitation email
	emailService := utils.NewEmailService()
	trainerName := trainer.FirstName + " " + trainer.LastName
	if err := emailService.SendTrainerInvitation(email, trainerName, token); err != nil {
		// Log error but don't fail the request
		// Invitation is still created in database
		// Could add to a queue for retry
	}

	// Load trainer relationship for response
	database.DB.Preload("Trainer").First(&invitation, invitation.ID)

	utils.CreatedResponse(c, "Invitation sent successfully", invitation.ToResponse())
}

// GetEmailInvitations gets all email invitations sent by the trainer
func GetEmailInvitations(c *gin.Context) {
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
		utils.ForbiddenResponse(c, "You must have a trainer profile to view invitations")
		return
	}

	var invitations []models.TrainerInvitation
	if err := database.DB.Where("trainer_id = ?", userID).
		Order("created_at DESC").
		Find(&invitations).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve invitations")
		return
	}

	// Convert to response format
	responses := make([]models.TrainerInvitationResponse, len(invitations))
	for i, inv := range invitations {
		responses[i] = inv.ToResponse()
	}

	utils.SuccessResponse(c, "Invitations retrieved successfully", responses)
}

// CancelEmailInvitation cancels a pending email invitation
func CancelEmailInvitation(c *gin.Context) {
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

	invitationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid invitation ID", nil)
		return
	}

	var invitation models.TrainerInvitation
	if err := database.DB.Where("id = ? AND trainer_id = ?", invitationID, userID).First(&invitation).Error; err != nil {
		utils.NotFoundResponse(c, "Invitation not found")
		return
	}

	if invitation.Status != models.InvitationStatusPending {
		utils.BadRequestResponse(c, "Only pending invitations can be cancelled", nil)
		return
	}

	if err := database.DB.Delete(&invitation).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to cancel invitation")
		return
	}

	utils.NoContentResponse(c)
}

// VerifyInvitationToken verifies an invitation token and returns invitation details
func VerifyInvitationToken(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		utils.BadRequestResponse(c, "Token is required", nil)
		return
	}

	var invitation models.TrainerInvitation
	if err := database.DB.Preload("Trainer").Where("invitation_token = ?", token).First(&invitation).Error; err != nil {
		utils.SuccessResponse(c, "Invalid invitation", models.VerifyInvitationResponse{
			Valid:   false,
			Message: "This invitation link is invalid or has already been used",
		})
		return
	}

	// Check if expired
	if invitation.IsExpired() {
		utils.SuccessResponse(c, "Expired invitation", models.VerifyInvitationResponse{
			Valid:   false,
			Message: "This invitation has expired",
		})
		return
	}

	// Check if already used
	if invitation.Status != models.InvitationStatusPending {
		utils.SuccessResponse(c, "Used invitation", models.VerifyInvitationResponse{
			Valid:   false,
			Message: "This invitation has already been used",
		})
		return
	}

	utils.SuccessResponse(c, "Valid invitation", models.VerifyInvitationResponse{
		Valid:        true,
		InviteeEmail: invitation.InviteeEmail,
		Trainer: &models.TrainerInfoResponse{
			ID:        invitation.Trainer.ID,
			Email:     invitation.Trainer.Email,
			FirstName: invitation.Trainer.FirstName,
			LastName:  invitation.Trainer.LastName,
		},
		ExpiresAt: invitation.ExpiresAt,
	})
}

// ResendEmailInvitation resends an invitation email
func ResendEmailInvitation(c *gin.Context) {
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

	invitationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid invitation ID", nil)
		return
	}

	var invitation models.TrainerInvitation
	if err := database.DB.Preload("Trainer").Where("id = ? AND trainer_id = ?", invitationID, userID).First(&invitation).Error; err != nil {
		utils.NotFoundResponse(c, "Invitation not found")
		return
	}

	if invitation.Status != models.InvitationStatusPending {
		utils.BadRequestResponse(c, "Only pending invitations can be resent", nil)
		return
	}

	// Generate new token and extend expiration
	newToken, err := utils.GenerateInvitationToken()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to generate new token")
		return
	}

	invitation.InvitationToken = newToken
	invitation.ExpiresAt = time.Now().AddDate(0, 0, 7)

	if err := database.DB.Save(&invitation).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update invitation")
		return
	}

	// Resend email
	emailService := utils.NewEmailService()
	trainerName := invitation.Trainer.FirstName + " " + invitation.Trainer.LastName
	emailService.SendTrainerInvitation(invitation.InviteeEmail, trainerName, newToken)

	utils.SuccessResponse(c, "Invitation resent successfully", invitation.ToResponse())
}
