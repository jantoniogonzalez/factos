package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/jantoniogonzalez/factos/internal/models"
	"github.com/jantoniogonzalez/factos/internal/validator"
)

type googleData struct {
	Id             string `form:"id"`
	Email          string `form:"email"`
	Verified_email bool   `form:"verified_email"`
	Picture        string `form:"picture"`
}

type userCreateForm struct {
	Username string
	GoogleId string
	validator.Validator
}

func (app *application) viewLandingPage(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
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

	user := app.checkUserExists(googleData.Id)

	fmt.Printf("UserData is %v and also %v \n", string(userData), googleData.Id)

	if user == nil {
		// New Account
		http.Redirect(w, r, "/signup?user="+googleData.Id, http.StatusSeeOther)
		return
	}

	// Login
	http.Redirect(w, r, "/", http.StatusSeeOther)
	// Create session
}

func (app *application) viewSignUp(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/signup_modal.tmpl",
	}

	data := app.newTemplateData(r, false)
	data.Form = userCreateForm{}

	app.render(w, files, data)
}

func (app *application) postSignUp(w http.ResponseWriter, r *http.Request) {
	// Read postbody
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	username := r.Form.Get("Username")
	fmt.Printf("Username entered is: %v", username)
	googleId := r.URL.Query().Get("user")

	userForm := userCreateForm{
		Username: username,
		GoogleId: googleId,
	}

	userForm.ValidateField(validator.NotEmpty(userForm.Username), "username", "This field cannot be blank")
	userForm.ValidateField(validator.MaxCharacters(userForm.Username, 32), "username", "This field cannot exceed 32 characters")

	if !userForm.Valid() {
		data := app.newTemplateData(r, false)
		data.Form = userForm
		app.render(w, []string{"user_model.tmpl"}, data)
		return
	}

	// make db calls
	err = app.users.Insert(username, googleId)
	if err != nil {
		// Check if user is unique
		if errors.Is(err, models.ErrDuplicateUsername) {
			userForm.AddFieldError("username", "This username is already in use")
			data := app.newTemplateData(r, false)
			data.Form = userForm
			app.render(w, []string{"user_model.tmpl"}, data)
		}
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
