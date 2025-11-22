-- Add primary_goal column back to user_fitness_profiles
ALTER TABLE user_fitness_profiles ADD COLUMN IF NOT EXISTS primary_goal VARCHAR(50) NOT NULL DEFAULT 'general_fitness';

-- Drop new constraints and indexes
DROP INDEX IF EXISTS idx_unique_profile_fitness_goal_combo;
ALTER TABLE user_fitness_goals DROP CONSTRAINT IF EXISTS fk_user_fitness_goals_profile;

-- Add user_id column back
ALTER TABLE user_fitness_goals ADD COLUMN IF NOT EXISTS user_id UUID;

-- Migrate data back: get user_id from the profile
UPDATE user_fitness_goals ufg
SET user_id = ufp.user_id
FROM user_fitness_profiles ufp
WHERE ufg.user_fitness_profile_id = ufp.id;

-- Remove user_fitness_profile_id column
ALTER TABLE user_fitness_goals DROP COLUMN IF EXISTS user_fitness_profile_id;
ALTER TABLE user_fitness_goals ALTER COLUMN user_id SET NOT NULL;

-- Restore old constraint
ALTER TABLE user_fitness_goals
ADD CONSTRAINT fk_user_fitness_goals_user
FOREIGN KEY (user_id)
REFERENCES users(id)
ON DELETE CASCADE;

-- Restore old unique index
CREATE UNIQUE INDEX unique_user_fitness_goal_combo
ON user_fitness_goals(user_id, fitness_goal_id)
WHERE deleted_at IS NULL;

-- Remove name_slug column from fitness_goals
DROP INDEX IF EXISTS idx_fitness_goals_name_slug;
ALTER TABLE fitness_goals DROP COLUMN IF EXISTS name_slug;
