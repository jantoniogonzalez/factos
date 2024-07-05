package main

import "net/http"

type templateData struct {
	// We need to get the fixtures and the facto too...
	Fixtures []*string
	Subnav   bool
	LoggedIn bool
	Form     any
}

func (app *application) newTemplateData(r *http.Request, hasSubnav bool) *templateData {
	return &templateData{
		LoggedIn: app.sessionManager.Exists(r.Context(), "authenticatedUserID"),
		Subnav:   hasSubnav,
	}
}
