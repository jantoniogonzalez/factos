package models

import "time"

type Team struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Logo         string    `json:"logo"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
	ApiTeamId    int       `json:"apiTeamId"`
}
