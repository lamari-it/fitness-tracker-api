package database

import (
	"fit-flow-api/config"
	"fit-flow-api/models"
	"fmt"
	"log"

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
	
	// Add unique constraint for translation combinations
	DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS unique_translation_combo ON translations (resource_type, resource_id, field_name, language)")
	
	log.Println("Database migration completed")
}