package test

import (
	"fit-flow-api/database"
	"testing"
)

func TestRPEScaleEndpoints(t *testing.T) {
	e := SetupTestApp(t)

	// Seed global RPE scale
	database.SeedGlobalRPEScale()

	// Register a user for testing
	e.POST("/api/v1/auth/register").
		WithJSON(map[string]interface{}{
			"email":            "rpe_test@example.com",
			"password":         "TestPassword123!",
			"password_confirm": "TestPassword123!",
			"first_name":       "RPE",
			"last_name":        "Tester",
		}).
		Expect().
		Status(201)

	token := GetAuthToken(e, "rpe_test@example.com", "TestPassword123!")

	t.Run("Get Global RPE Scale", func(t *testing.T) {
		resp := e.GET("/api/v1/rpe/scales/global").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(200).
			JSON().
			Object()

		resp.Value("success").Boolean().IsTrue()
		data := resp.Value("data").Object()
		data.Value("name").String().IsEqual("Standard RPE Scale")
		data.Value("is_global").Boolean().IsTrue()
		data.Value("min_value").Number().IsEqual(1)
		data.Value("max_value").Number().IsEqual(10)
		data.Value("values").Array().Length().IsEqual(10)
	})

	t.Run("List RPE Scales", func(t *testing.T) {
		resp := e.GET("/api/v1/rpe/scales").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(200).
			JSON().
			Object()

		resp.Value("success").Boolean().IsTrue()
		// Should include at least the global scale
		resp.Value("data").Array().NotEmpty()
	})

	t.Run("Create Custom RPE Scale", func(t *testing.T) {
		t.Run("Full Scale With Values", func(t *testing.T) {
			resp := e.POST("/api/v1/rpe/scales").
				WithHeader("Authorization", "Bearer "+token).
				WithJSON(map[string]interface{}{
					"name":        "My Custom Scale",
					"description": "A 1-5 scale for beginners",
					"min_value":   1,
					"max_value":   5,
					"values": []map[string]interface{}{
						{"value": 1, "label": "Easy", "description": "Minimal effort"},
						{"value": 2, "label": "Light", "description": "Light effort"},
						{"value": 3, "label": "Moderate", "description": "Medium effort"},
						{"value": 4, "label": "Hard", "description": "Hard effort"},
						{"value": 5, "label": "Maximum", "description": "Max effort"},
					},
				}).
				Expect().
				Status(201).
				JSON().
				Object()

			resp.Value("success").Boolean().IsTrue()
			data := resp.Value("data").Object()
			data.Value("name").String().IsEqual("My Custom Scale")
			data.Value("is_global").Boolean().IsFalse()
			data.Value("min_value").Number().IsEqual(1)
			data.Value("max_value").Number().IsEqual(5)
			data.Value("values").Array().Length().IsEqual(5)
		})

		t.Run("Scale Without Initial Values", func(t *testing.T) {
			resp := e.POST("/api/v1/rpe/scales").
				WithHeader("Authorization", "Bearer "+token).
				WithJSON(map[string]interface{}{
					"name":        "Empty Scale",
					"description": "A scale to add values later",
					"max_value":   100,
				}).
				Expect().
				Status(201).
				JSON().
				Object()

			resp.Value("success").Boolean().IsTrue()
			data := resp.Value("data").Object()
			data.Value("name").String().IsEqual("Empty Scale")
			data.Value("min_value").Number().IsEqual(1) // defaults to 1
			data.Value("max_value").Number().IsEqual(100)
		})

		t.Run("Invalid Scale - Min >= Max", func(t *testing.T) {
			e.POST("/api/v1/rpe/scales").
				WithHeader("Authorization", "Bearer "+token).
				WithJSON(map[string]interface{}{
					"name":      "Invalid Scale",
					"min_value": 10,
					"max_value": 5,
				}).
				Expect().
				Status(400)
		})
	})

	t.Run("Get Custom Scale", func(t *testing.T) {
		// First create a scale
		createResp := e.POST("/api/v1/rpe/scales").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"name":      "Get Test Scale",
				"min_value": 1,
				"max_value": 3,
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		scaleID := createResp.Value("data").Object().Value("id").String().Raw()

		// Get the scale
		resp := e.GET("/api/v1/rpe/scales/" + scaleID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(200).
			JSON().
			Object()

		resp.Value("success").Boolean().IsTrue()
		resp.Value("data").Object().Value("name").String().IsEqual("Get Test Scale")
	})

	t.Run("Update Custom Scale", func(t *testing.T) {
		// Create a scale
		createResp := e.POST("/api/v1/rpe/scales").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"name":        "Update Test Scale",
				"description": "Original description",
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		scaleID := createResp.Value("data").Object().Value("id").String().Raw()

		// Update the scale
		resp := e.PUT("/api/v1/rpe/scales/" + scaleID).
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"name":        "Updated Scale Name",
				"description": "New description",
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		resp.Value("success").Boolean().IsTrue()
		resp.Value("data").Object().Value("name").String().IsEqual("Updated Scale Name")
		resp.Value("data").Object().Value("description").String().IsEqual("New description")
	})

	t.Run("Add Value to Scale", func(t *testing.T) {
		// Create a scale without values
		createResp := e.POST("/api/v1/rpe/scales").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"name":      "Add Value Test",
				"min_value": 1,
				"max_value": 5,
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		scaleID := createResp.Value("data").Object().Value("id").String().Raw()

		// Add a value
		resp := e.POST("/api/v1/rpe/scales/" + scaleID + "/values").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"value":       3,
				"label":       "Medium",
				"description": "Medium effort level",
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		resp.Value("success").Boolean().IsTrue()
		resp.Value("data").Object().Value("value").Number().IsEqual(3)
		resp.Value("data").Object().Value("label").String().IsEqual("Medium")
	})

	t.Run("Delete Custom Scale", func(t *testing.T) {
		// Create a scale
		createResp := e.POST("/api/v1/rpe/scales").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"name": "Delete Test Scale",
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		scaleID := createResp.Value("data").Object().Value("id").String().Raw()

		// Delete the scale
		e.DELETE("/api/v1/rpe/scales/" + scaleID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(204)

		// Verify it's deleted
		e.GET("/api/v1/rpe/scales/" + scaleID).
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(404)
	})

	t.Run("Cannot Modify Global Scale", func(t *testing.T) {
		// Get global scale ID
		globalResp := e.GET("/api/v1/rpe/scales/global").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(200).
			JSON().
			Object()

		globalID := globalResp.Value("data").Object().Value("id").String().Raw()

		t.Run("Cannot Update Global Scale", func(t *testing.T) {
			e.PUT("/api/v1/rpe/scales/" + globalID).
				WithHeader("Authorization", "Bearer "+token).
				WithJSON(map[string]interface{}{
					"name": "Hacked Global Scale",
				}).
				Expect().
				Status(403)
		})

		t.Run("Cannot Delete Global Scale", func(t *testing.T) {
			e.DELETE("/api/v1/rpe/scales/" + globalID).
				WithHeader("Authorization", "Bearer "+token).
				Expect().
				Status(403)
		})

		t.Run("Cannot Add Value to Global Scale", func(t *testing.T) {
			e.POST("/api/v1/rpe/scales/" + globalID + "/values").
				WithHeader("Authorization", "Bearer "+token).
				WithJSON(map[string]interface{}{
					"value": 11,
					"label": "Super Max",
				}).
				Expect().
				Status(403)
		})
	})

	t.Run("Value Validation", func(t *testing.T) {
		// Create a scale
		createResp := e.POST("/api/v1/rpe/scales").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"name":      "Validation Test",
				"min_value": 1,
				"max_value": 5,
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		scaleID := createResp.Value("data").Object().Value("id").String().Raw()

		t.Run("Value Out of Range", func(t *testing.T) {
			e.POST("/api/v1/rpe/scales/" + scaleID + "/values").
				WithHeader("Authorization", "Bearer "+token).
				WithJSON(map[string]interface{}{
					"value": 10, // Max is 5
					"label": "Invalid",
				}).
				Expect().
				Status(400)
		})
	})

	t.Run("Authorization", func(t *testing.T) {
		t.Run("Cannot Access Without Auth", func(t *testing.T) {
			e.GET("/api/v1/rpe/scales").
				Expect().
				Status(401)
		})
	})
}

func TestTrainerClientRPEScaleAccess(t *testing.T) {
	e := SetupTestApp(t)

	// Seed global RPE scale and specialties for trainer profile
	database.SeedGlobalRPEScale()
	SeedTestSpecialties(t)

	// Get specialty IDs for trainer profile
	specialtyIDs := GetSpecialtyIDs(t, "Strength Training")

	// Register trainer with trainer profile
	trainerResp := e.POST("/api/v1/auth/register").
		WithJSON(map[string]interface{}{
			"email":            "trainer_rpe_access@example.com",
			"password":         "TrainerPass123!",
			"password_confirm": "TrainerPass123!",
			"first_name":       "Trainer",
			"last_name":        "RPEAccess",
			"trainer_profile": map[string]interface{}{
				"bio":           "Test trainer for RPE access",
				"hourly_rate":   50,
				"location":      "Test City",
				"specialty_ids": specialtyIDs,
				"visibility":    "public",
			},
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	trainerToken := trainerResp.Value("data").Object().Value("token").String().Raw()

	// Register client
	e.POST("/api/v1/auth/register").
		WithJSON(map[string]interface{}{
			"email":            "client_rpe_access@example.com",
			"password":         "ClientPass123!",
			"password_confirm": "ClientPass123!",
			"first_name":       "Client",
			"last_name":        "RPEAccess",
		}).
		Expect().
		Status(201)

	clientToken := GetAuthToken(e, "client_rpe_access@example.com", "ClientPass123!")

	// Get client user ID
	clientProfile := e.GET("/api/v1/auth/profile").
		WithHeader("Authorization", "Bearer "+clientToken).
		Expect().
		Status(200).
		JSON().
		Object()

	clientID := clientProfile.Value("data").Object().Value("id").String().Raw()

	// Trainer creates custom RPE scale
	scaleResp := e.POST("/api/v1/rpe/scales").
		WithHeader("Authorization", "Bearer "+trainerToken).
		WithJSON(map[string]interface{}{
			"name":        "Trainer's Custom Scale",
			"description": "A scale for my clients",
			"min_value":   1,
			"max_value":   10,
			"values": []map[string]interface{}{
				{"value": 1, "label": "Very Easy", "description": "Minimal effort"},
				{"value": 5, "label": "Moderate", "description": "Medium effort"},
				{"value": 10, "label": "Maximum", "description": "Max effort"},
			},
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	scaleID := scaleResp.Value("data").Object().Value("id").String().Raw()

	t.Run("Client Cannot Access Trainer Scale Before Relationship", func(t *testing.T) {
		// Client tries to get the scale - should fail
		e.GET("/api/v1/rpe/scales/" + scaleID).
			WithHeader("Authorization", "Bearer "+clientToken).
			Expect().
			Status(404)

		// Client lists scales - should not include trainer's scale
		listResp := e.GET("/api/v1/rpe/scales").
			WithHeader("Authorization", "Bearer "+clientToken).
			Expect().
			Status(200).
			JSON().
			Object()

		scales := listResp.Value("data").Array()
		for _, scale := range scales.Iter() {
			scaleName := scale.Object().Value("name").String().Raw()
			if scaleName == "Trainer's Custom Scale" {
				t.Fatalf("Client should not see trainer's scale before relationship")
			}
		}
	})

	// Trainer invites client
	inviteResp := e.POST("/api/v1/trainers/clients").
		WithHeader("Authorization", "Bearer "+trainerToken).
		WithJSON(map[string]interface{}{
			"client_id": clientID,
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	invitationID := inviteResp.Value("data").Object().Value("id").String().Raw()

	// Client accepts invitation
	e.PUT("/api/v1/me/trainer-invitations/" + invitationID).
		WithHeader("Authorization", "Bearer "+clientToken).
		WithJSON(map[string]interface{}{
			"action": "accept",
		}).
		Expect().
		Status(200)

	t.Run("Client Can Access Trainer Scale After Relationship", func(t *testing.T) {
		// Client can now get the scale
		resp := e.GET("/api/v1/rpe/scales/" + scaleID).
			WithHeader("Authorization", "Bearer "+clientToken).
			Expect().
			Status(200).
			JSON().
			Object()

		resp.Value("success").Boolean().IsTrue()
		resp.Value("data").Object().Value("name").String().IsEqual("Trainer's Custom Scale")
		resp.Value("data").Object().Value("values").Array().Length().IsEqual(3)
	})

	t.Run("Client Can List Trainer Scale After Relationship", func(t *testing.T) {
		listResp := e.GET("/api/v1/rpe/scales").
			WithHeader("Authorization", "Bearer "+clientToken).
			Expect().
			Status(200).
			JSON().
			Object()

		scales := listResp.Value("data").Array()
		foundTrainerScale := false
		for _, scale := range scales.Iter() {
			scaleName := scale.Object().Value("name").String().Raw()
			if scaleName == "Trainer's Custom Scale" {
				foundTrainerScale = true
				break
			}
		}
		if !foundTrainerScale {
			t.Fatalf("Client should see trainer's scale after relationship is active")
		}
	})

	t.Run("Client Cannot Modify Trainer Scale", func(t *testing.T) {
		// Client cannot update trainer's scale
		e.PUT("/api/v1/rpe/scales/" + scaleID).
			WithHeader("Authorization", "Bearer "+clientToken).
			WithJSON(map[string]interface{}{
				"name": "Hacked Scale Name",
			}).
			Expect().
			Status(403)

		// Client cannot delete trainer's scale
		e.DELETE("/api/v1/rpe/scales/" + scaleID).
			WithHeader("Authorization", "Bearer "+clientToken).
			Expect().
			Status(403)

		// Client cannot add values to trainer's scale
		e.POST("/api/v1/rpe/scales/" + scaleID + "/values").
			WithHeader("Authorization", "Bearer "+clientToken).
			WithJSON(map[string]interface{}{
				"value": 7,
				"label": "Hard",
			}).
			Expect().
			Status(403)
	})

	t.Run("Trainer Still Has Full Access", func(t *testing.T) {
		// Trainer can still update their scale
		e.PUT("/api/v1/rpe/scales/" + scaleID).
			WithHeader("Authorization", "Bearer "+trainerToken).
			WithJSON(map[string]interface{}{
				"description": "Updated description",
			}).
			Expect().
			Status(200)

		// Trainer can add values
		e.POST("/api/v1/rpe/scales/" + scaleID + "/values").
			WithHeader("Authorization", "Bearer "+trainerToken).
			WithJSON(map[string]interface{}{
				"value": 7,
				"label": "Hard",
			}).
			Expect().
			Status(201)
	})
}
