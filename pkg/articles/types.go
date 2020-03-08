package articles

// Article to be persisted in db
type Article struct {
	ID    string
	Title string
	Date  string
	Body  string
	Tags  []string
}
