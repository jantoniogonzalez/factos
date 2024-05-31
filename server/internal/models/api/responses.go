package api

import (
	"github.com/jantoniogonzalez/factos/internal/models"
)

type DefaultResponse struct {
	Code    int
	Message string
}

type FactoResponse struct {
	Code  int
	Facto models.Facto
}

type ErrorResponse struct {
	Code    int
	Message string
}
