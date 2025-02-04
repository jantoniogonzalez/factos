package models

import (
	"time"
)

type Factos struct {
	Id           int
	MatchId      int
	LeagueId     int
	Season       int
	GoalsHome    int
	GoalsAway    int
	Result       int
	LastModified time.Time
	Created      time.Time
	UserId       int
	ExtraTime    bool
	Penalties    bool
}
