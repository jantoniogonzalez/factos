package models

import "time"

type Fixture struct {
	ID               int       `json:"id"`
	ApiMatchId       int       `json:"apiMatchId"`
	Date             time.Time `json:"time"`
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
}

func GetFixtures(params map[string]string) (*FixtureResponse, error) {
	// Here we would have to get the parameters and make the sql query based on that.
	// We have to be careful of the parameters that are sent, we have to clear them to
	// avoid a SQL injection.
	return nil, nil
}
