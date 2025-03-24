package models

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
}

func MigrateBookTable(db *gorm.DB) error {
	return db.AutoMigrate(&Book{})
}
