package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UserID      string    `json:"user_id"`
	IsExample   bool      `json:"isExample"`
	Story       []Story
}

// Set a UUID as the primary key
func (project *Project) BeforeCreate(tx *gorm.DB) (err error) {
	project.ID = uuid.NewV4()
	return
}
