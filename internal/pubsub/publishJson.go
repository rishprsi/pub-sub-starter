package pubsub

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
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

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T),
) error {
	channel, queue, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}
	log.Printf("Consuming queue %s", queue.Name)

	deliveryChan, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	routine := func() {
		log.Println("Running the go routine for consuming queue")
		for chanMessage := range deliveryChan {
			var message T
			err := json.Unmarshal(chanMessage.Body, &message)
			if err != nil {
				log.Printf("Failed to extract message from queue: %v", err)
			}
			handler(message)
			err = chanMessage.Ack(false)
			if err != nil {
				log.Printf("Failed to ack for the required message: %v", err)
			}
		}
	}

	go routine()

	return nil
}
