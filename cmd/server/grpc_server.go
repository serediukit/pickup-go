package main

import (
	"log"
	"net"
	"os"
	"pickup-srv/internal/database"
	"pickup-srv/internal/repository"
	"pickup-srv/internal/service"
	"pickup-srv/proto"

	"google.golang.org/grpc"
)

func main() {
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)

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

	log.Printf("Starting pickup-srv on port %s", port)
	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
