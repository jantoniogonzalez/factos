package main

import (
	"net/http"
)

// * Authentication
// The aim of the function is to send the Google URL to sign up or login
func (app *application) auth(w http.ResponseWriter, r *http.Request) {
	state, err := app.generateRandomState()
	if err != nil {
		app.serverError(w, err, "Failed to generate state")
		return
	}
	app.sessionManager.Put(r.Context(), "state", state)
	app.logger.Debug("Added state to session in auth/",
		"state", state,
	)

	// Google url
	url := app.googleoauthconf.AuthCodeURL(state)
	app.logger.Debug("AuthCodeURL Generated",
		"url", url,
	)

	response := map[string]string{
		"url": url,
	}

	if err := app.writeJSON(w, 200, response, nil); err != nil {
		app.serverError(w, err, "Failed to write to JSON")
		return
	}
}

// After deciding to accept or deny the sign up request, do token exchange and get info
func (app *application) authCallback(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()
	app.logger.Info("Got url query stuff",
		"url", url,
	)

	sessionState := app.sessionManager.Get(r.Context(), "state")
	app.logger.Debug("State in session",
		"sessionState", sessionState,
	)

	// when rejected:  http://localhost:4000/auth/google/callback?error=access_denied&state=state

	// tok, err := app.googleoauthconf.Exchange(r.Context(), "authorization-code")

	// if err != nil {
	// 	app.serverError(w, err, "Failed token exchange")
	// }

	// client := app.googleoauthconf.Client(r.Context(), tok)

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
