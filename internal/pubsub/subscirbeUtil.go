package pubsub

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) AckType,
) error {
	unmarshaller := func(message []byte) (T, error) {
		var decodedMessage T
		err := json.Unmarshal(message, &decodedMessage)
		if err != nil {
			return decodedMessage, err
		}
		return decodedMessage, nil
	}
	err := subscribe(conn, exchange, queueName, key, queueType, handler, unmarshaller)
	if err != nil {
		return err
	}
	return nil
}

func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) AckType,
) error {
	unmarshaller := func(message []byte) (T, error) {
		var decodedMessage T
		buffer := bytes.NewBuffer(message)

		err := gob.NewDecoder(buffer).Decode(&decodedMessage)
		if err != nil {
			return decodedMessage, err
		}
		return decodedMessage, nil
	}
	err := subscribe(conn, exchange, queueName, key, queueType, handler, unmarshaller)
	if err != nil {
		return err
	}
	return nil
}

func subscribe[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) AckType,
	unmarshaller func([]byte) (T, error),
) error {
	channel, queue, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return err
	}
	err = channel.Qos(10, 0, true)
	if err != nil {
		log.Printf("Unable to set the prefetch: %v", err)
	}
	deliveryChan, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	routine := func() {
		log.Println("Running the go routine  for consuming queue")
		for chanMessage := range deliveryChan {
			message, err := unmarshaller(chanMessage.Body)
			if err != nil {
				log.Printf("Failed to extract message from queue: %v", err)
			}
			ackType := handler(message)

			switch ackType {
			case Ack:
				err = chanMessage.Ack(false)
				log.Println("Performing Ack")
			case NackRequeue:
				err = chanMessage.Nack(false, true)
				log.Println("Performing Nack requeue")
			case NackDiscard:
				err = chanMessage.Nack(false, false)
				log.Println("Performing Nack Discard")
			}
			if err != nil {
				log.Printf("Failed to send the required ack to the channel: %v", err)
			}
		}
	}
	go routine()
	return nil
}
