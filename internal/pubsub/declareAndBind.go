package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	SimpleQueueDurable SimpleQueueType = iota
	SimpleQueueTransient
)

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType) (*amqp.Channel, amqp.Queue, error) {
	channel, err := conn.Channel()
	if err != nil {
		return channel, amqp.Queue{}, err
	}

	isTransient := queueType == SimpleQueueTransient
	queue, err := channel.QueueDeclare(queueName, !isTransient, isTransient, isTransient, false, nil)
	if err != nil {
		return channel, queue, err
	}

	channel.QueueBind(queue.Name, key, exchange, false, nil)

	return channel, queue, nil
}
