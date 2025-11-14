package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.errorPage)

	dynamic := alice.New(addSecurityHeaders, app.sessionManager.LoadAndSave, noSurf)

	// Don't we need the home route? Like would we check if the user is logged in then? And add info about that?
	/*
		Don't we need the home route? Like would we check if the user is logged in then? And add info about that?
		The session can only be accessed from teh backend, i think, so even if we saved the username or info we want to include,
		i dont think the frontend would be able to access it?
		Maybe the home route can just connect with the server and check if the user is logged in, this can be done on a loading screen.
	*/

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
