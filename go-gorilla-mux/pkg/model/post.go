package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title string `json:"title"`
	Body  string `json:"body"`
}

func MigratePostTable(db *gorm.DB) error {
	return db.AutoMigrate(&Post{})
}
