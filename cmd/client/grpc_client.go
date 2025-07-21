package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"pickup-srv/proto"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewPickupServiceClient(conn)

	log.Println("=== Getting all users with limit 1 ===")
	resp, err := client.GetUsers(context.Background(), &proto.GetUsersRequest{
		UserSearchParams: &proto.UserSearchParams{
			Id:            33,
			Gender:        "m",
			SearchGender:  "f",
			Age:           25,
			SearchAgeFrom: 18,
			SearchAgeTo:   30,
			Location:      32,
		},
		Limit: 1,
	})
	if err != nil {
		log.Fatalf("GetUsers failed: %v", err)
	}

	log.Printf("Found %d users", resp.Total)
	for _, user := range resp.Users {
		log.Printf("User: %s - Age: %d, City: %s", user.Name, user.Age, user.City)
	}

	log.Println("=== Getting all users with limit 3 ===")
	resp, err = client.GetUsers(context.Background(), &proto.GetUsersRequest{
		UserSearchParams: &proto.UserSearchParams{
			Id:            109,
			Gender:        "f",
			SearchGender:  "m",
			Age:           20,
			SearchAgeFrom: 18,
			SearchAgeTo:   40,
			Location:      32,
		},
		Limit: 3,
	})
	if err != nil {
		log.Fatalf("GetUsers failed: %v", err)
	}

	log.Printf("Found %d users", resp.Total)
	for _, user := range resp.Users {
		log.Printf("User: %s - Age: %d, City: %s", user.Name, user.Age, user.City)
	}

	log.Println("=== Getting all users with limit 100 ===")
	resp, err = client.GetUsers(context.Background(), &proto.GetUsersRequest{
		UserSearchParams: &proto.UserSearchParams{
			Id:            109,
			Gender:        "f",
			SearchGender:  "m",
			Age:           20,
			SearchAgeFrom: 18,
			SearchAgeTo:   40,
			Location:      32,
		},
		Limit: 10,
	})
	if err != nil {
		log.Fatalf("GetUsers failed: %v", err)
	}

	log.Printf("Found %d users", resp.Total)
	for _, user := range resp.Users {
		log.Printf("User: %s - Age: %d, City: %s", user.Name, user.Age, user.City)
	}
}
