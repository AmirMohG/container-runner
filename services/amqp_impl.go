package services

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

type AMQPConnectionImpl struct {
	conn *amqp.Connection
}

func NewAMQPConnection() (*AMQPConnectionImpl, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://" + os.Getenv("MQ_USER") + ":" + os.Getenv("MQ_PASS") + "@" + os.Getenv("MQ_ADDRESS") + ":" + os.Getenv("MQ_PORT")))
	if err != nil {
		return nil, err
	}
	return &AMQPConnectionImpl{conn}, nil
}

func (a *AMQPConnectionImpl) Channel() (*amqp.Channel, error) {
	return a.conn.Channel()
}

func (a *AMQPConnectionImpl) Close() error {
	return a.conn.Close()
}

func (a *AMQPConnectionImpl) CreateQueue(name string) error {
	ch, err := a.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func (a *AMQPConnectionImpl) CreateConsumer(channel *amqp.Channel, queueName string) (<-chan amqp.Delivery, error) {
	msgs, err := channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
