package config

import (
	"time"

	"github.com/lexyu-golang-project-collection/go-frameworks-collection/go-sqlc/config/constants"
)

type Config struct {
	DB *DBConfig
}

type DBConfig struct {
	EnabledDatabases map[string]bool
	SQLite           *DatabaseConfig
	Postgres         *DatabaseConfig
	MySQL            *DatabaseConfig
	RetryPolicy      RetryPolicy
}

type DatabaseConfig struct {
	Driver       string
	DSN          string
	InMemory     bool
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

type RetryPolicy struct {
	Attempts int
	Delay    time.Duration
}

type DBOption struct {
	env         *EnvConfig
	dbType      string
	prefixDB    string
	dsn         string
	inMemory    bool
	maxOpen     int
	maxIdle     int
	maxLifetime time.Duration
}

type ConfigOption func(*Config)

func WithDatabase(dbType string, env *EnvConfig) *DBOption {
	return &DBOption{
		env:         env,
		dbType:      dbType,
		prefixDB:    prefixDB(dbType),
		dsn:         "",
		maxOpen:     constants.DefaultMaxOpenConns,
		maxIdle:     constants.DefaultMaxIdleConns,
		maxLifetime: constants.DefaultMaxLifetime,
	}
}

func prefixDB(dbType string) string {
	switch dbType {
	case constants.MySQL:
		return constants.EnvPrefixMySQL
	case constants.Postgres:
		return constants.EnvPrefixPostgres
	default:
		return constants.EnvPrefixSQLite
	}
}

func (o *DBOption) DSN() *DBOption {
	if value, ok := o.env.GetString(o.prefixDB + "DSN"); ok {
		o.dsn = value
	}
	return o
}

func (o *DBOption) InMemory() *DBOption {
	if value, ok := o.env.GetBool(o.prefixDB + "IN_MEMORY"); ok {
		o.inMemory = value
	}
	return o
}

func (o *DBOption) MaxOpenConns() *DBOption {
	if value, ok := o.env.GetInt(o.prefixDB + "MAX_OPEN_CONNS"); ok {
		o.maxOpen = value
	}
	return o
}

func (o *DBOption) MaxIdleConns() *DBOption {
	if value, ok := o.env.GetInt(o.prefixDB + "MAX_IDLE_CONNS"); ok {
		o.maxIdle = value
	}
	return o
}

func (o *DBOption) MaxLifetime() *DBOption {
	if value, ok := o.env.GetDuration(o.prefixDB + "MAX_LIFETIME"); ok {
		o.maxLifetime = value
	}
	return o
}

// Builder
func (o *DBOption) Build() ConfigOption {
	return func(c *Config) {
		dbConfig := &DatabaseConfig{
			Driver:       o.dbType,
			DSN:          o.dsn,
			InMemory:     o.inMemory,
			MaxOpenConns: o.maxOpen,
			MaxIdleConns: o.maxIdle,
			MaxLifetime:  o.maxLifetime,
		}

		switch o.dbType {
		case constants.SQLite:
			c.DB.SQLite = dbConfig
		case constants.MySQL:
			c.DB.MySQL = dbConfig
		case constants.Postgres:
			c.DB.Postgres = dbConfig
		}
		c.DB.EnabledDatabases[o.dbType] = true
	}
}

func WithSQLite(dsn string, inMemory bool) ConfigOption {
	return func(c *Config) {
		c.DB.SQLite = &DatabaseConfig{
			Driver:       constants.SQLite,
			DSN:          dsn,
			InMemory:     inMemory,
			MaxOpenConns: constants.DefaultMaxOpenConns,
			MaxIdleConns: constants.DefaultMaxIdleConns,
			MaxLifetime:  constants.DefaultMaxLifetime,
		}
		c.DB.EnabledDatabases[constants.SQLite] = true
	}
}

func WithMySQL(dsn string) ConfigOption {
	return func(c *Config) {
		c.DB.MySQL = &DatabaseConfig{
			Driver:       constants.MySQL,
			DSN:          dsn,
			MaxOpenConns: constants.DefaultMaxOpenConns,
			MaxIdleConns: constants.DefaultMaxIdleConns,
			MaxLifetime:  constants.DefaultMaxLifetime,
		}
		c.DB.EnabledDatabases[constants.MySQL] = true
	}
}

func WithPostgres(dsn string) ConfigOption {
	return func(c *Config) {
		c.DB.Postgres = &DatabaseConfig{
			Driver:       constants.Postgres,
			DSN:          dsn,
			MaxOpenConns: constants.DefaultMaxOpenConns,
			MaxIdleConns: constants.DefaultMaxIdleConns,
			MaxLifetime:  constants.DefaultMaxLifetime,
		}
		c.DB.EnabledDatabases[constants.Postgres] = true
	}
}

func WithMaxOpenConns(dbType string, count int) ConfigOption {
	return func(c *Config) {
		switch dbType {
		case constants.SQLite:
			if c.DB.SQLite != nil {
				c.DB.SQLite.MaxOpenConns = count
			}
		case constants.MySQL:
			if c.DB.MySQL != nil {
				c.DB.MySQL.MaxOpenConns = count
			}
		case constants.Postgres:
			if c.DB.Postgres != nil {
				c.DB.Postgres.MaxOpenConns = count
			}
		}
	}
}

func WithMaxIdleConns(dbType string, count int) ConfigOption {
	return func(c *Config) {
		switch dbType {
		case constants.SQLite:
			if c.DB.SQLite != nil {
				c.DB.SQLite.MaxIdleConns = count
			}
		case constants.MySQL:
			if c.DB.MySQL != nil {
				c.DB.MySQL.MaxIdleConns = count
			}
		case constants.Postgres:
			if c.DB.Postgres != nil {
				c.DB.Postgres.MaxIdleConns = count
			}
		}
	}
}

func WithMaxLifetime(dbType string, duration time.Duration) ConfigOption {
	return func(c *Config) {
		switch dbType {
		case constants.SQLite:
			if c.DB.SQLite != nil {
				c.DB.SQLite.MaxLifetime = duration
			}
		case constants.MySQL:
			if c.DB.MySQL != nil {
				c.DB.MySQL.MaxLifetime = duration
			}
		case constants.Postgres:
			if c.DB.Postgres != nil {
				c.DB.Postgres.MaxLifetime = duration
			}
		}
	}
}

type RetryOption struct {
	env      *EnvConfig
	attempts int
	delay    time.Duration
}

func WithRetryPolicy(env *EnvConfig) *RetryOption {
	return &RetryOption{
		env:      env,
		attempts: constants.DefaultRetryAttempts, // 默認值
		delay:    constants.DefaultRetryDelay,    // 默認值
	}
}

// Builder
func (r *RetryOption) Build() ConfigOption {
	if attempts, ok := r.env.GetInt("DB_RETRY_ATTEMPTS"); ok {
		r.attempts = attempts
	}
	if delay, ok := r.env.GetDuration("DB_RETRY_DELAY"); ok {
		r.delay = delay
	}

	return func(c *Config) {
		c.DB.RetryPolicy = RetryPolicy{
			Attempts: r.attempts,
			Delay:    r.delay,
		}
	}
}

func NewConfig(opts ...ConfigOption) *Config {
	cfg := &Config{
		DB: &DBConfig{
			EnabledDatabases: make(map[string]bool),
			RetryPolicy: RetryPolicy{
				Attempts: constants.DefaultRetryAttempts,
				Delay:    constants.DefaultRetryDelay,
			},
		},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}
