package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PadmavathiSundaram/ArticleAPI/pkg/articles"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBClient client to perform database operations
type DBClient interface {
	CloseConnectionPool()
	OpenConnectionPool(DBProperties *articles.DBProperties) (*mongo.Client, error)
}

// NewDBClient a client connected to DB to perform curd operations
func NewDBClient() DBClient {
	return &mongoClient{}
}

type mongoClient struct {
	session *mongo.Client
}

func (mc *mongoClient) CloseConnectionPool() {
	err := mc.session.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

// OpenConnectionPool for mongo db
func (mc *mongoClient) OpenConnectionPool(DBProperties *articles.DBProperties) (*mongo.Client, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(DBProperties.MaxTimeOut)*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(DBProperties.URL)

	clientOptions.SetMaxPoolSize(uint64(DBProperties.MaxThreadPoolSize))
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")
	mc.session = client
	return mc.session, nil
}
