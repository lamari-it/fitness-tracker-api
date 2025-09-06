package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Translation represents a translatable content entry
type Translation struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	ResourceType string         `gorm:"size:50;not null;uniqueIndex:unique_translation_combo" json:"resource_type"` // e.g., "exercise", "muscle_group", "workout_plan"
	ResourceID   uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:unique_translation_combo" json:"resource_id"`
	FieldName    string         `gorm:"size:50;not null;uniqueIndex:unique_translation_combo" json:"field_name"`    // e.g., "name", "description"
	Language     string         `gorm:"size:5;not null;uniqueIndex:unique_translation_combo" json:"language"`       // e.g., "en", "es", "fr"
	Content      string         `gorm:"type:text;not null" json:"content"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index;uniqueIndex:unique_translation_combo,where:deleted_at IS NULL" json:"deleted_at,omitempty"`
}

// BeforeCreate sets the UUID before creating the translation
func (t *Translation) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for Translation
func (Translation) TableName() string {
	return "translations"
}

// Validate validates the translation data
func (t *Translation) Validate() error {
	if t.ResourceType == "" {
		return gorm.ErrInvalidValue
	}
	if t.ResourceID == uuid.Nil {
		return gorm.ErrInvalidValue
	}
	if t.FieldName == "" {
		return gorm.ErrInvalidValue
	}
	if t.Language == "" {
		return gorm.ErrInvalidValue
	}
	if t.Content == "" {
		return gorm.ErrInvalidValue
	}
	return nil
}

// TranslationResponse represents the response structure for translations
type TranslationResponse struct {
	ID           uuid.UUID `json:"id"`
	ResourceType string    `json:"resource_type"`
	ResourceID   uuid.UUID `json:"resource_id"`
	FieldName    string    `json:"field_name"`
	Language     string    `json:"language"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ToResponse converts Translation to TranslationResponse
func (t *Translation) ToResponse() TranslationResponse {
	return TranslationResponse{
		ID:           t.ID,
		ResourceType: t.ResourceType,
		ResourceID:   t.ResourceID,
		FieldName:    t.FieldName,
		Language:     t.Language,
		Content:      t.Content,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
}

// CreateTranslationRequest represents the request structure for creating translations
type CreateTranslationRequest struct {
	ResourceType string    `json:"resource_type" binding:"required"`
	ResourceID   uuid.UUID `json:"resource_id" binding:"required"`
	FieldName    string    `json:"field_name" binding:"required"`
	Language     string    `json:"language" binding:"required"`
	Content      string    `json:"content" binding:"required"`
}

// UpdateTranslationRequest represents the request structure for updating translations
type UpdateTranslationRequest struct {
	Content string `json:"content" binding:"required"`
}

// MultilingualContent represents content with multiple language versions
type MultilingualContent struct {
	Default     string            `json:"default"`
	Translations map[string]string `json:"translations,omitempty"`
}

// GetContent returns the content for a specific language or default
func (mc *MultilingualContent) GetContent(language string) string {
	if mc.Translations != nil {
		if content, exists := mc.Translations[language]; exists && content != "" {
			return content
		}
	}
	return mc.Default
}

// SetContent sets the content for a specific language
func (mc *MultilingualContent) SetContent(language string, content string) {
	if mc.Translations == nil {
		mc.Translations = make(map[string]string)
	}
	mc.Translations[language] = content
}