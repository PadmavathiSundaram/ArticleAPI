package rest

import (
	"net/http"

	"github.com/PadmavathiSundaram/ArticleAPI/pkg/client"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// SetupRoutes sets up pet service routes for the given router
func SetupRoutes(r chi.Router, s Service) {
	r.Route("/api/articles", func(r chi.Router) {
		r.Post("/", s.PostArticle)
		r.Get("/{id}", s.GetArticle)
	})
}

// NewArticleService creates a new article service
func NewArticleService(dbClient client.DBClient) Service {
	return &service{dbClient: dbClient}
}

// Service defines a rest api for interaction
type Service interface {
	GetArticle(w http.ResponseWriter, r *http.Request)
	PostArticle(w http.ResponseWriter, r *http.Request)
}

type service struct {
	dbClient client.DBClient
}

// GetArticle handles a GET request to retrieve a Article
func (ps *service) GetArticle(w http.ResponseWriter, r *http.Request) {
	ps.dbClient.Select()
	render.Status(r, http.StatusOK)
}

// PostArticle handles a POST request to add a new Article
func (ps *service) PostArticle(w http.ResponseWriter, r *http.Request) {
	ps.dbClient.Insert()
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, nil)
}
