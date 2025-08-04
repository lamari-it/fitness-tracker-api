-- Drop index
DROP INDEX IF EXISTS idx_workout_exercises_set_group_id;

-- Drop foreign key constraint and column
ALTER TABLE workout_exercises 
    DROP CONSTRAINT IF EXISTS fk_workout_exercises_set_group,
    DROP COLUMN IF EXISTS set_group_id;