package test

import (
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestUserSettingsEndpoints(t *testing.T) {
	e := SetupTestApp(t)

	t.Run("User Settings", func(t *testing.T) {
		CleanDatabase(t)
		SeedTestRoles(t)
		testUserSettings(t, e)
	})
}

func testUserSettings(t *testing.T, e *httpexpect.Expect) {
	// Setup: Create a user
	userData := map[string]interface{}{
		"email":            "settings@example.com",
		"password":         "SettingsPass123!",
		"password_confirm": "SettingsPass123!",
		"first_name":       "Settings",
		"last_name":        "User",
	}

	e.POST("/api/v1/auth/register").
		WithJSON(userData).
		Expect().
		Status(201)

	token := GetAuthToken(e, "settings@example.com", "SettingsPass123!")

	t.Run("Get User Settings", func(t *testing.T) {
		response := e.GET("/api/v1/user/settings").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("retrieved")

		data := response.Value("data").Object()
		data.Value("email").String().IsEqual("settings@example.com")
		data.Value("first_name").String().IsEqual("Settings")
		data.Value("last_name").String().IsEqual("User")
		// Check default unit preferences
		data.Value("preferred_weight_unit").String().IsEqual("kg")
		data.Value("preferred_height_unit").String().IsEqual("cm")
		data.Value("preferred_distance_unit").String().IsEqual("km")
	})

	t.Run("Update Preferred Weight Unit to lb", func(t *testing.T) {
		response := e.PUT("/api/v1/user/settings").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"preferred_weight_unit": "lb",
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("updated")

		data := response.Value("data").Object()
		data.Value("preferred_weight_unit").String().IsEqual("lb")
		// Other units should remain unchanged
		data.Value("preferred_height_unit").String().IsEqual("cm")
		data.Value("preferred_distance_unit").String().IsEqual("km")
	})

	t.Run("Update Preferred Height Unit to ft", func(t *testing.T) {
		response := e.PUT("/api/v1/user/settings").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"preferred_height_unit": "ft",
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Object()
		data.Value("preferred_height_unit").String().IsEqual("ft")
		// Previous update should persist
		data.Value("preferred_weight_unit").String().IsEqual("lb")
	})

	t.Run("Update Preferred Distance Unit to mi", func(t *testing.T) {
		response := e.PUT("/api/v1/user/settings").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"preferred_distance_unit": "mi",
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Object()
		data.Value("preferred_distance_unit").String().IsEqual("mi")
	})

	t.Run("Update Multiple Settings At Once", func(t *testing.T) {
		response := e.PUT("/api/v1/user/settings").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"preferred_weight_unit":   "kg",
				"preferred_height_unit":   "cm",
				"preferred_distance_unit": "km",
				"first_name":              "Updated",
				"last_name":               "Name",
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Object()
		data.Value("preferred_weight_unit").String().IsEqual("kg")
		data.Value("preferred_height_unit").String().IsEqual("cm")
		data.Value("preferred_distance_unit").String().IsEqual("km")
		data.Value("first_name").String().IsEqual("Updated")
		data.Value("last_name").String().IsEqual("Name")
	})

	t.Run("Update With Invalid Weight Unit", func(t *testing.T) {
		response := e.PUT("/api/v1/user/settings").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"preferred_weight_unit": "stones",
			}).
			Expect().
			Status(400).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("Update With Invalid Height Unit", func(t *testing.T) {
		response := e.PUT("/api/v1/user/settings").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"preferred_height_unit": "inches",
			}).
			Expect().
			Status(400).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("Update With Invalid Distance Unit", func(t *testing.T) {
		response := e.PUT("/api/v1/user/settings").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"preferred_distance_unit": "yards",
			}).
			Expect().
			Status(400).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("Update With Empty Request Body", func(t *testing.T) {
		// Empty body should succeed but not change anything
		response := e.PUT("/api/v1/user/settings").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
	})

	t.Run("Get Settings Without Auth", func(t *testing.T) {
		e.GET("/api/v1/user/settings").
			Expect().
			Status(401)
	})

	t.Run("Update Settings Without Auth", func(t *testing.T) {
		e.PUT("/api/v1/user/settings").
			WithJSON(map[string]interface{}{
				"preferred_weight_unit": "lb",
			}).
			Expect().
			Status(401)
	})
}
