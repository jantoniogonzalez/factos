package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/jantoniogonzalez/factos/internal/api"
	"github.com/jantoniogonzalez/factos/internal/constants"
	"github.com/jantoniogonzalez/factos/internal/models"
)

// * Authentication
// The aim of the function is to send the Google URL to sign up or login
func (app *application) auth(w http.ResponseWriter, r *http.Request) {
	type ResponseData struct {
		URL string `json:"url"`
	}

	var response api.Response

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

	response = api.Response{
		Status:  constants.StatusSuccess,
		Message: "url to google login",
		Error:   "",
		Data:    ResponseData{url},
	}

	if err := app.writeJSON(w, 200, response, nil); err != nil {
		app.serverError(w, err, "Failed to write to JSON")
		return
	}
}

// After deciding to accept or deny the sign up request, do token exchange and get info
func (app *application) authCallback(w http.ResponseWriter, r *http.Request) {
	type UserInfo struct {
		GoogleId string `json:"id"`
		Email    string `json:"email"`
	}

	type ResponseData struct {
		GoogleId string `json:"googleId"`
		Username string `json:"username"`
	}

	var response api.Response

	queryError := r.URL.Query().Get("error")
	if queryError != "" {
		app.serverError(w, api.ErrGoogleLoginFailed, queryError)
		return
	}

	reqState := r.URL.Query().Get("state")
	if reqState == "" {
		app.serverError(w, api.ErrReqMissingState, "Missing state in request")
		return
	}

	sessionState := app.sessionManager.Get(r.Context(), "state")

	// Either request is malicious or something went wrong
	if sessionState != reqState {
		app.serverError(w, api.ErrStatesMismatch, "State mismatch. Possible CSRF attack")
		return
	}

	googleCode := r.URL.Query().Get("code")
	if googleCode == "" {
		app.serverError(w, api.ErrMissingGoogleCode, "No google code present.")
		return
	}

	tok, err := app.googleoauthconf.Exchange(r.Context(), googleCode)

	if err != nil {
		app.serverError(w, err, "Failed token exchange")
		return
	}

	client := app.googleoauthconf.Client(r.Context(), tok)

	res, err := client.Get(os.Getenv("GOOGLE_API_USERINFO"))

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
		// Send to new user creation
		app.sessionManager.Put(r.Context(), "googleId", userInfo.GoogleId)

		response = api.Response{
			Status:  constants.StatusRedirect,
			Message: "New User",
			Data: ResponseData{
				GoogleId: userInfo.GoogleId,
				Username: "",
			},
		}

		err = app.writeJSON(w, http.StatusOK, response, nil)
		if err != nil {
			app.serverError(w, err, "Failed to write json")
		}
		return
	}

	app.logger.Debug("User information from DB Get call",
		"Id", user.Id,
		"Username", user.Username,
	)

	app.sessionManager.Put(r.Context(), "username", user.Username)

	response = api.Response{
		Status:  constants.StatusSuccess,
		Message: "Successful user login",
		Data: ResponseData{
			Username: user.Username,
		},
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverError(w, err, "Failed to write json")
	}
}

func (app *application) postSignUp(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err, "Failed to parse form")
		return
	}

	var newUser models.User
	err = app.decoder.Decode(&newUser, r.PostForm)
	if err != nil {
		app.serverError(w, err, "Failed to decode post form")
		return
	}

	_, err = app.users.Insert(newUser.Username, newUser.GoogleId)

	if err == models.ErrDuplicateUsername {

		return
	}

	// googleId := app.sessionManager.Pop(r.Context(), "googleId")
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	app.sessionManager.Clear(r.Context())
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
