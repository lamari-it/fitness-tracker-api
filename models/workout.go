package models

import (
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
	TemplateWeeks int       `gorm:"not null;default:1" json:"template_weeks"`             // allows multi-week blocks

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Items []WorkoutPlanItem `gorm:"foreignKey:PlanID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

type Workout struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID            uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Title             string    `gorm:"type:text;not null" json:"title"`
	Description       string    `gorm:"type:text" json:"description"`
	DifficultyLevel   string    `gorm:"type:varchar(20)" json:"difficulty_level"` // beginner, intermediate, advanced
	EstimatedDuration *int      `gorm:"type:int" json:"estimated_duration"`       // duration in minutes
	IsTemplate        bool      `gorm:"default:false" json:"is_template"`         // save as template
	Visibility        string    `gorm:"type:varchar(20);default:'private'" json:"visibility"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Unified prescription architecture
	Prescriptions   []WorkoutPrescription `gorm:"foreignKey:WorkoutID" json:"prescriptions,omitempty"`
	WorkoutSessions []WorkoutSession      `gorm:"foreignKey:WorkoutID" json:"sessions,omitempty"`
	SharedWorkouts  []SharedWorkout       `gorm:"foreignKey:WorkoutID" json:"shared_workouts,omitempty"`
}

type WorkoutPlanItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PlanID    uuid.UUID `gorm:"type:uuid;not null;index" json:"plan_id"`
	WorkoutID uuid.UUID `gorm:"type:uuid;not null;index" json:"workout_id"`

	WeekIndex int `gorm:"not null;default:0" json:"week_index"` // optional, for multi-week blocks

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Workout Workout `gorm:"constraint:OnDelete:CASCADE" json:"workout,omitempty"`
}

type PlanEnrollment struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PlanID uuid.UUID `gorm:"type:uuid;not null;index" json:"plan_id"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`

	StartDate   time.Time `gorm:"not null" json:"start_date"`
	DaysPerWeek int       `gorm:"not null;check:days_per_week_check,days_per_week >= 1 AND days_per_week <= 7" json:"days_per_week"`

	CurrentIndex int `gorm:"not null;default:0" json:"current_index"` // index in rolling mode

	ScheduleMode      string        `gorm:"type:varchar(20);default:'rolling'" json:"schedule_mode"` // rolling | calendar
	PreferredWeekdays pq.Int32Array `gorm:"type:int[];default:'{}'" json:"preferred_weekdays"`       // only used in calendar mode (0=Mon..6=Sun)

	Status    string         `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type Exercise struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Slug         string         `gorm:"type:varchar(255);not null;unique" json:"slug"`
	Name         string         `gorm:"type:text;not null;unique" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	IsBodyweight bool           `gorm:"default:false" json:"is_bodyweight"`
	Instructions string         `gorm:"type:text" json:"instructions"`
	VideoURL     string         `gorm:"type:text" json:"video_url"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Prescriptions    []WorkoutPrescription `gorm:"foreignKey:ExerciseID" json:"prescriptions,omitempty"`
	SessionExercises []SessionExercise     `gorm:"foreignKey:ExerciseID" json:"session_exercises,omitempty"`
	MuscleGroups     []ExerciseMuscleGroup `gorm:"foreignKey:ExerciseID" json:"muscle_groups,omitempty"`
	Equipment        []ExerciseEquipment   `gorm:"foreignKey:ExerciseID" json:"equipment,omitempty"`
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

func (wpi *WorkoutPlanItem) BeforeCreate(tx *gorm.DB) (err error) {
	if wpi.ID == uuid.Nil {
		wpi.ID = uuid.New()
	}
	return
}

func (pe *PlanEnrollment) BeforeCreate(tx *gorm.DB) (err error) {
	if pe.ID == uuid.Nil {
		pe.ID = uuid.New()
	}
	return
}
