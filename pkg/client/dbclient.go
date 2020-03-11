package client

import (
	"context"
	"time"

	model "github.com/PadmavathiSundaram/ArticleAPI/pkg/model"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// NewDBClient a client connected to DB to perform curd operations
func NewDBClient() DBClient {
	return &mongoClient{}
}

// ConvertMapToBsonD helps generate bson structed query compartible with mongo
func ConvertMapToBsonD(m map[string]string) bson.D {
	d := bson.D{}
	for k, v := range m {
		d = append(d, bson.E{Key: k, Value: v})
	}
	return d
}

// ConvertMapToBsonM helps generate bson structed query compartible with mongo
func ConvertMapToBsonM(fields map[string]string) bson.M {
	m := bson.M{}
	for k, v := range fields {
		m[k] = v
	}
	return m
}

// DBClient client to perform database operations
type DBClient interface {
	DBDestroy() error
	HealthCheck() bool
	DBInit(DBProperties model.DBProperties) (err error)
	Read(key string, value string) interface{}
	Write(document interface{}) (interface{}, error)
	SimpleQuery(filterFields map[string]string, sortFields map[string]string) (interface{}, error)
	AdvancedQuery(pipeLine interface{}) (interface{}, error)
	Delete()
}

type mongoClient struct {
	session    *mongo.Client
	collection *mongo.Collection
}

func (mc *mongoClient) DBInit(DBProperties model.DBProperties) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(DBProperties.MaxTimeOut)*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(DBProperties.URL)
	clientOptions.SetMaxPoolSize(uint64(DBProperties.MaxThreadPoolSize))

	mc.session, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	if err = mc.session.Ping(ctx, nil); err != nil {
		return err
	}

	log.Infoln("Connected to MongoDB!")
	mc.collection = mc.session.Database(DBProperties.DatabaseName).Collection(DBProperties.CollectionName)

	if len(DBProperties.Indexes) > 0 {
		if err = mc.setUpIndexes(DBProperties.Indexes); err != nil {
			return err
		}
	}
	log.Infoln("Loaded Collection ", DBProperties.CollectionName)
	return nil
}

func (mc *mongoClient) HealthCheck() bool {
	if err := mc.session.Ping(context.Background(), nil); err != nil {
		return false
	}
	return true
}

func (mc *mongoClient) setUpIndexes(indexes map[string]bool) error {
	// Declare an array of bsonx models for the indexes
	models := []mongo.IndexModel{}
	collation := options.Collation{Strength: 2, Locale: "en"}
	for k, v := range indexes {
		index := mongo.IndexModel{
			Keys:    bsonx.Doc{{Key: k, Value: bsonx.Int32(-1)}},
			Options: options.Index().SetCollation(&collation).SetUnique(v),
		}
		models = append(models, index)
	}

	if _, err := mc.collection.Indexes().CreateMany(context.Background(), models); err != nil {
		log.Errorln("Indexes().CreateMany() ERROR:", err)
		return err
	}

	log.Infoln("Created Indexes:", models)
	return nil
}

func (mc *mongoClient) Read(key string, value string) interface{} {
	filter := bson.D{{Key: key, Value: value}}
	result := mc.collection.FindOne(context.Background(), filter)

	log.Infoln("Found a single document: \n", result)
	return result
}
func (mc *mongoClient) Write(document interface{}) (interface{}, error) {
	insertResult, err := mc.collection.InsertOne(context.Background(), document)
	if err != nil {
		return nil, err
	}
	log.Infoln("Inserted a single document: ", insertResult.InsertedID)
	return insertResult.InsertedID, nil
}
func (mc *mongoClient) SimpleQuery(filterFields map[string]string, sortFields map[string]string) (interface{}, error) {
	findOptions := options.Find()
	if len(sortFields) > 0 {
		sort := ConvertMapToBsonD(sortFields)
		findOptions.SetSort(sort)
	}
	filter := ConvertMapToBsonD(filterFields)

	return mc.collection.Find(context.Background(), filter, findOptions)
}

// pipeLine example: []bson.M anything that the mongo db aggregate can process
func (mc *mongoClient) AdvancedQuery(pipeLine interface{}) (interface{}, error) {
	return mc.collection.Aggregate(context.Background(), pipeLine)
}

// To Do implement for future use case
func (mc *mongoClient) Delete() {
	log.Infoln("Delete yet to be implemented")
}
func (mc *mongoClient) DBDestroy() error {
	err := mc.session.Disconnect(context.Background())
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	log.Infoln("Connection to MongoDB closed.")
	return nil
}
