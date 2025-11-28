package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode"

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

	// Handle JSON type mismatch errors (e.g., string instead of object)
	if unmarshalErr, ok := err.(*json.UnmarshalTypeError); ok {
		field := toSnakeCase(unmarshalErr.Field)
		// Extract the last part of the field path for nested fields
		if strings.Contains(field, ".") {
			parts := strings.Split(field, ".")
			field = parts[len(parts)-1]
		}

		var message string
		switch unmarshalErr.Type.Kind().String() {
		case "struct", "ptr":
			// Handle location field specifically
			if strings.Contains(strings.ToLower(unmarshalErr.Field), "location") {
				message = "Location must be an object with fields like city, region, country_code, latitude, longitude. Use null or omit the field if no location is needed."
			} else {
				message = fmt.Sprintf("Expected an object, but received %s. Please provide a valid object or omit this field.", unmarshalErr.Value)
			}
		case "slice":
			message = fmt.Sprintf("Expected an array, but received %s.", unmarshalErr.Value)
		case "int", "int64", "float64":
			message = fmt.Sprintf("Expected a number, but received %s.", unmarshalErr.Value)
		case "bool":
			message = fmt.Sprintf("Expected a boolean (true/false), but received %s.", unmarshalErr.Value)
		default:
			message = fmt.Sprintf("Invalid type: expected %s, but received %s.", unmarshalErr.Type.String(), unmarshalErr.Value)
		}

		validationErrors[field] = []string{message}
		ValidationErrorResponse(c, validationErrors)
		return
	}

	// Handle JSON syntax errors
	if syntaxErr, ok := err.(*json.SyntaxError); ok {
		validationErrors["json"] = []string{fmt.Sprintf("Invalid JSON syntax at position %d.", syntaxErr.Offset)}
		ValidationErrorResponse(c, validationErrors)
		return
	}

	// Parse field validation errors from binding
	if bindingErr, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range bindingErr {
			// Convert struct field name (e.g., "LastName") to snake_case (e.g., "last_name")
			field := toSnakeCase(fieldErr.Field())
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

			validationErrors[field] = []string{message}
		}
	}

	// If no specific errors were captured, provide a generic message
	if len(validationErrors) == 0 {
		validationErrors["request"] = []string{fmt.Sprintf("Invalid request format: %s", err.Error())}
	}

	ValidationErrorResponse(c, validationErrors)
}

// toSnakeCase converts CamelCase to snake_case
func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && unicode.IsUpper(r) {
			// Check if previous character was lowercase or if next is lowercase
			// to properly handle cases like "LastName" -> "last_name"
			if (i > 0 && unicode.IsLower(rune(str[i-1]))) ||
				(i < len(str)-1 && unicode.IsLower(rune(str[i+1]))) {
				result = append(result, '_')
			}
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}
