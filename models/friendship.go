package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Friendship struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	FriendID  uuid.UUID      `gorm:"type:uuid;not null" json:"friend_id"`
	Status    string         `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, accepted, blocked
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	User   User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Friend User `gorm:"foreignKey:FriendID;constraint:OnDelete:CASCADE" json:"friend,omitempty"`
}

func (f *Friendship) BeforeCreate(tx *gorm.DB) (err error) {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return
}

type FriendshipResponse struct {
	ID        uuid.UUID    `json:"id"`
	UserID    uuid.UUID    `json:"user_id"`
	FriendID  uuid.UUID    `json:"friend_id"`
	Status    string       `json:"status"`
	Friend    UserResponse `json:"friend"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func (f *Friendship) ToResponse() FriendshipResponse {
	return FriendshipResponse{
		ID:        f.ID,
		UserID:    f.UserID,
		FriendID:  f.FriendID,
		Status:    f.Status,
		Friend:    f.Friend.ToResponse(),
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}
