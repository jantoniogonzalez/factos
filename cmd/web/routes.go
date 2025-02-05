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
	authenticationRequired := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.viewLandingPage))

	router.Handler(http.MethodGet, "/factos/view/:id", dynamic.ThenFunc(app.viewFactosById))
	// I think we might not really need this one as it is loaded from the same page as the Tournament stuff
	router.Handler(http.MethodGet, "/factos/create/:matchId", dynamic.ThenFunc(app.createFactos))
	router.Handler(http.MethodPost, "/factos/create", authenticationRequired.ThenFunc(app.createFactosPost))

	router.Handler(http.MethodGet, "/results/:tournamentId/:season", dynamic.ThenFunc(app.viewTournamentResults))
	router.Handler(http.MethodGet, "/upcoming/:tournamentId/:season", dynamic.ThenFunc(app.viewTournamentFutureFixtures))
	router.Handler(http.MethodGet, "/predictions/:tournamentId/:season", dynamic.ThenFunc(app.viewTournamentPredictions))

	router.Handler(http.MethodGet, "/authenticate", dynamic.ThenFunc(app.auth))
	router.Handler(http.MethodGet, "/auth/google/callback", dynamic.ThenFunc(app.authCallback))
	router.Handler(http.MethodGet, "/signup", dynamic.ThenFunc(app.viewSignUp))
	router.Handler(http.MethodPost, "/signup", dynamic.ThenFunc(app.postSignUp))
	router.Handler(http.MethodPost, "/logout", dynamic.ThenFunc(app.logout))

	return router
}
