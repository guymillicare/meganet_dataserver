package routes

import (
	"sportsbook-backend/internal/controllers"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func SetupRouter(router *chi.Mux) {
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/datafeed", func(r chi.Router) {
		r.Get("/outcomes/{eventRefId}", controllers.GetOutcomesByEventId)
		r.Get("/events/{betType}/{sportId}", controllers.GetSportEventsWithOdds)
		r.Get("/events/{betType}/{sportId}/{countryId}", controllers.GetSportEventsWithOdds)
		r.Get("/events/{betType}/{sportId}/{countryId}/{leagueId}", controllers.GetSportEventsWithOdds)
		r.Post("/events", controllers.GetSportEventsWithLiveOdds)
	})
}
