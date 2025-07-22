package messaging

import (
	"fmt"
	"time"

	"pickup-srv/util"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ() (*RabbitMQ, error) {
	host := util.GetEnv("RABBITMQ_HOST", "localhost")
	port := util.GetEnv("RABBITMQ_PORT", "5672")
	user := util.GetEnv("RABBITMQ_USER", "guest")
	password := util.GetEnv("RABBITMQ_PASSWORD", "guest")

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)

	var conn *amqp.Connection
	var err error

	retriesCount := 10

	for i := 0; i < retriesCount; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}

		fmt.Printf("Failed to connect to RabbitMQ (attempt %d/%d): %v\n", i+1, retriesCount, err)
		time.Sleep(time.Duration(i+1) * 3 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ after 5 attempts: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
	}, nil
}

func (r *RabbitMQ) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

func (r *RabbitMQ) DeclareQueue(queueName string) error {
	_, err := r.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}

func (r *RabbitMQ) ConsumeMessages(queueName string) (<-chan amqp.Delivery, error) {
	msgs, err := r.channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	return msgs, err
}

func (r *RabbitMQ) GetChannel() *amqp.Channel {
	return r.channel
}
