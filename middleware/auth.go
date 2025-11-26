package middleware

import (
	"lamari-fit-api/database"
	"lamari-fit-api/models"
	"lamari-fit-api/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Authorization header is required")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			utils.UnauthorizedResponse(c, "Bearer token is required")
			c.Abort()
			return
		}

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			utils.UnauthorizedResponse(c, "Invalid token")
			c.Abort()
			return
		}

		// Verify user exists in database and is not deleted
		var user models.User
		if err := database.DB.First(&user, "id = ?", claims.UserID).Error; err != nil {
			utils.UnauthorizedResponse(c, "Invalid token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}
