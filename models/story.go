package models

import (
	"time"

	uuid "github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Story struct {
	UID                uuid.UUID `json:"id" gorm:"primary_key;"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
	UserID             string         `json:"user_id"`
	Headline           string         `gorm:"type:text" json:"headline"`
	UserStory          string         `gorm:"type:text" json:"userStory"`
	AcceptanceCriteria datatypes.JSON `json:"acceptanceCriteria"`
	Product            Product
	ProductID          uuid.UUID `json:"productId"`
	JiraIssueID        string    `json:"jiraIssueId"`
}

type GetStoryInput struct {
	ID uint `json:"id"`
}

// Set a UUID as the primary key
func (story *Story) BeforeCreate(tx *gorm.DB) (err error) {
	story.UID = uuid.New()
	return
}
