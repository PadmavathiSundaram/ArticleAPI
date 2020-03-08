package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/PadmavathiSundaram/ArticleAPI/pkg/articles"
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

// maps from internal errors to response status codes
// renderErrorResponse defaults to internal server error
// if a specific error code is not defined.
var errStatusMap = map[int]int{
	ErrInvalidInput: http.StatusBadRequest,
	ErrNotFound:     http.StatusNotFound,
}

// renderErrorResponse handles http responses in the case of an error
func renderErrorResponse(w http.ResponseWriter, err error) {
	message := err.Error()
	responseStatus := http.StatusInternalServerError
	// artricle service Errors store more specific response information
	if specificError, ok := err.(*Error); ok {
		message = specificError.Message
		// Attempt to get a more specific status code
		if status, ok := errStatusMap[specificError.Code]; ok {
			responseStatus = status
		}
	}
	http.Error(w, message, responseStatus)
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
	articleID, err := readArticleID(r)
	if err != nil {
		renderErrorResponse(w, err)
		return
	}
	article, err := ps.dbClient.Select(articleID)
	if err != nil {
		if "mongo: no documents in result" == err.Error() {
			renderErrorResponse(w, ErrorEf(ErrNotFound, err, "Article Not Found"))
			return
		}
		renderErrorResponse(w, err)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, article)
}

// PostArticle handles a POST request to add a new Article
func (ps *service) PostArticle(w http.ResponseWriter, r *http.Request) {
	article, err := readArticleBody(r)
	if err != nil {
		renderErrorResponse(w, err)
		return
	}

	err = ps.dbClient.Insert(article)
	if err != nil {
		renderErrorResponse(w, err)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, nil)
}

func readArticleID(r *http.Request) (string, error) {
	articleID := chi.URLParam(r, "id")
	if articleID == "" {
		// Reaching this indicates a bug. At this point, request context should contain an id
		return "", Errorf(ErrUnknown, "article ID was lost somewhere")
	}
	_, err := strconv.ParseInt(articleID, 10, 32)
	if err != nil {
		return "", Errorf(ErrInvalidInput, "Invalid article ID %v. ID should be a number", articleID)
	}
	return articleID, nil
}

func readArticleBody(r *http.Request) (*articles.Article, error) {
	if r.Body == nil {
		return nil, Errorf(ErrInvalidInput, "No request body")
	}
	articleData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, ErrorEf(ErrInvalidInput, err, "Bad request body")
	}
	var article articles.Article
	if err = json.Unmarshal(articleData, &article); err != nil {
		return nil, ErrorEf(ErrInvalidInput, err, "Invalid Article data")
	}
	return &article, nil
}
