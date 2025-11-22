-- Add name_slug column to fitness_goals table
ALTER TABLE fitness_goals ADD COLUMN IF NOT EXISTS name_slug VARCHAR(100) NOT NULL DEFAULT '';
CREATE UNIQUE INDEX IF NOT EXISTS idx_fitness_goals_name_slug ON fitness_goals(name_slug) WHERE deleted_at IS NULL;

-- Update existing fitness goals with slugs
UPDATE fitness_goals SET name_slug = 'weight_loss' WHERE name = 'Weight Loss';
UPDATE fitness_goals SET name_slug = 'muscle_gain' WHERE name = 'Muscle Gain';
UPDATE fitness_goals SET name_slug = 'endurance' WHERE name = 'Endurance';
UPDATE fitness_goals SET name_slug = 'strength' WHERE name = 'Strength';
UPDATE fitness_goals SET name_slug = 'flexibility' WHERE name = 'Flexibility';
UPDATE fitness_goals SET name_slug = 'general_fitness' WHERE name = 'General Fitness';
UPDATE fitness_goals SET name_slug = 'athletic_performance' WHERE name = 'Athletic Performance';
UPDATE fitness_goals SET name_slug = 'rehabilitation' WHERE name = 'Rehabilitation';
UPDATE fitness_goals SET name_slug = 'body_recomposition' WHERE name = 'Body Recomposition';
UPDATE fitness_goals SET name_slug = 'stress_relief' WHERE name = 'Stress Relief';

-- Alter user_fitness_goals table to reference user_fitness_profiles instead of users
-- First, add the new column
ALTER TABLE user_fitness_goals ADD COLUMN IF NOT EXISTS user_fitness_profile_id UUID;

-- Migrate existing data: link through user's fitness profile
UPDATE user_fitness_goals ufg
SET user_fitness_profile_id = ufp.id
FROM user_fitness_profiles ufp
WHERE ufg.user_id = ufp.user_id;

-- Drop old constraints and indexes
DROP INDEX IF EXISTS unique_user_fitness_goal_combo;
ALTER TABLE user_fitness_goals DROP CONSTRAINT IF EXISTS fk_user_fitness_goals_user;

-- Remove user_id column and add new constraint
ALTER TABLE user_fitness_goals DROP COLUMN IF EXISTS user_id;
ALTER TABLE user_fitness_goals ALTER COLUMN user_fitness_profile_id SET NOT NULL;

-- Add foreign key constraint
ALTER TABLE user_fitness_goals
ADD CONSTRAINT fk_user_fitness_goals_profile
FOREIGN KEY (user_fitness_profile_id)
REFERENCES user_fitness_profiles(id)
ON DELETE CASCADE;

-- Create new unique index
CREATE UNIQUE INDEX idx_unique_profile_fitness_goal_combo
ON user_fitness_goals(user_fitness_profile_id, fitness_goal_id)
WHERE deleted_at IS NULL;

-- Remove primary_goal column from user_fitness_profiles
ALTER TABLE user_fitness_profiles DROP COLUMN IF EXISTS primary_goal;
