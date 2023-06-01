package config

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

var (
	db *gorm.DB
)

func Connect() {
	envMap, _ := godotenv.Read("../.env")
	d, err := gorm.Open("mysql", envMap["DB_URL"])
	if err != nil {
		panic(err)
	}
	db = d
}

func GetDB() *gorm.DB {
	return db
}
