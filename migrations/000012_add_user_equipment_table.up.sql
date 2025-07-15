-- Add slug field to equipment table
ALTER TABLE equipment ADD COLUMN IF NOT EXISTS slug VARCHAR(100) UNIQUE;

-- Create index on slug for faster lookups
CREATE INDEX IF NOT EXISTS idx_equipment_slug ON equipment(slug);

-- Create user_equipment table to track what equipment users have at home and gym
CREATE TABLE IF NOT EXISTS user_equipment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    equipment_id UUID NOT NULL REFERENCES equipment(id) ON DELETE CASCADE,
    location_type VARCHAR(10) NOT NULL CHECK (location_type IN ('home', 'gym')),
    gym_location TEXT, -- Optional: specific gym location/branch
    notes TEXT, -- Optional: notes about the equipment (e.g., "20lb dumbbells", "adjustable 5-50lbs")
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure a user can't have duplicate equipment entries for the same location
    UNIQUE(user_id, equipment_id, location_type)
);

-- Create indexes for faster queries
CREATE INDEX idx_user_equipment_user_id ON user_equipment(user_id);
CREATE INDEX idx_user_equipment_location_type ON user_equipment(location_type);
CREATE INDEX idx_user_equipment_user_location ON user_equipment(user_id, location_type);

-- Add trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_user_equipment_updated_at BEFORE UPDATE
    ON user_equipment FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();