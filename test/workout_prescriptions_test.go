package test

import (
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestWorkoutPrescriptionEndpoints(t *testing.T) {
	e := SetupTestApp(t)

	t.Run("Workout Fields", func(t *testing.T) {
		CleanDatabase(t)
		testWorkoutFields(t, e)
	})

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

	t.Run("Isometric Hold Support", func(t *testing.T) {
		CleanDatabase(t)
		testIsometricHoldSupport(t, e)
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
						"target_weight": map[string]interface{}{
							"weight_value": 60.0,
							"weight_unit":  "kg",
						},
						"rpe_value_id": rpeValueID,
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
		// Check target_weight in WeightOutput format
		targetWeight := exercises.Value(0).Object().Value("target_weight").Object()
		targetWeight.Value("weight_value").Number().IsEqual(60.0)
		targetWeight.Value("weight_unit").String().IsEqual("kg")
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
						"target_weight": map[string]interface{}{
							"weight_value": 80.0,
							"weight_unit":  "kg",
						},
						"rpe_value_id": rpeValue9ID,
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
		targetWeight := exercises.Value(0).Object().Value("target_weight").Object()
		targetWeight.Value("weight_value").Number().IsEqual(80.0)
		targetWeight.Value("weight_unit").String().IsEqual("kg")
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
						"target_weight":  map[string]interface{}{"weight_value": 20.0, "weight_unit": "kg"},
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
						"target_weight":  map[string]interface{}{"weight_value": 50.0, "weight_unit": "kg"},
					},
					{
						"exercise_id":    exercise2ID,
						"exercise_order": 2,
						"sets":           3,
						"reps":           12,
						"target_weight":  map[string]interface{}{"weight_value": 10.0, "weight_unit": "kg"},
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
						"target_weight":  map[string]interface{}{"weight_value": 80.0, "weight_unit": "kg"},
					},
					{
						"exercise_id":    exercise1ID,
						"exercise_order": 2,
						"reps":           10,
						"target_weight":  map[string]interface{}{"weight_value": 60.0, "weight_unit": "kg"},
					},
					{
						"exercise_id":    exercise1ID,
						"exercise_order": 3,
						"reps":           12,
						"target_weight":  map[string]interface{}{"weight_value": 40.0, "weight_unit": "kg"},
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
						"exercise_id":    exercise3ID,
						"exercise_order": 1,
						"reps":           10,
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
				"exercise_id":   exercise2ID,
				"sets":          3,
				"reps":          12,
				"target_weight": map[string]interface{}{"weight_value": 25.0, "weight_unit": "kg"},
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

func testIsometricHoldSupport(t *testing.T, e *httpexpect.Expect) {
	// Create user
	userToken := createTestUserAndGetToken(e, "user@example.com", "UserPass123!", "Test", "User")

	// Create exercises - one for reps, one for isometric holds
	exercise1Response := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"name":        "Plank",
			"description": "Core isometric exercise",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	plankExerciseID := exercise1Response.Value("data").Object().Value("id").String().Raw()

	exercise2Response := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"name":        "Wall Sit",
			"description": "Leg isometric exercise",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	wallSitExerciseID := exercise2Response.Value("data").Object().Value("id").String().Raw()

	exercise3Response := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"name":        "Squat",
			"description": "Compound leg exercise with reps",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	squatExerciseID := exercise3Response.Value("data").Object().Value("id").String().Raw()

	// Create workout
	workoutResponse := e.POST("/api/v1/workouts/").
		WithHeader("Authorization", "Bearer "+userToken).
		WithJSON(map[string]interface{}{
			"title":       "Core and Legs",
			"description": "Mixed workout with isometric holds",
			"visibility":  "private",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	workoutID := workoutResponse.Value("data").Object().Value("id").String().Raw()

	t.Run("Create Prescription With Hold Seconds", func(t *testing.T) {
		response := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"type":        "straight",
				"group_order": 1,
				"exercises": []map[string]interface{}{
					{
						"exercise_id":    plankExerciseID,
						"exercise_order": 1,
						"sets":           3,
						"hold_seconds":   60,
					},
				},
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("successfully")

		// Verify hold_seconds is returned in response
		data := response.Value("data").Object()
		exercises := data.Value("exercises").Array()
		exercises.Length().IsEqual(1)
		exercises.Value(0).Object().Value("hold_seconds").Number().IsEqual(60)
		exercises.Value(0).Object().Value("sets").Number().IsEqual(3)
	})

	t.Run("Create Prescription With Reps", func(t *testing.T) {
		response := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"type":        "straight",
				"group_order": 2,
				"exercises": []map[string]interface{}{
					{
						"exercise_id":    squatExerciseID,
						"exercise_order": 1,
						"sets":           4,
						"reps":           12,
					},
				},
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		// Verify reps is returned in response
		data := response.Value("data").Object()
		exercises := data.Value("exercises").Array()
		exercises.Length().IsEqual(1)
		exercises.Value(0).Object().Value("reps").Number().IsEqual(12)
		exercises.Value(0).Object().Value("sets").Number().IsEqual(4)
	})

	t.Run("Create Mixed Prescription Group", func(t *testing.T) {
		// A circuit with both rep-based and time-based exercises
		response := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"type":         "circuit",
				"group_order":  3,
				"group_rounds": 3,
				"group_name":   "Core Circuit",
				"exercises": []map[string]interface{}{
					{
						"exercise_id":    plankExerciseID,
						"exercise_order": 1,
						"sets":           1,
						"hold_seconds":   30,
					},
					{
						"exercise_id":    wallSitExerciseID,
						"exercise_order": 2,
						"sets":           1,
						"hold_seconds":   45,
					},
					{
						"exercise_id":    squatExerciseID,
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

		// Verify mixed exercises
		data := response.Value("data").Object()
		exercises := data.Value("exercises").Array()
		exercises.Length().IsEqual(3)

		// First exercise - plank with hold_seconds
		exercises.Value(0).Object().Value("hold_seconds").Number().IsEqual(30)

		// Second exercise - wall sit with hold_seconds
		exercises.Value(1).Object().Value("hold_seconds").Number().IsEqual(45)

		// Third exercise - squat with reps
		exercises.Value(2).Object().Value("reps").Number().IsEqual(15)
	})

	t.Run("Reject Both Reps And Hold Seconds", func(t *testing.T) {
		response := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"type":        "straight",
				"group_order": 4,
				"exercises": []map[string]interface{}{
					{
						"exercise_id":    plankExerciseID,
						"exercise_order": 1,
						"sets":           3,
						"reps":           10,
						"hold_seconds":   30,
					},
				},
			}).
			Expect().
			Status(400).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("Reject Neither Reps Nor Hold Seconds", func(t *testing.T) {
		response := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"type":        "straight",
				"group_order": 4,
				"exercises": []map[string]interface{}{
					{
						"exercise_id":    plankExerciseID,
						"exercise_order": 1,
						"sets":           3,
					},
				},
			}).
			Expect().
			Status(400).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("Get Workout With Hold Seconds Prescriptions", func(t *testing.T) {
		response := e.GET("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		// Verify prescription groups
		groups := response.Value("data").Array()
		groups.Length().IsEqual(3)

		// First group - plank with hold_seconds
		firstGroup := groups.Value(0).Object()
		firstGroup.Value("type").String().IsEqual("straight")
		firstExercises := firstGroup.Value("exercises").Array()
		firstExercises.Value(0).Object().Value("hold_seconds").Number().IsEqual(60)

		// Second group - squat with reps
		secondGroup := groups.Value(1).Object()
		secondExercises := secondGroup.Value("exercises").Array()
		secondExercises.Value(0).Object().Value("reps").Number().IsEqual(12)

		// Third group - circuit with mixed exercises
		thirdGroup := groups.Value(2).Object()
		thirdGroup.Value("type").String().IsEqual("circuit")
		thirdExercises := thirdGroup.Value("exercises").Array()
		thirdExercises.Length().IsEqual(3)
	})

	t.Run("Hold Seconds With Weight", func(t *testing.T) {
		// Weighted plank or similar exercise
		response := e.POST("/api/v1/workouts/"+workoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"type":        "straight",
				"group_order": 4,
				"exercises": []map[string]interface{}{
					{
						"exercise_id":    plankExerciseID,
						"exercise_order": 1,
						"sets":           3,
						"hold_seconds":   45,
						"target_weight":  map[string]interface{}{"weight_value": 10.0, "weight_unit": "kg"},
					},
				},
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		// Verify both hold_seconds and target_weight are returned
		data := response.Value("data").Object()
		exercises := data.Value("exercises").Array()
		exercises.Value(0).Object().Value("hold_seconds").Number().IsEqual(45)
		targetWeight := exercises.Value(0).Object().Value("target_weight").Object()
		targetWeight.Value("weight_value").Number().IsEqual(10.0)
		targetWeight.Value("weight_unit").String().IsEqual("kg")
	})

	t.Run("Duplicate Workout Preserves Hold Seconds", func(t *testing.T) {
		// Duplicate the workout
		response := e.POST("/api/v1/workouts/"+workoutID+"/duplicate").
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		// Get the duplicated workout ID
		duplicatedWorkoutID := response.Value("data").Object().Value("id").String().Raw()

		// Get prescriptions from duplicated workout
		prescriptionsResponse := e.GET("/api/v1/workouts/"+duplicatedWorkoutID+"/prescriptions").
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		// Verify hold_seconds preserved
		groups := prescriptionsResponse.Value("data").Array()
		groups.Length().IsEqual(4)

		// Check first group has hold_seconds preserved
		firstGroupExercises := groups.Value(0).Object().Value("exercises").Array()
		firstGroupExercises.Value(0).Object().Value("hold_seconds").Number().IsEqual(60)
	})
}

func testWorkoutFields(t *testing.T, e *httpexpect.Expect) {
	// Create user
	userToken := createTestUserAndGetToken(e, "user@example.com", "UserPass123!", "Test", "User")

	var workoutID string

	t.Run("Create Workout With All Fields", func(t *testing.T) {
		response := e.POST("/api/v1/workouts/").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"title":              "Full Body Workout",
				"description":        "A complete full body workout routine",
				"difficulty_level":   "intermediate",
				"estimated_duration": 45,
				"is_template":        true,
				"visibility":         "public",
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("created")

		data := response.Value("data").Object()
		workoutID = data.Value("id").String().Raw()
		data.Value("title").String().IsEqual("Full Body Workout")
		data.Value("description").String().IsEqual("A complete full body workout routine")
		data.Value("difficulty_level").String().IsEqual("intermediate")
		data.Value("estimated_duration").Number().IsEqual(45)
		data.Value("is_template").Boolean().IsTrue()
		data.Value("visibility").String().IsEqual("public")
	})

	t.Run("Get Workout Returns All Fields", func(t *testing.T) {
		response := e.GET("/api/v1/workouts/"+workoutID).
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Object()
		data.Value("id").String().IsEqual(workoutID)
		data.Value("title").String().IsEqual("Full Body Workout")
		data.Value("difficulty_level").String().IsEqual("intermediate")
		data.Value("estimated_duration").Number().IsEqual(45)
		data.Value("is_template").Boolean().IsTrue()
		data.Value("visibility").String().IsEqual("public")
	})

	t.Run("Update Workout Fields", func(t *testing.T) {
		response := e.PUT("/api/v1/workouts/"+workoutID).
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"difficulty_level":   "advanced",
				"estimated_duration": 60,
				"is_template":        false,
				"visibility":         "private",
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Object()
		data.Value("difficulty_level").String().IsEqual("advanced")
		data.Value("estimated_duration").Number().IsEqual(60)
		data.Value("is_template").Boolean().IsFalse()
		data.Value("visibility").String().IsEqual("private")
	})

	t.Run("Create Workout With Minimal Fields", func(t *testing.T) {
		response := e.POST("/api/v1/workouts/").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"title": "Quick Workout",
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Object()
		data.Value("title").String().IsEqual("Quick Workout")
		data.Value("visibility").String().IsEqual("private") // default value
		data.Value("is_template").Boolean().IsFalse()        // default value
	})

	t.Run("Invalid Difficulty Level Returns Error", func(t *testing.T) {
		response := e.POST("/api/v1/workouts/").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"title":            "Test Workout",
				"difficulty_level": "invalid_level",
			}).
			Expect().
			Status(400).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("Invalid Visibility Returns Error", func(t *testing.T) {
		response := e.POST("/api/v1/workouts/").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"title":      "Test Workout",
				"visibility": "invalid_visibility",
			}).
			Expect().
			Status(400).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("Estimated Duration Too High Returns Error", func(t *testing.T) {
		response := e.POST("/api/v1/workouts/").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"title":              "Test Workout",
				"estimated_duration": 700, // max is 600
			}).
			Expect().
			Status(400).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("Duplicate Workout Copies All Fields", func(t *testing.T) {
		// First create a workout with all fields
		createResponse := e.POST("/api/v1/workouts/").
			WithHeader("Authorization", "Bearer "+userToken).
			WithJSON(map[string]interface{}{
				"title":              "Original Workout",
				"description":        "Original description",
				"difficulty_level":   "beginner",
				"estimated_duration": 30,
				"is_template":        true,
				"visibility":         "friends",
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		originalID := createResponse.Value("data").Object().Value("id").String().Raw()

		// Duplicate the workout
		duplicateResponse := e.POST("/api/v1/workouts/"+originalID+"/duplicate").
			WithHeader("Authorization", "Bearer "+userToken).
			Expect().
			Status(201).
			JSON().
			Object()

		duplicateResponse.Value("success").Boolean().IsTrue()

		data := duplicateResponse.Value("data").Object()
		data.Value("title").String().IsEqual("Original Workout (Copy)")
		data.Value("description").String().IsEqual("Original description")
		data.Value("difficulty_level").String().IsEqual("beginner")
		data.Value("estimated_duration").Number().IsEqual(30)
		data.Value("is_template").Boolean().IsTrue()
		data.Value("visibility").String().IsEqual("friends")
	})

	t.Run("Create Workout With Each Difficulty Level", func(t *testing.T) {
		levels := []string{"beginner", "intermediate", "advanced"}

		for _, level := range levels {
			response := e.POST("/api/v1/workouts/").
				WithHeader("Authorization", "Bearer "+userToken).
				WithJSON(map[string]interface{}{
					"title":            level + " Workout",
					"difficulty_level": level,
				}).
				Expect().
				Status(201).
				JSON().
				Object()

			response.Value("success").Boolean().IsTrue()
			response.Value("data").Object().Value("difficulty_level").String().IsEqual(level)
		}
	})
}
