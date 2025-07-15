package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/middleware"
	"fit-flow-api/models"
	"fit-flow-api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required,min=6"`
	PasswordConfirm  string `json:"password_confirm" binding:"required"`
	FirstName        string `json:"first_name" binding:"required"`
	LastName         string `json:"last_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User  models.UserResponse `json:"user"`
	Token string              `json:"token"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
		return
	}

	if req.Password != req.PasswordConfirm {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "auth.password_mismatch", nil)
		return
	}

	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		middleware.TranslateErrorResponse(c, http.StatusConflict, "user.email_exists", nil)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "general.internal_error", nil)
		return
	}

	user := models.User{
		Email:     strings.ToLower(req.Email),
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Provider:  "local",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "auth.register_failed", nil)
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "general.internal_error", nil)
		return
	}

	response := AuthResponse{
		User:  user.ToResponse(),
		Token: token,
	}

	middleware.TranslateResponse(c, http.StatusCreated, "user.register_success", response)
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_format", err.Error())
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", strings.ToLower(req.Email)).First(&user).Error; err != nil {
		middleware.TranslateErrorResponse(c, http.StatusUnauthorized, "auth.login_failed", nil)
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		middleware.TranslateErrorResponse(c, http.StatusUnauthorized, "auth.login_failed", nil)
		return
	}

	if !user.IsActive {
		middleware.TranslateErrorResponse(c, http.StatusUnauthorized, "auth.unauthorized", nil)
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		middleware.TranslateErrorResponse(c, http.StatusInternalServerError, "general.internal_error", nil)
		return
	}

	response := AuthResponse{
		User:  user.ToResponse(),
		Token: token,
	}

	middleware.TranslateResponse(c, http.StatusOK, "auth.login_success", response)
}

func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", userID.(uuid.UUID)).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}
