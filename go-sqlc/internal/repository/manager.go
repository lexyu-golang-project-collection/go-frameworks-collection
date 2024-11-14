package repository

import (
	"database/sql"

	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/config"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/db/mysql"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/db/postgres"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/db/sqlite"
)

type DBManager struct {
	sqlite   *sql.DB
	postgres *sql.DB
	mysql    *sql.DB

	SqliteQueries   *sqlite.Queries
	PostgresQueries *postgres.Queries
	MysqlQueries    *mysql.Queries
}

// DB Manager
func NewDBManager(cfg *config.Config) (*DBManager, error) {
	mgr := &DBManager{}

	// init SQLite
	sqliteDB, err := sql.Open(cfg.SQLite.Driver, cfg.SQLite.DSN)
	if err != nil {
		return nil, err
	}
	configureDB(sqliteDB, cfg.SQLite)
	mgr.sqlite = sqliteDB
	mgr.SqliteQueries = sqlite.New(sqliteDB)

	// init Postgres
	postgresDB, err := sql.Open(cfg.Postgres.Driver, cfg.Postgres.DSN)
	if err != nil {
		return nil, err
	}
	configureDB(postgresDB, cfg.Postgres)
	mgr.postgres = postgresDB
	mgr.PostgresQueries = postgres.New(postgresDB)

	// init MySQL
	mysqlDB, err := sql.Open(cfg.MySQL.Driver, cfg.MySQL.DSN)
	if err != nil {
		return nil, err
	}
	configureDB(mysqlDB, cfg.MySQL)
	mgr.mysql = mysqlDB
	mgr.MysqlQueries = mysql.New(mysqlDB)

	return mgr, nil
}

func configureDB(db *sql.DB, cfg config.DatabaseConfig) {
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxLifetime)
}

// Close 關閉所有數據庫連接
func (m *DBManager) Close() error {
	if m.SqliteQueries != nil {
		m.SqliteQueries.Close()
	}
	if m.PostgresQueries != nil {
		m.PostgresQueries.Close()
	}
	if m.MysqlQueries != nil {
		m.MysqlQueries.Close()
	}

	var err error
	if m.sqlite != nil {
		err = m.sqlite.Close()
	}
	if m.postgres != nil {
		err = m.postgres.Close()
	}
	if m.mysql != nil {
		err = m.mysql.Close()
	}
	return err
}
