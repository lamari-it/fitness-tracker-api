-- Exercises Table
CREATE TABLE exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    equipment TEXT,
    is_bodyweight BOOLEAN DEFAULT false,
    instructions TEXT,
    video_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Exercise Muscle Groups Junction Table
CREATE TABLE exercise_muscle_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exercise_id UUID NOT NULL,
    muscle_group_id UUID NOT NULL,
    "primary" BOOLEAN DEFAULT false,
    intensity VARCHAR(20) DEFAULT 'moderate',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE,
    FOREIGN KEY (muscle_group_id) REFERENCES muscle_groups(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_exercises_name ON exercises(name);
CREATE INDEX idx_exercises_equipment ON exercises(equipment);
CREATE INDEX idx_exercises_is_bodyweight ON exercises(is_bodyweight);
CREATE INDEX idx_exercise_muscle_groups_exercise_id ON exercise_muscle_groups(exercise_id);
CREATE INDEX idx_exercise_muscle_groups_muscle_group_id ON exercise_muscle_groups(muscle_group_id);
CREATE INDEX idx_exercise_muscle_groups_primary ON exercise_muscle_groups("primary");
CREATE INDEX idx_exercise_muscle_groups_intensity ON exercise_muscle_groups(intensity);

-- Unique constraint for exercise-muscle group combination
CREATE UNIQUE INDEX unique_exercise_muscle_combo ON exercise_muscle_groups(exercise_id, muscle_group_id);