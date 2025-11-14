package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
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

	w.WriteHeader(http.StatusInternalServerError)
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
