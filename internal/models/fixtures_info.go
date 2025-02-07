package models

type FixtureInfo struct {
	Fixture  Fixture `json:"fixture"`
	HomeTeam Team    `json:"homeTeam"`
	AwayTeam Team    `json:"awayTeam"`
	League   League  `json:"league"`
}

func NewFixtureInfo(f Fixture, homeTeam, awayTeam Team, league League) *FixtureInfo {
	return &FixtureInfo{
		Fixture:  f,
		HomeTeam: homeTeam,
		AwayTeam: awayTeam,
		League:   league,
	}
}
