package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerWar(ch *amqp.Channel, gs *gamelogic.GameState) func(rw gamelogic.RecognitionOfWar) pubsub.AckType {
	defer fmt.Print("> ")
	return func(rw gamelogic.RecognitionOfWar) pubsub.AckType {
		warOutcome, winner, loser := gs.HandleWar(rw)
		message := ""
		log.Printf("war outcome was %v", warOutcome)
		switch warOutcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			message = fmt.Sprintf("%v won a war against %v", winner, loser)
		case gamelogic.WarOutcomeYouWon:
			message = fmt.Sprintf("%v won a war against %v", winner, loser)
		case gamelogic.WarOutcomeDraw:
			message = fmt.Sprintf("A war between %v and %v resulted in a draw", winner, loser)
		default:
			log.Println("No valid outcome of the war")
			return pubsub.NackDiscard
		}
		err := pubsub.PublishGameLog(ch, message, gs)
		if err != nil {
			return pubsub.NackRequeue
		} else {
			return pubsub.Ack
		}
	}
}
