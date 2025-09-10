package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID      uint32 `gorm:"primaryKey"`
	Phone     int    `json:"phone" gorm:"unique"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Password  string `json:"password"`
	Code      string
	ExpiresAt time.Time
	Used      bool `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
