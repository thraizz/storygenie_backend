package models

import (
	uuid "github.com/google/uuid"
	"gorm.io/gorm"
)

type Feedback struct {
	gorm.Model
	UID     uuid.UUID `json:"id"`
	Liked   bool      `json:"liked"`
	UserID  string    `json:"user_id"`
	Story   Story
	StoryID uuid.UUID `json:"storyId"`
}

// Set a UUID as the primary key
func (feedback *Feedback) BeforeCreate(tx *gorm.DB) (err error) {
	feedback.UID = uuid.New()
	return
}
