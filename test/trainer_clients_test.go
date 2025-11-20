package test

import (
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestTrainerClientEndpoints(t *testing.T) {
	e := SetupTestApp(t)

	t.Run("Trainer Client Flow", func(t *testing.T) {
		CleanDatabase(t)
		testTrainerClientFlow(t, e)
	})

	t.Run("Invitation Responses", func(t *testing.T) {
		CleanDatabase(t)
		testInvitationResponses(t, e)
	})

	t.Run("Authorization Checks", func(t *testing.T) {
		CleanDatabase(t)
		testTrainerClientAuthorization(t, e)
	})
}

func testTrainerClientFlow(t *testing.T, e *httpexpect.Expect) {
	// Seed specialties for trainer profile creation
	SeedTestSpecialties(t)
	specialtyIDs := GetSpecialtyIDs(t, "Strength Training", "Weight Loss")

	// Create trainer user
	trainerToken := createTestUserAndGetToken(e, "trainer@example.com", "TrainerPass123!", "John", "Trainer")

	// Create client user
	clientToken := createTestUserAndGetToken(e, "client@example.com", "ClientPass123!", "Jane", "Client")

	// Get client user ID from the response
	clientResponse := e.GET("/api/v1/auth/profile").
		WithHeader("Authorization", "Bearer "+clientToken).
		Expect().
		Status(200).
		JSON().
		Object()
	clientID := clientResponse.Value("data").Object().Value("id").String().Raw()

	t.Run("Trainer Must Have Profile To Invite", func(t *testing.T) {
		// Try to invite without trainer profile
		response := e.POST("/api/v1/trainers/clients").
			WithHeader("Authorization", "Bearer "+trainerToken).
			WithJSON(map[string]interface{}{
				"client_id": clientID,
			}).
			Expect().
			Status(403).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("trainer profile")
	})

	// Create trainer profile
	e.POST("/api/v1/trainers/profile").
		WithHeader("Authorization", "Bearer "+trainerToken).
		WithJSON(map[string]interface{}{
			"bio":           "Certified personal trainer with 5+ years experience.",
			"specialty_ids": specialtyIDs,
			"hourly_rate":   75.00,
			"location":      "New York, NY",
		}).
		Expect().
		Status(201)

	var invitationID string

	t.Run("Invite Client Successfully", func(t *testing.T) {
		response := e.POST("/api/v1/trainers/clients").
			WithHeader("Authorization", "Bearer "+trainerToken).
			WithJSON(map[string]interface{}{
				"client_id": clientID,
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("invitation sent")

		data := response.Value("data").Object()
		invitationID = data.Value("id").String().Raw()
		data.Value("status").String().IsEqual("pending")
		data.Value("client").Object().Value("first_name").String().IsEqual("Jane")
	})

	t.Run("Cannot Invite Same Client Twice", func(t *testing.T) {
		response := e.POST("/api/v1/trainers/clients").
			WithHeader("Authorization", "Bearer "+trainerToken).
			WithJSON(map[string]interface{}{
				"client_id": clientID,
			}).
			Expect().
			Status(409).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("pending")
	})

	t.Run("Client Sees Pending Invitation", func(t *testing.T) {
		response := e.GET("/api/v1/me/trainer-invitations").
			WithHeader("Authorization", "Bearer "+clientToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Array()
		data.Length().IsEqual(1)
		data.Value(0).Object().Value("status").String().IsEqual("pending")
		data.Value(0).Object().Value("trainer").Object().Value("first_name").String().IsEqual("John")
	})

	t.Run("Trainer Sees Pending Client", func(t *testing.T) {
		response := e.GET("/api/v1/trainers/clients").
			WithHeader("Authorization", "Bearer "+trainerToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Array()
		data.Length().IsEqual(1)
		data.Value(0).Object().Value("status").String().IsEqual("pending")
	})

	t.Run("Trainer Can Filter Clients By Status", func(t *testing.T) {
		// Filter by active (should be empty)
		response := e.GET("/api/v1/trainers/clients").
			WithHeader("Authorization", "Bearer "+trainerToken).
			WithQuery("status", "active").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("data").Array().Length().IsEqual(0)

		// Filter by pending (should have 1)
		response = e.GET("/api/v1/trainers/clients").
			WithHeader("Authorization", "Bearer "+trainerToken).
			WithQuery("status", "pending").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("data").Array().Length().IsEqual(1)
	})

	t.Run("Client Accepts Invitation", func(t *testing.T) {
		response := e.PUT("/api/v1/me/trainer-invitations/" + invitationID).
			WithHeader("Authorization", "Bearer "+clientToken).
			WithJSON(map[string]interface{}{
				"action": "accept",
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("data").Object().Value("status").String().IsEqual("active")
	})

	t.Run("Client Sees Active Trainer", func(t *testing.T) {
		response := e.GET("/api/v1/me/trainers").
			WithHeader("Authorization", "Bearer "+clientToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Array()
		data.Length().IsEqual(1)
		data.Value(0).Object().Value("status").String().IsEqual("active")
	})

	t.Run("No Pending Invitations After Accept", func(t *testing.T) {
		response := e.GET("/api/v1/me/trainer-invitations").
			WithHeader("Authorization", "Bearer "+clientToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("data").Array().Length().IsEqual(0)
	})

	t.Run("Trainer Removes Client", func(t *testing.T) {
		// Get the link ID from trainer's client list
		clientsResponse := e.GET("/api/v1/trainers/clients").
			WithHeader("Authorization", "Bearer "+trainerToken).
			Expect().
			Status(200).
			JSON().
			Object()

		linkID := clientsResponse.Value("data").Array().Value(0).Object().Value("id").String().Raw()

		response := e.DELETE("/api/v1/trainers/clients/" + linkID).
			WithHeader("Authorization", "Bearer "+trainerToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
	})

	t.Run("Client No Longer Has Active Trainer", func(t *testing.T) {
		response := e.GET("/api/v1/me/trainers").
			WithHeader("Authorization", "Bearer "+clientToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("data").Array().Length().IsEqual(0)
	})
}

func testInvitationResponses(t *testing.T, e *httpexpect.Expect) {
	// Seed specialties
	SeedTestSpecialties(t)
	specialtyIDs := GetSpecialtyIDs(t, "Yoga", "Mobility")

	// Create trainer
	trainerToken := createTestUserAndGetToken(e, "trainer2@example.com", "TrainerPass123!", "Bob", "Trainer")

	// Create trainer profile
	e.POST("/api/v1/trainers/profile").
		WithHeader("Authorization", "Bearer "+trainerToken).
		WithJSON(map[string]interface{}{
			"bio":           "Yoga and mobility specialist.",
			"specialty_ids": specialtyIDs,
			"hourly_rate":   60.00,
			"location":      "Los Angeles, CA",
		}).
		Expect().
		Status(201)

	// Create client
	clientToken := createTestUserAndGetToken(e, "client2@example.com", "ClientPass123!", "Alice", "Client")

	// Get client ID
	clientResponse := e.GET("/api/v1/auth/profile").
		WithHeader("Authorization", "Bearer "+clientToken).
		Expect().
		Status(200).
		JSON().
		Object()
	clientID := clientResponse.Value("data").Object().Value("id").String().Raw()

	// Send invitation
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

	t.Run("Reject Invitation", func(t *testing.T) {
		response := e.PUT("/api/v1/me/trainer-invitations/" + invitationID).
			WithHeader("Authorization", "Bearer "+clientToken).
			WithJSON(map[string]interface{}{
				"action": "reject",
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("rejected")
	})

	t.Run("No Pending Invitations After Reject", func(t *testing.T) {
		response := e.GET("/api/v1/me/trainer-invitations").
			WithHeader("Authorization", "Bearer "+clientToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("data").Array().Length().IsEqual(0)
	})

	t.Run("Trainer No Longer Sees Rejected Client", func(t *testing.T) {
		response := e.GET("/api/v1/trainers/clients").
			WithHeader("Authorization", "Bearer "+trainerToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("data").Array().Length().IsEqual(0)
	})
}

func testTrainerClientAuthorization(t *testing.T, e *httpexpect.Expect) {
	// Seed specialties
	SeedTestSpecialties(t)
	specialtyIDs := GetSpecialtyIDs(t, "HIIT", "Cardio")

	// Create two trainers
	trainer1Token := createTestUserAndGetToken(e, "trainer1@example.com", "TrainerPass123!", "Trainer", "One")
	trainer2Token := createTestUserAndGetToken(e, "trainer2@example.com", "TrainerPass123!", "Trainer", "Two")

	// Create client
	clientToken := createTestUserAndGetToken(e, "client@example.com", "ClientPass123!", "Test", "Client")

	// Get client ID
	clientResponse := e.GET("/api/v1/auth/profile").
		WithHeader("Authorization", "Bearer "+clientToken).
		Expect().
		Status(200).
		JSON().
		Object()
	clientID := clientResponse.Value("data").Object().Value("id").String().Raw()

	// Create trainer profile for trainer 1
	e.POST("/api/v1/trainers/profile").
		WithHeader("Authorization", "Bearer "+trainer1Token).
		WithJSON(map[string]interface{}{
			"bio":           "First trainer.",
			"specialty_ids": specialtyIDs,
			"hourly_rate":   50.00,
			"location":      "Chicago, IL",
		}).
		Expect().
		Status(201)

	// Trainer 1 invites client
	inviteResponse := e.POST("/api/v1/trainers/clients").
		WithHeader("Authorization", "Bearer "+trainer1Token).
		WithJSON(map[string]interface{}{
			"client_id": clientID,
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	invitationID := inviteResponse.Value("data").Object().Value("id").String().Raw()

	t.Run("Cannot Invite Self", func(t *testing.T) {
		// Get trainer 1's user ID
		trainer1Response := e.GET("/api/v1/auth/profile").
			WithHeader("Authorization", "Bearer "+trainer1Token).
			Expect().
			Status(200).
			JSON().
			Object()
		trainer1ID := trainer1Response.Value("data").Object().Value("id").String().Raw()

		response := e.POST("/api/v1/trainers/clients").
			WithHeader("Authorization", "Bearer "+trainer1Token).
			WithJSON(map[string]interface{}{
				"client_id": trainer1ID,
			}).
			Expect().
			Status(400).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("yourself")
	})

	t.Run("Other Trainer Cannot Remove Client", func(t *testing.T) {
		// Create profile for trainer 2
		e.POST("/api/v1/trainers/profile").
			WithHeader("Authorization", "Bearer "+trainer2Token).
			WithJSON(map[string]interface{}{
				"bio":           "Second trainer.",
				"specialty_ids": specialtyIDs,
				"hourly_rate":   55.00,
				"location":      "Miami, FL",
			}).
			Expect().
			Status(201)

		// Trainer 2 tries to remove trainer 1's client
		response := e.DELETE("/api/v1/trainers/clients/" + invitationID).
			WithHeader("Authorization", "Bearer "+trainer2Token).
			Expect().
			Status(403).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("own clients")
	})

	t.Run("Other User Cannot Accept Invitation", func(t *testing.T) {
		// Trainer 2 tries to accept invitation meant for client
		response := e.PUT("/api/v1/me/trainer-invitations/" + invitationID).
			WithHeader("Authorization", "Bearer "+trainer2Token).
			WithJSON(map[string]interface{}{
				"action": "accept",
			}).
			Expect().
			Status(403).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("not for you")
	})

	t.Run("Invalid Client ID Returns 404", func(t *testing.T) {
		// Use a valid UUID format that doesn't exist in the database
		response := e.POST("/api/v1/trainers/clients").
			WithHeader("Authorization", "Bearer "+trainer1Token).
			WithJSON(map[string]interface{}{
				"client_id": "11111111-1111-1111-1111-111111111111",
			}).
			Expect().
			Status(404).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("not found")
	})
}
