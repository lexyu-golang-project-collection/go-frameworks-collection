package main

import (
	"go-crud-api-with-gin-gorm-postgres/initializers"
	"go-crud-api-with-gin-gorm-postgres/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.Post{})
}
