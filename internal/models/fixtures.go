package models

import "time"

type Fixture struct {
	ID               int       `json:"id"`
	Date             time.Time `json:"date"`
	LeagueId         int       `json:"leagueId"`
	HomeGoals        int       `json:"homeGoals"`
	AwayGoals        int       `json:"awayGoals"`
	HomePenalties    int       `json:"homePenalties"`
	AwayPenalties    int       `json:"awayPenalties"`
	HomeId           int       `json:"homeId"`
	AwayId           int       `json:"awayId"`
	Created          time.Time `json:"created"`
	LastModified     time.Time `json:"lastModified"`
	MatchStatusShort string    `json:"matchStatusShort"`
	ApiMatchId       int       `json:"apiMatchId"`
}
