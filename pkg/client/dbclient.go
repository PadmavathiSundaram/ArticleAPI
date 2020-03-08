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
	Insert(article *articles.Article) error
	Select(articleID string) (*articles.Article, error)
	Search(date string, tag string) ([]*articles.Article, error)
}

// NewDBClient a client connected to DB to perform curd operations
func NewDBClient(URL string) DBClient {
	return &mongoClient{session: getConnection(URL)}
	//	return &mongoClient{session: nil}
}

type mongoClient struct {
	session *mongo.Client
}

func (mongo *mongoClient) Insert(article *articles.Article) error {
	log.Infoln("Entered Insert")

	collection := mongo.session.Database("articlestore").Collection("articles")
	insertResult, err := collection.InsertOne(context.TODO(), article)

	if err != nil {
		// log.Fatal(err)
		return err
	}

	log.Infoln("Inserted a single document: ", insertResult.InsertedID)
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return nil
}

func (mongo *mongoClient) Select(articleID string) (result *articles.Article, e error) {
	log.Infoln("Entered select")
	filter := bson.D{{"ArticleID", articleID}}

	collection := mongo.session.Database("articlestore").Collection("articles")

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		// log.Fatal(err)
		return nil, err
	}
	mongo.Search("10-2-09", "health")
	log.Infoln("Found a single document: %+v\n", result)
	fmt.Printf("Found a single document: %+v\n", result)
	return result, nil
}

func (mongo *mongoClient) Search(date string, tag string) (results []*articles.Article, e error) {
	log.Infoln("Entered search")

	findOptions := options.Find()
	filter := bson.D{{"Date", date}, {"Tags", tag}}

	collection := mongo.session.Database("articlestore").Collection("articles")

	cur, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		// log.Fatal(err)
		return nil, err
	}
	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem articles.Article
		err := cur.Decode(&elem)
		if err != nil {
			// log.Fatal(err)
			return nil, err
		}
		fmt.Printf("Found multiple documents (array of pointers): %+v\n", &elem)

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		// log.Fatal(err)
		return nil, err
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)

	return results, nil
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
