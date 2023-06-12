package models

import (
	"time"

	uuid "github.com/google/uuid"
	"gorm.io/gorm"
)

type JiraRefreshToken struct {
	UID          uuid.UUID `json:"id" gorm:"primary_key;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	UserID       string         `json:"user_id"`
	RefreshToken string         `gorm:"type:text" json:"refreshToken"`
}

func (jiraRefreshToken *JiraRefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	jiraRefreshToken.UID = uuid.New()
	return
}
