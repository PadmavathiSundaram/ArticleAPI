package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/PadmavathiSundaram/ArticleAPI/pkg/client"
	model "github.com/PadmavathiSundaram/ArticleAPI/pkg/model"
	"github.com/PadmavathiSundaram/ArticleAPI/pkg/rest"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	defer close(quit)

	dbClient := newDBClient()
	server := newServer(dbClient)
	go gracefullShutdown(server, dbClient, quit, done)

	log.Infoln("Server is ready to handle requests at", *port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalln("Could not listen on", *port, err)
	}

	<-done
	log.Infoln("Server stopped")

}

// REF: https://marcofranssen.nl/go-webserver-with-graceful-shutdown/
func gracefullShutdown(server *http.Server, dbClient client.DBClient, quit <-chan os.Signal, done chan<- bool) {
	<-quit
	log.Infoln("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Fatalln(dbClient.DBDestroy())
	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalln("Could not gracefully shutdown the server: ", err)
	}
	close(done)
}
func newDBClient() client.DBClient {
	config, err := model.LoadConfig(*config)
	if err != nil {
		log.Fatalln("Could not Load Configurations.", err)
	}
	DBClient := client.NewDBClient()
	if err := DBClient.DBInit(config.DBProperties); err != nil {
		log.Fatalln("Could not Connect to the Database.", err)
	}
	return DBClient
}
func newServer(DBClient client.DBClient) *http.Server {
	router := chi.NewRouter()
	router.Use(mw.Logger)
	router.Use(middleware.Timeout(5 * time.Second))

	articlestore := client.NewArticleStore(DBClient)
	articleDelegate := rest.NewArticleDelegate(articlestore)
	rest.SetupRoutes(router, articleDelegate)
	server := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", *port),
	}
	return server
}
