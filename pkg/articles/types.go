package articles

// Article to be persisted in db
type Article struct {
	ArticleID string   `json:"id,omitempty" bson:"ArticleID,omitempty"`
	Title     string   `json:"title,omitempty" bson:"Title,omitempty"`
	Date      string   `json:"date,omitempty" bson:"Date,omitempty"`
	Body      string   `json:"body,omitempty" bson:"Body,omitempty"`
	Tags      []string `json:"tags,omitempty" bson:"Tags,omitempty"`
}
