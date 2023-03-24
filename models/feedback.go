package models

import (
	"time"

	uuid "github.com/google/uuid"
	"gorm.io/gorm"
)

type Feedback struct {
	UID       uuid.UUID `json:"id" gorm:"primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Liked     bool           `json:"liked"`
	UserID    string         `json:"user_id"`
	Story     Story
	StoryID   uuid.UUID `json:"storyId"`
}

// Set a UUID as the primary key
func (feedback *Feedback) BeforeCreate(tx *gorm.DB) (err error) {
	feedback.UID = uuid.New()
	return
}
