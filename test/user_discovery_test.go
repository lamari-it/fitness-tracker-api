package test

import (
	"lamari-fit-api/database"
	"lamari-fit-api/models"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestUserDiscoveryEndpoints(t *testing.T) {
	e := SetupTestApp(t)

	t.Run("User Discovery", func(t *testing.T) {
		CleanDatabase(t)
		SeedTestRoles(t)
		testUserDiscovery(t, e)
	})
}

func testUserDiscovery(t *testing.T, e *httpexpect.Expect) {
	// Setup: Create multiple users with different visibility settings

	// User 1: Public profile, looking for trainer
	user1Data := map[string]interface{}{
		"email":            "public_user@example.com",
		"password":         "PublicPass123!",
		"password_confirm": "PublicPass123!",
		"first_name":       "John",
		"last_name":        "Public",
	}
	e.POST("/api/v1/auth/register").WithJSON(user1Data).Expect().Status(201)
	token1 := GetAuthToken(e, "public_user@example.com", "PublicPass123!")

	// Set user1 to public visibility and looking for trainer
	e.PUT("/api/v1/user/settings").
		WithHeader("Authorization", "Bearer "+token1).
		WithJSON(map[string]interface{}{
			"profile_visibility":     "public",
			"is_looking_for_trainer": true,
			"bio":                    "Looking for a fitness trainer!",
		}).
		Expect().
		Status(200)

	// User 2: Private profile
	user2Data := map[string]interface{}{
		"email":            "private_user@example.com",
		"password":         "PrivatePass123!",
		"password_confirm": "PrivatePass123!",
		"first_name":       "Jane",
		"last_name":        "Private",
	}
	e.POST("/api/v1/auth/register").WithJSON(user2Data).Expect().Status(201)
	token2 := GetAuthToken(e, "private_user@example.com", "PrivatePass123!")

	// User2 stays private (default)

	// User 3: Public profile, NOT looking for trainer
	user3Data := map[string]interface{}{
		"email":            "public_no_trainer@example.com",
		"password":         "PublicPass123!",
		"password_confirm": "PublicPass123!",
		"first_name":       "Bob",
		"last_name":        "Smith",
	}
	e.POST("/api/v1/auth/register").WithJSON(user3Data).Expect().Status(201)
	token3 := GetAuthToken(e, "public_no_trainer@example.com", "PublicPass123!")

	// Set user3 to public but not looking for trainer
	e.PUT("/api/v1/user/settings").
		WithHeader("Authorization", "Bearer "+token3).
		WithJSON(map[string]interface{}{
			"profile_visibility":     "public",
			"is_looking_for_trainer": false,
			"bio":                    "Just here to workout",
		}).
		Expect().
		Status(200)

	// User 4: Friends only profile
	user4Data := map[string]interface{}{
		"email":            "friends_only@example.com",
		"password":         "FriendsPass123!",
		"password_confirm": "FriendsPass123!",
		"first_name":       "Alice",
		"last_name":        "FriendsOnly",
	}
	e.POST("/api/v1/auth/register").WithJSON(user4Data).Expect().Status(201)
	token4 := GetAuthToken(e, "friends_only@example.com", "FriendsPass123!")

	e.PUT("/api/v1/user/settings").
		WithHeader("Authorization", "Bearer "+token4).
		WithJSON(map[string]interface{}{
			"profile_visibility": "friends_only",
		}).
		Expect().
		Status(200)

	// Get user IDs for later tests
	var user1, user2, user4 models.User
	database.DB.Where("email = ?", "public_user@example.com").First(&user1)
	database.DB.Where("email = ?", "private_user@example.com").First(&user2)
	database.DB.Where("email = ?", "friends_only@example.com").First(&user4)

	// Tests for DiscoverUsers endpoint

	t.Run("Discover Public Users Returns Only Public Profiles", func(t *testing.T) {
		response := e.GET("/api/v1/users").
			WithHeader("Authorization", "Bearer "+token2).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()

		// Should only see public users (user1 and user3)
		data.Length().IsEqual(2)

		// Verify returned data has limited fields (UserDiscoveryResponse)
		for _, item := range data.Iter() {
			obj := item.Object()
			obj.ContainsKey("id")
			obj.ContainsKey("first_name")
			obj.ContainsKey("last_name")
			obj.ContainsKey("bio")
			obj.ContainsKey("is_looking_for_trainer")
			obj.ContainsKey("created_at")
			// Should NOT have sensitive fields
			obj.NotContainsKey("email")
			obj.NotContainsKey("password")
			obj.NotContainsKey("provider")
		}
	})

	t.Run("Search Users By Name", func(t *testing.T) {
		response := e.GET("/api/v1/users").
			WithHeader("Authorization", "Bearer "+token2).
			WithQuery("search", "John").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()

		// Should find John Public
		data.Length().IsEqual(1)
		data.Value(0).Object().Value("first_name").String().IsEqual("John")
	})

	t.Run("Search Users By Email", func(t *testing.T) {
		response := e.GET("/api/v1/users").
			WithHeader("Authorization", "Bearer "+token2).
			WithQuery("search", "public_user@").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()

		// Should find user with matching email pattern
		data.Length().IsEqual(1)
	})

	t.Run("Filter By Looking For Trainer True", func(t *testing.T) {
		response := e.GET("/api/v1/users").
			WithHeader("Authorization", "Bearer "+token2).
			WithQuery("is_looking_for_trainer", "true").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()

		// Should only find user1 who is looking for trainer
		data.Length().IsEqual(1)
		data.Value(0).Object().Value("is_looking_for_trainer").Boolean().IsTrue()
	})

	t.Run("Filter By Looking For Trainer False", func(t *testing.T) {
		response := e.GET("/api/v1/users").
			WithHeader("Authorization", "Bearer "+token2).
			WithQuery("is_looking_for_trainer", "false").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()

		// Should only find user3 who is not looking for trainer
		data.Length().IsEqual(1)
		data.Value(0).Object().Value("is_looking_for_trainer").Boolean().IsFalse()
	})

	t.Run("Pagination Works", func(t *testing.T) {
		response := e.GET("/api/v1/users").
			WithHeader("Authorization", "Bearer "+token2).
			WithQuery("page", "1").
			WithQuery("limit", "1").
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Array()
		data.Length().IsEqual(1)

		meta := response.Value("meta").Object()
		meta.Value("current_page").Number().IsEqual(1)
		meta.Value("per_page").Number().IsEqual(1)
		meta.Value("total_items").Number().IsEqual(2)
	})

	// Tests for GetUserPublicProfile endpoint

	t.Run("Get Public User Profile Returns Discovery Response", func(t *testing.T) {
		response := e.GET("/api/v1/users/"+user1.ID.String()).
			WithHeader("Authorization", "Bearer "+token2).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()

		// Should return UserDiscoveryResponse (limited fields)
		data.Value("id").String().IsEqual(user1.ID.String())
		data.Value("first_name").String().IsEqual("John")
		data.Value("last_name").String().IsEqual("Public")
		data.Value("bio").String().Contains("Looking for")
		data.Value("is_looking_for_trainer").Boolean().IsTrue()
		// Should NOT have email (it's a discovery response)
		data.NotContainsKey("email")
	})

	t.Run("Get Own Profile Returns Full Response", func(t *testing.T) {
		response := e.GET("/api/v1/users/"+user1.ID.String()).
			WithHeader("Authorization", "Bearer "+token1).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()

		// Should return full UserResponse when viewing self
		data.Value("id").String().IsEqual(user1.ID.String())
		data.Value("email").String().IsEqual("public_user@example.com")
		data.Value("first_name").String().IsEqual("John")
	})

	t.Run("Get Private User Profile Returns 404", func(t *testing.T) {
		e.GET("/api/v1/users/"+user2.ID.String()).
			WithHeader("Authorization", "Bearer "+token1).
			Expect().
			Status(404)
	})

	t.Run("Get Friends Only Profile Without Friendship Returns 404", func(t *testing.T) {
		e.GET("/api/v1/users/"+user4.ID.String()).
			WithHeader("Authorization", "Bearer "+token1).
			Expect().
			Status(404)
	})

	t.Run("Get Friends Only Profile With Friendship Returns Profile", func(t *testing.T) {
		// Create a friendship between user1 and user4
		database.DB.Create(&models.Friendship{
			UserID:   user1.ID,
			FriendID: user4.ID,
			Status:   "accepted",
		})

		response := e.GET("/api/v1/users/"+user4.ID.String()).
			WithHeader("Authorization", "Bearer "+token1).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		data := response.Value("data").Object()
		data.Value("first_name").String().IsEqual("Alice")
	})

	t.Run("Get Nonexistent User Returns 404", func(t *testing.T) {
		e.GET("/api/v1/users/00000000-0000-0000-0000-000000000000").
			WithHeader("Authorization", "Bearer "+token1).
			Expect().
			Status(404)
	})

	t.Run("Discover Users Without Auth Returns 401", func(t *testing.T) {
		e.GET("/api/v1/users").
			Expect().
			Status(401)
	})

	t.Run("Get User Profile Without Auth Returns 401", func(t *testing.T) {
		e.GET("/api/v1/users/" + user1.ID.String()).
			Expect().
			Status(401)
	})
}

func TestUserSettingsPrivacyFields(t *testing.T) {
	e := SetupTestApp(t)

	t.Run("Privacy Fields in User Settings", func(t *testing.T) {
		CleanDatabase(t)
		SeedTestRoles(t)
		testUserSettingsPrivacyFields(t, e)
	})
}

func testUserSettingsPrivacyFields(t *testing.T, e *httpexpect.Expect) {
	// Setup: Create a user
	userData := map[string]interface{}{
		"email":            "privacy_test@example.com",
		"password":         "PrivacyPass123!",
		"password_confirm": "PrivacyPass123!",
		"first_name":       "Privacy",
		"last_name":        "Test",
	}
	e.POST("/api/v1/auth/register").WithJSON(userData).Expect().Status(201)
	token := GetAuthToken(e, "privacy_test@example.com", "PrivacyPass123!")

	t.Run("Default Privacy Settings", func(t *testing.T) {
		response := e.GET("/api/v1/user/settings").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(200).
			JSON().
			Object()

		data := response.Value("data").Object()
		// Check default values
		data.Value("profile_visibility").String().IsEqual("private")
		data.Value("is_looking_for_trainer").Boolean().IsFalse()
		// bio is empty by default and may be omitted from response
	})

	t.Run("Update Profile Visibility to Public", func(t *testing.T) {
		response := e.PUT("/api/v1/user/settings").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"profile_visibility": "public",
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		data := response.Value("data").Object()
		data.Value("profile_visibility").String().IsEqual("public")
	})

	t.Run("Update Profile Visibility to Friends Only", func(t *testing.T) {
		response := e.PUT("/api/v1/user/settings").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"profile_visibility": "friends_only",
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		data := response.Value("data").Object()
		data.Value("profile_visibility").String().IsEqual("friends_only")
	})

	t.Run("Update Is Looking For Trainer", func(t *testing.T) {
		response := e.PUT("/api/v1/user/settings").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"is_looking_for_trainer": true,
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		data := response.Value("data").Object()
		data.Value("is_looking_for_trainer").Boolean().IsTrue()
	})

	t.Run("Update Bio", func(t *testing.T) {
		response := e.PUT("/api/v1/user/settings").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"bio": "I love fitness and working out!",
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		data := response.Value("data").Object()
		data.Value("bio").String().IsEqual("I love fitness and working out!")
	})

	t.Run("Invalid Profile Visibility Rejected", func(t *testing.T) {
		e.PUT("/api/v1/user/settings").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"profile_visibility": "invalid_value",
			}).
			Expect().
			Status(400)
	})

	t.Run("Update All Privacy Fields At Once", func(t *testing.T) {
		response := e.PUT("/api/v1/user/settings").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"profile_visibility":     "public",
				"is_looking_for_trainer": true,
				"bio":                    "Updated bio text",
			}).
			Expect().
			Status(200).
			JSON().
			Object()

		data := response.Value("data").Object()
		data.Value("profile_visibility").String().IsEqual("public")
		data.Value("is_looking_for_trainer").Boolean().IsTrue()
		data.Value("bio").String().IsEqual("Updated bio text")
	})
}
