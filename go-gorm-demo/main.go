package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Category struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Product struct {
	gorm.Model
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Customer struct {
	customer_id int32  `json:"customer_id"`
	first_name  string `json:"first_name"`
	last_name   string `json:"last_name"`
	referral_id int32  `json:"referral_id"`
	email       string `json:"email"`
}

type UserInfo struct {
	Id        int32  `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func main() {
	// MySQL
	dsn := "root:00112125@tcp(127.0.0.1:3306)/staffs?charset=utf8&parseTime=True&loc=Local"
	dial := mysql.Open(dsn)
	DB, err := gorm.Open(dial, &gorm.Config{})

	if err != nil {
		panic("failed to connect to db")
	}
	fmt.Println("Connection Successful")

	// Migrate the schema
	// DB.AutoMigrate(&User{})
	// DB.Migrator().CreateTable(&Category{})
	// DB.Migrator().CreateTable(&Product{})

	query := "SELECT * FROM userinfo WHERE id=?"
	var userinfo UserInfo
	DB.Raw(query, 1).Scan(&userinfo)
	fmt.Println(userinfo.Id, userinfo.Firstname, userinfo.Lastname, userinfo.Username, userinfo.Password)

	// for _, userinfo := range userinfos {
	// 	fmt.Println(userinfo.id, userinfo.firstname, userinfo.lastname, userinfo.username, userinfo.password)
	// }

	// ----------------------------------------------

	// MSSQL
	/*
			sql_db, err := sql.Open("sqlserver", "sqlserver://sa:p@ssw0rd@127.0.0.1/instance?database")
			err = sql_db.Ping()
			if err != nil {
				panic(err)
			}

		query := "select * from customer"
		rows, err := sql_db.Query(query)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		customers := []Customer{}

		for rows.Next() {
			customer := Customer{}

			rows.Scan(&customer.customer_id, &customer.first_name,
				&customer.last_name, &customer.referral_id, &customer.email)
			// fmt.Printf("customer_id=%v, first_name=%s, last_name=%s, referral_id=%v, email=%s ",
			// 	customer_id, first_name, last_name, referral_id, email)
			customers = append(customers, customer)
		}

		fmt.Printf("%#v", customers)
	*/
}
