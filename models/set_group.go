package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SetGroupType string

const (
	SetGroupTypeStraight  SetGroupType = "straight"
	SetGroupTypeSuperset  SetGroupType = "superset"
	SetGroupTypeCircuit   SetGroupType = "circuit"
	SetGroupTypeGiantSet  SetGroupType = "giant_set"
	SetGroupTypeDropSet   SetGroupType = "drop_set"
	SetGroupTypePyramid   SetGroupType = "pyramid"
	SetGroupTypeRestPause SetGroupType = "rest_pause"
)

type SetGroup struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	WorkoutID       uuid.UUID      `gorm:"type:uuid;not null" json:"workout_id"`
	GroupType       SetGroupType   `gorm:"type:varchar(20);not null;default:'straight'" json:"group_type"`
	Name            string         `gorm:"type:varchar(255)" json:"name"`
	Notes           string         `gorm:"type:text" json:"notes"`
	OrderNumber     int            `gorm:"not null" json:"order_number"`
	RestBetweenSets int            `json:"rest_between_sets"`
	Rounds          int            `gorm:"default:1" json:"rounds"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Workout          Workout           `gorm:"foreignKey:WorkoutID;constraint:OnDelete:CASCADE" json:"workout,omitempty"`
	WorkoutExercises []WorkoutExercise `gorm:"foreignKey:SetGroupID" json:"exercises,omitempty"`
}

func (sg *SetGroup) BeforeCreate(tx *gorm.DB) (err error) {
	if sg.ID == uuid.Nil {
		sg.ID = uuid.New()
	}
	return
}
