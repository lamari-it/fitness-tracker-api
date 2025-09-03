-- Drop indexes
DROP INDEX IF EXISTS unique_active_enrollment;
DROP INDEX IF EXISTS idx_workouts_user_id;
DROP INDEX IF EXISTS idx_plan_enrollments_status;
DROP INDEX IF EXISTS idx_plan_enrollments_user_id;
DROP INDEX IF EXISTS idx_plan_enrollments_plan_id;
DROP INDEX IF EXISTS idx_workout_plan_items_week_index;
DROP INDEX IF EXISTS idx_workout_plan_items_workout_id;
DROP INDEX IF EXISTS idx_workout_plan_items_plan_id;

-- Drop plan_enrollments table
DROP TABLE IF EXISTS plan_enrollments;

-- Restore workouts table to original structure
-- First add back the old columns
ALTER TABLE workouts
ADD COLUMN IF NOT EXISTS plan_id UUID,
ADD COLUMN IF NOT EXISTS day_number INTEGER NOT NULL DEFAULT 1,
ADD COLUMN IF NOT EXISTS notes TEXT;

-- Restore plan_id values from workout_plan_items
UPDATE workouts w
SET plan_id = wpi.plan_id
FROM workout_plan_items wpi
WHERE w.id = wpi.workout_id
AND w.plan_id IS NULL;

-- Add back the foreign key constraint
ALTER TABLE workouts
ADD CONSTRAINT workouts_plan_id_fkey 
FOREIGN KEY (plan_id) REFERENCES workout_plans(id) ON DELETE CASCADE;

-- Drop the workout_plan_items table
DROP TABLE IF EXISTS workout_plan_items;

-- Remove new columns from workouts
ALTER TABLE workouts
DROP CONSTRAINT IF EXISTS workouts_user_id_fkey,
DROP COLUMN IF EXISTS user_id,
DROP COLUMN IF EXISTS description,
DROP COLUMN IF EXISTS visibility;

-- Remove template_weeks from workout_plans
ALTER TABLE workout_plans
DROP COLUMN IF EXISTS template_weeks;