package test

import (
	"fit-flow-api/models"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
)

func TestEmailInvitationEndpoints(t *testing.T) {
	e := SetupTestApp(t)

	t.Run("Email Invitation Flow", func(t *testing.T) {
		CleanDatabase(t)
		testEmailInvitationFlow(t, e)
	})

	t.Run("Email Invitation Validation", func(t *testing.T) {
		CleanDatabase(t)
		testEmailInvitationValidation(t, e)
	})

	t.Run("Token Verification", func(t *testing.T) {
		CleanDatabase(t)
		testTokenVerification(t, e)
	})

	t.Run("Registration With Pending Invitation", func(t *testing.T) {
		CleanDatabase(t)
		testRegistrationWithPendingInvitation(t, e)
	})

	t.Run("Resend Invitation", func(t *testing.T) {
		CleanDatabase(t)
		testResendInvitation(t, e)
	})
}

func testEmailInvitationFlow(t *testing.T, e *httpexpect.Expect) {
	// Seed specialties for trainer profile creation
	SeedTestSpecialties(t)
	specialtyIDs := GetSpecialtyIDs(t, "Strength Training", "Weight Loss")

	// Create trainer user
	trainerToken := createTestUserAndGetToken(e, "trainer@example.com", "TrainerPass123!", "John", "Trainer")

	t.Run("Trainer Must Have Profile To Create Email Invitation", func(t *testing.T) {
		response := e.POST("/api/v1/trainers/email-invitations").
			WithHeader("Authorization", "Bearer "+trainerToken).
			WithJSON(map[string]interface{}{
				"email": "newclient@example.com",
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

	t.Run("Create Email Invitation Successfully", func(t *testing.T) {
		response := e.POST("/api/v1/trainers/email-invitations").
			WithHeader("Authorization", "Bearer "+trainerToken).
			WithJSON(map[string]interface{}{
				"email": "newclient@example.com",
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("sent successfully")

		data := response.Value("data").Object()
		invitationID = data.Value("id").String().Raw()
		data.Value("status").String().IsEqual("pending")
		data.Value("invitee_email").String().IsEqual("newclient@example.com")
		data.Value("expires_at").String().NotEmpty()
	})

	t.Run("Cannot Send Duplicate Email Invitation", func(t *testing.T) {
		response := e.POST("/api/v1/trainers/email-invitations").
			WithHeader("Authorization", "Bearer "+trainerToken).
			WithJSON(map[string]interface{}{
				"email": "newclient@example.com",
			}).
			Expect().
			Status(409).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("pending invitation")
	})

	t.Run("Get Email Invitations", func(t *testing.T) {
		response := e.GET("/api/v1/trainers/email-invitations").
			WithHeader("Authorization", "Bearer "+trainerToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Array()
		data.Length().IsEqual(1)
		data.Value(0).Object().Value("status").String().IsEqual("pending")
		data.Value(0).Object().Value("invitee_email").String().IsEqual("newclient@example.com")
	})

	t.Run("Cancel Email Invitation", func(t *testing.T) {
		response := e.DELETE("/api/v1/trainers/email-invitations/" + invitationID).
			WithHeader("Authorization", "Bearer "+trainerToken).
			Expect().
			Status(204)

		// Response body should be empty for 204
		response.Body().IsEmpty()
	})

	t.Run("No Invitations After Cancel", func(t *testing.T) {
		response := e.GET("/api/v1/trainers/email-invitations").
			WithHeader("Authorization", "Bearer "+trainerToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("data").Array().Length().IsEqual(0)
	})
}

func testEmailInvitationValidation(t *testing.T, e *httpexpect.Expect) {
	// Seed specialties
	SeedTestSpecialties(t)
	specialtyIDs := GetSpecialtyIDs(t, "Yoga", "Mobility")

	// Create trainer with profile
	trainerToken := createTestUserAndGetToken(e, "trainer@example.com", "TrainerPass123!", "Bob", "Trainer")

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

	t.Run("Invalid Email Format", func(t *testing.T) {
		response := e.POST("/api/v1/trainers/email-invitations").
			WithHeader("Authorization", "Bearer "+trainerToken).
			WithJSON(map[string]interface{}{
				"email": "not-an-email",
			}).
			Expect().
			Status(400).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("Empty Email", func(t *testing.T) {
		response := e.POST("/api/v1/trainers/email-invitations").
			WithHeader("Authorization", "Bearer "+trainerToken).
			WithJSON(map[string]interface{}{
				"email": "",
			}).
			Expect().
			Status(400).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("Cannot Invite Self", func(t *testing.T) {
		response := e.POST("/api/v1/trainers/email-invitations").
			WithHeader("Authorization", "Bearer "+trainerToken).
			WithJSON(map[string]interface{}{
				"email": "trainer@example.com",
			}).
			Expect().
			Status(400).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("yourself")
	})

	t.Run("Cannot Invite Existing Client", func(t *testing.T) {
		// Create client user
		clientToken := createTestUserAndGetToken(e, "client@example.com", "ClientPass123!", "Jane", "Client")

		// Get client ID
		clientResponse := e.GET("/api/v1/auth/profile").
			WithHeader("Authorization", "Bearer "+clientToken).
			Expect().
			Status(200).
			JSON().
			Object()
		clientID := clientResponse.Value("data").Object().Value("id").String().Raw()

		// Invite via ID-based invitation first
		e.POST("/api/v1/trainers/clients").
			WithHeader("Authorization", "Bearer "+trainerToken).
			WithJSON(map[string]interface{}{
				"client_id": clientID,
			}).
			Expect().
			Status(201)

		// Now try email invitation to same user
		response := e.POST("/api/v1/trainers/email-invitations").
			WithHeader("Authorization", "Bearer "+trainerToken).
			WithJSON(map[string]interface{}{
				"email": "client@example.com",
			}).
			Expect().
			Status(409).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("Cancel Non-Pending Invitation Fails", func(t *testing.T) {
		// Create an invitation
		createResponse := e.POST("/api/v1/trainers/email-invitations").
			WithHeader("Authorization", "Bearer "+trainerToken).
			WithJSON(map[string]interface{}{
				"email": "newuser@example.com",
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		invitationID := createResponse.Value("data").Object().Value("id").String().Raw()

		// Update invitation status to accepted directly in DB
		testDB.Exec("UPDATE trainer_invitations SET status = 'accepted' WHERE id = ?", invitationID)

		// Try to cancel
		response := e.DELETE("/api/v1/trainers/email-invitations/" + invitationID).
			WithHeader("Authorization", "Bearer "+trainerToken).
			Expect().
			Status(400).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("pending")
	})

	t.Run("Cancel Non-Existent Invitation", func(t *testing.T) {
		response := e.DELETE("/api/v1/trainers/email-invitations/11111111-1111-1111-1111-111111111111").
			WithHeader("Authorization", "Bearer "+trainerToken).
			Expect().
			Status(404).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})
}

func testTokenVerification(t *testing.T, e *httpexpect.Expect) {
	// Seed specialties
	SeedTestSpecialties(t)
	specialtyIDs := GetSpecialtyIDs(t, "HIIT", "Cardio")

	// Create trainer with profile
	trainerToken := createTestUserAndGetToken(e, "trainer@example.com", "TrainerPass123!", "Test", "Trainer")

	e.POST("/api/v1/trainers/profile").
		WithHeader("Authorization", "Bearer "+trainerToken).
		WithJSON(map[string]interface{}{
			"bio":           "HIIT and cardio specialist.",
			"specialty_ids": specialtyIDs,
			"hourly_rate":   80.00,
			"location":      "Chicago, IL",
		}).
		Expect().
		Status(201)

	// Create invitation
	e.POST("/api/v1/trainers/email-invitations").
		WithHeader("Authorization", "Bearer "+trainerToken).
		WithJSON(map[string]interface{}{
			"email": "newuser@example.com",
		}).
		Expect().
		Status(201)

	// Get token from database
	var invitation models.TrainerInvitation
	testDB.Where("invitee_email = ?", "newuser@example.com").First(&invitation)
	token := invitation.InvitationToken

	t.Run("Verify Valid Token", func(t *testing.T) {
		response := e.GET("/api/v1/invitations/verify/" + token).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()

		data := response.Value("data").Object()
		data.Value("valid").Boolean().IsTrue()
		data.Value("invitee_email").String().IsEqual("newuser@example.com")
		data.Value("trainer").Object().Value("first_name").String().IsEqual("Test")
		data.Value("trainer").Object().Value("last_name").String().IsEqual("Trainer")
		data.Value("expires_at").String().NotEmpty()
	})

	t.Run("Verify Invalid Token", func(t *testing.T) {
		response := e.GET("/api/v1/invitations/verify/invalidtoken123").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("data").Object().Value("valid").Boolean().IsFalse()
		response.Value("data").Object().Value("message").String().Contains("invalid")
	})

	t.Run("Verify Expired Token", func(t *testing.T) {
		// Create another invitation and expire it
		e.POST("/api/v1/trainers/email-invitations").
			WithHeader("Authorization", "Bearer "+trainerToken).
			WithJSON(map[string]interface{}{
				"email": "expired@example.com",
			}).
			Expect().
			Status(201)

		// Expire it in the database
		testDB.Exec("UPDATE trainer_invitations SET expires_at = ? WHERE invitee_email = ?",
			time.Now().Add(-24*time.Hour), "expired@example.com")

		// Get token
		var expiredInv models.TrainerInvitation
		testDB.Where("invitee_email = ?", "expired@example.com").First(&expiredInv)

		response := e.GET("/api/v1/invitations/verify/" + expiredInv.InvitationToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("data").Object().Value("valid").Boolean().IsFalse()
		response.Value("data").Object().Value("message").String().Contains("expired")
	})

	t.Run("Verify Already Used Token", func(t *testing.T) {
		// Create another invitation and mark it as accepted
		e.POST("/api/v1/trainers/email-invitations").
			WithHeader("Authorization", "Bearer "+trainerToken).
			WithJSON(map[string]interface{}{
				"email": "used@example.com",
			}).
			Expect().
			Status(201)

		// Mark as accepted
		testDB.Exec("UPDATE trainer_invitations SET status = 'accepted' WHERE invitee_email = ?", "used@example.com")

		// Get token
		var usedInv models.TrainerInvitation
		testDB.Where("invitee_email = ?", "used@example.com").First(&usedInv)

		response := e.GET("/api/v1/invitations/verify/" + usedInv.InvitationToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("data").Object().Value("valid").Boolean().IsFalse()
		response.Value("data").Object().Value("message").String().Contains("already been used")
	})
}

func testRegistrationWithPendingInvitation(t *testing.T, e *httpexpect.Expect) {
	// Seed specialties
	SeedTestSpecialties(t)
	specialtyIDs := GetSpecialtyIDs(t, "Weight Loss", "Functional Fitness")

	// Create trainer with profile
	trainerToken := createTestUserAndGetToken(e, "trainer@example.com", "TrainerPass123!", "Bob", "Smith")

	e.POST("/api/v1/trainers/profile").
		WithHeader("Authorization", "Bearer "+trainerToken).
		WithJSON(map[string]interface{}{
			"bio":           "Nutrition and performance specialist.",
			"specialty_ids": specialtyIDs,
			"hourly_rate":   90.00,
			"location":      "Miami, FL",
		}).
		Expect().
		Status(201)

	// Send email invitation
	e.POST("/api/v1/trainers/email-invitations").
		WithHeader("Authorization", "Bearer "+trainerToken).
		WithJSON(map[string]interface{}{
			"email": "newclient@example.com",
		}).
		Expect().
		Status(201)

	t.Run("Registration Creates Pending Client Link", func(t *testing.T) {
		// Register with the invited email
		response := e.POST("/api/v1/auth/register").
			WithJSON(map[string]interface{}{
				"email":            "newclient@example.com",
				"password":         "ClientPass123!",
				"password_confirm": "ClientPass123!",
				"first_name":       "New",
				"last_name":        "Client",
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		token := response.Value("data").Object().Value("token").String().Raw()

		// Check that client now has a pending trainer invitation
		invitationsResponse := e.GET("/api/v1/me/trainer-invitations").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(200).
			JSON().
			Object()

		data := invitationsResponse.Value("data").Array()
		data.Length().IsEqual(1)
		data.Value(0).Object().Value("status").String().IsEqual("pending")
		data.Value(0).Object().Value("trainer").Object().Value("first_name").String().IsEqual("Bob")
	})

	t.Run("Email Invitation Marked As Accepted", func(t *testing.T) {
		// Check that the email invitation was marked as accepted
		response := e.GET("/api/v1/trainers/email-invitations").
			WithHeader("Authorization", "Bearer "+trainerToken).
			Expect().
			Status(200).
			JSON().
			Object()

		data := response.Value("data").Array()
		data.Length().IsEqual(1)
		data.Value(0).Object().Value("status").String().IsEqual("accepted")
	})

	t.Run("Multiple Trainers Can Invite Same Email", func(t *testing.T) {
		// Create second trainer
		trainer2Token := createTestUserAndGetToken(e, "trainer2@example.com", "TrainerPass123!", "Jane", "Doe")

		e.POST("/api/v1/trainers/profile").
			WithHeader("Authorization", "Bearer "+trainer2Token).
			WithJSON(map[string]interface{}{
				"bio":           "Second trainer.",
				"specialty_ids": specialtyIDs,
				"hourly_rate":   70.00,
				"location":      "Boston, MA",
			}).
			Expect().
			Status(201)

		// Invite a new email
		e.POST("/api/v1/trainers/email-invitations").
			WithHeader("Authorization", "Bearer "+trainer2Token).
			WithJSON(map[string]interface{}{
				"email": "anotherclient@example.com",
			}).
			Expect().
			Status(201)

		// First trainer also invites same email
		e.POST("/api/v1/trainers/email-invitations").
			WithHeader("Authorization", "Bearer "+trainerToken).
			WithJSON(map[string]interface{}{
				"email": "anotherclient@example.com",
			}).
			Expect().
			Status(201)

		// Register with that email
		response := e.POST("/api/v1/auth/register").
			WithJSON(map[string]interface{}{
				"email":            "anotherclient@example.com",
				"password":         "ClientPass123!",
				"password_confirm": "ClientPass123!",
				"first_name":       "Another",
				"last_name":        "Client",
			}).
			Expect().
			Status(201).
			JSON().
			Object()

		token := response.Value("data").Object().Value("token").String().Raw()

		// Should have pending invitations from both trainers
		invitationsResponse := e.GET("/api/v1/me/trainer-invitations").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(200).
			JSON().
			Object()

		invitationsResponse.Value("data").Array().Length().IsEqual(2)
	})
}

func testResendInvitation(t *testing.T, e *httpexpect.Expect) {
	// Seed specialties
	SeedTestSpecialties(t)
	specialtyIDs := GetSpecialtyIDs(t, "Rehabilitation", "Bodybuilding")

	// Create trainer with profile
	trainerToken := createTestUserAndGetToken(e, "trainer@example.com", "TrainerPass123!", "Test", "Trainer")

	e.POST("/api/v1/trainers/profile").
		WithHeader("Authorization", "Bearer "+trainerToken).
		WithJSON(map[string]interface{}{
			"bio":           "Rehabilitation specialist.",
			"specialty_ids": specialtyIDs,
			"hourly_rate":   85.00,
			"location":      "Seattle, WA",
		}).
		Expect().
		Status(201)

	// Create invitation
	createResponse := e.POST("/api/v1/trainers/email-invitations").
		WithHeader("Authorization", "Bearer "+trainerToken).
		WithJSON(map[string]interface{}{
			"email": "client@example.com",
		}).
		Expect().
		Status(201).
		JSON().
		Object()

	invitationID := createResponse.Value("data").Object().Value("id").String().Raw()

	// Get original token
	var originalInv models.TrainerInvitation
	testDB.Where("invitee_email = ?", "client@example.com").First(&originalInv)
	originalToken := originalInv.InvitationToken
	originalExpiry := originalInv.ExpiresAt

	t.Run("Resend Invitation Successfully", func(t *testing.T) {
		response := e.POST("/api/v1/trainers/email-invitations/" + invitationID + "/resend").
			WithHeader("Authorization", "Bearer "+trainerToken).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("resent")
		response.Value("data").Object().Value("status").String().IsEqual("pending")
	})

	t.Run("Token Changed After Resend", func(t *testing.T) {
		var updatedInv models.TrainerInvitation
		testDB.Where("invitee_email = ?", "client@example.com").First(&updatedInv)

		// Token should be different
		if updatedInv.InvitationToken == originalToken {
			t.Error("Token should have changed after resend")
		}

		// Expiry should be extended
		if !updatedInv.ExpiresAt.After(originalExpiry) {
			t.Error("Expiration should be extended after resend")
		}
	})

	t.Run("Cannot Resend Non-Pending Invitation", func(t *testing.T) {
		// Mark as accepted
		testDB.Exec("UPDATE trainer_invitations SET status = 'accepted' WHERE id = ?", invitationID)

		response := e.POST("/api/v1/trainers/email-invitations/" + invitationID + "/resend").
			WithHeader("Authorization", "Bearer "+trainerToken).
			Expect().
			Status(400).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
		response.Value("message").String().Contains("pending")
	})

	t.Run("Cannot Resend Non-Existent Invitation", func(t *testing.T) {
		response := e.POST("/api/v1/trainers/email-invitations/11111111-1111-1111-1111-111111111111/resend").
			WithHeader("Authorization", "Bearer "+trainerToken).
			Expect().
			Status(404).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})

	t.Run("Other Trainer Cannot Resend Invitation", func(t *testing.T) {
		// Restore status to pending
		testDB.Exec("UPDATE trainer_invitations SET status = 'pending' WHERE id = ?", invitationID)

		// Create another trainer
		trainer2Token := createTestUserAndGetToken(e, "trainer2@example.com", "TrainerPass123!", "Other", "Trainer")

		e.POST("/api/v1/trainers/profile").
			WithHeader("Authorization", "Bearer "+trainer2Token).
			WithJSON(map[string]interface{}{
				"bio":           "Another trainer.",
				"specialty_ids": specialtyIDs,
				"hourly_rate":   65.00,
				"location":      "Portland, OR",
			}).
			Expect().
			Status(201)

		response := e.POST("/api/v1/trainers/email-invitations/" + invitationID + "/resend").
			WithHeader("Authorization", "Bearer "+trainer2Token).
			Expect().
			Status(404).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})
}
