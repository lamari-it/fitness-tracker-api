package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type WorkoutPlan struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Title         string    `gorm:"type:text;not null" json:"title"`
	Description   string    `gorm:"type:text" json:"description"`
	Visibility    string    `gorm:"type:varchar(20);default:'private'" json:"visibility"` // private, public, friends
	TemplateWeeks int       `gorm:"not null;default:1" json:"template_weeks"`                 // allows multi-week blocks

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Items []WorkoutPlanItem `gorm:"foreignKey:PlanID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

type Workout struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Title       string    `gorm:"type:text;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Visibility  string    `gorm:"type:varchar(20);default:'private'" json:"visibility"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Exercises inside workout define their own order (not handled by plan)
	// Exercises []Exercise ...
	SetGroups        []SetGroup        `gorm:"foreignKey:WorkoutID" json:"set_groups,omitempty"`
	WorkoutExercises []WorkoutExercise `gorm:"foreignKey:WorkoutID" json:"exercises,omitempty"`
	WorkoutSessions  []WorkoutSession  `gorm:"foreignKey:WorkoutID" json:"sessions,omitempty"`
	SharedWorkouts   []SharedWorkout   `gorm:"foreignKey:WorkoutID" json:"shared_workouts,omitempty"`
}

type WorkoutPlanItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PlanID    uuid.UUID `gorm:"type:uuid;not null;index" json:"plan_id"`
	WorkoutID uuid.UUID `gorm:"type:uuid;not null;index" json:"workout_id"`

	WeekIndex int `gorm:"not null;default:0" json:"week_index"` // optional, for multi-week blocks

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Workout Workout `gorm:"constraint:OnDelete:CASCADE" json:"workout,omitempty"`
}

type PlanEnrollment struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PlanID      uuid.UUID `gorm:"type:uuid;not null;index" json:"plan_id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`

	StartDate   time.Time `gorm:"not null" json:"start_date"`
	DaysPerWeek int       `gorm:"not null;check:days_per_week_check,days_per_week >= 1 AND days_per_week <= 7" json:"days_per_week"`

	CurrentIndex int `gorm:"not null;default:0" json:"current_index"` // index in rolling mode

	ScheduleMode      string         `gorm:"type:varchar(20);default:'rolling'" json:"schedule_mode"` // rolling | calendar
	PreferredWeekdays pq.Int32Array  `gorm:"type:int[];default:'{}'" json:"preferred_weekdays"`       // only used in calendar mode (0=Mon..6=Sun)

	Status    string    `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Exercise struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Slug         string    `gorm:"type:varchar(255);not null;unique" json:"slug"`
	Name         string    `gorm:"type:text;not null;unique" json:"name"`
	Description  string    `gorm:"type:text" json:"description"`
	IsBodyweight bool      `gorm:"default:false" json:"is_bodyweight"`
	Instructions string    `gorm:"type:text" json:"instructions"`
	VideoURL     string    `gorm:"type:text" json:"video_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	WorkoutExercises []WorkoutExercise     `gorm:"foreignKey:ExerciseID" json:"workout_exercises,omitempty"`
	ExerciseLogs     []ExerciseLog         `gorm:"foreignKey:ExerciseID" json:"exercise_logs,omitempty"`
	MuscleGroups     []ExerciseMuscleGroup `gorm:"foreignKey:ExerciseID" json:"muscle_groups,omitempty"`
	Equipment        []ExerciseEquipment   `gorm:"foreignKey:ExerciseID" json:"equipment,omitempty"`
}

type WorkoutExercise struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	WorkoutID         uuid.UUID `gorm:"type:uuid;not null" json:"workout_id"`
	SetGroupID        uuid.UUID `gorm:"type:uuid;not null" json:"set_group_id"`
	ExerciseID        uuid.UUID `gorm:"type:uuid;not null" json:"exercise_id"`
	OrderNumber       int       `gorm:"not null" json:"order_number"`
	TargetSets        int       `json:"target_sets"`
	TargetReps        int       `json:"target_reps"`
	TargetWeight      float64   `gorm:"type:numeric(10,2)" json:"target_weight"`
	TargetRestSec     int       `json:"target_rest_sec"`
	Prescription      string    `gorm:"type:varchar(20);not null;default:'reps'" json:"prescription"` // reps | time
	TargetDurationSec int       `gorm:"default:0" json:"target_duration_sec"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Workout  Workout  `gorm:"foreignKey:WorkoutID;constraint:OnDelete:CASCADE" json:"workout,omitempty"`
	SetGroup SetGroup `gorm:"foreignKey:SetGroupID;constraint:OnDelete:CASCADE" json:"set_group,omitempty"`
	Exercise Exercise `gorm:"foreignKey:ExerciseID;constraint:OnDelete:CASCADE" json:"exercise,omitempty"`
}

func (wp *WorkoutPlan) BeforeCreate(tx *gorm.DB) (err error) {
	if wp.ID == uuid.Nil {
		wp.ID = uuid.New()
	}
	return
}

func (w *Workout) BeforeCreate(tx *gorm.DB) (err error) {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return
}

func (e *Exercise) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return
}

func (we *WorkoutExercise) BeforeCreate(tx *gorm.DB) (err error) {
	if we.ID == uuid.Nil {
		we.ID = uuid.New()
	}
	return
}

// Optional but recommended: keep reps vs time mutually exclusive and valid
func (we *WorkoutExercise) BeforeSave(tx *gorm.DB) (err error) {
	if we.Prescription == "" {
		we.Prescription = "reps"
	}
	switch we.Prescription {
	case "reps":
		if we.TargetReps <= 0 {
			return errors.New("target_reps must be > 0 when prescription = 'reps'")
		}
		we.TargetDurationSec = 0
	case "time":
		if we.TargetDurationSec <= 0 {
			return errors.New("target_duration_sec must be > 0 when prescription = 'time'")
		}
		we.TargetReps = 0
	default:
		return errors.New("invalid prescription: must be 'reps' or 'time'")
	}
	return nil
}
