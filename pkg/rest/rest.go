package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/PadmavathiSundaram/ArticleAPI/pkg/client"
	"github.com/PadmavathiSundaram/ArticleAPI/pkg/model"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// SetupRoutes sets up Article service routes for the given router
func SetupRoutes(r chi.Router, d Delegate) {
	r.Use(recoverHandler, apiLogger)
	r.Route("/api", func(r chi.Router) {
		r.Get("/healthcheck", d.HealthCheck)
		r.Get("/tags/{tagName}/{date}", d.SearchTags)
		r.Route("/articles", func(r chi.Router) {
			r.Post("/", d.PostArticle)
			r.Get("/{id}", d.GetArticle)
		})
	})

}

// maps from internal errors to response status codes
// renderErrorResponse defaults to internal server error
// if a specific error code is not defined.
var errStatusMap = map[int]int{
	model.ErrDuplicate:    http.StatusConflict,
	model.ErrInvalidInput: http.StatusBadRequest,
	model.ErrNotFound:     http.StatusNotFound,
}

// renderErrorResponse handles http responses in the case of an error
func renderErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	responseStatus := http.StatusInternalServerError
	// artricle service model.Errors store more specific response information
	if specificError, ok := err.(*model.Error); ok {
		if status, ok := errStatusMap[specificError.Code]; ok {
			responseStatus = status
		}
	}

	render.Status(r, responseStatus)
	render.JSON(w, r, err)
}

func validateRequestArticle(article model.Article) error {
	if article.ArticleID == "" || article.Date == "" || len(article.Tags) <= 0 {
		return model.Errorf(model.ErrInvalidInput, "Invalid Article data - id, date and atleast 1 tag is mandatory")

	}
	layout := "2006-01-02"
	if _, err := time.Parse(layout, article.Date); err != nil {
		return model.ErrorEf(model.ErrInvalidInput, err, "Invalid date Format expected format YYYY-MM-DD 2016-09-22")
	}
	return nil

}

func readArticleBody(r *http.Request) (*model.Article, error) {
	if r.Body == nil {
		return nil, model.Errorf(model.ErrInvalidInput, "No request body")
	}
	articleData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, model.ErrorEf(model.ErrInvalidInput, err, "Bad request body")
	}
	var article model.Article
	if err = json.Unmarshal(articleData, &article); err != nil {
		return nil, model.ErrorEf(model.ErrInvalidInput, err, "Invalid Article data")
	}

	if err = validateRequestArticle(article); err != nil {
		return nil, err
	}
	return &article, nil
}

// NewArticleDelegate creates a new article service
func NewArticleDelegate(articleStore client.ArticleStore) Delegate {
	return &delegate{articleStore: articleStore}
}

// Delegate defines a rest api for interaction
type Delegate interface {
	GetArticle(w http.ResponseWriter, r *http.Request)
	SearchTags(w http.ResponseWriter, r *http.Request)
	PostArticle(w http.ResponseWriter, r *http.Request)
	HealthCheck(w http.ResponseWriter, r *http.Request)
}

type delegate struct {
	articleStore client.ArticleStore
}

func (d *delegate) HealthCheck(w http.ResponseWriter, r *http.Request) {
	health := model.Health{Status: "success"}
	if !d.articleStore.HealthCheck() {
		render.Status(r, http.StatusInternalServerError)
		health.Status = "failure"
		render.JSON(w, r, health)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, health)
}

// GetArticle handles a GET request to retrieve a Article
func (d *delegate) GetArticle(w http.ResponseWriter, r *http.Request) {
	articleID, err := readArticleID(r)
	if err != nil {
		renderErrorResponse(w, r, err)
		return
	}
	article, err := d.articleStore.ReadArticleByID(articleID)
	if err != nil {
		if "mongo: no documents in result" == err.Error() {
			renderErrorResponse(w, r, model.ErrorEf(model.ErrNotFound, err, "Article Not Found"))
			return
		}
		renderErrorResponse(w, r, err)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, article)
}

func (d *delegate) SearchTags(w http.ResponseWriter, r *http.Request) {
	date, tagName := readArticleParams(r)
	if date == "" || tagName == "" {
		renderErrorResponse(w, r, model.Errorf(model.ErrInvalidInput, "date and tagName are mandatory"))
		return
	}
	tags, err := d.articleStore.SearchTagsByDate(date, tagName)
	if err != nil {
		if "mongo: no documents in result" == err.Error() {
			renderErrorResponse(w, r, model.ErrorEf(model.ErrNotFound, err, "Article Not Found"))
			return
		}
		renderErrorResponse(w, r, err)
		return
	}
	if len(tags) == 0 {
		renderErrorResponse(w, r, model.ErrorEf(model.ErrNotFound, err, "No Matching Articles Found"))
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, tags[0])
}

// PostArticle handles a POST request to add a new Article
func (d *delegate) PostArticle(w http.ResponseWriter, r *http.Request) {
	article, err := readArticleBody(r)
	if err != nil {
		renderErrorResponse(w, r, err)
		return
	}

	err = d.articleStore.CreatArticle(article)
	if err != nil {
		renderErrorResponse(w, r, err)
		return
	}
	render.Status(r, http.StatusCreated)
}

func readArticleParams(r *http.Request) (string, string) {
	date := chi.URLParam(r, "date")
	tagName := chi.URLParam(r, "tagName")
	return date, tagName
}
func readArticleID(r *http.Request) (string, error) {
	articleID := chi.URLParam(r, "id")
	if articleID == "" {
		// Reaching this indicates a bug. At this point, request context should contain an id
		return "", model.Errorf(model.ErrUnknown, "article ID was lost somewhere")
	}
	return articleID, nil
}
