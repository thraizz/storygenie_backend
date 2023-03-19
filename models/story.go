package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Story struct {
	gorm.Model
	ID                 uuid.UUID      `json:"id"`
	UserID             string         `json:"user_id"`
	Headline           string         `gorm:"type:text" json:"headline"`
	UserStory          string         `gorm:"type:text" json:"userStory"`
	AcceptanceCriteria datatypes.JSON `json:"acceptanceCriteria"`
	ProjectID          uint           `json:"projectId"`
}

type GetStoryInput struct {
	ID uint `json:"id"`
}

// Set a UUID as the primary key
func (story *Story) BeforeCreate(tx *gorm.DB) (err error) {
	story.ID = uuid.NewV4()
	return
}
