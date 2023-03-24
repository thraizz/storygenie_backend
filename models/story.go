package models

import (
	uuid "github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Story struct {
	gorm.Model
	UID                uuid.UUID                    `json:"id"`
	UserID             string                       `json:"user_id"`
	Headline           string                       `gorm:"type:text" json:"headline"`
	UserStory          string                       `gorm:"type:text" json:"userStory"`
	AcceptanceCriteria datatypes.JSONType[[]string] `json:"acceptanceCriteria"`
	Product            Product
	ProductID          uint `json:"productId"`
}

type GetStoryInput struct {
	ID uint `json:"id"`
}

// Set a UUID as the primary key
func (story *Story) BeforeCreate(tx *gorm.DB) (err error) {
	story.UID = uuid.New()
	return
}
