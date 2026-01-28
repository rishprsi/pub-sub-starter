package pubsub

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int

const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	jsonVal, err := json.Marshal(val)
	if err != nil {
		return err
	}

	publishing := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonVal,
	}
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, publishing)
	if err != nil {
		return err
	}
	return nil
}
