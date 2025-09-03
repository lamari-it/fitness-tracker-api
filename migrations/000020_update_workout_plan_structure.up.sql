-- Add new columns to workout_plans table
ALTER TABLE workout_plans 
ADD COLUMN IF NOT EXISTS template_weeks INTEGER NOT NULL DEFAULT 1;

-- Update workouts table structure
-- First, add the new columns
ALTER TABLE workouts
ADD COLUMN IF NOT EXISTS user_id UUID,
ADD COLUMN IF NOT EXISTS description TEXT,
ADD COLUMN IF NOT EXISTS visibility VARCHAR(20) DEFAULT 'private';

-- Set user_id from the workout_plan before dropping plan_id
UPDATE workouts w
SET user_id = wp.user_id
FROM workout_plans wp
WHERE w.plan_id = wp.id
AND w.user_id IS NULL;

-- Make user_id NOT NULL after populating it (only if there are rows)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM workouts LIMIT 1) THEN
        ALTER TABLE workouts ALTER COLUMN user_id SET NOT NULL;
    END IF;
END $$;

-- Add foreign key constraint for user_id
ALTER TABLE workouts
ADD CONSTRAINT workouts_user_id_fkey 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Create workout_plan_items table to link workouts to plans
CREATE TABLE IF NOT EXISTS workout_plan_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL,
    workout_id UUID NOT NULL,
    week_index INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (plan_id) REFERENCES workout_plans(id) ON DELETE CASCADE,
    FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE,
    UNIQUE(plan_id, workout_id)
);

-- Migrate existing workout-plan relationships to workout_plan_items
INSERT INTO workout_plan_items (plan_id, workout_id, week_index)
SELECT DISTINCT plan_id, id as workout_id, 0 as week_index
FROM workouts
WHERE plan_id IS NOT NULL
ON CONFLICT DO NOTHING;

-- Now drop the old columns and constraints from workouts
ALTER TABLE workouts
DROP CONSTRAINT IF EXISTS workouts_plan_id_fkey,
DROP COLUMN IF EXISTS plan_id,
DROP COLUMN IF EXISTS day_number,
DROP COLUMN IF EXISTS notes;

-- Create plan_enrollments table
CREATE TABLE IF NOT EXISTS plan_enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL,
    user_id UUID NOT NULL,
    start_date TIMESTAMP NOT NULL,
    days_per_week INTEGER NOT NULL CHECK (days_per_week >= 1 AND days_per_week <= 7),
    current_index INTEGER NOT NULL DEFAULT 0,
    schedule_mode VARCHAR(20) DEFAULT 'rolling',
    preferred_weekdays INTEGER[] DEFAULT '{}',
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (plan_id) REFERENCES workout_plans(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_workout_plan_items_plan_id ON workout_plan_items(plan_id);
CREATE INDEX IF NOT EXISTS idx_workout_plan_items_workout_id ON workout_plan_items(workout_id);
CREATE INDEX IF NOT EXISTS idx_workout_plan_items_week_index ON workout_plan_items(week_index);
CREATE INDEX IF NOT EXISTS idx_plan_enrollments_plan_id ON plan_enrollments(plan_id);
CREATE INDEX IF NOT EXISTS idx_plan_enrollments_user_id ON plan_enrollments(user_id);
CREATE INDEX IF NOT EXISTS idx_plan_enrollments_status ON plan_enrollments(status);
CREATE INDEX IF NOT EXISTS idx_workouts_user_id ON workouts(user_id);

-- Add unique constraint for user-plan enrollment (one active enrollment per user per plan)
CREATE UNIQUE INDEX IF NOT EXISTS unique_active_enrollment 
ON plan_enrollments(user_id, plan_id) 
WHERE status = 'active';