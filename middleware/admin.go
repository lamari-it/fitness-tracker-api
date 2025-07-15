package middleware

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AdminMiddleware ensures that only admin users can access protected routes
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthMiddleware)
		userIDStr, exists := c.Get("user_id")
		if !exists {
			TranslateErrorResponse(c, http.StatusUnauthorized, "auth.unauthorized", nil)
			c.Abort()
			return
		}

		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			TranslateErrorResponse(c, http.StatusBadRequest, "validation.invalid_uuid", nil)
			c.Abort()
			return
		}

		// Get user from database
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			TranslateErrorResponse(c, http.StatusUnauthorized, "auth.unauthorized", nil)
			c.Abort()
			return
		}

		// Check if user is admin
		if !user.IsAdmin {
			TranslateErrorResponse(c, http.StatusForbidden, "general.admin_required", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}