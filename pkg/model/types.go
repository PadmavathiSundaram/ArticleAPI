package model

import (
	"encoding/json"
	"fmt"
	"os"
)

// Health - API health
type Health struct {
	Status string `json:"status,omitempty"`
}

// Tagsview view of Tags stats
type Tagsview struct {
	Tag         string    `json:"tag,omitempty" bson:"Tag,omitempty"`
	Articles    []*string `json:"articles" bson:"Articles"`
	RelatedTags []*string `json:"related_tags" bson:"RelatedTags"`
	Count       int       `json:"count,omitempty" bson:"Count,omitempty"`
}

// Article to be persisted in db
type Article struct {
	ArticleID string    `json:"id,omitempty" bson:"ArticleID,omitempty"`
	Title     string    `json:"title,omitempty" bson:"Title,omitempty"`
	Date      string    `json:"date,omitempty" bson:"Date,omitempty"`
	Body      string    `json:"body,omitempty" bson:"Body,omitempty"`
	Tags      []*string `json:"tags,omitempty" bson:"Tags,omitempty"`
}

// Config data of the application
type Config struct {
	DBProperties DBProperties
}

// DBProperties settings
type DBProperties struct {
	URL               string
	DatabaseName      string
	CollectionName    string
	MaxThreadPoolSize int
	MaxTimeOut        int
	Indexes           map[string]bool // {indexfieldName : isUnique}
}

// LoadConfig from path specified
func LoadConfig(file *os.File) (*Config, error) {
	decoder := json.NewDecoder(file)
	Config := Config{}
	err := decoder.Decode(&Config)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &Config, err
}
