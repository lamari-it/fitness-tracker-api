package test

import (
	"testing"

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

	t.Run("Session With Prescriptions Auto-Creates Structure", func(t *testing.T) {
		CleanDatabase(t)
		testSessionWithPrescriptions(t, e)
	})

	t.Run("Session Block Operations", func(t *testing.T) {
		CleanDatabase(t)
		testSessionBlockOperations(t, e)
	})

	t.Run("Session Exercise Operations", func(t *testing.T) {
		CleanDatabase(t)
		testSessionExerciseOperations(t, e)
	})

	t.Run("Session Set Operations", func(t *testing.T) {
		CleanDatabase(t)
		testSessionSetOperations(t, e)
	})

	t.Run("Authorization Checks", func(t *testing.T) {
		CleanDatabase(t)
		testWorkoutSessionAuthorization(t, e)
	})

	t.Run("Session Logging Authorization", func(t *testing.T) {
		CleanDatabase(t)
		testSessionLoggingAuthorization(t, e)
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
		data.Value("completed").Boolean().IsFalse()
		data.Value("blocks").Array().Length().IsEqual(0) // No prescriptions, empty blocks
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
				"notes":              "Completed workout",
				"perceived_intensity": 7,
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("ended_at").String().NotEmpty()
		data.Value("completed").Boolean().IsTrue()
		data.Value("perceived_intensity").Number().IsEqual(7)
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

func testSessionWithPrescriptions(t *testing.T, e *httpexpect.Expect) {
	// Create user
	userToken := createTestUserAndGetToken(e, "user@example.com", "UserPass123!", "Test", "User")

	// Seed RPE scale
	SeedTestGlobalRPEScale(t)

	// Create an exercise
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

	// Create another exercise
	exercise2Response := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"name":        "Squat",
			"description": "Compound leg exercise",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	exercise2ID := exercise2Response.Value("data").Object().Value("id").String().Raw()

	// Create a workout
	workoutResponse := e.POST("/api/v1/workouts/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"title":       "Full Body Workout",
			"description": "Test workout with prescriptions",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	workoutID := workoutResponse.Value("data").Object().Value("id").String().Raw()

	// Add prescription group with exercises
	rpeValueID := GetRPEValueID(t, 7)

	e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"type":              "straight",
			"group_order":       1,
			"group_rounds":      1,
			"rest_between_sets": 90,
			"exercises": []map[string]interface{}{
				{
					"exercise_id":    exerciseID,
					"exercise_order": 1,
					"sets":           3,
					"reps":           10,
					"target_weight_kg": 80.0,
					"rpe_value_id":   rpeValueID,
				},
			},
		}).
		Expect().
		Status(201)

	// Add another prescription group
	e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"type":        "straight",
			"group_order": 2,
			"exercises": []map[string]interface{}{
				{
					"exercise_id":    exercise2ID,
					"exercise_order": 1,
					"sets":           4,
					"reps":           8,
					"target_weight_kg": 100.0,
				},
			},
		}).
		Expect().
		Status(201)

	var sessionID string
	var blockID string
	var exerciseLogID string
	var setID string

	t.Run("Create Session From Workout Auto-Creates Structure", func(t *testing.T) {
		response := e.POST("/api/v1/workout-sessions").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"workout_id": workoutID,
				"notes":      "Testing auto-creation",
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Object()
		sessionID = data.Value("id").String().Raw()
		data.Value("workout_id").String().IsEqual(workoutID)
		data.Value("workout_title").String().IsEqual("Full Body Workout")

		// Verify blocks were created
		blocks := data.Value("blocks").Array()
		blocks.Length().IsEqual(2)

		// First block
		block1 := blocks.Value(0).Object()
		block1.Value("block_order").Number().IsEqual(1)
		blockID = block1.Value("id").String().Raw()

		exercises := block1.Value("exercises").Array()
		exercises.Length().IsEqual(1)

		ex1 := exercises.Value(0).Object()
		ex1.Value("exercise_name").String().IsEqual("Bench Press")
		exerciseLogID = ex1.Value("id").String().Raw()

		sets := ex1.Value("sets").Array()
		sets.Length().IsEqual(3) // 3 sets as prescribed

		set1 := sets.Value(0).Object()
		set1.Value("set_number").Number().IsEqual(1)
		set1.Value("actual_reps").Number().IsEqual(10) // Pre-filled from prescription
		set1.Value("actual_weight_kg").Number().IsEqual(80.0) // Pre-filled from prescription
		setID = set1.Value("id").String().Raw()

		// Second block
		block2 := blocks.Value(1).Object()
		block2.Value("block_order").Number().IsEqual(2)
		block2Ex := block2.Value("exercises").Array().Value(0).Object()
		block2Ex.Value("exercise_name").String().IsEqual("Squat")
		block2Ex.Value("sets").Array().Length().IsEqual(4) // 4 sets as prescribed
	})

	t.Run("Complete Session Block", func(t *testing.T) {
		response := e.PUT("/api/v1/session-blocks/"+blockID+"/complete").
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("completed_at").String().NotEmpty()
		data.Value("skipped").Boolean().IsFalse()
	})

	t.Run("Complete Session Exercise", func(t *testing.T) {
		response := e.PUT("/api/v1/session-exercises/"+exerciseLogID+"/complete").
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("completed_at").String().NotEmpty()
	})

	t.Run("Update Session Set", func(t *testing.T) {
		response := e.PUT("/api/v1/session-sets/"+setID).
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"actual_reps":      12, // Did more reps than prescribed
				"actual_weight_kg": 85.0, // Heavier weight
				"was_failure":      false,
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("actual_reps").Number().IsEqual(12)
		data.Value("actual_weight_kg").Number().IsEqual(85.0)
	})

	t.Run("Complete Session Set", func(t *testing.T) {
		response := e.PUT("/api/v1/session-sets/"+setID+"/complete").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"actual_reps": 12,
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("completed").Boolean().IsTrue()
	})

	t.Run("Add Extra Set To Exercise", func(t *testing.T) {
		response := e.POST("/api/v1/session-exercises/"+exerciseLogID+"/sets").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"actual_reps":      8,
				"actual_weight_kg": 90.0,
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("set_number").Number().IsEqual(4) // 4th set (extra)
		data.Value("actual_reps").Number().IsEqual(8)
		data.Value("actual_weight_kg").Number().IsEqual(90.0)
	})

	t.Run("End Session With Complete", func(t *testing.T) {
		response := e.PUT("/api/v1/workout-sessions/"+sessionID+"/end").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"perceived_intensity": 8,
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("completed").Boolean().IsTrue()
		data.Value("perceived_intensity").Number().IsEqual(8)
	})
}

func testSessionSetOperations(t *testing.T, e *httpexpect.Expect) {
	// Create user
	userToken := createTestUserAndGetToken(e, "user@example.com", "UserPass123!", "Test", "User")

	// Create an exercise
	exerciseResponse := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"name":        "Deadlift",
			"description": "Compound back exercise",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	exerciseID := exerciseResponse.Value("data").Object().Value("id").String().Raw()

	// Create a workout with prescription
	workoutResponse := e.POST("/api/v1/workouts/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"title": "Back Day",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	workoutID := workoutResponse.Value("data").Object().Value("id").String().Raw()

	e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"type":        "straight",
			"group_order": 1,
			"exercises": []map[string]interface{}{
				{
					"exercise_id":    exerciseID,
					"exercise_order": 1,
					"sets":           3,
					"reps":           5,
					"target_weight_kg": 140.0,
				},
			},
		}).
		Expect().
		Status(201)

	// Create session
	sessionResponse := e.POST("/api/v1/workout-sessions").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"workout_id": workoutID,
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	// Get set ID from created session
	sets := sessionResponse.Value("data").Object().
		Value("blocks").Array().Value(0).Object().
		Value("exercises").Array().Value(0).Object().
		Value("sets").Array()

	setID := sets.Value(0).Object().Value("id").String().Raw()

	t.Run("Get Session Set", func(t *testing.T) {
		response := e.GET("/api/v1/session-sets/"+setID).
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("set_number").Number().IsEqual(1)
	})

	t.Run("Update Set With Weight Conversion", func(t *testing.T) {
		// Update with lbs, should convert to kg
		response := e.PUT("/api/v1/session-sets/"+setID).
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"actual_weight_kg": 315.0,
				"weight_unit":      "lb",
				"actual_reps":      5,
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		// 315 lbs ≈ 142.88 kg
		data.Value("actual_weight_kg").Number().InDelta(142.88, 1.0)
	})

	t.Run("Delete Session Set", func(t *testing.T) {
		// Delete the last set
		lastSetID := sets.Value(2).Object().Value("id").String().Raw()

		e.DELETE("/api/v1/session-sets/"+lastSetID).
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(204)

		// Verify it's deleted
		e.GET("/api/v1/session-sets/"+lastSetID).
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(404)
	})
}

func testSessionBlockOperations(t *testing.T, e *httpexpect.Expect) {
	// Create user and setup
	userToken := createTestUserAndGetToken(e, "user@example.com", "UserPass123!", "Test", "User")

	// Create exercise
	exerciseResponse := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"name":        "Bench Press",
			"description": "Chest exercise",
		}).
		Expect().
		Status(201).
		JSON().
		Object()
	exerciseID := exerciseResponse.Value("data").Object().Value("id").String().Raw()

	// Create workout with prescription
	workoutResponse := e.POST("/api/v1/workouts/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"title": "Test Workout",
		}).
		Expect().
		Status(201).
		JSON().
		Object()
	workoutID := workoutResponse.Value("data").Object().Value("id").String().Raw()

	e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"type":        "straight",
			"group_order": 1,
			"exercises": []map[string]interface{}{
				{
					"exercise_id":    exerciseID,
					"exercise_order": 1,
					"sets":           3,
					"reps":           10,
				},
			},
		}).
		Expect().
		Status(201)

	// Create session from workout
	sessionResponse := e.POST("/api/v1/workout-sessions").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"workout_id": workoutID,
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	blockID := sessionResponse.Value("data").Object().
		Value("blocks").Array().Value(0).Object().
		Value("id").String().Raw()

	t.Run("Get Session Block", func(t *testing.T) {
		response := e.GET("/api/v1/session-blocks/"+blockID).
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("id").String().IsEqual(blockID)
		data.Value("block_order").Number().IsEqual(1)
		data.Value("exercises").Array().Length().IsEqual(1)
	})

	t.Run("Skip Session Block", func(t *testing.T) {
		response := e.PUT("/api/v1/session-blocks/"+blockID+"/skip").
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("skipped").Boolean().IsTrue()
		data.Value("completed_at").IsNull()
	})

	t.Run("Complete Session Block After Skip", func(t *testing.T) {
		response := e.PUT("/api/v1/session-blocks/"+blockID+"/complete").
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("skipped").Boolean().IsFalse()
		data.Value("completed_at").String().NotEmpty()
	})

	t.Run("Update Session Block RPE", func(t *testing.T) {
		response := e.PUT("/api/v1/session-blocks/"+blockID+"/rpe").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"perceived_exertion": 8,
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("perceived_exertion").Number().IsEqual(8)
	})

	t.Run("Invalid Block ID Returns 404", func(t *testing.T) {
		e.GET("/api/v1/session-blocks/11111111-1111-1111-1111-111111111111").
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(404)
	})
}

func testSessionExerciseOperations(t *testing.T, e *httpexpect.Expect) {
	// Create user and setup
	userToken := createTestUserAndGetToken(e, "user@example.com", "UserPass123!", "Test", "User")
	SeedTestGlobalRPEScale(t)

	// Create exercise
	exerciseResponse := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"name":        "Squat",
			"description": "Leg exercise",
		}).
		Expect().
		Status(201).
		JSON().
		Object()
	exerciseID := exerciseResponse.Value("data").Object().Value("id").String().Raw()

	// Create workout with prescription
	workoutResponse := e.POST("/api/v1/workouts/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"title": "Leg Day",
		}).
		Expect().
		Status(201).
		JSON().
		Object()
	workoutID := workoutResponse.Value("data").Object().Value("id").String().Raw()

	rpeValueID := GetRPEValueID(t, 7)

	e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"type":        "straight",
			"group_order": 1,
			"exercises": []map[string]interface{}{
				{
					"exercise_id":      exerciseID,
					"exercise_order":   1,
					"sets":             3,
					"reps":             8,
					"target_weight_kg": 100.0,
					"rpe_value_id":     rpeValueID,
				},
			},
		}).
		Expect().
		Status(201)

	// Create session from workout
	sessionResponse := e.POST("/api/v1/workout-sessions").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"workout_id": workoutID,
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	exerciseLogID := sessionResponse.Value("data").Object().
		Value("blocks").Array().Value(0).Object().
		Value("exercises").Array().Value(0).Object().
		Value("id").String().Raw()

	t.Run("Get Session Exercise", func(t *testing.T) {
		response := e.GET("/api/v1/session-exercises/"+exerciseLogID).
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("id").String().IsEqual(exerciseLogID)
		data.Value("exercise_name").String().IsEqual("Squat")
		data.Value("sets").Array().Length().IsEqual(3)
	})

	t.Run("Skip Session Exercise", func(t *testing.T) {
		response := e.PUT("/api/v1/session-exercises/"+exerciseLogID+"/skip").
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("skipped").Boolean().IsTrue()
	})

	t.Run("Complete Session Exercise", func(t *testing.T) {
		response := e.PUT("/api/v1/session-exercises/"+exerciseLogID+"/complete").
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("skipped").Boolean().IsFalse()
		data.Value("completed_at").String().NotEmpty()
	})

	t.Run("Update Session Exercise Notes", func(t *testing.T) {
		response := e.PUT("/api/v1/session-exercises/"+exerciseLogID+"/notes").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"notes": "Felt strong today - PR attempt next week",
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("notes").String().IsEqual("Felt strong today - PR attempt next week")
	})

	t.Run("Add Extra Set To Exercise", func(t *testing.T) {
		response := e.POST("/api/v1/session-exercises/"+exerciseLogID+"/sets").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"actual_reps":      6,
				"actual_weight_kg": 110.0,
				"rpe_value_id":     rpeValueID,
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("set_number").Number().IsEqual(4) // 4th set
		data.Value("actual_reps").Number().IsEqual(6)
		data.Value("actual_weight_kg").Number().IsEqual(110.0)
	})

	t.Run("Add Set With Weight In Lbs", func(t *testing.T) {
		response := e.POST("/api/v1/session-exercises/"+exerciseLogID+"/sets").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"actual_reps":      5,
				"actual_weight_kg": 225.0,
				"weight_unit":      "lb",
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("set_number").Number().IsEqual(5)
		// 225 lbs ≈ 102.06 kg
		data.Value("actual_weight_kg").Number().InDelta(102.06, 1.0)
	})

	t.Run("Invalid Exercise ID Returns 404", func(t *testing.T) {
		e.GET("/api/v1/session-exercises/11111111-1111-1111-1111-111111111111").
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(404)
	})
}

func testWorkoutSessionAuthorization(t *testing.T, e *httpexpect.Expect) {
	// Create two users
	user1Token := createTestUserAndGetToken(e, "user1@example.com", "User1Pass123!", "User", "One")
	user2Token := createTestUserAndGetToken(e, "user2@example.com", "User2Pass123!", "User", "Two")

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

	t.Run("Invalid Session ID Returns 404", func(t *testing.T) {
		e.GET("/api/v1/workout-sessions/11111111-1111-1111-1111-111111111111").
			WithHeader("Authorization", "Bearer "+user1Token).
			Expect().
			Status(404)
	})
}

func testSessionLoggingAuthorization(t *testing.T, e *httpexpect.Expect) {
	// Create two users
	user1Token := createTestUserAndGetToken(e, "user1@example.com", "User1Pass123!", "User", "One")
	user2Token := createTestUserAndGetToken(e, "user2@example.com", "User2Pass123!", "User", "Two")

	// Create exercise
	exerciseResponse := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+user1Token).
		WithJSON(map[string]interface{}{
			"name":        "Deadlift",
			"description": "Back exercise",
		}).
		Expect().
		Status(201).
		JSON().
		Object()
	exerciseID := exerciseResponse.Value("data").Object().Value("id").String().Raw()

	// User 1 creates workout with prescription
	workoutResponse := e.POST("/api/v1/workouts/").
		WithHeader("Authorization", "Bearer "+user1Token).
		WithJSON(map[string]interface{}{
			"title": "User 1 Workout",
		}).
		Expect().
		Status(201).
		JSON().
		Object()
	workoutID := workoutResponse.Value("data").Object().Value("id").String().Raw()

	e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
		WithHeader("Authorization", "Bearer "+user1Token).
		WithJSON(map[string]interface{}{
			"type":        "straight",
			"group_order": 1,
			"exercises": []map[string]interface{}{
				{
					"exercise_id":    exerciseID,
					"exercise_order": 1,
					"sets":           3,
					"reps":           5,
				},
			},
		}).
		Expect().
		Status(201)

	// User 1 creates session
	sessionResponse := e.POST("/api/v1/workout-sessions").
		WithHeader("Authorization", "Bearer "+user1Token).
		WithJSON(map[string]interface{}{
			"workout_id": workoutID,
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	blocks := sessionResponse.Value("data").Object().Value("blocks").Array()
	blockID := blocks.Value(0).Object().Value("id").String().Raw()
	exerciseLogID := blocks.Value(0).Object().Value("exercises").Array().Value(0).Object().Value("id").String().Raw()
	setID := blocks.Value(0).Object().Value("exercises").Array().Value(0).Object().Value("sets").Array().Value(0).Object().Value("id").String().Raw()

	t.Run("User Cannot Access Other User's Block", func(t *testing.T) {
		e.GET("/api/v1/session-blocks/"+blockID).
			WithHeader("Authorization", "Bearer "+user2Token).
			Expect().
			Status(404)
	})

	t.Run("User Cannot Complete Other User's Block", func(t *testing.T) {
		e.PUT("/api/v1/session-blocks/"+blockID+"/complete").
			WithHeader("Authorization", "Bearer "+user2Token).
			Expect().
			Status(403)
	})

	t.Run("User Cannot Skip Other User's Block", func(t *testing.T) {
		e.PUT("/api/v1/session-blocks/"+blockID+"/skip").
			WithHeader("Authorization", "Bearer "+user2Token).
			Expect().
			Status(403)
	})

	t.Run("User Cannot Update Other User's Block RPE", func(t *testing.T) {
		e.PUT("/api/v1/session-blocks/"+blockID+"/rpe").
			WithHeader("Authorization", "Bearer "+user2Token).
			WithJSON(map[string]interface{}{
				"perceived_exertion": 5,
			}).
			Expect().
			Status(403)
	})

	t.Run("User Cannot Access Other User's Exercise", func(t *testing.T) {
		e.GET("/api/v1/session-exercises/"+exerciseLogID).
			WithHeader("Authorization", "Bearer "+user2Token).
			Expect().
			Status(404)
	})

	t.Run("User Cannot Complete Other User's Exercise", func(t *testing.T) {
		e.PUT("/api/v1/session-exercises/"+exerciseLogID+"/complete").
			WithHeader("Authorization", "Bearer "+user2Token).
			Expect().
			Status(403)
	})

	t.Run("User Cannot Skip Other User's Exercise", func(t *testing.T) {
		e.PUT("/api/v1/session-exercises/"+exerciseLogID+"/skip").
			WithHeader("Authorization", "Bearer "+user2Token).
			Expect().
			Status(403)
	})

	t.Run("User Cannot Update Other User's Exercise Notes", func(t *testing.T) {
		e.PUT("/api/v1/session-exercises/"+exerciseLogID+"/notes").
			WithHeader("Authorization", "Bearer "+user2Token).
			WithJSON(map[string]interface{}{
				"notes": "Hacked!",
			}).
			Expect().
			Status(403)
	})

	t.Run("User Cannot Add Set To Other User's Exercise", func(t *testing.T) {
		e.POST("/api/v1/session-exercises/"+exerciseLogID+"/sets").
			WithHeader("Authorization", "Bearer "+user2Token).
			WithJSON(map[string]interface{}{
				"actual_reps":      10,
				"actual_weight_kg": 50.0,
			}).
			Expect().
			Status(403)
	})

	t.Run("User Cannot Access Other User's Set", func(t *testing.T) {
		e.GET("/api/v1/session-sets/"+setID).
			WithHeader("Authorization", "Bearer "+user2Token).
			Expect().
			Status(404)
	})

	t.Run("User Cannot Update Other User's Set", func(t *testing.T) {
		e.PUT("/api/v1/session-sets/"+setID).
			WithHeader("Authorization", "Bearer "+user2Token).
			WithJSON(map[string]interface{}{
				"actual_reps": 100,
			}).
			Expect().
			Status(403)
	})

	t.Run("User Cannot Complete Other User's Set", func(t *testing.T) {
		e.PUT("/api/v1/session-sets/"+setID+"/complete").
			WithHeader("Authorization", "Bearer "+user2Token).
			Expect().
			Status(403)
	})

	t.Run("User Cannot Delete Other User's Set", func(t *testing.T) {
		e.DELETE("/api/v1/session-sets/"+setID).
			WithHeader("Authorization", "Bearer "+user2Token).
			Expect().
			Status(403)
	})
}
