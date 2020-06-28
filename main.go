package main

import (
	"log"
	"net/http"
	"os"

	expogluster "github.com/aminueza/docker-gluester-exporter/expogluster"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	promExp := expogluster.NewPromExporter()
	router := mux.NewRouter()
	router.Use(cacheMiddleware)

	server := &expogluster.Exporter{
		Router:   router,
		Hostname: promExp.Hostname,
		Volumes:  promExp.Volumes,
		Profile:  promExp.Profile,
		Quota:    promExp.Quota,
	}

	prometheus.MustRegister(server)

	log.Println("Server is listening: http://" + server.Hostname)
	err := http.ListenAndServe(server.Hostname, handlers.LoggingHandler(os.Stdout, router))
	if err != nil {
		log.Fatal(err)
	}

}

func cacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")

		next.ServeHTTP(w, r)
	})
}
