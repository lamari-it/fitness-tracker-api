-- Add equipment column back to exercises table
ALTER TABLE exercises ADD COLUMN IF NOT EXISTS equipment TEXT;

-- Drop tables in reverse order
DROP TABLE IF EXISTS exercise_equipment;
DROP TABLE IF EXISTS equipment;