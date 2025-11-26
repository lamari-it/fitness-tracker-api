package controllers

import (
	"lamari-fit-api/config"
	"lamari-fit-api/database"
	"lamari-fit-api/models"
	"lamari-fit-api/utils"
	"strings"

	"github.com/Timothylock/go-signin-with-apple/apple"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AppleLoginRequest struct {
	IdentityToken string `json:"identity_token" binding:"required,min=1"`
	AuthCode      string `json:"auth_code" binding:"omitempty,min=1"`
	FirstName     string `json:"first_name" binding:"omitempty,min=1,max=100"`
	LastName      string `json:"last_name" binding:"omitempty,min=1,max=100"`
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

		appleID := subject
		user = models.User{
			Email:     strings.ToLower(email),
			FirstName: firstName,
			LastName:  lastName,
			Provider:  "apple",
			AppleID:   &appleID,
			Password:  uuid.New().String(),
		}

		if err := database.DB.Create(&user).Error; err != nil {
			utils.InternalServerErrorResponse(c, "Failed to create user.")
			return
		}
	} else {
		if user.AppleID == nil || *user.AppleID == "" {
			appleID := subject
			user.AppleID = &appleID
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
