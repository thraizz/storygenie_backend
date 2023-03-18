package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID string `gorm:"primaryKey"`
	// the firebase uid, it is used to uniquely identify a user
	UID         string `gorm:"unique;not null"`
	IsAnonymous bool
}
