package database

// GetExerciseEquipmentMappings returns equipment mappings for each exercise
func GetExerciseEquipmentMappings() map[string][]struct {
	Slug     string
	Optional bool
	Notes    string
} {
	return map[string][]struct {
		Slug     string
		Optional bool
		Notes    string
	}{
		// Chest Exercises
		"Push-ups": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"Bench Press": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "flat_bench", Optional: false, Notes: "Required for the exercise"},
			{Slug: "squat_rack", Optional: true, Notes: "Optional for safety when lifting heavy"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading the barbell"},
		},
		"Incline Bench Press": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "incline_bench", Optional: false, Notes: "Or adjustable bench set to incline"},
			{Slug: "squat_rack", Optional: true, Notes: "Optional for safety when lifting heavy"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading the barbell"},
		},
		"Decline Bench Press": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "decline_bench", Optional: false, Notes: "Or adjustable bench set to decline"},
			{Slug: "squat_rack", Optional: true, Notes: "Optional for safety when lifting heavy"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading the barbell"},
		},
		"Dumbbell Bench Press": {
			{Slug: "dumbbells", Optional: false, Notes: "Required for the exercise"},
			{Slug: "flat_bench", Optional: false, Notes: "Or adjustable bench"},
		},
		"Dumbbell Flyes": {
			{Slug: "dumbbells", Optional: false, Notes: "Required for the exercise"},
			{Slug: "flat_bench", Optional: false, Notes: "Or adjustable bench"},
		},
		"Chest Dips": {
			{Slug: "dip_station", Optional: false, Notes: "Parallel bars required"},
			{Slug: "dip_belt", Optional: true, Notes: "For adding weight"},
		},
		"Cable Chest Flyes": {
			{Slug: "cable_machine", Optional: false, Notes: "Cable machine with adjustable pulleys"},
			{Slug: "d_handle", Optional: false, Notes: "Single handles for each side"},
		},
		"Pec Deck": {
			{Slug: "pec_deck", Optional: false, Notes: "Pec deck machine required"},
		},

		// Back Exercises
		"Pull-ups": {
			{Slug: "pull_up_bar", Optional: false, Notes: "Pull-up bar required"},
			{Slug: "dip_belt", Optional: true, Notes: "For adding weight"},
			{Slug: "resistance_bands", Optional: true, Notes: "For assistance"},
		},
		"Chin-ups": {
			{Slug: "pull_up_bar", Optional: false, Notes: "Pull-up bar required"},
			{Slug: "dip_belt", Optional: true, Notes: "For adding weight"},
			{Slug: "resistance_bands", Optional: true, Notes: "For assistance"},
		},
		"Lat Pulldowns": {
			{Slug: "lat_pulldown", Optional: false, Notes: "Lat pulldown machine required"},
			{Slug: "lat_pulldown_bar", Optional: false, Notes: "Wide grip bar"},
		},
		"Barbell Rows": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading"},
		},
		"Dumbbell Rows": {
			{Slug: "dumbbells", Optional: false, Notes: "Single dumbbell required"},
			{Slug: "flat_bench", Optional: true, Notes: "Optional for support"},
		},
		"T-Bar Rows": {
			{Slug: "t_bar_row", Optional: false, Notes: "T-bar or landmine attachment"},
			{Slug: "barbell", Optional: false, Notes: "If using landmine"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading"},
			{Slug: "v_bar", Optional: true, Notes: "Optional handle attachment"},
		},
		"Cable Rows": {
			{Slug: "cable_machine", Optional: false, Notes: "Cable machine required"},
			{Slug: "straight_bar", Optional: false, Notes: "Or V-bar attachment"},
		},
		"Inverted Rows": {
			{Slug: "squat_rack", Optional: false, Notes: "Or smith machine set at appropriate height"},
			{Slug: "barbell", Optional: false, Notes: "Bar to hold onto"},
		},
		"Shrugs": {
			{Slug: "dumbbells", Optional: false, Notes: "Or barbell"},
			{Slug: "barbell", Optional: false, Notes: "Alternative to dumbbells"},
			{Slug: "plates", Optional: true, Notes: "If using barbell"},
		},
		"Face Pulls": {
			{Slug: "cable_machine", Optional: false, Notes: "Cable machine required"},
			{Slug: "rope_attachment", Optional: false, Notes: "Rope attachment for cables"},
		},

		// Leg Exercises
		"Squats": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for marking position"},
		},
		"Barbell Squats": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "squat_rack", Optional: false, Notes: "Or power rack for safety"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading"},
		},
		"Front Squats": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "squat_rack", Optional: false, Notes: "Or power rack for safety"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading"},
		},
		"Goblet Squats": {
			{Slug: "dumbbells", Optional: false, Notes: "Single dumbbell or kettlebell"},
			{Slug: "kettlebell", Optional: false, Notes: "Alternative to dumbbell"},
		},
		"Lunges": {
			{Slug: "dumbbells", Optional: true, Notes: "Optional for added resistance"},
		},
		"Bulgarian Split Squats": {
			{Slug: "flat_bench", Optional: false, Notes: "Or box for rear foot elevation"},
			{Slug: "dumbbells", Optional: true, Notes: "Optional for added resistance"},
		},
		"Leg Press": {
			{Slug: "leg_press", Optional: false, Notes: "Leg press machine required"},
		},
		"Leg Extensions": {
			{Slug: "leg_extension", Optional: false, Notes: "Leg extension machine required"},
		},
		"Deadlifts": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading"},
		},
		"Romanian Deadlifts": {
			{Slug: "barbell", Optional: false, Notes: "Or dumbbells"},
			{Slug: "dumbbells", Optional: false, Notes: "Alternative to barbell"},
			{Slug: "plates", Optional: true, Notes: "If using barbell"},
		},
		"Sumo Deadlifts": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading"},
		},
		"Stiff Leg Deadlifts": {
			{Slug: "barbell", Optional: false, Notes: "Or dumbbells"},
			{Slug: "dumbbells", Optional: false, Notes: "Alternative to barbell"},
			{Slug: "plates", Optional: true, Notes: "If using barbell"},
		},
		"Leg Curls": {
			{Slug: "leg_curl", Optional: false, Notes: "Leg curl machine required"},
		},
		"Hip Thrusts": {
			{Slug: "flat_bench", Optional: false, Notes: "Bench for shoulder support"},
			{Slug: "barbell", Optional: false, Notes: "Or dumbbell for resistance"},
			{Slug: "plates", Optional: true, Notes: "If using barbell"},
			{Slug: "exercise_mat", Optional: true, Notes: "For padding"},
		},
		"Glute Bridges": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"Calf Raises": {
			{Slug: "step_platform", Optional: true, Notes: "Optional for greater range of motion"},
			{Slug: "dumbbells", Optional: true, Notes: "Optional for added resistance"},
		},
		"Seated Calf Raises": {
			{Slug: "calf_raise", Optional: false, Notes: "Seated calf raise machine"},
			{Slug: "plates", Optional: true, Notes: "If machine uses plates"},
		},

		// Shoulder Exercises
		"Overhead Press": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "squat_rack", Optional: true, Notes: "For getting bar into position"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading"},
		},
		"Dumbbell Shoulder Press": {
			{Slug: "dumbbells", Optional: false, Notes: "Required for the exercise"},
			{Slug: "adjustable_bench", Optional: true, Notes: "Optional for seated variation"},
		},
		"Lateral Raises": {
			{Slug: "dumbbells", Optional: false, Notes: "Required for the exercise"},
		},
		"Rear Delt Flyes": {
			{Slug: "dumbbells", Optional: false, Notes: "Required for the exercise"},
			{Slug: "flat_bench", Optional: true, Notes: "Optional for chest support"},
		},
		"Front Raises": {
			{Slug: "dumbbells", Optional: false, Notes: "Or barbell/plate"},
			{Slug: "barbell", Optional: false, Notes: "Alternative to dumbbells"},
			{Slug: "plates", Optional: false, Notes: "Can use single plate"},
		},
		"Arnold Press": {
			{Slug: "dumbbells", Optional: false, Notes: "Required for the exercise"},
			{Slug: "adjustable_bench", Optional: true, Notes: "Optional for seated variation"},
		},
		"Upright Rows": {
			{Slug: "barbell", Optional: false, Notes: "Or EZ-bar"},
			{Slug: "ez_bar", Optional: false, Notes: "Alternative to straight bar"},
			{Slug: "plates", Optional: true, Notes: "If using barbell"},
		},
		"Pike Push-ups": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"Handstand Push-ups": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
			{Slug: "wall_bars", Optional: true, Notes: "Wall for support"},
		},

		// Arm Exercises
		"Bicep Curls": {
			{Slug: "dumbbells", Optional: false, Notes: "Or barbell"},
			{Slug: "barbell", Optional: false, Notes: "Alternative to dumbbells"},
			{Slug: "plates", Optional: true, Notes: "If using barbell"},
		},
		"Hammer Curls": {
			{Slug: "dumbbells", Optional: false, Notes: "Required for the exercise"},
		},
		"Preacher Curls": {
			{Slug: "preacher_bench", Optional: false, Notes: "Or preacher curl machine"},
			{Slug: "barbell", Optional: false, Notes: "Or EZ-bar or dumbbells"},
			{Slug: "ez_bar", Optional: false, Notes: "Alternative to straight bar"},
			{Slug: "plates", Optional: true, Notes: "If using barbell"},
		},
		"Concentration Curls": {
			{Slug: "dumbbells", Optional: false, Notes: "Single dumbbell required"},
			{Slug: "flat_bench", Optional: true, Notes: "For seated position"},
		},
		"Tricep Dips": {
			{Slug: "dip_station", Optional: false, Notes: "Or bench for bench dips"},
			{Slug: "flat_bench", Optional: false, Notes: "Alternative for bench dips"},
			{Slug: "dip_belt", Optional: true, Notes: "For adding weight"},
		},
		"Tricep Pushdowns": {
			{Slug: "cable_machine", Optional: false, Notes: "Cable machine required"},
			{Slug: "rope_attachment", Optional: false, Notes: "Or straight bar attachment"},
			{Slug: "straight_bar", Optional: false, Notes: "Alternative to rope"},
		},
		"Overhead Tricep Extension": {
			{Slug: "dumbbells", Optional: false, Notes: "Single dumbbell or EZ-bar"},
			{Slug: "ez_bar", Optional: false, Notes: "Alternative to dumbbell"},
			{Slug: "plates", Optional: true, Notes: "If using barbell"},
		},
		"Close-Grip Bench Press": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "flat_bench", Optional: false, Notes: "Required for the exercise"},
			{Slug: "squat_rack", Optional: true, Notes: "For safety"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading"},
		},
		"Diamond Push-ups": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"Wrist Curls": {
			{Slug: "barbell", Optional: false, Notes: "Or dumbbells"},
			{Slug: "dumbbells", Optional: false, Notes: "Alternative to barbell"},
			{Slug: "flat_bench", Optional: true, Notes: "For forearm support"},
			{Slug: "plates", Optional: true, Notes: "If using barbell"},
		},
		"Reverse Curls": {
			{Slug: "barbell", Optional: false, Notes: "Or EZ-bar"},
			{Slug: "ez_bar", Optional: false, Notes: "More comfortable grip"},
			{Slug: "plates", Optional: true, Notes: "If using barbell"},
		},
		"Farmer's Walks": {
			{Slug: "dumbbells", Optional: false, Notes: "Or farmer's walk handles"},
			{Slug: "farmers_walk", Optional: false, Notes: "Specialized handles"},
			{Slug: "kettlebell", Optional: false, Notes: "Alternative option"},
		},

		// Core Exercises
		"Planks": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"Crunches": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"Bicycle Crunches": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"Russian Twists": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
			{Slug: "medicine_ball", Optional: true, Notes: "Or dumbbell for added resistance"},
			{Slug: "dumbbells", Optional: true, Notes: "Alternative to medicine ball"},
		},
		"Side Planks": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"Leg Raises": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"Hanging Leg Raises": {
			{Slug: "pull_up_bar", Optional: false, Notes: "Pull-up bar required"},
			{Slug: "ab_straps", Optional: true, Notes: "For easier grip"},
		},
		"Mountain Climbers": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"Dead Bug": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"Hollow Body Hold": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"Ab Wheel Rollouts": {
			{Slug: "ab_wheel", Optional: false, Notes: "Ab wheel required"},
			{Slug: "exercise_mat", Optional: true, Notes: "For knee comfort"},
		},
		"Superman": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"Good Mornings": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "squat_rack", Optional: true, Notes: "For getting bar into position"},
			{Slug: "plates", Optional: true, Notes: "For added resistance"},
		},
		"Hyperextensions": {
			{Slug: "hyperextension_bench", Optional: false, Notes: "Or back extension machine"},
			{Slug: "plates", Optional: true, Notes: "For added resistance"},
		},

		// Full Body and Cardio
		"Burpees": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"Thrusters": {
			{Slug: "dumbbells", Optional: false, Notes: "Or barbell"},
			{Slug: "barbell", Optional: false, Notes: "Alternative to dumbbells"},
			{Slug: "plates", Optional: true, Notes: "If using barbell"},
		},
		"Man Makers": {
			{Slug: "dumbbells", Optional: false, Notes: "Required for the exercise"},
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"Turkish Get-ups": {
			{Slug: "dumbbells", Optional: false, Notes: "Or kettlebell"},
			{Slug: "kettlebell", Optional: false, Notes: "Traditional option"},
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"Kettlebell Swings": {
			{Slug: "kettlebell", Optional: false, Notes: "Required for the exercise"},
		},
		"Kettlebell Snatches": {
			{Slug: "kettlebell", Optional: false, Notes: "Required for the exercise"},
		},
		"Kettlebell Clean and Press": {
			{Slug: "kettlebell", Optional: false, Notes: "Required for the exercise"},
		},
		"Box Jumps": {
			{Slug: "plyo_box", Optional: false, Notes: "Plyometric box required"},
		},
		"Jump Squats": {
			{Slug: "dumbbells", Optional: true, Notes: "Optional for added resistance"},
		},
		"High Knees": {},
		"Jumping Jacks": {},
		"Bear Crawls": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"Crab Walks": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"Treadmill Running": {
			{Slug: "treadmill", Optional: false, Notes: "Treadmill required"},
		},
		"Stationary Bike": {
			{Slug: "stationary_bike", Optional: false, Notes: "Stationary bike required"},
		},
		"Rowing Machine": {
			{Slug: "rowing_machine", Optional: false, Notes: "Rowing machine required"},
		},
		"Elliptical Machine": {
			{Slug: "elliptical", Optional: false, Notes: "Elliptical machine required"},
		},
		"Jump Rope": {
			{Slug: "jump_rope", Optional: false, Notes: "Jump rope required"},
		},
		"Stair Climbing": {
			{Slug: "stair_climber", Optional: false, Notes: "Stair climber machine or actual stairs"},
		},

		// Olympic Lifts
		"Clean and Jerk": {
			{Slug: "barbell", Optional: false, Notes: "Olympic barbell recommended"},
			{Slug: "plates", Optional: false, Notes: "Bumper plates recommended"},
			{Slug: "lifting_belt", Optional: true, Notes: "For heavy lifts"},
			{Slug: "chalk", Optional: true, Notes: "For better grip"},
		},
		"Snatch": {
			{Slug: "barbell", Optional: false, Notes: "Olympic barbell recommended"},
			{Slug: "plates", Optional: false, Notes: "Bumper plates recommended"},
			{Slug: "lifting_belt", Optional: true, Notes: "For heavy lifts"},
			{Slug: "chalk", Optional: true, Notes: "For better grip"},
		},
		"Power Clean": {
			{Slug: "barbell", Optional: false, Notes: "Olympic barbell recommended"},
			{Slug: "plates", Optional: false, Notes: "Bumper plates recommended"},
			{Slug: "lifting_belt", Optional: true, Notes: "For heavy lifts"},
		},
		"Hang Clean": {
			{Slug: "barbell", Optional: false, Notes: "Olympic barbell recommended"},
			{Slug: "plates", Optional: false, Notes: "Bumper plates recommended"},
		},
		"Push Press": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "squat_rack", Optional: true, Notes: "For getting bar into position"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading"},
		},
		"Push Jerk": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "squat_rack", Optional: true, Notes: "For getting bar into position"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading"},
		},

		// Isometric and Stability
		"Wall Sit": {},
		"Glute Bridge Hold": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"Single-Leg Glute Bridge": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"Single-Leg Deadlift": {
			{Slug: "dumbbells", Optional: true, Notes: "Optional for added resistance"},
			{Slug: "kettlebell", Optional: true, Notes: "Alternative to dumbbells"},
		},
		"Pistol Squats": {},
		"Single-Leg Calf Raises": {
			{Slug: "step_platform", Optional: true, Notes: "For greater range of motion"},
		},
		"Bird Dog": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},

		// Stretching and Mobility (most don't need equipment)
		"Cat-Cow Stretch": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"Child's Pose": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"Downward Dog": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"Pigeon Pose": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"Cobra Stretch": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"Figure-4 Stretch": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"Seated Forward Fold": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
			{Slug: "stretching_strap", Optional: true, Notes: "For assistance"},
		},
		"Standing Quad Stretch": {},
		"Standing Calf Stretch": {},
		"Shoulder Rolls": {},
		"Arm Circles": {},
		"Neck Rolls": {},
		"Hip Circles": {},
		"Leg Swings": {},
	}
}