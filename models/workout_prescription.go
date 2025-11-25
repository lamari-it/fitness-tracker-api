package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PrescriptionType represents the type of exercise grouping/execution pattern
type PrescriptionType string

const (
	PrescriptionTypeStraight  PrescriptionType = "straight"
	PrescriptionTypeSuperset  PrescriptionType = "superset"
	PrescriptionTypeCircuit   PrescriptionType = "circuit"
	PrescriptionTypeGiantSet  PrescriptionType = "giant_set"
	PrescriptionTypeDropSet   PrescriptionType = "drop_set"
	PrescriptionTypePyramid   PrescriptionType = "pyramid"
	PrescriptionTypeRestPause PrescriptionType = "rest_pause"
	PrescriptionTypeAMRAP     PrescriptionType = "amrap"
	PrescriptionTypeEMOM      PrescriptionType = "emom"
	PrescriptionTypeHIIT      PrescriptionType = "hiit"
	PrescriptionTypeWarmup    PrescriptionType = "warmup"
	PrescriptionTypeCooldown  PrescriptionType = "cooldown"
)

// ValidPrescriptionTypes contains all valid prescription types
var ValidPrescriptionTypes = []PrescriptionType{
	PrescriptionTypeStraight,
	PrescriptionTypeSuperset,
	PrescriptionTypeCircuit,
	PrescriptionTypeGiantSet,
	PrescriptionTypeDropSet,
	PrescriptionTypePyramid,
	PrescriptionTypeRestPause,
	PrescriptionTypeAMRAP,
	PrescriptionTypeEMOM,
	PrescriptionTypeHIIT,
	PrescriptionTypeWarmup,
	PrescriptionTypeCooldown,
}

// IsValidPrescriptionType checks if the given type is valid
func IsValidPrescriptionType(t PrescriptionType) bool {
	for _, valid := range ValidPrescriptionTypes {
		if t == valid {
			return true
		}
	}
	return false
}

// WorkoutPrescription represents a unified prescription for exercises in a workout
// This single table handles all set types: straight sets, supersets, circuits, drop sets, etc.
type WorkoutPrescription struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign keys
	WorkoutID  uuid.UUID  `gorm:"type:uuid;not null;index" json:"workout_id"`
	ExerciseID uuid.UUID  `gorm:"type:uuid;not null;index" json:"exercise_id"`
	RPEValueID *uuid.UUID `gorm:"type:uuid" json:"rpe_value_id,omitempty"`

	// Group-level fields (shared across all rows with same GroupID)
	GroupID         uuid.UUID        `gorm:"type:uuid;not null;index" json:"group_id"`
	Type            PrescriptionType `gorm:"type:varchar(50);not null;default:'straight'" json:"type"`
	GroupOrder      int              `gorm:"not null" json:"group_order"`
	GroupRounds     *int             `gorm:"default:1" json:"group_rounds,omitempty"`
	RestBetweenSets *int             `gorm:"" json:"rest_between_sets,omitempty"`
	GroupName       *string          `gorm:"type:varchar(255)" json:"group_name,omitempty"`
	GroupNotes      *string          `gorm:"type:text" json:"group_notes,omitempty"`

	// Exercise-level fields (individual prescription row inside a group)
	ExerciseOrder              int      `gorm:"not null" json:"exercise_order"`
	Sets                       *int     `gorm:"" json:"sets,omitempty"`
	Reps                       *int     `gorm:"" json:"reps,omitempty"`
	HoldSeconds                *int     `gorm:"" json:"hold_seconds,omitempty"`
	TargetWeightKg             *float64 `gorm:"type:decimal(6,2)" json:"-"`
	OriginalTargetWeightValue  *float64 `gorm:"type:decimal(6,2)" json:"-"`
	OriginalTargetWeightUnit   *string  `gorm:"type:varchar(2)" json:"-"`
	Notes                      *string  `gorm:"type:text" json:"notes,omitempty"`

	// Relationships
	Workout  Workout        `gorm:"foreignKey:WorkoutID" json:"-"`
	Exercise Exercise       `gorm:"foreignKey:ExerciseID" json:"exercise,omitempty"`
	RPEValue *RPEScaleValue `gorm:"foreignKey:RPEValueID" json:"rpe_value,omitempty"`
}

// TableName specifies the table name for WorkoutPrescription
func (WorkoutPrescription) TableName() string {
	return "workout_prescriptions"
}

// BeforeCreate validates and generates UUIDs before creating
func (wp *WorkoutPrescription) BeforeCreate(tx *gorm.DB) error {
	// Generate UUIDs if not set
	if wp.ID == uuid.Nil {
		wp.ID = uuid.New()
	}
	if wp.GroupID == uuid.Nil {
		wp.GroupID = uuid.New()
	}

	// Validate prescription type
	if !IsValidPrescriptionType(wp.Type) {
		return errors.New("invalid prescription type")
	}

	// Validate mutual exclusivity of reps and hold_seconds
	hasReps := wp.Reps != nil && *wp.Reps > 0
	hasHold := wp.HoldSeconds != nil && *wp.HoldSeconds > 0

	if hasReps && hasHold {
		return errors.New("prescription cannot have both reps and hold_seconds")
	}

	if !hasReps && !hasHold {
		return errors.New("prescription must have either reps or hold_seconds")
	}

	// Validate group_order is positive
	if wp.GroupOrder < 1 {
		return errors.New("group_order must be at least 1")
	}

	// Validate exercise_order is positive
	if wp.ExerciseOrder < 1 {
		return errors.New("exercise_order must be at least 1")
	}

	return nil
}

// ===== Request DTOs =====

// PrescriptionExerciseRequest represents a single exercise within a prescription group
type PrescriptionExerciseRequest struct {
	ExerciseID    uuid.UUID    `json:"exercise_id" binding:"required"`
	ExerciseOrder int          `json:"exercise_order" binding:"required,min=1"`
	Sets          *int         `json:"sets,omitempty"`
	Reps          *int         `json:"reps,omitempty"`
	HoldSeconds   *int         `json:"hold_seconds,omitempty"`
	TargetWeight  *WeightInput `json:"target_weight,omitempty"`
	RPEValueID    *uuid.UUID   `json:"rpe_value_id,omitempty"`
	Notes         *string      `json:"notes,omitempty"`
}

// CreatePrescriptionGroupRequest represents the request to create a prescription group
type CreatePrescriptionGroupRequest struct {
	GroupID         *uuid.UUID                    `json:"group_id,omitempty"`
	Type            PrescriptionType              `json:"type" binding:"required"`
	GroupOrder      int                           `json:"group_order" binding:"required,min=1"`
	GroupRounds     *int                          `json:"group_rounds,omitempty"`
	RestBetweenSets *int                          `json:"rest_between_sets,omitempty"`
	GroupName       *string                       `json:"group_name,omitempty"`
	GroupNotes      *string                       `json:"group_notes,omitempty"`
	Exercises       []PrescriptionExerciseRequest `json:"exercises" binding:"required,min=1,dive"`
}

// UpdatePrescriptionGroupRequest represents the request to update a prescription group
type UpdatePrescriptionGroupRequest struct {
	Type            *PrescriptionType             `json:"type,omitempty"`
	GroupOrder      *int                          `json:"group_order,omitempty"`
	GroupRounds     *int                          `json:"group_rounds,omitempty"`
	RestBetweenSets *int                          `json:"rest_between_sets,omitempty"`
	GroupName       *string                       `json:"group_name,omitempty"`
	GroupNotes      *string                       `json:"group_notes,omitempty"`
	Exercises       []PrescriptionExerciseRequest `json:"exercises,omitempty"`
}

// ReorderPrescriptionGroupsRequest represents the request to reorder groups
type ReorderPrescriptionGroupsRequest struct {
	GroupOrders []GroupOrderItem `json:"group_orders" binding:"required,min=1,dive"`
}

// GroupOrderItem represents a single group order update
type GroupOrderItem struct {
	GroupID    uuid.UUID `json:"group_id" binding:"required"`
	GroupOrder int       `json:"group_order" binding:"required,min=1"`
}

// AddExerciseToPrescriptionRequest represents adding an exercise to an existing group
type AddExerciseToPrescriptionRequest struct {
	ExerciseID   uuid.UUID    `json:"exercise_id" binding:"required"`
	Sets         *int         `json:"sets,omitempty"`
	Reps         *int         `json:"reps,omitempty"`
	HoldSeconds  *int         `json:"hold_seconds,omitempty"`
	TargetWeight *WeightInput `json:"target_weight,omitempty"`
	RPEValueID   *uuid.UUID   `json:"rpe_value_id,omitempty"`
	Notes        *string      `json:"notes,omitempty"`
}

// ===== Response DTOs =====

// PrescriptionExerciseResponse represents a single exercise in the response
type PrescriptionExerciseResponse struct {
	ID            uuid.UUID      `json:"id"`
	ExerciseID    uuid.UUID      `json:"exercise_id"`
	ExerciseOrder int            `json:"exercise_order"`
	Sets          *int           `json:"sets,omitempty"`
	Reps          *int           `json:"reps,omitempty"`
	HoldSeconds   *int           `json:"hold_seconds,omitempty"`
	TargetWeight  *WeightOutput  `json:"target_weight,omitempty"`
	RPEValueID    *uuid.UUID     `json:"rpe_value_id,omitempty"`
	Notes         *string        `json:"notes,omitempty"`
	Exercise      *ExerciseBrief `json:"exercise,omitempty"`
	RPEValue      *RPEValueBrief `json:"rpe_value,omitempty"`
}

// PrescriptionGroupResponse represents a group of prescriptions in the response
type PrescriptionGroupResponse struct {
	GroupID         uuid.UUID                      `json:"group_id"`
	Type            PrescriptionType               `json:"type"`
	GroupOrder      int                            `json:"group_order"`
	GroupRounds     *int                           `json:"group_rounds,omitempty"`
	RestBetweenSets *int                           `json:"rest_between_sets,omitempty"`
	GroupName       *string                        `json:"group_name,omitempty"`
	GroupNotes      *string                        `json:"group_notes,omitempty"`
	Exercises       []PrescriptionExerciseResponse `json:"exercises"`
}

// ExerciseBrief is a brief representation of an exercise for responses
type ExerciseBrief struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
}

// RPEValueBrief is a brief representation of an RPE value for responses
type RPEValueBrief struct {
	ID          uuid.UUID `json:"id"`
	Value       int       `json:"value"`
	Label       string    `json:"label"`
	Description string    `json:"description,omitempty"`
}

// ===== Helper Functions =====

// GroupPrescriptionsByGroupID groups prescriptions by their GroupID for response formatting
func GroupPrescriptionsByGroupID(prescriptions []WorkoutPrescription) []PrescriptionGroupResponse {
	if len(prescriptions) == 0 {
		return []PrescriptionGroupResponse{}
	}

	groupMap := make(map[uuid.UUID]*PrescriptionGroupResponse)
	groupOrder := make([]uuid.UUID, 0)

	for _, p := range prescriptions {
		if _, exists := groupMap[p.GroupID]; !exists {
			groupMap[p.GroupID] = &PrescriptionGroupResponse{
				GroupID:         p.GroupID,
				Type:            p.Type,
				GroupOrder:      p.GroupOrder,
				GroupRounds:     p.GroupRounds,
				RestBetweenSets: p.RestBetweenSets,
				GroupName:       p.GroupName,
				GroupNotes:      p.GroupNotes,
				Exercises:       []PrescriptionExerciseResponse{},
			}
			groupOrder = append(groupOrder, p.GroupID)
		}

		exerciseResp := PrescriptionExerciseResponse{
			ID:            p.ID,
			ExerciseID:    p.ExerciseID,
			ExerciseOrder: p.ExerciseOrder,
			Sets:          p.Sets,
			Reps:          p.Reps,
			HoldSeconds:   p.HoldSeconds,
			TargetWeight:  nil, // Will be populated by controller with user's preferred unit
			RPEValueID:    p.RPEValueID,
			Notes:         p.Notes,
		}

		// Add exercise brief if loaded
		if p.Exercise.ID != uuid.Nil {
			exerciseResp.Exercise = &ExerciseBrief{
				ID:          p.Exercise.ID,
				Name:        p.Exercise.Name,
				Slug:        p.Exercise.Slug,
				Description: p.Exercise.Description,
			}
		}

		// Add RPE value brief if loaded
		if p.RPEValue != nil {
			exerciseResp.RPEValue = &RPEValueBrief{
				ID:          p.RPEValue.ID,
				Value:       p.RPEValue.Value,
				Label:       p.RPEValue.Label,
				Description: p.RPEValue.Description,
			}
		}

		groupMap[p.GroupID].Exercises = append(groupMap[p.GroupID].Exercises, exerciseResp)
	}

	// Build result in order
	result := make([]PrescriptionGroupResponse, 0, len(groupOrder))
	for _, gid := range groupOrder {
		result = append(result, *groupMap[gid])
	}

	return result
}

// GetNextGroupOrder returns the next available group order for a workout
func GetNextGroupOrder(db *gorm.DB, workoutID uuid.UUID) (int, error) {
	var maxOrder int
	err := db.Model(&WorkoutPrescription{}).
		Where("workout_id = ?", workoutID).
		Select("COALESCE(MAX(group_order), 0)").
		Scan(&maxOrder).Error
	if err != nil {
		return 0, err
	}
	return maxOrder + 1, nil
}

// GetNextExerciseOrder returns the next available exercise order for a group
func GetNextExerciseOrder(db *gorm.DB, groupID uuid.UUID) (int, error) {
	var maxOrder int
	err := db.Model(&WorkoutPrescription{}).
		Where("group_id = ?", groupID).
		Select("COALESCE(MAX(exercise_order), 0)").
		Scan(&maxOrder).Error
	if err != nil {
		return 0, err
	}
	return maxOrder + 1, nil
}

// ValidateGroupOrderContinuity checks that group orders are continuous (1, 2, 3, ...)
func ValidateGroupOrderContinuity(db *gorm.DB, workoutID uuid.UUID) error {
	var orders []int
	err := db.Model(&WorkoutPrescription{}).
		Where("workout_id = ?", workoutID).
		Distinct("group_order").
		Order("group_order ASC").
		Pluck("group_order", &orders).Error
	if err != nil {
		return err
	}

	for i, order := range orders {
		if order != i+1 {
			return errors.New("group_order values must be continuous starting from 1")
		}
	}
	return nil
}

// ValidateExerciseOrderContinuity checks that exercise orders within a group are continuous
func ValidateExerciseOrderContinuity(db *gorm.DB, groupID uuid.UUID) error {
	var orders []int
	err := db.Model(&WorkoutPrescription{}).
		Where("group_id = ?", groupID).
		Order("exercise_order ASC").
		Pluck("exercise_order", &orders).Error
	if err != nil {
		return err
	}

	for i, order := range orders {
		if order != i+1 {
			return errors.New("exercise_order values must be continuous starting from 1")
		}
	}
	return nil
}
