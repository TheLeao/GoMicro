package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // name of exchange
		"topic", // kind
		true, // is durable
		false, // is auto-deleted
		false, // is internal only (no, because microservices)
		false, // no wait
		nil, // no more args
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",
		false, // is durable
		false, // is auto-delete
		true, // is exclusive
		false, // no-wait
		nil,
	)
}