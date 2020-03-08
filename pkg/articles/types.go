package articles

// Article to be persisted in db
type Article struct {
	ArticleID    string
	Title string
	Date  string
	Body  string
	Tags  []string
}
