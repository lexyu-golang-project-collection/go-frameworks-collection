package models

import (
	"database/sql"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lex/go-crud-api-mysql/pkg/config"
)

var db *gorm.DB

type Book struct {
	gorm.Model
	Name        string `gorm:""json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
}

type Others struct {
	Name        string `json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
}

type User struct {
	ID           uint
	Name         string
	Email        *string
	Age          uint8
	Birthday     *time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type TestTable struct {
	ID         int64
	Username   string `gorm:"column:username"`
	Password   string `gorm:"column:password"`
	CreateTime int64  `gorm:"column:createtime"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Book{})
	db.AutoMigrate(&Others{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&TestTable{})
}

func (b *Book) CreateBook() *Book {
	db.NewRecord(b)
	db.Create(&b)
	return b
}

func GetAllBooks() []Book {
	var Books []Book
	db.Find(&Books)
	return Books
}

func GetBookById(Id int64) (*Book, *gorm.DB) {
	var getBook Book
	db := db.Where("ID=?", Id).Find(&getBook)
	return &getBook, db
}

func DeleteBook(ID int64) Book {
	var book Book
	db.Where("ID=?", ID).Delete(book)
	return book
}
