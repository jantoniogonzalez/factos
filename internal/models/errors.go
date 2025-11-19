package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")

	// Fixtures
	ErrNoMatchingMatchStatusShort = errors.New("models: no matching match status short")

	ErrDuplicateApiMatchId = errors.New("models: duplicate fixture with this apiMatchId")

	// Leagues
	ErrDuplicateApiLeagueId = errors.New("models: duplicate fixture with this apiLeagueId")
	// Teams
	ErrDuplicateApiTeamId = errors.New("models: duplicate fixture with this apiTeamId")
	// Users
	ErrDuplicateUsername = errors.New("models: duplicate username")

	ErrDuplicateGoogleId = errors.New("models: duplicate google account")

	ErrDuplicatePrimaryKey = errors.New("models: duplicate of primary key")
)
