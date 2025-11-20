package test

import (
	"net/http"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
)

func TestWeightLogEndpoints(t *testing.T) {
	e := SetupTestApp(t)

	t.Run("Create Weight Log", func(t *testing.T) {
		CleanDatabase(t)
		testCreateWeightLog(t, e)
	})

	t.Run("Get Weight Logs", func(t *testing.T) {
		CleanDatabase(t)
		testGetWeightLogs(t, e)
	})

	t.Run("Get Single Weight Log", func(t *testing.T) {
		CleanDatabase(t)
		testGetSingleWeightLog(t, e)
	})

	t.Run("Update Weight Log", func(t *testing.T) {
		CleanDatabase(t)
		testUpdateWeightLog(t, e)
	})

	t.Run("Get Weight Stats", func(t *testing.T) {
		CleanDatabase(t)
		testGetWeightStats(t, e)
	})

	t.Run("Delete Weight Log", func(t *testing.T) {
		CleanDatabase(t)
		testDeleteWeightLog(t, e)
	})

	t.Run("User Isolation", func(t *testing.T) {
		CleanDatabase(t)
		testWeightLogUserIsolation(t, e)
	})
}

func createWeightTestUserAndGetToken(e *httpexpect.Expect, email, password, firstName, lastName string) string {
	// Register user
	e.POST("/api/v1/auth/register").
		WithJSON(map[string]interface{}{
			"email":            email,
			"first_name":       firstName,
			"last_name":        lastName,
			"password":         password,
			"password_confirm": password,
		}).
		Expect().
		Status(http.StatusCreated)

	// Login and get token
	resp := e.POST("/api/v1/auth/login").
		WithJSON(map[string]interface{}{
			"email":    email,
			"password": password,
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	return resp.Value("data").Object().Value("token").String().Raw()
}

func testCreateWeightLog(t *testing.T, e *httpexpect.Expect) {
	token := createWeightTestUserAndGetToken(e, "weightlog@example.com", "WeightLog123!", "John", "Doe")

	t.Run("Successful Creation", func(t *testing.T) {
		resp := e.POST("/api/v1/user/weight-logs").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"weight_kg": 80.5,
				"notes":     "Morning weigh-in",
			}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()

		resp.Value("success").Boolean().IsTrue()
		resp.Value("message").String().Contains("Weight logged successfully")
		data := resp.Value("data").Object()
		data.Value("weight_kg").Number().IsEqual(80.5)
		data.Value("weight_lbs").Number().Gt(0)
		data.Value("notes").String().IsEqual("Morning weigh-in")
	})

	t.Run("Create Without Notes", func(t *testing.T) {
		resp := e.POST("/api/v1/user/weight-logs").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"weight_kg": 81.0,
			}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()

		resp.Value("success").Boolean().IsTrue()
		resp.Value("data").Object().Value("weight_kg").Number().IsEqual(81.0)
	})

	t.Run("Create Without Auth", func(t *testing.T) {
		e.POST("/api/v1/user/weight-logs").
			WithJSON(map[string]interface{}{
				"weight_kg": 80.0,
			}).
			Expect().
			Status(http.StatusUnauthorized)
	})

	t.Run("Invalid Weight - Too Low", func(t *testing.T) {
		e.POST("/api/v1/user/weight-logs").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"weight_kg": 15.0,
			}).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("Invalid Weight - Too High", func(t *testing.T) {
		e.POST("/api/v1/user/weight-logs").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"weight_kg": 600.0,
			}).
			Expect().
			Status(http.StatusBadRequest)
	})
}

func testGetWeightLogs(t *testing.T, e *httpexpect.Expect) {
	token := createWeightTestUserAndGetToken(e, "getlogs@example.com", "GetLogs123!", "Get", "Logs")

	// Create some logs
	for i := 0; i < 3; i++ {
		e.POST("/api/v1/user/weight-logs").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"weight_kg": 80.0 + float64(i),
			}).
			Expect().
			Status(http.StatusCreated)
	}

	t.Run("Get All Logs", func(t *testing.T) {
		resp := e.GET("/api/v1/user/weight-logs").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		resp.Value("success").Boolean().IsTrue()
		resp.Value("data").Array().Length().IsEqual(3)
		resp.Value("meta").Object().Value("total_items").Number().IsEqual(3)
	})

	t.Run("Get Logs with Pagination", func(t *testing.T) {
		resp := e.GET("/api/v1/user/weight-logs").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("page", 1).
			WithQuery("limit", 2).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		resp.Value("data").Array().Length().IsEqual(2)
		resp.Value("meta").Object().Value("current_page").Number().IsEqual(1)
		resp.Value("meta").Object().Value("per_page").Number().IsEqual(2)
	})

	t.Run("Get Logs with Date Filter", func(t *testing.T) {
		today := time.Now().Format("2006-01-02")

		resp := e.GET("/api/v1/user/weight-logs").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("start_date", today).
			WithQuery("end_date", today).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		resp.Value("success").Boolean().IsTrue()
		resp.Value("data").Array().Length().Ge(1)
	})

	t.Run("Get Logs Without Auth", func(t *testing.T) {
		e.GET("/api/v1/user/weight-logs").
			Expect().
			Status(http.StatusUnauthorized)
	})
}

func testGetSingleWeightLog(t *testing.T, e *httpexpect.Expect) {
	token := createWeightTestUserAndGetToken(e, "singlelog@example.com", "SingleLog123!", "Single", "Log")

	// Create a log
	resp := e.POST("/api/v1/user/weight-logs").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]interface{}{
			"weight_kg": 80.5,
			"notes":     "Test log",
		}).
		Expect().
		Status(http.StatusCreated).
		JSON().Object()

	logID := resp.Value("data").Object().Value("id").String().Raw()

	t.Run("Get Existing Log", func(t *testing.T) {
		resp := e.GET("/api/v1/user/weight-logs/{id}", logID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		resp.Value("success").Boolean().IsTrue()
		resp.Value("data").Object().Value("id").String().IsEqual(logID)
		resp.Value("data").Object().Value("weight_kg").Number().IsEqual(80.5)
	})

	t.Run("Get Non-existent Log", func(t *testing.T) {
		e.GET("/api/v1/user/weight-logs/{id}", "00000000-0000-0000-0000-000000000000").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("Get Log Without Auth", func(t *testing.T) {
		e.GET("/api/v1/user/weight-logs/{id}", logID).
			Expect().
			Status(http.StatusUnauthorized)
	})
}

func testUpdateWeightLog(t *testing.T, e *httpexpect.Expect) {
	token := createWeightTestUserAndGetToken(e, "updatelog@example.com", "UpdateLog123!", "Update", "Log")

	// Create a log
	resp := e.POST("/api/v1/user/weight-logs").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]interface{}{
			"weight_kg": 80.5,
			"notes":     "Original notes",
		}).
		Expect().
		Status(http.StatusCreated).
		JSON().Object()

	logID := resp.Value("data").Object().Value("id").String().Raw()

	t.Run("Successful Update", func(t *testing.T) {
		resp := e.PUT("/api/v1/user/weight-logs/{id}", logID).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"weight_kg": 80.0,
				"notes":     "Updated notes",
			}).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		resp.Value("success").Boolean().IsTrue()
		resp.Value("data").Object().Value("weight_kg").Number().IsEqual(80.0)
		resp.Value("data").Object().Value("notes").String().IsEqual("Updated notes")
	})

	t.Run("Update Non-existent Log", func(t *testing.T) {
		e.PUT("/api/v1/user/weight-logs/{id}", "00000000-0000-0000-0000-000000000000").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"weight_kg": 80.0,
			}).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("Update Without Auth", func(t *testing.T) {
		e.PUT("/api/v1/user/weight-logs/{id}", logID).
			WithJSON(map[string]interface{}{
				"weight_kg": 80.0,
			}).
			Expect().
			Status(http.StatusUnauthorized)
	})
}

func testGetWeightStats(t *testing.T, e *httpexpect.Expect) {
	token := createWeightTestUserAndGetToken(e, "stats@example.com", "Stats123!", "Stats", "User")

	// Create several weight entries
	for i := 0; i < 5; i++ {
		e.POST("/api/v1/user/weight-logs").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"weight_kg": 79.0 + float64(i)*0.5,
			}).
			Expect().
			Status(http.StatusCreated)
	}

	t.Run("Get Stats Default Period", func(t *testing.T) {
		resp := e.GET("/api/v1/user/weight-logs/stats").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		resp.Value("success").Boolean().IsTrue()
		data := resp.Value("data").Object()
		data.Value("total_entries").Number().IsEqual(5)
		data.Value("latest_weight_kg").Number().Gt(0)
		data.Value("min_weight_kg").Number().Gt(0)
		data.Value("max_weight_kg").Number().Gt(0)
		data.Value("avg_weight_kg").Number().Gt(0)
		data.Value("period_days").Number().IsEqual(30)
	})

	t.Run("Get Stats Custom Period", func(t *testing.T) {
		resp := e.GET("/api/v1/user/weight-logs/stats").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("days", 7).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		resp.Value("data").Object().Value("period_days").Number().IsEqual(7)
	})

	t.Run("Get Stats Without Auth", func(t *testing.T) {
		e.GET("/api/v1/user/weight-logs/stats").
			Expect().
			Status(http.StatusUnauthorized)
	})
}

func testDeleteWeightLog(t *testing.T, e *httpexpect.Expect) {
	token := createWeightTestUserAndGetToken(e, "deletelog@example.com", "DeleteLog123!", "Delete", "Log")

	// Create a log to delete
	resp := e.POST("/api/v1/user/weight-logs").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]interface{}{
			"weight_kg": 82.0,
		}).
		Expect().
		Status(http.StatusCreated).
		JSON().Object()

	logID := resp.Value("data").Object().Value("id").String().Raw()

	t.Run("Successful Delete", func(t *testing.T) {
		e.DELETE("/api/v1/user/weight-logs/{id}", logID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNoContent)

		// Verify it's deleted
		e.GET("/api/v1/user/weight-logs/{id}", logID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("Delete Non-existent Log", func(t *testing.T) {
		e.DELETE("/api/v1/user/weight-logs/{id}", "00000000-0000-0000-0000-000000000000").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("Delete Without Auth", func(t *testing.T) {
		// Create another log for this test
		resp := e.POST("/api/v1/user/weight-logs").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"weight_kg": 83.0,
			}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()

		deleteID := resp.Value("data").Object().Value("id").String().Raw()

		e.DELETE("/api/v1/user/weight-logs/{id}", deleteID).
			Expect().
			Status(http.StatusUnauthorized)
	})
}

func testWeightLogUserIsolation(t *testing.T, e *httpexpect.Expect) {
	// Create two users
	token1 := createWeightTestUserAndGetToken(e, "user1@weight.com", "UserOne123!", "User", "One")
	token2 := createWeightTestUserAndGetToken(e, "user2@weight.com", "UserTwo123!", "User", "Two")

	// User 1 creates a weight log
	resp := e.POST("/api/v1/user/weight-logs").
		WithHeader("Authorization", "Bearer "+token1).
		WithJSON(map[string]interface{}{
			"weight_kg": 75.0,
		}).
		Expect().
		Status(http.StatusCreated).
		JSON().Object()

	user1LogID := resp.Value("data").Object().Value("id").String().Raw()

	t.Run("User Cannot Access Other User's Log", func(t *testing.T) {
		e.GET("/api/v1/user/weight-logs/{id}", user1LogID).
			WithHeader("Authorization", "Bearer "+token2).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("User Cannot Update Other User's Log", func(t *testing.T) {
		e.PUT("/api/v1/user/weight-logs/{id}", user1LogID).
			WithHeader("Authorization", "Bearer "+token2).
			WithJSON(map[string]interface{}{
				"weight_kg": 100.0,
			}).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("User Cannot Delete Other User's Log", func(t *testing.T) {
		e.DELETE("/api/v1/user/weight-logs/{id}", user1LogID).
			WithHeader("Authorization", "Bearer "+token2).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("Users See Only Their Own Logs", func(t *testing.T) {
		// User 2 creates their own log
		e.POST("/api/v1/user/weight-logs").
			WithHeader("Authorization", "Bearer "+token2).
			WithJSON(map[string]interface{}{
				"weight_kg": 85.0,
			}).
			Expect().
			Status(http.StatusCreated)

		// User 1 should only see their log
		resp1 := e.GET("/api/v1/user/weight-logs").
			WithHeader("Authorization", "Bearer "+token1).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		resp1.Value("meta").Object().Value("total_items").Number().IsEqual(1)

		// User 2 should only see their log
		resp2 := e.GET("/api/v1/user/weight-logs").
			WithHeader("Authorization", "Bearer "+token2).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		resp2.Value("meta").Object().Value("total_items").Number().IsEqual(1)
	})
}
