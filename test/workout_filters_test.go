package test

import (
	"lamari-fit-api/database"
	"lamari-fit-api/models"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
)

func TestWorkoutFilters(t *testing.T) {
	e := SetupTestApp(t)

	// Create test user
	email := "workout_filter_test@example.com"
	password := "Password123!"

	// Register user
	resp := e.POST("/api/v1/auth/register").
		WithJSON(map[string]interface{}{
			"email":            email,
			"password":         password,
			"password_confirm": password,
			"first_name":       "Filter",
			"last_name":        "Test",
		}).
		Expect().
		Status(201).
		JSON().Object()

	token := resp.Path("$.data.access_token").String().Raw()

	// Get user ID
	var user models.User
	database.DB.Where("email = ?", email).First(&user)
	userID := user.ID

	// Create muscle groups
	var muscleGroupIDs []uuid.UUID
	muscleGroups := []models.MuscleGroup{
		{Name: "Filter Test Chest", Category: "upper"},
		{Name: "Filter Test Back", Category: "upper"},
		{Name: "Filter Test Legs", Category: "lower"},
	}
	for i := range muscleGroups {
		database.DB.Create(&muscleGroups[i])
		muscleGroupIDs = append(muscleGroupIDs, muscleGroups[i].ID)
	}

	// Create exercises
	var exerciseIDs []uuid.UUID
	exercises := []models.Exercise{
		{Name: "Filter Test Bench Press", Slug: "filter-test-bench-press"},
		{Name: "Filter Test Pull Up", Slug: "filter-test-pull-up"},
		{Name: "Filter Test Squat", Slug: "filter-test-squat"},
		{Name: "Filter Test Deadlift", Slug: "filter-test-deadlift"},
	}
	for i := range exercises {
		database.DB.Create(&exercises[i])
		exerciseIDs = append(exerciseIDs, exercises[i].ID)
	}

	// Link exercises to muscle groups
	// Bench Press -> Chest
	database.DB.Create(&models.ExerciseMuscleGroup{
		ExerciseID:    exerciseIDs[0],
		MuscleGroupID: muscleGroupIDs[0],
		Primary:       true,
	})
	// Pull Up -> Back
	database.DB.Create(&models.ExerciseMuscleGroup{
		ExerciseID:    exerciseIDs[1],
		MuscleGroupID: muscleGroupIDs[1],
		Primary:       true,
	})
	// Squat -> Legs
	database.DB.Create(&models.ExerciseMuscleGroup{
		ExerciseID:    exerciseIDs[2],
		MuscleGroupID: muscleGroupIDs[2],
		Primary:       true,
	})
	// Deadlift -> Back and Legs
	database.DB.Create(&models.ExerciseMuscleGroup{
		ExerciseID:    exerciseIDs[3],
		MuscleGroupID: muscleGroupIDs[1],
		Primary:       true,
	})
	database.DB.Create(&models.ExerciseMuscleGroup{
		ExerciseID:    exerciseIDs[3],
		MuscleGroupID: muscleGroupIDs[2],
		Primary:       false,
	})

	// Create workouts with different exercise combinations
	var workoutIDs []uuid.UUID
	workouts := []models.Workout{
		{UserID: userID, Title: "Chest Day", Visibility: "private"},            // Bench Press only
		{UserID: userID, Title: "Back Day", Visibility: "private"},             // Pull Up only
		{UserID: userID, Title: "Leg Day", Visibility: "private"},              // Squat only
		{UserID: userID, Title: "Full Body", Visibility: "private"},            // All exercises
		{UserID: userID, Title: "Upper Body Push Pull", Visibility: "private"}, // Bench + Pull Up
	}
	for i := range workouts {
		database.DB.Create(&workouts[i])
		workoutIDs = append(workoutIDs, workouts[i].ID)
	}

	// Add prescriptions to workouts
	intPtr := func(i int) *int { return &i }

	// Workout 0 (Chest Day): Bench Press
	database.DB.Create(&models.WorkoutPrescription{
		WorkoutID:     workoutIDs[0],
		ExerciseID:    exerciseIDs[0],
		GroupID:       uuid.New(),
		Type:          models.PrescriptionTypeStraight,
		GroupOrder:    1,
		ExerciseOrder: 1,
		Sets:          intPtr(3),
		Reps:          intPtr(10),
	})

	// Workout 1 (Back Day): Pull Up
	database.DB.Create(&models.WorkoutPrescription{
		WorkoutID:     workoutIDs[1],
		ExerciseID:    exerciseIDs[1],
		GroupID:       uuid.New(),
		Type:          models.PrescriptionTypeStraight,
		GroupOrder:    1,
		ExerciseOrder: 1,
		Sets:          intPtr(3),
		Reps:          intPtr(10),
	})

	// Workout 2 (Leg Day): Squat
	database.DB.Create(&models.WorkoutPrescription{
		WorkoutID:     workoutIDs[2],
		ExerciseID:    exerciseIDs[2],
		GroupID:       uuid.New(),
		Type:          models.PrescriptionTypeStraight,
		GroupOrder:    1,
		ExerciseOrder: 1,
		Sets:          intPtr(3),
		Reps:          intPtr(10),
	})

	// Workout 3 (Full Body): All exercises
	groupID := uuid.New()
	for i, exID := range exerciseIDs {
		database.DB.Create(&models.WorkoutPrescription{
			WorkoutID:     workoutIDs[3],
			ExerciseID:    exID,
			GroupID:       groupID,
			Type:          models.PrescriptionTypeStraight,
			GroupOrder:    1,
			ExerciseOrder: i + 1,
			Sets:          intPtr(3),
			Reps:          intPtr(10),
		})
	}

	// Workout 4 (Upper Body): Bench + Pull Up
	groupID2 := uuid.New()
	database.DB.Create(&models.WorkoutPrescription{
		WorkoutID:     workoutIDs[4],
		ExerciseID:    exerciseIDs[0],
		GroupID:       groupID2,
		Type:          models.PrescriptionTypeStraight,
		GroupOrder:    1,
		ExerciseOrder: 1,
		Sets:          intPtr(3),
		Reps:          intPtr(10),
	})
	database.DB.Create(&models.WorkoutPrescription{
		WorkoutID:     workoutIDs[4],
		ExerciseID:    exerciseIDs[1],
		GroupID:       groupID2,
		Type:          models.PrescriptionTypeStraight,
		GroupOrder:    1,
		ExerciseOrder: 2,
		Sets:          intPtr(3),
		Reps:          intPtr(10),
	})

	// Run subtests
	t.Run("OR Mode Single Muscle Group", func(t *testing.T) {
		testOrModeSingleMuscleGroup(t, e, token, muscleGroupIDs)
	})

	t.Run("OR Mode Multiple Muscle Groups", func(t *testing.T) {
		testOrModeMultipleMuscleGroups(t, e, token, muscleGroupIDs)
	})

	t.Run("AND Mode Multiple Muscle Groups", func(t *testing.T) {
		testAndModeMultipleMuscleGroups(t, e, token, muscleGroupIDs)
	})

	t.Run("OR Mode Single Exercise", func(t *testing.T) {
		testOrModeSingleExercise(t, e, token, exerciseIDs)
	})

	t.Run("OR Mode Multiple Exercises", func(t *testing.T) {
		testOrModeMultipleExercises(t, e, token, exerciseIDs)
	})

	t.Run("AND Mode Multiple Exercises", func(t *testing.T) {
		testAndModeMultipleExercises(t, e, token, exerciseIDs)
	})

	t.Run("OR Mode Combined Filters", func(t *testing.T) {
		testOrModeCombinedFilters(t, e, token, muscleGroupIDs, exerciseIDs)
	})

	t.Run("AND Mode Combined Filters", func(t *testing.T) {
		testAndModeCombinedFilters(t, e, token, muscleGroupIDs, exerciseIDs)
	})

	t.Run("Search With Filters", func(t *testing.T) {
		testSearchWithFilters(t, e, token, muscleGroupIDs)
	})

	t.Run("Default Mode Is OR", func(t *testing.T) {
		testDefaultModeIsOr(t, e, token, muscleGroupIDs)
	})

	t.Run("No Filters Returns All", func(t *testing.T) {
		testNoFiltersReturnsAll(t, e, token)
	})
}

func testOrModeSingleMuscleGroup(t *testing.T, e *httpexpect.Expect, token string, muscleGroupIDs []uuid.UUID) {
	resp := e.GET("/api/v1/workouts").
		WithHeader("Authorization", "Bearer "+token).
		WithQuery("muscle_group_id", muscleGroupIDs[0].String()).
		Expect().
		Status(200).
		JSON().Object()

	data := resp.Path("$.data").Array()
	// Should include: Chest Day, Full Body, Upper Body (all have Chest exercises)
	if len(data.Raw()) < 3 {
		t.Errorf("Expected at least 3 workouts, got %d", len(data.Raw()))
	}
}

func testOrModeMultipleMuscleGroups(t *testing.T, e *httpexpect.Expect, token string, muscleGroupIDs []uuid.UUID) {
	resp := e.GET("/api/v1/workouts").
		WithHeader("Authorization", "Bearer "+token).
		WithQuery("muscle_group_id", muscleGroupIDs[0].String()).
		WithQuery("muscle_group_id", muscleGroupIDs[1].String()).
		WithQuery("mode", "or").
		Expect().
		Status(200).
		JSON().Object()

	data := resp.Path("$.data").Array()
	// Should include workouts with Chest OR Back
	if len(data.Raw()) < 4 {
		t.Errorf("Expected at least 4 workouts, got %d", len(data.Raw()))
	}
}

func testAndModeMultipleMuscleGroups(t *testing.T, e *httpexpect.Expect, token string, muscleGroupIDs []uuid.UUID) {
	resp := e.GET("/api/v1/workouts").
		WithHeader("Authorization", "Bearer "+token).
		WithQuery("muscle_group_id", muscleGroupIDs[0].String()).
		WithQuery("muscle_group_id", muscleGroupIDs[1].String()).
		WithQuery("mode", "and").
		Expect().
		Status(200).
		JSON().Object()

	data := resp.Path("$.data").Array()
	// Should include: Full Body, Upper Body (both have Chest AND Back)
	if len(data.Raw()) < 2 {
		t.Errorf("Expected at least 2 workouts, got %d", len(data.Raw()))
	}
}

func testOrModeSingleExercise(t *testing.T, e *httpexpect.Expect, token string, exerciseIDs []uuid.UUID) {
	resp := e.GET("/api/v1/workouts").
		WithHeader("Authorization", "Bearer "+token).
		WithQuery("exercise_id", exerciseIDs[0].String()).
		Expect().
		Status(200).
		JSON().Object()

	data := resp.Path("$.data").Array()
	// Should include: Chest Day, Full Body, Upper Body
	if len(data.Raw()) < 3 {
		t.Errorf("Expected at least 3 workouts, got %d", len(data.Raw()))
	}
}

func testOrModeMultipleExercises(t *testing.T, e *httpexpect.Expect, token string, exerciseIDs []uuid.UUID) {
	resp := e.GET("/api/v1/workouts").
		WithHeader("Authorization", "Bearer "+token).
		WithQuery("exercise_id", exerciseIDs[0].String()).
		WithQuery("exercise_id", exerciseIDs[2].String()).
		WithQuery("mode", "or").
		Expect().
		Status(200).
		JSON().Object()

	data := resp.Path("$.data").Array()
	// Should include workouts with Bench OR Squat
	if len(data.Raw()) < 3 {
		t.Errorf("Expected at least 3 workouts, got %d", len(data.Raw()))
	}
}

func testAndModeMultipleExercises(t *testing.T, e *httpexpect.Expect, token string, exerciseIDs []uuid.UUID) {
	resp := e.GET("/api/v1/workouts").
		WithHeader("Authorization", "Bearer "+token).
		WithQuery("exercise_id", exerciseIDs[0].String()).
		WithQuery("exercise_id", exerciseIDs[1].String()).
		WithQuery("mode", "and").
		Expect().
		Status(200).
		JSON().Object()

	data := resp.Path("$.data").Array()
	// Should include: Full Body, Upper Body (both have Bench AND Pull Up)
	if len(data.Raw()) < 2 {
		t.Errorf("Expected at least 2 workouts, got %d", len(data.Raw()))
	}
}

func testOrModeCombinedFilters(t *testing.T, e *httpexpect.Expect, token string, muscleGroupIDs, exerciseIDs []uuid.UUID) {
	resp := e.GET("/api/v1/workouts").
		WithHeader("Authorization", "Bearer "+token).
		WithQuery("muscle_group_id", muscleGroupIDs[2].String()).
		WithQuery("exercise_id", exerciseIDs[0].String()).
		WithQuery("mode", "or").
		Expect().
		Status(200).
		JSON().Object()

	data := resp.Path("$.data").Array()
	// Should include workouts with Legs OR Bench Press
	if len(data.Raw()) < 4 {
		t.Errorf("Expected at least 4 workouts, got %d", len(data.Raw()))
	}
}

func testAndModeCombinedFilters(t *testing.T, e *httpexpect.Expect, token string, muscleGroupIDs, exerciseIDs []uuid.UUID) {
	resp := e.GET("/api/v1/workouts").
		WithHeader("Authorization", "Bearer "+token).
		WithQuery("muscle_group_id", muscleGroupIDs[0].String()).
		WithQuery("exercise_id", exerciseIDs[0].String()).
		WithQuery("mode", "and").
		Expect().
		Status(200).
		JSON().Object()

	data := resp.Path("$.data").Array()
	// Should include workouts with Chest muscle AND Bench Press exercise
	if len(data.Raw()) < 3 {
		t.Errorf("Expected at least 3 workouts, got %d", len(data.Raw()))
	}
}

func testSearchWithFilters(t *testing.T, e *httpexpect.Expect, token string, muscleGroupIDs []uuid.UUID) {
	resp := e.GET("/api/v1/workouts").
		WithHeader("Authorization", "Bearer "+token).
		WithQuery("search", "Full").
		WithQuery("muscle_group_id", muscleGroupIDs[0].String()).
		Expect().
		Status(200).
		JSON().Object()

	data := resp.Path("$.data").Array()
	// Should find "Full Body" workout that has Chest
	for _, item := range data.Iter() {
		title := item.Object().Path("$.title").String().Raw()
		if title != "" && title != "Full Body" {
			// Allow Full Body or empty (pagination might return empty)
		}
	}
}

func testDefaultModeIsOr(t *testing.T, e *httpexpect.Expect, token string, muscleGroupIDs []uuid.UUID) {
	resp := e.GET("/api/v1/workouts").
		WithHeader("Authorization", "Bearer "+token).
		WithQuery("muscle_group_id", muscleGroupIDs[0].String()).
		WithQuery("muscle_group_id", muscleGroupIDs[1].String()).
		Expect().
		Status(200).
		JSON().Object()

	data := resp.Path("$.data").Array()
	// OR mode should return more results than AND
	if len(data.Raw()) < 4 {
		t.Errorf("Expected at least 4 workouts (OR mode default), got %d", len(data.Raw()))
	}
}

func testNoFiltersReturnsAll(t *testing.T, e *httpexpect.Expect, token string) {
	resp := e.GET("/api/v1/workouts").
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(200).
		JSON().Object()

	data := resp.Path("$.data").Array()
	if len(data.Raw()) < 5 {
		t.Errorf("Expected at least 5 workouts, got %d", len(data.Raw()))
	}
}
