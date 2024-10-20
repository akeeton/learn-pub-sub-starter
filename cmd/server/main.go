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

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatal("Error creating channel:", err)
	}

	fmt.Println("Peril game server connected to RabbitMQ")
	gamelogic.PrintServerHelp()

	for {
		input := gamelogic.GetInput()

		if len(input) < 1 {
			continue
		}

		switch input[0] {
		case "pause":
			sendPauseMessage(publishCh, true)
		case "resume":
			sendPauseMessage(publishCh, false)
		case "quit":
			fmt.Println("Exiting server")
			return
		default:
			fmt.Println("Unknown command")
		}
	}
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
