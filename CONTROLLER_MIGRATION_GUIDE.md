# Controller Migration Guide: Standard Response Format

This guide helps you migrate existing controllers to use the new standard response format.

## Overview

All API responses must follow the standard format:
```json
{
  "success": boolean,
  "message": string,
  "data": any,
  "errors": object | null,
  "meta": object | null  // Only for paginated responses
}
```

## Step-by-Step Migration

### 1. Import the Response Utils

Add this import to your controller:
```go
import "lamari-fit-api/utils"
```

### 2. Replace Direct JSON Responses

#### Before:
```go
c.JSON(http.StatusOK, user)
```

#### After:
```go
utils.SuccessResponse(c, "User fetched successfully.", user)
```

### 3. Handle Validation Errors

#### Before:
```go
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
}
```

#### After:
```go
if err := c.ShouldBindJSON(&req); err != nil {
    utils.HandleBindingError(c, err)
    return
}
```

### 4. Handle Custom Validation

#### Before:
```go
if req.Password != req.PasswordConfirm {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
    return
}
```

#### After:
```go
if req.Password != req.PasswordConfirm {
    validationErrors := utils.ValidationErrors{
        "password_confirm": []string{"Passwords do not match."},
    }
    utils.ValidationErrorResponse(c, validationErrors)
    return
}
```

### 5. Handle Paginated Responses

#### Before:
```go
c.JSON(http.StatusOK, gin.H{
    "plans": plans,
    "total": total,
    "page":  page,
    "limit": limit,
})
```

#### After:
```go
utils.PaginatedResponse(c, "Workout plans fetched successfully.", plans, page, limit, int(total))
```

### 6. Common Response Patterns

#### Success Responses:
```go
// GET request success
utils.SuccessResponse(c, "Resource fetched successfully.", resource)

// POST request success (creation)
utils.CreatedResponse(c, "Resource created successfully.", resource)

// DELETE request success
utils.DeletedResponse(c, "Resource deleted successfully.")

// With translation support
utils.SuccessResponseWithTranslation(c, "resource.fetched", resource)
```

#### Error Responses:
```go
// Bad Request (400)
utils.BadRequestResponse(c, "Invalid request format.", nil)

// Unauthorized (401)
utils.UnauthorizedResponse(c, "User not authenticated.")

// Forbidden (403)
utils.ForbiddenResponse(c, "You don't have permission to access this resource.")

// Not Found (404)
utils.NotFoundResponse(c, "Resource not found.")

// Conflict (409)
utils.ConflictResponse(c, "A resource with this name already exists.")

// Internal Server Error (500)
utils.InternalServerErrorResponse(c, "Failed to process request.")

// With translation support
utils.ErrorResponseWithTranslation(c, http.StatusNotFound, "resource.not_found", nil)
```

## Complete Example: Auth Controller

### Before:
```go
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

    // ... create user ...

    c.JSON(http.StatusCreated, gin.H{
        "user": user,
        "token": token,
    })
}
```

### After:
```go
func Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.HandleBindingError(c, err)
        return
    }

    if req.Password != req.PasswordConfirm {
        validationErrors := utils.ValidationErrors{
            "password_confirm": []string{"Passwords do not match."},
        }
        utils.ValidationErrorResponse(c, validationErrors)
        return
    }

    // ... create user ...

    authResponse := AuthResponse{
        User:  user.ToResponse(),
        Token: token,
    }

    utils.CreatedResponse(c, "User registered successfully.", authResponse)
}
```

## Controller Checklist

Use this checklist when migrating each controller:

- [ ] Import `lamari-fit-api/utils` package
- [ ] Replace all `c.JSON()` calls with appropriate utils functions
- [ ] Update validation error handling to use `HandleBindingError`
- [ ] Convert custom validations to use `ValidationErrors` type
- [ ] Update paginated endpoints to use `PaginatedResponse`
- [ ] Ensure all error messages are user-friendly
- [ ] Test all endpoints with Postman collection
- [ ] Update any frontend code that depends on the old response format

## Response Utils Reference

### Success Functions:
- `SuccessResponse(c, message, data)` - Standard success (200)
- `CreatedResponse(c, message, data)` - Creation success (201)
- `DeletedResponse(c, message)` - Deletion success (200)
- `PaginatedResponse(c, message, data, page, perPage, total)` - Paginated success (200)

### Error Functions:
- `ErrorResponse(c, statusCode, message, errors)` - Generic error
- `BadRequestResponse(c, message, errors)` - 400 error
- `UnauthorizedResponse(c, message)` - 401 error
- `ForbiddenResponse(c, message)` - 403 error
- `NotFoundResponse(c, message)` - 404 error
- `ConflictResponse(c, message)` - 409 error
- `InternalServerErrorResponse(c, message)` - 500 error
- `ValidationErrorResponse(c, errors)` - Validation error (400)
- `HandleBindingError(c, err)` - Auto-parse Gin binding errors

### With Translation Support:
- `SuccessResponseWithTranslation(c, messageKey, data, args...)`
- `CreatedResponseWithTranslation(c, messageKey, data, args...)`
- `ErrorResponseWithTranslation(c, statusCode, messageKey, errors, args...)`
- `BadRequestResponseWithTranslation(c, messageKey, errors, args...)`

## Testing

After migration, test each endpoint to ensure:
1. Response structure matches the standard format
2. HTTP status codes are correct
3. Error messages are properly formatted
4. Pagination metadata is included where applicable
5. All fields (success, message, data, errors) are present

Use the updated Postman collection to verify all endpoints work correctly with the new format.