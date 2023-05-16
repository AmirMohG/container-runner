package services

import "github.com/streadway/amqp"

type AMQPConnection interface {
	Channel() (*amqp.Channel, error)
	Close() error
}
