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
	const connectionString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatal("Failed to connect to the rabbitmq server ", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to create a channel for pub sub")
	}
	defer channel.Close()

	// topicChannel, _, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, "game_logs", "game_logs.*", pubsub.SimpleQueueDurable)
	// if err != nil {
	// 	log.Fatal("Could not create the topic channel and queue: ", err)
	// }
	// defer topicChannel.Close()

	log.Println("Successfully connected to the Rabbitmq server")

	gamelogic.PrintServerHelp()

	pubsub.SubscribeGob(conn, routing.ExchangePerilTopic, "game_logs", "game_logs.*", pubsub.SimpleQueueDurable, handlerLogs())

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "pause":
			fmt.Println("Publishing paused game state")
			err = pubsub.PublishJSON(
				channel,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: true,
				},
			)
			if err != nil {
				log.Printf("could not publish time: %v", err)
			}
		case "resume":
			fmt.Println("Publishing resumes game state")
			err = pubsub.PublishJSON(
				channel,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: false,
				},
			)
			if err != nil {
				log.Printf("could not publish time: %v", err)
			}
		case "quit":
			log.Println("goodbye")
			return
		default:
			fmt.Println("unknown command")
		}
	}
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// <-signalChan
	// log.Println("Recieved interrupt closing connection and closing the server")
}
