package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	// amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// conn, err := amqp.Dial("amqp:localhost:5672")
	// if err != nil {
	// 	panic(err)
	// }
	// defer conn.Close()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
