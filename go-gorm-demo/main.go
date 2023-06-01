package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Category struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Product struct {
	gorm.Model
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func main() {
	dsn := "root:00112125@tcp(127.0.0.1:3306)/gorm_test_db?charset=utf8&parseTime=True&loc=Local"
	dial := mysql.Open(dsn)
	db, err := gorm.Open(dial, &gorm.Config{})
	if err != nil {
		panic("failed to connect to db")
	}

	// Migrate the schema
	db.AutoMigrate(&User{})

	db.Migrator().CreateTable(&Category{})

	db.Migrator().CreateTable(&Product{})

}
