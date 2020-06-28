package expogluster

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//API routes to commands
func API(server *Exporter) {

	apiRouter := server.Router.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(corsMiddleware)

	// routes with no auth (need to be listed in checkIfPathHasNoAuth method)
	apiRouter.Handle("/metrics", promhttp.Handler())

	// catch all - not found
	apiRouter.PathPrefix("/").HandlerFunc(routeNotFound)

}

func routeNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	httpError(errors.New("Route not found"), http.StatusNotFound, w)
}

func httpError(err error, status int, w http.ResponseWriter) {
	w.WriteHeader(status)

	payload := map[string]string{
		"error": err.Error(),
	}

	js, _ := json.MarshalIndent(payload, "", " ")
	w.Write(js)
}

func setupCorsResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, CONNECT, HEAD, PATCH, TRACE")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setupCorsResponse(w, r)

		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}
		next.ServeHTTP(w, r)
	})
}
