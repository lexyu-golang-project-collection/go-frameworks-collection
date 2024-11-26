package main

import amqp "github.com/rabbitmq/amqp091-go"

func main() {
	conn, err := amqp.Dial("amqp:localhost:5672")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
}
