package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	api "buffup/GolangTechTask/api"
)

func createVoteable(client api.VotingServiceClient) {
	fmt.Println("createVoteable")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.CreateVoteable(ctx, &api.CreateVoteableRequest{
		Question: "Q1", 
		Answers: []string{"Q1A1", "Q1A2"},
	})
	if err != nil {
		log.Fatalf("fatal: %v", err)
	}
	log.Printf("createVoteableResp: %s", r)
}

func listVoteables(client api.VotingServiceClient) {
	fmt.Println("listVoteables")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.ListVoteables(ctx, &api.ListVoteableRequest{Page: 1, Size: 5})

	if err != nil {
		log.Fatalf("fatal: %v", err)
	}

	for i, c := range r.Votables {
        fmt.Println(i, c)
    }
}

func main() {
	fmt.Println("Client Start")

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := api.NewVotingServiceClient(conn)

	createVoteable(client)

	listVoteables(client)
}