package controllers

import (
	"lamari-fit-api/database"
	"lamari-fit-api/models"
	"lamari-fit-api/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TrainerProfileData struct {
	Bio          string      `json:"bio" binding:"omitempty,max=1000"`
	SpecialtyIDs []uuid.UUID `json:"specialty_ids" binding:"omitempty,max=20"`
	HourlyRate   *float64    `json:"hourly_rate,omitempty" binding:"omitempty,gte=0,lte=9999.99"`
	Location     string      `json:"location" binding:"omitempty,max=500"`
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
	Email      string `json:"email" binding:"required,email,max=255"`
	Password   string `json:"password" binding:"required,min=1,max=128"`
	DeviceInfo string `json:"device_info" binding:"omitempty,max=255"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	DeviceInfo   string `json:"device_info" binding:"omitempty,max=255"`
}

type AuthResponse struct {
	User           models.UserResponse            `json:"user"`
	AccessToken    string                         `json:"access_token"`
	RefreshToken   string                         `json:"refresh_token"`
	ExpiresIn      int64                          `json:"expires_in"`
	TokenType      string                         `json:"token_type"`
	Token          string                         `json:"token"` // Deprecated: use access_token
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
			visibility = "private"
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

		// Assign "user" and "trainer" roles
		var userRole, trainerRole models.Role
		if err := tx.Where("name = ?", "user").First(&userRole).Error; err == nil {
			err := tx.Model(&user).Association("Roles").Append(&userRole)
			if err != nil {
				tx.Rollback()
				utils.InternalServerErrorResponse(c, "Failed to assign user role.")
				return
			}
		}
		if err := tx.Where("name = ?", "trainer").First(&trainerRole).Error; err == nil {
			err := tx.Model(&user).Association("Roles").Append(&trainerRole)
			if err != nil {
				tx.Rollback()
				utils.InternalServerErrorResponse(c, "Failed to assign user role.")
				return
			}
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

		// Assign "user" role
		var userRole models.Role
		if err := database.DB.Where("name = ?", "user").First(&userRole).Error; err == nil {
			database.DB.Model(&user).Association("Roles").Append(&userRole)
		}
	}

	// Preload roles for response
	database.DB.Preload("Roles").First(&user, "id = ?", user.ID)

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

	accessToken, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to generate authentication token.")
		return
	}

	// Generate refresh token
	refreshToken, tokenHash, err := utils.GenerateRefreshToken()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to generate refresh token.")
		return
	}

	// Store refresh token in database
	refreshTokenRecord := models.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		IPAddress: c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		ExpiresAt: utils.GetRefreshTokenExpiration(),
	}
	if err := database.DB.Create(&refreshTokenRecord).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create session.")
		return
	}

	response := AuthResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    utils.GetAccessTokenExpiresIn(),
		TokenType:    "Bearer",
		Token:        accessToken, // Deprecated: backward compatibility
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
	if err := database.DB.Preload("Roles").Where("email = ?", strings.ToLower(req.Email)).First(&user).Error; err != nil {
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

	accessToken, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to generate authentication token.")
		return
	}

	// Generate refresh token
	refreshToken, tokenHash, err := utils.GenerateRefreshToken()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to generate refresh token.")
		return
	}

	// Store refresh token in database
	refreshTokenRecord := models.RefreshToken{
		UserID:     user.ID,
		TokenHash:  tokenHash,
		DeviceInfo: req.DeviceInfo,
		IPAddress:  c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
		ExpiresAt:  utils.GetRefreshTokenExpiration(),
	}
	if err := database.DB.Create(&refreshTokenRecord).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create session.")
		return
	}

	response := AuthResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    utils.GetAccessTokenExpiresIn(),
		TokenType:    "Bearer",
		Token:        accessToken, // Deprecated: backward compatibility
	}

	utils.SuccessResponse(c, "Login successful.", response)
}

// RefreshToken exchanges a refresh token for a new access/refresh token pair
func RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Hash the provided token to look it up
	tokenHash := utils.HashRefreshToken(req.RefreshToken)

	// Find the refresh token
	var refreshTokenRecord models.RefreshToken
	if err := database.DB.Where("token_hash = ?", tokenHash).First(&refreshTokenRecord).Error; err != nil {
		utils.UnauthorizedResponse(c, "Invalid refresh token.")
		return
	}

	// Check if token was revoked (potential token theft - revoke all user tokens)
	if refreshTokenRecord.IsRevoked() {
		// Security: revoke ALL tokens for this user (reuse detection)
		database.DB.Model(&models.RefreshToken{}).
			Where("user_id = ? AND revoked_at IS NULL", refreshTokenRecord.UserID).
			Update("revoked_at", time.Now())
		utils.UnauthorizedResponse(c, "Refresh token has been revoked. Please login again.")
		return
	}

	// Check if token is expired
	if refreshTokenRecord.IsExpired() {
		utils.UnauthorizedResponse(c, "Refresh token has expired. Please login again.")
		return
	}

	// Get the user
	var user models.User
	if err := database.DB.Preload("Roles").First(&user, "id = ?", refreshTokenRecord.UserID).Error; err != nil {
		utils.UnauthorizedResponse(c, "User not found.")
		return
	}

	if !user.IsActive {
		utils.ForbiddenResponse(c, "Your account has been deactivated.")
		return
	}

	// Generate new access token
	accessToken, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to generate access token.")
		return
	}

	// Generate new refresh token (rotation)
	newRefreshToken, newTokenHash, err := utils.GenerateRefreshToken()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to generate refresh token.")
		return
	}

	// Revoke old refresh token and create new one
	tx := database.DB.Begin()

	// Revoke the old token
	refreshTokenRecord.Revoke()
	if err := tx.Save(&refreshTokenRecord).Error; err != nil {
		tx.Rollback()
		utils.InternalServerErrorResponse(c, "Failed to rotate refresh token.")
		return
	}

	// Create new refresh token record
	deviceInfo := req.DeviceInfo
	if deviceInfo == "" {
		deviceInfo = refreshTokenRecord.DeviceInfo // Keep existing device info
	}

	newRefreshTokenRecord := models.RefreshToken{
		UserID:     user.ID,
		TokenHash:  newTokenHash,
		DeviceInfo: deviceInfo,
		IPAddress:  c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
		ExpiresAt:  utils.GetRefreshTokenExpiration(),
	}
	if err := tx.Create(&newRefreshTokenRecord).Error; err != nil {
		tx.Rollback()
		utils.InternalServerErrorResponse(c, "Failed to create new session.")
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to complete token refresh.")
		return
	}

	response := AuthResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    utils.GetAccessTokenExpiresIn(),
		TokenType:    "Bearer",
		Token:        accessToken, // Deprecated: backward compatibility
	}

	utils.SuccessResponse(c, "Token refreshed successfully.", response)
}

// Logout revokes the current session's refresh token
func Logout(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Hash the provided token
	tokenHash := utils.HashRefreshToken(req.RefreshToken)

	// Find and revoke the refresh token
	result := database.DB.Model(&models.RefreshToken{}).
		Where("token_hash = ? AND revoked_at IS NULL", tokenHash).
		Update("revoked_at", time.Now())

	if result.RowsAffected == 0 {
		utils.NotFoundResponse(c, "Session not found or already revoked.")
		return
	}

	utils.SuccessResponse(c, "Logged out successfully.", nil)
}

// LogoutAll revokes all refresh tokens for the current user
func LogoutAll(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	// Revoke all refresh tokens for this user
	result := database.DB.Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID.(uuid.UUID)).
		Update("revoked_at", time.Now())

	utils.SuccessResponse(c, "Logged out from all devices successfully.", map[string]int64{
		"sessions_revoked": result.RowsAffected,
	})
}

// GetSessions returns all active sessions for the current user
func GetSessions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	// Get current token hash from header (if available) to mark current session
	currentTokenHash := ""
	if authHeader := c.GetHeader("X-Refresh-Token"); authHeader != "" {
		currentTokenHash = utils.HashRefreshToken(authHeader)
	}

	var refreshTokens []models.RefreshToken
	if err := database.DB.Where("user_id = ? AND revoked_at IS NULL AND expires_at > ?",
		userID.(uuid.UUID), time.Now()).
		Order("created_at DESC").
		Find(&refreshTokens).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch sessions.")
		return
	}

	sessions := make([]models.SessionResponse, len(refreshTokens))
	for i, rt := range refreshTokens {
		sessions[i] = rt.ToSessionResponse(currentTokenHash)
	}

	utils.SuccessResponse(c, "Sessions fetched successfully.", sessions)
}

// RevokeSession revokes a specific session by its ID
func RevokeSession(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid session ID.", nil)
		return
	}

	// Find and revoke the session (only if it belongs to the current user)
	result := database.DB.Model(&models.RefreshToken{}).
		Where("id = ? AND user_id = ? AND revoked_at IS NULL", sessionID, userID.(uuid.UUID)).
		Update("revoked_at", time.Now())

	if result.RowsAffected == 0 {
		utils.NotFoundResponse(c, "Session not found or already revoked.")
		return
	}

	utils.SuccessResponse(c, "Session revoked successfully.", nil)
}

func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated.")
		return
	}

	var user models.User
	if err := database.DB.Preload("Roles").Where("id = ?", userID.(uuid.UUID)).First(&user).Error; err != nil {
		utils.NotFoundResponse(c, "User not found.")
		return
	}

	utils.SuccessResponse(c, "Profile fetched successfully.", user.ToResponse())
}
