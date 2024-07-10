package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/jantoniogonzalez/factos/internal/models"
	"github.com/justinas/nosurf"
)

type templateData struct {
	// We need to get the fixtures and the facto too...
	Fixtures         []models.Response
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

func newTemplateCache() (map[string]*template.Template, error) {
	paths, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	cachedFiles := make(map[string]*template.Template)

	for _, page := range paths {
		pagename := filepath.Base(page)

		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			"./ui/html/partials/subnav.tmpl",
			page,
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cachedFiles[pagename] = ts
	}

	return cachedFiles, nil
}
