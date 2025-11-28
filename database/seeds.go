package database

import (
	"lamari-fit-api/models"
	"lamari-fit-api/utils"
	"log"
	"strings"

	"github.com/google/uuid"
)

// generateSlug creates a URL-friendly slug from a name
func generateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)
	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "_")
	// Remove apostrophes
	slug = strings.ReplaceAll(slug, "'", "")
	// Replace multiple hyphens with single hyphen
	slug = strings.ReplaceAll(slug, "-", "_")
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

func SeedExerciseTypes() {
	exerciseTypes := []models.ExerciseType{
		{Slug: "compound", Name: "Compound", Description: "Multi-joint movements"},
		{Slug: "isolation", Name: "Isolation", Description: "Targeted single-joint exercises"},
		{Slug: "isometric", Name: "Isometric", Description: "Static holds"},
		{Slug: "plyometric", Name: "Plyometric", Description: "Explosive / jump-based"},
		{Slug: "cardio", Name: "Cardio", Description: "Heart rate / conditioning"},
		{Slug: "mobility", Name: "Mobility", Description: "Movement quality / joint prep"},
		{Slug: "stretching", Name: "Stretching", Description: "Flexibility / cool-down"},
		{Slug: "push", Name: "Push", Description: "Exercises where force is pushed away from the body"},
		{Slug: "pull", Name: "Pull", Description: "Exercises where force is pulled toward the body"},
	}

	for _, exerciseType := range exerciseTypes {
		var existing models.ExerciseType
		if err := DB.Where("slug = ?", exerciseType.Slug).First(&existing).Error; err != nil {
			if err := DB.Create(&exerciseType).Error; err != nil {
				log.Printf("Failed to create exercise type %s: %v", exerciseType.Name, err)
			} else {
				log.Printf("Created exercise type: %s", exerciseType.Name)
			}
		} else {
			log.Printf("Exercise type already exists: %s", exerciseType.Name)
		}
	}
}

// GetExerciseTypeMappings returns a map of exercise slug to exercise type slugs
func GetExerciseTypeMappings() map[string][]string {
	return map[string][]string{
		// Chest Exercises
		"push_ups":             {"compound", "push"},
		"bench_press":          {"compound", "push"},
		"incline_bench_press":  {"compound", "push"},
		"decline_bench_press":  {"compound", "push"},
		"dumbbell_bench_press": {"compound", "push"},
		"dumbbell_flyes":       {"isolation", "push"},
		"chest_dips":           {"compound", "push"},
		"cable_crossovers":     {"isolation", "push"},

		// Back Exercises
		"pull_ups":      {"compound", "pull"},
		"chin_ups":      {"compound", "pull"},
		"lat_pulldowns": {"compound", "pull"},
		"barbell_rows":  {"compound", "pull"},
		"dumbbell_rows": {"compound", "pull"},
		"cable_rows":    {"compound", "pull"},
		"face_pulls":    {"isolation", "pull"},
		"deadlifts":     {"compound", "pull"},
		"rack_pulls":    {"compound", "pull"},
		"inverted_rows": {"compound", "pull"},

		// Shoulder Exercises
		"overhead_press":          {"compound", "push"},
		"dumbbell_shoulder_press": {"compound", "push"},
		"lateral_raises":          {"isolation", "push"},
		"rear_delt_flyes":         {"isolation", "pull"},
		"front_raises":            {"isolation", "push"},
		"arnold_press":            {"compound", "push"},
		"upright_rows":            {"compound", "pull"},
		"pike_push_ups":           {"compound", "push"},
		"handstand_push_ups":      {"compound", "push"},

		// Arm Exercises
		"bicep_curls":               {"isolation", "pull"},
		"hammer_curls":              {"isolation", "pull"},
		"preacher_curls":            {"isolation", "pull"},
		"concentration_curls":       {"isolation", "pull"},
		"tricep_dips":               {"compound", "push"},
		"tricep_pushdowns":          {"isolation", "push"},
		"skull_crushers":            {"isolation", "push"},
		"overhead_tricep_extension": {"isolation", "push"},
		"close_grip_bench_press":    {"compound", "push"},

		// Leg Exercises
		"squats":                 {"compound", "push"},
		"front_squats":           {"compound", "push"},
		"goblet_squats":          {"compound", "push"},
		"leg_press":              {"compound", "push"},
		"lunges":                 {"compound", "push"},
		"walking_lunges":         {"compound", "push"},
		"bulgarian_split_squats": {"compound", "push"},
		"step_ups":               {"compound", "push"},
		"leg_extensions":         {"isolation", "push"},
		"leg_curls":              {"isolation", "pull"},
		"romanian_deadlifts":     {"compound", "pull"},
		"hip_thrusts":            {"compound", "push"},
		"glute_bridges":          {"compound", "push"},
		"calf_raises":            {"isolation", "push"},
		"seated_calf_raises":     {"isolation", "push"},
		"box_jumps":              {"plyometric", "push"},
		"jump_squats":            {"plyometric", "push"},

		// Core Exercises
		"crunches":           {"isolation"},
		"sit_ups":            {"compound"},
		"planks":             {"isometric"},
		"side_planks":        {"isometric"},
		"russian_twists":     {"isolation"},
		"bicycle_crunches":   {"compound"},
		"leg_raises":         {"isolation"},
		"hanging_leg_raises": {"isolation"},
		"ab_wheel_rollouts":  {"compound"},
		"mountain_climbers":  {"cardio", "compound"},
		"dead_bugs":          {"isolation"},
		"hollow_body_holds":  {"isometric"},
		"superman":           {"isometric"},
		"back_extensions":    {"compound"},
		"hyperextensions":    {"compound"},
		"good_mornings":      {"compound"},
		"cable_woodchops":    {"compound"},
		"pallof_press":       {"isometric"},

		// Cardio Exercises
		"burpees":            {"cardio", "plyometric", "compound"},
		"high_knees":         {"cardio", "plyometric"},
		"jumping_jacks":      {"cardio"},
		"bear_crawls":        {"cardio", "compound"},
		"crab_walks":         {"cardio", "compound"},
		"treadmill_running":  {"cardio"},
		"stationary_bike":    {"cardio"},
		"rowing_machine":     {"cardio", "compound", "pull"},
		"elliptical_machine": {"cardio"},
		"jump_rope":          {"cardio", "plyometric"},
		"stair_climbing":     {"cardio"},

		// Olympic Lifts
		"clean_and_jerk": {"compound", "pull", "push"},
		"snatch":         {"compound", "pull"},
		"power_cleans":   {"compound", "pull"},
		"hang_cleans":    {"compound", "pull"},
		"clean_pulls":    {"compound", "pull"},
		"push_press":     {"compound", "push"},

		// Functional Exercises
		"farmers_walk":        {"compound"},
		"kettlebell_swings":   {"compound", "pull"},
		"turkish_get_ups":     {"compound"},
		"battle_ropes":        {"cardio", "compound"},
		"medicine_ball_slams": {"plyometric", "compound"},
		"box_step_overs":      {"cardio", "compound"},
		"wall_balls":          {"plyometric", "compound", "push"},

		// Unilateral Exercises
		"single_leg_glute_bridge": {"compound", "push"},
		"single_leg_deadlift":     {"compound", "pull"},
		"pistol_squats":           {"compound", "push"},
		"single_leg_calf_raises":  {"isolation", "push"},
		"bird_dog":                {"isometric"},

		// Stretching and Mobility
		"cat_cow_stretch":            {"mobility", "stretching"},
		"childs_pose":                {"mobility", "stretching"},
		"downward_dog":               {"mobility", "stretching"},
		"pigeon_pose":                {"mobility", "stretching"},
		"seated_forward_fold":        {"stretching"},
		"standing_hamstring_stretch": {"stretching"},
		"hip_flexor_stretch":         {"stretching", "mobility"},
		"chest_stretch":              {"stretching"},
		"tricep_stretch":             {"stretching"},
		"shoulder_stretch":           {"stretching"},
		"standing_quad_stretch":      {"stretching"},
		"standing_calf_stretch":      {"stretching"},
		"shoulder_rolls":             {"mobility"},
		"arm_circles":                {"mobility"},
		"neck_rolls":                 {"mobility"},
		"hip_circles":                {"mobility"},
		"leg_swings":                 {"mobility"},
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

	// Get equipment IDs
	equipmentMap := make(map[string]uuid.UUID)
	var equipment []models.Equipment
	DB.Find(&equipment)
	for _, eq := range equipment {
		equipmentMap[eq.Slug] = eq.ID
	}

	// Get exercise type IDs
	exerciseTypeMap := make(map[string]uuid.UUID)
	var exerciseTypes []models.ExerciseType
	DB.Find(&exerciseTypes)
	for _, et := range exerciseTypes {
		exerciseTypeMap[et.Slug] = et.ID
	}

	// Get equipment mappings
	exerciseEquipmentMappings := GetExerciseEquipmentMappings()

	// Get exercise type mappings
	exerciseTypeMappings := GetExerciseTypeMappings()

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
		if err := DB.Where("slug = ?", exerciseData.Exercise.Name).First(&existingExercise).Error; err != nil {
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

			// Assign equipment
			if equipmentList, exists := exerciseEquipmentMappings[exerciseData.Exercise.Slug]; exists {
				for _, eqData := range equipmentList {
					if equipmentID, exists := equipmentMap[eqData.Slug]; exists {
						assignment := models.ExerciseEquipment{
							ExerciseID:  exerciseData.Exercise.ID,
							EquipmentID: equipmentID,
							Optional:    eqData.Optional,
							Notes:       eqData.Notes,
						}

						if err := DB.Create(&assignment).Error; err != nil {
							log.Printf("Failed to assign equipment %s to exercise %s: %v", eqData.Slug, exerciseData.Exercise.Name, err)
						} else {
							log.Printf("Assigned equipment %s to exercise %s", eqData.Slug, exerciseData.Exercise.Name)
						}
					}
				}
			}

			// Assign exercise types
			if typeList, exists := exerciseTypeMappings[exerciseData.Exercise.Slug]; exists {
				for _, typeSlug := range typeList {
					if exerciseTypeID, exists := exerciseTypeMap[typeSlug]; exists {
						assignment := models.ExerciseExerciseType{
							ExerciseID:     exerciseData.Exercise.ID,
							ExerciseTypeID: exerciseTypeID,
						}

						if err := DB.Create(&assignment).Error; err != nil {
							log.Printf("Failed to assign exercise type %s to exercise %s: %v", typeSlug, exerciseData.Exercise.Name, err)
						} else {
							log.Printf("Assigned exercise type %s to exercise %s", typeSlug, exerciseData.Exercise.Name)
						}
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
		{Name: "Weight Loss", NameSlug: "weight_loss", Description: "Reduce body weight and body fat percentage", Category: "body_composition", IconName: "scale"},
		{Name: "Muscle Gain", NameSlug: "muscle_gain", Description: "Build lean muscle mass and increase strength", Category: "body_composition", IconName: "dumbbell"},
		{Name: "Endurance", NameSlug: "endurance", Description: "Improve cardiovascular fitness and stamina", Category: "performance", IconName: "running"},
		{Name: "Strength", NameSlug: "strength", Description: "Increase maximum strength and power output", Category: "performance", IconName: "weight"},
		{Name: "Flexibility", NameSlug: "flexibility", Description: "Improve range of motion and mobility", Category: "wellness", IconName: "stretch"},
		{Name: "General Fitness", NameSlug: "general_fitness", Description: "Overall health and wellness improvement", Category: "wellness", IconName: "heart"},
		{Name: "Athletic Performance", NameSlug: "athletic_performance", Description: "Sport-specific performance enhancement", Category: "performance", IconName: "trophy"},
		{Name: "Rehabilitation", NameSlug: "rehabilitation", Description: "Recover from injury or medical condition", Category: "wellness", IconName: "medical"},
		{Name: "Body Recomposition", NameSlug: "body_recomposition", Description: "Simultaneous fat loss and muscle gain", Category: "body_composition", IconName: "transform"},
		{Name: "Stress Relief", NameSlug: "stress_relief", Description: "Mental health and stress management through exercise", Category: "wellness", IconName: "mindfulness"},
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

func SeedSpecialties() {
	specialties := []models.Specialty{
		{
			Name:        "Strength Training",
			Description: "Build muscle mass, increase strength, and improve overall power through resistance training and weightlifting.",
		},
		{
			Name:        "Weight Loss",
			Description: "Achieve sustainable fat loss through customized workout programs and nutritional guidance tailored to your goals.",
		},
		{
			Name:        "Bodybuilding",
			Description: "Develop muscle definition, symmetry, and size with specialized training techniques for competitive or aesthetic goals.",
		},
		{
			Name:        "Functional Fitness",
			Description: "Improve everyday movement patterns, balance, and coordination through exercises that mimic real-life activities.",
		},
		{
			Name:        "HIIT",
			Description: "High-Intensity Interval Training for maximum calorie burn, cardiovascular fitness, and metabolic conditioning.",
		},
		{
			Name:        "Yoga",
			Description: "Enhance flexibility, mindfulness, and inner peace through traditional and modern yoga practices.",
		},
		{
			Name:        "Cardio",
			Description: "Improve cardiovascular health, endurance, and stamina through running, cycling, and other aerobic exercises.",
		},
		{
			Name:        "Rehabilitation",
			Description: "Recover from injuries and improve mobility through specialized corrective exercise and therapeutic movement.",
		},
		{
			Name:        "Mobility",
			Description: "Increase range of motion, reduce stiffness, and prevent injuries through targeted mobility work and stretching.",
		},
		{
			Name:        "CrossFit",
			Description: "Build overall fitness through constantly varied, high-intensity functional movements and competitive workouts.",
		},
	}

	for _, specialty := range specialties {
		var existing models.Specialty
		if err := DB.Where("name = ?", specialty.Name).First(&existing).Error; err != nil {
			if err := DB.Create(&specialty).Error; err != nil {
				log.Printf("Failed to create specialty: %v", err)
			} else {
				log.Printf("Created specialty: %s", specialty.Name)
			}
		} else {
			log.Printf("Specialty already exists: %s", specialty.Name)
		}
	}
}

func SeedTrainerProfiles() {
	// Get existing users (at least 2-3 should exist from previous seeds)
	var users []models.User
	DB.Limit(3).Find(&users)

	if len(users) < 2 {
		log.Println("Not enough users to seed trainer profiles")
		return
	}

	// Get specialties for trainer profiles
	var strengthTraining, weightLoss, bodybuilding, functionalFitness, rehabilitation, mobility models.Specialty
	DB.Where("name = ?", "Strength Training").First(&strengthTraining)
	DB.Where("name = ?", "Weight Loss").First(&weightLoss)
	DB.Where("name = ?", "Bodybuilding").First(&bodybuilding)
	DB.Where("name = ?", "Functional Fitness").First(&functionalFitness)
	DB.Where("name = ?", "Rehabilitation").First(&rehabilitation)
	DB.Where("name = ?", "Mobility").First(&mobility)

	// First trainer profile
	var existing1 models.TrainerProfile
	hourlyRate1 := 75.00
	if err := DB.Where("user_id = ?", users[0].ID).First(&existing1).Error; err != nil {
		nyCity := "New York"
		nyRegion := "NY"
		nyCountry := "US"
		nyLat := 40.7128
		nyLng := -74.0060
		profile1 := models.TrainerProfile{
			UserID:     users[0].ID,
			Bio:        "Certified personal trainer with 5+ years experience in strength training, weight loss, and bodybuilding. Passionate about helping clients achieve their fitness goals through customized workout programs and nutrition guidance.",
			HourlyRate: &hourlyRate1,
			Location: models.Location{
				City:        &nyCity,
				Region:      &nyRegion,
				CountryCode: &nyCountry,
				Latitude:    &nyLat,
				Longitude:   &nyLng,
			},
			Visibility: "public",
		}
		if err := DB.Create(&profile1).Error; err != nil {
			log.Printf("Failed to create trainer profile: %v", err)
		} else {
			// Associate specialties
			specialties := []models.Specialty{strengthTraining, weightLoss, bodybuilding}
			DB.Model(&profile1).Association("Specialties").Append(&specialties)
			log.Printf("Created trainer profile for user %s", profile1.UserID)
		}
	} else {
		log.Printf("Trainer profile already exists for user %s", users[0].ID)
	}

	// Second trainer profile if we have enough users
	if len(users) >= 2 {
		var existing2 models.TrainerProfile
		hourlyRate2 := 60.00
		if err := DB.Where("user_id = ?", users[1].ID).First(&existing2).Error; err != nil {
			laCity := "Los Angeles"
			laRegion := "CA"
			laCountry := "US"
			laLat := 34.0522
			laLng := -118.2437
			profile2 := models.TrainerProfile{
				UserID:     users[1].ID,
				Bio:        "Specializing in functional fitness and injury prevention. I help clients improve mobility, recover from injuries, and build sustainable fitness habits. Certified in rehabilitation and corrective exercise.",
				HourlyRate: &hourlyRate2,
				Location: models.Location{
					City:        &laCity,
					Region:      &laRegion,
					CountryCode: &laCountry,
					Latitude:    &laLat,
					Longitude:   &laLng,
				},
				Visibility: "public",
			}
			if err := DB.Create(&profile2).Error; err != nil {
				log.Printf("Failed to create trainer profile: %v", err)
			} else {
				// Associate specialties
				specialties := []models.Specialty{functionalFitness, rehabilitation, mobility}
				DB.Model(&profile2).Association("Specialties").Append(&specialties)
				log.Printf("Created trainer profile for user %s", profile2.UserID)
			}
		} else {
			log.Printf("Trainer profile already exists for user %s", users[1].ID)
		}
	}
}

func SeedTrainerClientLinks() {
	// Get users and trainer profiles
	var users []models.User
	DB.Limit(4).Find(&users)

	if len(users) < 3 {
		log.Println("Not enough users to seed trainer-client links")
		return
	}

	// Get trainer profiles to find trainer user IDs
	var trainerProfiles []models.TrainerProfile
	DB.Limit(2).Find(&trainerProfiles)

	if len(trainerProfiles) < 1 {
		log.Println("No trainer profiles found to seed trainer-client links")
		return
	}

	// Create sample trainer-client relationships
	// Using the first trainer with users who aren't trainers
	trainerUserID := trainerProfiles[0].UserID

	// Find users who aren't this trainer
	var clientUsers []models.User
	for _, user := range users {
		if user.ID != trainerUserID {
			clientUsers = append(clientUsers, user)
		}
	}

	if len(clientUsers) < 2 {
		log.Println("Not enough client users to seed trainer-client links")
		return
	}

	// Create one active and one pending relationship
	links := []models.TrainerClientLink{
		{
			TrainerID: trainerUserID,
			ClientID:  clientUsers[0].ID,
			Status:    "active",
		},
	}

	// Add pending if we have another user
	if len(clientUsers) >= 2 {
		links = append(links, models.TrainerClientLink{
			TrainerID: trainerUserID,
			ClientID:  clientUsers[1].ID,
			Status:    "pending",
		})
	}

	for _, link := range links {
		var existing models.TrainerClientLink
		if err := DB.Where("trainer_id = ? AND client_id = ?", link.TrainerID, link.ClientID).First(&existing).Error; err != nil {
			if err := DB.Create(&link).Error; err != nil {
				log.Printf("Failed to create trainer-client link: %v", err)
			} else {
				log.Printf("Created trainer-client link: trainer=%s, client=%s, status=%s", link.TrainerID, link.ClientID, link.Status)
			}
		} else {
			log.Printf("Trainer-client link already exists: trainer=%s, client=%s", link.TrainerID, link.ClientID)
		}
	}
}

// SeedGlobalRPEScale creates the global standard RPE scale
func SeedGlobalRPEScale() {
	// Check if global scale already exists
	var existing models.RPEScale
	if err := DB.Where("is_global = ?", true).First(&existing).Error; err == nil {
		log.Println("Global RPE scale already exists")
		return
	}

	// Create the global RPE scale
	scale := models.RPEScale{
		Name:        "Standard RPE Scale",
		Description: "Standard Rate of Perceived Exertion scale used globally for resistance training",
		MinValue:    1,
		MaxValue:    10,
		IsGlobal:    true,
		TrainerID:   nil,
	}

	if err := DB.Create(&scale).Error; err != nil {
		log.Printf("Failed to create global RPE scale: %v", err)
		return
	}
	log.Printf("Created global RPE scale: %s", scale.Name)

	// Create scale values
	values := []models.RPEScaleValue{
		{ScaleID: scale.ID, Value: 1, Label: "Very Light", Description: "Minimal effort, easy breathing, could do this all day"},
		{ScaleID: scale.ID, Value: 2, Label: "Light", Description: "Comfortable, could maintain for hours"},
		{ScaleID: scale.ID, Value: 3, Label: "Moderate", Description: "Breathing harder but comfortable"},
		{ScaleID: scale.ID, Value: 4, Label: "Somewhat Hard", Description: "Sweating lightly, can hold a conversation"},
		{ScaleID: scale.ID, Value: 5, Label: "Hard", Description: "Deep breathing, short sentences only"},
		{ScaleID: scale.ID, Value: 6, Label: "Harder", Description: "Can speak a few words at a time, 4+ reps in reserve"},
		{ScaleID: scale.ID, Value: 7, Label: "Very Hard", Description: "3-4 reps left in tank, challenging"},
		{ScaleID: scale.ID, Value: 8, Label: "Very Hard+", Description: "2 reps left in tank, definitely working hard"},
		{ScaleID: scale.ID, Value: 9, Label: "Extremely Hard", Description: "1 rep left in tank, near maximum effort"},
		{ScaleID: scale.ID, Value: 10, Label: "Maximum", Description: "Could not do another rep, absolute maximum effort"},
	}

	for _, value := range values {
		if err := DB.Create(&value).Error; err != nil {
			log.Printf("Failed to create RPE value %d: %v", value.Value, err)
		} else {
			log.Printf("Created RPE value: %d - %s", value.Value, value.Label)
		}
	}
}

// SeedSearchableUsers creates 50 public users with locations for search testing
func SeedSearchableUsers() {
	log.Println("Seeding 50 searchable users...")

	// Helper to create string pointer
	strPtr := func(s string) *string { return &s }
	floatPtr := func(f float64) *float64 { return &f }

	// Diverse user data with global locations
	userData := []struct {
		FirstName           string
		LastName            string
		Email               string
		Bio                 string
		IsLookingForTrainer bool
		City                string
		Region              string
		CountryCode         string
		Latitude            float64
		Longitude           float64
	}{
		// US - East Coast
		{"Emma", "Johnson", "emma.johnson@example.com", "Fitness enthusiast looking to build muscle and improve overall health.", true, "New York", "NY", "US", 40.7128, -74.0060},
		{"Michael", "Williams", "michael.williams@example.com", "Marathon runner seeking strength training guidance.", true, "Boston", "MA", "US", 42.3601, -71.0589},
		{"Sarah", "Brown", "sarah.brown@example.com", "Yoga practitioner wanting to add weight training.", true, "Philadelphia", "PA", "US", 39.9526, -75.1652},
		{"James", "Davis", "james.davis@example.com", "Former athlete getting back into shape.", true, "Washington", "DC", "US", 38.9072, -77.0369},
		{"Jennifer", "Miller", "jennifer.miller@example.com", "New mom looking for postpartum fitness guidance.", true, "Miami", "FL", "US", 25.7617, -80.1918},
		{"Robert", "Wilson", "robert.wilson@example.com", "Senior looking for mobility and strength work.", false, "Atlanta", "GA", "US", 33.7490, -84.3880},
		{"Lisa", "Moore", "lisa.moore@example.com", "Busy professional needing efficient workouts.", true, "Charlotte", "NC", "US", 35.2271, -80.8431},
		{"David", "Taylor", "david.taylor@example.com", "Bodybuilding competitor in off-season.", false, "Orlando", "FL", "US", 28.5383, -81.3792},

		// US - West Coast
		{"Emily", "Anderson", "emily.anderson@example.com", "CrossFit athlete looking to improve Olympic lifts.", true, "Los Angeles", "CA", "US", 34.0522, -118.2437},
		{"Christopher", "Thomas", "christopher.thomas@example.com", "Tech worker combating sedentary lifestyle.", true, "San Francisco", "CA", "US", 37.7749, -122.4194},
		{"Ashley", "Jackson", "ashley.jackson@example.com", "Triathlete seeking performance optimization.", true, "Seattle", "WA", "US", 47.6062, -122.3321},
		{"Matthew", "White", "matthew.white@example.com", "Surfer wanting functional fitness training.", false, "San Diego", "CA", "US", 32.7157, -117.1611},
		{"Amanda", "Harris", "amanda.harris@example.com", "Hiking enthusiast building endurance.", true, "Portland", "OR", "US", 45.5152, -122.6784},
		{"Joshua", "Martin", "joshua.martin@example.com", "Basketball player improving vertical leap.", true, "Phoenix", "AZ", "US", 33.4484, -112.0740},
		{"Stephanie", "Garcia", "stephanie.garcia@example.com", "Dancer adding strength training.", true, "Las Vegas", "NV", "US", 36.1699, -115.1398},
		{"Andrew", "Martinez", "andrew.martinez@example.com", "MMA fighter in training camp.", false, "Denver", "CO", "US", 39.7392, -104.9903},

		// US - Central
		{"Nicole", "Robinson", "nicole.robinson@example.com", "Runner transitioning to ultra-marathons.", true, "Chicago", "IL", "US", 41.8781, -87.6298},
		{"Daniel", "Clark", "daniel.clark@example.com", "Powerlifter focusing on competition prep.", false, "Houston", "TX", "US", 29.7604, -95.3698},
		{"Michelle", "Rodriguez", "michelle.rodriguez@example.com", "Weight loss journey, down 50lbs so far.", true, "Dallas", "TX", "US", 32.7767, -96.7970},
		{"Kevin", "Lewis", "kevin.lewis@example.com", "Golf enthusiast improving core strength.", true, "Austin", "TX", "US", 30.2672, -97.7431},
		{"Rachel", "Lee", "rachel.lee@example.com", "Swimmer cross-training for competitions.", true, "Minneapolis", "MN", "US", 44.9778, -93.2650},
		{"Justin", "Walker", "justin.walker@example.com", "Football player in off-season training.", false, "Kansas City", "MO", "US", 39.0997, -94.5786},

		// UK
		{"Charlotte", "Hall", "charlotte.hall@example.com", "London professional seeking work-life balance through fitness.", true, "London", "England", "GB", 51.5074, -0.1278},
		{"Oliver", "Young", "oliver.young@example.com", "Rugby player improving conditioning.", false, "Manchester", "England", "GB", 53.4808, -2.2426},
		{"Sophie", "Allen", "sophie.allen@example.com", "Pilates instructor adding weight training.", true, "Birmingham", "England", "GB", 52.4862, -1.8904},
		{"Harry", "King", "harry.king@example.com", "Cricket player building explosive power.", true, "Leeds", "England", "GB", 53.8008, -1.5491},
		{"Grace", "Wright", "grace.wright@example.com", "Netball player improving agility.", true, "Glasgow", "Scotland", "GB", 55.8642, -4.2518},
		{"Jack", "Scott", "jack.scott@example.com", "Cyclist focusing on leg strength.", false, "Edinburgh", "Scotland", "GB", 55.9533, -3.1883},

		// Canada
		{"Olivia", "Green", "olivia.green@example.com", "Hockey player building off-ice strength.", true, "Toronto", "ON", "CA", 43.6532, -79.3832},
		{"Ethan", "Adams", "ethan.adams@example.com", "Skier preparing for winter season.", true, "Vancouver", "BC", "CA", 49.2827, -123.1207},
		{"Ava", "Nelson", "ava.nelson@example.com", "Figure skater improving flexibility and strength.", true, "Montreal", "QC", "CA", 45.5017, -73.5673},
		{"Noah", "Carter", "noah.carter@example.com", "Lacrosse player building endurance.", false, "Calgary", "AB", "CA", 51.0447, -114.0719},

		// Australia
		{"Mia", "Mitchell", "mia.mitchell@example.com", "Beach volleyball player training year-round.", true, "Sydney", "NSW", "AU", -33.8688, 151.2093},
		{"Liam", "Roberts", "liam.roberts@example.com", "AFL player in pre-season training.", false, "Melbourne", "VIC", "AU", -37.8136, 144.9631},
		{"Isabella", "Turner", "isabella.turner@example.com", "Swimmer training for nationals.", true, "Brisbane", "QLD", "AU", -27.4698, 153.0251},
		{"Mason", "Phillips", "mason.phillips@example.com", "Surfer building paddle strength.", true, "Perth", "WA", "AU", -31.9505, 115.8605},

		// Europe
		{"Sophia", "Campbell", "sophia.campbell@example.com", "Tennis player improving footwork.", true, "Paris", "Île-de-France", "FR", 48.8566, 2.3522},
		{"Lucas", "Evans", "lucas.evans@example.com", "Cyclist training for Tour amateur events.", true, "Lyon", "Auvergne-Rhône-Alpes", "FR", 45.7640, 4.8357},
		{"Amelia", "Edwards", "amelia.edwards@example.com", "Fencer building explosive speed.", true, "Berlin", "Berlin", "DE", 52.5200, 13.4050},
		{"Benjamin", "Collins", "benjamin.collins@example.com", "Handball player improving throwing power.", false, "Munich", "Bavaria", "DE", 48.1351, 11.5820},
		{"Chloe", "Stewart", "chloe.stewart@example.com", "Field hockey player building endurance.", true, "Amsterdam", "North Holland", "NL", 52.3676, 4.9041},
		{"Alexander", "Sanchez", "alexander.sanchez@example.com", "Soccer player in academy training.", true, "Barcelona", "Catalonia", "ES", 41.3851, 2.1734},
		{"Ella", "Morris", "ella.morris@example.com", "Padel player seeking agility training.", true, "Madrid", "Madrid", "ES", 40.4168, -3.7038},
		{"William", "Rogers", "william.rogers@example.com", "Water polo player building swim strength.", false, "Rome", "Lazio", "IT", 41.9028, 12.4964},

		// Asia
		{"Aria", "Reed", "aria.reed@example.com", "Martial artist training for competition.", true, "Tokyo", "Tokyo", "JP", 35.6762, 139.6503},
		{"Henry", "Cook", "henry.cook@example.com", "Baseball player improving batting power.", true, "Osaka", "Osaka", "JP", 34.6937, 135.5023},
		{"Scarlett", "Morgan", "scarlett.morgan@example.com", "Badminton player building court speed.", true, "Singapore", "Singapore", "SG", 1.3521, 103.8198},
		{"Sebastian", "Bell", "sebastian.bell@example.com", "Table tennis player improving reflexes.", false, "Hong Kong", "Hong Kong", "HK", 22.3193, 114.1694},
		{"Victoria", "Murphy", "victoria.murphy@example.com", "Yoga instructor expanding practice.", true, "Seoul", "Seoul", "KR", 37.5665, 126.9780},
		{"Jack", "Bailey", "jack.bailey2@example.com", "Esports player improving posture and health.", true, "Taipei", "Taiwan", "TW", 25.0330, 121.5654},
	}

	hashedPassword, _ := utils.HashPassword("Password123!")

	for _, u := range userData {
		var existing models.User
		if err := DB.Where("email = ?", u.Email).First(&existing).Error; err == nil {
			log.Printf("User already exists: %s", u.Email)
			continue
		}

		user := models.User{
			Email:               u.Email,
			Password:            hashedPassword,
			FirstName:           u.FirstName,
			LastName:            u.LastName,
			Provider:            "local",
			IsActive:            true,
			ProfileVisibility:   "public",
			IsLookingForTrainer: u.IsLookingForTrainer,
			Bio:                 u.Bio,
			Location: models.Location{
				City:        strPtr(u.City),
				Region:      strPtr(u.Region),
				CountryCode: strPtr(u.CountryCode),
				Latitude:    floatPtr(u.Latitude),
				Longitude:   floatPtr(u.Longitude),
			},
		}

		if err := DB.Create(&user).Error; err != nil {
			log.Printf("Failed to create user %s: %v", u.Email, err)
		} else {
			log.Printf("Created searchable user: %s %s (%s)", u.FirstName, u.LastName, u.City)
		}
	}

	log.Println("Finished seeding searchable users")
}

// SeedSearchableTrainers creates 50 public trainers with locations for search testing
func SeedSearchableTrainers() {
	log.Println("Seeding 50 searchable trainers...")

	// Helper to create string pointer
	strPtr := func(s string) *string { return &s }
	floatPtr := func(f float64) *float64 { return &f }

	// Get all specialties for random assignment
	var allSpecialties []models.Specialty
	DB.Find(&allSpecialties)
	if len(allSpecialties) == 0 {
		log.Println("No specialties found - run SeedSpecialties first")
		return
	}

	// Diverse trainer data with global locations
	trainerData := []struct {
		FirstName           string
		LastName            string
		Email               string
		Bio                 string
		HourlyRate          float64
		IsLookingForClients bool
		City                string
		Region              string
		CountryCode         string
		Latitude            float64
		Longitude           float64
		SpecialtyIndices    []int // indices into allSpecialties
	}{
		// US - East Coast
		{"Marcus", "Johnson", "marcus.trainer@example.com", "NASM certified personal trainer with 10+ years experience. Specializing in strength training and body transformation.", 85.00, true, "New York", "NY", "US", 40.7589, -73.9851, []int{0, 1, 2}},
		{"Samantha", "Williams", "samantha.trainer@example.com", "Former Division I athlete turned trainer. Expert in athletic performance and sports conditioning.", 95.00, true, "Brooklyn", "NY", "US", 40.6782, -73.9442, []int{1, 3, 4}},
		{"Derek", "Brown", "derek.trainer@example.com", "Bodybuilding coach with multiple competition wins. Specializing in contest prep and nutrition.", 120.00, true, "Boston", "MA", "US", 42.3601, -71.0589, []int{2, 5}},
		{"Christina", "Davis", "christina.trainer@example.com", "Women's fitness specialist focusing on strength training and confidence building.", 75.00, true, "Philadelphia", "PA", "US", 39.9526, -75.1652, []int{0, 1}},
		{"Marcus", "Miller", "marcus.miller.trainer@example.com", "Certified strength and conditioning specialist. Work with athletes of all levels.", 90.00, true, "Washington", "DC", "US", 38.9072, -77.0369, []int{1, 3}},
		{"Jasmine", "Wilson", "jasmine.trainer@example.com", "Holistic fitness coach combining strength training with mindfulness practices.", 80.00, true, "Miami", "FL", "US", 25.7617, -80.1918, []int{0, 4, 5}},
		{"Brandon", "Moore", "brandon.trainer@example.com", "Former NFL trainer now working with everyday athletes. Focus on functional fitness.", 110.00, true, "Atlanta", "GA", "US", 33.7490, -84.3880, []int{3, 4}},
		{"Alexis", "Taylor", "alexis.trainer@example.com", "Pre and postnatal fitness specialist. Helping moms stay strong through all stages.", 70.00, true, "Charlotte", "NC", "US", 35.2271, -80.8431, []int{0, 5}},

		// US - West Coast
		{"Tyler", "Anderson", "tyler.trainer@example.com", "CrossFit Level 3 trainer with Olympic lifting specialization.", 100.00, true, "Los Angeles", "CA", "US", 34.0522, -118.2437, []int{1, 2, 3}},
		{"Megan", "Thomas", "megan.trainer@example.com", "Celebrity trainer focusing on body transformation and lifestyle coaching.", 200.00, false, "Beverly Hills", "CA", "US", 34.0736, -118.4004, []int{0, 2}},
		{"Ryan", "Jackson", "ryan.trainer@example.com", "Endurance coach for runners and triathletes. Boston Marathon qualifier.", 85.00, true, "San Francisco", "CA", "US", 37.7749, -122.4194, []int{3, 4}},
		{"Brittany", "White", "brittany.trainer@example.com", "Yoga and strength fusion specialist. RYT-500 with CSCS certification.", 90.00, true, "Seattle", "WA", "US", 47.6062, -122.3321, []int{4, 5}},
		{"Jordan", "Harris", "jordan.trainer@example.com", "Functional movement specialist focusing on injury prevention.", 95.00, true, "Portland", "OR", "US", 45.5152, -122.6784, []int{3, 5}},
		{"Kayla", "Martin", "kayla.trainer@example.com", "HIIT and metabolic conditioning expert. Transform your fitness in 12 weeks.", 75.00, true, "San Diego", "CA", "US", 32.7157, -117.1611, []int{1, 4}},
		{"Austin", "Garcia", "austin.trainer@example.com", "MMA conditioning coach. Training fighters for over 8 years.", 110.00, true, "Las Vegas", "NV", "US", 36.1699, -115.1398, []int{1, 3}},
		{"Taylor", "Martinez", "taylor.trainer@example.com", "Outdoor fitness specialist. Hiking, climbing, and adventure training.", 80.00, true, "Denver", "CO", "US", 39.7392, -104.9903, []int{3, 4}},

		// US - Central
		{"Cameron", "Robinson", "cameron.trainer@example.com", "Powerlifting coach with multiple state records. Technique-focused training.", 100.00, true, "Chicago", "IL", "US", 41.8781, -87.6298, []int{0, 2}},
		{"Morgan", "Clark", "morgan.trainer@example.com", "Sports performance coach working with high school and college athletes.", 85.00, true, "Houston", "TX", "US", 29.7604, -95.3698, []int{1, 3}},
		{"Dakota", "Rodriguez", "dakota.trainer@example.com", "Weight loss transformation specialist. Lost 100lbs myself and now help others.", 70.00, true, "Dallas", "TX", "US", 32.7767, -96.7970, []int{0, 4}},
		{"Casey", "Lewis", "casey.trainer@example.com", "Golf fitness specialist improving your game through targeted training.", 95.00, true, "Austin", "TX", "US", 30.2672, -97.7431, []int{3, 5}},
		{"Avery", "Lee", "avery.trainer@example.com", "Swimming and aquatic fitness coach. Former Olympic trials qualifier.", 90.00, true, "Minneapolis", "MN", "US", 44.9778, -93.2650, []int{3, 4}},
		{"Riley", "Walker", "riley.trainer@example.com", "Youth sports performance specialist. Age-appropriate training for young athletes.", 65.00, true, "Kansas City", "MO", "US", 39.0997, -94.5786, []int{1, 3}},

		// UK
		{"Oliver", "Thompson", "oliver.trainer.uk@example.com", "Level 4 personal trainer specializing in strength and conditioning.", 65.00, true, "London", "England", "GB", 51.5074, -0.1278, []int{0, 1}},
		{"Charlotte", "Davies", "charlotte.trainer.uk@example.com", "Women's health and fitness specialist. Hormonal health focus.", 70.00, true, "Manchester", "England", "GB", 53.4808, -2.2426, []int{0, 5}},
		{"George", "Evans", "george.trainer.uk@example.com", "Boxing fitness coach bringing ring training to everyday fitness.", 60.00, true, "Birmingham", "England", "GB", 52.4862, -1.8904, []int{1, 4}},
		{"Amelia", "Wilson", "amelia.trainer.uk@example.com", "Rehabilitation specialist helping clients recover from injuries.", 75.00, true, "Leeds", "England", "GB", 53.8008, -1.5491, []int{5}},
		{"Harry", "Thomas", "harry.trainer.uk@example.com", "Rugby strength and conditioning coach for club and national level.", 80.00, true, "Glasgow", "Scotland", "GB", 55.8642, -4.2518, []int{1, 3}},
		{"Isla", "Roberts", "isla.trainer.uk@example.com", "Kettlebell specialist and StrongFirst certified instructor.", 55.00, true, "Edinburgh", "Scotland", "GB", 55.9533, -3.1883, []int{0, 1}},

		// Canada
		{"Liam", "Campbell", "liam.trainer.ca@example.com", "Hockey performance specialist. On and off-ice conditioning.", 90.00, true, "Toronto", "ON", "CA", 43.6532, -79.3832, []int{1, 3}},
		{"Emma", "Stewart", "emma.trainer.ca@example.com", "Ski and snowboard conditioning coach. Year-round athletic prep.", 85.00, true, "Vancouver", "BC", "CA", 49.2827, -123.1207, []int{3, 4}},
		{"Noah", "Anderson", "noah.trainer.ca@example.com", "Bilingual trainer offering services in English and French.", 75.00, true, "Montreal", "QC", "CA", 45.5017, -73.5673, []int{0, 1}},
		{"Olivia", "Brown", "olivia.trainer.ca@example.com", "Mountain athlete coach. Climbing, hiking, and alpine training.", 80.00, true, "Calgary", "AB", "CA", 51.0447, -114.0719, []int{3, 4}},

		// Australia
		{"Jack", "Smith", "jack.trainer.au@example.com", "Beach body specialist combining surfing fitness with strength training.", 85.00, true, "Sydney", "NSW", "AU", -33.8688, 151.2093, []int{0, 4}},
		{"Chloe", "Jones", "chloe.trainer.au@example.com", "AFL performance coach for amateur and semi-pro players.", 90.00, true, "Melbourne", "VIC", "AU", -37.8136, 144.9631, []int{1, 3}},
		{"William", "Taylor", "william.trainer.au@example.com", "Cricket fitness specialist improving batting and bowling power.", 75.00, true, "Brisbane", "QLD", "AU", -27.4698, 153.0251, []int{1, 3}},
		{"Sophie", "Williams", "sophie.trainer.au@example.com", "Outdoor bootcamp specialist. Group and individual training.", 60.00, true, "Perth", "WA", "AU", -31.9505, 115.8605, []int{1, 4}},

		// Europe
		{"Lucas", "Martin", "lucas.trainer.fr@example.com", "Triathlon coach with Ironman experience. Endurance specialist.", 95.00, true, "Paris", "Île-de-France", "FR", 48.8566, 2.3522, []int{3, 4}},
		{"Léa", "Bernard", "lea.trainer.fr@example.com", "Dance fitness fusion trainer. Ballet meets bootcamp.", 70.00, true, "Lyon", "Auvergne-Rhône-Alpes", "FR", 45.7640, 4.8357, []int{4, 5}},
		{"Maximilian", "Schmidt", "max.trainer.de@example.com", "German national team conditioning consultant. Elite performance training.", 150.00, false, "Berlin", "Berlin", "DE", 52.5200, 13.4050, []int{1, 3}},
		{"Anna", "Mueller", "anna.trainer.de@example.com", "Functional fitness and mobility specialist. Move better, feel better.", 80.00, true, "Munich", "Bavaria", "DE", 48.1351, 11.5820, []int{3, 5}},
		{"Daan", "de Vries", "daan.trainer.nl@example.com", "Cycling performance coach for road and track cyclists.", 85.00, true, "Amsterdam", "North Holland", "NL", 52.3676, 4.9041, []int{3, 4}},
		{"Carlos", "Fernandez", "carlos.trainer.es@example.com", "Football academy trainer. Youth development specialist.", 75.00, true, "Barcelona", "Catalonia", "ES", 41.3851, 2.1734, []int{1, 3}},
		{"Maria", "Garcia", "maria.trainer.es@example.com", "Padel and tennis fitness specialist. Court performance training.", 70.00, true, "Madrid", "Madrid", "ES", 40.4168, -3.7038, []int{1, 4}},
		{"Marco", "Rossi", "marco.trainer.it@example.com", "Italian national swimming coach. Aquatic excellence.", 100.00, true, "Rome", "Lazio", "IT", 41.9028, 12.4964, []int{3, 4}},

		// Asia
		{"Yuki", "Tanaka", "yuki.trainer.jp@example.com", "Martial arts conditioning specialist. Karate and judo background.", 90.00, true, "Tokyo", "Tokyo", "JP", 35.6762, 139.6503, []int{1, 4}},
		{"Kenji", "Suzuki", "kenji.trainer.jp@example.com", "Baseball performance coach. Pitching and batting power specialist.", 85.00, true, "Osaka", "Osaka", "JP", 34.6937, 135.5023, []int{1, 3}},
		{"Wei", "Chen", "wei.trainer.sg@example.com", "Corporate wellness specialist. Helping busy professionals stay fit.", 100.00, true, "Singapore", "Singapore", "SG", 1.3521, 103.8198, []int{0, 4}},
		{"Ming", "Wong", "ming.trainer.hk@example.com", "Functional training and TRX specialist. Small space workouts.", 95.00, true, "Hong Kong", "Hong Kong", "HK", 22.3193, 114.1694, []int{0, 3}},
		{"Ji-Young", "Kim", "jiyoung.trainer.kr@example.com", "K-pop idol training methodology. Dance and fitness fusion.", 110.00, true, "Seoul", "Seoul", "KR", 37.5665, 126.9780, []int{4, 5}},
		{"Mei-Ling", "Lin", "meiling.trainer.tw@example.com", "Holistic fitness combining Eastern and Western training methods.", 75.00, true, "Taipei", "Taiwan", "TW", 25.0330, 121.5654, []int{0, 5}},
	}

	hashedPassword, _ := utils.HashPassword("Password123!")

	for _, t := range trainerData {
		// Check if trainer user already exists
		var existingUser models.User
		if err := DB.Where("email = ?", t.Email).First(&existingUser).Error; err == nil {
			log.Printf("Trainer user already exists: %s", t.Email)
			continue
		}

		// Create user for trainer
		user := models.User{
			Email:             t.Email,
			Password:          hashedPassword,
			FirstName:         t.FirstName,
			LastName:          t.LastName,
			Provider:          "local",
			IsActive:          true,
			ProfileVisibility: "public",
		}

		if err := DB.Create(&user).Error; err != nil {
			log.Printf("Failed to create trainer user %s: %v", t.Email, err)
			continue
		}

		// Create trainer profile
		profile := models.TrainerProfile{
			UserID:              user.ID,
			Bio:                 t.Bio,
			HourlyRate:          floatPtr(t.HourlyRate),
			Visibility:          "public",
			IsLookingForClients: t.IsLookingForClients,
			Location: models.Location{
				City:        strPtr(t.City),
				Region:      strPtr(t.Region),
				CountryCode: strPtr(t.CountryCode),
				Latitude:    floatPtr(t.Latitude),
				Longitude:   floatPtr(t.Longitude),
			},
		}

		if err := DB.Create(&profile).Error; err != nil {
			log.Printf("Failed to create trainer profile for %s: %v", t.Email, err)
			continue
		}

		// Assign specialties
		var specialties []models.Specialty
		for _, idx := range t.SpecialtyIndices {
			if idx < len(allSpecialties) {
				specialties = append(specialties, allSpecialties[idx])
			}
		}
		if len(specialties) > 0 {
			DB.Model(&profile).Association("Specialties").Append(&specialties)
		}

		log.Printf("Created searchable trainer: %s %s (%s) - $%.2f/hr", t.FirstName, t.LastName, t.City, t.HourlyRate)
	}

	log.Println("Finished seeding searchable trainers")
}

func SeedDatabase() {
	log.Println("Starting database seeding...")
	SeedMuscleGroups()
	SeedEquipment()
	SeedExerciseTypes()
	SeedExercises()
	SeedFitnessLevels()
	SeedFitnessGoals()
	SeedGlobalRPEScale()

	if err := SeedRoles(DB); err != nil {
		log.Printf("Failed to seed Role data: %v", err)
	}

	// Seed RBAC data
	if err := SeedRBACData(DB); err != nil {
		log.Printf("Failed to seed RBAC data: %v", err)
	}

	// Migrate existing users to roles
	if err := MigrateExistingUsersToRoles(DB); err != nil {
		log.Printf("Failed to migrate users to roles: %v", err)
	}

	// Seed specialties (must be before trainer profiles)
	SeedSpecialties()

	// Seed trainer profiles
	SeedTrainerProfiles()

	// Seed trainer-client links (must be after trainer profiles)
	SeedTrainerClientLinks()

	// Seed searchable users and trainers for location search testing
	SeedSearchableUsers()
	SeedSearchableTrainers()

	log.Println("Database seeding completed!")
}
