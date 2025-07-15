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

func SeedFitnessLevels() {
	fitnessLevels := []models.FitnessLevel{
		{Name: "Beginner", Description: "New to fitness or returning after a long break", SortOrder: 1},
		{Name: "Intermediate", Description: "Regular exercise experience with good form knowledge", SortOrder: 2},
		{Name: "Advanced", Description: "Experienced athlete with years of consistent training", SortOrder: 3},
		{Name: "Elite", Description: "Competitive athlete or professional level fitness", SortOrder: 4},
	}

	for _, level := range fitnessLevels {
		var existing models.FitnessLevel
		if err := DB.Where("name = ?", level.Name).First(&existing).Error; err != nil {
			if err := DB.Create(&level).Error; err != nil {
				log.Printf("Failed to create fitness level %s: %v", level.Name, err)
			} else {
				log.Printf("Created fitness level: %s", level.Name)
			}
		} else {
			log.Printf("Fitness level already exists: %s", level.Name)
		}
	}
}

func SeedFitnessGoals() {
	fitnessGoals := []models.FitnessGoal{
		{Name: "Weight Loss", Description: "Reduce body weight and body fat percentage", Category: "body_composition", IconName: "scale"},
		{Name: "Muscle Gain", Description: "Build lean muscle mass and increase strength", Category: "body_composition", IconName: "dumbbell"},
		{Name: "Endurance", Description: "Improve cardiovascular fitness and stamina", Category: "performance", IconName: "running"},
		{Name: "Strength", Description: "Increase maximum strength and power output", Category: "performance", IconName: "weight"},
		{Name: "Flexibility", Description: "Improve range of motion and mobility", Category: "wellness", IconName: "stretch"},
		{Name: "General Fitness", Description: "Overall health and wellness improvement", Category: "wellness", IconName: "heart"},
		{Name: "Athletic Performance", Description: "Sport-specific performance enhancement", Category: "performance", IconName: "trophy"},
		{Name: "Rehabilitation", Description: "Recover from injury or medical condition", Category: "wellness", IconName: "medical"},
		{Name: "Body Recomposition", Description: "Simultaneous fat loss and muscle gain", Category: "body_composition", IconName: "transform"},
		{Name: "Stress Relief", Description: "Mental health and stress management through exercise", Category: "wellness", IconName: "mindfulness"},
	}

	for _, goal := range fitnessGoals {
		var existing models.FitnessGoal
		if err := DB.Where("name = ?", goal.Name).First(&existing).Error; err != nil {
			if err := DB.Create(&goal).Error; err != nil {
				log.Printf("Failed to create fitness goal %s: %v", goal.Name, err)
			} else {
				log.Printf("Created fitness goal: %s", goal.Name)
			}
		} else {
			log.Printf("Fitness goal already exists: %s", goal.Name)
		}
	}
}

func SeedEquipment() {
	equipment := []models.Equipment{
		// Free Weights
		{Slug: "dumbbells", Name: "Dumbbells", Description: "Adjustable or fixed weight dumbbells", Category: "free_weight"},
		{Slug: "barbell", Name: "Barbell", Description: "Olympic or standard barbell", Category: "free_weight"},
		{Slug: "kettlebell", Name: "Kettlebell", Description: "Cast iron or vinyl coated kettlebells", Category: "free_weight"},
		{Slug: "plates", Name: "Weight Plates", Description: "Standard or Olympic weight plates", Category: "free_weight"},
		{Slug: "ez_bar", Name: "EZ Curl Bar", Description: "Curved barbell for arm exercises", Category: "free_weight"},
		
		// Bodyweight Equipment
		{Slug: "pull_up_bar", Name: "Pull-up Bar", Description: "Doorway or wall-mounted pull-up bar", Category: "other"},
		{Slug: "dip_station", Name: "Dip Station", Description: "Parallel bars for dips", Category: "other"},
		{Slug: "gymnastic_rings", Name: "Gymnastic Rings", Description: "Suspension rings for advanced bodyweight training", Category: "other"},
		
		// Resistance Equipment
		{Slug: "resistance_bands", Name: "Resistance Bands", Description: "Elastic bands of varying resistance", Category: "other"},
		{Slug: "suspension_trainer", Name: "Suspension Trainer", Description: "TRX or similar suspension system", Category: "other"},
		
		// Cardio Equipment
		{Slug: "treadmill", Name: "Treadmill", Description: "Motorized running machine", Category: "cardio"},
		{Slug: "stationary_bike", Name: "Stationary Bike", Description: "Indoor cycling bike", Category: "cardio"},
		{Slug: "rowing_machine", Name: "Rowing Machine", Description: "Indoor rowing ergometer", Category: "cardio"},
		{Slug: "elliptical", Name: "Elliptical Machine", Description: "Low-impact cardio machine", Category: "cardio"},
		{Slug: "jump_rope", Name: "Jump Rope", Description: "Speed rope for cardio", Category: "cardio"},
		
		// Machines
		{Slug: "cable_machine", Name: "Cable Machine", Description: "Adjustable cable pulley system", Category: "cable"},
		{Slug: "smith_machine", Name: "Smith Machine", Description: "Guided barbell rack", Category: "machine"},
		{Slug: "leg_press", Name: "Leg Press Machine", Description: "Seated or lying leg press", Category: "machine"},
		{Slug: "lat_pulldown", Name: "Lat Pulldown Machine", Description: "Cable machine for back exercises", Category: "machine"},
		{Slug: "chest_press", Name: "Chest Press Machine", Description: "Seated or lying chest press", Category: "machine"},
		{Slug: "leg_curl", Name: "Leg Curl Machine", Description: "Hamstring curl machine", Category: "machine"},
		{Slug: "leg_extension", Name: "Leg Extension Machine", Description: "Quadriceps extension machine", Category: "machine"},
		
		// Benches and Racks
		{Slug: "flat_bench", Name: "Flat Bench", Description: "Standard flat weight bench", Category: "other"},
		{Slug: "adjustable_bench", Name: "Adjustable Bench", Description: "Incline/decline adjustable bench", Category: "other"},
		{Slug: "squat_rack", Name: "Squat Rack", Description: "Power rack or squat stand", Category: "other"},
		{Slug: "preacher_bench", Name: "Preacher Bench", Description: "Angled bench for bicep curls", Category: "other"},
		
		// Other Equipment
		{Slug: "foam_roller", Name: "Foam Roller", Description: "Myofascial release tool", Category: "other"},
		{Slug: "exercise_mat", Name: "Exercise Mat", Description: "Yoga or workout mat", Category: "other"},
		{Slug: "medicine_ball", Name: "Medicine Ball", Description: "Weighted ball for functional training", Category: "other"},
		{Slug: "ab_wheel", Name: "Ab Wheel", Description: "Core training wheel", Category: "other"},
		{Slug: "bosu_ball", Name: "BOSU Ball", Description: "Balance trainer dome", Category: "other"},
		{Slug: "stability_ball", Name: "Stability Ball", Description: "Large inflatable exercise ball", Category: "other"},
		
		// Bodyweight/No Equipment
		{Slug: "bodyweight", Name: "Bodyweight Only", Description: "No equipment required", Category: "other"},
	}

	for _, equip := range equipment {
		var existing models.Equipment
		if err := DB.Where("slug = ?", equip.Slug).First(&existing).Error; err != nil {
			if err := DB.Create(&equip).Error; err != nil {
				log.Printf("Failed to create equipment %s: %v", equip.Name, err)
			} else {
				log.Printf("Created equipment: %s", equip.Name)
			}
		} else {
			log.Printf("Equipment already exists: %s", equip.Name)
		}
	}
}

func SeedDatabase() {
	log.Println("Starting database seeding...")
	SeedMuscleGroups()
	SeedEquipment()
	SeedExercises()
	SeedFitnessLevels()
	SeedFitnessGoals()
	log.Println("Database seeding completed!")
}