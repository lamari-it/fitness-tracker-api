package utils

import (
	"fit-flow-api/models"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetAuthUserID extracts the authenticated user ID from Gin context.
// Returns the user ID and true if successful, or uuid.Nil and false if not authenticated.
// Automatically sends appropriate error responses when authentication fails.
func GetAuthUserID(c *gin.Context) (uuid.UUID, bool) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return uuid.Nil, false
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		InternalServerErrorResponse(c, "Invalid user ID type")
		return uuid.Nil, false
	}

	return userID, true
}

// ParseUUID parses a UUID string with consistent error handling.
// Returns the UUID and true if successful, or uuid.Nil and false if invalid.
// Automatically sends a BadRequest response with a descriptive error message.
func ParseUUID(c *gin.Context, idStr, resourceName string) (uuid.UUID, bool) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		BadRequestResponse(c, fmt.Sprintf("Invalid %s ID format", resourceName), nil)
		return uuid.Nil, false
	}
	return id, true
}

// ParseUUIDParam parses a UUID from a URL parameter.
// Returns the UUID and true if successful, or uuid.Nil and false if invalid.
// Automatically sends a BadRequest response with a descriptive error message.
func ParseUUIDParam(c *gin.Context, paramName, resourceName string) (uuid.UUID, bool) {
	idStr := c.Param(paramName)
	return ParseUUID(c, idStr, resourceName)
}

// RequireAdmin checks if the current user is an admin.
// Returns true if the user is an admin, false otherwise.
// Automatically sends a ForbiddenResponse if the user is not an admin.
func RequireAdmin(c *gin.Context) bool {
	user, exists := c.Get("user")
	if !exists {
		ForbiddenResponse(c, "Admin access required.")
		return false
	}
	if !user.(models.User).IsAdmin {
		ForbiddenResponse(c, "Admin access required.")
		return false
	}
	return true
}

// GetAuthUser extracts the full authenticated user object from Gin context.
// Returns the user and true if successful, or nil and false if not authenticated.
// Automatically sends an UnauthorizedResponse when authentication fails.
func GetAuthUser(c *gin.Context) (*models.User, bool) {
	userVal, exists := c.Get("user")
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return nil, false
	}
	user, ok := userVal.(models.User)
	if !ok {
		InternalServerErrorResponse(c, "Invalid user type")
		return nil, false
	}
	return &user, true
}
