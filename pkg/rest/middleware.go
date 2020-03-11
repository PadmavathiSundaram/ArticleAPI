package rest

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func apiLogger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Infoln("Audit log : Request", r.URL)
		start := time.Now()
		next.ServeHTTP(w, r)
		end := time.Now()
		log.Infoln("Audit Log : Request ", r.URL, " processed in ", end.Sub(start))
	}

	return http.HandlerFunc(fn)
}
