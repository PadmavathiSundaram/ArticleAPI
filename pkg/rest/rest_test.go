package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PadmavathiSundaram/ArticleAPI/pkg/model"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
)

type mockArticleDelegate struct{}

func (d *mockArticleDelegate) HealthCheck(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
}
func (d *mockArticleDelegate) GetArticle(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
}
func (d *mockArticleDelegate) SearchTags(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
}
func (d *mockArticleDelegate) PostArticle(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, nil)
}

func TestSetupRoutes(t *testing.T) {
	router := chi.NewRouter()
	mockArticleStore := &mockArticleStore{status: true}
	mockArticleDelegate := NewArticleDelegate(mockArticleStore)
	SetupRoutes(router, mockArticleDelegate)
	server := httptest.NewServer(router)
	defer server.Close()

	testScenarios := []struct {
		Description string
		URL         string
		StatusCode  int
	}{
		{"Route Get articles not found", "/api/articles/error", 404},
		{"Route mismatched Search tags", "/api/tags/oo/2019-10-02", 404},
		{"Route empty path params Search tags", "/api/tags/ /77", 404},
		{"Route Get articles", "/api/articles/1", 200},
		{"Route Search tags", "/api/tags/tagName/2019-10-02", 200},
		{"Route Post articles", "/api/healthcheck", 200},
	}
	for _, td := range testScenarios {
		t.Run(fmt.Sprintf("%s - method Get  Url %v : %v",
			td.Description, td.URL, td.StatusCode), func(t *testing.T) {
			resp, err := http.Get(server.URL + td.URL)
			assert.NoError(t, err)
			assert.Equal(t, td.StatusCode, resp.StatusCode, "Received unexpected response code: %d\n", resp.StatusCode)
		})
	}

}

func TestSetupRoutesFailedHealthCheck(t *testing.T) {

	router := chi.NewRouter()
	mockArticleStore := &mockArticleStore{status: false}
	mockArticleStore.HealthCheck()
	mockArticleDelegate := NewArticleDelegate(mockArticleStore)
	SetupRoutes(router, mockArticleDelegate)
	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/healthcheck")
	assert.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode, "Received unexpected response code: %d\n", resp.StatusCode)
}
func TestSetupRoutesPostNilArticle(t *testing.T) {

	router := chi.NewRouter()
	mockArticleStore := &mockArticleStore{status: true}
	mockArticleDelegate := NewArticleDelegate(mockArticleStore)
	SetupRoutes(router, mockArticleDelegate)
	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Post(server.URL+"/api/articles/", "application/json", bytes.NewBuffer(nil))
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode, "Received unexpected response code: %d\n", resp.StatusCode)
}

func TestSetupRoutesPostValidArticle(t *testing.T) {

	router := chi.NewRouter()
	mockArticleStore := &mockArticleStore{status: true}
	mockArticleDelegate := NewArticleDelegate(mockArticleStore)
	SetupRoutes(router, mockArticleDelegate)
	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Post(server.URL+"/api/articles/", "application/json", bytes.NewBuffer(nil))
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode, "Received unexpected response code: %d\n", resp.StatusCode)

}

func TestSetupRoutesPostNotAnArticle(t *testing.T) {

	router := chi.NewRouter()
	mockArticleStore := &mockArticleStore{status: true}
	mockArticleDelegate := NewArticleDelegate(mockArticleStore)
	SetupRoutes(router, mockArticleDelegate)
	server := httptest.NewServer(router)
	defer server.Close()

	article := &model.Article{}
	article.ArticleID = "11"
	s := "success"
	article.Tags = []*string{&s}
	article.Date = "2019-02-01"

	requestBody, parserErr := json.Marshal(article)
	if parserErr != nil {
		panic(fmt.Errorf("Error in test code, could not marshal testPet to json. %v", parserErr))
	}

	resp, err := http.Post(server.URL+"/api/articles/", "application/json", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode, "Received unexpected response code: %d\n", resp.StatusCode)

}

func TestSetupRoutesPostInvalidArticle(t *testing.T) {

	router := chi.NewRouter()
	mockArticleStore := &mockArticleStore{status: false}
	mockArticleDelegate := NewArticleDelegate(mockArticleStore)
	SetupRoutes(router, mockArticleDelegate)
	server := httptest.NewServer(router)
	defer server.Close()

	article := &model.Article{}
	article.ArticleID = "1"
	s := "success"
	article.Tags = []*string{&s}
	article.Date = "2019-02-01"

	duplicateArticle := &model.Article{}
	duplicateArticle.ArticleID = "1"
	duplicateArticle.Tags = []*string{&s}
	duplicateArticle.Date = "2019-02-01"

	invalidArticle := &model.Article{}
	invalidArticle.ArticleID = ""
	invalidArticle.Tags = []*string{}
	invalidArticle.Date = ""

	invalidIDArticle := &model.Article{}
	invalidIDArticle.ArticleID = ""
	invalidIDArticle.Tags = []*string{&s}
	invalidIDArticle.Date = "2019-02-02"

	invalidDateArticle := &model.Article{}
	invalidDateArticle.ArticleID = "2222"
	invalidDateArticle.Tags = []*string{&s}
	invalidDateArticle.Date = "21313223"

	testScenarios := []struct {
		Description string
		Body        *model.Article
		StatusCode  int
	}{
		{"Duplicate record", duplicateArticle, 409},
		{"Missing Mandatory fiels", invalidArticle, 400},
		{"Invalid Date", invalidDateArticle, 400},
		{"nil Request body", nil, 400},
		{"Invalid Date", invalidIDArticle, 400},
	}
	for _, td := range testScenarios {
		t.Run(fmt.Sprintf("%s - method post : %v",
			td.Description, td.StatusCode), func(t *testing.T) {
			requestBody, parserErr := json.Marshal(td.Body)
			if parserErr != nil {
				panic(fmt.Errorf("Error in test code, could not marshal testPet to json. %v", parserErr))
			}
			resp, _ := http.Post(server.URL+"/api/articles/", "application/json", bytes.NewBuffer(requestBody))
			assert.Equal(t, td.StatusCode, resp.StatusCode, "Received unexpected response code: %d\n", resp.StatusCode)
		})
	}
}

type mockArticleStore struct{ status bool }

func (store *mockArticleStore) HealthCheck() bool {
	return store.status
}

func (store *mockArticleStore) CreatArticle(article *model.Article) error {
	if article.ArticleID == "11" {
		return nil
	}

	return model.Errorf(model.ErrDuplicate, "Duplicate key error")
}

func (store *mockArticleStore) ReadArticleByID(articleID string) (*model.Article, error) {
	if articleID == "1" {
		return &model.Article{}, nil
	}
	return nil, model.Errorf(model.ErrNotFound, "Not found error")
}

func (store *mockArticleStore) SearchTagsByDate(date string, tag string) ([]*model.Tagsview, error) {
	var tags []*model.Tagsview
	if tag == "tagName" {
		tags = append(tags, &model.Tagsview{})
		return tags, nil
	}
	return nil, model.Errorf(model.ErrNotFound, "Not found error")
}
