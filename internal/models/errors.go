package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")

	ErrDuplicateUsername = errors.New("models: duplicate username")

	ErrDuplicateGoogleId = errors.New("models: duplicate google account")
)
