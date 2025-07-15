-- Drop trigger
DROP TRIGGER IF EXISTS update_user_equipment_updated_at ON user_equipment;

-- Drop table
DROP TABLE IF EXISTS user_equipment;

-- Drop indexes
DROP INDEX IF EXISTS idx_equipment_slug;

-- Remove slug column from equipment table
ALTER TABLE equipment DROP COLUMN IF EXISTS slug;