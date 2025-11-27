package test

import (
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestFavoritesEndpoints(t *testing.T) {
	e := SetupTestApp(t)

	t.Run("Favorites", func(t *testing.T) {
		CleanDatabase(t)
		SeedTestRoles(t)
		testFavorites(t, e)
	})
}

func testFavorites(t *testing.T, e *httpexpect.Expect) {
	// Setup: Create a user
	userData := map[string]interface{}{
		"email":            "favorites@example.com",
		"password":         "FavoritesPass123!",
		"password_confirm": "FavoritesPass123!",
		"first_name":       "Favorites",
		"last_name":        "User",
	}

	e.POST("/api/v1/auth/register").
		WithJSON(userData).
		Expect().
		Status(201)

	token := GetAuthToken(e, "favorites@example.com", "FavoritesPass123!")

	// Create a muscle group for exercises
	muscleGroupResp := e.POST("/api/v1/muscle-groups/").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]interface{}{
			"name":        "Test Muscle",
			"name_slug":   "test-muscle",
			"description": "Test muscle group",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	muscleGroupID := muscleGroupResp.Value("data").Object().Value("id").String().Raw()

	// Create an exercise for testing
	exerciseResp := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]interface{}{
			"name":              "Test Exercise",
			"name_slug":         "test-exercise",
			"description":       "A test exercise for favorites",
			"primary_muscle_id": muscleGroupID,
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	exerciseID := exerciseResp.Value("data").Object().Value("id").String().Raw()

	// Create a second exercise for testing
	exercise2Resp := e.POST("/api/v1/exercises/").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]interface{}{
			"name":              "Test Exercise 2",
			"name_slug":         "test-exercise-2",
			"description":       "Another test exercise",
			"primary_muscle_id": muscleGroupID,
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	exercise2ID := exercise2Resp.Value("data").Object().Value("id").String().Raw()

	// Create a workout for testing
	workoutResp := e.POST("/api/v1/workouts/").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]interface{}{
			"title":       "Test Workout",
			"description": "A test workout for favorites",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	workoutID := workoutResp.Value("data").Object().Value("id").String().Raw()

	// Create a second workout for testing
	workout2Resp := e.POST("/api/v1/workouts/").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]interface{}{
			"title":       "Test Workout 2",
			"description": "Another test workout",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	workout2ID := workout2Resp.Value("data").Object().Value("id").String().Raw()

	// =====================
	// Exercise Favorites Tests
	// =====================

	t.Run("Get Empty Exercise Favorites", func(t *testing.T) {
		response := e.GET("/api/v1/user/favorites").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "exercise").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("data").Array().Length().IsEqual(0)
		response.Value("meta").Object().Value("total_items").Number().IsEqual(0)
	})

	t.Run("Add Exercise to Favorites", func(t *testing.T) {
		response := e.POST("/api/v1/user/favorites").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "exercise").
			WithJSON(map[string]interface{}{
				"item_id": exerciseID,
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("added")

		data := response.Value("data").Object()
		data.Value("id").String().NotEmpty()
		data.Value("type").String().IsEqual("exercise")
		data.Value("item_id").String().IsEqual(exerciseID)
		data.Value("item").Object().Value("name").String().IsEqual("Test Exercise")
	})

	t.Run("Check Exercise is Favorited", func(t *testing.T) {
		response := e.GET("/api/v1/user/favorites/" + exerciseID + "/check").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "exercise").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("data").Object().Value("is_favorited").Boolean().IsTrue()
	})

	t.Run("Check Non-Favorited Exercise", func(t *testing.T) {
		response := e.GET("/api/v1/user/favorites/" + exercise2ID + "/check").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "exercise").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("data").Object().Value("is_favorited").Boolean().IsFalse()
	})

	t.Run("Add Duplicate Exercise Favorite", func(t *testing.T) {
		response := e.POST("/api/v1/user/favorites").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "exercise").
			WithJSON(map[string]interface{}{
				"item_id": exerciseID,
			}).
			Expect().
			Status(409).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("already")
	})

	t.Run("Add Second Exercise to Favorites", func(t *testing.T) {
		e.POST("/api/v1/user/favorites").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "exercise").
			WithJSON(map[string]interface{}{
				"item_id": exercise2ID,
			}).
			Expect().
			Status(201)
	})

	t.Run("Get Exercise Favorites with Pagination", func(t *testing.T) {
		response := e.GET("/api/v1/user/favorites").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "exercise").
			WithQuery("page", 1).
			WithQuery("limit", 10).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("data").Array().Length().IsEqual(2)
		response.Value("meta").Object().Value("total_items").Number().IsEqual(2)
		response.Value("meta").Object().Value("current_page").Number().IsEqual(1)
		response.Value("meta").Object().Value("per_page").Number().IsEqual(10)
	})

	t.Run("Get Exercise Favorites with Limit 1", func(t *testing.T) {
		response := e.GET("/api/v1/user/favorites").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "exercise").
			WithQuery("page", 1).
			WithQuery("limit", 1).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("data").Array().Length().IsEqual(1)
		response.Value("meta").Object().Value("total_items").Number().IsEqual(2)
		response.Value("meta").Object().Value("total_pages").Number().IsEqual(2)
	})

	t.Run("Remove Exercise from Favorites", func(t *testing.T) {
		response := e.DELETE("/api/v1/user/favorites/" + exerciseID).
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "exercise").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("removed")
	})

	t.Run("Check Exercise is No Longer Favorited", func(t *testing.T) {
		response := e.GET("/api/v1/user/favorites/" + exerciseID + "/check").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "exercise").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("data").Object().Value("is_favorited").Boolean().IsFalse()
	})

	t.Run("Remove Non-Existent Exercise Favorite", func(t *testing.T) {
		response := e.DELETE("/api/v1/user/favorites/" + exerciseID).
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "exercise").
			Expect().
			Status(404).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	// =====================
	// Workout Favorites Tests
	// =====================

	t.Run("Get Empty Workout Favorites", func(t *testing.T) {
		response := e.GET("/api/v1/user/favorites").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "workout").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("data").Array().Length().IsEqual(0)
	})

	t.Run("Add Workout to Favorites", func(t *testing.T) {
		response := e.POST("/api/v1/user/favorites").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "workout").
			WithJSON(map[string]interface{}{
				"item_id": workoutID,
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("added")

		data := response.Value("data").Object()
		data.Value("id").String().NotEmpty()
		data.Value("type").String().IsEqual("workout")
		data.Value("item_id").String().IsEqual(workoutID)
		data.Value("item").Object().Value("title").String().IsEqual("Test Workout")
	})

	t.Run("Check Workout is Favorited", func(t *testing.T) {
		response := e.GET("/api/v1/user/favorites/" + workoutID + "/check").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "workout").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("data").Object().Value("is_favorited").Boolean().IsTrue()
	})

	t.Run("Add Duplicate Workout Favorite", func(t *testing.T) {
		response := e.POST("/api/v1/user/favorites").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "workout").
			WithJSON(map[string]interface{}{
				"item_id": workoutID,
			}).
			Expect().
			Status(409).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("already")
	})

	t.Run("Add Second Workout to Favorites", func(t *testing.T) {
		e.POST("/api/v1/user/favorites").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "workout").
			WithJSON(map[string]interface{}{
				"item_id": workout2ID,
			}).
			Expect().
			Status(201)
	})

	t.Run("Get Workout Favorites", func(t *testing.T) {
		response := e.GET("/api/v1/user/favorites").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "workout").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("data").Array().Length().IsEqual(2)
	})

	t.Run("Remove Workout from Favorites", func(t *testing.T) {
		response := e.DELETE("/api/v1/user/favorites/" + workoutID).
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "workout").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("removed")
	})

	t.Run("Check Workout is No Longer Favorited", func(t *testing.T) {
		response := e.GET("/api/v1/user/favorites/" + workoutID + "/check").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "workout").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("data").Object().Value("is_favorited").Boolean().IsFalse()
	})

	// =====================
	// Validation Tests
	// =====================

	t.Run("Get Favorites Without Type Parameter", func(t *testing.T) {
		e.GET("/api/v1/user/favorites").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(400)
	})

	t.Run("Get Favorites With Invalid Type", func(t *testing.T) {
		e.GET("/api/v1/user/favorites").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "invalid").
			Expect().
			Status(400)
	})

	t.Run("Add Favorite Without Type Parameter", func(t *testing.T) {
		e.POST("/api/v1/user/favorites").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"item_id": exerciseID,
			}).
			Expect().
			Status(400)
	})

	t.Run("Add Favorite With Invalid Item ID", func(t *testing.T) {
		e.POST("/api/v1/user/favorites").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "exercise").
			WithJSON(map[string]interface{}{
				"item_id": "invalid-uuid",
			}).
			Expect().
			Status(400)
	})

	t.Run("Add Favorite For Non-Existent Exercise", func(t *testing.T) {
		// Use a valid random UUID that doesn't exist
		e.POST("/api/v1/user/favorites").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "exercise").
			WithJSON(map[string]interface{}{
				"item_id": "11111111-1111-1111-1111-111111111111",
			}).
			Expect().
			Status(404)
	})

	t.Run("Add Favorite For Non-Existent Workout", func(t *testing.T) {
		// Use a valid random UUID that doesn't exist
		e.POST("/api/v1/user/favorites").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "workout").
			WithJSON(map[string]interface{}{
				"item_id": "11111111-1111-1111-1111-111111111111",
			}).
			Expect().
			Status(404)
	})

	t.Run("Remove Favorite With Invalid UUID", func(t *testing.T) {
		e.DELETE("/api/v1/user/favorites/invalid-uuid").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "exercise").
			Expect().
			Status(400)
	})

	t.Run("Check Favorite With Invalid UUID", func(t *testing.T) {
		e.GET("/api/v1/user/favorites/invalid-uuid/check").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "exercise").
			Expect().
			Status(400)
	})

	// =====================
	// Authentication Tests
	// =====================

	t.Run("Get Favorites Without Auth", func(t *testing.T) {
		e.GET("/api/v1/user/favorites").
			WithQuery("type", "exercise").
			Expect().
			Status(401)
	})

	t.Run("Add Favorite Without Auth", func(t *testing.T) {
		e.POST("/api/v1/user/favorites").
			WithQuery("type", "exercise").
			WithJSON(map[string]interface{}{
				"item_id": exerciseID,
			}).
			Expect().
			Status(401)
	})

	t.Run("Remove Favorite Without Auth", func(t *testing.T) {
		e.DELETE("/api/v1/user/favorites/" + exerciseID).
			WithQuery("type", "exercise").
			Expect().
			Status(401)
	})

	t.Run("Check Favorite Without Auth", func(t *testing.T) {
		e.GET("/api/v1/user/favorites/" + exerciseID + "/check").
			WithQuery("type", "exercise").
			Expect().
			Status(401)
	})

	// =====================
	// Cross-User Isolation Tests
	// =====================

	// Create second user
	user2Data := map[string]interface{}{
		"email":            "favorites2@example.com",
		"password":         "FavoritesPass123!",
		"password_confirm": "FavoritesPass123!",
		"first_name":       "Favorites2",
		"last_name":        "User2",
	}

	e.POST("/api/v1/auth/register").
		WithJSON(user2Data).
		Expect().
		Status(201)

	token2 := GetAuthToken(e, "favorites2@example.com", "FavoritesPass123!")

	t.Run("User 2 Cannot See User 1 Favorites", func(t *testing.T) {
		// User 1 still has exercise2 favorited
		response := e.GET("/api/v1/user/favorites").
			WithHeader("Authorization", "Bearer "+token2).
			WithQuery("type", "exercise").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("data").Array().Length().IsEqual(0)
	})

	t.Run("User 2 Can Favorite Same Exercise", func(t *testing.T) {
		// User 2 can favorite exercise2 even though User 1 has it favorited
		response := e.POST("/api/v1/user/favorites").
			WithHeader("Authorization", "Bearer "+token2).
			WithQuery("type", "exercise").
			WithJSON(map[string]interface{}{
				"item_id": exercise2ID,
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
	})

	t.Run("Both Users Have Same Exercise Favorited Independently", func(t *testing.T) {
		// User 1's favorites
		resp1 := e.GET("/api/v1/user/favorites").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("type", "exercise").
			Expect().
			Status(200).
			JSON().
			Object()

		resp1.Value("data").Array().Length().IsEqual(1)

		// User 2's favorites
		resp2 := e.GET("/api/v1/user/favorites").
			WithHeader("Authorization", "Bearer "+token2).
			WithQuery("type", "exercise").
			Expect().
			Status(200).
			JSON().
			Object()

		resp2.Value("data").Array().Length().IsEqual(1)
	})
}
