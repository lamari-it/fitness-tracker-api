-- Drop check constraint
ALTER TABLE set_logs DROP CONSTRAINT IF EXISTS check_weight_unit;

-- Drop weight_unit column
ALTER TABLE set_logs DROP COLUMN IF EXISTS weight_unit;