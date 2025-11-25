-- Add unit preference fields to users table
ALTER TABLE users
ADD COLUMN IF NOT EXISTS preferred_weight_unit VARCHAR(2) DEFAULT 'kg',
ADD COLUMN IF NOT EXISTS preferred_height_unit VARCHAR(5) DEFAULT 'cm',
ADD COLUMN IF NOT EXISTS preferred_distance_unit VARCHAR(2) DEFAULT 'km';

-- Update user_fitness_profiles table for unified weight system
-- Remove preferred_weight_unit (moved to users table)
ALTER TABLE user_fitness_profiles
DROP COLUMN IF EXISTS preferred_weight_unit;

-- Update current_weight fields to canonical + original storage
ALTER TABLE user_fitness_profiles
ADD COLUMN IF NOT EXISTS original_current_weight_value DECIMAL(6,2),
ADD COLUMN IF NOT EXISTS original_current_weight_unit VARCHAR(2);

-- Rename current_weight_kg to match naming convention (if exists)
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'user_fitness_profiles'
        AND column_name = 'current_weight_kg'
    ) THEN
        -- Modify existing column type
        ALTER TABLE user_fitness_profiles
        ALTER COLUMN current_weight_kg TYPE DECIMAL(6,2);
    ELSE
        -- Create new column if it doesn't exist
        ALTER TABLE user_fitness_profiles
        ADD COLUMN current_weight_kg DECIMAL(6,2);
    END IF;
END $$;

-- Update target_weight fields to canonical + original storage
ALTER TABLE user_fitness_profiles
ADD COLUMN IF NOT EXISTS original_target_weight_value DECIMAL(6,2),
ADD COLUMN IF NOT EXISTS original_target_weight_unit VARCHAR(2);

-- Ensure target_weight_kg exists with correct type
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'user_fitness_profiles'
        AND column_name = 'target_weight_kg'
    ) THEN
        ALTER TABLE user_fitness_profiles
        ALTER COLUMN target_weight_kg TYPE DECIMAL(6,2);
    ELSE
        ALTER TABLE user_fitness_profiles
        ADD COLUMN target_weight_kg DECIMAL(6,2);
    END IF;
END $$;

-- Update workout_prescriptions table
-- Remove unused weight_kg column
ALTER TABLE workout_prescriptions
DROP COLUMN IF EXISTS weight_kg;

-- Update target_weight to canonical + original storage
ALTER TABLE workout_prescriptions
ADD COLUMN IF NOT EXISTS original_target_weight_value DECIMAL(6,2),
ADD COLUMN IF NOT EXISTS original_target_weight_unit VARCHAR(2);

-- Ensure target_weight_kg exists with correct type
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'workout_prescriptions'
        AND column_name = 'target_weight_kg'
    ) THEN
        ALTER TABLE workout_prescriptions
        ALTER COLUMN target_weight_kg TYPE DECIMAL(6,2);
    ELSE
        ALTER TABLE workout_prescriptions
        ADD COLUMN target_weight_kg DECIMAL(6,2);
    END IF;
END $$;

-- Update session_sets table for unified weight system
ALTER TABLE session_sets
ADD COLUMN IF NOT EXISTS original_actual_weight_value DECIMAL(6,2),
ADD COLUMN IF NOT EXISTS original_actual_weight_unit VARCHAR(2);

-- Ensure actual_weight_kg exists with correct type
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'session_sets'
        AND column_name = 'actual_weight_kg'
    ) THEN
        ALTER TABLE session_sets
        ALTER COLUMN actual_weight_kg TYPE DECIMAL(6,2);
    ELSE
        ALTER TABLE session_sets
        ADD COLUMN actual_weight_kg DECIMAL(6,2);
    END IF;
END $$;

-- Create index on preferred weight unit for performance
CREATE INDEX IF NOT EXISTS idx_users_preferred_weight_unit ON users(preferred_weight_unit);

-- Add comments for documentation
COMMENT ON COLUMN users.preferred_weight_unit IS 'User preference for weight display (kg or lb)';
COMMENT ON COLUMN users.preferred_height_unit IS 'User preference for height display';
COMMENT ON COLUMN users.preferred_distance_unit IS 'User preference for distance display';

COMMENT ON COLUMN user_fitness_profiles.current_weight_kg IS 'Canonical weight storage in kg';
COMMENT ON COLUMN user_fitness_profiles.original_current_weight_value IS 'Original weight value as entered by user';
COMMENT ON COLUMN user_fitness_profiles.original_current_weight_unit IS 'Original weight unit as entered by user (kg or lb)';

COMMENT ON COLUMN user_fitness_profiles.target_weight_kg IS 'Canonical target weight storage in kg';
COMMENT ON COLUMN user_fitness_profiles.original_target_weight_value IS 'Original target weight value as entered by user';
COMMENT ON COLUMN user_fitness_profiles.original_target_weight_unit IS 'Original target weight unit as entered by user (kg or lb)';

COMMENT ON COLUMN workout_prescriptions.target_weight_kg IS 'Canonical target weight storage in kg';
COMMENT ON COLUMN workout_prescriptions.original_target_weight_value IS 'Original target weight value as entered';
COMMENT ON COLUMN workout_prescriptions.original_target_weight_unit IS 'Original target weight unit as entered (kg or lb)';

COMMENT ON COLUMN session_sets.actual_weight_kg IS 'Canonical actual weight storage in kg';
COMMENT ON COLUMN session_sets.original_actual_weight_value IS 'Original actual weight value as logged';
COMMENT ON COLUMN session_sets.original_actual_weight_unit IS 'Original actual weight unit as logged (kg or lb)';
