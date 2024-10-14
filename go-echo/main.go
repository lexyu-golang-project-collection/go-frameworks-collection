package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type User struct {
	Name string
	ID   uuid.UUID
}

func main() {

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		str := fmt.Sprintf("Hello, %+v", User{Name: "Tester", ID: uuid.New()})
		return c.String(http.StatusOK, str)
	})

	e.Logger.Fatal(e.Start(":3333"))
}
