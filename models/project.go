package models

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	// The name of the project
	Name string `json:"name"`
	// The description of the project
	Description string `json:"description"`
	// The user id of the project
	UserID    string `json:"user_id"`
	IsExample bool   `json:"isExample"`
}
