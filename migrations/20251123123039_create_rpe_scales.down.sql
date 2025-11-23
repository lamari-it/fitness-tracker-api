-- Remove RPE value references from workout_exercises and set_logs
ALTER TABLE set_logs DROP COLUMN IF EXISTS rpe_value_id;
ALTER TABLE workout_exercises DROP COLUMN IF EXISTS target_rpe_value_id;

-- Drop RPE tables
DROP TABLE IF EXISTS rpe_scale_values;
DROP TABLE IF EXISTS rpe_scales;
