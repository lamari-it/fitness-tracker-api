package test

import (
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestSpecialtiesEndpoint(t *testing.T) {
	e := SetupTestApp(t)

	t.Run("List Specialties", func(t *testing.T) {
		CleanDatabase(t)
		testListSpecialties(t, e)
	})
}

func testListSpecialties(t *testing.T, e *httpexpect.Expect) {
	// Create a test user and get auth token
	token := createTestUserAndGetToken(e, "user@example.com", "UserPass123!", "Test", "User")

	// Seed specialties
	SeedTestSpecialties(t)

	t.Run("Get All Specialties", func(t *testing.T) {
		response := e.GET("/api/v1/specialties/").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(200).
			JSON().
			Object()

		response.Value("success").Boolean().IsTrue()
		response.Value("message").String().Contains("retrieved")

		data := response.Value("data").Array()
		data.Length().IsEqual(10) // We seed 10 specialties

		// Check that each specialty has required fields
		for _, item := range data.Iter() {
			specialty := item.Object()
			specialty.Value("id").String().NotEmpty()
			specialty.Value("name").String().NotEmpty()
			specialty.Value("description").String().NotEmpty()
			specialty.Value("created_at").String().NotEmpty()
			specialty.Value("updated_at").String().NotEmpty()
		}

		// Check that specialties are sorted alphabetically by name
		firstSpecialty := data.Value(0).Object()
		firstSpecialty.Value("name").String().IsEqual("Bodybuilding")
	})

	t.Run("Specialties Without Auth", func(t *testing.T) {
		// This should still require auth based on current routes
		response := e.GET("/api/v1/specialties/").
			Expect().
			Status(401).
			JSON().
			Object()

		response.Value("success").Boolean().IsFalse()
	})
}
