package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const connectionString = "amqp://guest:guest@localhost:5672/"

	connection, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatal("Error connecting to RabbitMQ:", err)
	}
	defer connection.Close()

	fmt.Println("Connected to RabbitMQ")

	// Wait for Ctrl+C
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh
	fmt.Println("Received Ctrl+C, shutting down connection to RabbitMQ")
	connection.Close()
}
