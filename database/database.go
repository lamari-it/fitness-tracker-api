package database

import (
	"fit-flow-api/config"
	"fit-flow-api/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	cfg := config.AppConfig

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")
}

// InitializeDB initializes the database with migrations or AutoMigrate fallback
func InitializeDB() {
	// Check if migrations should be used (recommended for production)
	useMigrations := os.Getenv("USE_MIGRATIONS") == "true"

	if useMigrations {
		log.Println("Using golang-migrate for database schema management")
		if err := MigrateUp(DB); err != nil {
			log.Printf("Migration failed: %v", err)
			log.Println("Falling back to AutoMigrate...")
			AutoMigrate()
		} else {
			log.Println("Database migrations completed successfully")
		}
	} else {
		log.Println("Using GORM AutoMigrate for database schema management")
		AutoMigrate()
	}
}

// AutoMigrate uses GORM's AutoMigrate feature (legacy/development mode)
func AutoMigrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
		&models.UserRole{},
		&models.TrainerProfile{},
		&models.TrainerReview{},
		&models.TrainerClientLink{},
		&models.Friendship{},
		&models.MuscleGroup{},
		&models.Exercise{},
		&models.ExerciseMuscleGroup{},
		&models.Equipment{},
		&models.ExerciseEquipment{},
		&models.UserEquipment{},
		&models.FitnessLevel{},
		&models.FitnessGoal{},
		&models.UserFitnessGoal{},
		&models.WorkoutPlan{},
		&models.WorkoutPlanItem{},
		&models.Workout{},
		&models.WorkoutExercise{},
		&models.PlanEnrollment{},
		&models.WorkoutSession{},
		&models.ExerciseLog{},
		&models.SetLog{},
		&models.SharedWorkout{},
		&models.WorkoutComment{},
		&models.WorkoutCommentReaction{},
		&models.Translation{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Add unique constraint for exercise-muscle group combination
	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS unique_exercise_muscle_combo ON exercise_muscle_groups (exercise_id, muscle_group_id)")

	// Add unique constraint for exercise-equipment combination
	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS unique_exercise_equipment_combo ON exercise_equipment (exercise_id, equipment_id)")

	// Add unique constraint for user-fitness goal combination
	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS unique_user_fitness_goal_combo ON user_fitness_goals (user_id, fitness_goal_id)")

	// Add unique constraint for translation combinations
	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS unique_translation_combo ON translations (resource_type, resource_id, field_name, language)")

	// Add unique constraint for user equipment combinations
	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS unique_user_equipment_combo ON user_equipment (user_id, equipment_id, location_type)")

	log.Println("Database AutoMigrate completed")
}

// Connect creates a new database connection for CLI commands
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// RunSeeders runs all database seeders
func RunSeeders(db *gorm.DB) error {
	// Set the global DB variable for seeders
	DB = db
	SeedDatabase()
	return nil
}

// DropAllData deletes all data from tables but keeps the schema
func DropAllData(db *gorm.DB) error {
	// Set the global DB variable
	DB = db

	// Delete in reverse order to respect foreign key constraints
	tables := []interface{}{
		&models.WorkoutCommentReaction{},
		&models.WorkoutComment{},
		&models.SharedWorkout{},
		&models.SetLog{},
		&models.ExerciseLog{},
		&models.WorkoutSession{},
		&models.PlanEnrollment{},
		&models.WorkoutExercise{},
		&models.WorkoutPlanItem{},
		&models.Workout{},
		&models.WorkoutPlan{},
		&models.UserRole{},
		&models.RolePermission{},
		&models.Role{},
		&models.Permission{},
		&models.UserFitnessGoal{},
		&models.UserEquipment{},
		&models.ExerciseEquipment{},
		&models.ExerciseMuscleGroup{},
		&models.Exercise{},
		&models.Equipment{},
		&models.MuscleGroup{},
		&models.FitnessGoal{},
		&models.FitnessLevel{},
		&models.Friendship{},
		&models.TrainerClientLink{},
		&models.TrainerReview{},
		&models.TrainerProfile{},
		&models.Translation{},
		&models.User{},
	}

	for _, table := range tables {
		if db.Migrator().HasTable(table) {
			if err := db.Unscoped().Where("1=1").Delete(table).Error; err != nil {
				return fmt.Errorf("failed to delete data from table: %w", err)
			}
		}
	}

	log.Println("All data dropped successfully")
	return nil
}

// DropAllTables drops all tables from the database
func DropAllTables(db *gorm.DB) error {
	// Drop tables in reverse order to respect foreign key constraints
	tables := []interface{}{
		&models.WorkoutCommentReaction{},
		&models.WorkoutComment{},
		&models.SharedWorkout{},
		&models.SetLog{},
		&models.ExerciseLog{},
		&models.WorkoutSession{},
		&models.WorkoutExercise{},
		&models.Workout{},
		&models.WorkoutPlan{},
		&models.UserFitnessGoal{},
		&models.UserEquipment{},
		&models.ExerciseEquipment{},
		&models.ExerciseMuscleGroup{},
		&models.Exercise{},
		&models.Equipment{},
		&models.MuscleGroup{},
		&models.FitnessGoal{},
		&models.FitnessLevel{},
		&models.Friendship{},
		&models.TrainerClientLink{},
		&models.TrainerReview{},
		&models.TrainerProfile{},
		&models.Translation{},
		&models.User{},
	}

	if err := db.Migrator().DropTable(tables...); err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}

	log.Println("All tables dropped successfully")
	return nil
}

// Migrate runs pending migrations
func Migrate(db *gorm.DB, cfg *config.Config) error {
	return MigrateUp(db)
}

// MigrateFresh drops all tables and re-runs migrations
func MigrateFresh(db *gorm.DB, cfg *config.Config) error {
	// First, drop all tables
	if err := DropAllTables(db); err != nil {
		return err
	}

	// Then run migrations
	return MigrateUp(db)
}

// RollbackMigration rolls back the last migration
func RollbackMigration(db *gorm.DB, cfg *config.Config) error {
	// Get current version
	version, dirty, err := GetMigrationVersion(db)
	if err != nil {
		return err
	}

	if dirty {
		return fmt.Errorf("database is in dirty state, please fix manually")
	}

	if version == 0 {
		return fmt.Errorf("no migrations to rollback")
	}

	// Rollback to previous version
	return MigrateToVersion(db, version-1)
}
