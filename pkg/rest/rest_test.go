package rest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
)

type mockArticleService struct{}

func (s *mockArticleService) HealthCheck(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
}
func (s *mockArticleService) GetArticle(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
}
func (s *mockArticleService) SearchTags(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
}
func (s *mockArticleService) PostArticle(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, nil)
}

func TestSetupRoutes(t *testing.T) {
	router := chi.NewRouter()
	mockArticleService := &mockArticleService{}
	SetupRoutes(router, mockArticleService)
	server := httptest.NewServer(router)
	defer server.Close()

	var resp *http.Response
	var err error
	testScenarios := []struct {
		Description string
		Method      string
		URL         string
		StatusCode  int
	}{
		{"Route HealthCheck", "Get", "/api/healthcheck", 200},
		{"Route Get articles", "Get", "/api/articles/1", 200},
		{"Route Search tags", "Get", "/api/tags/tagName/2019-10-02", 200},
		{"Route Post articles", "Post", "/api/articles/", 201},
	}
	for _, td := range testScenarios {
		t.Run(fmt.Sprintf("%s - method %v  Url %v : %v",
			td.Description, td.Method, td.URL, td.StatusCode), func(t *testing.T) {
			switch td.Method {
			case "Get":
				resp, err = http.Get(server.URL + td.URL)
			case "Post":
				resp, err = http.Post(server.URL+td.URL, "application/json", nil)
			}
			assert.NoError(t, err)
			assert.Equal(t, td.StatusCode, resp.StatusCode, "Received unexpected response code: %d\n", resp.StatusCode)
		})
	}

}
