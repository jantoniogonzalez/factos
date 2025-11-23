package models

import (
	"time"
)

type Factos struct {
	Id            int       `json:"id"`
	MatchId       int       `json:"matchId"`
	LastModified  time.Time `json:"lastModified"`
	Created       time.Time `json:"created"`
	UserId        int       `json:"userId"`
	ExtraTime     bool      `json:"extraTime"`
	Penalties     bool      `json:"penalties"`
	HomeGoals     int       `json:"homeGoals"`
	AwayGoals     int       `json:"awayGoals"`
	HomePenalties int       `json:"homePenalties"`
	AwayPenalties int       `json:"awayPenalties"`
	Result        string    `json:"result"`
}
