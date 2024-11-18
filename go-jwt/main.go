package main

import (
	"log"

	"github.com/go-sql-driver/mysql"
)

func main() {
	cfg := mysql.Config{
		User:                 Envs.DBUser,
		Passwd:               Envs.DBPassword,
		Addr:                 Envs.DBAddress,
		DBName:               Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	sqlRepository := NewMySQLRespository(cfg)
	db, err := sqlRepository.Init()
	if err != nil {
		log.Fatal(err)
	}
	store := NewStore(db)
	api := NewApiServer(":8080", store)

	api.Serve()
}
