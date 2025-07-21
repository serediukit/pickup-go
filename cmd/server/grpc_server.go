package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"pickup-srv/internal/consumer"
	"pickup-srv/internal/database"
	"pickup-srv/internal/messaging"
	"pickup-srv/internal/repository"
	"pickup-srv/internal/service"
	"pickup-srv/proto"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)

	rabbitmq, err := messaging.NewRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitmq.Close()

	userConsumer := consumer.NewUserConsumer(rabbitmq, userRepo)
	if err := userConsumer.Start(ctx); err != nil {
		log.Fatalf("Failed to start user consumer: %v", err)
	}

	pickupService := service.NewPickupService(userRepo)

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "8080"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterPickupServiceServer(s, pickupService)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down...")
		cancel()
		s.GracefulStop()
	}()

	log.Println("RabbitMQ consumer started for user registration events")

	log.Printf("Starting pickup-srv on port %s", port)
	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
