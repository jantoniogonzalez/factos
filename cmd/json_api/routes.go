package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.errorPage)

	// dynamic := alice.New(addSecurityHeaders, app.sessionManager.LoadAndSave, noSurf)
	dynamic := alice.New(addSecurityHeaders, app.sessionManager.LoadAndSave)

	// Routes
	// *Authentication
	router.Handler(http.MethodGet, "/authenticate", dynamic.ThenFunc(app.auth))
	router.Handler(http.MethodGet, "/auth/google/callback", dynamic.ThenFunc(app.authCallback))
	router.Handler(http.MethodGet, "/signup", dynamic.ThenFunc(app.signUp))
	router.Handler(http.MethodPost, "/signup", dynamic.ThenFunc(app.postSignUp))
	router.Handler(http.MethodPost, "/logout", dynamic.ThenFunc(app.logout))
	//Leagues
	router.Handler(http.MethodPost, "/leagues/create", dynamic.ThenFunc(app.createLeaguebyApiIdAndSeasonPost))
	// *RapidApi
	router.Handler(http.MethodGet, "/rapid-api/leagues", dynamic.ThenFunc(app.getRapidApiLeaguesbyApiIdAndSeason))
	// *Factos
	// @returns: []FactosResponse
	router.Handler(http.MethodGet, "/factos/latest/:userId", dynamic.ThenFunc(app.viewLatestFactosByUserId))
	router.Handler(http.MethodGet, "/factos/all/:userId", dynamic.ThenFunc(app.viewAllFactosByUserId))
	router.Handler(http.MethodGet, "/factos/league/:leagueId/:userId", dynamic.ThenFunc(app.viewAllLeagueFactosByUserId))
	router.Handler(http.MethodGet, "/factos/league-latest/:leagueId/:userId", dynamic.ThenFunc(app.viewLatestLeagueFactosByUserId))
	router.Handler(http.MethodPost, "/new-facto/:matchId", dynamic.ThenFunc(app.createFacto))
	router.Handler(http.MethodPost, "/edit-facto/:factoId", dynamic.ThenFunc(app.editFacto))
	// *Fixtures
	// @returns: []FixturesResponse
	router.Handler(http.MethodGet, "/fixtures/latest/:leagueId", dynamic.ThenFunc(app.viewLatestFixturesbyLeagueId))
	router.Handler(http.MethodGet, "/fixtures/upcoming/:leagueId", dynamic.ThenFunc(app.viewUpcomingFixturesbyLeagueId))

	return router
}
