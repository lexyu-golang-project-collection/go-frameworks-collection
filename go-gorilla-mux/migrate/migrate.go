package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/config"
	models "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/model"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: 無法載入 .env 檔案，使用預設值")
	}

	if err := config.ConnectToDatabases(); err != nil {
		log.Fatalf("無法連接到資料庫：%v", err)
	}
	defer config.CloseConnections()

	log.Println("開始資料庫遷移...")

	if err := models.MigrateBookTable(config.GetMySQLDB()); err != nil {
		log.Printf("Book 表格遷移失敗：%v", err)
	} else {
		log.Println("Book 表格遷移成功")
	}

	if err := models.MigratePostTable(config.GetPostgresDB()); err != nil {
		log.Printf("Post 表格遷移失敗：%v", err)
	} else {
		log.Println("Post 表格遷移成功")
	}

	models.InitMovies()
	log.Println("電影資料初始化成功")

	log.Println("遷移完成")
}
