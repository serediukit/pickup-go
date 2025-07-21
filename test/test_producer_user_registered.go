package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"pickup-srv/internal/messaging"
	"pickup-srv/internal/models"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitmq, err := messaging.NewRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitmq.Close()

	queueName := "user.registration"
	if err := rabbitmq.DeclareQueue(queueName); err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Sample user registration events
	events := []models.UserRegistrationEvent{
		{
			Name:          "Mike Johnson",
			Age:           28,
			City:          "Seattle",
			Gender:        "m",
			SearchGender:  "f",
			SearchAgeFrom: 20,
			SearchAgeTo:   30,
			Location:      32,
		},
		{
			Name:          "Sarah Connor",
			Age:           35,
			City:          "Los Angeles",
			Gender:        "f",
			SearchGender:  "m",
			SearchAgeFrom: 30,
			SearchAgeTo:   40,
			Location:      12,
		},
		{
			Name:          "David Smith",
			Age:           42,
			City:          "Boston",
			Gender:        "m",
			SearchGender:  "m",
			SearchAgeFrom: 18,
			SearchAgeTo:   20,
			Location:      17,
		},
	}

	for i := 0; i < 100; i++ {
		randEventIndex := rand.Intn(len(events))
		event := events[randEventIndex]
		if err := publishEvent(rabbitmq, queueName, event); err != nil {
			log.Printf("Failed to publish event: %v", err)
		} else {
			log.Printf("Published registration event for user: %s", event.Name)
		}
		time.Sleep(10 * time.Second)
	}

	log.Println("All events published successfully!")
}

func publishEvent(rabbitmq *messaging.RabbitMQ, queueName string, event models.UserRegistrationEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return rabbitmq.GetChannel().PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
