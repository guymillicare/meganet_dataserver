package routes

import (
	"sportsbook-backend/internal/controllers"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func SetupRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/datafeed", func(r chi.Router) {
		r.Get("/outcomes/{eventId}", controllers.GetOutcomesByEventId)
	})

	return router
}
