package main

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/jantoniogonzalez/factos/internal/models"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, page string, data *templateData) {
	ts := app.cachedFiles[page]
	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) checkUserExistsByGoogleId(googleId string) *models.User {
	user, err := app.users.Get(googleId)

	if err != nil {
		return nil
	}

	return user
}

func (app *application) parseForm() {

}
