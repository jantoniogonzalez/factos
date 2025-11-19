package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/jantoniogonzalez/factos/internal/api"
	"github.com/jantoniogonzalez/factos/internal/constants"
	"github.com/jantoniogonzalez/factos/internal/models"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, data api.Response, headers http.Header) error {
	b, err := json.Marshal(data)

	if err != nil {
		return err
	}

	b = append(b, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)

	return nil
}

func (app *application) serverError(w http.ResponseWriter, err error, msg string) {
	app.logger.Error("Server Error",
		"message", msg,
		"error", err,
	)

	response := api.Response{
		Status:  constants.StatusError,
		Message: msg,
		Error:   err.Error(),
		Data:    nil,
	}

	_ = app.writeJSON(w, http.StatusInternalServerError, response, nil)
}

func (app *application) clientError(w http.ResponseWriter) {

}

func (app *application) generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// Returns models.ErrNoRecord if user does not exists
func (app *application) userExists(googleId string) (*models.User, error) {
	user, err := app.users.Get(googleId)
	if err != nil {
		if err == models.ErrNoRecord {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return user, err
}
