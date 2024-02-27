package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	api "buffup/GolangTechTask/api"

	"google.golang.org/grpc"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoUri        = "mongodb://localhost"
	mongoDb         = "voteDb"
	mongoCollection = "voteable"
)

type server struct {
	api.UnimplementedVotingServiceServer
}

type Voteable struct {
	Uuid     string   `bson:"uuid,omitempty"`
	Question string   `bson:"question,omitempty"`
	Answers  []string `bson:"answers,omitempty"`
}

func getMongoClient() (*mongo.Client, error) {
	var clientInstance *mongo.Client
	var clientInstanceError error
	var mongoOnce sync.Once

	mongoOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(mongoUri)

		client, err := mongo.Connect(context.TODO(), clientOptions)

		if err != nil {
			clientInstanceError = err
		}

		clientInstance = client
	})
	return clientInstance, clientInstanceError
}

func getCollection() (*mongo.Collection, *mongo.Client, error) {
	client, err := getMongoClient()
	if err != nil {
		return nil, nil, err
	}

	return client.Database(mongoDb).Collection(mongoCollection), client, nil
}

func (s *server) CreateVoteable(ctx context.Context, in *api.CreateVoteableRequest) (*api.CreateVoteableResponse, error) {
	log.Printf("Create Voteable: %v", in)

	coll, client, err := getCollection()
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	newUuid := uuid.New()

	newVoteable := Voteable{
		Uuid:     newUuid.String(),
		Question: in.GetQuestion(),
		Answers:  in.GetAnswers(),
	}

	result, err := coll.InsertOne(context.TODO(), newVoteable)

	if err != nil {
		panic(err)
	}

	log.Println("Result", result)

	return &api.CreateVoteableResponse{Uuid: newUuid.String()}, nil
}

func (s *server) ListVoteables(ctx context.Context, in *api.ListVoteableRequest) (*api.ListVoteableResponse, error) {
	log.Printf("List Voteables: %d %d", in.GetPage(), in.GetSize())

	coll, client, err := getCollection()
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	limit := int64(in.Size)
	skip := int64(in.Page)*limit - limit
	fOpt := options.FindOptions{Limit: &limit, Skip: &skip}

	cursor, err := coll.Find(context.TODO(), bson.D{{}}, &fOpt)

	if err != nil {
		panic(err)
	}

	log.Println("cursor:", cursor)

	var voteables []*api.Voteable
	if err = cursor.All(context.TODO(), &voteables); err != nil {
		panic(err)
	}

	log.Println("voteables:", voteables)

	return &api.ListVoteableResponse{Votables: voteables}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	api.RegisterVotingServiceServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
