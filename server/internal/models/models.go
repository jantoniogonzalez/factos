package models

import "time"

type User struct {
	Username string
	Factos   int64
	UserId   string
}

type Facto struct {
	MatchId       string
	GoalsHomeTeam int32
	GoalsAwayTeam int32
	LastModified  time.Time
	UserId        string
	ExtraTime     bool
	Penalties     bool
}

type Friend struct {
	UserId1         string
	UserId2         string
	DateEstablished time.Time
}

type FriendRequest struct {
	Sender   string
	Receiver string
	DataSent time.Time
	Status   int32
}

type Challenge struct {
	Challenger   string
	Opponent     string
	Status       int16
	FactosBetted int32
}
