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
		"push_ups": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"bench_press": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "flat_bench", Optional: false, Notes: "Required for the exercise"},
			{Slug: "squat_rack", Optional: true, Notes: "Optional for safety when lifting heavy"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading the barbell"},
		},
		"incline_bench_press": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "incline_bench", Optional: false, Notes: "Or adjustable bench set to incline"},
			{Slug: "squat_rack", Optional: true, Notes: "Optional for safety when lifting heavy"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading the barbell"},
		},
		"decline_bench_press": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "decline_bench", Optional: false, Notes: "Or adjustable bench set to decline"},
			{Slug: "squat_rack", Optional: true, Notes: "Optional for safety when lifting heavy"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading the barbell"},
		},
		"dumbbell_bench_press": {
			{Slug: "dumbbells", Optional: false, Notes: "Required for the exercise"},
			{Slug: "flat_bench", Optional: false, Notes: "Or adjustable bench"},
		},
		"dumbbell_flyes": {
			{Slug: "dumbbells", Optional: false, Notes: "Required for the exercise"},
			{Slug: "flat_bench", Optional: false, Notes: "Or adjustable bench"},
		},
		"chest_dips": {
			{Slug: "dip_station", Optional: false, Notes: "Parallel bars required"},
			{Slug: "dip_belt", Optional: true, Notes: "For adding weight"},
		},
		"cable_chest_flyes": {
			{Slug: "cable_machine", Optional: false, Notes: "Cable machine with adjustable pulleys"},
			{Slug: "d_handle", Optional: false, Notes: "Single handles for each side"},
		},
		"pec_deck": {
			{Slug: "pec_deck", Optional: false, Notes: "Pec deck machine required"},
		},

		// Back Exercises
		"pull_ups": {
			{Slug: "pull_up_bar", Optional: false, Notes: "Pull-up bar required"},
			{Slug: "dip_belt", Optional: true, Notes: "For adding weight"},
			{Slug: "resistance_bands", Optional: true, Notes: "For assistance"},
		},
		"chin_ups": {
			{Slug: "pull_up_bar", Optional: false, Notes: "Pull-up bar required"},
			{Slug: "dip_belt", Optional: true, Notes: "For adding weight"},
			{Slug: "resistance_bands", Optional: true, Notes: "For assistance"},
		},
		"lat_pulldowns": {
			{Slug: "lat_pulldown", Optional: false, Notes: "Lat pulldown machine required"},
			{Slug: "lat_pulldown_bar", Optional: false, Notes: "Wide grip bar"},
		},
		"barbell_rows": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading"},
		},
		"dumbbell_rows": {
			{Slug: "dumbbells", Optional: false, Notes: "Single dumbbell required"},
			{Slug: "flat_bench", Optional: true, Notes: "Optional for support"},
		},
		"t_bar_rows": {
			{Slug: "t_bar_row", Optional: false, Notes: "T-bar or landmine attachment"},
			{Slug: "barbell", Optional: false, Notes: "If using landmine"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading"},
			{Slug: "v_bar", Optional: true, Notes: "Optional handle attachment"},
		},
		"cable_rows": {
			{Slug: "cable_machine", Optional: false, Notes: "Cable machine required"},
			{Slug: "straight_bar", Optional: false, Notes: "Or V-bar attachment"},
		},
		"inverted_rows": {
			{Slug: "squat_rack", Optional: false, Notes: "Or smith machine set at appropriate height"},
			{Slug: "barbell", Optional: false, Notes: "Bar to hold onto"},
		},
		"shrugs": {
			{Slug: "dumbbells", Optional: false, Notes: "Or barbell"},
			{Slug: "barbell", Optional: false, Notes: "Alternative to dumbbells"},
			{Slug: "plates", Optional: true, Notes: "If using barbell"},
		},
		"face_pulls": {
			{Slug: "cable_machine", Optional: false, Notes: "Cable machine required"},
			{Slug: "rope_attachment", Optional: false, Notes: "Rope attachment for cables"},
		},

		// Leg Exercises
		"squats": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for marking position"},
		},
		"barbell_squats": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "squat_rack", Optional: false, Notes: "Or power rack for safety"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading"},
		},
		"front_squats": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "squat_rack", Optional: false, Notes: "Or power rack for safety"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading"},
		},
		"goblet_squats": {
			{Slug: "dumbbells", Optional: false, Notes: "Single dumbbell or kettlebell"},
			{Slug: "kettlebell", Optional: false, Notes: "Alternative to dumbbell"},
		},
		"lunges": {
			{Slug: "dumbbells", Optional: true, Notes: "Optional for added resistance"},
		},
		"bulgarian_split_squats": {
			{Slug: "flat_bench", Optional: false, Notes: "Or box for rear foot elevation"},
			{Slug: "dumbbells", Optional: true, Notes: "Optional for added resistance"},
		},
		"leg_press": {
			{Slug: "leg_press", Optional: false, Notes: "Leg press machine required"},
		},
		"leg_extensions": {
			{Slug: "leg_extension", Optional: false, Notes: "Leg extension machine required"},
		},
		"deadlifts": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading"},
		},
		"romanian_deadlifts": {
			{Slug: "barbell", Optional: false, Notes: "Or dumbbells"},
			{Slug: "dumbbells", Optional: false, Notes: "Alternative to barbell"},
			{Slug: "plates", Optional: true, Notes: "If using barbell"},
		},
		"sumo_deadlifts": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading"},
		},
		"stiff_leg_deadlifts": {
			{Slug: "barbell", Optional: false, Notes: "Or dumbbells"},
			{Slug: "dumbbells", Optional: false, Notes: "Alternative to barbell"},
			{Slug: "plates", Optional: true, Notes: "If using barbell"},
		},
		"leg_curls": {
			{Slug: "leg_curl", Optional: false, Notes: "Leg curl machine required"},
		},
		"hip_thrusts": {
			{Slug: "flat_bench", Optional: false, Notes: "Bench for shoulder support"},
			{Slug: "barbell", Optional: false, Notes: "Or dumbbell for resistance"},
			{Slug: "plates", Optional: true, Notes: "If using barbell"},
			{Slug: "exercise_mat", Optional: true, Notes: "For padding"},
		},
		"glute_bridges": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"calf_raises": {
			{Slug: "step_platform", Optional: true, Notes: "Optional for greater range of motion"},
			{Slug: "dumbbells", Optional: true, Notes: "Optional for added resistance"},
		},
		"seated_calf_raises": {
			{Slug: "calf_raise", Optional: false, Notes: "Seated calf raise machine"},
			{Slug: "plates", Optional: true, Notes: "If machine uses plates"},
		},

		// Shoulder Exercises
		"overhead_press": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "squat_rack", Optional: true, Notes: "For getting bar into position"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading"},
		},
		"dumbbell_shoulder_press": {
			{Slug: "dumbbells", Optional: false, Notes: "Required for the exercise"},
			{Slug: "adjustable_bench", Optional: true, Notes: "Optional for seated variation"},
		},
		"lateral_raises": {
			{Slug: "dumbbells", Optional: false, Notes: "Required for the exercise"},
		},
		"rear_delt_flyes": {
			{Slug: "dumbbells", Optional: false, Notes: "Required for the exercise"},
			{Slug: "flat_bench", Optional: true, Notes: "Optional for chest support"},
		},
		"front_raises": {
			{Slug: "dumbbells", Optional: false, Notes: "Or barbell/plate"},
			{Slug: "barbell", Optional: false, Notes: "Alternative to dumbbells"},
			{Slug: "plates", Optional: false, Notes: "Can use single plate"},
		},
		"arnold_press": {
			{Slug: "dumbbells", Optional: false, Notes: "Required for the exercise"},
			{Slug: "adjustable_bench", Optional: true, Notes: "Optional for seated variation"},
		},
		"upright_rows": {
			{Slug: "barbell", Optional: false, Notes: "Or EZ-bar"},
			{Slug: "ez_bar", Optional: false, Notes: "Alternative to straight bar"},
			{Slug: "plates", Optional: true, Notes: "If using barbell"},
		},
		"pike_push_ups": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"handstand_push_ups": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
			{Slug: "wall_bars", Optional: true, Notes: "Wall for support"},
		},

		// Arm Exercises
		"bicep_curls": {
			{Slug: "dumbbells", Optional: false, Notes: "Or barbell"},
			{Slug: "barbell", Optional: false, Notes: "Alternative to dumbbells"},
			{Slug: "plates", Optional: true, Notes: "If using barbell"},
		},
		"hammer_curls": {
			{Slug: "dumbbells", Optional: false, Notes: "Required for the exercise"},
		},
		"preacher_curls": {
			{Slug: "preacher_bench", Optional: false, Notes: "Or preacher curl machine"},
			{Slug: "barbell", Optional: false, Notes: "Or EZ-bar or dumbbells"},
			{Slug: "ez_bar", Optional: false, Notes: "Alternative to straight bar"},
			{Slug: "plates", Optional: true, Notes: "If using barbell"},
		},
		"concentration_curls": {
			{Slug: "dumbbells", Optional: false, Notes: "Single dumbbell required"},
			{Slug: "flat_bench", Optional: true, Notes: "For seated position"},
		},
		"tricep_dips": {
			{Slug: "dip_station", Optional: false, Notes: "Or bench for bench dips"},
			{Slug: "flat_bench", Optional: false, Notes: "Alternative for bench dips"},
			{Slug: "dip_belt", Optional: true, Notes: "For adding weight"},
		},
		"tricep_pushdowns": {
			{Slug: "cable_machine", Optional: false, Notes: "Cable machine required"},
			{Slug: "rope_attachment", Optional: false, Notes: "Or straight bar attachment"},
			{Slug: "straight_bar", Optional: false, Notes: "Alternative to rope"},
		},
		"overhead_tricep_extension": {
			{Slug: "dumbbells", Optional: false, Notes: "Single dumbbell or EZ-bar"},
			{Slug: "ez_bar", Optional: false, Notes: "Alternative to dumbbell"},
			{Slug: "plates", Optional: true, Notes: "If using barbell"},
		},
		"close_grip_bench_press": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "flat_bench", Optional: false, Notes: "Required for the exercise"},
			{Slug: "squat_rack", Optional: true, Notes: "For safety"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading"},
		},
		"diamond_push_ups": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"wrist_curls": {
			{Slug: "barbell", Optional: false, Notes: "Or dumbbells"},
			{Slug: "dumbbells", Optional: false, Notes: "Alternative to barbell"},
			{Slug: "flat_bench", Optional: true, Notes: "For forearm support"},
			{Slug: "plates", Optional: true, Notes: "If using barbell"},
		},
		"reverse_curls": {
			{Slug: "barbell", Optional: false, Notes: "Or EZ-bar"},
			{Slug: "ez_bar", Optional: false, Notes: "More comfortable grip"},
			{Slug: "plates", Optional: true, Notes: "If using barbell"},
		},
		"farmers_walks": {
			{Slug: "dumbbells", Optional: false, Notes: "Or farmer's walk handles"},
			{Slug: "farmers_walk", Optional: false, Notes: "Specialized handles"},
			{Slug: "kettlebell", Optional: false, Notes: "Alternative option"},
		},

		// Core Exercises
		"planks": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"crunches": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"bicycle_crunches": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"russian_twists": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
			{Slug: "medicine_ball", Optional: true, Notes: "Or dumbbell for added resistance"},
			{Slug: "dumbbells", Optional: true, Notes: "Alternative to medicine ball"},
		},
		"side_planks": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"leg_raises": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"hanging_leg_raises": {
			{Slug: "pull_up_bar", Optional: false, Notes: "Pull-up bar required"},
			{Slug: "ab_straps", Optional: true, Notes: "For easier grip"},
		},
		"mountain_climbers": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"dead_bug": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"hollow_body_hold": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"ab_wheel_rollouts": {
			{Slug: "ab_wheel", Optional: false, Notes: "Ab wheel required"},
			{Slug: "exercise_mat", Optional: true, Notes: "For knee comfort"},
		},
		"superman": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"good_mornings": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "squat_rack", Optional: true, Notes: "For getting bar into position"},
			{Slug: "plates", Optional: true, Notes: "For added resistance"},
		},
		"hyperextensions": {
			{Slug: "hyperextension_bench", Optional: false, Notes: "Or back extension machine"},
			{Slug: "plates", Optional: true, Notes: "For added resistance"},
		},

		// Full Body and Cardio
		"burpees": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"thrusters": {
			{Slug: "dumbbells", Optional: false, Notes: "Or barbell"},
			{Slug: "barbell", Optional: false, Notes: "Alternative to dumbbells"},
			{Slug: "plates", Optional: true, Notes: "If using barbell"},
		},
		"man_makers": {
			{Slug: "dumbbells", Optional: false, Notes: "Required for the exercise"},
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"turkish_get_ups": {
			{Slug: "dumbbells", Optional: false, Notes: "Or kettlebell"},
			{Slug: "kettlebell", Optional: false, Notes: "Traditional option"},
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"kettlebell_swings": {
			{Slug: "kettlebell", Optional: false, Notes: "Required for the exercise"},
		},
		"kettlebell_snatches": {
			{Slug: "kettlebell", Optional: false, Notes: "Required for the exercise"},
		},
		"kettlebell_clean_and_press": {
			{Slug: "kettlebell", Optional: false, Notes: "Required for the exercise"},
		},
		"box_jumps": {
			{Slug: "plyo_box", Optional: false, Notes: "Plyometric box required"},
		},
		"jump_squats": {
			{Slug: "dumbbells", Optional: true, Notes: "Optional for added resistance"},
		},
		"high_knees": {},
		"jumping_jacks": {},
		"bear_crawls": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"crab_walks": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"treadmill_running": {
			{Slug: "treadmill", Optional: false, Notes: "Treadmill required"},
		},
		"stationary_bike": {
			{Slug: "stationary_bike", Optional: false, Notes: "Stationary bike required"},
		},
		"rowing_machine": {
			{Slug: "rowing_machine", Optional: false, Notes: "Rowing machine required"},
		},
		"elliptical_machine": {
			{Slug: "elliptical", Optional: false, Notes: "Elliptical machine required"},
		},
		"jump_rope": {
			{Slug: "jump_rope", Optional: false, Notes: "Jump rope required"},
		},
		"stair_climbing": {
			{Slug: "stair_climber", Optional: false, Notes: "Stair climber machine or actual stairs"},
		},

		// Olympic Lifts
		"clean_and_jerk": {
			{Slug: "barbell", Optional: false, Notes: "Olympic barbell recommended"},
			{Slug: "plates", Optional: false, Notes: "Bumper plates recommended"},
			{Slug: "lifting_belt", Optional: true, Notes: "For heavy lifts"},
			{Slug: "chalk", Optional: true, Notes: "For better grip"},
		},
		"snatch": {
			{Slug: "barbell", Optional: false, Notes: "Olympic barbell recommended"},
			{Slug: "plates", Optional: false, Notes: "Bumper plates recommended"},
			{Slug: "lifting_belt", Optional: true, Notes: "For heavy lifts"},
			{Slug: "chalk", Optional: true, Notes: "For better grip"},
		},
		"power_clean": {
			{Slug: "barbell", Optional: false, Notes: "Olympic barbell recommended"},
			{Slug: "plates", Optional: false, Notes: "Bumper plates recommended"},
			{Slug: "lifting_belt", Optional: true, Notes: "For heavy lifts"},
		},
		"hang_clean": {
			{Slug: "barbell", Optional: false, Notes: "Olympic barbell recommended"},
			{Slug: "plates", Optional: false, Notes: "Bumper plates recommended"},
		},
		"push_press": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "squat_rack", Optional: true, Notes: "For getting bar into position"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading"},
		},
		"push_jerk": {
			{Slug: "barbell", Optional: false, Notes: "Required for the exercise"},
			{Slug: "squat_rack", Optional: true, Notes: "For getting bar into position"},
			{Slug: "plates", Optional: false, Notes: "Weight plates for loading"},
		},

		// Isometric and Stability
		"wall_sit": {},
		"glute_bridge_hold": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"single_leg_glute_bridge": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"single_leg_deadlift": {
			{Slug: "dumbbells", Optional: true, Notes: "Optional for added resistance"},
			{Slug: "kettlebell", Optional: true, Notes: "Alternative to dumbbells"},
		},
		"pistol_squats": {},
		"single_leg_calf_raises": {
			{Slug: "step_platform", Optional: true, Notes: "For greater range of motion"},
		},
		"bird_dog": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},

		// Stretching and Mobility (most don't need equipment)
		"cat_cow_stretch": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"childs_pose": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"downward_dog": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"pigeon_pose": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"cobra_stretch": {
			{Slug: "exercise_mat", Optional: true, Notes: "Recommended for comfort"},
		},
		"figure_4_stretch": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
		},
		"seated_forward_fold": {
			{Slug: "exercise_mat", Optional: true, Notes: "Optional for comfort"},
			{Slug: "stretching_strap", Optional: true, Notes: "For assistance"},
		},
		"standing_quad_stretch": {},
		"standing_calf_stretch": {},
		"shoulder_rolls": {},
		"arm_circles": {},
		"neck_rolls": {},
		"hip_circles": {},
		"leg_swings": {},
	}
}