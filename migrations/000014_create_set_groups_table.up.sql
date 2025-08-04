-- Set Groups Table
CREATE TABLE set_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_id UUID NOT NULL,
    group_type VARCHAR(20) NOT NULL DEFAULT 'straight',
    name VARCHAR(255),
    notes TEXT,
    order_number INTEGER NOT NULL,
    rest_between_sets INTEGER,
    rounds INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE,
    CONSTRAINT check_group_type CHECK (group_type IN ('straight', 'superset', 'circuit', 'giant_set', 'drop_set', 'pyramid', 'rest_pause'))
);

-- Create indexes for set_groups
CREATE INDEX idx_set_groups_workout_id ON set_groups(workout_id);
CREATE INDEX idx_set_groups_order_number ON set_groups(order_number);
CREATE INDEX idx_set_groups_group_type ON set_groups(group_type);