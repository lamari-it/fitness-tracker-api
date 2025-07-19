package database

import (
	"fit-flow-api/models"
	"log"
	"strings"

	"github.com/google/uuid"
)

// generateSlug creates a URL-friendly slug from a name
func generateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)
	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove apostrophes
	slug = strings.ReplaceAll(slug, "'", "")
	// Replace multiple hyphens with single hyphen
	slug = strings.ReplaceAll(slug, "--", "-")
	return slug
}

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
		// Chest Exercises
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
				Name:         "Bench Press",
				Description:  "Classic chest exercise using a barbell",
				IsBodyweight: false,
				Instructions: "Lie on bench, grip barbell slightly wider than shoulders, lower to chest, press up.",
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
				Name:         "Incline Bench Press",
				Description:  "Chest exercise targeting upper pecs",
				IsBodyweight: false,
				Instructions: "On incline bench, press barbell from upper chest to arms length.",
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
				Name:         "Decline Bench Press",
				Description:  "Chest exercise targeting lower pecs",
				IsBodyweight: false,
				Instructions: "On decline bench, press barbell from lower chest to arms length.",
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
				Name:         "Dumbbell Bench Press",
				Description:  "Chest exercise using dumbbells",
				IsBodyweight: false,
				Instructions: "Lie on bench with dumbbells, lower to chest level, press up.",
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
				Name:         "Dumbbell Flyes",
				Description:  "Chest isolation exercise",
				IsBodyweight: false,
				Instructions: "Lie on bench, lower dumbbells in arc motion, squeeze chest to bring back up.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Chest", true, "high"},
				{"Shoulders", false, "low"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Chest Dips",
				Description:  "Bodyweight chest exercise",
				IsBodyweight: true,
				Instructions: "On dip bars, lean forward, lower body, push back up.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Chest", true, "high"},
				{"Triceps", false, "moderate"},
				{"Shoulders", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Cable Chest Flyes",
				Description:  "Cable chest isolation exercise",
				IsBodyweight: false,
				Instructions: "Stand between cables, bring handles together in arc motion.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Chest", true, "high"},
				{"Shoulders", false, "low"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Pec Deck",
				Description:  "Machine chest fly exercise",
				IsBodyweight: false,
				Instructions: "Sit on machine, bring pads together in front of chest.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Chest", true, "high"},
			},
		},

		// Back Exercises
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
				Name:         "Chin-ups",
				Description:  "Bodyweight exercise with underhand grip",
				IsBodyweight: true,
				Instructions: "Hang from bar with palms facing you, pull up until chin clears bar.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Back", true, "high"},
				{"Biceps", false, "high"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Lat Pulldowns",
				Description:  "Cable exercise targeting the lats",
				IsBodyweight: false,
				Instructions: "Sit at lat pulldown machine, pull bar down to chest level.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Back", true, "high"},
				{"Biceps", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Barbell Rows",
				Description:  "Bent-over rowing exercise",
				IsBodyweight: false,
				Instructions: "Bend over holding barbell, pull to lower chest, squeeze shoulder blades.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Back", true, "high"},
				{"Biceps", false, "moderate"},
				{"Shoulders", false, "low"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Dumbbell Rows",
				Description:  "Single-arm rowing exercise",
				IsBodyweight: false,
				Instructions: "Support body with one hand, row dumbbell with other to hip level.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Back", true, "high"},
				{"Biceps", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "T-Bar Rows",
				Description:  "Rowing exercise using T-bar",
				IsBodyweight: false,
				Instructions: "Straddle T-bar, pull handle to chest, squeeze shoulder blades.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Back", true, "high"},
				{"Biceps", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Cable Rows",
				Description:  "Seated cable rowing exercise",
				IsBodyweight: false,
				Instructions: "Sit at cable machine, pull handle to lower chest, squeeze shoulder blades.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Back", true, "high"},
				{"Biceps", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Inverted Rows",
				Description:  "Bodyweight rowing exercise",
				IsBodyweight: true,
				Instructions: "Hang under bar, pull chest to bar, keep body straight.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Back", true, "high"},
				{"Biceps", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Shrugs",
				Description:  "Trapezius exercise",
				IsBodyweight: false,
				Instructions: "Hold weight, lift shoulders straight up, squeeze at top.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Back", true, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Face Pulls",
				Description:  "Rear delt and upper back exercise",
				IsBodyweight: false,
				Instructions: "Pull cable to face level, separate hands, squeeze rear delts.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Back", true, "moderate"},
				{"Shoulders", false, "moderate"},
			},
		},

		// Leg Exercises
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
				Name:         "Barbell Squats",
				Description:  "Loaded squat exercise",
				IsBodyweight: false,
				Instructions: "Bar on shoulders, squat down keeping chest up, drive through heels.",
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
				Name:         "Front Squats",
				Description:  "Squat with bar in front position",
				IsBodyweight: false,
				Instructions: "Bar across front delts, squat down keeping torso upright.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Quadriceps", true, "high"},
				{"Glutes", false, "high"},
				{"Abs", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Goblet Squats",
				Description:  "Squat holding weight at chest",
				IsBodyweight: false,
				Instructions: "Hold dumbbell at chest, squat down keeping torso upright.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Quadriceps", true, "high"},
				{"Glutes", false, "high"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Lunges",
				Description:  "Single-leg exercise",
				IsBodyweight: true,
				Instructions: "Step forward, lower back knee toward ground, push back to start.",
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
				Name:         "Bulgarian Split Squats",
				Description:  "Rear-foot elevated split squat",
				IsBodyweight: true,
				Instructions: "Rear foot elevated, lower front leg into lunge position.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Quadriceps", true, "high"},
				{"Glutes", false, "high"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Leg Press",
				Description:  "Machine leg exercise",
				IsBodyweight: false,
				Instructions: "Sit on machine, lower weight to 90 degrees, press back up.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Quadriceps", true, "high"},
				{"Glutes", false, "high"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Leg Extensions",
				Description:  "Quadriceps isolation exercise",
				IsBodyweight: false,
				Instructions: "Sit on machine, extend legs against resistance.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Quadriceps", true, "high"},
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
				Name:         "Romanian Deadlifts",
				Description:  "Hip-hinge deadlift variation",
				IsBodyweight: false,
				Instructions: "Keep legs relatively straight, hinge at hips, lower bar to mid-shin.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Hamstrings", true, "high"},
				{"Glutes", false, "high"},
				{"Lower Back", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Sumo Deadlifts",
				Description:  "Wide-stance deadlift variation",
				IsBodyweight: false,
				Instructions: "Wide stance, hands inside legs, pull bar up close to body.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Hamstrings", true, "high"},
				{"Glutes", false, "high"},
				{"Quadriceps", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Stiff Leg Deadlifts",
				Description:  "Hamstring-focused deadlift",
				IsBodyweight: false,
				Instructions: "Keep legs straight, hinge at hips, feel stretch in hamstrings.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Hamstrings", true, "high"},
				{"Glutes", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Leg Curls",
				Description:  "Hamstring isolation exercise",
				IsBodyweight: false,
				Instructions: "Lie on machine, curl heels toward glutes against resistance.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Hamstrings", true, "high"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Hip Thrusts",
				Description:  "Glute-focused exercise",
				IsBodyweight: false,
				Instructions: "Shoulders on bench, drive hips up, squeeze glutes at top.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Glutes", true, "high"},
				{"Hamstrings", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Glute Bridges",
				Description:  "Bodyweight glute exercise",
				IsBodyweight: true,
				Instructions: "Lie on back, drive hips up, squeeze glutes at top.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Glutes", true, "high"},
				{"Hamstrings", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Calf Raises",
				Description:  "Calf muscle exercise",
				IsBodyweight: true,
				Instructions: "Rise up on toes, squeeze calves at top, lower slowly.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Calves", true, "high"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Seated Calf Raises",
				Description:  "Seated calf exercise",
				IsBodyweight: false,
				Instructions: "Sit with weight on thighs, rise up on toes, squeeze calves.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Calves", true, "high"},
			},
		},

		// Shoulder Exercises
		{
			Exercise: models.Exercise{
				Name:         "Overhead Press",
				Description:  "Standing shoulder press",
				IsBodyweight: false,
				Instructions: "Press barbell from shoulders to overhead, keep core tight.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Shoulders", true, "high"},
				{"Triceps", false, "moderate"},
				{"Abs", false, "low"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Dumbbell Shoulder Press",
				Description:  "Seated or standing dumbbell press",
				IsBodyweight: false,
				Instructions: "Press dumbbells from shoulder level to overhead.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Shoulders", true, "high"},
				{"Triceps", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Lateral Raises",
				Description:  "Side deltoid isolation",
				IsBodyweight: false,
				Instructions: "Raise dumbbells to sides until parallel to floor.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Shoulders", true, "high"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Rear Delt Flyes",
				Description:  "Rear deltoid isolation",
				IsBodyweight: false,
				Instructions: "Bend over, raise dumbbells to sides, squeeze rear delts.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Shoulders", true, "high"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Front Raises",
				Description:  "Front deltoid isolation",
				IsBodyweight: false,
				Instructions: "Raise dumbbell in front of body to shoulder height.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Shoulders", true, "high"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Arnold Press",
				Description:  "Rotating shoulder press",
				IsBodyweight: false,
				Instructions: "Start with palms facing you, rotate while pressing overhead.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Shoulders", true, "high"},
				{"Triceps", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Upright Rows",
				Description:  "Shoulder and trap exercise",
				IsBodyweight: false,
				Instructions: "Pull bar up to chest level, elbows high.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Shoulders", true, "high"},
				{"Back", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Pike Push-ups",
				Description:  "Bodyweight shoulder exercise",
				IsBodyweight: true,
				Instructions: "In downward dog position, lower head toward ground, push back up.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Shoulders", true, "high"},
				{"Triceps", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Handstand Push-ups",
				Description:  "Advanced bodyweight shoulder exercise",
				IsBodyweight: true,
				Instructions: "In handstand position, lower head toward ground, push back up.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Shoulders", true, "high"},
				{"Triceps", false, "moderate"},
				{"Abs", false, "moderate"},
			},
		},

		// Arm Exercises
		{
			Exercise: models.Exercise{
				Name:         "Bicep Curls",
				Description:  "Basic bicep exercise",
				IsBodyweight: false,
				Instructions: "Hold dumbbells, curl up to shoulders, squeeze biceps.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Biceps", true, "high"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Hammer Curls",
				Description:  "Neutral grip bicep exercise",
				IsBodyweight: false,
				Instructions: "Hold dumbbells with neutral grip, curl up to shoulders.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Biceps", true, "high"},
				{"Forearms", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Preacher Curls",
				Description:  "Bicep exercise on preacher bench",
				IsBodyweight: false,
				Instructions: "Sit at preacher bench, curl weight up, control negative.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Biceps", true, "high"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Concentration Curls",
				Description:  "Isolated bicep exercise",
				IsBodyweight: false,
				Instructions: "Sit, elbow on thigh, curl dumbbell up, focus on bicep.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Biceps", true, "high"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Tricep Dips",
				Description:  "Bodyweight tricep exercise",
				IsBodyweight: true,
				Instructions: "On parallel bars or chair, lower body, push back up.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Triceps", true, "high"},
				{"Shoulders", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Tricep Pushdowns",
				Description:  "Cable tricep exercise",
				IsBodyweight: false,
				Instructions: "At cable machine, push rope down, squeeze triceps.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Triceps", true, "high"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Overhead Tricep Extension",
				Description:  "Tricep exercise overhead",
				IsBodyweight: false,
				Instructions: "Hold weight overhead, lower behind head, extend back up.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Triceps", true, "high"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Close-Grip Bench Press",
				Description:  "Tricep-focused bench press",
				IsBodyweight: false,
				Instructions: "Narrow grip on barbell, press focusing on triceps.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Triceps", true, "high"},
				{"Chest", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Diamond Push-ups",
				Description:  "Tricep-focused push-up variation",
				IsBodyweight: true,
				Instructions: "Hands in diamond shape, push up focusing on triceps.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Triceps", true, "high"},
				{"Chest", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Wrist Curls",
				Description:  "Forearm exercise",
				IsBodyweight: false,
				Instructions: "Forearms on thighs, curl wrists up and down.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Forearms", true, "high"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Reverse Curls",
				Description:  "Forearm and bicep exercise",
				IsBodyweight: false,
				Instructions: "Overhand grip, curl barbell up, focus on forearms.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Forearms", true, "high"},
				{"Biceps", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Farmer's Walks",
				Description:  "Grip and forearm exercise",
				IsBodyweight: false,
				Instructions: "Hold heavy weights, walk maintaining good posture.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Forearms", true, "high"},
				{"Back", false, "moderate"},
				{"Abs", false, "moderate"},
			},
		},

		// Core Exercises
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
				Name:         "Crunches",
				Description:  "Basic abdominal exercise",
				IsBodyweight: true,
				Instructions: "Lie on back, curl shoulders toward knees, squeeze abs.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Abs", true, "high"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Bicycle Crunches",
				Description:  "Dynamic ab exercise",
				IsBodyweight: true,
				Instructions: "Lie on back, alternate bringing elbow to opposite knee.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Abs", true, "high"},
				{"Obliques", false, "high"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Russian Twists",
				Description:  "Oblique exercise",
				IsBodyweight: true,
				Instructions: "Sit with feet up, twist torso side to side.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Obliques", true, "high"},
				{"Abs", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Side Planks",
				Description:  "Oblique stability exercise",
				IsBodyweight: true,
				Instructions: "Lie on side, prop up on elbow, hold straight line.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Obliques", true, "high"},
				{"Abs", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Leg Raises",
				Description:  "Lower ab exercise",
				IsBodyweight: true,
				Instructions: "Lie on back, raise legs straight up, lower slowly.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Abs", true, "high"},
				{"Hip Flexors", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Hanging Leg Raises",
				Description:  "Advanced ab exercise",
				IsBodyweight: true,
				Instructions: "Hang from bar, raise legs up, control the movement.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Abs", true, "high"},
				{"Hip Flexors", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Mountain Climbers",
				Description:  "Dynamic core exercise",
				IsBodyweight: true,
				Instructions: "In plank position, alternate bringing knees to chest.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Abs", true, "high"},
				{"Cardiovascular", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Dead Bug",
				Description:  "Core stability exercise",
				IsBodyweight: true,
				Instructions: "Lie on back, extend opposite arm and leg, alternate.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Abs", true, "high"},
				{"Transverse Abdominis", false, "high"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Hollow Body Hold",
				Description:  "Isometric core exercise",
				IsBodyweight: true,
				Instructions: "Lie on back, press lower back down, hold hollow position.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Abs", true, "high"},
				{"Transverse Abdominis", false, "high"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Ab Wheel Rollouts",
				Description:  "Advanced core exercise",
				IsBodyweight: false,
				Instructions: "Kneel with ab wheel, roll forward maintaining core tension.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Abs", true, "high"},
				{"Lower Back", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Superman",
				Description:  "Lower back exercise",
				IsBodyweight: true,
				Instructions: "Lie face down, lift chest and legs off ground.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Lower Back", true, "high"},
				{"Glutes", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Good Mornings",
				Description:  "Lower back and hamstring exercise",
				IsBodyweight: false,
				Instructions: "Barbell on shoulders, hinge at hips, return to standing.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Lower Back", true, "high"},
				{"Hamstrings", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Hyperextensions",
				Description:  "Lower back exercise on machine",
				IsBodyweight: false,
				Instructions: "On hyperextension bench, lower torso, raise back up.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Lower Back", true, "high"},
				{"Glutes", false, "moderate"},
			},
		},

		// Full Body and Cardio
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
		{
			Exercise: models.Exercise{
				Name:         "Thrusters",
				Description:  "Full-body exercise combining squat and press",
				IsBodyweight: false,
				Instructions: "Squat with dumbbells, stand and press overhead in one motion.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Full Body", true, "high"},
				{"Cardiovascular", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Man Makers",
				Description:  "Complex full-body exercise",
				IsBodyweight: false,
				Instructions: "Dumbbell burpee with rows, squat to press combination.",
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
		{
			Exercise: models.Exercise{
				Name:         "Turkish Get-ups",
				Description:  "Full-body functional exercise",
				IsBodyweight: false,
				Instructions: "Complex movement from lying to standing while holding weight.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Full Body", true, "high"},
				{"Abs", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Kettlebell Swings",
				Description:  "Hip hinge exercise with kettlebell",
				IsBodyweight: false,
				Instructions: "Swing kettlebell between legs, drive hips forward to shoulder height.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Glutes", true, "high"},
				{"Hamstrings", false, "high"},
				{"Cardiovascular", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Kettlebell Snatches",
				Description:  "Explosive full-body exercise",
				IsBodyweight: false,
				Instructions: "Swing kettlebell from between legs to overhead in one motion.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Full Body", true, "high"},
				{"Cardiovascular", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Kettlebell Clean and Press",
				Description:  "Two-part kettlebell exercise",
				IsBodyweight: false,
				Instructions: "Clean kettlebell to shoulder, then press overhead.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Full Body", true, "high"},
				{"Shoulders", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Box Jumps",
				Description:  "Plyometric jumping exercise",
				IsBodyweight: true,
				Instructions: "Jump onto box, land softly, step down.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Quadriceps", true, "high"},
				{"Glutes", false, "high"},
				{"Calves", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Jump Squats",
				Description:  "Plyometric squat variation",
				IsBodyweight: true,
				Instructions: "Squat down, explode up into jump, land softly.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Quadriceps", true, "high"},
				{"Glutes", false, "high"},
				{"Calves", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "High Knees",
				Description:  "Cardio exercise",
				IsBodyweight: true,
				Instructions: "Run in place bringing knees up to chest level.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Cardiovascular", true, "high"},
				{"Hip Flexors", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Jumping Jacks",
				Description:  "Basic cardio exercise",
				IsBodyweight: true,
				Instructions: "Jump feet apart while raising arms, jump back to start.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Cardiovascular", true, "high"},
				{"Shoulders", false, "low"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Bear Crawls",
				Description:  "Full-body crawling exercise",
				IsBodyweight: true,
				Instructions: "Crawl forward on hands and feet, keep knees slightly off ground.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Abs", true, "high"},
				{"Shoulders", false, "moderate"},
				{"Quadriceps", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Crab Walks",
				Description:  "Reverse crawling exercise",
				IsBodyweight: true,
				Instructions: "Walk backward on hands and feet, belly facing up.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Triceps", true, "high"},
				{"Shoulders", false, "moderate"},
				{"Abs", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Treadmill Running",
				Description:  "Cardio exercise on treadmill",
				IsBodyweight: true,
				Instructions: "Run at steady pace or intervals on treadmill.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Cardiovascular", true, "high"},
				{"Quadriceps", false, "moderate"},
				{"Calves", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Stationary Bike",
				Description:  "Cardio exercise on bike",
				IsBodyweight: true,
				Instructions: "Pedal at steady pace or intervals on stationary bike.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Cardiovascular", true, "high"},
				{"Quadriceps", false, "moderate"},
				{"Calves", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Rowing Machine",
				Description:  "Full-body cardio exercise",
				IsBodyweight: true,
				Instructions: "Row with proper form, drive with legs, pull with arms.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Cardiovascular", true, "high"},
				{"Back", false, "moderate"},
				{"Quadriceps", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Elliptical Machine",
				Description:  "Low-impact cardio exercise",
				IsBodyweight: true,
				Instructions: "Move arms and legs in elliptical motion.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Cardiovascular", true, "high"},
				{"Quadriceps", false, "moderate"},
				{"Shoulders", false, "low"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Jump Rope",
				Description:  "Cardio exercise with rope",
				IsBodyweight: true,
				Instructions: "Jump over rope with both feet, maintain rhythm.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Cardiovascular", true, "high"},
				{"Calves", false, "moderate"},
				{"Shoulders", false, "low"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Stair Climbing",
				Description:  "Cardio exercise on stairs",
				IsBodyweight: true,
				Instructions: "Climb stairs at steady pace, use handrails for balance only.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Cardiovascular", true, "high"},
				{"Quadriceps", false, "high"},
				{"Glutes", false, "moderate"},
			},
		},

		// Olympic Lifts and Variations
		{
			Exercise: models.Exercise{
				Name:         "Clean and Jerk",
				Description:  "Olympic lifting exercise",
				IsBodyweight: false,
				Instructions: "Clean barbell to shoulders, then jerk overhead.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Full Body", true, "high"},
				{"Shoulders", false, "high"},
				{"Quadriceps", false, "high"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Snatch",
				Description:  "Olympic lifting exercise",
				IsBodyweight: false,
				Instructions: "Pull barbell from floor to overhead in one motion.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Full Body", true, "high"},
				{"Shoulders", false, "high"},
				{"Back", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Power Clean",
				Description:  "Explosive pulling exercise",
				IsBodyweight: false,
				Instructions: "Pull barbell from floor to shoulder level explosively.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Full Body", true, "high"},
				{"Back", false, "moderate"},
				{"Quadriceps", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Hang Clean",
				Description:  "Clean from hang position",
				IsBodyweight: false,
				Instructions: "From hang position, clean barbell to shoulders.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Full Body", true, "high"},
				{"Back", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Push Press",
				Description:  "Explosive overhead press",
				IsBodyweight: false,
				Instructions: "Use leg drive to help press barbell overhead.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Shoulders", true, "high"},
				{"Quadriceps", false, "moderate"},
				{"Triceps", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Push Jerk",
				Description:  "Explosive overhead lift",
				IsBodyweight: false,
				Instructions: "Drive barbell overhead, drop under bar, stand up.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Shoulders", true, "high"},
				{"Quadriceps", false, "moderate"},
				{"Triceps", false, "moderate"},
			},
		},

		// Isometric and Stability Exercises
		{
			Exercise: models.Exercise{
				Name:         "Wall Sit",
				Description:  "Isometric leg exercise",
				IsBodyweight: true,
				Instructions: "Sit against wall with thighs parallel to floor, hold position.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Quadriceps", true, "high"},
				{"Glutes", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Glute Bridge Hold",
				Description:  "Isometric glute exercise",
				IsBodyweight: true,
				Instructions: "Hold bridge position, squeeze glutes throughout.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Glutes", true, "high"},
				{"Hamstrings", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Single-Leg Glute Bridge",
				Description:  "Unilateral glute exercise",
				IsBodyweight: true,
				Instructions: "Bridge with one leg, hold or perform reps.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Glutes", true, "high"},
				{"Hamstrings", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Single-Leg Deadlift",
				Description:  "Unilateral balance exercise",
				IsBodyweight: true,
				Instructions: "Balance on one leg, hinge at hip, reach toward ground.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Hamstrings", true, "high"},
				{"Glutes", false, "high"},
				{"Abs", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Pistol Squats",
				Description:  "Single-leg squat exercise",
				IsBodyweight: true,
				Instructions: "Squat on one leg, other leg extended forward.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Quadriceps", true, "high"},
				{"Glutes", false, "high"},
				{"Abs", false, "moderate"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Single-Leg Calf Raises",
				Description:  "Unilateral calf exercise",
				IsBodyweight: true,
				Instructions: "Rise up on one foot, squeeze calf at top.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Calves", true, "high"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Bird Dog",
				Description:  "Core stability exercise",
				IsBodyweight: true,
				Instructions: "On hands and knees, extend opposite arm and leg.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Abs", true, "high"},
				{"Lower Back", false, "moderate"},
				{"Glutes", false, "low"},
			},
		},

		// Stretching and Mobility
		{
			Exercise: models.Exercise{
				Name:         "Cat-Cow Stretch",
				Description:  "Spinal mobility exercise",
				IsBodyweight: true,
				Instructions: "On hands and knees, arch and round spine alternately.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Lower Back", true, "low"},
				{"Abs", false, "low"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Child's Pose",
				Description:  "Relaxation and stretch pose",
				IsBodyweight: true,
				Instructions: "Kneel, sit back on heels, reach arms forward.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Lower Back", true, "low"},
				{"Shoulders", false, "low"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Downward Dog",
				Description:  "Full-body stretch",
				IsBodyweight: true,
				Instructions: "Inverted V position, stretch hamstrings and calves.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Hamstrings", true, "low"},
				{"Calves", false, "low"},
				{"Shoulders", false, "low"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Pigeon Pose",
				Description:  "Hip flexibility exercise",
				IsBodyweight: true,
				Instructions: "One leg forward bent, other leg back straight.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Hip Flexors", true, "low"},
				{"Glutes", false, "low"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Cobra Stretch",
				Description:  "Back extension stretch",
				IsBodyweight: true,
				Instructions: "Lie face down, push up with arms, arch back.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Abs", true, "low"},
				{"Hip Flexors", false, "low"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Figure-4 Stretch",
				Description:  "Hip and glute stretch",
				IsBodyweight: true,
				Instructions: "Lie on back, ankle on opposite knee, pull thigh.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Glutes", true, "low"},
				{"Hip Flexors", false, "low"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Seated Forward Fold",
				Description:  "Hamstring stretch",
				IsBodyweight: true,
				Instructions: "Sit with legs extended, reach toward toes.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Hamstrings", true, "low"},
				{"Lower Back", false, "low"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Standing Quad Stretch",
				Description:  "Quadriceps stretch",
				IsBodyweight: true,
				Instructions: "Stand, pull heel toward glutes, hold.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Quadriceps", true, "low"},
				{"Hip Flexors", false, "low"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Standing Calf Stretch",
				Description:  "Calf muscle stretch",
				IsBodyweight: true,
				Instructions: "Step back, keep heel down, lean forward.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Calves", true, "low"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Shoulder Rolls",
				Description:  "Shoulder mobility exercise",
				IsBodyweight: true,
				Instructions: "Roll shoulders forward and backward in circles.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Shoulders", true, "low"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Arm Circles",
				Description:  "Shoulder warm-up exercise",
				IsBodyweight: true,
				Instructions: "Make circles with arms, forward and backward.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Shoulders", true, "low"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Neck Rolls",
				Description:  "Neck mobility exercise",
				IsBodyweight: true,
				Instructions: "Gently roll head in circles, both directions.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Shoulders", true, "low"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Hip Circles",
				Description:  "Hip mobility exercise",
				IsBodyweight: true,
				Instructions: "Hands on hips, rotate hips in circles.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Hip Flexors", true, "low"},
				{"Glutes", false, "low"},
			},
		},
		{
			Exercise: models.Exercise{
				Name:         "Leg Swings",
				Description:  "Dynamic hip stretch",
				IsBodyweight: true,
				Instructions: "Swing leg forward and back, side to side.",
			},
			MuscleGroups: []struct {
				Name      string
				Primary   bool
				Intensity string
			}{
				{"Hip Flexors", true, "low"},
				{"Hamstrings", false, "low"},
			},
		},
	}

	for _, exerciseData := range exercises {
		// Generate slug if not provided
		if exerciseData.Exercise.Slug == "" {
			exerciseData.Exercise.Slug = generateSlug(exerciseData.Exercise.Name)
		}
		
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
		{Slug: "trap_bar", Name: "Trap Bar", Description: "Hexagonal deadlift bar", Category: "free_weight"},
		{Slug: "curl_bar", Name: "Straight Curl Bar", Description: "Straight barbell for bicep exercises", Category: "free_weight"},
		{Slug: "swiss_bar", Name: "Swiss Bar", Description: "Multi-grip football bar", Category: "free_weight"},
		{Slug: "safety_squat_bar", Name: "Safety Squat Bar", Description: "Padded squat bar with handles", Category: "free_weight"},
		{Slug: "cambered_bar", Name: "Cambered Bar", Description: "Curved barbell for squats", Category: "free_weight"},
		{Slug: "log_bar", Name: "Log Bar", Description: "Strongman log for pressing", Category: "free_weight"},
		{Slug: "axle_bar", Name: "Axle Bar", Description: "Thick strongman bar", Category: "free_weight"},

		// Bodyweight Equipment
		{Slug: "pull_up_bar", Name: "Pull-up Bar", Description: "Doorway or wall-mounted pull-up bar", Category: "other"},
		{Slug: "dip_station", Name: "Dip Station", Description: "Parallel bars for dips", Category: "other"},
		{Slug: "gymnastic_rings", Name: "Gymnastic Rings", Description: "Suspension rings for advanced bodyweight training", Category: "other"},
		{Slug: "parallettes", Name: "Parallettes", Description: "Low parallel bars for handstand training", Category: "other"},
		{Slug: "wall_bars", Name: "Wall Bars", Description: "Swedish wall bars for stretching and bodyweight", Category: "other"},

		// Resistance Equipment
		{Slug: "resistance_bands", Name: "Resistance Bands", Description: "Elastic bands of varying resistance", Category: "other"},
		{Slug: "loop_bands", Name: "Loop Bands", Description: "Circular resistance bands", Category: "other"},
		{Slug: "suspension_trainer", Name: "Suspension Trainer", Description: "TRX or similar suspension system", Category: "other"},
		{Slug: "battle_ropes", Name: "Battle Ropes", Description: "Heavy ropes for cardio and strength", Category: "other"},
		{Slug: "chains", Name: "Chains", Description: "Heavy chains for variable resistance", Category: "other"},

		// Cardio Equipment
		{Slug: "treadmill", Name: "Treadmill", Description: "Motorized running machine", Category: "cardio"},
		{Slug: "stationary_bike", Name: "Stationary Bike", Description: "Indoor cycling bike", Category: "cardio"},
		{Slug: "rowing_machine", Name: "Rowing Machine", Description: "Indoor rowing ergometer", Category: "cardio"},
		{Slug: "elliptical", Name: "Elliptical Machine", Description: "Low-impact cardio machine", Category: "cardio"},
		{Slug: "jump_rope", Name: "Jump Rope", Description: "Speed rope for cardio", Category: "cardio"},
		{Slug: "air_bike", Name: "Air Bike", Description: "Fan bike for high-intensity cardio", Category: "cardio"},
		{Slug: "stair_climber", Name: "Stair Climber", Description: "Stair climbing machine", Category: "cardio"},
		{Slug: "ski_erg", Name: "Ski Erg", Description: "Upper body cardio machine", Category: "cardio"},
		{Slug: "versa_climber", Name: "VersaClimber", Description: "Vertical climbing machine", Category: "cardio"},

		// Machines
		{Slug: "cable_machine", Name: "Cable Machine", Description: "Adjustable cable pulley system", Category: "cable"},
		{Slug: "smith_machine", Name: "Smith Machine", Description: "Guided barbell rack", Category: "machine"},
		{Slug: "leg_press", Name: "Leg Press Machine", Description: "Seated or lying leg press", Category: "machine"},
		{Slug: "lat_pulldown", Name: "Lat Pulldown Machine", Description: "Cable machine for back exercises", Category: "machine"},
		{Slug: "chest_press", Name: "Chest Press Machine", Description: "Seated or lying chest press", Category: "machine"},
		{Slug: "leg_curl", Name: "Leg Curl Machine", Description: "Hamstring curl machine", Category: "machine"},
		{Slug: "leg_extension", Name: "Leg Extension Machine", Description: "Quadriceps extension machine", Category: "machine"},
		{Slug: "pec_deck", Name: "Pec Deck", Description: "Chest fly machine", Category: "machine"},
		{Slug: "shoulder_press", Name: "Shoulder Press Machine", Description: "Seated shoulder press", Category: "machine"},
		{Slug: "hack_squat", Name: "Hack Squat Machine", Description: "Angled squat machine", Category: "machine"},
		{Slug: "calf_raise", Name: "Calf Raise Machine", Description: "Standing or seated calf raise", Category: "machine"},
		{Slug: "hip_abduction", Name: "Hip Abduction Machine", Description: "Outer thigh machine", Category: "machine"},
		{Slug: "hip_adduction", Name: "Hip Adduction Machine", Description: "Inner thigh machine", Category: "machine"},
		{Slug: "back_extension", Name: "Back Extension Machine", Description: "Hyperextension machine", Category: "machine"},
		{Slug: "glute_ham_raise", Name: "Glute Ham Raise", Description: "GHR machine for posterior chain", Category: "machine"},
		{Slug: "reverse_hyper", Name: "Reverse Hyper", Description: "Reverse hyperextension machine", Category: "machine"},
		{Slug: "belt_squat", Name: "Belt Squat Machine", Description: "Belt-loaded squat machine", Category: "machine"},
		{Slug: "pendulum_squat", Name: "Pendulum Squat", Description: "Arc-motion squat machine", Category: "machine"},
		{Slug: "leg_curl_lying", Name: "Lying Leg Curl", Description: "Prone hamstring curl machine", Category: "machine"},
		{Slug: "leg_curl_seated", Name: "Seated Leg Curl", Description: "Seated hamstring curl machine", Category: "machine"},
		{Slug: "preacher_curl", Name: "Preacher Curl Machine", Description: "Seated bicep curl machine", Category: "machine"},
		{Slug: "tricep_dip", Name: "Tricep Dip Machine", Description: "Assisted dip machine", Category: "machine"},
		{Slug: "pull_up_assist", Name: "Assisted Pull-up Machine", Description: "Counterweight pull-up machine", Category: "machine"},

		// Benches and Racks
		{Slug: "flat_bench", Name: "Flat Bench", Description: "Standard flat weight bench", Category: "other"},
		{Slug: "adjustable_bench", Name: "Adjustable Bench", Description: "Incline/decline adjustable bench", Category: "other"},
		{Slug: "squat_rack", Name: "Squat Rack", Description: "Power rack or squat stand", Category: "other"},
		{Slug: "power_rack", Name: "Power Rack", Description: "Full power cage with safety bars", Category: "other"},
		{Slug: "half_rack", Name: "Half Rack", Description: "Wall-mounted squat rack", Category: "other"},
		{Slug: "preacher_bench", Name: "Preacher Bench", Description: "Angled bench for bicep curls", Category: "other"},
		{Slug: "scott_bench", Name: "Scott Bench", Description: "Angled arm curl bench", Category: "other"},
		{Slug: "decline_bench", Name: "Decline Bench", Description: "Fixed decline bench", Category: "other"},
		{Slug: "incline_bench", Name: "Incline Bench", Description: "Fixed incline bench", Category: "other"},
		{Slug: "utility_bench", Name: "Utility Bench", Description: "Multi-purpose bench", Category: "other"},
		{Slug: "hyperextension_bench", Name: "Hyperextension Bench", Description: "Back extension bench", Category: "other"},
		{Slug: "ab_bench", Name: "Ab Bench", Description: "Decline bench for ab exercises", Category: "other"},
		{Slug: "dumbbell_rack", Name: "Dumbbell Rack", Description: "Storage rack for dumbbells", Category: "other"},
		{Slug: "barbell_rack", Name: "Barbell Rack", Description: "Storage rack for barbells", Category: "other"},
		{Slug: "plate_rack", Name: "Plate Rack", Description: "Storage rack for weight plates", Category: "other"},

		// Functional Training
		{Slug: "medicine_ball", Name: "Medicine Ball", Description: "Weighted ball for functional training", Category: "other"},
		{Slug: "slam_ball", Name: "Slam Ball", Description: "Non-bouncing weighted ball", Category: "other"},
		{Slug: "wall_ball", Name: "Wall Ball", Description: "Soft weighted ball for wall throws", Category: "other"},
		{Slug: "stability_ball", Name: "Stability Ball", Description: "Large inflatable exercise ball", Category: "other"},
		{Slug: "bosu_ball", Name: "BOSU Ball", Description: "Balance trainer dome", Category: "other"},
		{Slug: "balance_board", Name: "Balance Board", Description: "Wobble board for balance training", Category: "other"},
		{Slug: "balance_pad", Name: "Balance Pad", Description: "Foam pad for balance training", Category: "other"},
		{Slug: "agility_ladder", Name: "Agility Ladder", Description: "Ladder for footwork drills", Category: "other"},
		{Slug: "speed_hurdles", Name: "Speed Hurdles", Description: "Low hurdles for agility training", Category: "other"},
		{Slug: "cones", Name: "Cones", Description: "Marker cones for agility drills", Category: "other"},
		{Slug: "plyo_box", Name: "Plyo Box", Description: "Box for plyometric exercises", Category: "other"},
		{Slug: "step_platform", Name: "Step Platform", Description: "Aerobic step platform", Category: "other"},
		{Slug: "tire", Name: "Tire", Description: "Large tire for flipping and hitting", Category: "other"},
		{Slug: "sledgehammer", Name: "Sledgehammer", Description: "Heavy hammer for tire workouts", Category: "other"},
		{Slug: "sled", Name: "Sled", Description: "Weighted sled for pushing/pulling", Category: "other"},
		{Slug: "prowler", Name: "Prowler Sled", Description: "Low push/pull sled", Category: "other"},
		{Slug: "farmers_walk", Name: "Farmer's Walk Handles", Description: "Handles for farmer's walks", Category: "other"},
		{Slug: "yoke", Name: "Yoke", Description: "Strongman yoke for carries", Category: "other"},
		{Slug: "atlas_stones", Name: "Atlas Stones", Description: "Heavy round stones for lifting", Category: "other"},
		{Slug: "sandbag", Name: "Sandbag", Description: "Weighted sandbag for functional training", Category: "other"},
		{Slug: "bulgarian_bag", Name: "Bulgarian Bag", Description: "Crescent-shaped weighted bag", Category: "other"},
		{Slug: "mace", Name: "Mace", Description: "Steel mace for rotational training", Category: "other"},
		{Slug: "club_bells", Name: "Club Bells", Description: "Weighted clubs for swinging", Category: "other"},

		// Recovery and Mobility
		{Slug: "foam_roller", Name: "Foam Roller", Description: "Myofascial release tool", Category: "other"},
		{Slug: "massage_ball", Name: "Massage Ball", Description: "Lacrosse ball for trigger points", Category: "other"},
		{Slug: "massage_stick", Name: "Massage Stick", Description: "Roller stick for muscle release", Category: "other"},
		{Slug: "stretching_strap", Name: "Stretching Strap", Description: "Yoga strap for assisted stretching", Category: "other"},
		{Slug: "resistance_bands_therapy", Name: "Therapy Bands", Description: "Light resistance bands for rehab", Category: "other"},
		{Slug: "compression_boots", Name: "Compression Boots", Description: "Pneumatic compression for recovery", Category: "other"},
		{Slug: "percussion_massager", Name: "Percussion Massager", Description: "Electric massage gun", Category: "other"},
		{Slug: "inversion_table", Name: "Inversion Table", Description: "Table for spinal decompression", Category: "other"},

		// Mats and Surfaces
		{Slug: "exercise_mat", Name: "Exercise Mat", Description: "Yoga or workout mat", Category: "other"},
		{Slug: "yoga_mat", Name: "Yoga Mat", Description: "Non-slip yoga mat", Category: "other"},
		{Slug: "pilates_mat", Name: "Pilates Mat", Description: "Thick mat for Pilates", Category: "other"},
		{Slug: "gym_mat", Name: "Gym Mat", Description: "Large exercise mat", Category: "other"},
		{Slug: "crash_mat", Name: "Crash Mat", Description: "Thick safety mat", Category: "other"},
		{Slug: "puzzle_mat", Name: "Puzzle Mat", Description: "Interlocking floor tiles", Category: "other"},

		// Accessories
		{Slug: "ab_wheel", Name: "Ab Wheel", Description: "Core training wheel", Category: "other"},
		{Slug: "ab_straps", Name: "Ab Straps", Description: "Hanging ab straps", Category: "other"},
		{Slug: "dip_belt", Name: "Dip Belt", Description: "Belt for adding weight to dips", Category: "other"},
		{Slug: "weight_vest", Name: "Weight Vest", Description: "Vest for adding bodyweight resistance", Category: "other"},
		{Slug: "ankle_weights", Name: "Ankle Weights", Description: "Weighted cuffs for legs", Category: "other"},
		{Slug: "wrist_weights", Name: "Wrist Weights", Description: "Weighted cuffs for arms", Category: "other"},
		{Slug: "lifting_straps", Name: "Lifting Straps", Description: "Wrist straps for grip assistance", Category: "other"},
		{Slug: "lifting_belt", Name: "Lifting Belt", Description: "Leather belt for core support", Category: "other"},
		{Slug: "knee_sleeves", Name: "Knee Sleeves", Description: "Compression sleeves for knees", Category: "other"},
		{Slug: "wrist_wraps", Name: "Wrist Wraps", Description: "Supportive wraps for wrists", Category: "other"},
		{Slug: "elbow_sleeves", Name: "Elbow Sleeves", Description: "Compression sleeves for elbows", Category: "other"},
		{Slug: "gloves", Name: "Gloves", Description: "Workout gloves for grip", Category: "other"},
		{Slug: "chalk", Name: "Chalk", Description: "Magnesium carbonate for grip", Category: "other"},
		{Slug: "liquid_chalk", Name: "Liquid Chalk", Description: "Liquid grip enhancer", Category: "other"},
		{Slug: "rosin_bag", Name: "Rosin Bag", Description: "Grip enhancer bag", Category: "other"},

		// Specialty Equipment
		{Slug: "landmine", Name: "Landmine", Description: "Pivoting barbell attachment", Category: "other"},
		{Slug: "t_bar_row", Name: "T-Bar Row", Description: "T-shaped rowing handle", Category: "other"},
		{Slug: "cable_handles", Name: "Cable Handles", Description: "Various cable attachment handles", Category: "cable"},
		{Slug: "lat_pulldown_bar", Name: "Lat Pulldown Bar", Description: "Wide grip pulldown bar", Category: "cable"},
		{Slug: "straight_bar", Name: "Straight Bar Cable", Description: "Straight bar for cable exercises", Category: "cable"},
		{Slug: "rope_attachment", Name: "Rope Attachment", Description: "Rope for cable exercises", Category: "cable"},
		{Slug: "tricep_rope", Name: "Tricep Rope", Description: "Rope specifically for tricep work", Category: "cable"},
		{Slug: "v_bar", Name: "V-Bar", Description: "V-shaped cable attachment", Category: "cable"},
		{Slug: "mag_grip", Name: "MAG Grip", Description: "Neutral grip cable attachment", Category: "cable"},
		{Slug: "d_handle", Name: "D-Handle", Description: "Single handle for cable work", Category: "cable"},
		{Slug: "ankle_strap", Name: "Ankle Strap", Description: "Strap for cable leg exercises", Category: "cable"},

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
