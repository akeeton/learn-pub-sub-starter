package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatal("Error connecting to RabbitMQ:", err)
	}
	defer conn.Close()

	_, logsQueue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		"game_logs",
		"game_logs.*",
		pubsub.SimpleQueueDurable,
	)
	if err != nil {
		log.Fatal("Error declaring and binding queue:", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", logsQueue.Name)

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatal("Error creating channel:", err)
	}

	fmt.Println("Peril game server connected to RabbitMQ")
	gamelogic.PrintServerHelp()

	// Infinite REPL loop
	for {
		words := gamelogic.GetInput()

		if len(words) < 1 {
			continue
		}

		switch words[0] {
		case "pause":
			handlePause(publishCh)
		case "resume":
			handleResume(publishCh)
		case "quit":
			fallthrough
		case "q":
			fmt.Println("Exiting server")
			return
		default:
			fmt.Println("Unknown command")
		}
	}
}

func handlePause(publishCh *amqp.Channel) {
	sendPauseMessage(publishCh, true)
}

func handleResume(publishCh *amqp.Channel) {
	sendPauseMessage(publishCh, false)
}

func sendPauseMessage(publishCh *amqp.Channel, isPaused bool) {
	err := pubsub.PublishJSON(
		publishCh,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{
			IsPaused: isPaused,
		},
	)
	if err != nil {
		log.Println("Error publishing JSON to channel:", err)
	}

	kind := "Resume"
	if isPaused {
		kind = "Pause"
	}
	fmt.Printf("%s message sent!\n", kind)
}
