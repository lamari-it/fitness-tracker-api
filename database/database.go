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
		&models.Workout{},
		&models.WorkoutExercise{},
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