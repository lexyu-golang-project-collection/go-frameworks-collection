package main

import (
	_ "embed"
	"log"

	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/config"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/config/constants"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/repository"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/service"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/router"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	env := config.GetEnvConfig()

	// Create config with all database options
	cfg := config.NewConfig(
		// SQLite configuration
		config.WithDatabase(constants.SQLite, env).
			DSN().InMemory().MaxOpenConns().MaxIdleConns().MaxLifetime().
			Build(),

		// MySQL configuration
		config.WithDatabase(constants.MySQL, env).
			DSN().MaxOpenConns().MaxIdleConns().MaxLifetime().
			Build(),

		// Postgres configuration
		config.WithDatabase(constants.Postgres, env).
			DSN().MaxOpenConns().MaxIdleConns().MaxLifetime().
			Build(),

		// Global retry policy
		config.WithRetryPolicy(env).Build(),
	)

	// init DB Manager
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
