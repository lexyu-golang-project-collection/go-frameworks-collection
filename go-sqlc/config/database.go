package config

import (
	"time"
)

type DatabaseConfig struct {
	Driver       string
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

type Config struct {
	SQLite   DatabaseConfig
	Postgres DatabaseConfig
	MySQL    DatabaseConfig
}

func NewConfig() *Config {
	return &Config{
		SQLite: DatabaseConfig{
			Driver:       "sqlite3",
			DSN:          "file:authors.db?cache=shared&mode=rwc",
			MaxOpenConns: 10,
			MaxIdleConns: 5,
			MaxLifetime:  15 * time.Minute,
		},
		Postgres: DatabaseConfig{
			Driver:       "postgres",
			DSN:          "postgresql://user:pass@localhost:5432/authors?sslmode=disable",
			MaxOpenConns: 25,
			MaxIdleConns: 5,
			MaxLifetime:  30 * time.Minute,
		},
		MySQL: DatabaseConfig{
			Driver:       "mysql",
			DSN:          "user:pass@tcp(localhost:3306)/authors",
			MaxOpenConns: 25,
			MaxIdleConns: 5,
			MaxLifetime:  30 * time.Minute,
		},
	}
}
