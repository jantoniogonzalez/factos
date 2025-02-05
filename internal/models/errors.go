package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")

	// Fixtures
	ErrNoMatchingMatchStatusShort = errors.New("models: no matching match status short")

	ErrDuplicateApiMatchId = errors.New("models: duplicate fixture with this apiMatchId")

	// Users
	ErrDuplicateUsername = errors.New("models: duplicate username")

	ErrDuplicateGoogleId = errors.New("models: duplicate google account")
)
