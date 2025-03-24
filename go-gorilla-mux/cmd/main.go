package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/docs"

	"github.com/joho/godotenv"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/config"
	controllers "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/controller"
	models "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/model"
	base "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/repository/interfaces"
	repo_mysql "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/repository/mysql"
	repo_postgres "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/repository/postgres"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/routes"
	services "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-gorilla-mux/pkg/service"
)

// @title 整合 API
// @version 1.0
// @description 電影、書籍和帖子的綜合 RESTful API
// @termsOfService http://swagger.io/terms/
// @contact.name API 支援
// @contact.url http://www.example.com/support
// @contact.email support@example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: 無法載入 .env 檔案，使用預設值")
	}

	if err := config.ConnectToDatabases(); err != nil {
		log.Fatalf("無法連接到資料庫：%v", err)
	}

	defer config.CloseConnections()

	models.InitMovies()

	mysqlDB := config.GetMySQLDB()
	postgresDB := config.GetPostgresDB()

	mysqlTxManager := base.NewTxManager(mysqlDB)
	postgresTxManager := base.NewTxManager(postgresDB)

	movieRepo := repo_mysql.NewMovieRepository()
	bookRepo := repo_mysql.NewBookRepository(mysqlDB)
	postRepo := repo_postgres.NewPostRepository(postgresDB)

	movieService := services.NewMovieService(movieRepo)
	bookService := services.NewBookService(bookRepo, mysqlTxManager)
	postService := services.NewPostService(postRepo, postgresTxManager)

	movieController := controllers.NewMovieController(movieService)
	bookController := controllers.NewBookController(bookService)
	postController := controllers.NewPostController(postService)

	router := routes.NewRouter(movieController, bookController, postController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		fmt.Printf("伺服器啟動於 :%s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("無法啟動伺服器：%v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在關閉伺服器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("伺服器強制關閉：%v", err)
	}

	log.Println("伺服器優雅地關閉")
}
