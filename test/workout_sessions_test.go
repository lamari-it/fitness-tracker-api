package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
)

func TestWorkoutSessionEndpoints(t *testing.T) {
	e := SetupTestApp(t)

	t.Run("User Workout Session Flow", func(t *testing.T) {
		CleanDatabase(t)
		testUserWorkoutSessionFlow(t, e)
	})

	t.Run("Trainer Logs Session For Client", func(t *testing.T) {
		CleanDatabase(t)
		testTrainerLogsForClient(t, e)
	})

	t.Run("Exercise and Set Logging", func(t *testing.T) {
		CleanDatabase(t)
		testExerciseAndSetLogging(t, e)
	})

	t.Run("Authorization Checks", func(t *testing.T) {
		CleanDatabase(t)
		testWorkoutSessionAuthorization(t, e)
	})

	t.Run("Weight Conversion", func(t *testing.T) {
		CleanDatabase(t)
		testWeightConversion(t, e)
	})
}

func testUserWorkoutSessionFlow(t *testing.T, e *httpexpect.Expect) {
	// Create user
	userToken := createTestUserAndGetToken(e, "user@example.com", "UserPass123!", "Test", "User")

	var sessionID string

	t.Run("Create Workout Session", func(t *testing.T) {
		response := e.POST("/api/v1/workout-sessions").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"notes": "Morning workout",
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("created")

		data := response.Value("data").Object()
		sessionID = data.Value("id").String().Raw()
		data.Value("notes").String().IsEqual("Morning workout")
		data.Value("started_at").String().NotEmpty()
		data.Value("ended_at").IsNull()
	})

	t.Run("Get Workout Session", func(t *testing.T) {
		response := e.GET("/api/v1/workout-sessions/"+sessionID).
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("data").Object().Value("id").String().IsEqual(sessionID)
	})

	t.Run("List Workout Sessions", func(t *testing.T) {
		response := e.GET("/api/v1/workout-sessions").
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("data").Array().Length().IsEqual(1)
	})

	t.Run("Update Workout Session", func(t *testing.T) {
		response := e.PUT("/api/v1/workout-sessions/"+sessionID).
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"notes": "Updated notes - felt great today",
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("data").Object().Value("notes").String().IsEqual("Updated notes - felt great today")
	})

	t.Run("End Workout Session", func(t *testing.T) {
		response := e.PUT("/api/v1/workout-sessions/"+sessionID+"/end").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"notes": "Completed workout",
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("ended_at").String().NotEmpty()
		data.Value("duration_minutes").Number().Ge(0)
	})

	t.Run("Delete Workout Session", func(t *testing.T) {
		e.DELETE("/api/v1/workout-sessions/"+sessionID).
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(204)

		// Verify it's gone
		e.GET("/api/v1/workout-sessions/"+sessionID).
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(404)
	})
}

func testTrainerLogsForClient(t *testing.T, e *httpexpect.Expect) {
	// Seed specialties
	SeedTestSpecialties(t)
	specialtyIDs := GetSpecialtyIDs(t, "Strength Training", "Weight Loss")

	// Create trainer
	trainerToken := createTestUserAndGetToken(e, "trainer@example.com", "TrainerPass123!", "John", "Trainer")

	// Create trainer profile
	e.POST("/api/v1/trainers/profile").
		WithHeader("Authorization", "Bearer "+trainerToken).
		WithJSON(map[string]interface{}{
			"bio":           "Certified personal trainer.",
			"specialty_ids": specialtyIDs,
			"hourly_rate":   75.00,
			"location":      "New York, NY",
		}).
		Expect().
		Status(201)

	// Create client
	clientToken := createTestUserAndGetToken(e, "client@example.com", "ClientPass123!", "Jane", "Client")

	// Get client user ID
	clientResponse := e.GET("/api/v1/auth/profile").
		WithHeader("Authorization", "Bearer "+clientToken).
		Expect().
		Status(200).
		JSON().
		Object()
	clientID := clientResponse.Value("data").Object().Value("id").String().Raw()

	// Trainer invites client
	inviteResponse := e.POST("/api/v1/trainers/clients").
		WithHeader("Authorization", "Bearer "+trainerToken).
		WithJSON(map[string]interface{}{
			"client_id": clientID,
		}).
		Expect().
		Status(201).
		JSON().
		Object()
	invitationID := inviteResponse.Value("data").Object().Value("id").String().Raw()

	// Client accepts
	e.PUT("/api/v1/me/trainer-invitations/"+invitationID).
		WithHeader("Authorization", "Bearer "+clientToken).
		WithJSON(map[string]interface{}{
			"action": "accept",
		}).
		Expect().
		Status(200)

	var sessionID string

	t.Run("Trainer Creates Session For Client", func(t *testing.T) {
		response := e.POST("/api/v1/workout-sessions").
			WithHeader("Authorization", "Bearer "+trainerToken).
			WithJSON(map[string]interface{}{
				"user_id": clientID,
				"notes":   "PT session - leg day",
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Object()
		sessionID = data.Value("id").String().Raw()
		data.Value("user_id").String().IsEqual(clientID)
		data.Value("created_by_id").String().NotEmpty()
		data.Value("created_by_name").String().IsEqual("John Trainer")
	})

	t.Run("Trainer Can Access Session They Created", func(t *testing.T) {
		response := e.GET("/api/v1/workout-sessions/"+sessionID).
			WithHeader("Authorization", "Bearer "+trainerToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
	})

	t.Run("Client Can Access Their Own Session", func(t *testing.T) {
		response := e.GET("/api/v1/workout-sessions/"+sessionID).
			WithHeader("Authorization", "Bearer "+clientToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
	})

	t.Run("Trainer Can View Client Sessions They Created", func(t *testing.T) {
		response := e.GET("/api/v1/workout-sessions").
			WithHeader("Authorization", "Bearer "+trainerToken).
			WithQuery("client_id", clientID).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("data").Array().Length().IsEqual(1)
	})

	t.Run("Trainer Cannot Create Session For Non-Client", func(t *testing.T) {
		// Create another user who is not trainer's client
		createTestUserAndGetToken(e, "other@example.com", "OtherPass123!", "Other", "User")

		otherResponse := e.POST("/api/v1/auth/login").
			WithJSON(map[string]interface{}{
				"email":    "other@example.com",
				"password": "OtherPass123!",
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		// Get the other user's ID by logging in again
		e.POST("/api/v1/auth/login").
			WithJSON(map[string]interface{}{
				"email":    "other@example.com",
				"password": "OtherPass123!",
			}).
			Expect().
			Status(200)

		// Get token for other user to retrieve their ID
		otherToken := otherResponse.Value("data").Object().Value("token").String().Raw()
		otherUserResponse := e.GET("/api/v1/auth/profile").
			WithHeader("Authorization", "Bearer "+otherToken).
			Expect().
			Status(200).
			JSON().
			Object()
		otherUserID := otherUserResponse.Value("data").Object().Value("id").String().Raw()

		// Try to create session for non-client
		response := e.POST("/api/v1/workout-sessions").
			WithHeader("Authorization", "Bearer "+trainerToken).
			WithJSON(map[string]interface{}{
				"user_id": otherUserID,
				"notes":   "Should fail",
			}).
			Expect().
			Status(403).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("Not authorized")
	})
}

func testExerciseAndSetLogging(t *testing.T, e *httpexpect.Expect) {
	// Create user
	userToken := createTestUserAndGetToken(e, "user@example.com", "UserPass123!", "Test", "User")

	// Create an exercise first
	exerciseResponse := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"name":        "Bench Press",
			"description": "Compound chest exercise",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	exerciseID := exerciseResponse.Value("data").Object().Value("id").String().Raw()

	// Create workout session
	sessionResponse := e.POST("/api/v1/workout-sessions").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"notes": "Chest day",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	sessionID := sessionResponse.Value("data").Object().Value("id").String().Raw()

	var exerciseLogID string
	var setLogID string

	t.Run("Create Exercise Log", func(t *testing.T) {
		response := e.POST("/api/v1/exercise-logs").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"session_id":   sessionID,
				"exercise_id":  exerciseID,
				"order_number": 1,
				"notes":        "Felt strong today",
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Object()
		exerciseLogID = data.Value("id").String().Raw()
		data.Value("exercise_name").String().IsEqual("Bench Press")
		data.Value("order_number").Number().IsEqual(1)
	})

	t.Run("Get Exercise Log", func(t *testing.T) {
		response := e.GET("/api/v1/exercise-logs/"+exerciseLogID).
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("data").Object().Value("id").String().IsEqual(exerciseLogID)
	})

	t.Run("Create Set Log", func(t *testing.T) {
		response := e.POST("/api/v1/set-logs").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"exercise_log_id": exerciseLogID,
				"set_number":      1,
				"weight":          100.0,
				"weight_unit":     "kg",
				"reps":            8,
				"rpe":             7.5,
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Object()
		setLogID = data.Value("id").String().Raw()
		data.Value("weight_kg").Number().IsEqual(100.0)
		data.Value("weight_display").Number().IsEqual(100.0)
		data.Value("weight_display_unit").String().IsEqual("kg")
		data.Value("input_weight").Number().IsEqual(100.0)
		data.Value("input_weight_unit").String().IsEqual("kg")
		data.Value("reps").Number().IsEqual(8)
		data.Value("rpe").Number().IsEqual(7.5)
	})

	t.Run("Create Multiple Sets", func(t *testing.T) {
		// Set 2
		e.POST("/api/v1/set-logs").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"exercise_log_id": exerciseLogID,
				"set_number":      2,
				"weight":          100.0,
				"weight_unit":     "kg",
				"reps":            7,
				"rpe":             8.0,
			}).
			Expect().
			Status(201)

		// Set 3
		e.POST("/api/v1/set-logs").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"exercise_log_id": exerciseLogID,
				"set_number":      3,
				"weight":          100.0,
				"weight_unit":     "kg",
				"reps":            6,
				"rpe":             9.0,
			}).
			Expect().
			Status(201)
	})

	t.Run("Get Exercise Log With Sets", func(t *testing.T) {
		response := e.GET("/api/v1/exercise-logs/"+exerciseLogID).
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		data := response.Value("data").Object()
		data.Value("set_logs").Array().Length().IsEqual(3)
	})

	t.Run("Update Set Log", func(t *testing.T) {
		response := e.PUT("/api/v1/set-logs/"+setLogID).
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"weight":      105.0,
				"reps":        8,
				"weight_unit": "kg",
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("weight_kg").Number().IsEqual(105.0)
		data.Value("input_weight").Number().IsEqual(105.0)
	})

	t.Run("Update Exercise Log", func(t *testing.T) {
		response := e.PUT("/api/v1/exercise-logs/"+exerciseLogID).
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"notes":             "PR on bench press!",
				"difficulty_rating": 8,
				"difficulty_type":   "hard",
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("difficulty_rating").Number().IsEqual(8)
		data.Value("difficulty_type").String().IsEqual("hard")
	})

	t.Run("Delete Set Log", func(t *testing.T) {
		e.DELETE("/api/v1/set-logs/"+setLogID).
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(204)

		// Verify it's deleted
		e.GET("/api/v1/set-logs/"+setLogID).
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(404)
	})

	t.Run("Delete Exercise Log", func(t *testing.T) {
		e.DELETE("/api/v1/exercise-logs/"+exerciseLogID).
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(204)

		// Verify it's deleted
		e.GET("/api/v1/exercise-logs/"+exerciseLogID).
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(404)
	})
}

func testWorkoutSessionAuthorization(t *testing.T, e *httpexpect.Expect) {
	// Create two users
	user1Token := createTestUserAndGetToken(e, "user1@example.com", "User1Pass123!", "User", "One")
	user2Token := createTestUserAndGetToken(e, "user2@example.com", "User2Pass123!", "User", "Two")

	// Create an exercise for user1
	exerciseResponse := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+user1Token).
		WithJSON(map[string]interface{}{
			"name":        "Squat",
			"description": "Compound leg exercise",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	exerciseID := exerciseResponse.Value("data").Object().Value("id").String().Raw()

	// User 1 creates a session
	sessionResponse := e.POST("/api/v1/workout-sessions").
		WithHeader("Authorization", "Bearer "+user1Token).
		WithJSON(map[string]interface{}{
			"notes": "User 1's workout",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	sessionID := sessionResponse.Value("data").Object().Value("id").String().Raw()

	// User 1 creates an exercise log
	exerciseLogResponse := e.POST("/api/v1/exercise-logs").
		WithHeader("Authorization", "Bearer "+user1Token).
		WithJSON(map[string]interface{}{
			"session_id":   sessionID,
			"exercise_id":  exerciseID,
			"order_number": 1,
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	exerciseLogID := exerciseLogResponse.Value("data").Object().Value("id").String().Raw()

	// User 1 creates a set log
	setLogResponse := e.POST("/api/v1/set-logs").
		WithHeader("Authorization", "Bearer "+user1Token).
		WithJSON(map[string]interface{}{
			"exercise_log_id": exerciseLogID,
			"set_number":      1,
			"weight":          100.0,
			"reps":            5,
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	setLogID := setLogResponse.Value("data").Object().Value("id").String().Raw()

	t.Run("User Cannot Access Other User's Session", func(t *testing.T) {
		e.GET("/api/v1/workout-sessions/"+sessionID).
			WithHeader("Authorization", "Bearer "+user2Token).
			Expect().
			Status(404)
	})

	t.Run("User Cannot Update Other User's Session", func(t *testing.T) {
		response := e.PUT("/api/v1/workout-sessions/"+sessionID).
			WithHeader("Authorization", "Bearer "+user2Token).
			WithJSON(map[string]interface{}{
				"notes": "Hacked!",
			}).
			Expect().
			Status(403).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("User Cannot Delete Other User's Session", func(t *testing.T) {
		response := e.DELETE("/api/v1/workout-sessions/"+sessionID).
			WithHeader("Authorization", "Bearer "+user2Token).
			Expect().
			Status(403).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("User Cannot Access Other User's Exercise Log", func(t *testing.T) {
		e.GET("/api/v1/exercise-logs/"+exerciseLogID).
			WithHeader("Authorization", "Bearer "+user2Token).
			Expect().
			Status(404)
	})

	t.Run("User Cannot Update Other User's Exercise Log", func(t *testing.T) {
		response := e.PUT("/api/v1/exercise-logs/"+exerciseLogID).
			WithHeader("Authorization", "Bearer "+user2Token).
			WithJSON(map[string]interface{}{
				"notes": "Hacked!",
			}).
			Expect().
			Status(403).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("User Cannot Access Other User's Set Log", func(t *testing.T) {
		e.GET("/api/v1/set-logs/"+setLogID).
			WithHeader("Authorization", "Bearer "+user2Token).
			Expect().
			Status(404)
	})

	t.Run("User Cannot Update Other User's Set Log", func(t *testing.T) {
		response := e.PUT("/api/v1/set-logs/"+setLogID).
			WithHeader("Authorization", "Bearer "+user2Token).
			WithJSON(map[string]interface{}{
				"weight": 200.0,
			}).
			Expect().
			Status(403).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("User Cannot Create Exercise Log In Other User's Session", func(t *testing.T) {
		response := e.POST("/api/v1/exercise-logs").
			WithHeader("Authorization", "Bearer "+user2Token).
			WithJSON(map[string]interface{}{
				"session_id":   sessionID,
				"exercise_id":  exerciseID,
				"order_number": 2,
			}).
			Expect().
			Status(403).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("Not authorized")
	})

	t.Run("User Cannot Create Set Log In Other User's Exercise Log", func(t *testing.T) {
		response := e.POST("/api/v1/set-logs").
			WithHeader("Authorization", "Bearer "+user2Token).
			WithJSON(map[string]interface{}{
				"exercise_log_id": exerciseLogID,
				"set_number":      2,
				"weight":          50.0,
				"reps":            10,
			}).
			Expect().
			Status(403).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("Not authorized")
	})

	t.Run("Invalid Session ID Returns 404", func(t *testing.T) {
		e.GET("/api/v1/workout-sessions/11111111-1111-1111-1111-111111111111").
			WithHeader("Authorization", "Bearer "+user1Token).
			Expect().
			Status(404)
	})

	t.Run("Invalid Exercise Log ID Returns 404", func(t *testing.T) {
		e.GET("/api/v1/exercise-logs/11111111-1111-1111-1111-111111111111").
			WithHeader("Authorization", "Bearer "+user1Token).
			Expect().
			Status(404)
	})

	t.Run("Invalid Set Log ID Returns 404", func(t *testing.T) {
		e.GET("/api/v1/set-logs/11111111-1111-1111-1111-111111111111").
			WithHeader("Authorization", "Bearer "+user1Token).
			Expect().
			Status(404)
	})
}

func testWeightConversion(t *testing.T, e *httpexpect.Expect) {
	// Create a user and get token
	userEmail := fmt.Sprintf("weight_test_%d@example.com", time.Now().UnixNano())
	userToken := createTestUserAndGetToken(e, userEmail, "TestPassword123!", "Weight", "Tester")

	// Create exercise for testing
	exerciseResponse := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"name":        "Weight Test Squat",
			"description": "Test exercise for weight conversion",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	exerciseID := exerciseResponse.Value("data").Object().Value("id").String().Raw()

	// Create a session
	sessionResponse := e.POST("/api/v1/workout-sessions").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"started_at": time.Now().Format(time.RFC3339),
			"notes":      "Weight conversion test session",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	sessionID := sessionResponse.Value("data").Object().Value("id").String().Raw()

	// Create an exercise log
	exerciseLogResponse := e.POST("/api/v1/exercise-logs").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"session_id":   sessionID,
			"exercise_id":  exerciseID,
			"order_number": 1,
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	exerciseLogID := exerciseLogResponse.Value("data").Object().Value("id").String().Raw()

	t.Run("Create Set Log With Kg", func(t *testing.T) {
		response := e.POST("/api/v1/set-logs").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"exercise_log_id": exerciseLogID,
				"set_number":      1,
				"weight":          100.0,
				"weight_unit":     "kg",
				"reps":            8,
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		data := response.Value("data").Object()
		// Weight should be stored as 100 kg
		data.Value("weight_kg").Number().IsEqual(100.0)
		data.Value("weight_display").Number().IsEqual(100.0)
		data.Value("weight_display_unit").String().IsEqual("kg")
		data.Value("input_weight").Number().IsEqual(100.0)
		data.Value("input_weight_unit").String().IsEqual("kg")
	})

	t.Run("Create Set Log With Lbs", func(t *testing.T) {
		response := e.POST("/api/v1/set-logs").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"exercise_log_id": exerciseLogID,
				"set_number":      2,
				"weight":          225.0,
				"weight_unit":     "lb",
				"reps":            5,
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		data := response.Value("data").Object()
		// 225 lbs should be converted to ~102.06 kg
		data.Value("weight_kg").Number().InDelta(102.06, 0.1)
		data.Value("input_weight").Number().IsEqual(225.0)
		data.Value("input_weight_unit").String().IsEqual("lb")
	})

	t.Run("Create Set Log With Lbs Alias", func(t *testing.T) {
		response := e.POST("/api/v1/set-logs").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"exercise_log_id": exerciseLogID,
				"set_number":      3,
				"weight":          135.0,
				"weight_unit":     "lbs",
				"reps":            10,
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		data := response.Value("data").Object()
		// 135 lbs should be converted to ~61.23 kg
		data.Value("weight_kg").Number().InDelta(61.23, 0.1)
		data.Value("input_weight").Number().IsEqual(135.0)
		// lbs should be normalized to lb
		data.Value("input_weight_unit").String().IsEqual("lb")
	})

	t.Run("Update Set Log With Different Unit", func(t *testing.T) {
		// First create a set log
		createResponse := e.POST("/api/v1/set-logs").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"exercise_log_id": exerciseLogID,
				"set_number":      4,
				"weight":          50.0,
				"weight_unit":     "kg",
				"reps":            12,
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		setLogID := createResponse.Value("data").Object().Value("id").String().Raw()

		// Update it with lbs
		updateResponse := e.PUT("/api/v1/set-logs/"+setLogID).
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"weight":      110.0,
				"weight_unit": "lb",
				"reps":        10,
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		data := updateResponse.Value("data").Object()
		// 110 lbs should be converted to ~49.90 kg
		data.Value("weight_kg").Number().InDelta(49.90, 0.1)
		data.Value("input_weight").Number().IsEqual(110.0)
		data.Value("input_weight_unit").String().IsEqual("lb")
	})

	t.Run("Get Set Log Returns Converted Values", func(t *testing.T) {
		// Create a set with lbs
		createResponse := e.POST("/api/v1/set-logs").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"exercise_log_id": exerciseLogID,
				"set_number":      5,
				"weight":          200.0,
				"weight_unit":     "lb",
				"reps":            3,
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		setLogID := createResponse.Value("data").Object().Value("id").String().Raw()

		// Get the set log and verify conversion
		getResponse := e.GET("/api/v1/set-logs/"+setLogID).
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		data := getResponse.Value("data").Object()
		// 200 lbs = ~90.72 kg
		data.Value("weight_kg").Number().InDelta(90.72, 0.1)
		data.Value("input_weight").Number().IsEqual(200.0)
		data.Value("input_weight_unit").String().IsEqual("lb")
		// Default user preference is kg (no fitness profile), so display should be kg
		data.Value("weight_display_unit").String().IsEqual("kg")
	})

	t.Run("Zero Weight Is Allowed", func(t *testing.T) {
		response := e.POST("/api/v1/set-logs").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"exercise_log_id": exerciseLogID,
				"set_number":      6,
				"weight":          0,
				"weight_unit":     "kg",
				"reps":            20,
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		data := response.Value("data").Object()
		data.Value("weight_kg").Number().IsEqual(0)
		data.Value("input_weight").Number().IsEqual(0)
	})
}
