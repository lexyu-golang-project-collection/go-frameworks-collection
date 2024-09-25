package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

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
	CustomerID int32  `json:"customer_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ReferralID int32  `json:"referral_id"`
	Email      string `json:"email"`
}

type UserInfo struct {
	Id        int32  `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type JSONDates []time.Time

func (j JSONDates) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONDates) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &j)
}

type Combination struct {
	ID      int       `json:"id" gorm:"primaryKey"`
	SCodeID string    `json:"scode_id"`
	CCodeID string    `json:"ccode_id"`
	Dates   JSONDates `json:"dates" gorm:"type:json"` // 原生 JSON 類型
}

func main() {
	// MySQL
	password := "P@ssw0rd"
	dial := mysql.Open("root:" + password + "@tcp(127.0.0.1:3306)/staffs?charset=utf8&parseTime=True")
	DB, err := gorm.Open(dial, &gorm.Config{})

	if err != nil {
		panic("failed to connect to db")
	}
	fmt.Println("Connection Successful")

	// Migrate the schema
	// DB.AutoMigrate(&User{})
	// DB.Migrator().CreateTable(&Category{})
	// DB.Migrator().CreateTable(&Product{})
	// DB.Migrator().CreateTable(&Combination{})

	// ---------------------------------------------
	createFakeCombinations(DB, 10)

	// ---------------------------------------------
	// query := "SELECT * FROM userinfo WHERE id=?"
	// var userinfo UserInfo
	// DB.Raw(query, 1).Scan(&userinfo)
	// fmt.Println(userinfo.Id, userinfo.Firstname, userinfo.Lastname, userinfo.Username, userinfo.Password)

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

// 創建（Create）
func CreateCombination(db *gorm.DB, scode, ccode string, dates []time.Time) (*Combination, error) {
	combination := &Combination{
		SCodeID: scode,
		CCodeID: ccode,
		Dates:   JSONDates(dates),
	}
	result := db.Create(combination)
	if result.Error != nil {
		return nil, result.Error
	}
	return combination, nil
}

func createFakeCombinations(db *gorm.DB, count int) {
	for i := 0; i < count; i++ {
		combination := Combination{
			SCodeID: fmt.Sprintf("S%03d", i+1),
			CCodeID: fmt.Sprintf("C%03d", i+1),
			Dates: JSONDates{
				time.Now().AddDate(0, 0, i),
				time.Now().AddDate(0, 0, i+1),
				time.Now().AddDate(0, 0, i+2),
			},
		}
		db.Create(&combination)
	}
	fmt.Printf("Created %d fake combinations\n", count)
}

// 讀取（Read）
func GetCombination(db *gorm.DB, id int) (*Combination, error) {
	var combination Combination
	result := db.First(&combination, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &combination, nil
}

// 更新（Update）
func UpdateCombination(db *gorm.DB, id int, newDates []time.Time) (*Combination, error) {
	var combination Combination
	result := db.First(&combination, id)
	if result.Error != nil {
		return nil, result.Error
	}

	combination.Dates = JSONDates(newDates)
	result = db.Save(&combination)
	if result.Error != nil {
		return nil, result.Error
	}
	return &combination, nil
}

// 刪除（Delete）
func DeleteCombination(db *gorm.DB, id int) error {
	result := db.Delete(&Combination{}, id)
	return result.Error
}

// 查詢特定日期的組合
func GetCombinationsWithDate(db *gorm.DB, date time.Time) ([]Combination, error) {
	var combinations []Combination
	dateStr, _ := json.Marshal(date.Format("2006-01-02"))
	result := db.Where("dates @> ?", dateStr).Find(&combinations)
	if result.Error != nil {
		return nil, result.Error
	}
	return combinations, nil
}
