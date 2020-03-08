package articles

import (
	"encoding/json"
	"fmt"
	"os"
)

// Article to be persisted in db
type Article struct {
	ArticleID string   `json:"id,omitempty" bson:"ArticleID,omitempty"`
	Title     string   `json:"title,omitempty" bson:"Title,omitempty"`
	Date      string   `json:"date,omitempty" bson:"Date,omitempty"`
	Body      string   `json:"body,omitempty" bson:"Body,omitempty"`
	Tags      []string `json:"tags,omitempty" bson:"Tags,omitempty"`
}

// DBProperties settings
type DBProperties struct {
	URL               string
	DatabaseName      string
	CollectionName    string
	MaxThreadPoolSize int
	MaxTimeOut        int
}

// LoadDBProperties from path specified
func LoadDBProperties(file *os.File) (*DBProperties, error) {
	decoder := json.NewDecoder(file)
	DBProperties := DBProperties{}
	err := decoder.Decode(&DBProperties)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &DBProperties, err
}
