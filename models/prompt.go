// A prompt look like {Description: requestData.StoryDescription, Version: 1.1}

package models

import "gorm.io/gorm"

type Prompt struct {
	gorm.Model
	// The description of the prompt
	Description string `json:"description"`
	// The version of the prompt
	Version string `json:"version"`
	UserID  string `json:"user_id"`
}
