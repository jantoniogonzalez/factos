package rapidapi

import "errors"

var (
	ErrGeneric        = errors.New("Something went wrong with the Rapid Api call")
	ErrResponseErrors = errors.New("There may be one or multiple response errors")
)
