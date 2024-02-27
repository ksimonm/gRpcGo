package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	api "buffup/GolangTechTask/api"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	api.RegisterVotingServiceServer(s, &server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
}

func bufDialer(ctx context.Context, address string) (net.Conn, error) {
	return lis.Dial()
}

func TestMain(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	values := []string{"1", "2", "3", "4", "5", "6"}
	for _, v := range values {
		t.Run("CreateVoteable", func(t *testing.T) {
			client := api.NewVotingServiceClient(conn)

			question := "Q" + v

			resp, err := client.CreateVoteable(ctx, &api.CreateVoteableRequest{
				Question: question,
				Answers:  []string{question + "A1", question + "A2"},
			})

			fmt.Println("UUID:", resp.GetUuid())

			if err != nil {
				t.Fatal(err)
			}

			if resp.GetUuid() == "" {
				t.Fatal("uuid is nil")
			}
		})
	}

	page := []int32{1, 2}
	for _, v := range page {
		t.Run("ListVoteables", func(t *testing.T) {
			client := api.NewVotingServiceClient(conn)

			resp, err := client.ListVoteables(ctx, &api.ListVoteableRequest{Page: v, Size: 3})

			fmt.Println("Votables:", resp.GetVotables())

			if err != nil {
				t.Fatal(err)
			}

			if len(resp.GetVotables()) != 3 {
				t.Fatal("uuid is nil")
			}
		})
	}
}
