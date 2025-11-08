package main

import "net/http"

// * Authentication
func (app *application) auth(w http.ResponseWriter, r *http.Request) {}

func (app *application) authCallback(w http.ResponseWriter, r *http.Request) {}

func (app *application) viewSignUp(w http.ResponseWriter, r *http.Request) {}

func (app *application) postSignUp(w http.ResponseWriter, r *http.Request) {}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {}

// * Factos
func (app *application) viewLatestFactosByUserId(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement this function

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
