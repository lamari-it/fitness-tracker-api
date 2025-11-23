package test

import (
	"fit-flow-api/config"
	"fit-flow-api/database"
	"fit-flow-api/routes"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	testDB *gorm.DB
)

// SetupTestApp creates a test application instance
func SetupTestApp(t *testing.T) *httpexpect.Expect {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Load test environment
	_ = godotenv.Load("../.env.test")
	if os.Getenv("DB_NAME") == "" {
		// Use default test database if .env.test doesn't exist
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_PORT", "5467")
		os.Setenv("DB_USER", "postgres")
		os.Setenv("DB_PASSWORD", "password")
		os.Setenv("DB_NAME", "fitflow_test")
		os.Setenv("DB_SSLMODE", "disable")
		os.Setenv("JWT_SECRET", "test_secret_key_for_testing_only")
		os.Setenv("JWT_EXPIRES_IN", "24h")
		os.Setenv("APP_ENV", "test")
	}

	// Initialize config
	config.LoadConfig()

	// Setup test database
	SetupTestDatabase(t)

	// Create Gin router
	router := gin.New()
	router.Use(gin.Recovery())

	// Setup routes
	routes.SetupRoutes(router)

	// Create test server
	server := httptest.NewServer(router)

	// Create httpexpect instance
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  server.URL,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}

// SetupTestDatabase creates a test database connection
func SetupTestDatabase(t *testing.T) {
	cfg := config.AppConfig
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)

	var err error
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Set global database
	database.DB = testDB

	// Drop all tables and re-run AutoMigrate to ensure clean schema
	// This is necessary when model changes remove columns (AutoMigrate doesn't drop columns)
	if err := database.DropAllTables(testDB); err != nil {
		// Tables might not exist on first run, which is fine
		t.Logf("Note: DropAllTables returned: %v", err)
	}

	// Run migrations
	database.AutoMigrate()

	// Clear all data before tests
	CleanDatabase(t)
}

// CleanDatabase clears all data from the test database
func CleanDatabase(t *testing.T) {
	if testDB == nil {
		t.Fatal("Test database not initialized")
	}

	// Delete all data in reverse order to respect foreign key constraints
	tables := []string{
		"workout_comment_reactions",
		"workout_comments",
		"shared_workouts",
		"set_logs",
		"exercise_logs",
		"workout_sessions",
		"plan_enrollments",
		"workout_exercises",
		"set_groups",
		"workout_plan_items",
		"workouts",
		"workout_plans",
		"user_roles",
		"role_permissions",
		"roles",
		"permissions",
		"user_fitness_goals",
		"user_equipments",
		"exercise_equipments",
		"exercise_muscle_groups",
		"exercises",
		"equipment",
		"muscle_groups",
		"fitness_goals",
		"fitness_levels",
		"friendships",
		"trainer_invitations",
		"trainer_client_links",
		"trainer_reviews",
		"trainer_specialties",
		"trainer_profiles",
		"specialties",
		"translations",
		"users",
	}

	for _, table := range tables {
		if err := testDB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			// Table might not exist, continue
			continue
		}
	}
}

// CreateTestUser creates a test user for authentication tests
func CreateTestUser(t *testing.T) map[string]interface{} {
	userData := map[string]interface{}{
		"email":      "test@example.com",
		"password":   "TestPassword123!",
		"first_name": "Test",
		"last_name":  "User",
	}
	return userData
}

// GetAuthToken performs login and returns the JWT token
func GetAuthToken(e *httpexpect.Expect, email, password string) string {
	response := e.POST("/api/v1/auth/login").
		WithJSON(map[string]interface{}{
			"email":    email,
			"password": password,
		}).
		Expect().
		Status(200).
		JSON().
		Object()

	response.Value("success").Boolean().IsTrue()
	return response.Value("data").Object().Value("token").String().Raw()
}

// SeedTestSpecialties creates test specialties and returns them for use in tests
func SeedTestSpecialties(t *testing.T) *httpexpect.Expect {
	if testDB == nil {
		t.Fatal("Test database not initialized")
	}

	// Seed specialties that tests will use
	database.SeedSpecialties()

	return nil
}

// GetSpecialtyIDs retrieves specialty IDs by names for use in tests
func GetSpecialtyIDs(t *testing.T, names ...string) []string {
	if testDB == nil {
		t.Fatal("Test database not initialized")
	}

	var ids []string
	for _, name := range names {
		var specialty struct {
			ID string
		}
		if err := testDB.Table("specialties").Select("id").Where("name = ?", name).First(&specialty).Error; err != nil {
			t.Fatalf("Failed to get specialty ID for %s: %v", name, err)
		}
		ids = append(ids, specialty.ID)
	}

	return ids
}

// SeedTestFitnessGoals creates test fitness goals for use in tests
func SeedTestFitnessGoals(t *testing.T) {
	if testDB == nil {
		t.Fatal("Test database not initialized")
	}

	// Seed fitness goals that tests will use
	database.SeedFitnessGoals()
}

// GetFitnessGoalIDs retrieves fitness goal IDs by name_slug for use in tests
func GetFitnessGoalIDs(t *testing.T, slugs ...string) []string {
	if testDB == nil {
		t.Fatal("Test database not initialized")
	}

	var ids []string
	for _, slug := range slugs {
		var goal struct {
			ID string
		}
		if err := testDB.Table("fitness_goals").Select("id").Where("name_slug = ?", slug).First(&goal).Error; err != nil {
			t.Fatalf("Failed to get fitness goal ID for %s: %v", slug, err)
		}
		ids = append(ids, goal.ID)
	}

	return ids
}

// SeedTestRoles creates test roles for use in auth tests
func SeedTestRoles(t *testing.T) {
	if testDB == nil {
		t.Fatal("Test database not initialized")
	}

	// Seed roles that auth tests will use
	if err := database.SeedRoles(testDB); err != nil {
		t.Logf("Note: SeedRoles returned: %v", err)
	}
}
