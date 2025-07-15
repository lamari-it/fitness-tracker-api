package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SharedWorkout struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	WorkoutID     uuid.UUID `gorm:"type:uuid;not null" json:"workout_id"`
	SharedByID    uuid.UUID `gorm:"type:uuid;not null" json:"shared_by_id"`
	SharedWithID  uuid.UUID `gorm:"type:uuid;not null" json:"shared_with_id"`
	Permission    string    `gorm:"type:varchar(20);default:'view'" json:"permission"` // view, edit, copy
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Workout      Workout           `gorm:"foreignKey:WorkoutID;constraint:OnDelete:CASCADE" json:"workout,omitempty"`
	SharedBy     User              `gorm:"foreignKey:SharedByID;constraint:OnDelete:CASCADE" json:"shared_by,omitempty"`
	SharedWith   User              `gorm:"foreignKey:SharedWithID;constraint:OnDelete:CASCADE" json:"shared_with,omitempty"`
	Comments     []WorkoutComment  `gorm:"foreignKey:SharedWorkoutID" json:"comments,omitempty"`
}

type WorkoutComment struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SharedWorkoutID  uuid.UUID `gorm:"type:uuid;not null" json:"shared_workout_id"`
	UserID           uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	ParentID         *uuid.UUID `gorm:"type:uuid" json:"parent_id"` // for threaded comments
	Content          string    `gorm:"type:text;not null" json:"content"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	SharedWorkout SharedWorkout              `gorm:"foreignKey:SharedWorkoutID;constraint:OnDelete:CASCADE" json:"shared_workout,omitempty"`
	User          User                       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Parent        *WorkoutComment            `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"parent,omitempty"`
	Replies       []WorkoutComment           `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
	Reactions     []WorkoutCommentReaction   `gorm:"foreignKey:CommentID" json:"reactions,omitempty"`
}

type WorkoutCommentReaction struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CommentID uuid.UUID `gorm:"type:uuid;not null" json:"comment_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Reaction  string    `gorm:"type:varchar(20);not null" json:"reaction"` // like, love, laugh, wow, sad, angry
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Comment WorkoutComment `gorm:"foreignKey:CommentID;constraint:OnDelete:CASCADE" json:"comment,omitempty"`
	User    User           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

func (sw *SharedWorkout) BeforeCreate(tx *gorm.DB) (err error) {
	if sw.ID == uuid.Nil {
		sw.ID = uuid.New()
	}
	return
}

func (wc *WorkoutComment) BeforeCreate(tx *gorm.DB) (err error) {
	if wc.ID == uuid.Nil {
		wc.ID = uuid.New()
	}
	return
}

func (wcr *WorkoutCommentReaction) BeforeCreate(tx *gorm.DB) (err error) {
	if wcr.ID == uuid.Nil {
		wcr.ID = uuid.New()
	}
	return
}

// Response DTOs
type SharedWorkoutResponse struct {
	ID           uuid.UUID                `json:"id"`
	WorkoutID    uuid.UUID                `json:"workout_id"`
	WorkoutTitle string                   `json:"workout_title"`
	SharedBy     UserResponse             `json:"shared_by"`
	SharedWith   UserResponse             `json:"shared_with"`
	Permission   string                   `json:"permission"`
	CommentsCount int                     `json:"comments_count"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
}

type WorkoutCommentResponse struct {
	ID              uuid.UUID                      `json:"id"`
	UserID          uuid.UUID                      `json:"user_id"`
	User            UserResponse                   `json:"user"`
	ParentID        *uuid.UUID                     `json:"parent_id"`
	Content         string                         `json:"content"`
	RepliesCount    int                            `json:"replies_count"`
	ReactionsCount  map[string]int                 `json:"reactions_count"`
	UserReaction    *string                        `json:"user_reaction,omitempty"`
	CreatedAt       time.Time                      `json:"created_at"`
	UpdatedAt       time.Time                      `json:"updated_at"`
}

func (sw *SharedWorkout) ToResponse() SharedWorkoutResponse {
	return SharedWorkoutResponse{
		ID:           sw.ID,
		WorkoutID:    sw.WorkoutID,
		WorkoutTitle: sw.Workout.Title,
		SharedBy:     sw.SharedBy.ToResponse(),
		SharedWith:   sw.SharedWith.ToResponse(),
		Permission:   sw.Permission,
		CommentsCount: len(sw.Comments),
		CreatedAt:    sw.CreatedAt,
		UpdatedAt:    sw.UpdatedAt,
	}
}

func (wc *WorkoutComment) ToResponse(currentUserID uuid.UUID) WorkoutCommentResponse {
	response := WorkoutCommentResponse{
		ID:             wc.ID,
		UserID:         wc.UserID,
		User:           wc.User.ToResponse(),
		ParentID:       wc.ParentID,
		Content:        wc.Content,
		RepliesCount:   len(wc.Replies),
		ReactionsCount: make(map[string]int),
		CreatedAt:      wc.CreatedAt,
		UpdatedAt:      wc.UpdatedAt,
	}

	// Count reactions by type
	for _, reaction := range wc.Reactions {
		response.ReactionsCount[reaction.Reaction]++
		if reaction.UserID == currentUserID {
			response.UserReaction = &reaction.Reaction
		}
	}

	return response
}