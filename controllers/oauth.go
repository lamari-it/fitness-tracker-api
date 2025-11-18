package controllers

import (
	"context"
	"encoding/json"
	"fit-flow-api/config"
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

func getGoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.AppConfig.GoogleClientID,
		ClientSecret: config.AppConfig.GoogleClientSecret,
		RedirectURL:  config.AppConfig.GoogleRedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

func GoogleLogin(c *gin.Context) {
	oauthConfig := getGoogleOAuthConfig()

	state := uuid.New().String()
	c.SetCookie("oauth_state", state, 600, "/", "", false, true)

	url := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	storedState, err := c.Cookie("oauth_state")
	if err != nil || state != storedState {
		utils.BadRequestResponse(c, "Invalid state parameter.", nil)
		return
	}

	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	oauthConfig := getGoogleOAuthConfig()
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to exchange code for token.")
		return
	}

	client := oauthConfig.Client(context.Background(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get user info.")
		return
	}
	defer response.Body.Close()

	var googleUser GoogleUser
	if err := json.NewDecoder(response.Body).Decode(&googleUser); err != nil {
		utils.InternalServerErrorResponse(c, "Failed to decode user info.")
		return
	}

	var user models.User
	err = database.DB.Where("google_id = ? OR email = ?", googleUser.ID, strings.ToLower(googleUser.Email)).First(&user).Error

	if err != nil {
		googleID := googleUser.ID
		user = models.User{
			Email:     strings.ToLower(googleUser.Email),
			FirstName: googleUser.GivenName,
			LastName:  googleUser.FamilyName,
			Provider:  "google",
			GoogleID:  &googleID,
			Password:  uuid.New().String(),
		}

		if err := database.DB.Create(&user).Error; err != nil {
			utils.InternalServerErrorResponse(c, "Failed to create user.")
			return
		}
	} else {
		if user.GoogleID == nil || *user.GoogleID == "" {
			googleID := googleUser.ID
			user.GoogleID = &googleID
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

	response2 := AuthResponse{
		User:  user.ToResponse(),
		Token: jwtToken,
	}

	utils.SuccessResponse(c, "Google login successful.", response2)
}
