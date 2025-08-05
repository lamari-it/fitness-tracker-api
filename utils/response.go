package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// StandardResponse represents the standard API response structure
type StandardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta represents pagination metadata
type Meta struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	TotalPages  int `json:"total_pages"`
	TotalItems  int `json:"total_items"`
}

// ValidationErrors represents field-specific validation errors
type ValidationErrors map[string][]string

// SuccessResponse sends a successful response
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
		Errors:  nil,
	})
}

// CreatedResponse sends a successful creation response
func CreatedResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
		Errors:  nil,
	})
}

// PaginatedResponse sends a successful response with pagination metadata
func PaginatedResponse(c *gin.Context, message string, data interface{}, currentPage, perPage, totalItems int) {
	totalPages := (totalItems + perPage - 1) / perPage // Ceiling division
	
	c.JSON(http.StatusOK, StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
		Errors:  nil,
		Meta: &Meta{
			CurrentPage: currentPage,
			PerPage:     perPage,
			TotalPages:  totalPages,
			TotalItems:  totalItems,
		},
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, message string, errors interface{}) {
	c.JSON(statusCode, StandardResponse{
		Success: false,
		Message: message,
		Data:    nil,
		Errors:  errors,
	})
}

// BadRequestResponse sends a bad request error response
func BadRequestResponse(c *gin.Context, message string, errors interface{}) {
	ErrorResponse(c, http.StatusBadRequest, message, errors)
}

// UnauthorizedResponse sends an unauthorized error response
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, message, nil)
}

// ForbiddenResponse sends a forbidden error response
func ForbiddenResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, message, nil)
}

// NotFoundResponse sends a not found error response
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message, nil)
}

// InternalServerErrorResponse sends an internal server error response
func InternalServerErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, message, nil)
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, errors ValidationErrors) {
	ErrorResponse(c, http.StatusBadRequest, "Validation failed.", errors)
}

// ConflictResponse sends a conflict error response
func ConflictResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusConflict, message, nil)
}

// NoContentResponse sends a no content response (for successful deletions)
func NoContentResponse(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

// DeletedResponse sends a successful deletion response
func DeletedResponse(c *gin.Context, message string) {
	c.JSON(http.StatusOK, StandardResponse{
		Success: true,
		Message: message,
		Data:    nil,
		Errors:  nil,
	})
}


// HandleBindingError processes Gin binding errors and sends appropriate response
func HandleBindingError(c *gin.Context, err error) {
	validationErrors := make(ValidationErrors)
	
	// Parse field validation errors from binding
	if bindingErr, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range bindingErr {
			field := fieldErr.Field()
			tag := fieldErr.Tag()
			
			var message string
			switch tag {
			case "required":
				message = "This field is required."
			case "email":
				message = "Please provide a valid email address."
			case "min":
				message = fmt.Sprintf("This field must be at least %s characters.", fieldErr.Param())
			case "max":
				message = fmt.Sprintf("This field must not exceed %s characters.", fieldErr.Param())
			default:
				message = fmt.Sprintf("This field failed %s validation.", tag)
			}
			
			validationErrors[strings.ToLower(field)] = []string{message}
		}
	}
	
	ValidationErrorResponse(c, validationErrors)
}