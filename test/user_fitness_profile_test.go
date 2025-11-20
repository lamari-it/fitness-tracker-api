package test

import (
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestUserFitnessProfileEndpoints(t *testing.T) {
	e := SetupTestApp(t)

	t.Run("Create Fitness Profile", func(t *testing.T) {
		CleanDatabase(t)
		testCreateFitnessProfile(t, e)
	})

	t.Run("Get Fitness Profile", func(t *testing.T) {
		CleanDatabase(t)
		testGetFitnessProfile(t, e)
	})

	t.Run("Update Fitness Profile", func(t *testing.T) {
		CleanDatabase(t)
		testUpdateFitnessProfile(t, e)
	})

	t.Run("Delete Fitness Profile", func(t *testing.T) {
		CleanDatabase(t)
		testDeleteFitnessProfile(t, e)
	})

	t.Run("Log Weight", func(t *testing.T) {
		CleanDatabase(t)
		testLogWeight(t, e)
	})

	t.Run("Validation Errors", func(t *testing.T) {
		CleanDatabase(t)
		testFitnessProfileValidation(t, e)
	})
}

func testCreateFitnessProfile(t *testing.T, e *httpexpect.Expect) {
	token := createTestUserAndGetToken(e, "profile@example.com", "ProfilePass123!", "John", "Doe")

	t.Run("Successful Creation with Required Fields", func(t *testing.T) {
		profileData := map[string]interface{}{
			"date_of_birth":     "1990-05-15",
			"gender":            "male",
			"height_cm":         180.5,
			"current_weight_kg": 80.0,
			"primary_goal":      "muscle_gain",
		}

		response := e.POST("/api/v1/user/fitness-profile").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(profileData).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("created")

		data := response.Value("data").Object()
		data.Value("id").String().NotEmpty()
		data.Value("date_of_birth").String().IsEqual("1990-05-15")
		data.Value("age").Number().Gt(0)
		data.Value("gender").String().IsEqual("male")
		data.Value("height_cm").Number().IsEqual(180.5)
		data.Value("height_ft_in").String().NotEmpty()
		data.Value("current_weight_kg").Number().IsEqual(80.0)
		data.Value("current_weight_lbs").Number().Gt(0)
		data.Value("primary_goal").String().IsEqual("muscle_gain")
		data.Value("preferred_unit_system").String().IsEqual("metric")
		data.Value("target_weekly_workouts").Number().IsEqual(3)
		data.Value("activity_level").String().IsEqual("moderate")
		data.Value("training_locations").Array().Length().IsEqual(1)
		data.Value("preferred_workout_duration_mins").Number().IsEqual(45)
		data.Value("available_days").Array().Length().IsEqual(3)
	})

	t.Run("Duplicate Profile Creation", func(t *testing.T) {
		profileData := map[string]interface{}{
			"date_of_birth":     "1990-05-15",
			"gender":            "female",
			"height_cm":         165.0,
			"current_weight_kg": 60.0,
			"primary_goal":      "weight_loss",
		}

		response := e.POST("/api/v1/user/fitness-profile").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(profileData).
			Expect().
			Status(409).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("already exists")
	})

	t.Run("Successful Creation with All Fields", func(t *testing.T) {
		fullToken := createTestUserAndGetToken(e, "fullprofile@example.com", "FullPass123!", "Jane", "Smith")

		profileData := map[string]interface{}{
			"date_of_birth":                  "1985-08-20",
			"gender":                         "female",
			"height_cm":                      165.0,
			"current_weight_kg":              70.0,
			"primary_goal":                   "weight_loss",
			"preferred_unit_system":          "imperial",
			"target_weight_kg":               60.0,
			"target_weekly_workouts":         5,
			"activity_level":                 "active",
			"training_locations":             []string{"home", "gym"},
			"preferred_workout_duration_mins": 60,
			"available_days":                 []string{"monday", "tuesday", "thursday", "saturday"},
			"health_conditions":              "None",
			"injuries_notes":                 "Minor knee issue",
		}

		response := e.POST("/api/v1/user/fitness-profile").
			WithHeader("Authorization", "Bearer "+fullToken).
			WithJSON(profileData).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Object()
		data.Value("preferred_unit_system").String().IsEqual("imperial")
		data.Value("target_weight_kg").Number().IsEqual(60.0)
		data.Value("target_weight_lbs").Number().Gt(0)
		data.Value("target_weekly_workouts").Number().IsEqual(5)
		data.Value("activity_level").String().IsEqual("active")
		data.Value("training_locations").Array().Length().IsEqual(2)
		data.Value("preferred_workout_duration_mins").Number().IsEqual(60)
		data.Value("available_days").Array().Length().IsEqual(4)
		data.Value("health_conditions").String().IsEqual("None")
		data.Value("injuries_notes").String().IsEqual("Minor knee issue")
	})

	t.Run("Create Profile Without Auth", func(t *testing.T) {
		profileData := map[string]interface{}{
			"date_of_birth":     "1990-05-15",
			"gender":            "male",
			"height_cm":         180.0,
			"current_weight_kg": 80.0,
			"primary_goal":      "general_fitness",
		}

		response := e.POST("/api/v1/user/fitness-profile").
			WithJSON(profileData).
			Expect().
			Status(401).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})
}

func testGetFitnessProfile(t *testing.T, e *httpexpect.Expect) {
	token := createTestUserAndGetToken(e, "getprofile@example.com", "GetPass123!", "Get", "Profile")

	t.Run("Get Non-existent Profile", func(t *testing.T) {
		response := e.GET("/api/v1/user/fitness-profile").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(404).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("not found")
	})

	// Create profile first
	profileData := map[string]interface{}{
		"date_of_birth":     "1992-03-10",
		"gender":            "other",
		"height_cm":         170.0,
		"current_weight_kg": 65.0,
		"primary_goal":      "strength",
	}

	e.POST("/api/v1/user/fitness-profile").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(profileData).
		Expect().
		Status(201)

	t.Run("Successful Get Profile", func(t *testing.T) {
		response := e.GET("/api/v1/user/fitness-profile").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("retrieved")

		data := response.Value("data").Object()
		data.Value("id").String().NotEmpty()
		data.Value("date_of_birth").String().IsEqual("1992-03-10")
		data.Value("gender").String().IsEqual("other")
		data.Value("height_cm").Number().IsEqual(170.0)
		data.Value("current_weight_kg").Number().IsEqual(65.0)
		data.Value("primary_goal").String().IsEqual("strength")
	})

	t.Run("Get Profile Without Auth", func(t *testing.T) {
		response := e.GET("/api/v1/user/fitness-profile").
			Expect().
			Status(401).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})
}

func testUpdateFitnessProfile(t *testing.T, e *httpexpect.Expect) {
	token := createTestUserAndGetToken(e, "updateprofile@example.com", "UpdatePass123!", "Update", "Profile")

	// Create profile first
	profileData := map[string]interface{}{
		"date_of_birth":     "1988-12-25",
		"gender":            "male",
		"height_cm":         175.0,
		"current_weight_kg": 85.0,
		"primary_goal":      "general_fitness",
	}

	e.POST("/api/v1/user/fitness-profile").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(profileData).
		Expect().
		Status(201)

	t.Run("Successful Full Update", func(t *testing.T) {
		updateData := map[string]interface{}{
			"current_weight_kg":              82.0,
			"primary_goal":                   "muscle_gain",
			"target_weight_kg":               78.0,
			"target_weekly_workouts":         4,
			"activity_level":                 "very_active",
			"training_locations":             []string{"gym", "outdoors"},
			"preferred_workout_duration_mins": 90,
			"available_days":                 []string{"monday", "wednesday", "friday", "sunday"},
		}

		response := e.PUT("/api/v1/user/fitness-profile").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(updateData).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("updated")

		data := response.Value("data").Object()
		data.Value("current_weight_kg").Number().IsEqual(82.0)
		data.Value("primary_goal").String().IsEqual("muscle_gain")
		data.Value("target_weight_kg").Number().IsEqual(78.0)
		data.Value("target_weekly_workouts").Number().IsEqual(4)
		data.Value("activity_level").String().IsEqual("very_active")
		data.Value("training_locations").Array().Length().IsEqual(2)
		data.Value("preferred_workout_duration_mins").Number().IsEqual(90)
		data.Value("available_days").Array().Length().IsEqual(4)
	})

	t.Run("Partial Update - Weight Only", func(t *testing.T) {
		updateData := map[string]interface{}{
			"current_weight_kg": 80.0,
		}

		response := e.PUT("/api/v1/user/fitness-profile").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(updateData).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("current_weight_kg").Number().IsEqual(80.0)
		// Other fields should remain unchanged
		data.Value("primary_goal").String().IsEqual("muscle_gain")
		data.Value("activity_level").String().IsEqual("very_active")
	})

	t.Run("Update Non-existent Profile", func(t *testing.T) {
		newToken := createTestUserAndGetToken(e, "nonexistent@example.com", "NoProfile123!", "No", "Profile")

		updateData := map[string]interface{}{
			"current_weight_kg": 70.0,
		}

		response := e.PUT("/api/v1/user/fitness-profile").
			WithHeader("Authorization", "Bearer "+newToken).
			WithJSON(updateData).
			Expect().
			Status(404).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("not found")
	})

	t.Run("Update Without Auth", func(t *testing.T) {
		updateData := map[string]interface{}{
			"current_weight_kg": 75.0,
		}

		response := e.PUT("/api/v1/user/fitness-profile").
			WithJSON(updateData).
			Expect().
			Status(401).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})
}

func testDeleteFitnessProfile(t *testing.T, e *httpexpect.Expect) {
	token := createTestUserAndGetToken(e, "deleteprofile@example.com", "DeletePass123!", "Delete", "Profile")

	t.Run("Delete Non-existent Profile", func(t *testing.T) {
		response := e.DELETE("/api/v1/user/fitness-profile").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(404).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("not found")
	})

	// Create profile first
	profileData := map[string]interface{}{
		"date_of_birth":     "1995-07-04",
		"gender":            "female",
		"height_cm":         162.0,
		"current_weight_kg": 55.0,
		"primary_goal":      "flexibility",
	}

	e.POST("/api/v1/user/fitness-profile").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(profileData).
		Expect().
		Status(201)

	t.Run("Successful Delete", func(t *testing.T) {
		e.DELETE("/api/v1/user/fitness-profile").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(204)

		// Verify profile is deleted
		response := e.GET("/api/v1/user/fitness-profile").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(404).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("Delete Without Auth", func(t *testing.T) {
		response := e.DELETE("/api/v1/user/fitness-profile").
			Expect().
			Status(401).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})
}

func testLogWeight(t *testing.T, e *httpexpect.Expect) {
	token := createTestUserAndGetToken(e, "logweight@example.com", "LogWeightPass123!", "Log", "Weight")

	// Create profile first
	profileData := map[string]interface{}{
		"date_of_birth":     "1990-01-01",
		"gender":            "male",
		"height_cm":         180.0,
		"current_weight_kg": 85.0,
		"primary_goal":      "weight_loss",
		"target_weight_kg":  75.0,
	}

	e.POST("/api/v1/user/fitness-profile").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(profileData).
		Expect().
		Status(201)

	t.Run("Successful Weight Log", func(t *testing.T) {
		weightData := map[string]interface{}{
			"weight_kg": 83.5,
		}

		response := e.POST("/api/v1/user/fitness-profile/log-weight").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(weightData).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("logged")

		data := response.Value("data").Object()
		data.Value("current_weight_kg").Number().IsEqual(83.5)
		data.Value("current_weight_lbs").Number().Gt(0)
	})

	t.Run("Log Weight Without Profile", func(t *testing.T) {
		noProfileToken := createTestUserAndGetToken(e, "noprofile@example.com", "NoProfile123!", "No", "Profile")

		weightData := map[string]interface{}{
			"weight_kg": 70.0,
		}

		response := e.POST("/api/v1/user/fitness-profile/log-weight").
			WithHeader("Authorization", "Bearer "+noProfileToken).
			WithJSON(weightData).
			Expect().
			Status(404).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("Log Weight Without Auth", func(t *testing.T) {
		weightData := map[string]interface{}{
			"weight_kg": 70.0,
		}

		response := e.POST("/api/v1/user/fitness-profile/log-weight").
			WithJSON(weightData).
			Expect().
			Status(401).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})
}

func testFitnessProfileValidation(t *testing.T, e *httpexpect.Expect) {
	token := createTestUserAndGetToken(e, "validation@example.com", "ValidationPass123!", "Valid", "Test")

	testCases := []struct {
		name        string
		profileData map[string]interface{}
	}{
		{
			name: "Missing Date of Birth",
			profileData: map[string]interface{}{
				"gender":            "male",
				"height_cm":         180.0,
				"current_weight_kg": 80.0,
				"primary_goal":      "general_fitness",
			},
		},
		{
			name: "Invalid Gender",
			profileData: map[string]interface{}{
				"date_of_birth":     "1990-05-15",
				"gender":            "invalid",
				"height_cm":         180.0,
				"current_weight_kg": 80.0,
				"primary_goal":      "general_fitness",
			},
		},
		{
			name: "Height Too Low",
			profileData: map[string]interface{}{
				"date_of_birth":     "1990-05-15",
				"gender":            "male",
				"height_cm":         40.0,
				"current_weight_kg": 80.0,
				"primary_goal":      "general_fitness",
			},
		},
		{
			name: "Height Too High",
			profileData: map[string]interface{}{
				"date_of_birth":     "1990-05-15",
				"gender":            "male",
				"height_cm":         350.0,
				"current_weight_kg": 80.0,
				"primary_goal":      "general_fitness",
			},
		},
		{
			name: "Weight Too Low",
			profileData: map[string]interface{}{
				"date_of_birth":     "1990-05-15",
				"gender":            "male",
				"height_cm":         180.0,
				"current_weight_kg": 10.0,
				"primary_goal":      "general_fitness",
			},
		},
		{
			name: "Invalid Primary Goal",
			profileData: map[string]interface{}{
				"date_of_birth":     "1990-05-15",
				"gender":            "male",
				"height_cm":         180.0,
				"current_weight_kg": 80.0,
				"primary_goal":      "invalid_goal",
			},
		},
		{
			name: "Invalid Activity Level",
			profileData: map[string]interface{}{
				"date_of_birth":     "1990-05-15",
				"gender":            "male",
				"height_cm":         180.0,
				"current_weight_kg": 80.0,
				"primary_goal":      "general_fitness",
				"activity_level":    "super_active",
			},
		},
		{
			name: "Invalid Training Location",
			profileData: map[string]interface{}{
				"date_of_birth":       "1990-05-15",
				"gender":              "male",
				"height_cm":           180.0,
				"current_weight_kg":   80.0,
				"primary_goal":        "general_fitness",
				"training_locations":  []string{"space"},
			},
		},
		{
			name: "Invalid Day",
			profileData: map[string]interface{}{
				"date_of_birth":     "1990-05-15",
				"gender":            "male",
				"height_cm":         180.0,
				"current_weight_kg": 80.0,
				"primary_goal":      "general_fitness",
				"available_days":    []string{"funday"},
			},
		},
		{
			name: "Weekly Workouts Too High",
			profileData: map[string]interface{}{
				"date_of_birth":          "1990-05-15",
				"gender":                 "male",
				"height_cm":              180.0,
				"current_weight_kg":      80.0,
				"primary_goal":           "general_fitness",
				"target_weekly_workouts": 10,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := e.POST("/api/v1/user/fitness-profile").
				WithHeader("Authorization", "Bearer "+token).
				WithJSON(tc.profileData).
				Expect().
				Status(400).
				JSON().
				Object()

			response.Value("success").Boolean().IsFalse()
		})
	}
}
