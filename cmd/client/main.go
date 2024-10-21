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
		log.Fatal("Error connecting to RabbitMQ:", err)
	}
	defer conn.Close()
	fmt.Println("Peril game client connected to RabbitMQ")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatal("Error welcoming client:", err)
	}

	_, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		fmt.Sprintf("%s.%s", routing.PauseKey, username),
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
	)
	if err != nil {
		log.Fatal("Error declaring and binding queue:", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	gameState := gamelogic.NewGameState(username)

	// Infinite REPL loop
	for {
		words := gamelogic.GetInput()

		if len(words) < 1 {
			continue
		}

		switch words[0] {
		case "spawn":
			handleSpawn(gameState, words)
		case "move":
			handleMove(gameState, words)
		case "status":
			handleStatus(gameState)
		case "help":
			handleHelp()
		case "spam":
			handleSpam()
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("Unknown command")
		}
	}
}

func handleSpawn(gameState *gamelogic.GameState, words []string) {
	err := gameState.CommandSpawn(words)
	if err != nil {
		fmt.Println("Error spawning unit:", err)
	}
}

func handleMove(gameState *gamelogic.GameState, words []string) {
	_, err := gameState.CommandMove(words)
	if err != nil {
		fmt.Println("Error moving:", err)
	}
}

func handleStatus(gameState *gamelogic.GameState) {
	gameState.CommandStatus()
}

func handleHelp() {
	gamelogic.PrintClientHelp()
}

func handleSpam() {
	fmt.Println("Spamming not allowed yet!")
}
