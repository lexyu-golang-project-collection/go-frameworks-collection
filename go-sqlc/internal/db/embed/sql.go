package embedSQL

import (
	"embed"
)

//go:embed sqlc/*/*.sql
var SQLFiles embed.FS

//go:embed sqlc/sqlite/schema.sql
var SQLiteDDL string

func GetSQLiteDDL() string {
	return SQLiteDDL
}
