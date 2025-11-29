package test

import (
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestAuthEndpoints(t *testing.T) {
	// Setup test app
	e := SetupTestApp(t)

	t.Run("Registration", func(t *testing.T) {
		CleanDatabase(t) // Clean before registration tests
		SeedTestRoles(t) // Seed roles for role assignment
		testRegistration(t, e)
	})

	t.Run("Registration With Trainer Profile", func(t *testing.T) {
		CleanDatabase(t)       // Clean before registration with trainer profile tests
		SeedTestRoles(t)       // Seed roles for role assignment
		SeedTestSpecialties(t) // Seed specialties for trainer profile
		testRegistrationWithTrainerProfile(t, e)
	})

	t.Run("Login", func(t *testing.T) {
		CleanDatabase(t) // Clean before login tests
		SeedTestRoles(t) // Seed roles for role assignment
		testLogin(t, e)
	})
}

func testRegistration(t *testing.T, e *httpexpect.Expect) {
	// Test successful registration
	t.Run("Successful Registration", func(t *testing.T) {
		userData := map[string]interface{}{
			"email":            "newuser@example.com",
			"password":         "SecurePassword123!",
			"password_confirm": "SecurePassword123!",
			"first_name":       "John",
			"last_name":        "Doe",
		}

		response := e.POST("/api/v1/auth/register").
			WithJSON(userData).
			Expect().
			Status(201).
			JSON().
			Object()

		// Check response structure
		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().NotEmpty()

		// Check user data
		data := response.Value("data").Object()
		data.Value("user").Object().NotEmpty()
		data.Value("token").String().NotEmpty()

		// Verify user fields
		user := data.Value("user").Object()
		user.Value("id").String().NotEmpty()
		user.Value("email").String().IsEqual("newuser@example.com")
		user.Value("first_name").String().IsEqual("John")
		user.Value("last_name").String().IsEqual("Doe")
		user.NotContainsKey("password") // Password should not be returned

		// Verify user role is assigned
		roles := user.Value("roles").Array()
		roles.Length().IsEqual(1)
		roles.Value(0).Object().Value("name").String().IsEqual("user")
	})

	// Test registration with existing email
	t.Run("Duplicate Email Registration", func(t *testing.T) {
		// First registration
		userData := map[string]interface{}{
			"email":            "duplicate@example.com",
			"password":         "Password123!",
			"password_confirm": "Password123!",
			"first_name":       "Jane",
			"last_name":        "Smith",
		}

		e.POST("/api/v1/auth/register").
			WithJSON(userData).
			Expect().
			Status(201)

		// Attempt duplicate registration
		response := e.POST("/api/v1/auth/register").
			WithJSON(userData).
			Expect().
			Status(409).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("already exists")
	})

	// Test registration with invalid data
	t.Run("Invalid Registration Data", func(t *testing.T) {
		testCases := []struct {
			name     string
			userData map[string]interface{}
			field    string
		}{
			{
				name: "Missing Email",
				userData: map[string]interface{}{
					"password":   "Password123!",
					"first_name": "John",
					"last_name":  "Doe",
				},
				field: "email",
			},
			{
				name: "Invalid Email Format",
				userData: map[string]interface{}{
					"email":            "invalid-email",
					"password":         "Password123!",
					"password_confirm": "Password123!",
					"first_name":       "John",
					"last_name":        "Doe",
				},
				field: "email",
			},
			{
				name: "Missing Password",
				userData: map[string]interface{}{
					"email":      "test@example.com",
					"first_name": "John",
					"last_name":  "Doe",
				},
				field: "password",
			},
			{
				name: "Short Password",
				userData: map[string]interface{}{
					"email":            "test@example.com",
					"password":         "123",
					"password_confirm": "123",
					"first_name":       "John",
					"last_name":        "Doe",
				},
				field: "password",
			},
			{
				name: "Missing First Name",
				userData: map[string]interface{}{
					"email":            "test@example.com",
					"password":         "Password123!",
					"password_confirm": "Password123!",
					"last_name":        "Doe",
				},
				field: "first_name",
			},
			{
				name: "Missing Last Name",
				userData: map[string]interface{}{
					"email":            "test@example.com",
					"password":         "Password123!",
					"password_confirm": "Password123!",
					"first_name":       "John",
				},
				field: "last_name",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				response := e.POST("/api/v1/auth/register").
					WithJSON(tc.userData).
					Expect().
					Status(400).
					JSON().
					Object()

				response.Value("success").Boolean().IsFalse()
				// Check that error is related to the expected field
				if tc.field != "" {
					errors := response.Value("errors").Object()
					errors.ContainsKey(tc.field)
				}
			})
		}
	})
}

func testRegistrationWithTrainerProfile(t *testing.T, e *httpexpect.Expect) {
	// Get specialty IDs for trainer profile
	specialtyIDs := GetSpecialtyIDs(t, "Strength Training", "Weight Loss", "Cardio")

	t.Run("Successful Registration With Trainer Profile", func(t *testing.T) {
		userData := map[string]interface{}{
			"email":            "traineruser@example.com",
			"password":         "TrainerPassword123!",
			"password_confirm": "TrainerPassword123!",
			"first_name":       "Trainer",
			"last_name":        "User",
			"trainer_profile": map[string]interface{}{
				"bio":           "Experienced fitness trainer with focus on strength training and nutrition.",
				"specialty_ids": specialtyIDs,
				"hourly_rate":   75.00,
				"location": map[string]interface{}{
					"city":         "New York",
					"region":       "NY",
					"country_code": "US",
				},
				"visibility": "public",
			},
		}

		response := e.POST("/api/v1/auth/register").
			WithJSON(userData).
			Expect().
			Status(201).
			JSON().
			Object()

		// Check response structure
		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().NotEmpty()

		// Check user data
		data := response.Value("data").Object()
		data.Value("user").Object().NotEmpty()
		data.Value("token").String().NotEmpty()
		data.Value("trainer_profile").Object().NotEmpty()

		// Verify user fields
		user := data.Value("user").Object()
		user.Value("id").String().NotEmpty()
		user.Value("email").String().IsEqual("traineruser@example.com")
		user.Value("first_name").String().IsEqual("Trainer")
		user.Value("last_name").String().IsEqual("User")
		user.NotContainsKey("password")

		// Verify user and trainer roles are assigned
		roles := user.Value("roles").Array()
		roles.Length().IsEqual(2)
		// Check that both "user" and "trainer" roles are present
		roleNames := []string{}
		for i := 0; i < 2; i++ {
			roleNames = append(roleNames, roles.Value(i).Object().Value("name").String().Raw())
		}
		// Verify both roles are present (order may vary)
		hasUser := false
		hasTrainer := false
		for _, name := range roleNames {
			if name == "user" {
				hasUser = true
			}
			if name == "trainer" {
				hasTrainer = true
			}
		}
		if !hasUser || !hasTrainer {
			t.Errorf("Expected both 'user' and 'trainer' roles, got: %v", roleNames)
		}

		// Verify trainer profile fields
		trainerProfile := data.Value("trainer_profile").Object()
		trainerProfile.Value("id").String().NotEmpty()
		trainerProfile.Value("user_id").String().NotEmpty()
		trainerProfile.Value("bio").String().Contains("Experienced fitness trainer")
		trainerProfile.Value("specialties").Array().Length().IsEqual(3)
		trainerProfile.Value("hourly_rate").Number().IsEqual(75.00)
		trainerProfile.Value("location").Object().Value("city").String().IsEqual("New York")
		trainerProfile.Value("visibility").String().IsEqual("public")
	})

	t.Run("Registration With Trainer Profile Private Visibility", func(t *testing.T) {
		// Get specialty IDs
		privateSpecialtyIDs := GetSpecialtyIDs(t, "Functional Fitness", "HIIT")

		userData := map[string]interface{}{
			"email":            "privatetrainer@example.com",
			"password":         "PrivatePass123!",
			"password_confirm": "PrivatePass123!",
			"first_name":       "Private",
			"last_name":        "Trainer",
			"trainer_profile": map[string]interface{}{
				"bio":           "Private trainer with exclusive clientele.",
				"specialty_ids": privateSpecialtyIDs,
				"hourly_rate":   200.00,
				"location": map[string]interface{}{
					"city":         "Beverly Hills",
					"region":       "CA",
					"country_code": "US",
				},
				"visibility": "private",
			},
		}

		response := e.POST("/api/v1/auth/register").
			WithJSON(userData).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		// Verify visibility is set correctly
		trainerProfile := response.Value("data").Object().Value("trainer_profile").Object()
		trainerProfile.Value("visibility").String().IsEqual("private")
	})

	t.Run("Registration With Trainer Profile Default Visibility", func(t *testing.T) {
		// Get specialty IDs
		defaultSpecialtyIDs := GetSpecialtyIDs(t, "Yoga")

		userData := map[string]interface{}{
			"email":            "defaultvis@example.com",
			"password":         "DefaultPass123!",
			"password_confirm": "DefaultPass123!",
			"first_name":       "Default",
			"last_name":        "Visibility",
			"trainer_profile": map[string]interface{}{
				"bio":           "Trainer profile with default visibility setting.",
				"specialty_ids": defaultSpecialtyIDs,
				"hourly_rate":   50.00,
				"location": map[string]interface{}{
					"city":         "Chicago",
					"region":       "IL",
					"country_code": "US",
				},
				// No visibility specified - should default to "private"
			},
		}

		response := e.POST("/api/v1/auth/register").
			WithJSON(userData).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		// Verify default visibility is "private"
		trainerProfile := response.Value("data").Object().Value("trainer_profile").Object()
		trainerProfile.Value("visibility").String().IsEqual("private")
	})

	t.Run("Registration With Minimal Trainer Profile", func(t *testing.T) {
		// Test that empty/zero values are now allowed
		userData := map[string]interface{}{
			"email":            "minimaltrainer@example.com",
			"password":         "MinimalPass123!",
			"password_confirm": "MinimalPass123!",
			"first_name":       "Minimal",
			"last_name":        "Trainer",
			"trainer_profile": map[string]interface{}{
				"bio":           "",         // Empty bio allowed
				"specialty_ids": []string{}, // Empty array allowed
				"hourly_rate":   0,          // Zero allowed
				// Empty/nil location allowed
				// No visibility - defaults to "private"
			},
		}

		response := e.POST("/api/v1/auth/register").
			WithJSON(userData).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		// Verify trainer profile with defaults
		trainerProfile := response.Value("data").Object().Value("trainer_profile").Object()
		trainerProfile.Value("id").String().NotEmpty()
		trainerProfile.Value("bio").String().IsEqual("")
		trainerProfile.Value("specialties").Array().Length().IsEqual(0)
		trainerProfile.Value("hourly_rate").Number().IsEqual(0)
		// Location should be nil/not present when empty (omitempty)
		trainerProfile.NotContainsKey("location")
		trainerProfile.Value("visibility").String().IsEqual("private")
	})

	t.Run("Registration With Invalid Trainer Profile", func(t *testing.T) {
		// Test that bio over 1000 characters fails
		longBio := ""
		for i := 0; i < 1001; i++ {
			longBio += "a"
		}

		userData := map[string]interface{}{
			"email":            "invalidprofile@example.com",
			"password":         "InvalidPass123!",
			"password_confirm": "InvalidPass123!",
			"first_name":       "Invalid",
			"last_name":        "Profile",
			"trainer_profile": map[string]interface{}{
				"bio": longBio, // Over 1000 characters
			},
		}

		response := e.POST("/api/v1/auth/register").
			WithJSON(userData).
			Expect().
			Status(400).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("Registration Without Trainer Profile Still Works", func(t *testing.T) {
		userData := map[string]interface{}{
			"email":            "regularuser@example.com",
			"password":         "RegularPass123!",
			"password_confirm": "RegularPass123!",
			"first_name":       "Regular",
			"last_name":        "User",
			// No trainer_profile field
		}

		response := e.POST("/api/v1/auth/register").
			WithJSON(userData).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		// Verify no trainer_profile in response
		data := response.Value("data").Object()
		data.Value("user").Object().NotEmpty()
		data.Value("token").String().NotEmpty()
		data.NotContainsKey("trainer_profile")
	})
}

func testLogin(t *testing.T, e *httpexpect.Expect) {
	// Setup: Create a user for login tests
	setupUserData := map[string]interface{}{
		"email":            "logintest@example.com",
		"password":         "LoginPassword123!",
		"password_confirm": "LoginPassword123!",
		"first_name":       "Login",
		"last_name":        "Test",
	}

	e.POST("/api/v1/auth/register").
		WithJSON(setupUserData).
		Expect().
		Status(201)

	// Test successful login
	t.Run("Successful Login", func(t *testing.T) {
		loginData := map[string]interface{}{
			"email":    "logintest@example.com",
			"password": "LoginPassword123!",
		}

		response := e.POST("/api/v1/auth/login").
			WithJSON(loginData).
			Expect().
			Status(200).
			JSON().
			Object()

		// Check response structure
		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("success")

		// Check data
		data := response.Value("data").Object()
		data.Value("token").String().NotEmpty()
		data.Value("user").Object().NotEmpty()

		// Verify user fields
		user := data.Value("user").Object()
		user.Value("id").String().NotEmpty()
		user.Value("email").String().IsEqual("logintest@example.com")
		user.Value("first_name").String().IsEqual("Login")
		user.Value("last_name").String().IsEqual("Test")
		user.NotContainsKey("password")

		// Verify roles are returned
		roles := user.Value("roles").Array()
		roles.Length().IsEqual(1)
		roles.Value(0).Object().Value("name").String().IsEqual("user")
	})

	// Test login with wrong password
	t.Run("Wrong Password", func(t *testing.T) {
		loginData := map[string]interface{}{
			"email":    "logintest@example.com",
			"password": "WrongPassword123!",
		}

		response := e.POST("/api/v1/auth/login").
			WithJSON(loginData).
			Expect().
			Status(401).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("Invalid")
	})

	// Test login with non-existent user
	t.Run("Non-existent User", func(t *testing.T) {
		loginData := map[string]interface{}{
			"email":    "nonexistent@example.com",
			"password": "Password123!",
		}

		response := e.POST("/api/v1/auth/login").
			WithJSON(loginData).
			Expect().
			Status(401).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("Invalid")
	})

	// Test login with invalid data
	t.Run("Invalid Login Data", func(t *testing.T) {
		testCases := []struct {
			name      string
			loginData map[string]interface{}
		}{
			{
				name: "Missing Email",
				loginData: map[string]interface{}{
					"password": "Password123!",
				},
			},
			{
				name: "Missing Password",
				loginData: map[string]interface{}{
					"email": "test@example.com",
				},
			},
			{
				name: "Invalid Email Format",
				loginData: map[string]interface{}{
					"email":    "invalid-email",
					"password": "Password123!",
				},
			},
			{
				name:      "Empty Body",
				loginData: map[string]interface{}{},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				response := e.POST("/api/v1/auth/login").
					WithJSON(tc.loginData).
					Expect().
					Status(400).
					JSON().
					Object()

				response.Value("success").Boolean().IsFalse()
			})
		}
	})
}

func TestRefreshTokenEndpoints(t *testing.T) {
	e := SetupTestApp(t)
	CleanDatabase(t)
	SeedTestRoles(t)

	t.Run("Refresh Token Flow", func(t *testing.T) {
		testRefreshTokenFlow(t, e)
	})
}

func testRefreshTokenFlow(t *testing.T, e *httpexpect.Expect) {
	// First create and login a user
	userData := map[string]interface{}{
		"email":            "refreshtest@example.com",
		"password":         "RefreshPassword123!",
		"password_confirm": "RefreshPassword123!",
		"first_name":       "Refresh",
		"last_name":        "Test",
	}

	registerResponse := e.POST("/api/v1/auth/register").
		WithJSON(userData).
		Expect().
		Status(201).
		JSON().
		Object()

	// Check that we get both tokens
	data := registerResponse.Value("data").Object()
	data.Value("access_token").String().NotEmpty()
	data.Value("refresh_token").String().NotEmpty()
	data.Value("expires_in").Number().Gt(0)
	data.Value("token_type").String().IsEqual("Bearer")
	data.Value("token").String().NotEmpty() // Backward compatibility

	refreshToken := data.Value("refresh_token").String().Raw()
	_ = data.Value("access_token").String().Raw() // Used for verification only

	t.Run("Login Returns Both Tokens", func(t *testing.T) {
		loginData := map[string]interface{}{
			"email":       "refreshtest@example.com",
			"password":    "RefreshPassword123!",
			"device_info": "Test Device",
		}

		loginResponse := e.POST("/api/v1/auth/login").
			WithJSON(loginData).
			Expect().
			Status(200).
			JSON().
			Object()

		loginData2 := loginResponse.Value("data").Object()
		loginData2.Value("access_token").String().NotEmpty()
		loginData2.Value("refresh_token").String().NotEmpty()
		loginData2.Value("expires_in").Number().Gt(0)
		loginData2.Value("token_type").String().IsEqual("Bearer")
	})

	t.Run("Successful Token Refresh", func(t *testing.T) {
		refreshData := map[string]interface{}{
			"refresh_token": refreshToken,
		}

		response := e.POST("/api/v1/auth/refresh").
			WithJSON(refreshData).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("refresh")

		// Check new tokens are returned
		newData := response.Value("data").Object()
		newData.Value("access_token").String().NotEmpty()
		newData.Value("refresh_token").String().NotEmpty()
		newData.Value("expires_in").Number().Gt(0)

		// Old token should no longer work (rotation)
		oldRefreshData := map[string]interface{}{
			"refresh_token": refreshToken,
		}

		e.POST("/api/v1/auth/refresh").
			WithJSON(oldRefreshData).
			Expect().
			Status(401)
	})

	t.Run("Invalid Refresh Token", func(t *testing.T) {
		refreshData := map[string]interface{}{
			"refresh_token": "invalid_token_here",
		}

		response := e.POST("/api/v1/auth/refresh").
			WithJSON(refreshData).
			Expect().
			Status(401).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("Logout Single Session", func(t *testing.T) {
		// Login to get a new refresh token
		loginData := map[string]interface{}{
			"email":    "refreshtest@example.com",
			"password": "RefreshPassword123!",
		}

		loginResponse := e.POST("/api/v1/auth/login").
			WithJSON(loginData).
			Expect().
			Status(200).
			JSON().
			Object()

		newRefreshToken := loginResponse.Value("data").Object().Value("refresh_token").String().Raw()

		// Logout using the refresh token
		logoutData := map[string]interface{}{
			"refresh_token": newRefreshToken,
		}

		logoutResponse := e.POST("/api/v1/auth/logout").
			WithJSON(logoutData).
			Expect().
			Status(200).
			JSON().
			Object()

		logoutResponse.Value("success").Boolean().IsTrue()

		// Verify token no longer works
		refreshData := map[string]interface{}{
			"refresh_token": newRefreshToken,
		}

		e.POST("/api/v1/auth/refresh").
			WithJSON(refreshData).
			Expect().
			Status(401)
	})

	t.Run("Logout All Sessions", func(t *testing.T) {
		// Login multiple times to create multiple sessions
		loginData := map[string]interface{}{
			"email":       "refreshtest@example.com",
			"password":    "RefreshPassword123!",
			"device_info": "Device 1",
		}

		loginResponse1 := e.POST("/api/v1/auth/login").
			WithJSON(loginData).
			Expect().
			Status(200).
			JSON().
			Object()

		accessToken1 := loginResponse1.Value("data").Object().Value("access_token").String().Raw()
		refreshToken1 := loginResponse1.Value("data").Object().Value("refresh_token").String().Raw()

		loginData["device_info"] = "Device 2"
		loginResponse2 := e.POST("/api/v1/auth/login").
			WithJSON(loginData).
			Expect().
			Status(200).
			JSON().
			Object()

		refreshToken2 := loginResponse2.Value("data").Object().Value("refresh_token").String().Raw()

		// Logout all using the access token
		logoutAllResponse := e.POST("/api/v1/auth/logout-all").
			WithHeader("Authorization", "Bearer "+accessToken1).
			Expect().
			Status(200).
			JSON().
			Object()

		logoutAllResponse.Value("success").Boolean().IsTrue()
		logoutAllResponse.Value("data").Object().Value("sessions_revoked").Number().Ge(2)

		// Verify both refresh tokens no longer work
		e.POST("/api/v1/auth/refresh").
			WithJSON(map[string]interface{}{"refresh_token": refreshToken1}).
			Expect().
			Status(401)

		e.POST("/api/v1/auth/refresh").
			WithJSON(map[string]interface{}{"refresh_token": refreshToken2}).
			Expect().
			Status(401)
	})

	t.Run("Get Sessions", func(t *testing.T) {
		// Login to create a session
		loginData := map[string]interface{}{
			"email":       "refreshtest@example.com",
			"password":    "RefreshPassword123!",
			"device_info": "Test Sessions Device",
		}

		loginResponse := e.POST("/api/v1/auth/login").
			WithJSON(loginData).
			Expect().
			Status(200).
			JSON().
			Object()

		accessToken := loginResponse.Value("data").Object().Value("access_token").String().Raw()

		// Get sessions
		sessionsResponse := e.GET("/api/v1/auth/sessions").
			WithHeader("Authorization", "Bearer "+accessToken).
			Expect().
			Status(200).
			JSON().
			Object()

		sessionsResponse.Value("success").Boolean().IsTrue()
		sessionsResponse.Value("data").Array().NotEmpty()

		// Check session structure
		session := sessionsResponse.Value("data").Array().Value(0).Object()
		session.Value("id").String().NotEmpty()
		session.Value("device_info").String().NotEmpty()
		session.Value("created_at").String().NotEmpty()
	})

	t.Run("Revoke Specific Session", func(t *testing.T) {
		// Login twice
		loginData := map[string]interface{}{
			"email":       "refreshtest@example.com",
			"password":    "RefreshPassword123!",
			"device_info": "Main Device",
		}

		loginResponse1 := e.POST("/api/v1/auth/login").
			WithJSON(loginData).
			Expect().
			Status(200).
			JSON().
			Object()

		accessToken := loginResponse1.Value("data").Object().Value("access_token").String().Raw()

		loginData["device_info"] = "Other Device"
		e.POST("/api/v1/auth/login").
			WithJSON(loginData).
			Expect().
			Status(200)

		// Get sessions
		sessionsResponse := e.GET("/api/v1/auth/sessions").
			WithHeader("Authorization", "Bearer "+accessToken).
			Expect().
			Status(200).
			JSON().
			Object()

		// Get the second session ID
		sessions := sessionsResponse.Value("data").Array()
		if sessions.Length().Raw() > 1 {
			sessionID := sessions.Value(1).Object().Value("id").String().Raw()

			// Revoke that session
			revokeResponse := e.DELETE("/api/v1/auth/sessions/"+sessionID).
				WithHeader("Authorization", "Bearer "+accessToken).
				Expect().
				Status(200).
				JSON().
				Object()

			revokeResponse.Value("success").Boolean().IsTrue()
		}
	})

	t.Run("Access Token Still Works After Refresh", func(t *testing.T) {
		// Login to get fresh tokens
		loginData := map[string]interface{}{
			"email":    "refreshtest@example.com",
			"password": "RefreshPassword123!",
		}

		loginResponse := e.POST("/api/v1/auth/login").
			WithJSON(loginData).
			Expect().
			Status(200).
			JSON().
			Object()

		newAccessToken := loginResponse.Value("data").Object().Value("access_token").String().Raw()

		// Access token should still work (JWT is stateless)
		e.GET("/api/v1/auth/profile").
			WithHeader("Authorization", "Bearer "+newAccessToken).
			Expect().
			Status(200)
	})
}

// Test protected endpoint with token
func TestProtectedEndpoint(t *testing.T) {
	e := SetupTestApp(t)
	CleanDatabase(t) // Clean before test

	// Create a user and get token
	userData := map[string]interface{}{
		"email":            "protected@example.com",
		"password":         "Password123!",
		"password_confirm": "Password123!",
		"first_name":       "Protected",
		"last_name":        "User",
	}

	// Register user
	e.POST("/api/v1/auth/register").
		WithJSON(userData).
		Expect().
		Status(201)

	// Get token
	token := GetAuthToken(e, "protected@example.com", "Password123!")

	t.Run("Access Protected Route With Token", func(t *testing.T) {
		response := e.GET("/api/v1/auth/profile").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("email").String().IsEqual("protected@example.com")
	})

	t.Run("Access Protected Route Without Token", func(t *testing.T) {
		response := e.GET("/api/v1/auth/profile").
			Expect().
			Status(401).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("Authorization")
	})

	t.Run("Access Protected Route With Invalid Token", func(t *testing.T) {
		response := e.GET("/api/v1/auth/profile").
			WithHeader("Authorization", "Bearer invalid_token_here").
			Expect().
			Status(401).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})
}
