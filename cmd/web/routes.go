package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.clientError(w, http.StatusNotFound)
	})

	fileServer := http.FileServer(http.Dir("./ui/static"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(addSecurityHeaders, app.sessionManager.LoadAndSave, noSurf)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.viewLandingPage))

	router.Handler(http.MethodGet, "/factos/view/:id", dynamic.ThenFunc(app.viewFactosById))
	router.Handler(http.MethodGet, "/factos/create", dynamic.ThenFunc(app.createFactos))
	router.Handler(http.MethodPost, "/factos/create", dynamic.ThenFunc(app.createFactosPost))

	router.Handler(http.MethodGet, "/results/:tournamentId", dynamic.ThenFunc(app.viewTournamentResults))
	router.Handler(http.MethodGet, "/upcoming/:tournamentId", dynamic.ThenFunc(app.viewFutureFixtures))
	router.Handler(http.MethodGet, "/predictions/:tournamentId", dynamic.ThenFunc(app.viewTournamentPredictions))

	router.Handler(http.MethodGet, "/authenticate", dynamic.ThenFunc(app.auth))
	router.Handler(http.MethodGet, "/auth/google/callback", dynamic.ThenFunc(app.authCallback))
	router.Handler(http.MethodGet, "/signup", dynamic.ThenFunc(app.viewSignUp))
	router.Handler(http.MethodPost, "/signup", dynamic.ThenFunc(app.postSignUp))

	return router
}
