package crud

import (
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func InitDB(db *gorm.DB) Database {
	return Database{db}
}
