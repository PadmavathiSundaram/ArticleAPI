package client

import (
	"context"
	"fmt"

	"github.com/PadmavathiSundaram/ArticleAPI/pkg/articles"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Store is the interface that performs DB operations
type Store interface {
	Insert(article *articles.Article) error
	Select(articleID string) (*articles.Article, error)
	Search(date string, tag string) ([]*articles.Article, error)
}

// NewArticleStore a client connected to DB to perform curd operations
func NewArticleStore(DBProperties articles.DBProperties) (Store, error) {
	mongo := NewDBClient()
	client, err := mongo.OpenConnectionPool(DBProperties)
	if err != nil {
		return nil, err
	}
	collection := client.Database(DBProperties.DatabaseName).Collection(DBProperties.CollectionName)
	return &articleStore{collection: collection}, nil
}

type articleStore struct {
	collection *mongo.Collection
}

func (store *articleStore) Insert(article *articles.Article) error {
	log.Infoln("Entered Insert")

	insertResult, err := store.collection.InsertOne(context.TODO(), article)

	if err != nil {

		return err
	}

	log.Infoln("Inserted a single document: ", insertResult.InsertedID)
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return nil
}

func (store *articleStore) Select(articleID string) (result *articles.Article, e error) {
	log.Infoln("Entered select")
	filter := bson.D{{"ArticleID", articleID}}

	err := store.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {

		return nil, err
	}
	log.Infoln("Found a single document: %+v\n", result)
	fmt.Printf("Found a single document: %+v\n", result)
	return result, nil
}

func (store *articleStore) Search(date string, tag string) (results []*articles.Article, e error) {
	log.Infoln("Entered search")

	findOptions := options.Find()
	filter := bson.D{{"Date", date}, {"Tags", tag}}

	cur, err := store.collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem articles.Article
		err := cur.Decode(&elem)
		if err != nil {

			return nil, err
		}
		fmt.Printf("Found multiple documents (array of pointers): %+v\n", &elem)

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {

		return nil, err
	}

	// Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)

	return results, nil
}
