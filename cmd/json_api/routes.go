package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	dynamic := alice.New(addSecurityHeaders, app.sessionManager.LoadAndSave, noSurf)

	// Routes
	// *Authentication
	router.Handler(http.MethodGet, "/authenticate", dynamic.ThenFunc(app.auth))
	router.Handler(http.MethodGet, "/auth/google/callback", dynamic.ThenFunc(app.authCallback))
	// TODO: Check if this is really necessary as this is really more like rendered in the front end
	router.Handler(http.MethodGet, "/signup", dynamic.ThenFunc(app.viewSignUp))
	router.Handler(http.MethodPost, "/signup", dynamic.ThenFunc(app.postSignUp))
	router.Handler(http.MethodPost, "/logout", dynamic.ThenFunc(app.logout))
	// *Factos
	// @returns: []FactosResponse
	router.Handler(http.MethodGet, "/factos/latest/:userId", dynamic.ThenFunc(app.viewLatestFactosByUserId))
	router.Handler(http.MethodGet, "/factos/:userId", dynamic.ThenFunc(app.viewAllFactosByUserId))
	router.Handler(http.MethodGet, "/factos/:leagueId/:userId", dynamic.ThenFunc(app.viewAllLeagueFactosByUserId))
	router.Handler(http.MethodGet, "/factos/latest/:leagueId/:userId", dynamic.ThenFunc(app.viewLatestLeagueFactosByUserId))
	router.Handler(http.MethodPost, "/new-facto/:matchId", dynamic.ThenFunc(app.createFacto))
	router.Handler(http.MethodPost, "/edit-facto/:factoId", dynamic.ThenFunc(app.editFacto))
	// *Fixtures
	// @returns: []FactosResponse
	router.Handler(http.MethodGet, "/fixtures/latest/:leagueId", dynamic.ThenFunc(app.viewLatestFixturesbyLeagueId))
	router.Handler(http.MethodGet, "/fixtures/upcoming/:leagueId", dynamic.ThenFunc(app.viewUpcomingFixturesbyLeagueId))

	return router
}
