package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func Router(handler *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/api/audiofiles/requests/{id}", handler.GetByRequestID)
	r.Post("/api/audiofiles", handler.SaveFile)
	r.Get("/api/audio/{name}", handler.GetFile)
	return r
}
