package middleware

import (
	"net/http"
	"strings"

	"lamari-fit-api/utils"

	"github.com/gin-gonic/gin"
)

// I18nMiddleware handles language detection and sets the language in the context
func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		i18n := utils.GetI18n()

		// Get language from multiple sources (priority order)
		lang := getLanguageFromSources(c, i18n)

		// Set language in context for use in controllers
		c.Set("language", lang)

		// Set language in header for response
		c.Header("Content-Language", lang)

		c.Next()
	}
}

// getLanguageFromSources detects language from various sources
func getLanguageFromSources(c *gin.Context, i18n *utils.I18n) string {
	// 1. Check query parameter 'lang'
	if lang := c.Query("lang"); lang != "" {
		if i18n.IsLanguageSupported(lang) {
			return lang
		}
	}

	// 2. Check custom header 'X-Language'
	if lang := c.GetHeader("X-Language"); lang != "" {
		if i18n.IsLanguageSupported(lang) {
			return lang
		}
	}

	// 3. Check Accept-Language header
	if acceptLang := c.GetHeader("Accept-Language"); acceptLang != "" {
		if lang := parseAcceptLanguage(acceptLang, i18n); lang != "" {
			return lang
		}
	}

	// 4. Check user preference from JWT token (if authenticated)
	if userID, exists := c.Get("userID"); exists {
		// For now, we'll skip user preference from database
		// This can be implemented later by fetching user's preferred language
		_ = userID
	}

	// 5. Default language
	return i18n.GetDefaultLanguage()
}

// parseAcceptLanguage parses Accept-Language header and returns the best match
func parseAcceptLanguage(acceptLang string, i18n *utils.I18n) string {
	languages := strings.Split(acceptLang, ",")

	for _, lang := range languages {
		// Clean up the language tag (remove quality values, etc.)
		lang = strings.TrimSpace(lang)
		if idx := strings.Index(lang, ";"); idx != -1 {
			lang = lang[:idx]
		}

		// Try exact match first
		if i18n.IsLanguageSupported(lang) {
			return lang
		}

		// Try language without region (e.g., "en-US" -> "en")
		if idx := strings.Index(lang, "-"); idx != -1 {
			baseLang := lang[:idx]
			if i18n.IsLanguageSupported(baseLang) {
				return baseLang
			}
		}
	}

	return ""
}

// GetLanguage is a helper function to get language from context
func GetLanguage(c *gin.Context) string {
	if lang, exists := c.Get("language"); exists {
		return lang.(string)
	}
	return utils.GetI18n().GetDefaultLanguage()
}

// Translate is a helper function to translate messages in controllers
func Translate(c *gin.Context, key string, args ...interface{}) string {
	lang := GetLanguage(c)
	return utils.GetI18n().T(lang, key, args...)
}

// TranslateResponse creates a standardized response with translation
func TranslateResponse(c *gin.Context, statusCode int, messageKey string, data interface{}, args ...interface{}) {
	message := Translate(c, messageKey, args...)

	response := gin.H{
		"status":  http.StatusText(statusCode),
		"message": message,
	}

	if data != nil {
		response["data"] = data
	}

	c.JSON(statusCode, response)
}

// TranslateErrorResponse creates a standardized error response with translation
func TranslateErrorResponse(c *gin.Context, statusCode int, messageKey string, errors interface{}, args ...interface{}) {
	// message := Translate(c, messageKey, args...)

	response := gin.H{
		"status":  http.StatusText(statusCode),
		"message": messageKey,
		"error":   true,
	}

	if errors != nil {
		response["errors"] = errors
	}

	c.JSON(statusCode, response)
}
