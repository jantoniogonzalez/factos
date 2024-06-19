package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	dynamic := alice.New()

	router.Handler(http.MethodGet, "/factos/:id", dynamic.ThenFunc(app.viewFactosById))
	router.Handler(http.MethodGet, "/factos/create", dynamic.ThenFunc(app.createFactos))
	router.Handler(http.MethodPost, "/factos/create", dynamic.ThenFunc(app.createFactosPost))

	return router
}
