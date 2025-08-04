-- Now make set_group_id NOT NULL and add foreign key constraint
ALTER TABLE workout_exercises 
    ALTER COLUMN set_group_id SET NOT NULL,
    ADD CONSTRAINT fk_workout_exercises_set_group
    FOREIGN KEY (set_group_id) REFERENCES set_groups(id) ON DELETE CASCADE;

-- Create index for set_group_id
CREATE INDEX idx_workout_exercises_set_group_id ON workout_exercises(set_group_id);