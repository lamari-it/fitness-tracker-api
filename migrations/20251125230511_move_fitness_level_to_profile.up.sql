-- Move fitness_level_id from users table to user_fitness_profiles table

-- Step 1: Add fitness_level_id column to user_fitness_profiles
ALTER TABLE user_fitness_profiles
ADD COLUMN fitness_level_id UUID REFERENCES fitness_levels(id) ON DELETE SET NULL;

-- Step 2: Migrate existing fitness_level_id values from users to their fitness profiles
UPDATE user_fitness_profiles ufp
SET fitness_level_id = u.fitness_level_id
FROM users u
WHERE ufp.user_id = u.id
AND u.fitness_level_id IS NOT NULL;

-- Step 3: Drop the foreign key constraint from users table
ALTER TABLE users
DROP CONSTRAINT IF EXISTS fk_users_fitness_level;

-- Step 4: Drop the fitness_level_id column from users table
ALTER TABLE users
DROP COLUMN IF EXISTS fitness_level_id;
