package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

type templateData struct {
	// We need to get the fixtures and the facto too...
	Fixtures         []*string
	Subnav           bool
	LoggedIn         bool
	LoggedInUsername string
	Form             any
	CSRFToken        string
}

func (app *application) newTemplateData(r *http.Request, hasSubnav bool) *templateData {
	return &templateData{
		LoggedIn:         app.sessionManager.Exists(r.Context(), "authenticatedUsername"),
		LoggedInUsername: app.sessionManager.GetString(r.Context(), "authenticatedUsername"),
		Subnav:           hasSubnav,
		CSRFToken:        nosurf.Token(r),
	}
}
