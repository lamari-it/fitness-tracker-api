-- Rollback: Move fitness_level_id back from user_fitness_profiles to users table

-- Step 1: Add fitness_level_id column back to users table
ALTER TABLE users
ADD COLUMN fitness_level_id UUID REFERENCES fitness_levels(id) ON DELETE SET NULL;

-- Step 2: Migrate fitness_level_id values from profiles back to users
UPDATE users u
SET fitness_level_id = ufp.fitness_level_id
FROM user_fitness_profiles ufp
WHERE u.id = ufp.user_id
AND ufp.fitness_level_id IS NOT NULL;

-- Step 3: Drop the fitness_level_id column from user_fitness_profiles
ALTER TABLE user_fitness_profiles
DROP COLUMN IF EXISTS fitness_level_id;
