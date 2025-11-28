package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                    uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email                 string    `gorm:"unique;not null" json:"email"`
	Password              string    `gorm:"not null" json:"-"`
	FirstName             string    `gorm:"not null" json:"first_name"`
	LastName              string    `gorm:"not null" json:"last_name"`
	Provider              string    `gorm:"default:'local'" json:"provider"`
	GoogleID              *string   `gorm:"unique" json:"google_id,omitempty"`
	AppleID               *string   `gorm:"unique" json:"apple_id,omitempty"`
	IsActive              bool      `gorm:"default:true" json:"is_active"`
	IsAdmin               bool      `gorm:"default:false" json:"is_admin"`
	PreferredWeightUnit   string    `gorm:"type:varchar(2);default:'kg'" json:"preferred_weight_unit"`
	PreferredHeightUnit   string    `gorm:"type:varchar(5);default:'cm'" json:"preferred_height_unit"`
	PreferredDistanceUnit string    `gorm:"type:varchar(2);default:'km'" json:"preferred_distance_unit"`
	// Profile & Privacy
	ProfileVisibility   string         `gorm:"type:varchar(20);default:'private'" json:"profile_visibility"` // public, private, friends_only
	IsLookingForTrainer bool           `gorm:"default:false" json:"is_looking_for_trainer"`
	Bio                 string         `gorm:"type:text" json:"bio,omitempty"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Roles []Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

type UserResponse struct {
	ID                    uuid.UUID `json:"id"`
	Email                 string    `json:"email"`
	FirstName             string    `json:"first_name"`
	LastName              string    `json:"last_name"`
	Provider              string    `json:"provider"`
	Roles                 []Role    `json:"roles,omitempty"`
	IsActive              bool      `json:"is_active"`
	IsAdmin               bool      `json:"is_admin"`
	PreferredWeightUnit   string    `json:"preferred_weight_unit"`
	PreferredHeightUnit   string    `json:"preferred_height_unit"`
	PreferredDistanceUnit string    `json:"preferred_distance_unit"`
	ProfileVisibility     string    `json:"profile_visibility"`
	IsLookingForTrainer   bool      `json:"is_looking_for_trainer"`
	Bio                   string    `json:"bio,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// UserDiscoveryResponse is a limited response for public user discovery
type UserDiscoveryResponse struct {
	ID                  uuid.UUID `json:"id"`
	FirstName           string    `json:"first_name"`
	LastName            string    `json:"last_name"`
	Bio                 string    `json:"bio,omitempty"`
	IsLookingForTrainer bool      `json:"is_looking_for_trainer"`
	CreatedAt           time.Time `json:"created_at"`
}

// UpdateUserSettingsRequest is used for updating user preferences
type UpdateUserSettingsRequest struct {
	PreferredWeightUnit   string `json:"preferred_weight_unit" binding:"omitempty,oneof=kg lb"`
	PreferredHeightUnit   string `json:"preferred_height_unit" binding:"omitempty,oneof=cm ft"`
	PreferredDistanceUnit string `json:"preferred_distance_unit" binding:"omitempty,oneof=km mi"`
	FirstName             string `json:"first_name" binding:"omitempty,min=1,max=100"`
	LastName              string `json:"last_name" binding:"omitempty,min=1,max=100"`
	ProfileVisibility     string `json:"profile_visibility" binding:"omitempty,oneof=public private friends_only"`
	IsLookingForTrainer   *bool  `json:"is_looking_for_trainer"`
	Bio                   string `json:"bio" binding:"omitempty,max=500"`
}

func (u *User) ToResponse() UserResponse {
	response := UserResponse{
		ID:                    u.ID,
		Email:                 u.Email,
		FirstName:             u.FirstName,
		LastName:              u.LastName,
		Provider:              u.Provider,
		IsActive:              u.IsActive,
		IsAdmin:               u.IsAdmin,
		PreferredWeightUnit:   u.PreferredWeightUnit,
		PreferredHeightUnit:   u.PreferredHeightUnit,
		PreferredDistanceUnit: u.PreferredDistanceUnit,
		ProfileVisibility:     u.ProfileVisibility,
		IsLookingForTrainer:   u.IsLookingForTrainer,
		Bio:                   u.Bio,
		CreatedAt:             u.CreatedAt,
		UpdatedAt:             u.UpdatedAt,
	}

	if len(u.Roles) > 0 {
		response.Roles = u.Roles
	}

	return response
}

// ToDiscoveryResponse returns a limited response for public user discovery
func (u *User) ToDiscoveryResponse() UserDiscoveryResponse {
	return UserDiscoveryResponse{
		ID:                  u.ID,
		FirstName:           u.FirstName,
		LastName:            u.LastName,
		Bio:                 u.Bio,
		IsLookingForTrainer: u.IsLookingForTrainer,
		CreatedAt:           u.CreatedAt,
	}
}
