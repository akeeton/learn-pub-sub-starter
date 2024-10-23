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
	fmt.Println("Starting Peril client...")

	const connectionString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatalln("Error connecting to RabbitMQ:", err)
	}
	defer conn.Close()
	fmt.Println("Peril game client connected to RabbitMQ")

	err = pubsub.DeclareAndBindCommonExchangesAndQueues(conn)
	if err != nil {
		log.Fatalln("Error declaring common exchanges and queues:", err)
	}

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalln("Error creating channel:", err)
	}

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalln("Error welcoming client:", err)
	}

	gs := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+gs.GetUsername(),
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
		handlerPause(gs),
	)
	if err != nil {
		log.Fatalln("Error subscribing to pause queue:", err)
	}

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+gs.GetUsername(),
		routing.ArmyMovesPrefix+".*",
		pubsub.SimpleQueueTransient,
		handlerMove(gs, publishCh),
	)
	if err != nil {
		log.Fatalln("Error subscribing to army_move queue:", err)
	}

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.WarRecognitionsPrefix,
		routing.WarRecognitionsPrefix+".*",
		pubsub.SimpleQueueDurable,
		handlerWar(gs, publishCh),
	)
	if err != nil {
		log.Fatalln("Error subscribing to war queue:", err)
	}

	// Infinite REPL loop
	for {
		words := gamelogic.GetInput()

		if len(words) < 1 {
			continue
		}

		switch words[0] {
		case "spawn":
			handleSpawn(gs, words)
		case "move":
			handleMove(gs, publishCh, words)
		case "status":
			handleStatus(gs)
		case "help":
			handleHelp()
		case "spam":
			handleSpam()
		case "quit":
			fallthrough
		case "q":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("Unknown command")
		}
	}
}

func handleSpawn(gs *gamelogic.GameState, words []string) {
	err := gs.CommandSpawn(words)
	if err != nil {
		fmt.Println("Error spawning unit:", err)
	}
}

func handleMove(gs *gamelogic.GameState, publishCh *amqp.Channel, words []string) {
	mv, err := gs.CommandMove(words)
	if err != nil {
		fmt.Println("Error moving:", err)
	}

	err = pubsub.PublishJSON(
		publishCh,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+gs.GetUsername(),
		mv,
	)
	if err != nil {
		log.Println("Error publishing move message:", err)
	}

	fmt.Println("Move published successfully")
}

func handleStatus(gs *gamelogic.GameState) {
	gs.CommandStatus()
}

func handleHelp() {
	gamelogic.PrintClientHelp()
}

func handleSpam() {
	fmt.Println("Spamming not allowed yet!")
}
