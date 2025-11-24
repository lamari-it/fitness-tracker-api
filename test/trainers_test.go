package test

import (
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestTrainerProfileEndpoints(t *testing.T) {
	e := SetupTestApp(t)

	t.Run("Create Trainer Profile", func(t *testing.T) {
		CleanDatabase(t)
		testCreateTrainerProfile(t, e)
	})

	t.Run("Get Trainer Profile", func(t *testing.T) {
		CleanDatabase(t)
		testGetTrainerProfile(t, e)
	})

	t.Run("Update Trainer Profile", func(t *testing.T) {
		CleanDatabase(t)
		testUpdateTrainerProfile(t, e)
	})

	t.Run("Delete Trainer Profile", func(t *testing.T) {
		CleanDatabase(t)
		testDeleteTrainerProfile(t, e)
	})

	t.Run("List Trainers", func(t *testing.T) {
		CleanDatabase(t)
		testListTrainers(t, e)
	})

	t.Run("Get Trainer Public Profile", func(t *testing.T) {
		CleanDatabase(t)
		testGetTrainerPublicProfile(t, e)
	})

	t.Run("Visibility Access Control", func(t *testing.T) {
		CleanDatabase(t)
		testVisibilityAccessControl(t, e)
	})
}

func createTestUserAndGetToken(e *httpexpect.Expect, email, password, firstName, lastName string) string {
	userData := map[string]interface{}{
		"email":            email,
		"password":         password,
		"password_confirm": password,
		"first_name":       firstName,
		"last_name":        lastName,
	}

	e.POST("/api/v1/auth/register").
		WithJSON(userData).
		Expect().
		Status(201)

	return GetAuthToken(e, email, password)
}

func testCreateTrainerProfile(t *testing.T, e *httpexpect.Expect) {
	// Seed specialties for tests
	SeedTestSpecialties(t)
	specialtyIDs := GetSpecialtyIDs(t, "Strength Training", "Weight Loss", "Bodybuilding")

	token := createTestUserAndGetToken(e, "trainer@example.com", "TrainerPass123!", "John", "Trainer")

	t.Run("Successful Profile Creation", func(t *testing.T) {
		profileData := map[string]interface{}{
			"bio":           "Certified personal trainer with 5+ years experience in strength training.",
			"specialty_ids": specialtyIDs,
			"hourly_rate":   75.00,
			"location":      "New York, NY",
		}

		response := e.POST("/api/v1/trainers/profile").
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
		data.Value("user_id").String().NotEmpty()
		data.Value("bio").String().IsEqual("Certified personal trainer with 5+ years experience in strength training.")
		data.Value("specialties").Array().Length().IsEqual(3)
		data.Value("hourly_rate").Number().IsEqual(75.00)
		data.Value("location").String().IsEqual("New York, NY")
		data.Value("visibility").String().IsEqual("private") // Default visibility
		data.Value("created_at").String().NotEmpty()
		data.Value("updated_at").String().NotEmpty()
	})

	t.Run("Successful Profile Creation With Custom Visibility", func(t *testing.T) {
		customToken := createTestUserAndGetToken(e, "custom@example.com", "CustomPass123!", "Custom", "Trainer")
		customSpecialtyIDs := GetSpecialtyIDs(t, "Yoga", "Mobility")

		profileData := map[string]interface{}{
			"bio":           "Private trainer with exclusive clientele.",
			"specialty_ids": customSpecialtyIDs,
			"hourly_rate":   200.00,
			"location":      "Beverly Hills, CA",
			"visibility":    "private",
		}

		response := e.POST("/api/v1/trainers/profile").
			WithHeader("Authorization", "Bearer "+customToken).
			WithJSON(profileData).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("visibility").String().IsEqual("private")
	})

	t.Run("Duplicate Profile Creation", func(t *testing.T) {
		cardioIDs := GetSpecialtyIDs(t, "Cardio")

		profileData := map[string]interface{}{
			"bio":           "Another trainer bio",
			"specialty_ids": cardioIDs,
			"hourly_rate":   50.00,
			"location":      "Los Angeles, CA",
		}

		response := e.POST("/api/v1/trainers/profile").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(profileData).
			Expect().
			Status(409).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("already exists")
	})

	t.Run("Invalid Profile Data", func(t *testing.T) {
		newToken := createTestUserAndGetToken(e, "newtrainer@example.com", "NewPass123!", "New", "Trainer")
		strengthIDs := GetSpecialtyIDs(t, "Strength Training")

		// Test cases that should fail validation
		testCases := []struct {
			name        string
			profileData map[string]interface{}
		}{
			{
				name: "Negative Hourly Rate",
				profileData: map[string]interface{}{
					"bio":           "Certified personal trainer with experience.",
					"specialty_ids": strengthIDs,
					"hourly_rate":   -10.00,
					"location":      "New York, NY",
				},
			},
			{
				name: "Bio Too Long",
				profileData: map[string]interface{}{
					"bio":           string(make([]byte, 1001)), // Over 1000 chars
					"specialty_ids": strengthIDs,
					"hourly_rate":   75.00,
					"location":      "New York, NY",
				},
			},
			{
				name: "Hourly Rate Too High",
				profileData: map[string]interface{}{
					"bio":           "Valid bio",
					"specialty_ids": strengthIDs,
					"hourly_rate":   10000.00, // Over 9999.99
					"location":      "New York, NY",
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				response := e.POST("/api/v1/trainers/profile").
					WithHeader("Authorization", "Bearer "+newToken).
					WithJSON(tc.profileData).
					Expect().
					Status(400).
					JSON().
					Object()

				response.Value("success").Boolean().IsFalse()
			})
		}
	})

	t.Run("Create Profile Without Auth", func(t *testing.T) {
		authStrengthIDs := GetSpecialtyIDs(t, "Strength Training")

		profileData := map[string]interface{}{
			"bio":           "Certified personal trainer with experience.",
			"specialty_ids": authStrengthIDs,
			"hourly_rate":   75.00,
			"location":      "New York, NY",
		}

		response := e.POST("/api/v1/trainers/profile").
			WithJSON(profileData).
			Expect().
			Status(401).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})
}

func testGetTrainerProfile(t *testing.T, e *httpexpect.Expect) {
	SeedTestSpecialties(t)
	token := createTestUserAndGetToken(e, "getprofile@example.com", "GetPass123!", "Get", "Trainer")

	t.Run("Get Non-existent Profile", func(t *testing.T) {
		response := e.GET("/api/v1/trainers/profile").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(404).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("not found")
	})

	// Create profile first
	specialtyIDs := GetSpecialtyIDs(t, "Functional Fitness", "Mobility")
	profileData := map[string]interface{}{
		"bio":           "Experienced trainer specializing in functional fitness.",
		"specialty_ids": specialtyIDs,
		"hourly_rate":   60.00,
		"location":      "San Francisco, CA",
	}

	e.POST("/api/v1/trainers/profile").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(profileData).
		Expect().
		Status(201)

	t.Run("Successful Get Profile", func(t *testing.T) {
		response := e.GET("/api/v1/trainers/profile").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("retrieved")

		data := response.Value("data").Object()
		data.Value("id").String().NotEmpty()
		data.Value("bio").String().IsEqual("Experienced trainer specializing in functional fitness.")
		data.Value("specialties").Array().Length().IsEqual(2)
		data.Value("hourly_rate").Number().IsEqual(60.00)
		data.Value("location").String().IsEqual("San Francisco, CA")

		// Check user is preloaded
		user := data.Value("user").Object()
		user.Value("id").String().NotEmpty()
		user.Value("email").String().IsEqual("getprofile@example.com")
		user.Value("first_name").String().IsEqual("Get")
		user.Value("last_name").String().IsEqual("Trainer")
	})

	t.Run("Get Profile Without Auth", func(t *testing.T) {
		response := e.GET("/api/v1/trainers/profile").
			Expect().
			Status(401).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})
}

func testUpdateTrainerProfile(t *testing.T, e *httpexpect.Expect) {
	SeedTestSpecialties(t)
	token := createTestUserAndGetToken(e, "updateprofile@example.com", "UpdatePass123!", "Update", "Trainer")

	// Create profile first
	initialSpecialtyIDs := GetSpecialtyIDs(t, "Strength Training")
	profileData := map[string]interface{}{
		"bio":           "Original bio for the trainer profile.",
		"specialty_ids": initialSpecialtyIDs,
		"hourly_rate":   50.00,
		"location":      "Miami, FL",
	}

	e.POST("/api/v1/trainers/profile").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(profileData).
		Expect().
		Status(201)

	// Get specialty IDs for update tests - accessible to all subtests
	updateSpecialtyIDs := GetSpecialtyIDs(t, "Strength Training", "HIIT", "Cardio")

	t.Run("Successful Full Update", func(t *testing.T) {
		updateData := map[string]interface{}{
			"bio":           "Updated bio with more experience and certifications.",
			"specialty_ids": updateSpecialtyIDs,
			"hourly_rate":   85.00,
			"location":      "Miami Beach, FL",
		}

		response := e.PUT("/api/v1/trainers/profile").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(updateData).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("updated")

		data := response.Value("data").Object()
		data.Value("bio").String().IsEqual("Updated bio with more experience and certifications.")
		data.Value("specialties").Array().Length().IsEqual(3)
		data.Value("hourly_rate").Number().IsEqual(85.00)
		data.Value("location").String().IsEqual("Miami Beach, FL")
	})

	t.Run("Partial Update - Bio Only", func(t *testing.T) {
		updateData := map[string]interface{}{
			"bio":           "Another updated bio content here.",
			"specialty_ids": updateSpecialtyIDs, // Required for validation
		}

		response := e.PUT("/api/v1/trainers/profile").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(updateData).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("bio").String().IsEqual("Another updated bio content here.")
		// Other fields should remain unchanged
		data.Value("hourly_rate").Number().IsEqual(85.00)
		data.Value("location").String().IsEqual("Miami Beach, FL")
	})

	t.Run("Partial Update - Hourly Rate Only", func(t *testing.T) {
		updateData := map[string]interface{}{
			"hourly_rate":   100.00,
			"specialty_ids": updateSpecialtyIDs, // Required for validation
		}

		response := e.PUT("/api/v1/trainers/profile").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(updateData).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("hourly_rate").Number().IsEqual(100.00)
	})

	t.Run("Update Visibility", func(t *testing.T) {
		updateData := map[string]interface{}{
			"visibility":    "link_only",
			"specialty_ids": updateSpecialtyIDs, // Required for validation
		}

		response := e.PUT("/api/v1/trainers/profile").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(updateData).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("visibility").String().IsEqual("link_only")
	})

	t.Run("Update Non-existent Profile", func(t *testing.T) {
		newToken := createTestUserAndGetToken(e, "nonexistent@example.com", "NoProfile123!", "No", "Profile")

		updateData := map[string]interface{}{
			"bio": "Trying to update non-existent profile.",
		}

		response := e.PUT("/api/v1/trainers/profile").
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
			"bio": "Unauthorized update attempt.",
		}

		response := e.PUT("/api/v1/trainers/profile").
			WithJSON(updateData).
			Expect().
			Status(401).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})
}

func testDeleteTrainerProfile(t *testing.T, e *httpexpect.Expect) {
	SeedTestSpecialties(t)
	token := createTestUserAndGetToken(e, "deleteprofile@example.com", "DeletePass123!", "Delete", "Trainer")

	t.Run("Delete Non-existent Profile", func(t *testing.T) {
		response := e.DELETE("/api/v1/trainers/profile").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(404).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("not found")
	})

	// Create profile first
	cardioIDs := GetSpecialtyIDs(t, "Cardio")
	profileData := map[string]interface{}{
		"bio":           "Profile to be deleted.",
		"specialty_ids": cardioIDs,
		"hourly_rate":   40.00,
		"location":      "Boston, MA",
	}

	e.POST("/api/v1/trainers/profile").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(profileData).
		Expect().
		Status(201)

	t.Run("Successful Delete", func(t *testing.T) {
		e.DELETE("/api/v1/trainers/profile").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(204)

		// Verify profile is deleted
		response := e.GET("/api/v1/trainers/profile").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(404).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("Delete Without Auth", func(t *testing.T) {
		response := e.DELETE("/api/v1/trainers/profile").
			Expect().
			Status(401).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})
}

func testListTrainers(t *testing.T, e *httpexpect.Expect) {
	SeedTestSpecialties(t)

	// Create multiple trainers
	trainer1Token := createTestUserAndGetToken(e, "trainer1@example.com", "Pass123!", "Alice", "Smith")
	trainer2Token := createTestUserAndGetToken(e, "trainer2@example.com", "Pass123!", "Bob", "Jones")
	trainer3Token := createTestUserAndGetToken(e, "trainer3@example.com", "Pass123!", "Carol", "Williams")
	trainer4Token := createTestUserAndGetToken(e, "trainer4@example.com", "Pass123!", "David", "Brown")
	trainer5Token := createTestUserAndGetToken(e, "trainer5@example.com", "Pass123!", "Eve", "Davis")
	regularUserToken := createTestUserAndGetToken(e, "regular@example.com", "Pass123!", "Regular", "User")

	// Get specialty IDs for profiles
	strengthBodybuildingIDs := GetSpecialtyIDs(t, "Strength Training", "Bodybuilding")
	yogaMobilityIDs := GetSpecialtyIDs(t, "Yoga", "Mobility")
	cardioHIITIDs := GetSpecialtyIDs(t, "Cardio", "HIIT", "Weight Loss")
	functionalIDs := GetSpecialtyIDs(t, "Functional Fitness")
	rehabIDs := GetSpecialtyIDs(t, "Rehabilitation")

	// Create trainer profiles with different visibility settings
	profile1 := map[string]interface{}{
		"bio":           "Strength training expert with certifications.",
		"specialty_ids": strengthBodybuildingIDs,
		"hourly_rate":   80.00,
		"location":      "New York, NY",
		"visibility":    "public",
	}
	e.POST("/api/v1/trainers/profile").
		WithHeader("Authorization", "Bearer "+trainer1Token).
		WithJSON(profile1).
		Expect().
		Status(201)

	profile2 := map[string]interface{}{
		"bio":           "Yoga and mobility specialist.",
		"specialty_ids": yogaMobilityIDs,
		"hourly_rate":   60.00,
		"location":      "Los Angeles, CA",
		"visibility":    "public",
	}
	e.POST("/api/v1/trainers/profile").
		WithHeader("Authorization", "Bearer "+trainer2Token).
		WithJSON(profile2).
		Expect().
		Status(201)

	profile3 := map[string]interface{}{
		"bio":           "Cardio and HIIT training expert.",
		"specialty_ids": cardioHIITIDs,
		"hourly_rate":   70.00,
		"location":      "New York, NY",
		"visibility":    "public",
	}
	e.POST("/api/v1/trainers/profile").
		WithHeader("Authorization", "Bearer "+trainer3Token).
		WithJSON(profile3).
		Expect().
		Status(201)

	// Create link_only profile - should NOT appear in list
	profile4 := map[string]interface{}{
		"bio":           "Exclusive trainer for link-only access.",
		"specialty_ids": functionalIDs,
		"hourly_rate":   150.00,
		"location":      "New York, NY",
		"visibility":    "link_only",
	}
	e.POST("/api/v1/trainers/profile").
		WithHeader("Authorization", "Bearer "+trainer4Token).
		WithJSON(profile4).
		Expect().
		Status(201)

	// Create private profile - should NOT appear in list
	profile5 := map[string]interface{}{
		"bio":           "Private trainer not listed publicly.",
		"specialty_ids": rehabIDs,
		"hourly_rate":   250.00,
		"location":      "New York, NY",
		"visibility":    "private",
	}
	e.POST("/api/v1/trainers/profile").
		WithHeader("Authorization", "Bearer "+trainer5Token).
		WithJSON(profile5).
		Expect().
		Status(201)

	t.Run("List All Trainers", func(t *testing.T) {
		response := e.GET("/api/v1/trainers").
			WithHeader("Authorization", "Bearer "+regularUserToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("retrieved")

		data := response.Value("data").Array()
		data.Length().IsEqual(3)

		// Check that each trainer has public fields
		firstTrainer := data.Value(0).Object()
		firstTrainer.Value("id").String().NotEmpty()
		firstTrainer.Value("bio").String().NotEmpty()
		firstTrainer.Value("specialties").Array().NotEmpty()
		firstTrainer.Value("hourly_rate").Number().Gt(0)
		firstTrainer.Value("location").String().NotEmpty()
		firstTrainer.Value("user").Object().NotEmpty()
		firstTrainer.Value("review_count").Number().IsEqual(0)
		firstTrainer.Value("average_rating").Number().IsEqual(0)

		// Check pagination meta
		meta := response.Value("meta").Object()
		meta.Value("current_page").Number().IsEqual(1)
		meta.Value("per_page").Number().IsEqual(10)
		meta.Value("total_items").Number().IsEqual(3)
	})

	t.Run("List Trainers With Pagination", func(t *testing.T) {
		response := e.GET("/api/v1/trainers").
			WithQuery("page", 1).
			WithQuery("limit", 2).
			WithHeader("Authorization", "Bearer "+regularUserToken).
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

	t.Run("Search Trainers By Name", func(t *testing.T) {
		response := e.GET("/api/v1/trainers").
			WithQuery("search", "Alice").
			WithHeader("Authorization", "Bearer "+regularUserToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		data.Length().IsEqual(1)

		trainer := data.Value(0).Object()
		trainer.Value("user").Object().Value("first_name").String().IsEqual("Alice")
	})

	t.Run("Filter Trainers By Location", func(t *testing.T) {
		response := e.GET("/api/v1/trainers").
			WithQuery("location", "New York").
			WithHeader("Authorization", "Bearer "+regularUserToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		data.Length().IsEqual(2) // Alice and Carol are in New York
	})

	t.Run("Filter Trainers By Specialty", func(t *testing.T) {
		response := e.GET("/api/v1/trainers").
			WithQuery("specialty", "Yoga").
			WithHeader("Authorization", "Bearer "+regularUserToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		data.Length().IsEqual(1) // Only Bob has Yoga specialty

		trainer := data.Value(0).Object()
		trainer.Value("user").Object().Value("first_name").String().IsEqual("Bob")
	})

	t.Run("Sort Trainers By Rate", func(t *testing.T) {
		response := e.GET("/api/v1/trainers").
			WithQuery("sort_by", "rate").
			WithHeader("Authorization", "Bearer "+regularUserToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		data.Length().IsEqual(3)

		// Should be sorted by hourly_rate ascending (60, 70, 80)
		data.Value(0).Object().Value("hourly_rate").Number().IsEqual(60.00)
		data.Value(1).Object().Value("hourly_rate").Number().IsEqual(70.00)
		data.Value(2).Object().Value("hourly_rate").Number().IsEqual(80.00)
	})

	t.Run("List Trainers Without Auth", func(t *testing.T) {
		response := e.GET("/api/v1/trainers").
			Expect().
			Status(401).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})
}

func testGetTrainerPublicProfile(t *testing.T, e *httpexpect.Expect) {
	SeedTestSpecialties(t)
	trainerToken := createTestUserAndGetToken(e, "publictrainer@example.com", "PublicPass123!", "Public", "Trainer")
	clientToken := createTestUserAndGetToken(e, "client@example.com", "ClientPass123!", "Client", "User")

	// Create trainer profile with public visibility
	specialtyIDs := GetSpecialtyIDs(t, "Strength Training", "Cardio")
	profileData := map[string]interface{}{
		"bio":           "Public trainer profile for viewing.",
		"specialty_ids": specialtyIDs,
		"hourly_rate":   65.00,
		"location":      "Chicago, IL",
		"visibility":    "public",
	}

	createResponse := e.POST("/api/v1/trainers/profile").
		WithHeader("Authorization", "Bearer "+trainerToken).
		WithJSON(profileData).
		Expect().
		Status(201).
		JSON().
		Object()

	trainerID := createResponse.Value("data").Object().Value("id").String().Raw()

	t.Run("Get Public Profile Successfully", func(t *testing.T) {
		response := e.GET("/api/v1/trainers/"+trainerID).
			WithHeader("Authorization", "Bearer "+clientToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("retrieved")

		data := response.Value("data").Object()
		data.Value("id").String().IsEqual(trainerID)
		data.Value("bio").String().IsEqual("Public trainer profile for viewing.")
		data.Value("specialties").Array().Length().IsEqual(2)
		data.Value("hourly_rate").Number().IsEqual(65.00)
		data.Value("location").String().IsEqual("Chicago, IL")
		data.Value("review_count").Number().IsEqual(0)
		data.Value("average_rating").Number().IsEqual(0)

		// Check user has limited fields (public response)
		user := data.Value("user").Object()
		user.Value("id").String().NotEmpty()
		user.Value("first_name").String().IsEqual("Public")
		user.Value("last_name").String().IsEqual("Trainer")
		// Should NOT contain sensitive info
		user.NotContainsKey("email")
		user.NotContainsKey("password")
	})

	t.Run("Get Non-existent Trainer", func(t *testing.T) {
		fakeID := "00000000-0000-0000-0000-000000000000"
		response := e.GET("/api/v1/trainers/"+fakeID).
			WithHeader("Authorization", "Bearer "+clientToken).
			Expect().
			Status(404).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("not found")
	})

	t.Run("Get Trainer With Invalid UUID", func(t *testing.T) {
		response := e.GET("/api/v1/trainers/invalid-uuid").
			WithHeader("Authorization", "Bearer "+clientToken).
			Expect().
			Status(400).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("Get Public Profile Without Auth", func(t *testing.T) {
		response := e.GET("/api/v1/trainers/" + trainerID).
			Expect().
			Status(401).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})
}

func testVisibilityAccessControl(t *testing.T, e *httpexpect.Expect) {
	SeedTestSpecialties(t)

	// Create trainers with different visibility settings
	publicTrainerToken := createTestUserAndGetToken(e, "public_trainer@example.com", "Pass123!", "Public", "Trainer")
	linkOnlyTrainerToken := createTestUserAndGetToken(e, "linkonly_trainer@example.com", "Pass123!", "LinkOnly", "Trainer")
	privateTrainerToken := createTestUserAndGetToken(e, "private_trainer@example.com", "Pass123!", "Private", "Trainer")
	viewerToken := createTestUserAndGetToken(e, "viewer@example.com", "Pass123!", "Viewer", "User")

	// Get specialty IDs
	strengthIDs := GetSpecialtyIDs(t, "Strength Training")
	yogaIDs := GetSpecialtyIDs(t, "Yoga")
	cardioIDs := GetSpecialtyIDs(t, "Cardio")

	// Create public profile
	publicProfile := map[string]interface{}{
		"bio":           "Public trainer visible to everyone.",
		"specialty_ids": strengthIDs,
		"hourly_rate":   50.00,
		"location":      "Public City",
		"visibility":    "public",
	}
	publicResp := e.POST("/api/v1/trainers/profile").
		WithHeader("Authorization", "Bearer "+publicTrainerToken).
		WithJSON(publicProfile).
		Expect().
		Status(201).
		JSON().
		Object()
	publicTrainerID := publicResp.Value("data").Object().Value("id").String().Raw()

	// Create link_only profile
	linkOnlyProfile := map[string]interface{}{
		"bio":           "Link-only trainer viewable with direct link.",
		"specialty_ids": yogaIDs,
		"hourly_rate":   100.00,
		"location":      "Link City",
		"visibility":    "link_only",
	}
	linkOnlyResp := e.POST("/api/v1/trainers/profile").
		WithHeader("Authorization", "Bearer "+linkOnlyTrainerToken).
		WithJSON(linkOnlyProfile).
		Expect().
		Status(201).
		JSON().
		Object()
	linkOnlyTrainerID := linkOnlyResp.Value("data").Object().Value("id").String().Raw()

	// Create private profile
	privateProfile := map[string]interface{}{
		"bio":           "Private trainer only visible to owner.",
		"specialty_ids": cardioIDs,
		"hourly_rate":   200.00,
		"location":      "Private City",
		"visibility":    "private",
	}
	privateResp := e.POST("/api/v1/trainers/profile").
		WithHeader("Authorization", "Bearer "+privateTrainerToken).
		WithJSON(privateProfile).
		Expect().
		Status(201).
		JSON().
		Object()
	privateTrainerID := privateResp.Value("data").Object().Value("id").String().Raw()

	t.Run("Public Profile Accessible By Anyone", func(t *testing.T) {
		response := e.GET("/api/v1/trainers/"+publicTrainerID).
			WithHeader("Authorization", "Bearer "+viewerToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("visibility").String().IsEqual("public")
		data.Value("bio").String().Contains("Public trainer")
	})

	t.Run("Link Only Profile Accessible With Direct Link", func(t *testing.T) {
		response := e.GET("/api/v1/trainers/"+linkOnlyTrainerID).
			WithHeader("Authorization", "Bearer "+viewerToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("visibility").String().IsEqual("link_only")
		data.Value("bio").String().Contains("Link-only trainer")
	})

	t.Run("Private Profile Not Accessible By Others", func(t *testing.T) {
		response := e.GET("/api/v1/trainers/"+privateTrainerID).
			WithHeader("Authorization", "Bearer "+viewerToken).
			Expect().
			Status(404).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("not found")
	})

	t.Run("Private Profile Accessible By Owner", func(t *testing.T) {
		response := e.GET("/api/v1/trainers/"+privateTrainerID).
			WithHeader("Authorization", "Bearer "+privateTrainerToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("visibility").String().IsEqual("private")
		data.Value("bio").String().Contains("Private trainer")
	})

	t.Run("Only Public Trainers In Search Results", func(t *testing.T) {
		response := e.GET("/api/v1/trainers").
			WithHeader("Authorization", "Bearer "+viewerToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()

		// Should only contain the public trainer, not link_only or private
		data.Length().IsEqual(1)
		firstTrainer := data.Value(0).Object()
		firstTrainer.Value("visibility").String().IsEqual("public")
	})

	t.Run("Owner Can Edit Own Private Profile", func(t *testing.T) {
		updateData := map[string]interface{}{
			"bio":           "Updated private trainer bio content.",
			"specialty_ids": cardioIDs, // Required for validation
		}

		response := e.PUT("/api/v1/trainers/profile").
			WithHeader("Authorization", "Bearer "+privateTrainerToken).
			WithJSON(updateData).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("bio").String().Contains("Updated private trainer")
	})

	t.Run("Owner Can Delete Own Profile", func(t *testing.T) {
		// Create a new trainer to delete
		deleteToken := createTestUserAndGetToken(e, "delete_test@example.com", "Pass123!", "Delete", "Test")
		deleteSpecialtyIDs := GetSpecialtyIDs(t, "HIIT")
		deleteProfile := map[string]interface{}{
			"bio":           "Profile to be deleted by owner.",
			"specialty_ids": deleteSpecialtyIDs,
			"hourly_rate":   50.00,
			"location":      "Delete City",
			"visibility":    "private",
		}
		e.POST("/api/v1/trainers/profile").
			WithHeader("Authorization", "Bearer "+deleteToken).
			WithJSON(deleteProfile).
			Expect().
			Status(201)

		e.DELETE("/api/v1/trainers/profile").
			WithHeader("Authorization", "Bearer "+deleteToken).
			Expect().
			Status(204)
	})

	t.Run("Change Visibility From Private To Public", func(t *testing.T) {
		// Include specialty_ids to satisfy validation requirement
		updateData := map[string]interface{}{
			"visibility":    "public",
			"specialty_ids": cardioIDs,
		}

		response := e.PUT("/api/v1/trainers/profile").
			WithHeader("Authorization", "Bearer "+privateTrainerToken).
			WithJSON(updateData).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("visibility").String().IsEqual("public")

		// Now should appear in search results
		listResponse := e.GET("/api/v1/trainers").
			WithHeader("Authorization", "Bearer "+viewerToken).
			Expect().
			Status(200).
			JSON().
			Object()

		listResponse.Value("data").Array().Length().IsEqual(2) // Now 2 public trainers
	})
}
