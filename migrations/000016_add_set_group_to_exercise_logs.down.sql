-- Drop index
DROP INDEX IF EXISTS idx_exercise_logs_set_group_id;

-- Drop foreign key constraint and column
ALTER TABLE exercise_logs 
    DROP CONSTRAINT IF EXISTS fk_exercise_logs_set_group,
    DROP COLUMN IF EXISTS set_group_id;