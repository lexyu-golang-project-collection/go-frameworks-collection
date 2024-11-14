package embedSQL

import (
	"embed"
	"fmt"
)

//go:embed sqlc/*/*.sql
var SQLFiles embed.FS

//go:embed sqlc/sqlite/schema.sql
var SQLiteDDL string

//go:embed sqlc/mysql/schema.sql
var MySQLDDL string

//go:embed sqlc/postgres/schema.sql
var PostgresDDL string

// 導出獲取 DDL 的函數
func GetSQLiteDDL() string {
	return SQLiteDDL
}

func GetMySQLDDL() string {
	return MySQLDDL
}

func GetPostgresDDL() string {
	return PostgresDDL
}

// 獲取特定的 SQL 文件內容
func GetSQL(dbType, filename string) (string, error) {
	content, err := SQLFiles.ReadFile(fmt.Sprintf("schema/%s/%s", dbType, filename))
	if err != nil {
		return "", err
	}
	return string(content), nil
}
