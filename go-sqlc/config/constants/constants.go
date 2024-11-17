package constants

import "time"

const (
	// Database Types
	SQLite   = "sqlite3"
	Postgres = "postgres"
	MySQL    = "mysql"

	// Environment Prefixes
	EnvPrefixSQLite   = "SQLITE_"
	EnvPrefixPostgres = "PG_"
	EnvPrefixMySQL    = "MYSQL_"
	EnvPrefixRedis    = "REDIS_"
	EnvPrefixServer   = "SERVER_"

	// Default Values
	DefaultMaxOpenConns  = 10
	DefaultMaxIdleConns  = 5
	DefaultMaxLifetime   = 15 * time.Minute
	DefaultRetryAttempts = 3
	DefaultRetryDelay    = time.Second * 3
)
