package test

import (
	"lamari-fit-api/database"
	"lamari-fit-api/models"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestLocationSearchEndpoints(t *testing.T) {
	e := SetupTestApp(t)

	t.Run("User Location Search", func(t *testing.T) {
		CleanDatabase(t)
		SeedTestRoles(t)
		testUserLocationSearch(t, e)
	})

	t.Run("Trainer Location Search", func(t *testing.T) {
		CleanDatabase(t)
		SeedTestRoles(t)
		SeedTestSpecialties(t)
		testTrainerLocationSearch(t, e)
	})
}

func testUserLocationSearch(t *testing.T, e *httpexpect.Expect) {
	// Create users with different locations

	// User 1: New York (lat: 40.7128, lng: -74.0060)
	user1Data := map[string]interface{}{
		"email":            "user1_ny@example.com",
		"password":         "Password123!",
		"password_confirm": "Password123!",
		"first_name":       "John",
		"last_name":        "NewYork",
	}
	e.POST("/api/v1/auth/register").WithJSON(user1Data).Expect().Status(201)
	token1 := GetAuthToken(e, "user1_ny@example.com", "Password123!")

	// Set location for user1
	e.PUT("/api/v1/user/settings").
		WithHeader("Authorization", "Bearer "+token1).
		WithJSON(map[string]interface{}{
			"profile_visibility":     "public",
			"is_looking_for_trainer": true,
			"location": map[string]interface{}{
				"latitude":     40.7128,
				"longitude":    -74.0060,
				"city":         "New York",
				"region":       "NY",
				"country_code": "US",
			},
		}).
		Expect().
		Status(200)

	// User 2: Los Angeles (lat: 34.0522, lng: -118.2437)
	user2Data := map[string]interface{}{
		"email":            "user2_la@example.com",
		"password":         "Password123!",
		"password_confirm": "Password123!",
		"first_name":       "Jane",
		"last_name":        "LosAngeles",
	}
	e.POST("/api/v1/auth/register").WithJSON(user2Data).Expect().Status(201)
	token2 := GetAuthToken(e, "user2_la@example.com", "Password123!")

	e.PUT("/api/v1/user/settings").
		WithHeader("Authorization", "Bearer "+token2).
		WithJSON(map[string]interface{}{
			"profile_visibility":     "public",
			"is_looking_for_trainer": false,
			"location": map[string]interface{}{
				"latitude":     34.0522,
				"longitude":    -118.2437,
				"city":         "Los Angeles",
				"region":       "CA",
				"country_code": "US",
			},
		}).
		Expect().
		Status(200)

	// User 3: London (lat: 51.5074, lng: -0.1278)
	user3Data := map[string]interface{}{
		"email":            "user3_london@example.com",
		"password":         "Password123!",
		"password_confirm": "Password123!",
		"first_name":       "Bob",
		"last_name":        "London",
	}
	e.POST("/api/v1/auth/register").WithJSON(user3Data).Expect().Status(201)
	token3 := GetAuthToken(e, "user3_london@example.com", "Password123!")

	e.PUT("/api/v1/user/settings").
		WithHeader("Authorization", "Bearer "+token3).
		WithJSON(map[string]interface{}{
			"profile_visibility":     "public",
			"is_looking_for_trainer": true,
			"location": map[string]interface{}{
				"latitude":     51.5074,
				"longitude":    -0.1278,
				"city":         "London",
				"region":       "England",
				"country_code": "GB",
			},
		}).
		Expect().
		Status(200)

	// User 4: No location, public
	user4Data := map[string]interface{}{
		"email":            "user4_noloc@example.com",
		"password":         "Password123!",
		"password_confirm": "Password123!",
		"first_name":       "Alice",
		"last_name":        "NoLocation",
	}
	e.POST("/api/v1/auth/register").WithJSON(user4Data).Expect().Status(201)
	token4 := GetAuthToken(e, "user4_noloc@example.com", "Password123!")

	e.PUT("/api/v1/user/settings").
		WithHeader("Authorization", "Bearer "+token4).
		WithJSON(map[string]interface{}{
			"profile_visibility":     "public",
			"is_looking_for_trainer": true,
		}).
		Expect().
		Status(200)

	// Test cases

	t.Run("No Location Filter Returns All Public Users", func(t *testing.T) {
		response := e.GET("/api/v1/search/users").
			WithHeader("Authorization", "Bearer "+token1).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		data.Length().IsEqual(4) // All public users
	})

	t.Run("Radius Search From New York 50km", func(t *testing.T) {
		response := e.GET("/api/v1/search/users").
			WithHeader("Authorization", "Bearer "+token1).
			WithQuery("lat", "40.7128").
			WithQuery("lng", "-74.0060").
			WithQuery("radius_km", "50").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		// Should only find user1 (New York)
		data.Length().IsEqual(1)
		data.Value(0).Object().Value("first_name").String().IsEqual("John")
		data.Value(0).Object().Value("distance_km").Number().Le(50)
	})

	t.Run("Radius Search Large Radius Includes Multiple Users", func(t *testing.T) {
		response := e.GET("/api/v1/search/users").
			WithHeader("Authorization", "Bearer "+token1).
			WithQuery("lat", "40.7128").
			WithQuery("lng", "-74.0060").
			WithQuery("radius_km", "500").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		// Should find user1 (New York) but not LA or London
		data.Length().Ge(1)
	})

	t.Run("Structured Search By Country Code", func(t *testing.T) {
		response := e.GET("/api/v1/search/users").
			WithHeader("Authorization", "Bearer "+token1).
			WithQuery("country_code", "US").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		// Should find user1 (NY) and user2 (LA)
		data.Length().IsEqual(2)
	})

	t.Run("Structured Search By City", func(t *testing.T) {
		response := e.GET("/api/v1/search/users").
			WithHeader("Authorization", "Bearer "+token1).
			WithQuery("city", "London").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		data.Length().IsEqual(1)
		data.Value(0).Object().Value("first_name").String().IsEqual("Bob")
	})

	t.Run("Structured Search By Region", func(t *testing.T) {
		response := e.GET("/api/v1/search/users").
			WithHeader("Authorization", "Bearer "+token1).
			WithQuery("region", "CA").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		data.Length().IsEqual(1)
		data.Value(0).Object().Value("first_name").String().IsEqual("Jane")
	})

	t.Run("Free Text Search By Location Name", func(t *testing.T) {
		response := e.GET("/api/v1/search/users").
			WithHeader("Authorization", "Bearer "+token1).
			WithQuery("q", "New York").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		data.Length().IsEqual(1)
		data.Value(0).Object().Value("first_name").String().IsEqual("John")
	})

	t.Run("Free Text Search By User Name", func(t *testing.T) {
		response := e.GET("/api/v1/search/users").
			WithHeader("Authorization", "Bearer "+token1).
			WithQuery("q", "Jane").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		data.Length().IsEqual(1)
		data.Value(0).Object().Value("first_name").String().IsEqual("Jane")
	})

	t.Run("Combined Location and Is Looking For Trainer Filter", func(t *testing.T) {
		response := e.GET("/api/v1/search/users").
			WithHeader("Authorization", "Bearer "+token1).
			WithQuery("country_code", "US").
			WithQuery("is_looking_for_trainer", "true").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		// Should only find user1 (NY, looking for trainer)
		data.Length().IsEqual(1)
		data.Value(0).Object().Value("first_name").String().IsEqual("John")
		data.Value(0).Object().Value("is_looking_for_trainer").Boolean().IsTrue()
	})

	t.Run("Pagination Works", func(t *testing.T) {
		response := e.GET("/api/v1/search/users").
			WithHeader("Authorization", "Bearer "+token1).
			WithQuery("page", "1").
			WithQuery("limit", "2").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		data.Length().IsEqual(2)

		meta := response.Value("meta").Object()
		meta.Value("current_page").Number().IsEqual(1)
		meta.Value("per_page").Number().IsEqual(2)
		meta.Value("total_items").Number().IsEqual(4)
	})

	t.Run("Response Includes Location Data", func(t *testing.T) {
		response := e.GET("/api/v1/search/users").
			WithHeader("Authorization", "Bearer "+token1).
			WithQuery("city", "New York").
			Expect().
			Status(200).
			JSON().
			Object()

		data := response.Value("data").Array()
		data.Length().IsEqual(1)

		location := data.Value(0).Object().Value("location").Object()
		location.Value("city").String().IsEqual("New York")
		location.Value("region").String().IsEqual("NY")
		location.Value("country_code").String().IsEqual("US")
	})

	t.Run("Search Without Auth Returns 401", func(t *testing.T) {
		e.GET("/api/v1/search/users").
			Expect().
			Status(401)
	})

	_ = token4 // Suppress unused variable warning
}

func testTrainerLocationSearch(t *testing.T, e *httpexpect.Expect) {
	// Create trainers with different locations

	// Trainer 1: New York
	trainer1Data := map[string]interface{}{
		"email":            "trainer1_ny@example.com",
		"password":         "Password123!",
		"password_confirm": "Password123!",
		"first_name":       "Mike",
		"last_name":        "NewYorkTrainer",
	}
	e.POST("/api/v1/auth/register").WithJSON(trainer1Data).Expect().Status(201)
	token1 := GetAuthToken(e, "trainer1_ny@example.com", "Password123!")

	// Create trainer profile with location
	specialtyIDs := GetSpecialtyIDs(t, "Strength Training")
	e.POST("/api/v1/trainers/profile").
		WithHeader("Authorization", "Bearer "+token1).
		WithJSON(map[string]interface{}{
			"bio":                    "NYC based strength trainer",
			"specialty_ids":          specialtyIDs,
			"hourly_rate":            75.00,
			"visibility":             "public",
			"is_looking_for_clients": true,
			"location": map[string]interface{}{
				"latitude":     40.7128,
				"longitude":    -74.0060,
				"city":         "New York",
				"region":       "NY",
				"country_code": "US",
			},
		}).
		Expect().
		Status(201)

	// Trainer 2: Los Angeles
	trainer2Data := map[string]interface{}{
		"email":            "trainer2_la@example.com",
		"password":         "Password123!",
		"password_confirm": "Password123!",
		"first_name":       "Sarah",
		"last_name":        "LATrainer",
	}
	e.POST("/api/v1/auth/register").WithJSON(trainer2Data).Expect().Status(201)
	token2 := GetAuthToken(e, "trainer2_la@example.com", "Password123!")

	yogaIDs := GetSpecialtyIDs(t, "Yoga")
	e.POST("/api/v1/trainers/profile").
		WithHeader("Authorization", "Bearer "+token2).
		WithJSON(map[string]interface{}{
			"bio":                    "LA based yoga instructor",
			"specialty_ids":          yogaIDs,
			"hourly_rate":            60.00,
			"visibility":             "public",
			"is_looking_for_clients": true,
			"location": map[string]interface{}{
				"latitude":     34.0522,
				"longitude":    -118.2437,
				"city":         "Los Angeles",
				"region":       "CA",
				"country_code": "US",
			},
		}).
		Expect().
		Status(201)

	// Trainer 3: London
	trainer3Data := map[string]interface{}{
		"email":            "trainer3_london@example.com",
		"password":         "Password123!",
		"password_confirm": "Password123!",
		"first_name":       "James",
		"last_name":        "LondonTrainer",
	}
	e.POST("/api/v1/auth/register").WithJSON(trainer3Data).Expect().Status(201)
	token3 := GetAuthToken(e, "trainer3_london@example.com", "Password123!")

	cardioIDs := GetSpecialtyIDs(t, "Cardio")
	e.POST("/api/v1/trainers/profile").
		WithHeader("Authorization", "Bearer "+token3).
		WithJSON(map[string]interface{}{
			"bio":                    "London based cardio instructor",
			"specialty_ids":          cardioIDs,
			"hourly_rate":            50.00,
			"visibility":             "public",
			"is_looking_for_clients": false,
			"location": map[string]interface{}{
				"latitude":     51.5074,
				"longitude":    -0.1278,
				"city":         "London",
				"region":       "England",
				"country_code": "GB",
			},
		}).
		Expect().
		Status(201)

	// Test cases

	t.Run("No Location Filter Returns All Public Trainers", func(t *testing.T) {
		response := e.GET("/api/v1/search/trainers").
			WithHeader("Authorization", "Bearer "+token1).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		data.Length().IsEqual(3)
	})

	t.Run("Radius Search From New York", func(t *testing.T) {
		response := e.GET("/api/v1/search/trainers").
			WithHeader("Authorization", "Bearer "+token1).
			WithQuery("lat", "40.7128").
			WithQuery("lng", "-74.0060").
			WithQuery("radius_km", "50").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		// Should only find trainer1 (New York)
		data.Length().IsEqual(1)
		data.Value(0).Object().Value("user").Object().Value("first_name").String().IsEqual("Mike")
		data.Value(0).Object().Value("distance_km").Number().Le(50)
	})

	t.Run("Structured Search By Country Code", func(t *testing.T) {
		response := e.GET("/api/v1/search/trainers").
			WithHeader("Authorization", "Bearer "+token1).
			WithQuery("country_code", "US").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		// Should find trainer1 (NY) and trainer2 (LA)
		data.Length().IsEqual(2)
	})

	t.Run("Structured Search By City", func(t *testing.T) {
		response := e.GET("/api/v1/search/trainers").
			WithHeader("Authorization", "Bearer "+token1).
			WithQuery("city", "London").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		data.Length().IsEqual(1)
		data.Value(0).Object().Value("user").Object().Value("first_name").String().IsEqual("James")
	})

	t.Run("Free Text Search", func(t *testing.T) {
		response := e.GET("/api/v1/search/trainers").
			WithHeader("Authorization", "Bearer "+token1).
			WithQuery("q", "Los Angeles").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		data.Length().IsEqual(1)
		data.Value(0).Object().Value("user").Object().Value("first_name").String().IsEqual("Sarah")
	})

	t.Run("Filter By Specialty", func(t *testing.T) {
		response := e.GET("/api/v1/search/trainers").
			WithHeader("Authorization", "Bearer "+token1).
			WithQuery("specialty", "Yoga").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		data.Length().IsEqual(1)
		data.Value(0).Object().Value("user").Object().Value("first_name").String().IsEqual("Sarah")
	})

	t.Run("Filter By Is Looking For Clients", func(t *testing.T) {
		response := e.GET("/api/v1/search/trainers").
			WithHeader("Authorization", "Bearer "+token1).
			WithQuery("is_looking_for_clients", "true").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		// Should find trainer1 and trainer2
		data.Length().IsEqual(2)
	})

	t.Run("Combined Location Specialty and Availability Filter", func(t *testing.T) {
		response := e.GET("/api/v1/search/trainers").
			WithHeader("Authorization", "Bearer "+token1).
			WithQuery("country_code", "US").
			WithQuery("is_looking_for_clients", "true").
			WithQuery("specialty", "Strength").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		// Should only find trainer1 (NY, Strength Training, looking for clients)
		data.Length().IsEqual(1)
		data.Value(0).Object().Value("user").Object().Value("first_name").String().IsEqual("Mike")
	})

	t.Run("Sort By Rate", func(t *testing.T) {
		response := e.GET("/api/v1/search/trainers").
			WithHeader("Authorization", "Bearer "+token1).
			WithQuery("sort_by", "rate").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		data.Length().IsEqual(3)
		// First should be cheapest (London trainer at 50)
		data.Value(0).Object().Value("hourly_rate").Number().IsEqual(50)
	})

	t.Run("Response Includes Location And Review Data", func(t *testing.T) {
		response := e.GET("/api/v1/search/trainers").
			WithHeader("Authorization", "Bearer "+token1).
			WithQuery("city", "New York").
			Expect().
			Status(200).
			JSON().
			Object()

		data := response.Value("data").Array()
		data.Length().IsEqual(1)

		trainer := data.Value(0).Object()
		trainer.ContainsKey("id")
		trainer.ContainsKey("user")
		trainer.ContainsKey("specialties")
		trainer.ContainsKey("hourly_rate")
		trainer.ContainsKey("review_count")
		trainer.ContainsKey("average_rating")

		location := trainer.Value("location").Object()
		location.Value("city").String().IsEqual("New York")
		location.Value("region").String().IsEqual("NY")
		location.Value("country_code").String().IsEqual("US")
	})

	t.Run("Pagination Works", func(t *testing.T) {
		response := e.GET("/api/v1/search/trainers").
			WithHeader("Authorization", "Bearer "+token1).
			WithQuery("page", "1").
			WithQuery("limit", "2").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		data.Length().IsEqual(2)

		meta := response.Value("meta").Object()
		meta.Value("current_page").Number().IsEqual(1)
		meta.Value("per_page").Number().IsEqual(2)
		meta.Value("total_items").Number().IsEqual(3)
	})

	t.Run("Search Without Auth Returns 401", func(t *testing.T) {
		e.GET("/api/v1/search/trainers").
			Expect().
			Status(401)
	})

	_ = token2 // Suppress unused variable warning
	_ = token3 // Suppress unused variable warning
}

func TestHaversineDistance(t *testing.T) {
	e := SetupTestApp(t)

	t.Run("Distance Calculation", func(t *testing.T) {
		CleanDatabase(t)
		SeedTestRoles(t)
		testDistanceCalculation(t, e)
	})
}

func testDistanceCalculation(t *testing.T, e *httpexpect.Expect) {
	// Create a user with known location
	userData := map[string]interface{}{
		"email":            "distance_test@example.com",
		"password":         "Password123!",
		"password_confirm": "Password123!",
		"first_name":       "Distance",
		"last_name":        "Test",
	}
	e.POST("/api/v1/auth/register").WithJSON(userData).Expect().Status(201)
	token := GetAuthToken(e, "distance_test@example.com", "Password123!")

	// Set location to New York
	e.PUT("/api/v1/user/settings").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]interface{}{
			"profile_visibility": "public",
			"location": map[string]interface{}{
				"latitude":     40.7128,
				"longitude":    -74.0060,
				"city":         "New York",
				"country_code": "US",
			},
		}).
		Expect().
		Status(200)

	t.Run("Distance Calculated Correctly From Same Location", func(t *testing.T) {
		response := e.GET("/api/v1/search/users").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("lat", "40.7128").
			WithQuery("lng", "-74.0060").
			WithQuery("radius_km", "100").
			Expect().
			Status(200).
			JSON().
			Object()

		data := response.Value("data").Array()
		data.Length().IsEqual(1)
		// Distance from same point should be ~0
		data.Value(0).Object().Value("distance_km").Number().Lt(1)
	})

	t.Run("Default Radius Used When Not Specified", func(t *testing.T) {
		// Create another user far away
		user2Data := map[string]interface{}{
			"email":            "far_user@example.com",
			"password":         "Password123!",
			"password_confirm": "Password123!",
			"first_name":       "Far",
			"last_name":        "User",
		}
		e.POST("/api/v1/auth/register").WithJSON(user2Data).Expect().Status(201)
		token2 := GetAuthToken(e, "far_user@example.com", "Password123!")

		// Set location to London (far from NY)
		e.PUT("/api/v1/user/settings").
			WithHeader("Authorization", "Bearer "+token2).
			WithJSON(map[string]interface{}{
				"profile_visibility": "public",
				"location": map[string]interface{}{
					"latitude":     51.5074,
					"longitude":    -0.1278,
					"city":         "London",
					"country_code": "GB",
				},
			}).
			Expect().
			Status(200)

		// Search from NY with no radius (should use default 25km)
		response := e.GET("/api/v1/search/users").
			WithHeader("Authorization", "Bearer "+token).
			WithQuery("lat", "40.7128").
			WithQuery("lng", "-74.0060").
			Expect().
			Status(200).
			JSON().
			Object()

		data := response.Value("data").Array()
		// Should only find the NY user (London is too far for 25km default)
		data.Length().IsEqual(1)
		data.Value(0).Object().Value("first_name").String().IsEqual("Distance")
	})
}

func TestUserLocationUpdate(t *testing.T) {
	e := SetupTestApp(t)

	t.Run("User Location Update", func(t *testing.T) {
		CleanDatabase(t)
		SeedTestRoles(t)
		testUserLocationUpdate(t, e)
	})
}

func testUserLocationUpdate(t *testing.T, e *httpexpect.Expect) {
	// Create a user
	userData := map[string]interface{}{
		"email":            "location_update@example.com",
		"password":         "Password123!",
		"password_confirm": "Password123!",
		"first_name":       "Location",
		"last_name":        "Update",
	}
	e.POST("/api/v1/auth/register").WithJSON(userData).Expect().Status(201)
	token := GetAuthToken(e, "location_update@example.com", "Password123!")

	t.Run("Set Initial Location", func(t *testing.T) {
		response := e.PUT("/api/v1/user/settings").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"location": map[string]interface{}{
					"latitude":     40.7128,
					"longitude":    -74.0060,
					"city":         "New York",
					"region":       "NY",
					"country_code": "US",
				},
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		location := response.Value("data").Object().Value("location").Object()
		location.Value("city").String().IsEqual("New York")
		location.Value("region").String().IsEqual("NY")
		location.Value("country_code").String().IsEqual("US")
	})

	t.Run("Update Location", func(t *testing.T) {
		response := e.PUT("/api/v1/user/settings").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"location": map[string]interface{}{
					"city":   "Brooklyn",
					"region": "NY",
				},
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		location := response.Value("data").Object().Value("location").Object()
		// City and region should be updated
		location.Value("city").String().IsEqual("Brooklyn")
		location.Value("region").String().IsEqual("NY")
		// Country code should be preserved from previous update
		location.Value("country_code").String().IsEqual("US")
	})

	t.Run("Location Persists Across Settings Requests", func(t *testing.T) {
		// Get current settings
		response := e.GET("/api/v1/user/settings").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(200).
			JSON().
			Object()

		location := response.Value("data").Object().Value("location").Object()
		location.Value("city").String().IsEqual("Brooklyn")
		location.Value("region").String().IsEqual("NY")
		location.Value("country_code").String().IsEqual("US")
	})
}

func TestTrainerLocationUpdate(t *testing.T) {
	e := SetupTestApp(t)

	t.Run("Trainer Location Update", func(t *testing.T) {
		CleanDatabase(t)
		SeedTestRoles(t)
		SeedTestSpecialties(t)
		testTrainerLocationUpdate(t, e)
	})
}

func testTrainerLocationUpdate(t *testing.T, e *httpexpect.Expect) {
	// Create a trainer
	trainerData := map[string]interface{}{
		"email":            "trainer_loc@example.com",
		"password":         "Password123!",
		"password_confirm": "Password123!",
		"first_name":       "Trainer",
		"last_name":        "Loc",
	}
	e.POST("/api/v1/auth/register").WithJSON(trainerData).Expect().Status(201)
	token := GetAuthToken(e, "trainer_loc@example.com", "Password123!")

	specialtyIDs := GetSpecialtyIDs(t, "Strength Training")

	t.Run("Create Trainer Profile With Location", func(t *testing.T) {
		response := e.POST("/api/v1/trainers/profile").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"bio":           "Test trainer",
				"specialty_ids": specialtyIDs,
				"hourly_rate":   50.00,
				"visibility":    "public",
				"location": map[string]interface{}{
					"latitude":     40.7128,
					"longitude":    -74.0060,
					"city":         "New York",
					"region":       "NY",
					"country_code": "US",
				},
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		location := response.Value("data").Object().Value("location").Object()
		location.Value("city").String().IsEqual("New York")
	})

	t.Run("Update Trainer Profile Location", func(t *testing.T) {
		response := e.PUT("/api/v1/trainers/profile").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"location": map[string]interface{}{
					"city":   "Brooklyn",
					"region": "NY",
				},
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		location := response.Value("data").Object().Value("location").Object()
		location.Value("city").String().IsEqual("Brooklyn")
	})

	// Get trainer ID for public profile test
	var trainer models.TrainerProfile
	database.DB.Where("bio = ?", "Test trainer").First(&trainer)

	t.Run("Trainer Public Profile Shows Location", func(t *testing.T) {
		response := e.GET("/api/v1/trainers/"+trainer.ID.String()).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		location := response.Value("data").Object().Value("location").Object()
		location.Value("city").String().IsEqual("Brooklyn")
	})
}
