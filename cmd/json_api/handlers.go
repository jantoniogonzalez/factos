package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/jantoniogonzalez/factos/internal/models"
)

type UserInfo struct {
	GoogleId string `json:"id"`
	Email    string `json:"email"`
}

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
	queryError := r.URL.Query().Get("error")
	if queryError != "" {
		app.serverError(w, nil, queryError)
		return
	}

	reqState := r.URL.Query().Get("state")
	if reqState == "" {
		app.logger.Error("No state received")
		return
	}

	sessionState := app.sessionManager.Get(r.Context(), "state")

	// Either request is malicious or something went wrong
	if sessionState != reqState {
		app.serverError(w, nil, "State mismatch. Possible CSRF attack")
		return
	}

	googleCode := r.URL.Query().Get("code")
	if googleCode == "" {
		app.logger.Error("No google code received")
		return
	}

	url := r.URL.Query()
	// Need to check the state and the code
	app.logger.Debug("Got url query stuff",
		"url", url,
	)

	tok, err := app.googleoauthconf.Exchange(r.Context(), googleCode)

	if err != nil {
		app.serverError(w, err, "Failed token exchange")
		return
	}

	client := app.googleoauthconf.Client(r.Context(), tok)

	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

	app.logger.Debug("Response status",
		"status", res.Status,
	)

	app.logger.Debug("Response headers",
		"headers", res.Header,
	)

	if err != nil {
		app.serverError(w, err, "Failed retrieve user information from Google api")
		return
	}

	defer res.Body.Close()

	resData, err := io.ReadAll(res.Body)

	if err != nil {
		app.serverError(w, err, "Failed to read response body")
		return
	}

	var userInfo UserInfo
	err = json.Unmarshal(resData, &userInfo)
	if err != nil {
		app.serverError(w, err, "Failed to unmarshal data")
		return
	}

	app.logger.Debug("User Info",
		"Google Id", userInfo.GoogleId,
		"Email", userInfo.Email,
	)

	user, err := app.users.Get(userInfo.GoogleId)
	if err != nil && err != models.ErrNoRecord {
		app.serverError(w, err, "Failed DB call to Get Users")
		return
	}

	if err == models.ErrNoRecord {
		//Create a new user
		app.logger.Debug("Need to create new user")
	} else {
		app.logger.Debug("User information from DB Get call",
			"Id", user.Id,
			"Username", user.Username,
		)
	}
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
