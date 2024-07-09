package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

type templateData struct {
	// We need to get the fixtures and the facto too...
	Fixtures         any
	Subnav           bool
	LoggedIn         bool
	LoggedInUsername string
	Form             any
	CSRFToken        string
}

func (app *application) newTemplateData(r *http.Request, hasSubnav bool) *templateData {
	fmt.Printf("Is user logged in? %v\nWhat is their username? %v\n", app.sessionManager.Exists(r.Context(), "authenticatedUsername"), app.sessionManager.GetString(r.Context(), "authenticatedUsername"))
	return &templateData{
		LoggedIn:         app.sessionManager.Exists(r.Context(), "authenticatedUsername"),
		LoggedInUsername: app.sessionManager.GetString(r.Context(), "authenticatedUsername"),
		Subnav:           hasSubnav,
		CSRFToken:        nosurf.Token(r),
	}
}
