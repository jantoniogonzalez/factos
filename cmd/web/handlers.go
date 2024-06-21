package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (app *application) viewLandingPage(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/partials/subnav.tmpl",
		"./ui/html/pages/landing.tmpl",
	}

	data := &templateData{
		Subnav: false,
	}

	app.render(w, files, data)

}

func (app *application) viewFactosById(w http.ResponseWriter, r *http.Request) {

}

func (app *application) createFactos(w http.ResponseWriter, r *http.Request) {

}

func (app *application) createFactosPost(w http.ResponseWriter, r *http.Request) {

}

func (app *application) viewTournamentPredictions(w http.ResponseWriter, r *http.Request) {

}

func (app *application) viewTournamentResults(w http.ResponseWriter, r *http.Request) {

}

func (app *application) viewFutureFixtures(w http.ResponseWriter, r *http.Request) {

}

func (app *application) auth(w http.ResponseWriter, r *http.Request) {
	url := app.oauthConfig.AuthCodeURL("state")

	fmt.Printf("The url is: %s\n", url)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (app *application) authCallback(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()

	if url.Get("code") == "" {
		// If user decides to cancel
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	tok, err := app.oauthConfig.Exchange(context.TODO(), url.Get("code"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tok.AccessToken)
	if err != nil {
		app.serverError(w, err)
		return
	}

	userData, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("Couldn't get userdata")
	}

	fmt.Printf("UserData is %v\n", string(userData))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
