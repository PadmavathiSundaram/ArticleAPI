package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBClient client to perform database operations
type DBClient interface {
}

// NewDBClient a client connected to DB to perform curd operations
func NewDBClient(URL string) DBClient {
	return &mongoClient{session: getConnection(URL)}
}

type mongoClient struct {
	session *mongo.Client
}

func getConnection(URL string) *mongo.Client {
	// to do to move it to config
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(URL)
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}
