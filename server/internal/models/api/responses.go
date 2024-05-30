package api

import (
	"time"
)

type FactoResponse struct {
	Code          int
	MatchId       string
	GoalsHomeTeam int32
	GoalsAwayTeam int32
	LastModified  time.Time
}
