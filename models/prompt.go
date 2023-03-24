// A prompt look like {Description: requestData.StoryDescription, Version: 1.1}

package models

import (
	"time"

	uuid "github.com/google/uuid"
	"gorm.io/gorm"
)

type Prompt struct {
	UID       uuid.UUID `json:"id" gorm:"primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	// The description of the prompt
	Description string `json:"description"`
	// The version of the prompt
	Version string `json:"version"`
	UserID  string `json:"user_id"`
}
