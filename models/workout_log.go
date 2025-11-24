package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkoutSession struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	CreatedByID *uuid.UUID     `gorm:"type:uuid" json:"created_by_id,omitempty"` // Who logged the session (trainer or self)
	WorkoutID   *uuid.UUID     `gorm:"type:uuid" json:"workout_id"`              // nullable for free-form workouts
	StartedAt   time.Time      `gorm:"not null" json:"started_at"`
	EndedAt     *time.Time     `json:"ended_at"`
	Notes       string         `gorm:"type:text" json:"notes"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	User         User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	CreatedBy    *User         `gorm:"foreignKey:CreatedByID;constraint:OnDelete:SET NULL" json:"created_by,omitempty"`
	Workout      *Workout      `gorm:"foreignKey:WorkoutID;constraint:OnDelete:SET NULL" json:"workout,omitempty"`
	ExerciseLogs []ExerciseLog `gorm:"foreignKey:SessionID" json:"exercise_logs,omitempty"`
}

type ExerciseLog struct {
	ID               uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SessionID        uuid.UUID      `gorm:"type:uuid;not null" json:"session_id"`
	PrescriptionID   *uuid.UUID     `gorm:"type:uuid" json:"prescription_id"` // Reference to the prescription this log is for
	ExerciseID       uuid.UUID      `gorm:"type:uuid;not null" json:"exercise_id"`
	OrderNumber      int            `gorm:"not null" json:"order_number"`
	Notes            string         `gorm:"type:text" json:"notes"`
	DifficultyRating int            `gorm:"check:difficulty_rating >= 0 AND difficulty_rating <= 10" json:"difficulty_rating"`
	DifficultyType   string         `gorm:"type:varchar(20)" json:"difficulty_type"` // easy, moderate, hard, very_hard
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Session      WorkoutSession       `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE" json:"session,omitempty"`
	Prescription *WorkoutPrescription `gorm:"foreignKey:PrescriptionID;constraint:OnDelete:SET NULL" json:"prescription,omitempty"`
	Exercise     Exercise             `gorm:"foreignKey:ExerciseID;constraint:OnDelete:CASCADE" json:"exercise,omitempty"`
	SetLogs      []SetLog             `gorm:"foreignKey:ExerciseLogID" json:"set_logs,omitempty"`
}

type SetLog struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ExerciseLogID   uuid.UUID      `gorm:"type:uuid;not null" json:"exercise_log_id"`
	SetNumber       int            `gorm:"not null" json:"set_number"`
	Weight          float64        `gorm:"type:numeric(10,2)" json:"weight"`                      // Canonical weight in kg
	WeightUnit      string         `gorm:"type:varchar(5);default:'kg'" json:"weight_unit"`       // Display unit (deprecated, use InputWeightUnit)
	InputWeight     float64        `gorm:"type:numeric(10,2)" json:"input_weight"`                // Original weight value entered by user
	InputWeightUnit string         `gorm:"type:varchar(5);default:'kg'" json:"input_weight_unit"` // Unit of original input (kg/lb)
	Reps            int            `json:"reps"`
	RestAfterSec    int            `json:"rest_after_sec"`
	Tempo           string         `gorm:"type:varchar(10)" json:"tempo"`                             // e.g., "3-1-2-1"
	RPE             float64        `gorm:"type:numeric(3,1);check:rpe >= 0 AND rpe <= 10" json:"rpe"` // Rate of Perceived Exertion (legacy, 0 = not rated)
	RPEValueID      *uuid.UUID     `gorm:"type:uuid" json:"rpe_value_id,omitempty"`                   // Reference to custom RPE scale value
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	ExerciseLog ExerciseLog    `gorm:"foreignKey:ExerciseLogID;constraint:OnDelete:CASCADE" json:"exercise_log,omitempty"`
	RPEValue    *RPEScaleValue `gorm:"foreignKey:RPEValueID;constraint:OnDelete:SET NULL" json:"rpe_value,omitempty"`
}

func (ws *WorkoutSession) BeforeCreate(tx *gorm.DB) (err error) {
	if ws.ID == uuid.Nil {
		ws.ID = uuid.New()
	}
	return
}

func (el *ExerciseLog) BeforeCreate(tx *gorm.DB) (err error) {
	if el.ID == uuid.Nil {
		el.ID = uuid.New()
	}
	return
}

func (sl *SetLog) BeforeCreate(tx *gorm.DB) (err error) {
	if sl.ID == uuid.Nil {
		sl.ID = uuid.New()
	}
	return
}

// Response DTOs
type WorkoutSessionResponse struct {
	ID            uuid.UUID          `json:"id"`
	UserID        uuid.UUID          `json:"user_id"`
	CreatedByID   *uuid.UUID         `json:"created_by_id,omitempty"`
	CreatedByName string             `json:"created_by_name,omitempty"`
	WorkoutID     *uuid.UUID         `json:"workout_id"`
	StartedAt     time.Time          `json:"started_at"`
	EndedAt       *time.Time         `json:"ended_at"`
	Notes         string             `json:"notes"`
	Duration      *int               `json:"duration_minutes,omitempty"`
	ExerciseLogs  []ExerciseLogBrief `json:"exercise_logs,omitempty"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

type ExerciseLogBrief struct {
	ID               uuid.UUID        `json:"id"`
	PrescriptionID   *uuid.UUID       `json:"prescription_id,omitempty"`
	GroupID          *uuid.UUID       `json:"group_id,omitempty"`
	GroupName        string           `json:"group_name,omitempty"`
	GroupType        PrescriptionType `json:"group_type,omitempty"`
	ExerciseID       uuid.UUID        `json:"exercise_id"`
	ExerciseName     string           `json:"exercise_name"`
	OrderNumber      int              `json:"order_number"`
	TotalSets        int              `json:"total_sets"`
	DifficultyRating int              `json:"difficulty_rating"`
}

func (ws *WorkoutSession) ToResponse() WorkoutSessionResponse {
	response := WorkoutSessionResponse{
		ID:          ws.ID,
		UserID:      ws.UserID,
		CreatedByID: ws.CreatedByID,
		WorkoutID:   ws.WorkoutID,
		StartedAt:   ws.StartedAt,
		EndedAt:     ws.EndedAt,
		Notes:       ws.Notes,
		CreatedAt:   ws.CreatedAt,
		UpdatedAt:   ws.UpdatedAt,
	}

	// Include creator's name if available
	if ws.CreatedBy != nil {
		response.CreatedByName = ws.CreatedBy.FirstName + " " + ws.CreatedBy.LastName
	}

	if ws.EndedAt != nil {
		duration := int(ws.EndedAt.Sub(ws.StartedAt).Minutes())
		response.Duration = &duration
	}

	for _, log := range ws.ExerciseLogs {
		brief := ExerciseLogBrief{
			ID:               log.ID,
			PrescriptionID:   log.PrescriptionID,
			ExerciseID:       log.ExerciseID,
			ExerciseName:     log.Exercise.Name,
			OrderNumber:      log.OrderNumber,
			TotalSets:        len(log.SetLogs),
			DifficultyRating: log.DifficultyRating,
		}
		if log.Prescription != nil {
			brief.GroupID = &log.Prescription.GroupID
			if log.Prescription.GroupName != nil {
				brief.GroupName = *log.Prescription.GroupName
			}
			brief.GroupType = log.Prescription.Type
		}
		response.ExerciseLogs = append(response.ExerciseLogs, brief)
	}

	return response
}
