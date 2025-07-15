package controllers

import (
	"fit-flow-api/config"
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"
	"net/http"
	"strings"

	"github.com/Timothylock/go-signin-with-apple/apple"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AppleLoginRequest struct {
	IdentityToken string `json:"identity_token" binding:"required"`
	AuthCode      string `json:"auth_code,omitempty"`
	FirstName     string `json:"first_name,omitempty"`
	LastName      string `json:"last_name,omitempty"`
}

func AppleLogin(c *gin.Context) {
	var req AppleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := apple.GetClaims(req.IdentityToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Apple identity token"})
		return
	}

	email, ok := (*claims)["email"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email not found in token"})
		return
	}

	subject, ok := (*claims)["sub"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Subject not found in token"})
		return
	}

	issuer, ok := (*claims)["iss"].(string)
	if !ok || issuer != "https://appleid.apple.com" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token issuer"})
		return
	}

	audience, ok := (*claims)["aud"].(string)
	if !ok || audience != config.AppConfig.AppleClientID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token audience"})
		return
	}

	var user models.User
	err = database.DB.Where("apple_id = ? OR email = ?", subject, strings.ToLower(email)).First(&user).Error
	
	if err != nil {
		firstName := req.FirstName
		lastName := req.LastName
		if firstName == "" {
			firstName = "User"
		}
		if lastName == "" {
			lastName = ""
		}
		
		user = models.User{
			Email:     strings.ToLower(email),
			FirstName: firstName,
			LastName:  lastName,
			Provider:  "apple",
			AppleID:   subject,
			Password:  uuid.New().String(),
		}
		
		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	} else {
		if user.AppleID == "" {
			user.AppleID = subject
			database.DB.Save(&user)
		}
	}
	
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is deactivated"})
		return
	}
	
	jwtToken, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	
	response := AuthResponse{
		User:  user.ToResponse(),
		Token: jwtToken,
	}
	
	c.JSON(http.StatusOK, response)
}