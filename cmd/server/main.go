package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const connectionString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatal("Error connecting to RabbitMQ:", err)
	}
	defer conn.Close()

	fmt.Println("Connected to RabbitMQ")

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatal("Error creating RabbitMQ channel:", err)
	}

	err = pubsub.PublishJSON(
		publishCh,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{
			IsPaused: true,
		},
	)
	if err != nil {
		log.Println("Error publishing JSON to channel:", err)
	}

	fmt.Println("Pause message sent!")

	// Wait for Ctrl+C
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh
	fmt.Println("Received Ctrl+C, shutting down connection to RabbitMQ")
	conn.Close()
}
