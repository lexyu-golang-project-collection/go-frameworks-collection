package main

import (
	_ "embed"
	"log"

	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/config"
	embedSQL "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/db/embed"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/repository"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/service"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/router"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	embedSQL.GetSQLiteDDL()
	embedSQL.GetMySQLDDL()
	embedSQL.GetPostgresDDL()
}

func main() {
	cfg := config.NewConfig()

	// init DBs
	dbManager, err := repository.NewDBManager(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbManager.Close()

	// init service
	services := service.NewServices(dbManager)

	// start server
	srv := router.NewServer(services)
	srv.Start()
}
