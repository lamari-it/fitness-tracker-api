package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkoutPlan struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Title       string    `gorm:"type:text;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Visibility  string    `gorm:"type:varchar(20);default:'private'" json:"visibility"` // private, public, friends
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	User     User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Workouts []Workout `gorm:"foreignKey:PlanID" json:"workouts,omitempty"`
}

type Workout struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PlanID    uuid.UUID `gorm:"type:uuid;not null" json:"plan_id"`
	Title     string    `gorm:"type:text;not null" json:"title"`
	DayNumber int       `gorm:"not null" json:"day_number"`
	Notes     string    `gorm:"type:text" json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Plan              WorkoutPlan       `gorm:"foreignKey:PlanID;constraint:OnDelete:CASCADE" json:"plan,omitempty"`
	WorkoutExercises  []WorkoutExercise `gorm:"foreignKey:WorkoutID" json:"exercises,omitempty"`
	WorkoutSessions   []WorkoutSession  `gorm:"foreignKey:WorkoutID" json:"sessions,omitempty"`
	SharedWorkouts    []SharedWorkout   `gorm:"foreignKey:WorkoutID" json:"shared_workouts,omitempty"`
}

type Exercise struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name         string    `gorm:"type:text;not null;unique" json:"name"`
	Description  string    `gorm:"type:text" json:"description"`
	IsBodyweight bool      `gorm:"default:false" json:"is_bodyweight"`
	Instructions string    `gorm:"type:text" json:"instructions"`
	VideoURL     string    `gorm:"type:text" json:"video_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relationships
	WorkoutExercises []WorkoutExercise     `gorm:"foreignKey:ExerciseID" json:"workout_exercises,omitempty"`
	ExerciseLogs     []ExerciseLog         `gorm:"foreignKey:ExerciseID" json:"exercise_logs,omitempty"`
	MuscleGroups     []ExerciseMuscleGroup `gorm:"foreignKey:ExerciseID" json:"muscle_groups,omitempty"`
	Equipment        []ExerciseEquipment   `gorm:"foreignKey:ExerciseID" json:"equipment,omitempty"`
}

type WorkoutExercise struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	WorkoutID      uuid.UUID `gorm:"type:uuid;not null" json:"workout_id"`
	ExerciseID     uuid.UUID `gorm:"type:uuid;not null" json:"exercise_id"`
	OrderNumber    int       `gorm:"not null" json:"order_number"`
	TargetSets     int       `json:"target_sets"`
	TargetReps     int       `json:"target_reps"`
	TargetWeight   float64   `gorm:"type:numeric(10,2)" json:"target_weight"`
	TargetRestSec  int       `json:"target_rest_sec"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Workout  Workout  `gorm:"foreignKey:WorkoutID;constraint:OnDelete:CASCADE" json:"workout,omitempty"`
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