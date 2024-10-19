package main

import (
	"fmt"
	"net/http"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var user User

type User struct {
	Name string
	ID   uuid.UUID
}

func saveUser(c echo.Context) error {
	user = User{
		ID:   uuid.New(),
		Name: gofakeit.Name(),
	}
	return c.JSON(http.StatusCreated, user)
}

func getUser(c echo.Context) error {
	return c.JSON(http.StatusCreated, user)
}

func main() {

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		str := fmt.Sprintf("Hello, %+v", User{Name: "Tester", ID: uuid.New()})
		return c.String(http.StatusOK, str)
	})

	e.POST("/users", saveUser)
	e.GET("/users/:id", getUser)
	// e.PUT("/users/:id", updateUser)
	// e.DELETE("/users/:id", deleteUser)

	e.Logger.Fatal(e.Start(":3333"))
}
