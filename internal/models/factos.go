package models

import "time"

type Factos struct {
	MatchId       string
	GoalsHomeTeam int
	GoalsAwayTeam int
	LastModified  time.Time
	Created       time.Time
	UserId        string
	ExtraTime     bool
	Penalties     bool
}

func Insert(mathId string, goalsHomeTeam, goalsAwayTeam int, extraTime, penalties bool) error {

}

func Get() {

}

func Edit() {

}
