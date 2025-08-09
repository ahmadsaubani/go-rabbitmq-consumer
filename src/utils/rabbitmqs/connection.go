package rabbitmqs

import (
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conn    *amqp.Connection
	channel *amqp.Channel
	once    sync.Once
	initErr error
)

func InitRabbitMQ() error {
	once.Do(func() {
		var err error
		conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			initErr = fmt.Errorf("failed to connect to RabbitMQ: %w", err)
			return
		}

		channel, err = conn.Channel()
		if err != nil {
			initErr = fmt.Errorf("failed to open a channel: %w", err)
			return
		}
	})
	return initErr
}

func GetChannel() (*amqp.Channel, error) {
	if channel == nil {
		return nil, fmt.Errorf("RabbitMQ channel is not initialized")
	}
	return channel, nil
}

func GetConnection() (*amqp.Connection, error) {
	if conn == nil {
		return nil, fmt.Errorf("RabbitMQ connection is not initialized")
	}
	return conn, nil
}
