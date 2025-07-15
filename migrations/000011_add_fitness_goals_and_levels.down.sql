-- Drop indexes
DROP INDEX IF EXISTS idx_user_fitness_goals_priority;
DROP INDEX IF EXISTS idx_user_fitness_goals_fitness_goal_id;
DROP INDEX IF EXISTS idx_user_fitness_goals_user_id;
DROP INDEX IF EXISTS idx_users_fitness_level_id;
DROP INDEX IF EXISTS idx_fitness_goals_category;
DROP INDEX IF EXISTS idx_fitness_goals_name;
DROP INDEX IF EXISTS idx_fitness_levels_sort_order;
DROP INDEX IF EXISTS idx_fitness_levels_name;

-- Drop user_fitness_goals table
DROP TABLE IF EXISTS user_fitness_goals;

-- Remove fitness_level_id from users table
ALTER TABLE users DROP COLUMN IF EXISTS fitness_level_id;

-- Drop fitness_goals table
DROP TABLE IF EXISTS fitness_goals;

-- Drop fitness_levels table
DROP TABLE IF EXISTS fitness_levels;