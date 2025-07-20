-- Add slug column to exercises table
ALTER TABLE exercises ADD COLUMN slug VARCHAR(255);

-- Create unique index on slug
CREATE UNIQUE INDEX idx_exercises_slug ON exercises(slug);

-- Update existing exercises with slugs based on their names
UPDATE exercises 
SET slug = LOWER(
    REPLACE(
        REPLACE(
            REPLACE(name, ' ', '-'),
            '''', ''
        ),
        '--', '-'
    )
);

-- Make slug column NOT NULL after populating it
ALTER TABLE exercises ALTER COLUMN slug SET NOT NULL;