package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var buffer bytes.Buffer
	err := gob.NewEncoder(&buffer).Encode(val)
	if err != nil {
		return err
	}

	rawVal := buffer.Bytes()
	publishMessage := amqp.Publishing{
		ContentType: "application/gob",
		Body:        rawVal,
	}
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, publishMessage)
	if err != nil {
		return err
	}

	return nil
}
