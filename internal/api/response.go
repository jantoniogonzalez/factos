package api

import "github.com/jantoniogonzalez/factos/internal/constants"

type Response struct {
	Status  constants.ResponseStatus `json:"status"`
	Message string                   `json:"message"`
	Error   string                   `json:"error,omitempty"`
	Data    interface{}              `json:"data,omitempty"`
}
