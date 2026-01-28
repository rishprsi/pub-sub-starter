package pubsub

import (
	"log"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGameLog(ch *amqp.Channel, val string, gs *gamelogic.GameState) error {
	username := gs.GetUsername()
	logStruct := routing.GameLog{
		CurrentTime: time.Now(),
		Message:     val,
		Username:    username,
	}

	log.Printf("Entering log: %v from user %v", val, username)
	key := routing.GameLogSlug + "." + username
	exchange := routing.ExchangePerilTopic

	err := PublishGob(ch, exchange, key, logStruct)
	if err != nil {
		return err
	}
	return nil
}
