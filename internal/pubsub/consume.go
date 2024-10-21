package pubsub

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Acktype int

type SimpleQueueType int

const (
	SimpleQueueDurable SimpleQueueType = iota
	SimpleQueueTransient
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T),
) error {
	ch, queue, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return fmt.Errorf("error declaring and binding to queue: %w", err)
	}

	deliveryCh, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("error consuming queue: %w", err)
	}

	go func() {
		defer ch.Close()

		for delivery := range deliveryCh {
			var val T
			if err := json.Unmarshal(delivery.Body, &val); err != nil {
				log.Println("Error unmarshalling delivery body:", err)
				continue
			}

			handler(val)

			if err := delivery.Ack(false); err != nil {
				log.Println("Error acking delivery:", err)
				continue
			}
		}
	}()

	return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("error creating channel: %w", err)
	}

	durable := simpleQueueType == SimpleQueueDurable
	autoDelete := simpleQueueType == SimpleQueueTransient
	exclusive := simpleQueueType == SimpleQueueTransient
	noWait := false
	queue, err := ch.QueueDeclare(
		queueName,
		durable,
		autoDelete,
		exclusive,
		noWait,
		nil,
	)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("error declaring queue: %w", err)
	}

	err = ch.QueueBind(queue.Name, key, exchange, noWait, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("error binding queue: %w", err)
	}

	return ch, queue, nil
}
