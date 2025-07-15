package database

import (
	"fit-flow-api/models"
	"log"

	"github.com/google/uuid"
)

func SeedMuscleGroups() {
	muscleGroups := []models.MuscleGroup{
		// Upper Body
		{Name: "Chest", Description: "Pectoralis major and minor muscles", Category: "upper"},
		{Name: "Back", Description: "Latissimus dorsi, rhomboids, and trapezius", Category: "upper"},
		{Name: "Shoulders", Description: "Deltoids (anterior, lateral, posterior)", Category: "upper"},
		{Name: "Biceps", Description: "Biceps brachii and brachialis", Category: "upper"},
		{Name: "Triceps", Description: "Triceps brachii", Category: "upper"},
		{Name: "Forearms", Description: "Flexors and extensors of the forearm", Category: "upper"},
		
		// Lower Body
		{Name: "Quadriceps", Description: "Front thigh muscles", Category: "lower"},
		{Name: "Hamstrings", Description: "Back thigh muscles", Category: "lower"},
		{Name: "Glutes", Description: "Gluteus maximus, medius, and minimus", Category: "lower"},
		{Name: "Calves", Description: "Gastrocnemius and soleus", Category: "lower"},
		{Name: "Hip Flexors", Description: "Muscles that flex the hip", Category: "lower"},
		{Name: "Adductors", Description: "Inner thigh muscles", Category: "lower"},
		{Name: "Abductors", Description: "Outer thigh muscles", Category: "lower"},
		
		// Core
		{Name: "Abs", Description: "Rectus abdominis", Category: "core"},
		{Name: "Obliques", Description: "Internal and external obliques", Category: "core"},
		{Name: "Lower Back", Description: "Erector spinae and multifidus", Category: "core"},
		{Name: "Transverse Abdominis", Description: "Deep core stabilizing muscle", Category: "core"},
		
		// Cardio/Full Body
		{Name: "Full Body", Description: "Compound movements targeting multiple muscle groups", Category: "cardio"},
		{Name: "Cardiovascular", Description: "Heart and circulatory system", Category: "cardio"},
	}

	for _, muscleGroup := range muscleGroups {
		var existing models.MuscleGroup
		if err := DB.Where("name = ?", muscleGroup.Name).First(&existing).Error; err != nil {
			if err := DB.Create(&muscleGroup).Error; err != nil {
				log.Printf("Failed to create muscle group %s: %v", muscleGroup.Name, err)
			} else {
				log.Printf("Created muscle group: %s", muscleGroup.Name)
			}
		} else {
			log.Printf("Muscle group already exists: %s", muscleGroup.Name)
		}
	}
}

func SeedExercises() {
	// First, get muscle group IDs
	muscleGroupMap := make(map[string]uuid.UUID)
	var muscleGroups []models.MuscleGroup
	DB.Find(&muscleGroups)
	for _, mg := range muscleGroups {
		muscleGroupMap[mg.Name] = mg.ID
	}

	exercises := []struct {
		Exercise     models.Exercise
		MuscleGroups []struct {
			Name      string
			Primary   bool
			Intensity string
		}
	}{
		{
			Exercise: models.Exercise{
				Name:         "Push-ups",
				Description:  "A bodyweight exercise targeting chest, shoulders, and triceps",
				Equipment:    "none",
				IsBodyweight: true,
				Instructions: "Start in a plank position, lower your body until your chest nearly touches the floor, then push back up.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Chest", true, "high"},
				{"Shoulders", false, "moderate"},
				{"Triceps", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Squats",
				Description:  "A lower body exercise targeting quadriceps, glutes, and hamstrings",
				Equipment:    "none",
				IsBodyweight: true,
				Instructions: "Stand with feet shoulder-width apart, lower your body as if sitting back into a chair, then stand back up.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Quadriceps", true, "high"},
				{"Glutes", false, "high"},
				{"Hamstrings", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Deadlifts",
				Description:  "A compound exercise targeting multiple muscle groups",
				Equipment:    "barbell",
				IsBodyweight: false,
				Instructions: "Stand with feet hip-width apart, bend at hips and knees to lower down, then drive through heels to stand up.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Hamstrings", true, "high"},
				{"Glutes", false, "high"},
				{"Lower Back", false, "high"},
				{"Back", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Pull-ups",
				Description:  "A bodyweight exercise targeting back and biceps",
				Equipment:    "pull-up bar",
				IsBodyweight: true,
				Instructions: "Hang from a bar with palms facing away, pull your body up until your chin clears the bar, then lower down.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Back", true, "high"},
				{"Biceps", false, "high"},
				{"Shoulders", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Planks",
				Description:  "A core stability exercise",
				Equipment:    "none",
				IsBodyweight: true,
				Instructions: "Hold a straight line from head to heels, supporting your body on forearms and toes.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Abs", true, "high"},
				{"Transverse Abdominis", false, "high"},
				{"Lower Back", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Burpees",
				Description:  "A full-body exercise combining squat, push-up, and jump",
				Equipment:    "none",
				IsBodyweight: true,
				Instructions: "Start standing, drop into a squat, kick back into plank, do a push-up, jump feet back to squat, then jump up.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Full Body", true, "high"},
				{"Cardiovascular", false, "high"},
			},
		},
	}

	for _, exerciseData := range exercises {
		var existingExercise models.Exercise
		if err := DB.Where("name = ?", exerciseData.Exercise.Name).First(&existingExercise).Error; err != nil {
			// Create the exercise
			if err := DB.Create(&exerciseData.Exercise).Error; err != nil {
				log.Printf("Failed to create exercise %s: %v", exerciseData.Exercise.Name, err)
				continue
			}
			log.Printf("Created exercise: %s", exerciseData.Exercise.Name)

			// Assign muscle groups
			for _, mgData := range exerciseData.MuscleGroups {
				if muscleGroupID, exists := muscleGroupMap[mgData.Name]; exists {
					assignment := models.ExerciseMuscleGroup{
						ExerciseID:    exerciseData.Exercise.ID,
						MuscleGroupID: muscleGroupID,
						Primary:       mgData.Primary,
						Intensity:     mgData.Intensity,
					}
					
					if err := DB.Create(&assignment).Error; err != nil {
						log.Printf("Failed to assign muscle group %s to exercise %s: %v", mgData.Name, exerciseData.Exercise.Name, err)
					} else {
						log.Printf("Assigned muscle group %s to exercise %s", mgData.Name, exerciseData.Exercise.Name)
					}
				}
			}
		} else {
			log.Printf("Exercise already exists: %s", exerciseData.Exercise.Name)
		}
	}
}

func SeedDatabase() {
	log.Println("Starting database seeding...")
	SeedMuscleGroups()
	SeedExercises()
	log.Println("Database seeding completed!")
}