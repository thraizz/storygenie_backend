package models

import (
	uuid "github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type Product struct {
	gorm.Model
	UID         uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UserID      string    `json:"user_id"`
	IsExample   bool      `json:"isExample"`
	Story       []Story
}

// Set a UUID as the primary key
func (product *Product) BeforeCreate(tx *gorm.DB) (err error) {
	product.UID = uuid.New()
	return
}
