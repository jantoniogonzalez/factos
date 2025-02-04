package models

import "time"

type League struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	ApiLeagueId  int       `json:"apiLeagueId"`
	Country      string    `json:"country"`
	Season       int       `json:"season"`
	Logo         string    `json:"logo"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
}
