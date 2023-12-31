package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	//mux := c
	mux := chi.NewRouter()

	// specify who can connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: true,
		MaxAge: 300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		msg := jsonResponse {
			Message: "Wilkommen",
		}

		result, _ := json.MarshalIndent(msg, "", "\t")
		
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	})

	mux.Post("/", app.Broker)

	mux.Post("/handle", app.HandleSubmission)

	return mux
}