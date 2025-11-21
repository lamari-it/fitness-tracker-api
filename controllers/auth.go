package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TrainerProfileData struct {
	Bio          string      `json:"bio" binding:"required,min=10,max=1000"`
	SpecialtyIDs []uuid.UUID `json:"specialty_ids" binding:"required,min=1,max=20"`
	HourlyRate   float64     `json:"hourly_rate" binding:"required,gt=0,lte=9999.99"`
	Location     string      `json:"location" binding:"required,min=2,max=500"`
	Visibility   string      `json:"visibility" binding:"omitempty,oneof=public link_only private"`
}

type RegisterRequest struct {
	Email           string              `json:"email" binding:"required,email,max=255"`
	Password        string              `json:"password" binding:"required,min=8,max=128"`
	PasswordConfirm string              `json:"password_confirm" binding:"required"`
	FirstName       string              `json:"first_name" binding:"required,min=1,max=100"`
	LastName        string              `json:"last_name" binding:"required,min=1,max=100"`
	TrainerProfile  *TrainerProfileData `json:"trainer_profile" binding:"omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=1,max=128"`
}

type AuthResponse struct {
	User           models.UserResponse            `json:"user"`
	Token          string                         `json:"token"`
	TrainerProfile *models.TrainerProfileResponse `json:"trainer_profile,omitempty"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	if req.Password != req.PasswordConfirm {
		validationErrors := utils.ValidationErrors{
			"password_confirm": []string{"Passwords do not match."},
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		utils.ConflictResponse(c, "A user with this email already exists.")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to process password.")
		return
	}

	user := models.User{
		Email:     strings.ToLower(req.Email),
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Provider:  "local",
		IsActive:  true,
	}

	// Use transaction if trainer profile is provided
	var trainerProfile *models.TrainerProfile
	if req.TrainerProfile != nil {
		// Validate that all specialty IDs exist
		var specialties []models.Specialty
		if err := database.DB.Where("id IN ?", req.TrainerProfile.SpecialtyIDs).Find(&specialties).Error; err != nil {
			utils.InternalServerErrorResponse(c, "Failed to validate specialties")
			return
		}
		if len(specialties) != len(req.TrainerProfile.SpecialtyIDs) {
			utils.BadRequestResponse(c, "One or more specialty IDs are invalid", nil)
			return
		}

		// Validate trainer profile data
		visibility := req.TrainerProfile.Visibility
		if visibility == "" {
			visibility = "public"
		}

		tp := models.TrainerProfile{
			Bio:         req.TrainerProfile.Bio,
			Specialties: specialties,
			HourlyRate:  req.TrainerProfile.HourlyRate,
			Location:    req.TrainerProfile.Location,
			Visibility:  visibility,
		}

		if err := tp.Validate(); err != nil {
			utils.BadRequestResponse(c, "Trainer profile validation failed", err.Error())
			return
		}

		// Create user and trainer profile in transaction
		tx := database.DB.Begin()
		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			utils.InternalServerErrorResponse(c, "Failed to create user.")
			return
		}

		tp.UserID = user.ID
		if err := tx.Create(&tp).Error; err != nil {
			tx.Rollback()
			utils.InternalServerErrorResponse(c, "Failed to create trainer profile.")
			return
		}

		// Associate specialties with the trainer profile
		if err := tx.Model(&tp).Association("Specialties").Replace(&specialties); err != nil {
			tx.Rollback()
			utils.InternalServerErrorResponse(c, "Failed to associate specialties")
			return
		}

		if err := tx.Commit().Error; err != nil {
			utils.InternalServerErrorResponse(c, "Failed to complete registration.")
			return
		}

		// Preload user and specialties for response
		database.DB.Preload("User").Preload("Specialties").First(&tp, "id = ?", tp.ID)
		trainerProfile = &tp
	} else {
		if err := database.DB.Create(&user).Error; err != nil {
			utils.InternalServerErrorResponse(c, "Failed to create user.")
			return
		}
	}

	// Check for pending email invitations and create TrainerClientLinks
	var pendingInvitations []models.TrainerInvitation
	if err := database.DB.Where("invitee_email = ? AND status = ?", strings.ToLower(req.Email), models.InvitationStatusPending).Find(&pendingInvitations).Error; err == nil {
		for _, invitation := range pendingInvitations {
			// Skip expired invitations
			if invitation.IsExpired() {
				invitation.Status = models.InvitationStatusExpired
				database.DB.Save(&invitation)
				continue
			}

			// Create pending TrainerClientLink
			clientLink := models.TrainerClientLink{
				TrainerID: invitation.TrainerID,
				ClientID:  user.ID,
				Status:    "pending",
			}
			database.DB.Create(&clientLink)

			// Mark invitation as accepted
			invitation.Status = models.InvitationStatusAccepted
			now := database.DB.NowFunc()
			invitation.AcceptedAt = &now
			database.DB.Save(&invitation)
		}
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to generate authentication token.")
		return
	}

	response := AuthResponse{
		User:  user.ToResponse(),
		Token: token,
	}

	if trainerProfile != nil {
		profileResponse := trainerProfile.ToResponse()
		response.TrainerProfile = &profileResponse
	}

	utils.CreatedResponse(c, "User registered successfully.", response)
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", strings.ToLower(req.Email)).First(&user).Error; err != nil {
		utils.UnauthorizedResponse(c, "Invalid email or password.")
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		utils.UnauthorizedResponse(c, "Invalid email or password.")
		return
	}

	if !user.IsActive {
		utils.ForbiddenResponse(c, "Your account has been deactivated.")
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to generate authentication token.")
		return
	}

	response := AuthResponse{
		User:  user.ToResponse(),
		Token: token,
	}

	utils.SuccessResponse(c, "Login successful.", response)
}

func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", userID.(uuid.UUID)).First(&user).Error; err != nil {
		utils.NotFoundResponse(c, "User not found.")
		return
	}

	utils.SuccessResponse(c, "Profile fetched successfully.", user.ToResponse())
}
