-- Add weight_unit column to set_logs table
ALTER TABLE set_logs ADD COLUMN weight_unit VARCHAR(5) DEFAULT 'kg';

-- Add check constraint to ensure valid weight units
ALTER TABLE set_logs ADD CONSTRAINT check_weight_unit 
    CHECK (weight_unit IN ('kg', 'lbs', 'lb'));