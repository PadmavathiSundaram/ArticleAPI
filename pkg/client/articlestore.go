package client

import (
	"context"
	"strings"

	model "github.com/PadmavathiSundaram/ArticleAPI/pkg/model"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ArticleStore is the interface that performs DB operations
type ArticleStore interface {
	HealthCheck() bool
	CreatArticle(article *model.Article) error
	ReadArticleByID(articleID string) (*model.Article, error)
	SearchTagsByDate(date string, tag string) ([]*model.Tagsview, error)
}

// NewArticleStore a service to perform Article curd operations
func NewArticleStore(dbClient DBClient) ArticleStore {
	return &mongoArticleStore{dbClient: dbClient}
}

type mongoArticleStore struct {
	dbClient DBClient
}

func (store *mongoArticleStore) HealthCheck() bool {
	return store.dbClient.HealthCheck()
}
func (store *mongoArticleStore) CreatArticle(article *model.Article) error {
	if _, err := store.dbClient.Write(article); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return model.Errorf(model.ErrDuplicate, err.Error())
		}
		return model.Errorf(model.ErrUnknown, err.Error())
	}
	return nil
}

func (store *mongoArticleStore) ReadArticleByID(articleID string) (result *model.Article, e error) {

	res := store.dbClient.Read("ArticleID", articleID)
	singleResult, ok := res.(*mongo.SingleResult)
	if !ok {
		return nil, mongo.CommandError{Message: "Unable to parse Read Result"}
	}
	if err := singleResult.Decode(&result); err != nil {
		return nil, err
	}
	log.Infoln("Found a single article: \n", result)
	return result, nil
}

func (store *mongoArticleStore) SearchTagsByDate(date string, tag string) (results []*model.Tagsview, e error) {
	pipeLine := []bson.M{
		// filtes records based on date and tags
		{"$match": bson.M{"Date": date, "Tags": tag}},
		// sorts them based of creation order - descending
		{"$sort": bson.M{"_id": -1}},
		// groups the result set to aggregate
		// the articleIds related tags and count
		{"$group": bson.M{
			"_id": nil,
			// pushes all items into a list
			"articlesList": bson.M{"$push": "$ArticleID"},
			"relatedTags":  bson.M{"$push": "$Tags"},
			// identifys the total records impacted by the query
			"count": bson.M{"$sum": 1},
		}},
		// generates the projections/views
		{"$project": bson.M{
			// truncates the articles list to 10 entries
			"Articles": bson.M{"$slice": bson.A{
				"$articlesList",
				10},
			},
			// remopves duplicate entries and the tag in match from realted tags
			"RelatedTags": bson.M{"$reduce": bson.M{
				"input":        "$relatedTags",
				"initialValue": bson.A{},
				"in": bson.M{"$setDifference": bson.A{
					bson.M{"$setUnion": bson.A{"$$value", "$$this"}},
					bson.A{tag},
				}},
			}},
			// count of records on that day with the tag
			"Count": "$count",
			// the tag queried for
			"Tag": tag,
		}},
	}

	cur, err := store.dbClient.AdvancedQuery(pipeLine)
	if err != nil {
		return nil, err
	}

	cursor, ok := cur.(*mongo.Cursor)
	defer cursor.Close(context.Background())

	if !ok {
		return nil, mongo.CommandError{Message: "Invalid search operation"}
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cursor.Next(context.Background()) {
		// create a value into which the single document can be decoded
		var elem model.Tagsview
		if err := cursor.Decode(&elem); err != nil {
			return nil, err
		}
		results = append(results, &elem)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	log.Infoln("Search query Completed")
	return results, nil
}
