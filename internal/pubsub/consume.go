package pubsub

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Acktype int

const (
	Ack Acktype = iota
	NackRequeue
	NackDiscard
)

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
	handler func(T) Acktype,
) error {
	return subscribe(
		conn,
		exchange,
		queueName,
		key,
		simpleQueueType,
		handler,
		func(body []byte) (T, error) {
			var val T
			err := json.Unmarshal(body, &val)
			return val, err
		},
	)
}

func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) Acktype,
) error {
	return subscribe(
		conn,
		exchange,
		queueName,
		key,
		simpleQueueType,
		handler,
		func(body []byte) (T, error) {
			var val T
			err := gob.NewDecoder(bytes.NewBuffer(body)).Decode(&val)
			return val, err
		},
	)
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
		amqp.Table{
			"x-dead-letter-exchange": "peril_dlx",
		},
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

func subscribe[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) Acktype,
	unmarshaller func([]byte) (T, error),
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
			val, err := unmarshaller(delivery.Body)
			if err != nil {
				log.Println("Error unmarshalling delivery body:", err)
				continue
			}

			switch handler(val) {
			case Ack:
				delivery.Ack(false)
				log.Println("Ack")
			case NackDiscard:
				delivery.Nack(false, false)
				log.Println("NackDiscard")
			case NackRequeue:
				delivery.Nack(false, true)
				log.Println("NackRequeue")
			default:
				log.Println("Handler returned invalid Acktype")
			}
		}
	}()

	return nil
}
