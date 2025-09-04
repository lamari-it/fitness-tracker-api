package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email           string `json:"email" binding:"required,email,max=255"`
	Password        string `json:"password" binding:"required,min=8,max=128"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
	FirstName       string `json:"first_name" binding:"required,min=1,max=100"`
	LastName        string `json:"last_name" binding:"required,min=1,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=1,max=128"`
}

type AuthResponse struct {
	User  models.UserResponse `json:"user"`
	Token string              `json:"token"`
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

	if err := database.DB.Create(&user).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create user.")
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
