package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	mysqlDB      *gorm.DB
	postgresDB   *gorm.DB
	mysqlOnce    sync.Once
	postgresOnce sync.Once
)

func ConnectToDatabases() error {
	var mysqlErr, pgErr error

	mysqlOnce.Do(func() {
		dsn := os.Getenv("MYSQL_DB_URL")
		if dsn == "" {
			log.Println("MySQL DSN 未找到，使用預設值")
			dsn = "user:password@tcp(localhost:3306)/bookstore?charset=utf8mb4&parseTime=True&loc=Local"
		}

		mysqlDB, mysqlErr = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if mysqlErr == nil {
			log.Println("已連接到 MySQL 資料庫")
		}
	})

	postgresOnce.Do(func() {
		dsn := os.Getenv("PG_DB_URL")
		if dsn == "" {
			log.Println("PostgreSQL DSN 未找到，使用預設值")
			dsn = "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
		}

		postgresDB, pgErr = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if pgErr == nil {
			log.Println("已連接到 PostgreSQL 資料庫")
		}
	})

	if mysqlErr != nil {
		return fmt.Errorf("無法連接到 MySQL: %w", mysqlErr)
	}

	if pgErr != nil {
		return fmt.Errorf("無法連接到 PostgreSQL: %w", pgErr)
	}

	return nil
}

func GetMySQLDB() *gorm.DB {
	if mysqlDB == nil {
		if err := ConnectToDatabases(); err != nil {
			log.Printf("警告: 獲取資料庫連接時出錯: %v", err)
		}
	}
	return mysqlDB
}

func GetPostgresDB() *gorm.DB {
	if postgresDB == nil {
		if err := ConnectToDatabases(); err != nil {
			log.Printf("警告: 獲取資料庫連接時出錯: %v", err)
		}
	}
	return postgresDB
}

func CloseConnections() {
	if mysqlDB != nil {
		sqlDB, err := mysqlDB.DB()
		if err == nil {
			sqlDB.Close()
			log.Println("MySQL 連接已關閉")
		}
	}

	if postgresDB != nil {
		sqlDB, err := postgresDB.DB()
		if err == nil {
			sqlDB.Close()
			log.Println("PostgreSQL 連接已關閉")
		}
	}
}
