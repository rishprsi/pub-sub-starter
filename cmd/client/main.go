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
		log.Fatal("Failed to establish connection with the queue")
	}
	defer conn.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatal("Failed to load client name")
	}

	gameState := gamelogic.NewGameState(username)

	pauseQueueName := fmt.Sprintf("pause.%s", username)
	pauseHandler := handlerPause(gameState)
	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, pauseQueueName, routing.PauseKey, pubsub.SimpleQueueTransient, pauseHandler)
	if err != nil {
		log.Printf("Failed to create pause channel and queue with the error: %s", err)
	}

	armyMoveQueueName := fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, username)
	// pauseHandler := handlerPause(gameState)
	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, armyMoveQueueName, routing.ArmyMovesPrefix+".*", pubsub.SimpleQueueTransient, handlerMove(gameState))
	if err != nil {
		log.Printf("Failed to create pause channel and queue with the error: %s", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal("Couldn't create the publish channel")
	}

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "spawn":
			err := gameState.CommandSpawn(words)
			if err != nil {
				log.Printf("Failed to spawn unit with err: %v", err)
			}
		case "move":
			armyMove, err := gameState.CommandMove(words)
			if err != nil {
				log.Printf("Could not move unit: %v", err)
			}
			err = pubsub.PublishJSON(
				channel,
				routing.ExchangePerilTopic,
				armyMoveQueueName,
				armyMove,
			)
			if err != nil {
				log.Printf("Failed to publish army move command: %v\n", err)
			}

			log.Printf("Move completed")
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			log.Printf("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("unknown command")
		}
	}

	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// <-signalChan
	// log.Println("Received the interrupt signal closing connection")
}
