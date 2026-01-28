package main

import (
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerMove(gs *gamelogic.GameState, publishCh *amqp.Channel) func(move gamelogic.ArmyMove) pubsub.AckType {
	return func(move gamelogic.ArmyMove) pubsub.AckType {
		outcome := gs.HandleMove(move)
		log.Printf("Move outcome was %v", outcome)
		defer print("> ")
		if outcome == gamelogic.MoveOutcomeMakeWar {
			log.Println("The war outcome got triggered")
			key := routing.WarRecognitionsPrefix + "." + gs.GetUsername()
			message := gamelogic.RecognitionOfWar{
				Attacker: move.Player,
				Defender: gs.GetPlayerSnap(),
			}
			err := pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilTopic,
				key,
				message,
			)
			if err != nil {
				log.Printf("Failed to publish war declaration: %v", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		}
		if outcome == gamelogic.MoveOutComeSafe {
			return pubsub.Ack
		} else {
			return pubsub.NackDiscard
		}
	}
}
