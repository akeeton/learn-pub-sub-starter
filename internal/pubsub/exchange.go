package pubsub

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclarePerilExchanges(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("error creating channel: %w", err)
	}
	defer ch.Close()

	err = declareExchange(ch, routing.ExchangePerilDirect, amqp.ExchangeDirect)
	if err != nil {
		return err
	}

	err = declareExchange(ch, routing.ExchangePerilTopic, amqp.ExchangeTopic)
	if err != nil {
		return err
	}

	return nil
}

func declareExchange(ch *amqp.Channel, name, kind string) error {
	durable := true
	autoDelete := false
	internal := false
	noWait := false
	err := ch.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, nil)
	if err != nil {
		return fmt.Errorf("error declaring '%s' exchange: %w", name, err)
	}

	return nil
}
