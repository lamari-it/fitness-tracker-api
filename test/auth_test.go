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
		testRegistration(t, e)
	})

	t.Run("Registration With Trainer Profile", func(t *testing.T) {
		CleanDatabase(t) // Clean before registration with trainer profile tests
		testRegistrationWithTrainerProfile(t, e)
	})

	t.Run("Login", func(t *testing.T) {
		CleanDatabase(t) // Clean before login tests
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
	t.Run("Successful Registration With Trainer Profile", func(t *testing.T) {
		userData := map[string]interface{}{
			"email":            "traineruser@example.com",
			"password":         "TrainerPassword123!",
			"password_confirm": "TrainerPassword123!",
			"first_name":       "Trainer",
			"last_name":        "User",
			"trainer_profile": map[string]interface{}{
				"bio":         "Experienced fitness trainer with focus on strength training and nutrition.",
				"specialties": []string{"Strength Training", "Nutrition", "Weight Loss"},
				"hourly_rate": 75.00,
				"location":    "New York, NY",
				"visibility":  "public",
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

		// Verify trainer profile fields
		trainerProfile := data.Value("trainer_profile").Object()
		trainerProfile.Value("id").String().NotEmpty()
		trainerProfile.Value("user_id").String().NotEmpty()
		trainerProfile.Value("bio").String().Contains("Experienced fitness trainer")
		trainerProfile.Value("specialties").Array().Length().IsEqual(3)
		trainerProfile.Value("hourly_rate").Number().IsEqual(75.00)
		trainerProfile.Value("location").String().IsEqual("New York, NY")
		trainerProfile.Value("visibility").String().IsEqual("public")
	})

	t.Run("Registration With Trainer Profile Private Visibility", func(t *testing.T) {
		userData := map[string]interface{}{
			"email":            "privatetrainer@example.com",
			"password":         "PrivatePass123!",
			"password_confirm": "PrivatePass123!",
			"first_name":       "Private",
			"last_name":        "Trainer",
			"trainer_profile": map[string]interface{}{
				"bio":         "Private trainer with exclusive clientele.",
				"specialties": []string{"VIP Training", "Executive Coaching"},
				"hourly_rate": 200.00,
				"location":    "Beverly Hills, CA",
				"visibility":  "private",
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
		userData := map[string]interface{}{
			"email":            "defaultvis@example.com",
			"password":         "DefaultPass123!",
			"password_confirm": "DefaultPass123!",
			"first_name":       "Default",
			"last_name":        "Visibility",
			"trainer_profile": map[string]interface{}{
				"bio":         "Trainer profile with default visibility setting.",
				"specialties": []string{"General Fitness"},
				"hourly_rate": 50.00,
				"location":    "Chicago, IL",
				// No visibility specified - should default to "public"
			},
		}

		response := e.POST("/api/v1/auth/register").
			WithJSON(userData).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		// Verify default visibility is "public"
		trainerProfile := response.Value("data").Object().Value("trainer_profile").Object()
		trainerProfile.Value("visibility").String().IsEqual("public")
	})

	t.Run("Registration With Invalid Trainer Profile", func(t *testing.T) {
		userData := map[string]interface{}{
			"email":            "invalidprofile@example.com",
			"password":         "InvalidPass123!",
			"password_confirm": "InvalidPass123!",
			"first_name":       "Invalid",
			"last_name":        "Profile",
			"trainer_profile": map[string]interface{}{
				"bio":         "Short", // Too short, should be at least 10 characters
				"specialties": []string{"Strength Training"},
				"hourly_rate": 75.00,
				"location":    "New York, NY",
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

	t.Run("Registration With Missing Trainer Profile Fields", func(t *testing.T) {
		userData := map[string]interface{}{
			"email":            "missingfields@example.com",
			"password":         "MissingPass123!",
			"password_confirm": "MissingPass123!",
			"first_name":       "Missing",
			"last_name":        "Fields",
			"trainer_profile": map[string]interface{}{
				"bio": "This bio is long enough to pass validation.",
				// Missing specialties, hourly_rate, and location
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
