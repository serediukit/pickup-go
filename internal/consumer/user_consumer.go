package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"pickup-srv/internal/messaging"
	"pickup-srv/internal/models"
	"pickup-srv/internal/repository"
)

type UserConsumer struct {
	rabbitmq *messaging.RabbitMQ
	userRepo *repository.UserRepository
}

func NewUserConsumer(rabbitmq *messaging.RabbitMQ, userRepo *repository.UserRepository) *UserConsumer {
	return &UserConsumer{
		rabbitmq: rabbitmq,
		userRepo: userRepo,
	}
}

func (c *UserConsumer) Start(ctx context.Context) error {
	queueName := "user.registration"

	if err := c.rabbitmq.DeclareQueue(queueName); err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	msgs, err := c.rabbitmq.ConsumeMessages(queueName)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("Starting to consume messages from queue: %s", queueName)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Consumer context cancelled, stopping...")
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Println("Message channel closed")
					return
				}

				if err := c.processMessage(msg); err != nil {
					log.Printf("Error processing message: %v", err)
					msg.Nack(false, true)
				} else {
					msg.Ack(false)
				}
			}
		}
	}()

	return nil
}

func (c *UserConsumer) processMessage(msg amqp.Delivery) error {
	log.Printf("Received message: %s", string(msg.Body))

	var userEvent models.UserRegistrationEvent
	if err := json.Unmarshal(msg.Body, &userEvent); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	if err := c.validateEvent(&userEvent); err != nil {
		return fmt.Errorf("invalid event data: %w", err)
	}

	if err := c.userRepo.CreateUser(&userEvent); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	log.Printf("Successfully created user: %s", userEvent.Name)
	return nil
}

func (c *UserConsumer) validateEvent(event *models.UserRegistrationEvent) error {
	if event.Name == "" {
		return fmt.Errorf("name is required")
	}
	if event.Age <= 0 {
		return fmt.Errorf("age must be positive")
	}
	if event.City == "" {
		return fmt.Errorf("city is required")
	}
	if event.Gender != "m" && event.Gender != "f" {
		return fmt.Errorf("invalid gender")
	}
	if event.SearchGender != "m" && event.SearchGender != "f" {
		return fmt.Errorf("invalid search gender")
	}
	if event.SearchAgeFrom <= 0 || event.SearchAgeTo < event.SearchAgeFrom {
		return fmt.Errorf("invalid search age range")
	}
	return nil
}
