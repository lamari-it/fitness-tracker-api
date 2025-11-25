package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ===== WORKOUT SESSION =====

// WorkoutSession represents a single workout execution by a user
type WorkoutSession struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID             uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	CreatedByID        *uuid.UUID     `gorm:"type:uuid" json:"created_by_id,omitempty"` // Who logged the session (trainer or self)
	WorkoutID          *uuid.UUID     `gorm:"type:uuid" json:"workout_id"`              // nullable for free-form workouts
	StartedAt          time.Time      `gorm:"not null" json:"started_at"`
	EndedAt            *time.Time     `json:"ended_at"`
	DurationSeconds    *int           `json:"duration_seconds,omitempty"`
	PerceivedIntensity *int           `gorm:"check:perceived_intensity >= 1 AND perceived_intensity <= 10" json:"perceived_intensity,omitempty"`
	Notes              string         `gorm:"type:text" json:"notes"`
	Completed          bool           `gorm:"default:false" json:"completed"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	User          User           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	CreatedBy     *User          `gorm:"foreignKey:CreatedByID;constraint:OnDelete:SET NULL" json:"created_by,omitempty"`
	Workout       *Workout       `gorm:"foreignKey:WorkoutID;constraint:OnDelete:SET NULL" json:"workout,omitempty"`
	SessionBlocks []SessionBlock `gorm:"foreignKey:SessionID" json:"session_blocks,omitempty"`
}

func (ws *WorkoutSession) BeforeCreate(tx *gorm.DB) (err error) {
	if ws.ID == uuid.Nil {
		ws.ID = uuid.New()
	}
	return
}

// ===== SESSION BLOCK =====

// SessionBlock represents a block of exercises in a session (mirrors prescription group)
type SessionBlock struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SessionID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"session_id"`
	GroupID           uuid.UUID      `gorm:"type:uuid;not null;index" json:"group_id"` // From workout_prescriptions
	BlockOrder        int            `gorm:"not null" json:"block_order"`
	StartedAt         *time.Time     `json:"started_at,omitempty"`
	CompletedAt       *time.Time     `json:"completed_at,omitempty"`
	Skipped           bool           `gorm:"default:false" json:"skipped"`
	PerceivedExertion *int           `gorm:"check:perceived_exertion >= 1 AND perceived_exertion <= 10" json:"perceived_exertion,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Session          WorkoutSession    `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE" json:"session,omitempty"`
	SessionExercises []SessionExercise `gorm:"foreignKey:SessionBlockID" json:"session_exercises,omitempty"`
}

func (sb *SessionBlock) BeforeCreate(tx *gorm.DB) (err error) {
	if sb.ID == uuid.Nil {
		sb.ID = uuid.New()
	}
	return
}

// ===== SESSION EXERCISE =====

// SessionExercise represents an exercise instance within a session block
type SessionExercise struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SessionBlockID uuid.UUID      `gorm:"type:uuid;not null;index" json:"session_block_id"`
	PrescriptionID *uuid.UUID     `gorm:"type:uuid;index" json:"prescription_id,omitempty"` // Reference to workout_prescriptions
	ExerciseID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"exercise_id"`
	ExerciseOrder  int            `gorm:"not null" json:"exercise_order"`
	StartedAt      *time.Time     `json:"started_at,omitempty"`
	CompletedAt    *time.Time     `json:"completed_at,omitempty"`
	Skipped        bool           `gorm:"default:false" json:"skipped"`
	Notes          string         `gorm:"type:text" json:"notes"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	SessionBlock SessionBlock         `gorm:"foreignKey:SessionBlockID;constraint:OnDelete:CASCADE" json:"session_block,omitempty"`
	Prescription *WorkoutPrescription `gorm:"foreignKey:PrescriptionID;constraint:OnDelete:SET NULL" json:"prescription,omitempty"`
	Exercise     Exercise             `gorm:"foreignKey:ExerciseID;constraint:OnDelete:CASCADE" json:"exercise,omitempty"`
	SessionSets  []SessionSet         `gorm:"foreignKey:SessionExerciseID" json:"session_sets,omitempty"`
}

func (se *SessionExercise) BeforeCreate(tx *gorm.DB) (err error) {
	if se.ID == uuid.Nil {
		se.ID = uuid.New()
	}
	return
}

// ===== SESSION SET =====

// SessionSet represents an actual performed set within a session exercise
type SessionSet struct {
	ID                         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SessionExerciseID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"session_exercise_id"`
	SetNumber                  int            `gorm:"not null" json:"set_number"`
	Completed                  bool           `gorm:"default:false" json:"completed"`
	ActualReps                 *int           `json:"actual_reps,omitempty"`
	ActualWeightKg             *float64       `gorm:"type:decimal(6,2)" json:"-"`
	OriginalActualWeightValue  *float64       `gorm:"type:decimal(6,2)" json:"-"`
	OriginalActualWeightUnit   *string        `gorm:"type:varchar(2)" json:"-"`
	ActualDurationSeconds      *int           `json:"actual_duration_seconds,omitempty"`
	RPEValueID                 *uuid.UUID     `gorm:"type:uuid" json:"rpe_value_id,omitempty"`
	WasFailure                 bool           `gorm:"default:false" json:"was_failure"`
	Notes                      string         `gorm:"type:text" json:"notes"`
	CreatedAt                  time.Time      `json:"created_at"`
	UpdatedAt                  time.Time      `json:"updated_at"`
	DeletedAt                  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	SessionExercise SessionExercise `gorm:"foreignKey:SessionExerciseID;constraint:OnDelete:CASCADE" json:"session_exercise,omitempty"`
	RPEValue        *RPEScaleValue  `gorm:"foreignKey:RPEValueID;constraint:OnDelete:SET NULL" json:"rpe_value,omitempty"`
}

func (ss *SessionSet) BeforeCreate(tx *gorm.DB) (err error) {
	if ss.ID == uuid.Nil {
		ss.ID = uuid.New()
	}
	return
}

// ===== REQUEST DTOs =====

// CreateWorkoutSessionRequest represents the request to start a workout session
type CreateWorkoutSessionRequest struct {
	UserID    *uuid.UUID `json:"user_id,omitempty"`    // Optional: if not provided, uses authenticated user
	WorkoutID *uuid.UUID `json:"workout_id,omitempty"` // Optional: for free-form workouts
	StartedAt *time.Time `json:"started_at,omitempty"` // Optional: defaults to now
	Notes     *string    `json:"notes,omitempty"`
}

// EndWorkoutSessionRequest represents the request to end a workout session
type EndWorkoutSessionRequest struct {
	EndedAt            *time.Time `json:"ended_at,omitempty"` // Optional: defaults to now
	Notes              *string    `json:"notes,omitempty"`
	PerceivedIntensity *int       `json:"perceived_intensity,omitempty" binding:"omitempty,min=1,max=10"`
}

// UpdateWorkoutSessionRequest represents the request to update a session
type UpdateWorkoutSessionRequest struct {
	Notes              *string `json:"notes,omitempty"`
	PerceivedIntensity *int    `json:"perceived_intensity,omitempty" binding:"omitempty,min=1,max=10"`
}

// CompleteWorkoutSessionRequest represents the request to complete a workout session
type CompleteWorkoutSessionRequest struct {
	PerceivedIntensity *int   `json:"perceived_intensity,omitempty" binding:"omitempty,min=1,max=10"`
	Notes              string `json:"notes,omitempty"`
}

// UpdateSessionBlockRequest represents the request to update a session block
type UpdateSessionBlockRequest struct {
	PerceivedExertion *int `json:"perceived_exertion,omitempty" binding:"omitempty,min=1,max=10"`
}

// UpdateSessionExerciseRequest represents the request to update a session exercise
type UpdateSessionExerciseRequest struct {
	Notes string `json:"notes,omitempty"`
}

// CreateSessionSetRequest represents the request to add a set to an exercise
type CreateSessionSetRequest struct {
	ActualReps            *int         `json:"actual_reps,omitempty"`
	ActualWeight          *WeightInput `json:"actual_weight,omitempty"`
	ActualDurationSeconds *int         `json:"actual_duration_seconds,omitempty"`
	RPEValueID            *uuid.UUID   `json:"rpe_value_id,omitempty"`
	Notes                 *string      `json:"notes,omitempty"`
}

// UpdateSessionSetRequest represents the request to update a logged set
type UpdateSessionSetRequest struct {
	ActualReps            *int         `json:"actual_reps,omitempty"`
	ActualWeight          *WeightInput `json:"actual_weight,omitempty"`
	ActualDurationSeconds *int         `json:"actual_duration_seconds,omitempty"`
	RPEValueID            *uuid.UUID   `json:"rpe_value_id,omitempty"`
	WasFailure            *bool        `json:"was_failure,omitempty"`
	Completed             *bool        `json:"completed,omitempty"`
	Notes                 *string      `json:"notes,omitempty"`
}

// ===== RESPONSE DTOs =====

// SessionSetResponse represents a set in the response
type SessionSetResponse struct {
	ID                    uuid.UUID      `json:"id"`
	SetNumber             int            `json:"set_number"`
	Completed             bool           `json:"completed"`
	ActualReps            *int           `json:"actual_reps,omitempty"`
	ActualWeight          *WeightOutput  `json:"actual_weight,omitempty"`
	ActualDurationSeconds *int           `json:"actual_duration_seconds,omitempty"`
	RPEValueID            *uuid.UUID     `json:"rpe_value_id,omitempty"`
	RPEValue              *RPEValueBrief `json:"rpe_value,omitempty"`
	WasFailure            bool           `json:"was_failure"`
	Notes                 string         `json:"notes,omitempty"`
	CreatedAt             time.Time      `json:"created_at"`
}

// SessionExerciseResponse represents an exercise in the response
type SessionExerciseResponse struct {
	ID             uuid.UUID            `json:"id"`
	PrescriptionID *uuid.UUID           `json:"prescription_id,omitempty"`
	ExerciseID     uuid.UUID            `json:"exercise_id"`
	ExerciseName   string               `json:"exercise_name"`
	ExerciseOrder  int                  `json:"exercise_order"`
	StartedAt      *time.Time           `json:"started_at,omitempty"`
	CompletedAt    *time.Time           `json:"completed_at,omitempty"`
	Skipped        bool                 `json:"skipped"`
	Notes          string               `json:"notes,omitempty"`
	Prescription   *PrescriptionBrief   `json:"prescription,omitempty"`
	Sets           []SessionSetResponse `json:"sets"`
}

// PrescriptionBrief represents prescription details in responses
type PrescriptionBrief struct {
	Sets         *int          `json:"sets,omitempty"`
	Reps         *int          `json:"reps,omitempty"`
	HoldSeconds  *int          `json:"hold_seconds,omitempty"`
	TargetWeight *WeightOutput `json:"target_weight,omitempty"`
}

// SessionBlockResponse represents a block in the response
type SessionBlockResponse struct {
	ID                uuid.UUID                 `json:"id"`
	GroupID           uuid.UUID                 `json:"group_id"`
	BlockOrder        int                       `json:"block_order"`
	Type              PrescriptionType          `json:"type,omitempty"`
	GroupName         string                    `json:"group_name,omitempty"`
	GroupRounds       *int                      `json:"group_rounds,omitempty"`
	RestBetweenSets   *int                      `json:"rest_between_sets,omitempty"`
	StartedAt         *time.Time                `json:"started_at"`
	CompletedAt       *time.Time                `json:"completed_at"`
	Skipped           bool                      `json:"skipped"`
	PerceivedExertion *int                      `json:"perceived_exertion,omitempty"`
	Exercises         []SessionExerciseResponse `json:"exercises"`
}

// WorkoutSessionResponse represents the full session response
type WorkoutSessionResponse struct {
	ID                 uuid.UUID              `json:"id"`
	UserID             uuid.UUID              `json:"user_id"`
	CreatedByID        *uuid.UUID             `json:"created_by_id,omitempty"`
	CreatedByName      string                 `json:"created_by_name,omitempty"`
	WorkoutID          *uuid.UUID             `json:"workout_id,omitempty"`
	WorkoutTitle       string                 `json:"workout_title,omitempty"`
	StartedAt          time.Time              `json:"started_at"`
	EndedAt            *time.Time             `json:"ended_at,omitempty"`
	DurationSeconds    *int                   `json:"duration_seconds,omitempty"`
	DurationMinutes    *int                   `json:"duration_minutes,omitempty"`
	PerceivedIntensity *int                   `json:"perceived_intensity,omitempty"`
	Notes              string                 `json:"notes,omitempty"`
	Completed          bool                   `json:"completed"`
	Blocks             []SessionBlockResponse `json:"blocks"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
}

// WorkoutSessionListResponse represents a session in list view (without nested details)
type WorkoutSessionListResponse struct {
	ID                 uuid.UUID  `json:"id"`
	UserID             uuid.UUID  `json:"user_id"`
	WorkoutID          *uuid.UUID `json:"workout_id,omitempty"`
	WorkoutTitle       string     `json:"workout_title,omitempty"`
	StartedAt          time.Time  `json:"started_at"`
	EndedAt            *time.Time `json:"ended_at,omitempty"`
	DurationMinutes    *int       `json:"duration_minutes,omitempty"`
	PerceivedIntensity *int       `json:"perceived_intensity,omitempty"`
	Completed          bool       `json:"completed"`
	TotalBlocks        int        `json:"total_blocks"`
	CompletedBlocks    int        `json:"completed_blocks"`
	CreatedAt          time.Time  `json:"created_at"`
}

// ===== HELPER FUNCTIONS =====

// ToListResponse converts a WorkoutSession to a list response
func (ws *WorkoutSession) ToListResponse() WorkoutSessionListResponse {
	response := WorkoutSessionListResponse{
		ID:                 ws.ID,
		UserID:             ws.UserID,
		WorkoutID:          ws.WorkoutID,
		StartedAt:          ws.StartedAt,
		EndedAt:            ws.EndedAt,
		PerceivedIntensity: ws.PerceivedIntensity,
		Completed:          ws.Completed,
		CreatedAt:          ws.CreatedAt,
	}

	if ws.Workout != nil {
		response.WorkoutTitle = ws.Workout.Title
	}

	if ws.EndedAt != nil {
		duration := int(ws.EndedAt.Sub(ws.StartedAt).Minutes())
		response.DurationMinutes = &duration
	}

	response.TotalBlocks = len(ws.SessionBlocks)
	for _, block := range ws.SessionBlocks {
		if block.CompletedAt != nil || block.Skipped {
			response.CompletedBlocks++
		}
	}

	return response
}

// BuildSessionResponse builds a full session response with nested blocks, exercises, and sets
func BuildSessionResponse(session WorkoutSession, weightUnit string) WorkoutSessionResponse {
	response := WorkoutSessionResponse{
		ID:                 session.ID,
		UserID:             session.UserID,
		CreatedByID:        session.CreatedByID,
		WorkoutID:          session.WorkoutID,
		StartedAt:          session.StartedAt,
		EndedAt:            session.EndedAt,
		DurationSeconds:    session.DurationSeconds,
		PerceivedIntensity: session.PerceivedIntensity,
		Notes:              session.Notes,
		Completed:          session.Completed,
		Blocks:             []SessionBlockResponse{},
		CreatedAt:          session.CreatedAt,
		UpdatedAt:          session.UpdatedAt,
	}

	// Include creator's name if available
	if session.CreatedBy != nil {
		response.CreatedByName = session.CreatedBy.FirstName + " " + session.CreatedBy.LastName
	}

	// Include workout title if available
	if session.Workout != nil {
		response.WorkoutTitle = session.Workout.Title
	}

	// Calculate duration in minutes
	if session.EndedAt != nil {
		duration := int(session.EndedAt.Sub(session.StartedAt).Minutes())
		response.DurationMinutes = &duration
	}

	// Build nested block responses
	for _, block := range session.SessionBlocks {
		blockResp := SessionBlockResponse{
			ID:                block.ID,
			GroupID:           block.GroupID,
			BlockOrder:        block.BlockOrder,
			StartedAt:         block.StartedAt,
			CompletedAt:       block.CompletedAt,
			Skipped:           block.Skipped,
			PerceivedExertion: block.PerceivedExertion,
			Exercises:         []SessionExerciseResponse{},
		}

		// Get group metadata from first prescription in block (if available)
		if len(block.SessionExercises) > 0 && block.SessionExercises[0].Prescription != nil {
			p := block.SessionExercises[0].Prescription
			blockResp.Type = p.Type
			if p.GroupName != nil {
				blockResp.GroupName = *p.GroupName
			}
			blockResp.GroupRounds = p.GroupRounds
			blockResp.RestBetweenSets = p.RestBetweenSets
		}

		// Build nested exercise responses
		for _, exercise := range block.SessionExercises {
			exerciseResp := SessionExerciseResponse{
				ID:             exercise.ID,
				PrescriptionID: exercise.PrescriptionID,
				ExerciseID:     exercise.ExerciseID,
				ExerciseName:   exercise.Exercise.Name,
				ExerciseOrder:  exercise.ExerciseOrder,
				StartedAt:      exercise.StartedAt,
				CompletedAt:    exercise.CompletedAt,
				Skipped:        exercise.Skipped,
				Notes:          exercise.Notes,
				Sets:           []SessionSetResponse{},
			}

			// Include prescription details if available
			if exercise.Prescription != nil {
				exerciseResp.Prescription = &PrescriptionBrief{
					Sets:        exercise.Prescription.Sets,
					Reps:        exercise.Prescription.Reps,
					HoldSeconds: exercise.Prescription.HoldSeconds,
				}
				// Convert target weight to user's preferred unit
				if exercise.Prescription.TargetWeightKg != nil {
					convertedValue := *exercise.Prescription.TargetWeightKg
					if weightUnit == "lb" {
						convertedValue = convertedValue * 2.20462262185
					}
					unit := weightUnit
					exerciseResp.Prescription.TargetWeight = &WeightOutput{
						WeightValue: &convertedValue,
						WeightUnit:  &unit,
					}
				}
			}

			// Build nested set responses
			for _, set := range exercise.SessionSets {
				setResp := SessionSetResponse{
					ID:                    set.ID,
					SetNumber:             set.SetNumber,
					Completed:             set.Completed,
					ActualReps:            set.ActualReps,
					ActualDurationSeconds: set.ActualDurationSeconds,
					RPEValueID:            set.RPEValueID,
					WasFailure:            set.WasFailure,
					Notes:                 set.Notes,
					CreatedAt:             set.CreatedAt,
				}

				// Convert weight to user's preferred unit
				if set.ActualWeightKg != nil {
					convertedValue := *set.ActualWeightKg
					if weightUnit == "lb" {
						convertedValue = convertedValue * 2.20462262185
					}
					unit := weightUnit
					setResp.ActualWeight = &WeightOutput{
						WeightValue: &convertedValue,
						WeightUnit:  &unit,
					}
				}

				// Include RPE value details if available
				if set.RPEValue != nil {
					setResp.RPEValue = &RPEValueBrief{
						ID:          set.RPEValue.ID,
						Value:       set.RPEValue.Value,
						Label:       set.RPEValue.Label,
						Description: set.RPEValue.Description,
					}
				}

				exerciseResp.Sets = append(exerciseResp.Sets, setResp)
			}

			blockResp.Exercises = append(blockResp.Exercises, exerciseResp)
		}

		response.Blocks = append(response.Blocks, blockResp)
	}

	return response
}

// ToResponse converts a SessionBlock to a response with weight converted to user's preferred unit
func (sb *SessionBlock) ToResponse(preferredWeightUnit string) SessionBlockResponse {
	blockResp := SessionBlockResponse{
		ID:                sb.ID,
		GroupID:           sb.GroupID,
		BlockOrder:        sb.BlockOrder,
		StartedAt:         sb.StartedAt,
		CompletedAt:       sb.CompletedAt,
		Skipped:           sb.Skipped,
		PerceivedExertion: sb.PerceivedExertion,
		Exercises:         []SessionExerciseResponse{},
	}

	// Get group metadata from first prescription in block (if available)
	if len(sb.SessionExercises) > 0 && sb.SessionExercises[0].Prescription != nil {
		p := sb.SessionExercises[0].Prescription
		blockResp.Type = p.Type
		if p.GroupName != nil {
			blockResp.GroupName = *p.GroupName
		}
		blockResp.GroupRounds = p.GroupRounds
		blockResp.RestBetweenSets = p.RestBetweenSets
	}

	// Build nested exercise responses
	for _, exercise := range sb.SessionExercises {
		blockResp.Exercises = append(blockResp.Exercises, exercise.ToResponse(preferredWeightUnit))
	}

	return blockResp
}

// ToResponse converts a SessionExercise to a response with weight converted to user's preferred unit
func (se *SessionExercise) ToResponse(preferredWeightUnit string) SessionExerciseResponse {
	exerciseResp := SessionExerciseResponse{
		ID:             se.ID,
		PrescriptionID: se.PrescriptionID,
		ExerciseID:     se.ExerciseID,
		ExerciseName:   se.Exercise.Name,
		ExerciseOrder:  se.ExerciseOrder,
		StartedAt:      se.StartedAt,
		CompletedAt:    se.CompletedAt,
		Skipped:        se.Skipped,
		Notes:          se.Notes,
		Sets:           []SessionSetResponse{},
	}

	// Include prescription details if available
	if se.Prescription != nil {
		exerciseResp.Prescription = &PrescriptionBrief{
			Sets:        se.Prescription.Sets,
			Reps:        se.Prescription.Reps,
			HoldSeconds: se.Prescription.HoldSeconds,
		}
		// Convert target weight to user's preferred unit
		if se.Prescription.TargetWeightKg != nil {
			convertedValue := *se.Prescription.TargetWeightKg
			if preferredWeightUnit == "lb" {
				convertedValue = convertedValue * 2.20462262185
			}
			unit := preferredWeightUnit
			exerciseResp.Prescription.TargetWeight = &WeightOutput{
				WeightValue: &convertedValue,
				WeightUnit:  &unit,
			}
		}
	}

	// Build nested set responses
	for _, set := range se.SessionSets {
		exerciseResp.Sets = append(exerciseResp.Sets, set.ToResponse(preferredWeightUnit))
	}

	return exerciseResp
}

// ToResponse converts a SessionSet to a response (with default kg weight unit)
func (ss *SessionSet) ToResponse(preferredWeightUnit string) SessionSetResponse {
	setResp := SessionSetResponse{
		ID:                    ss.ID,
		SetNumber:             ss.SetNumber,
		Completed:             ss.Completed,
		ActualReps:            ss.ActualReps,
		ActualDurationSeconds: ss.ActualDurationSeconds,
		RPEValueID:            ss.RPEValueID,
		WasFailure:            ss.WasFailure,
		Notes:                 ss.Notes,
		CreatedAt:             ss.CreatedAt,
	}

	// Convert weight to user's preferred unit
	if ss.ActualWeightKg != nil {
		convertedValue := *ss.ActualWeightKg
		if preferredWeightUnit == "lb" {
			convertedValue = convertedValue * 2.20462262185
		}
		unit := preferredWeightUnit
		if unit == "" {
			unit = "kg" // default
		}
		setResp.ActualWeight = &WeightOutput{
			WeightValue: &convertedValue,
			WeightUnit:  &unit,
		}
	}

	// Include RPE value details if available
	if ss.RPEValue != nil {
		setResp.RPEValue = &RPEValueBrief{
			ID:          ss.RPEValue.ID,
			Value:       ss.RPEValue.Value,
			Label:       ss.RPEValue.Label,
			Description: ss.RPEValue.Description,
		}
	}

	return setResp
}
