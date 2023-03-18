package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Story struct {
	gorm.Model
	// The headline of the story
	Headline string `gorm:"type:text" json:"headline"`
	// The user story of the story
	UserStory string `gorm:"type:text" json:"userStory"`
	// The acceptance criteria of the story
	AcceptanceCriteria datatypes.JSON `gorm:"type:text[]"`
	// Each story is linked to a project, this foregin key is the project id
	ProjectID string `json:"projectId"`
	Project   Project
	// Each story has a prompt, this foregin key is the prompt id
	PromptID string `json:"promptId"`
	Prompt   Prompt `gorm:"foreignKey:PromptID"`
	UserID   string `json:"user_id"`
}

type GetStoryInput struct {
	ID uint `json:"id"`
}
