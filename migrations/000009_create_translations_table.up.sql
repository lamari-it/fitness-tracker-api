-- Translations Table
CREATE TABLE translations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_type VARCHAR(50) NOT NULL,
    resource_id UUID NOT NULL,
    field_name VARCHAR(50) NOT NULL,
    language VARCHAR(5) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_translations_resource_type ON translations(resource_type);
CREATE INDEX idx_translations_resource_id ON translations(resource_id);
CREATE INDEX idx_translations_field_name ON translations(field_name);
CREATE INDEX idx_translations_language ON translations(language);
CREATE INDEX idx_translations_deleted_at ON translations(deleted_at);

-- Unique constraint for translation combinations
CREATE UNIQUE INDEX unique_translation_combo ON translations(resource_type, resource_id, field_name, language) WHERE deleted_at IS NULL;