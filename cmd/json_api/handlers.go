package main

import (
	"net/http"
)

// * Authentication
func (app *application) auth(w http.ResponseWriter, r *http.Request) {
	// Google url
	url := app.googleoauthconf.AuthCodeURL("state")
	app.logger.Debug("AuthCodeURL Generated",
		"url", url,
	)

	response := map[string]string{
		"url": url,
	}

	if err := app.writeJSON(w, 200, response, nil); err != nil {
		app.serverError(w, err, "Failed to marshal to JSON")
	}
}

func (app *application) authCallback(w http.ResponseWriter, r *http.Request) {

}

func (app *application) viewSignUp(w http.ResponseWriter, r *http.Request) {
	// Send the user to the page where they can add a new username
}

func (app *application) postSignUp(w http.ResponseWriter, r *http.Request) {
	// Get the username and create a new account with user's email and username
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {

}

// * Factos
func (app *application) viewLatestFactosByUserId(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement this function
	//app.factos.GetByUser();
}

func (app *application) viewAllFactosByUserId(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement this function
}

func (app *application) viewAllLeagueFactosByUserId(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement this function
}

func (app *application) viewLatestLeagueFactosByUserId(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement this function
}

func (app *application) createFacto(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement this function
}

func (app *application) editFacto(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement this function
}

func (app *application) viewLatestFixturesbyLeagueId(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement this function
}

func (app *application) viewUpcomingFixturesbyLeagueId(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement this function
}

func (app *application) errorPage(w http.ResponseWriter, r *http.Request) {}
