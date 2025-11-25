package database

import (
	"fit-flow-api/config"
	"fit-flow-api/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	cfg := config.AppConfig

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)

	var logMode logger.Interface
	if cfg.Environment == "dev" {
		logMode = logger.Default.LogMode(logger.Info)
	} else {
		logMode = logger.Default.LogMode(logger.Silent)
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logMode,
	})
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
		// Core user & auth (no dependencies)
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
		&models.RoleInheritance{},
		&models.UserRole{},
		&models.Translation{},

		// Trainer system
		&models.Specialty{},
		&models.TrainerProfile{},
		&models.TrainerSpecialty{},
		&models.TrainerReview{},
		&models.TrainerClientLink{},
		&models.TrainerInvitation{},
		&models.Friendship{},

		// Exercise reference data
		&models.MuscleGroup{},
		&models.Equipment{},
		&models.Exercise{},
		&models.ExerciseMuscleGroup{},
		&models.ExerciseEquipment{},
		&models.UserEquipment{},

		// Fitness reference data
		&models.FitnessLevel{},
		&models.FitnessGoal{},

		// User fitness profile (must be before UserFitnessGoal)
		&models.UserFitnessProfile{},
		&models.UserFitnessGoal{},  // FK to UserFitnessProfile
		&models.WeightLog{},

		// RPE scales
		&models.RPEScale{},
		&models.RPEScaleValue{},

		// Workout plans & workouts
		&models.WorkoutPlan{},
		&models.WorkoutPlanItem{},
		&models.Workout{},
		&models.WorkoutPrescription{},
		&models.PlanEnrollment{},

		// Workout sessions
		&models.WorkoutSession{},
		&models.SessionBlock{},
		&models.SessionExercise{},
		&models.SessionSet{},

		// Social features
		&models.SharedWorkout{},
		&models.WorkoutComment{},
		&models.WorkoutCommentReaction{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
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
		&models.SessionSet{},
		&models.SessionExercise{},
		&models.SessionBlock{},
		&models.WorkoutSession{},
		&models.PlanEnrollment{},
		&models.WorkoutPrescription{},
		&models.WorkoutPlanItem{},
		&models.Workout{},
		&models.WorkoutPlan{},
		&models.RPEScaleValue{},
		&models.RPEScale{},
		&models.UserRole{},
		&models.RoleInheritance{},
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
		// Social
		&models.WorkoutCommentReaction{},
		&models.WorkoutComment{},
		&models.SharedWorkout{},
		// Workout logs
		&models.SessionSet{},
		&models.SessionExercise{},
		&models.SessionBlock{},
		&models.WorkoutSession{},
		// Workout structure
		&models.WorkoutPrescription{},
		&models.WorkoutPlanItem{},
		&models.PlanEnrollment{},
		&models.Workout{},
		&models.WorkoutPlan{},
		// RPE scales (after workout tables that reference them)
		&models.RPEScaleValue{},
		&models.RPEScale{},
		// User fitness
		&models.WeightLog{},
		&models.UserFitnessGoal{},
		&models.UserFitnessProfile{},
		// User equipment
		&models.UserEquipment{},
		// Exercise relationships
		&models.ExerciseEquipment{},
		&models.ExerciseMuscleGroup{},
		&models.Exercise{},
		// Reference data
		&models.Equipment{},
		&models.MuscleGroup{},
		&models.FitnessGoal{},
		&models.FitnessLevel{},
		// Social/Friends
		&models.Friendship{},
		// Trainer
		&models.TrainerInvitation{},
		&models.TrainerClientLink{},
		&models.TrainerReview{},
		&models.TrainerSpecialty{},
		&models.TrainerProfile{},
		&models.Specialty{},
		// RBAC
		&models.UserRole{},
		&models.RoleInheritance{},
		&models.RolePermission{},
		&models.Role{},
		&models.Permission{},
		// Translations
		&models.Translation{},
		// Users (last)
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
