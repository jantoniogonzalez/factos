package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/go-playground/form/v4"
	"github.com/jantoniogonzalez/factos/internal/api"
	"github.com/jantoniogonzalez/factos/internal/constants"
	"github.com/jantoniogonzalez/factos/internal/models"
	rapidapi "github.com/jantoniogonzalez/factos/internal/rapid_api"
	"github.com/jantoniogonzalez/factos/internal/validator"
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
		GoogleId string `json:"googleId,omitempty"`
		Username string `json:"username,omitempty"`
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

	sessionState := app.sessionManager.Pop(r.Context(), "state")

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

	user, err := app.userExists(userInfo.GoogleId)

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
			},
		}

		err = app.writeJSON(w, http.StatusOK, response, nil)
		if err != nil {
			app.serverError(w, err, "Failed to write json")
		}
		return
	}

	if err != nil {
		app.serverError(w, err, "Failed DB call to Get Users")
		return
	}

	app.logger.Debug("User information from DB Get call",
		"Id", user.Id,
		"Username", user.Username,
	)

	app.sessionManager.Put(r.Context(), "authenticatedUsername", user.Username)

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

func (app *application) signUp(w http.ResponseWriter, r *http.Request) {
	app.logger.Debug("Received signup from Get")
	app.writeJSON(w, http.StatusOK, api.Response{
		Status:  constants.StatusSuccess,
		Message: "Hola from GET /signup",
	}, nil)
}

func (app *application) postSignUp(w http.ResponseWriter, r *http.Request) {
	type CreateUserForm struct {
		Username string `form:"username"`
		GoogleId string `form:"googleId"`
		validator.Validator
	}

	var response api.Response

	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err, "Failed to parse form")
		return
	}

	var newUser CreateUserForm
	err = app.decoder.Decode(&newUser, r.Form)

	var invalidDecoderError *form.InvalidDecoderError

	if errors.As(err, &invalidDecoderError) {
		response = api.Response{
			Status:  constants.StatusError,
			Message: "Parsing error",
			Error:   "Form format does not match CreateUserForm",
		}
		app.writeJSON(w, http.StatusBadRequest, response, nil)
		return
	}

	if err != nil {
		app.serverError(w, err, "Failed to decode post form")
		return
	}

	app.logger.Debug("Form Received", "googleId", newUser.GoogleId, "username", newUser.Username)

	// Check google id is the same
	googleId := app.sessionManager.Get(r.Context(), "googleId")

	app.logger.Debug("Current Session", "googleId", googleId)

	if googleId != newUser.GoogleId {
		response = api.Response{
			Status:  constants.StatusError,
			Message: "Google Id in sent form and session do not match",
			Error:   "Form is invalid",
		}
		app.writeJSON(w, http.StatusBadRequest, response, nil)
		return
	}

	// Validate fields
	newUser.Validator.ValidateUsername(newUser.Username)
	if !newUser.Validator.Valid() {
		response = api.Response{
			Status:  constants.StatusError,
			Message: "Form fields are invalid for UserCreateForm",
			Error:   "One or multiple fields are invalid",
			Data:    newUser.Validator.FieldErrors,
		}
		app.writeJSON(w, http.StatusBadRequest, response, nil)
		return
	}

	_, err = app.users.Insert(newUser.Username, newUser.GoogleId)

	if err == models.ErrDuplicateGoogleId || err == models.ErrDuplicatePrimaryKey {
		response = api.Response{
			Status:  constants.StatusError,
			Message: "User already exists",
			Error:   "User already exists",
		}
		err = app.sessionManager.Clear(r.Context())

		if err != nil {
			app.serverError(w, err, "Failed to clear session")
			return
		}

		err = app.writeJSON(w, http.StatusBadRequest, response, nil)
		if err != nil {
			app.serverError(w, err, "Failed to writeJSON")
		}

		return
	}

	if err == models.ErrDuplicateUsername {
		newUser.AddFieldError("username", "Username already exists")
		response = api.Response{
			Status:  constants.StatusError,
			Message: "Form fields are invalid for UserCreateForm",
			Error:   "One of multiple fields are invalid",
			Data:    newUser.Validator.FieldErrors,
		}
		err = app.writeJSON(w, http.StatusConflict, response, nil)
		if err != nil {
			app.serverError(w, err, "Failed to writeJSON")
		}
		return
	}

	_ = app.sessionManager.Pop(r.Context(), "googleId")
	app.sessionManager.Put(r.Context(), "authenticatedUsername", newUser.Username)

	response = api.Response{
		Status:  constants.StatusSuccess,
		Message: "Sign up successful",
	}

	err = app.writeJSON(w, http.StatusCreated, response, nil)

	if err != nil {
		app.serverError(w, err, "Failed to writeJSON")
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	type ResponseData struct {
		Redirect string `json:"redirect"`
	}

	var response api.Response

	username := app.sessionManager.Pop(r.Context(), "authenticatedUsername")

	app.logger.Debug("Username", "username", username)

	if username == nil {
		response = api.Response{
			Status:  constants.StatusError,
			Message: "No user found",
		}
		err := app.writeJSON(w, http.StatusNotModified, response, nil)

		if err != nil {
			app.serverError(w, err, "Failed to writeJSON")
		}

		return
	}

	response = api.Response{
		Status:  constants.StatusSuccess,
		Message: "User logged out successfully",
		Data: ResponseData{
			Redirect: "/",
		},
	}

	err := app.writeJSON(w, http.StatusOK, response, nil)

	if err != nil {
		app.serverError(w, err, "Failed to writeJSON")
	}
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

// * Leagues
func (app *application) getRapidApiLeaguesbyApiIdAndSeason(w http.ResponseWriter, r *http.Request) {
	type ResponseData struct {
		RapidApiResponse *rapidapi.FullLeaguesResponse
	}

	apiLeagueId := r.URL.Query().Get("apiLeagueId")
	season := r.URL.Query().Get("season")

	var response api.Response

	var err error

	app.logger.Debug("Received the following from form",
		"apiLeagueId", apiLeagueId,
		"season", season)

	if len(apiLeagueId) < 1 || len(season) < 1 {
		response = api.Response{
			Status:  constants.StatusError,
			Message: "Fields cannot be empty",
			Error:   "Either ApiLeagueId or Season is empty blud",
		}
		err = app.writeJSON(w, http.StatusBadRequest, response, nil)
		if err != nil {
			app.serverError(w, err, "Failed to writeJSON")
		}
		return
	}

	apires, err := rapidapi.GetLeagueByApiIdAndSeason(apiLeagueId, season)

	if err != nil {
		response = api.Response{
			Status:  constants.StatusError,
			Message: "Received error from GetLEagueByApiIdAndSeason",
			Error:   err.Error(),
		}

		err = app.writeJSON(w, http.StatusBadRequest, response, nil)
		if err != nil {
			app.serverError(w, err, "Failed writeJSON")
		}
		return
	}

	// CHECK ERRORS

	response = api.Response{
		Status:  constants.StatusSuccess,
		Message: "Found the following responses",
		Data:    ResponseData{RapidApiResponse: apires},
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverError(w, err, "Failed writeJSON")
	}

}

func (app *application) createLeaguePost(w http.ResponseWriter, r *http.Request) {
	type CreateLeagueForm struct {
		Name        string `json:"name"`
		ApiLeagueId int    `json:"apiLeagueId"`
		Country     string `json:"country"`
		Season      int    `json:"season"`
		Logo        string `json:"logo"`
		validator.Validator
	}

	var response api.Response

	err := r.ParseForm()

	if err != nil {
		app.serverError(w, err, "Failed to ParseForm in createLeaguebyApiLeagueIdPost")
		return
	}

	var createLeagueForm CreateLeagueForm

	err = app.decoder.Decode(&createLeagueForm, r.PostForm)

	var invalidDecoderError *form.InvalidDecoderError

	if errors.As(err, &invalidDecoderError) {
		response = api.Response{
			Status:  constants.StatusError,
			Message: "Form does not provide expected format",
			Error:   err.Error(),
		}
		err = app.writeJSON(w, http.StatusBadRequest, response, nil)
		if err != nil {
			app.serverError(w, err, "Failed to writeJSON")
		}
		return
	}

	if err != nil {
		app.serverError(w, err, "Failed to decode form")
		return
	}

	app.logger.Debug("Received the following CreateLeagueForm",
		"Name", createLeagueForm.Name,
		"ApiLeagueId", createLeagueForm.ApiLeagueId,
		"Country", createLeagueForm.Country,
		"Season", createLeagueForm.Season,
		"Logo", createLeagueForm.Logo)

	// Hmmm the problem with this is that we have to remember to validate all fields, should be fine though?
	createLeagueForm.Validator.ValidateLeagueName(createLeagueForm.Name)
	createLeagueForm.Validator.ValidateLeagueCountry(createLeagueForm.Country)
	createLeagueForm.Validator.ValidateLeagueLogo(createLeagueForm.Logo)
	createLeagueForm.Validator.ValidateLeagueApiLeagueId(createLeagueForm.ApiLeagueId)
	createLeagueForm.Validator.ValidateLeagueSeason(createLeagueForm.Season)

	if !createLeagueForm.Validator.Valid() {
		response = api.Response{
			Status:  constants.StatusError,
			Message: "There may be one or several field errors",
			Error:   "Stuff",
			Data:    createLeagueForm.Validator.FieldErrors,
		}
		err = app.writeJSON(w, http.StatusBadRequest, response, nil)
		if err != nil {
			app.serverError(w, err, "Failed writeJSON")
		}
		return
	}

	newLeague := &models.League{
		Name:        createLeagueForm.Name,
		ApiLeagueId: createLeagueForm.ApiLeagueId,
		Country:     createLeagueForm.Country,
		Season:      createLeagueForm.Season,
		Logo:        createLeagueForm.Logo,
	}

	leagueId, err := app.leagues.InsertOne(newLeague)

	if err != nil {
		response = api.Response{
			Status:  constants.StatusError,
			Message: "Error in InsertOne",
			Error:   err.Error(),
		}

		if errors.Is(err, models.ErrDuplicateApiLeagueIdAndSeasonCombination) {
			response.Error = models.ErrDuplicateApiLeagueIdAndSeasonCombination.Error()
			err = app.writeJSON(w, http.StatusBadRequest, response, nil)
			if err != nil {
				app.serverError(w, err, "Failed writeJSON")
			}
			return
		}

		app.serverError(w, err, "Failed InsertOne")
		return
	}

	response = api.Response{
		Status:  constants.StatusSuccess,
		Message: "Id of inserted league",
		Data:    leagueId,
	}

	app.writeJSON(w, http.StatusOK, response, nil)
}

// * Fixtures
func (app *application) viewLatestFixturesbyLeagueId(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement this function
}

func (app *application) viewUpcomingFixturesbyLeagueId(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement this function
}

func (app *application) errorPage(w http.ResponseWriter, r *http.Request) {
	app.logger.Debug("Bad Request", "url", r.URL)
}
