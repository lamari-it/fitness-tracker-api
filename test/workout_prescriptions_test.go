package test

import (
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestWorkoutPrescriptionEndpoints(t *testing.T) {
	e := SetupTestApp(t)

	t.Run("Basic Workout and Prescription Flow", func(t *testing.T) {
		CleanDatabase(t)
		testBasicWorkoutPrescriptionFlow(t, e)
	})

	t.Run("Prescription Types", func(t *testing.T) {
		CleanDatabase(t)
		testPrescriptionTypes(t, e)
	})

	t.Run("Prescription Group Operations", func(t *testing.T) {
		CleanDatabase(t)
		testPrescriptionGroupOperations(t, e)
	})

	t.Run("Reorder Prescription Groups", func(t *testing.T) {
		CleanDatabase(t)
		testReorderPrescriptionGroups(t, e)
	})

	t.Run("Add Exercise To Prescription Group", func(t *testing.T) {
		CleanDatabase(t)
		testAddExerciseToPrescriptionGroup(t, e)
	})

	t.Run("Duplicate Workout With Prescriptions", func(t *testing.T) {
		CleanDatabase(t)
		testDuplicateWorkoutWithPrescriptions(t, e)
	})

	t.Run("Prescription Authorization", func(t *testing.T) {
		CleanDatabase(t)
		testPrescriptionAuthorization(t, e)
	})

	t.Run("Prescription Validation", func(t *testing.T) {
		CleanDatabase(t)
		testPrescriptionValidation(t, e)
	})
}

func testBasicWorkoutPrescriptionFlow(t *testing.T, e *httpexpect.Expect) {
	// Seed global RPE scale for RPE value tests
	SeedTestGlobalRPEScale(t)

	// Get RPE value ID for RPE 8 (2 reps left in tank)
	rpeValueID := GetRPEValueID(t, 8)

	// Create user
	userToken := createTestUserAndGetToken(e, "user@example.com", "UserPass123!", "Test", "User")

	// Create exercises
	exercise1Response := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"name":        "Bench Press",
			"description": "Compound chest exercise",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	exercise1ID := exercise1Response.Value("data").Object().Value("id").String().Raw()

	exercise2Response := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"name":        "Incline Dumbbell Press",
			"description": "Upper chest exercise",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	exercise2ID := exercise2Response.Value("data").Object().Value("id").String().Raw()

	var workoutID string
	var groupID string

	t.Run("Create Workout", func(t *testing.T) {
		response := e.POST("/api/v1/workouts/").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"title":       "Push Day",
				"description": "Chest, shoulders, and triceps",
				"visibility":  "private",
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("created")

		data := response.Value("data").Object()
		workoutID = data.Value("id").String().Raw()
		data.Value("title").String().IsEqual("Push Day")
		data.Value("visibility").String().IsEqual("private")
	})

	t.Run("Create Straight Set Prescription Group", func(t *testing.T) {
		response := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"type":              "straight",
				"group_order":       1,
				"group_name":        "Bench Press",
				"rest_between_sets": 120,
				"exercises": []map[string]interface{}{
					{
						"exercise_id":    exercise1ID,
						"exercise_order": 1,
						"sets":           4,
						"reps":           8,
						"weight_kg":      60.0,
						"rpe_value_id":   rpeValueID,
					},
				},
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Object()
		groupID = data.Value("group_id").String().Raw()
		data.Value("type").String().IsEqual("straight")
		data.Value("group_order").Number().IsEqual(1)
		data.Value("group_name").String().IsEqual("Bench Press")
		data.Value("rest_between_sets").Number().IsEqual(120)

		exercises := data.Value("exercises").Array()
		exercises.Length().IsEqual(1)
		exercises.Value(0).Object().Value("sets").Number().IsEqual(4)
		exercises.Value(0).Object().Value("reps").Number().IsEqual(8)
		exercises.Value(0).Object().Value("weight_kg").Number().IsEqual(60.0)
		exercises.Value(0).Object().Value("rpe_value_id").String().IsEqual(rpeValueID)
	})

	t.Run("Get Workout With Prescriptions", func(t *testing.T) {
		response := e.GET("/api/v1/workouts/"+workoutID).
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Object()
		data.Value("id").String().IsEqual(workoutID)
		data.Value("title").String().IsEqual("Push Day")

		prescriptions := data.Value("prescriptions").Array()
		prescriptions.Length().IsEqual(1)
		prescriptions.Value(0).Object().Value("group_id").String().IsEqual(groupID)
		prescriptions.Value(0).Object().Value("type").String().IsEqual("straight")
	})

	t.Run("Get Workout Prescriptions", func(t *testing.T) {
		response := e.GET("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Array()
		data.Length().IsEqual(1)
		data.Value(0).Object().Value("group_id").String().IsEqual(groupID)

		// Verify RPE value is returned with full details
		exercises := data.Value(0).Object().Value("exercises").Array()
		exercises.Value(0).Object().Value("rpe_value_id").String().IsEqual(rpeValueID)
		rpeValue := exercises.Value(0).Object().Value("rpe_value").Object()
		rpeValue.Value("id").String().IsEqual(rpeValueID)
		rpeValue.Value("value").Number().IsEqual(8)
		rpeValue.Value("label").String().IsEqual("Very Hard+")
	})

	t.Run("Update Prescription Group", func(t *testing.T) {
		// Get RPE value ID for RPE 9 (near max effort)
		rpeValue9ID := GetRPEValueID(t, 9)

		response := e.PUT("/api/v1/workouts/"+workoutID+"/prescriptions/"+groupID).
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"group_name":        "Heavy Bench Press",
				"rest_between_sets": 180,
				"exercises": []map[string]interface{}{
					{
						"exercise_id":    exercise1ID,
						"exercise_order": 1,
						"sets":           5,
						"reps":           5,
						"weight_kg":      80.0,
						"rpe_value_id":   rpeValue9ID,
					},
				},
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Object()
		data.Value("group_name").String().IsEqual("Heavy Bench Press")
		data.Value("rest_between_sets").Number().IsEqual(180)

		exercises := data.Value("exercises").Array()
		exercises.Length().IsEqual(1)
		exercises.Value(0).Object().Value("sets").Number().IsEqual(5)
		exercises.Value(0).Object().Value("reps").Number().IsEqual(5)
		exercises.Value(0).Object().Value("weight_kg").Number().IsEqual(80.0)
		exercises.Value(0).Object().Value("rpe_value_id").String().IsEqual(rpeValue9ID)
	})

	t.Run("Create Second Prescription Group", func(t *testing.T) {
		response := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"type":              "straight",
				"group_order":       2,
				"group_name":        "Incline Press",
				"rest_between_sets": 90,
				"exercises": []map[string]interface{}{
					{
						"exercise_id":    exercise2ID,
						"exercise_order": 1,
						"sets":           3,
						"reps":           10,
						"weight_kg":      20.0,
					},
				},
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
	})

	t.Run("Verify Two Prescription Groups", func(t *testing.T) {
		response := e.GET("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		data := response.Value("data").Array()
		data.Length().IsEqual(2)
	})

	t.Run("Delete Prescription Group", func(t *testing.T) {
		e.DELETE("/api/v1/workouts/"+workoutID+"/prescriptions/"+groupID).
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200)

		// Verify only one group remains
		response := e.GET("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("data").Array().Length().IsEqual(1)
	})

	t.Run("Delete Workout", func(t *testing.T) {
		e.DELETE("/api/v1/workouts/"+workoutID).
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200)

		// Verify it's gone
		e.GET("/api/v1/workouts/"+workoutID).
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(404)
	})
}

func testPrescriptionTypes(t *testing.T, e *httpexpect.Expect) {
	// Create user
	userToken := createTestUserAndGetToken(e, "user@example.com", "UserPass123!", "Test", "User")

	// Create exercises
	exercise1Response := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"name":        "Bench Press",
			"description": "Compound chest exercise",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	exercise1ID := exercise1Response.Value("data").Object().Value("id").String().Raw()

	exercise2Response := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"name":        "Dumbbell Fly",
			"description": "Chest isolation exercise",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	exercise2ID := exercise2Response.Value("data").Object().Value("id").String().Raw()

	exercise3Response := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"name":        "Push Ups",
			"description": "Bodyweight exercise",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	exercise3ID := exercise3Response.Value("data").Object().Value("id").String().Raw()

	// Create workout
	workoutResponse := e.POST("/api/v1/workouts/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"title":       "Test Workout",
			"description": "Testing different prescription types",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	workoutID := workoutResponse.Value("data").Object().Value("id").String().Raw()

	t.Run("Create Superset", func(t *testing.T) {
		response := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"type":              "superset",
				"group_order":       1,
				"group_name":        "Chest Superset",
				"group_rounds":      3,
				"rest_between_sets": 60,
				"exercises": []map[string]interface{}{
					{
						"exercise_id":    exercise1ID,
						"exercise_order": 1,
						"sets":           3,
						"reps":           10,
						"weight_kg":      50.0,
					},
					{
						"exercise_id":    exercise2ID,
						"exercise_order": 2,
						"sets":           3,
						"reps":           12,
						"weight_kg":      10.0,
					},
				},
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Object()
		data.Value("type").String().IsEqual("superset")
		data.Value("group_rounds").Number().IsEqual(3)
		data.Value("exercises").Array().Length().IsEqual(2)
	})

	t.Run("Create Circuit", func(t *testing.T) {
		response := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"type":              "circuit",
				"group_order":       2,
				"group_name":        "Chest Circuit",
				"group_rounds":      4,
				"rest_between_sets": 30,
				"exercises": []map[string]interface{}{
					{
						"exercise_id":    exercise1ID,
						"exercise_order": 1,
						"sets":           1,
						"reps":           10,
					},
					{
						"exercise_id":    exercise2ID,
						"exercise_order": 2,
						"sets":           1,
						"reps":           12,
					},
					{
						"exercise_id":    exercise3ID,
						"exercise_order": 3,
						"sets":           1,
						"reps":           15,
					},
				},
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Object()
		data.Value("type").String().IsEqual("circuit")
		data.Value("exercises").Array().Length().IsEqual(3)
	})

	t.Run("Create Drop Set", func(t *testing.T) {
		response := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"type":        "drop_set",
				"group_order": 3,
				"group_name":  "Bench Drop Set",
				"exercises": []map[string]interface{}{
					{
						"exercise_id":    exercise1ID,
						"exercise_order": 1,
						"reps":           8,
						"weight_kg":      80.0,
					},
					{
						"exercise_id":    exercise1ID,
						"exercise_order": 2,
						"reps":           10,
						"weight_kg":      60.0,
					},
					{
						"exercise_id":    exercise1ID,
						"exercise_order": 3,
						"reps":           12,
						"weight_kg":      40.0,
					},
				},
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Object()
		data.Value("type").String().IsEqual("drop_set")
		data.Value("exercises").Array().Length().IsEqual(3)
	})

	t.Run("Create AMRAP", func(t *testing.T) {
		response := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"type":        "amrap",
				"group_order": 4,
				"group_name":  "AMRAP Finisher",
				"group_notes": "As many rounds as possible in 10 minutes",
				"exercises": []map[string]interface{}{
					{
						"exercise_id":      exercise3ID,
						"exercise_order":   1,
						"duration_seconds": 600,
					},
				},
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Object()
		data.Value("type").String().IsEqual("amrap")
	})

	t.Run("Verify All Prescription Types", func(t *testing.T) {
		response := e.GET("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		data := response.Value("data").Array()
		data.Length().IsEqual(4)
	})
}

func testPrescriptionGroupOperations(t *testing.T, e *httpexpect.Expect) {
	// Create user
	userToken := createTestUserAndGetToken(e, "user@example.com", "UserPass123!", "Test", "User")

	// Create exercise
	exerciseResponse := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"name":        "Squat",
			"description": "Compound leg exercise",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	exerciseID := exerciseResponse.Value("data").Object().Value("id").String().Raw()

	// Create workout
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

	// Create prescription group
	prescriptionResponse := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
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
		Status(201).
		JSON().
		Object()

	groupID := prescriptionResponse.Value("data").Object().Value("group_id").String().Raw()

	t.Run("Update Group Type", func(t *testing.T) {
		response := e.PUT("/api/v1/workouts/"+workoutID+"/prescriptions/"+groupID).
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"type": "pyramid",
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("data").Object().Value("type").String().IsEqual("pyramid")
	})

	t.Run("Update Group Order", func(t *testing.T) {
		response := e.PUT("/api/v1/workouts/"+workoutID+"/prescriptions/"+groupID).
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"group_order": 5,
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("data").Object().Value("group_order").Number().IsEqual(5)
	})

	t.Run("Update Group Notes", func(t *testing.T) {
		response := e.PUT("/api/v1/workouts/"+workoutID+"/prescriptions/"+groupID).
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"group_notes": "Focus on depth",
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("data").Object().Value("group_notes").String().IsEqual("Focus on depth")
	})

	t.Run("Delete Non-Existent Group Returns 404", func(t *testing.T) {
		e.DELETE("/api/v1/workouts/"+workoutID+"/prescriptions/11111111-1111-1111-1111-111111111111").
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(404)
	})

	t.Run("Update Non-Existent Group Returns 404", func(t *testing.T) {
		e.PUT("/api/v1/workouts/"+workoutID+"/prescriptions/11111111-1111-1111-1111-111111111111").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"group_name": "Test",
			}).
			Expect().
			Status(404)
	})
}

func testReorderPrescriptionGroups(t *testing.T, e *httpexpect.Expect) {
	// Create user
	userToken := createTestUserAndGetToken(e, "user@example.com", "UserPass123!", "Test", "User")

	// Create exercise
	exerciseResponse := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"name":        "Test Exercise",
			"description": "Test",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	exerciseID := exerciseResponse.Value("data").Object().Value("id").String().Raw()

	// Create workout
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

	// Create three prescription groups
	group1Response := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"type":        "straight",
			"group_order": 1,
			"group_name":  "Group A",
			"exercises": []map[string]interface{}{
				{"exercise_id": exerciseID, "exercise_order": 1, "sets": 3, "reps": 10},
			},
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	group1ID := group1Response.Value("data").Object().Value("group_id").String().Raw()

	group2Response := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"type":        "straight",
			"group_order": 2,
			"group_name":  "Group B",
			"exercises": []map[string]interface{}{
				{"exercise_id": exerciseID, "exercise_order": 1, "sets": 3, "reps": 10},
			},
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	group2ID := group2Response.Value("data").Object().Value("group_id").String().Raw()

	group3Response := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"type":        "straight",
			"group_order": 3,
			"group_name":  "Group C",
			"exercises": []map[string]interface{}{
				{"exercise_id": exerciseID, "exercise_order": 1, "sets": 3, "reps": 10},
			},
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	group3ID := group3Response.Value("data").Object().Value("group_id").String().Raw()

	t.Run("Reorder Groups", func(t *testing.T) {
		// Swap order: C=1, A=2, B=3
		e.PUT("/api/v1/workouts/"+workoutID+"/prescriptions/reorder").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"group_orders": []map[string]interface{}{
					{"group_id": group3ID, "group_order": 1},
					{"group_id": group1ID, "group_order": 2},
					{"group_id": group2ID, "group_order": 3},
				},
			}).
			Expect().
			Status(200)

		// Verify new order
		response := e.GET("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		data := response.Value("data").Array()
		data.Length().IsEqual(3)

		// Groups should be returned in order
		data.Value(0).Object().Value("group_name").String().IsEqual("Group C")
		data.Value(1).Object().Value("group_name").String().IsEqual("Group A")
		data.Value(2).Object().Value("group_name").String().IsEqual("Group B")
	})
}

func testAddExerciseToPrescriptionGroup(t *testing.T, e *httpexpect.Expect) {
	// Create user
	userToken := createTestUserAndGetToken(e, "user@example.com", "UserPass123!", "Test", "User")

	// Create exercises
	exercise1Response := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"name":        "Bench Press",
			"description": "Chest exercise",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	exercise1ID := exercise1Response.Value("data").Object().Value("id").String().Raw()

	exercise2Response := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"name":        "Incline Press",
			"description": "Upper chest exercise",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	exercise2ID := exercise2Response.Value("data").Object().Value("id").String().Raw()

	// Create workout
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

	// Create prescription group with one exercise
	prescriptionResponse := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"type":        "superset",
			"group_order": 1,
			"exercises": []map[string]interface{}{
				{"exercise_id": exercise1ID, "exercise_order": 1, "sets": 3, "reps": 10},
			},
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	groupID := prescriptionResponse.Value("data").Object().Value("group_id").String().Raw()

	t.Run("Add Exercise To Group", func(t *testing.T) {
		response := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions/"+groupID+"/exercises").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"exercise_id": exercise2ID,
				"sets":        3,
				"reps":        12,
				"weight_kg":   25.0,
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Object()
		data.Value("exercise_order").Number().IsEqual(2) // Auto-assigned order
		data.Value("sets").Number().IsEqual(3)
		data.Value("reps").Number().IsEqual(12)
	})

	t.Run("Verify Group Has Two Exercises", func(t *testing.T) {
		response := e.GET("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		data := response.Value("data").Array()
		data.Length().IsEqual(1)
		data.Value(0).Object().Value("exercises").Array().Length().IsEqual(2)
	})

	t.Run("Add Exercise To Non-Existent Group Returns 404", func(t *testing.T) {
		e.POST("/api/v1/workouts/"+workoutID+"/prescriptions/11111111-1111-1111-1111-111111111111/exercises").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"exercise_id": exercise2ID,
				"sets":        3,
				"reps":        10,
			}).
			Expect().
			Status(404)
	})
}

func testDuplicateWorkoutWithPrescriptions(t *testing.T, e *httpexpect.Expect) {
	// Create user
	userToken := createTestUserAndGetToken(e, "user@example.com", "UserPass123!", "Test", "User")

	// Create exercises
	exercise1Response := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"name":        "Squat",
			"description": "Leg exercise",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	exercise1ID := exercise1Response.Value("data").Object().Value("id").String().Raw()

	exercise2Response := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"name":        "Leg Press",
			"description": "Machine leg exercise",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	exercise2ID := exercise2Response.Value("data").Object().Value("id").String().Raw()

	// Create workout
	workoutResponse := e.POST("/api/v1/workouts/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"title":       "Leg Day",
			"description": "Quad focused",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	workoutID := workoutResponse.Value("data").Object().Value("id").String().Raw()

	// Create prescription groups
	e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"type":        "straight",
			"group_order": 1,
			"group_name":  "Squats",
			"exercises": []map[string]interface{}{
				{"exercise_id": exercise1ID, "exercise_order": 1, "sets": 5, "reps": 5, "weight_kg": 100},
			},
		}).
		Expect().
		Status(201)

	e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"type":        "superset",
			"group_order": 2,
			"group_name":  "Accessory",
			"exercises": []map[string]interface{}{
				{"exercise_id": exercise1ID, "exercise_order": 1, "sets": 3, "reps": 12},
				{"exercise_id": exercise2ID, "exercise_order": 2, "sets": 3, "reps": 12},
			},
		}).
		Expect().
		Status(201)

	t.Run("Duplicate Workout", func(t *testing.T) {
		response := e.POST("/api/v1/workouts/"+workoutID+"/duplicate").
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Object()
		newWorkoutID := data.Value("id").String().Raw()

		// Verify it's a different ID
		if newWorkoutID == workoutID {
			t.Error("Duplicated workout has same ID as original")
		}

		data.Value("title").String().IsEqual("Leg Day (Copy)")
		data.Value("description").String().IsEqual("Quad focused")

		// Verify prescriptions were duplicated
		prescriptions := data.Value("prescriptions").Array()
		prescriptions.Length().IsEqual(2)

		// Check first group
		prescriptions.Value(0).Object().Value("group_name").String().IsEqual("Squats")
		prescriptions.Value(0).Object().Value("exercises").Array().Length().IsEqual(1)

		// Check second group
		prescriptions.Value(1).Object().Value("group_name").String().IsEqual("Accessory")
		prescriptions.Value(1).Object().Value("exercises").Array().Length().IsEqual(2)
	})
}

func testPrescriptionAuthorization(t *testing.T, e *httpexpect.Expect) {
	// Create two users
	user1Token := createTestUserAndGetToken(e, "user1@example.com", "User1Pass123!", "User", "One")
	user2Token := createTestUserAndGetToken(e, "user2@example.com", "User2Pass123!", "User", "Two")

	// User1 creates an exercise
	exerciseResponse := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+user1Token).
		WithJSON(map[string]interface{}{
			"name":        "Test Exercise",
			"description": "Test",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	exerciseID := exerciseResponse.Value("data").Object().Value("id").String().Raw()

	// User1 creates a workout
	workoutResponse := e.POST("/api/v1/workouts/").
		WithHeader("Authorization", "Bearer "+user1Token).
		WithJSON(map[string]interface{}{
			"title": "User1's Workout",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	workoutID := workoutResponse.Value("data").Object().Value("id").String().Raw()

	// User1 creates a prescription group
	prescriptionResponse := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
		WithHeader("Authorization", "Bearer "+user1Token).
		WithJSON(map[string]interface{}{
			"type":        "straight",
			"group_order": 1,
			"exercises": []map[string]interface{}{
				{"exercise_id": exerciseID, "exercise_order": 1, "sets": 3, "reps": 10},
			},
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	groupID := prescriptionResponse.Value("data").Object().Value("group_id").String().Raw()

	t.Run("User Cannot Access Other User's Workout", func(t *testing.T) {
		e.GET("/api/v1/workouts/"+workoutID).
			WithHeader("Authorization", "Bearer "+user2Token).
			Expect().
			Status(404)
	})

	t.Run("User Cannot Create Prescription In Other User's Workout", func(t *testing.T) {
		e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+user2Token).
			WithJSON(map[string]interface{}{
				"type":        "straight",
				"group_order": 2,
				"exercises": []map[string]interface{}{
					{"exercise_id": exerciseID, "exercise_order": 1, "sets": 3, "reps": 10},
				},
			}).
			Expect().
			Status(404)
	})

	t.Run("User Cannot Update Other User's Prescription", func(t *testing.T) {
		e.PUT("/api/v1/workouts/"+workoutID+"/prescriptions/"+groupID).
			WithHeader("Authorization", "Bearer "+user2Token).
			WithJSON(map[string]interface{}{
				"group_name": "Hacked",
			}).
			Expect().
			Status(404)
	})

	t.Run("User Cannot Delete Other User's Prescription", func(t *testing.T) {
		e.DELETE("/api/v1/workouts/"+workoutID+"/prescriptions/"+groupID).
			WithHeader("Authorization", "Bearer "+user2Token).
			Expect().
			Status(404)
	})

	t.Run("User Cannot Duplicate Other User's Workout", func(t *testing.T) {
		e.POST("/api/v1/workouts/"+workoutID+"/duplicate").
			WithHeader("Authorization", "Bearer "+user2Token).
			Expect().
			Status(404)
	})
}

func testPrescriptionValidation(t *testing.T, e *httpexpect.Expect) {
	// Create user
	userToken := createTestUserAndGetToken(e, "user@example.com", "UserPass123!", "Test", "User")

	// Create exercise
	exerciseResponse := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"name":        "Test Exercise",
			"description": "Test",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	exerciseID := exerciseResponse.Value("data").Object().Value("id").String().Raw()

	// Create workout
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

	t.Run("Invalid Prescription Type Returns Error", func(t *testing.T) {
		response := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"type":        "invalid_type",
				"group_order": 1,
				"exercises": []map[string]interface{}{
					{"exercise_id": exerciseID, "exercise_order": 1, "sets": 3, "reps": 10},
				},
			}).
			Expect().
			Status(400).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("Missing Required Type Returns Error", func(t *testing.T) {
		response := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"group_order": 1,
				"exercises": []map[string]interface{}{
					{"exercise_id": exerciseID, "exercise_order": 1, "sets": 3, "reps": 10},
				},
			}).
			Expect().
			Status(400).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("Missing Required Group Order Returns Error", func(t *testing.T) {
		response := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"type": "straight",
				"exercises": []map[string]interface{}{
					{"exercise_id": exerciseID, "exercise_order": 1, "sets": 3, "reps": 10},
				},
			}).
			Expect().
			Status(400).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("Empty Exercises Array Returns Error", func(t *testing.T) {
		response := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"type":        "straight",
				"group_order": 1,
				"exercises":   []map[string]interface{}{},
			}).
			Expect().
			Status(400).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("Invalid Exercise ID Returns Error", func(t *testing.T) {
		response := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"type":        "straight",
				"group_order": 1,
				"exercises": []map[string]interface{}{
					{"exercise_id": "11111111-1111-1111-1111-111111111111", "exercise_order": 1, "sets": 3, "reps": 10},
				},
			}).
			Expect().
			Status(400).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("Invalid Workout ID Returns Error", func(t *testing.T) {
		e.POST("/api/v1/workouts/invalid-uuid/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"type":        "straight",
				"group_order": 1,
				"exercises": []map[string]interface{}{
					{"exercise_id": exerciseID, "exercise_order": 1, "sets": 3, "reps": 10},
				},
			}).
			Expect().
			Status(400)
	})

	t.Run("Non-Existent Workout Returns 404", func(t *testing.T) {
		e.POST("/api/v1/workouts/11111111-1111-1111-1111-111111111111/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"type":        "straight",
				"group_order": 1,
				"exercises": []map[string]interface{}{
					{"exercise_id": exerciseID, "exercise_order": 1, "sets": 3, "reps": 10},
				},
			}).
			Expect().
			Status(404)
	})
}

