package client

import (
	"context"
	"fmt"
	"time"

	"github.com/PadmavathiSundaram/ArticleAPI/pkg/articles"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBClient client to perform database operations
type DBClient interface {
	CloseConnection()
	Insert()
	Select()
}

// NewDBClient a client connected to DB to perform curd operations
func NewDBClient(URL string) DBClient {
	return &mongoClient{session: getConnection(URL)}
	//	return &mongoClient{session: nil}
}

type mongoClient struct {
	session *mongo.Client
}

func (mongo *mongoClient) Insert() {
	log.Infoln("Entered")
	Tags := []string{"health", "fit"}
	ash := articles.Article{ArticleID: "1", Title: "Ash", Date: "10-2-09", Body: "Pallet Town", Tags:Tags	}
	
	collection := mongo.session.Database("articlestore").Collection("articles")
	insertResult, err := collection.InsertOne(context.TODO(), ash)
	if err != nil {
		log.Fatal(err)
	}
	log.Infoln("Inserted a single document: ", insertResult.InsertedID)
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

func (mongo *mongoClient) Select() {
	log.Infoln("Entered select")
	filter := bson.M{"ArticleID": "1"}
	collection := mongo.session.Database("articlestore").Collection("articles")
	var result articles.Article

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Infoln(filter)
		log.Infoln(result)
		log.Fatal(err)
	}
	log.Infoln("Found a single document: %+v\n", result)
	fmt.Printf("Found a single document: %+v\n", result)
}

func (mongo *mongoClient) CloseConnection() {
	err := mongo.session.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func getConnection(URL string) *mongo.Client {
	// to do to move it to config
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(URL)
	// to do move to config
	clientOptions.SetMaxPoolSize(10)
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
