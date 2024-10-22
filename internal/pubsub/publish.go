package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	body, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	return publishBytes(ch, exchange, key, "application/json", body)
}

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(val)
	if err != nil {
		return fmt.Errorf("error encoding gob: %w", err)
	}

	return publishBytes(ch, exchange, key, "application/gob", buffer.Bytes())
}

func publishBytes(ch *amqp.Channel, exchange, key, contentType string, body []byte) error {
	err := ch.PublishWithContext(
		context.Background(),
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("error publishing to channel: %w", err)
	}

	return nil
}
