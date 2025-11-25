-- Rollback unified weight system migration

-- Remove comments
COMMENT ON COLUMN users.preferred_weight_unit IS NULL;
COMMENT ON COLUMN users.preferred_height_unit IS NULL;
COMMENT ON COLUMN users.preferred_distance_unit IS NULL;

COMMENT ON COLUMN user_fitness_profiles.current_weight_kg IS NULL;
COMMENT ON COLUMN user_fitness_profiles.original_current_weight_value IS NULL;
COMMENT ON COLUMN user_fitness_profiles.original_current_weight_unit IS NULL;

COMMENT ON COLUMN user_fitness_profiles.target_weight_kg IS NULL;
COMMENT ON COLUMN user_fitness_profiles.original_target_weight_value IS NULL;
COMMENT ON COLUMN user_fitness_profiles.original_target_weight_unit IS NULL;

COMMENT ON COLUMN workout_prescriptions.target_weight_kg IS NULL;
COMMENT ON COLUMN workout_prescriptions.original_target_weight_value IS NULL;
COMMENT ON COLUMN workout_prescriptions.original_target_weight_unit IS NULL;

COMMENT ON COLUMN session_sets.actual_weight_kg IS NULL;
COMMENT ON COLUMN session_sets.original_actual_weight_value IS NULL;
COMMENT ON COLUMN session_sets.original_actual_weight_unit IS NULL;

-- Drop index
DROP INDEX IF EXISTS idx_users_preferred_weight_unit;

-- Revert session_sets table changes
ALTER TABLE session_sets
DROP COLUMN IF EXISTS original_actual_weight_value,
DROP COLUMN IF EXISTS original_actual_weight_unit;

-- Note: We keep actual_weight_kg as it may have existed before

-- Revert workout_prescriptions table changes
ALTER TABLE workout_prescriptions
DROP COLUMN IF EXISTS original_target_weight_value,
DROP COLUMN IF EXISTS original_target_weight_unit;

-- Restore weight_kg column (Note: Data will be lost)
ALTER TABLE workout_prescriptions
ADD COLUMN IF NOT EXISTS weight_kg DECIMAL(10,2);

-- Revert user_fitness_profiles table changes
ALTER TABLE user_fitness_profiles
DROP COLUMN IF EXISTS original_current_weight_value,
DROP COLUMN IF EXISTS original_current_weight_unit,
DROP COLUMN IF EXISTS original_target_weight_value,
DROP COLUMN IF EXISTS original_target_weight_unit;

-- Restore preferred_weight_unit to user_fitness_profiles
ALTER TABLE user_fitness_profiles
ADD COLUMN IF NOT EXISTS preferred_weight_unit VARCHAR(2) DEFAULT 'kg';

-- Remove unit preference fields from users table
ALTER TABLE users
DROP COLUMN IF EXISTS preferred_weight_unit,
DROP COLUMN IF EXISTS preferred_height_unit,
DROP COLUMN IF EXISTS preferred_distance_unit;
