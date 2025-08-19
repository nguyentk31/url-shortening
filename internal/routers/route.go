package routers

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nguyentk31/url-shortening/internal/database"
	"github.com/nguyentk31/url-shortening/internal/handlers"
)

func NewRouter(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	service := handlers.NewService(database.New(db))
	v1Router := chi.NewRouter()

	// Register shorten routes
	v1Router.Post("/shorten", service.CreateUrl)
	v1Router.Get("/shorten/{shortenID}", service.RetrieveUrl)
	v1Router.Put("/shorten/{shortenID}", service.UpdateUrl)
	v1Router.Delete("/shorten/{shortenID}", service.DeleteUrl)
	v1Router.Get("/shorten/{shortenID}/stats", service.StatsUrls)
	v1Router.Post("/shorten/{shortenID}/increment", service.IncrementAccessCount)

	// Register utility routes
	v1Router.Get("/health", handlers.HandlerReady)
	v1Router.Get("/err", handlers.HandleErr)

	r.Mount("/v1", v1Router)

	return r
}
