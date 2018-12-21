package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Routes() *chi.Mux {
	Route := chi.NewRouter()
	Route.Use(middleware.RequestID)
	Route.Use(middleware.Recoverer)

	Route.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hopes, Home"))
	})
	Route.Route("/pubsub", func(r chi.Router) {
		r.Mount("/", PubSubHandler())
	})

	return Route
}

func main() {
	router := Routes()

	walkFunc := func(method string, route string, handler http.Handler,
		middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route)
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error())
	}

	proj := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if proj == "" {
		fmt.Fprintf(os.Stderr, "GOOGLE_CLOUD_PROJECT environment variable must be set.\n")
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("No port found on Env, use Default %s", port)
	}
	log.Printf("Listening on port %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), router)

}
