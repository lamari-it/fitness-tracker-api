-- Create RPE scales table
CREATE TABLE rpe_scales (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    min_value INT NOT NULL DEFAULT 1,
    max_value INT NOT NULL DEFAULT 10,
    is_global BOOLEAN NOT NULL DEFAULT false,
    trainer_id UUID REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT chk_valid_range CHECK (min_value < max_value),
    CONSTRAINT chk_global_or_trainer CHECK (
        (is_global = true AND trainer_id IS NULL) OR
        (is_global = false AND trainer_id IS NOT NULL)
    )
);

-- Create RPE scale values table
CREATE TABLE rpe_scale_values (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scale_id UUID NOT NULL REFERENCES rpe_scales(id) ON DELETE CASCADE,
    value INT NOT NULL,
    label VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT uq_scale_value UNIQUE (scale_id, value)
);

-- Indexes
CREATE INDEX idx_rpe_scales_trainer ON rpe_scales(trainer_id);
CREATE INDEX idx_rpe_scales_global ON rpe_scales(is_global);
CREATE INDEX idx_rpe_scales_deleted_at ON rpe_scales(deleted_at);
CREATE INDEX idx_rpe_scale_values_scale ON rpe_scale_values(scale_id);

-- Add RPE value references to workout_exercises and set_logs
ALTER TABLE workout_exercises ADD COLUMN target_rpe_value_id UUID REFERENCES rpe_scale_values(id) ON DELETE SET NULL;
ALTER TABLE set_logs ADD COLUMN rpe_value_id UUID REFERENCES rpe_scale_values(id) ON DELETE SET NULL;

-- Indexes for the new columns
CREATE INDEX idx_workout_exercises_target_rpe ON workout_exercises(target_rpe_value_id);
CREATE INDEX idx_set_logs_rpe_value ON set_logs(rpe_value_id);
