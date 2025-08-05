package controllers

import (
	"fit-flow-api/config"
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"
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
		utils.HandleBindingError(c, err)
		return
	}

	claims, err := apple.GetClaims(req.IdentityToken)
	if err != nil {
		utils.UnauthorizedResponse(c, "Invalid Apple identity token.")
		return
	}

	email, ok := (*claims)["email"].(string)
	if !ok {
		utils.BadRequestResponse(c, "Email not found in token.", nil)
		return
	}

	subject, ok := (*claims)["sub"].(string)
	if !ok {
		utils.BadRequestResponse(c, "Subject not found in token.", nil)
		return
	}

	issuer, ok := (*claims)["iss"].(string)
	if !ok || issuer != "https://appleid.apple.com" {
		utils.UnauthorizedResponse(c, "Invalid token issuer.")
		return
	}

	audience, ok := (*claims)["aud"].(string)
	if !ok || audience != config.AppConfig.AppleClientID {
		utils.UnauthorizedResponse(c, "Invalid token audience.")
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
			utils.InternalServerErrorResponse(c, "Failed to create user.")
			return
		}
	} else {
		if user.AppleID == "" {
			user.AppleID = subject
			database.DB.Save(&user)
		}
	}
	
	if !user.IsActive {
		utils.UnauthorizedResponse(c, "Account is deactivated.")
		return
	}
	
	jwtToken, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to generate token.")
		return
	}
	
	response := AuthResponse{
		User:  user.ToResponse(),
		Token: jwtToken,
	}
	
	utils.SuccessResponse(c, "Apple login successful.", response)
}