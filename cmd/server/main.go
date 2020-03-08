package main

import (
	"fmt"
	"net/http"

	"github.com/PadmavathiSundaram/ArticleAPI/pkg/articles"
	"github.com/PadmavathiSundaram/ArticleAPI/pkg/client"
	"github.com/PadmavathiSundaram/ArticleAPI/pkg/rest"
	"github.com/go-chi/chi"
	mw "github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	port   = kingpin.Flag("port", "port").Short('p').Default("4852").Int()
	config = kingpin.Flag("config", "config").Short('c').Default("cmd/server/config/config.standalone.json").File()
)

func main() {
	kingpin.Parse()

	router := chi.NewRouter()
	router.Use(mw.Logger)

	DBProperties, err := articles.LoadDBProperties(*config)
	if err != nil {
		log.Fatalf("Could not Load Configurations. %v", err)
	}
	articlestore, err := client.NewArticleStore(DBProperties)
	if err != nil {
		log.Fatalf("Could not connect data storage. %v", err)
	}
	articleService := rest.NewArticleService(articlestore)
	rest.SetupRoutes(router, articleService)
	server := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", *port),
	}
	log.Infoln("Server listening on port", *port)
	log.Fatal(server.ListenAndServe())
}
