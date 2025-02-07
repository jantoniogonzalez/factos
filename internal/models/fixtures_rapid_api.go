package models

import (
	"time"
)

type FixtureResponse struct {
	Parameters struct {
		League string `json:"league"`
		Season string `json:"season"`
		Next   string `json:"next"`
	} `json:"parameters"`
	Errors  []any `json:"errors"`
	Results int   `json:"results"`
	Paging  struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"paging"`
	Response []*Response `json:"response"`
}

type Response struct {
	Fixture struct {
		ID        int       `json:"id"`
		Referee   any       `json:"referee"`
		Timezone  string    `json:"timezone"`
		Date      time.Time `json:"date"`
		Timestamp int       `json:"timestamp"`
		Periods   struct {
			First  any `json:"first"`
			Second any `json:"second"`
		} `json:"periods"`
		Venue struct {
			ID   any    `json:"id"`
			Name string `json:"name"`
			City string `json:"city"`
		} `json:"venue"`
		Status struct {
			Long    string `json:"long"`
			Short   string `json:"short"`
			Elapsed int    `json:"elapsed"`
		} `json:"status"`
	} `json:"fixture"`
	League LeagueResponse `json:"league"`
	Teams  struct {
		Home struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Logo   string `json:"logo"`
			Winner bool   `json:"winner"`
		} `json:"home"`
		Away struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Logo   string `json:"logo"`
			Winner bool   `json:"winner"`
		} `json:"away"`
	} `json:"teams"`
	Goals struct {
		Home int `json:"home"`
		Away int `json:"away"`
	} `json:"goals"`
	Score struct {
		Halftime struct {
			Home int `json:"home"`
			Away int `json:"away"`
		} `json:"halftime"`
		Fulltime struct {
			Home int `json:"home"`
			Away int `json:"away"`
		} `json:"fulltime"`
		Extratime struct {
			Home int `json:"home"`
			Away int `json:"away"`
		} `json:"extratime"`
		Penalty struct {
			Home int `json:"home"`
			Away int `json:"away"`
		} `json:"penalty"`
	} `json:"score"`
}

type LeagueResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Logo    string `json:"logo"`
	Flag    any    `json:"flag"`
	Season  int    `json:"season"`
	Round   string `json:"round"`
}
