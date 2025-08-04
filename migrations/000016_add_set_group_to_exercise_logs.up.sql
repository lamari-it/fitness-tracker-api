-- Add set_group_id column to exercise_logs table
ALTER TABLE exercise_logs ADD COLUMN set_group_id UUID;

-- Add foreign key constraint (nullable to allow free-form workouts)
ALTER TABLE exercise_logs 
    ADD CONSTRAINT fk_exercise_logs_set_group
    FOREIGN KEY (set_group_id) REFERENCES set_groups(id) ON DELETE SET NULL;

-- Create index for set_group_id
CREATE INDEX idx_exercise_logs_set_group_id ON exercise_logs(set_group_id);