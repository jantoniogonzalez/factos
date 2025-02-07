package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jantoniogonzalez/factos/internal/models"
	rapidapi "github.com/jantoniogonzalez/factos/internal/rapid_api"
	"github.com/jantoniogonzalez/factos/internal/validator"
)

type googleData struct {
	Id             string `form:"id"`
	Email          string `form:"email"`
	Verified_email bool   `form:"verified_email"`
	Picture        string `form:"picture"`
}

type userCreateForm struct {
	Username string `form:"username"`
	GoogleId string `form:"googleId"`
	validator.Validator
}

// TODO: We gotta add the other fields into FactoCreateForm
type factoCreateForm struct {
	MatchId   int  `form:"matchId"`
	HomeGoals int  `form:"homeGoals"`
	AwayGoals int  `form:"awayGoals"`
	ExtraTime bool `form:"extraTime"`
	Penalties bool `form:"penalties"`
	//TODO: Add a validator for the FactosCreateForm
}

// TODO: Add a fixtureCreateForm with validation

func (app *application) viewLandingPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r, false)
	app.render(w, "landing.tmpl", data)
}

func (app *application) viewFactosById(w http.ResponseWriter, r *http.Request) {

}

// We might not need this altogether
func (app *application) createFactos(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r, false)

	path := strings.Split(r.URL.Path, "/")
	params := make(map[string]string)
	params["id"] = path[3]

	// TODO: remove rapidapi thing
	fixturesResponse, err := rapidapi.GetFixturesRapidApi(params)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Before displaying the data, check if the match happened already?
	// What if they try to change the path
	fixture := fixturesResponse.Response[0]

	if !app.BeforeDate(fixture.Fixture.Date) {
		// Make a page that shows the game already happened
		return
	}

	// Would we need to check if a facto was already created??
	// Cuz then we could reuse the createFactos with the Edit Factos

	// Pass matchId in data form
	var form factoCreateForm
	form.MatchId = fixture.Fixture.ID

	data.Form = form
	// app.render(w, "factosModal.tmpl", data)
}

func (app *application) createFactosPost(w http.ResponseWriter, r *http.Request) {
	// Here we could add the stuff needed...
	var form factoCreateForm

	err := app.decodePostForm(r, form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	userId := app.sessionManager.GetInt(r.Context(), "userId")

	newFacto := &models.Factos{
		MatchId:   form.MatchId,
		HomeGoals: form.HomeGoals,
		AwayGoals: form.AwayGoals,
		UserId:    userId,
		ExtraTime: form.ExtraTime,
		Penalties: form.Penalties,
	}
	// TODO: Validate Facto

	_, err = app.factos.InsertOne(newFacto)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) viewTournamentPredictions(w http.ResponseWriter, r *http.Request) {

}

func (app *application) viewTournamentResults(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	params := make(map[string]string)
	params["league"] = path[2]
	params["season"] = path[3]
	params["last"] = "10"
	fmt.Printf("In viewTournamentResults, looking for league: %v, season: %v, last: %v\n", params["league"], params["season"], params["last"])

	res, err := rapidapi.GetFixturesRapidApi(params)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if len(res.Errors) > 0 {
		app.notFound(w)
		return
	}

	id, err := strconv.Atoi(path[2])
	if err != nil {
		app.notFound(w)
		return
	}
	season, err := strconv.Atoi(path[3])
	if err != nil {
		app.notFound(w)
		return
	}

	league := &models.LeagueResponse{
		ID:     id,
		Season: season,
	}

	data := app.newTemplateData(r, true)
	data.Fixtures = res.Response
	data.League = league
	app.render(w, "matches.tmpl", data)
}

func (app *application) viewTournamentFutureFixtures(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	params := make(map[string]string)
	params["league"] = path[2]
	params["season"] = path[3]
	params["next"] = "10"
	fmt.Printf("In viewTournamentResults, looking for league: %v, season: %v, next: %v\n", params["league"], params["season"], params["next"])

	res, err := rapidapi.GetFixturesRapidApi(params)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if len(res.Errors) > 0 {
		app.notFound(w)
		return
	}

	id, err := strconv.Atoi(path[2])
	if err != nil {
		app.notFound(w)
		return
	}
	season, err := strconv.Atoi(path[3])
	if err != nil {
		app.notFound(w)
		return
	}

	league := &models.LeagueResponse{
		ID:     id,
		Season: season,
	}

	data := app.newTemplateData(r, true)
	data.Fixtures = res.Response
	data.League = league
	app.render(w, "matches.tmpl", data)
}

func (app *application) auth(w http.ResponseWriter, r *http.Request) {
	googleUrl := app.oauthConfig.AuthCodeURL("state")

	parsedUrl, err := url.Parse(googleUrl)
	if err != nil {
		app.serverError(w, err)
		return
	}

	values := parsedUrl.Query()
	values.Set("prompt", "select_account")
	parsedUrl.RawQuery = values.Encode()

	fmt.Printf("The url is: %s\n", parsedUrl.String())
	http.Redirect(w, r, parsedUrl.String(), http.StatusSeeOther)
}

func (app *application) authCallback(w http.ResponseWriter, r *http.Request) {
	urlA := r.URL.Query()

	code := urlA.Get("code")
	if len(code) == 0 {
		// If user decides to cancel
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tok, err := app.oauthConfig.Exchange(context.TODO(), code)
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

	googleData := &googleData{}

	err = json.Unmarshal(userData, googleData)
	if err != nil {
		app.serverError(w, err)
		return
	}

	user := app.checkUserExistsByGoogleId(googleData.Id)

	fmt.Printf("UserData is %v and also %v \n", string(userData), googleData.Id)

	if user == nil {
		// New Account
		app.sessionManager.Put(r.Context(), "googleId", googleData.Id)
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return
	}

	// Login
	app.sessionManager.Put(r.Context(), "authenticatedUsername", user.Username)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	// Create session
}

func (app *application) viewSignUp(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r, false)
	data.Form = userCreateForm{}

	app.render(w, "signup_modal.tmpl", data)
}

func (app *application) postSignUp(w http.ResponseWriter, r *http.Request) {
	var form userCreateForm

	err := app.decodePostForm(r, form)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	googleId := app.sessionManager.GetString(r.Context(), "googleId")

	fmt.Printf("Username entered is: %v\nGoogleId entered is: %v\n", username, googleId)
	userForm := userCreateForm{
		Username: username,
		GoogleId: googleId,
	}

	userForm.ValidateField(validator.NotEmpty(userForm.GoogleId), "googleId", "No Google account was found.\nClick on 'Sign up' to select a Google account.")
	userForm.ValidateField(validator.NotEmpty(userForm.Username), "username", "This field cannot be blank")
	userForm.ValidateField(validator.MaxCharacters(userForm.Username, 32), "username", "This field cannot exceed 32 characters")

	if !userForm.Valid() {
		data := app.newTemplateData(r, false)
		data.Form = userForm
		app.render(w, "signup_modal.tmpl", data)
		return
	}

	fmt.Printf("Form is valid")

	// make db calls
	userId, err := app.users.Insert(username, googleId)
	if err != nil {
		// Check if user is unique
		if errors.Is(err, models.ErrDuplicateUsername) {
			userForm.AddFieldError("username", "This username is already in use")
			data := app.newTemplateData(r, false)
			data.Form = userForm
			app.render(w, "signup_modal.tmpl", data)
			return
		}
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "googleId")
	// Not displaying authenticated username
	app.sessionManager.Put(r.Context(), "authenticatedUsername", userForm.Username)
	app.sessionManager.Put(r.Context(), "userId", userId)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	app.sessionManager.Remove(r.Context(), "authenticatedUsername")
	app.sessionManager.Remove(r.Context(), "userId")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
