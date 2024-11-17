package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/config"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/config/constants"
	embedSQL "github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/db/embed"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/db/mysql"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/db/postgres"
	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/internal/db/sqlite"
)

type DBManager struct {
	config  *config.Config
	dbs     map[string]*sql.DB
	queries map[string]interface{}
}

func NewDBManager(cfg *config.Config) (*DBManager, error) {
	manager := &DBManager{
		config:  cfg,
		dbs:     make(map[string]*sql.DB),
		queries: make(map[string]interface{}),
	}

	if err := manager.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize database manager: %w", err)
	}

	return manager, nil
}

func (m *DBManager) initialize() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var initErrors []error

	// Initialize enabled databases
	for dbType := range m.config.DB.EnabledDatabases {
		if err := m.initDatabase(ctx, dbType); err != nil {
			initErrors = append(initErrors, fmt.Errorf("%s initialization failed: %w", dbType, err))
		}
	}

	// Check if any database was successfully initialized
	if len(m.dbs) == 0 && len(initErrors) > 0 {
		return fmt.Errorf("all database initializations failed: %v", initErrors)
	}

	if len(initErrors) > 0 {
		log.Printf("Some database initializations failed: %v", initErrors)
	}

	return nil
}

func (m *DBManager) initDatabase(ctx context.Context, dbType string) error {
	var db *sql.DB
	var err error
	var dbConfig *config.DatabaseConfig

	switch dbType {
	case constants.SQLite:
		dbConfig = m.config.DB.SQLite
		if dbConfig.InMemory {
			db, err = m.initSQLiteMemory(ctx)
		} else {
			db, err = m.initSQLiteFile(ctx, dbConfig)
		}
	case constants.MySQL:
		dbConfig = m.config.DB.MySQL
		db, err = m.initMySQL(ctx, dbConfig)
	case constants.Postgres:
		dbConfig = m.config.DB.Postgres
		db, err = m.initPostgres(ctx, dbConfig)
	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	if err != nil {
		return err
	}

	// Configure connection pool
	m.configureDB(db, dbConfig)
	m.dbs[dbType] = db

	// Initialize queries
	if err := m.initQueries(dbType, db); err != nil {
		db.Close()
		delete(m.dbs, dbType)
		return err
	}

	return nil
}

func (m *DBManager) initSQLiteMemory(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("failed to open in-memory SQLite: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping in-memory SQLite: %w", err)
	}

	if err := m.ensureTablesExist(ctx, db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func (m *DBManager) initSQLiteFile(ctx context.Context, config *config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", config.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite file: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping SQLite file: %w", err)
	}

	return db, nil
}

func (m *DBManager) initMySQL(ctx context.Context, config *config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open MySQL: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping MySQL: %w", err)
	}

	return db, nil
}

func (m *DBManager) initPostgres(ctx context.Context, config *config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open Postgres: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping Postgres: %w", err)
	}

	return db, nil
}

func (m *DBManager) initQueries(dbType string, db *sql.DB) error {
	switch dbType {
	case constants.SQLite:
		m.queries[dbType] = sqlite.New(db)
	case constants.MySQL:
		m.queries[dbType] = mysql.New(db)
	case constants.Postgres:
		m.queries[dbType] = postgres.New(db)
	default:
		return fmt.Errorf("unsupported database type for queries: %s", dbType)
	}
	return nil
}

func (m *DBManager) ensureTablesExist(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, embedSQL.GetSQLiteDDL()); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}
	return nil
}

func (m *DBManager) configureDB(db *sql.DB, cfg *config.DatabaseConfig) {
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxLifetime)
}

// GetQueries returns the queries for the specified database type
func (m *DBManager) GetQueries(dbType string) interface{} {
	return m.queries[dbType]
}

// Close closes all database connections
func (m *DBManager) Close() error {
	var errs []error

	// Close all database connections
	for dbType, db := range m.dbs {
		if err := db.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close %s: %w", dbType, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing databases: %v", errs)
	}
	return nil
}

// Helper methods for type assertions
func (m *DBManager) SQLiteQueries() *sqlite.Queries {
	if q, ok := m.queries[constants.SQLite].(*sqlite.Queries); ok {
		return q
	}
	return nil
}

func (m *DBManager) MySQLQueries() *mysql.Queries {
	if q, ok := m.queries[constants.MySQL].(*mysql.Queries); ok {
		return q
	}
	return nil
}

func (m *DBManager) PostgresQueries() *postgres.Queries {
	if q, ok := m.queries[constants.Postgres].(*postgres.Queries); ok {
		return q
	}
	return nil
}

// type DBManager struct {
// 	config   *config.DBConfig
// 	sqlite   *sql.DB
// 	postgres *sql.DB
// 	mysql    *sql.DB

//		SqliteQueries   *sqlite.Queries
//		PostgresQueries *postgres.Queries
//		MysqlQueries    *mysql.Queries
//	}

// type Option func(*config.DBConfig)

// func WithSQLite(enable bool, isMemory bool) Option {
// 	return func(c *config.DBConfig) {
// 		c.EnabledDatabases[constants.SQLite] = enable
// 		c.SQLite.InMemory = isMemory
// 	}
// }

// func WithMySQL(enable bool) Option {
// 	return func(c *config.DBConfig) {
// 		c.EnabledDatabases[constants.MySQL] = enable
// 	}
// }

// func WithPostgres(enable bool) Option {
// 	return func(c *config.DBConfig) {
// 		c.EnabledDatabases[constants.Postgres] = enable
// 	}
// }

// func WithRetry(attempts int, delay time.Duration) Option {
// 	return func(c *config.DBConfig) {
// 		c.RetryPolicy.Attempts = attempts
// 		c.RetryPolicy.Delay = delay
// 	}
// }

// func NewDBManager(cfg *config.Config, opts ...Option) (*DBManager, error) {
// 	for _, opt := range opts {
// 		opt(cfg.DB)
// 	}

// 	mgr := &DBManager{
// 		config: cfg.DB,
// 	}

// 	if err := mgr.initializeConnections(); err != nil {
// 		return nil, err
// 	}

// 	return mgr, nil
// }

// func (m *DBManager) initializeConnections() error {
// 	var errs []error

// 	// init SQLite
// 	if m.config.EnabledDatabases[constants.SQLite] {
// 		if err := m.initSQLite(); err != nil {
// 			log.Printf("SQLite initialization failed: %v", err)
// 			// file mode failed
// 			if !m.config.SQLite.InMemory {
// 				log.Printf("Switching to in-memory SQLite mode")
// 				m.config.SQLite.InMemory = true
// 				// try Memory mode
// 				if retryErr := m.initSQLite(); retryErr != nil {
// 					// Both Failed
// 					errs = append(errs, fmt.Errorf("SQLite initialization failed - File mode: %v, Memory mode: %v", err, retryErr))
// 				} else {
// 					log.Printf("Successfully switched to in-memory SQLite mode")
// 				}
// 			} else {
// 				errs = append(errs, fmt.Errorf("in-memory SQLite initialization failed: %v", err))
// 			}
// 		}
// 	}

// 	// init MySQL
// 	if m.config.EnabledDatabases[constants.MySQL] {
// 		if err := m.initMySQL(); err != nil {
// 			errs = append(errs, fmt.Errorf("mysql initialization failed: %v", err))
// 		}
// 	}

// 	// init Postgres
// 	if m.config.EnabledDatabases[constants.Postgres] {
// 		if err := m.initPostgres(); err != nil {
// 			errs = append(errs, fmt.Errorf("postgres initialization failed: %v", err))
// 		}
// 	}

// 	// init successfully at least one
// 	if !m.hasAnyValidConnection() {
// 		var errStr string
// 		for _, err := range errs {
// 			errStr += err.Error() + "; "
// 		}
// 		return fmt.Errorf("no valid database connections available: %s", errStr)
// 	}

// 	if len(errs) > 0 {
// 		log.Printf("Some database connections failed: %v", errs)
// 	}

// 	return nil
// }

// func (m *DBManager) hasAnyValidConnection() bool {
// 	return m.sqlite != nil || m.postgres != nil || m.mysql != nil
// }

// func (m *DBManager) initSQLite() error {
// 	var db *sql.DB
// 	var err error

// 	if m.config.SQLite.InMemory {
// 		db, err = sql.Open("sqlite3", ":memory:")
// 		if err != nil {
// 			return fmt.Errorf("memory sqlite open error: %w", err)
// 		}
// 	} else {
// 		db, err = sql.Open(m.config.SQLite.Driver, m.config.SQLite.DSN)
// 		if err != nil {
// 			return fmt.Errorf("sqlite file open error: %w", err)
// 		}
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	if err := db.PingContext(ctx); err != nil {
// 		db.Close()
// 		return fmt.Errorf("sqlite ping error: %w", err)
// 	}

// 	if m.config.SQLite.InMemory {
// 		if err := m.ensureTablesExist(ctx, db); err != nil {
// 			db.Close()
// 			return err
// 		}
// 	}

// 	configureDB(db, m.config.SQLite)
// 	m.sqlite = db
// 	m.SqliteQueries = sqlite.New(db)
// 	return nil
// }

// func (m *DBManager) ensureTablesExist(ctx context.Context, db *sql.DB) error {
// 	if _, err := db.ExecContext(ctx, embedSQL.GetSQLiteDDL()); err != nil {
// 		return fmt.Errorf("failed to create tables: %w", err)
// 	}
// 	return nil
// }

// func (m *DBManager) initMySQL() error {
// 	db, err := sql.Open(m.config.MySQL.Driver, m.config.MySQL.DSN)
// 	if err != nil {
// 		return fmt.Errorf("mysql open error: %w", err)
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	if err := db.PingContext(ctx); err != nil {
// 		db.Close()
// 		return fmt.Errorf("mysql ping error: %w", err)
// 	}

// 	configureDB(db, m.config.MySQL)
// 	m.mysql = db
// 	m.MysqlQueries = mysql.New(db)
// 	return nil
// }

// func (m *DBManager) initPostgres() error {
// 	db, err := sql.Open(m.config.Postgres.Driver, m.config.Postgres.DSN)
// 	if err != nil {
// 		return fmt.Errorf("postgres open error: %w", err)
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	if err := db.PingContext(ctx); err != nil {
// 		db.Close()
// 		return fmt.Errorf("postgres ping error: %w", err)
// 	}

// 	configureDB(db, m.config.Postgres)
// 	m.postgres = db
// 	m.PostgresQueries = postgres.New(db)
// 	return nil
// }

// func (m *DBManager) Close() error {
// 	var errs []error

// 	if m.SqliteQueries != nil {
// 		if err := m.SqliteQueries.Close(); err != nil {
// 			errs = append(errs, fmt.Errorf("sqlite queries close error: %w", err))
// 		}
// 	}
// 	if m.PostgresQueries != nil {
// 		if err := m.PostgresQueries.Close(); err != nil {
// 			errs = append(errs, fmt.Errorf("postgres queries close error: %w", err))
// 		}
// 	}
// 	if m.MysqlQueries != nil {
// 		if err := m.MysqlQueries.Close(); err != nil {
// 			errs = append(errs, fmt.Errorf("mysql queries close error: %w", err))
// 		}
// 	}

// 	// dbs closed
// 	if m.sqlite != nil {
// 		if err := m.sqlite.Close(); err != nil {
// 			errs = append(errs, fmt.Errorf("sqlite close error: %w", err))
// 		}
// 	}
// 	if m.postgres != nil {
// 		if err := m.postgres.Close(); err != nil {
// 			errs = append(errs, fmt.Errorf("postgres close error: %w", err))
// 		}
// 	}
// 	if m.mysql != nil {
// 		if err := m.mysql.Close(); err != nil {
// 			errs = append(errs, fmt.Errorf("mysql close error: %w", err))
// 		}
// 	}

// 	if len(errs) > 0 {
// 		return fmt.Errorf("multiple close errors: %v", errs)
// 	}
// 	return nil
// }

// func configureDB(db *sql.DB, cfg *config.DatabaseConfig) {
// 	db.SetMaxOpenConns(cfg.MaxOpenConns)
// 	db.SetMaxIdleConns(cfg.MaxIdleConns)
// 	db.SetConnMaxLifetime(cfg.MaxLifetime)
// }
