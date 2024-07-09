package models

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
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
	Response []struct {
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
				Elapsed any    `json:"elapsed"`
			} `json:"status"`
		} `json:"fixture"`
		League struct {
			ID      int    `json:"id"`
			Name    string `json:"name"`
			Country string `json:"country"`
			Logo    string `json:"logo"`
			Flag    any    `json:"flag"`
			Season  int    `json:"season"`
			Round   string `json:"round"`
		} `json:"league"`
		Teams struct {
			Home struct {
				ID     int    `json:"id"`
				Name   string `json:"name"`
				Logo   string `json:"logo"`
				Winner any    `json:"winner"`
			} `json:"home"`
			Away struct {
				ID     int    `json:"id"`
				Name   string `json:"name"`
				Logo   string `json:"logo"`
				Winner any    `json:"winner"`
			} `json:"away"`
		} `json:"teams"`
		Goals struct {
			Home any `json:"home"`
			Away any `json:"away"`
		} `json:"goals"`
		Score struct {
			Halftime struct {
				Home any `json:"home"`
				Away any `json:"away"`
			} `json:"halftime"`
			Fulltime struct {
				Home any `json:"home"`
				Away any `json:"away"`
			} `json:"fulltime"`
			Extratime struct {
				Home any `json:"home"`
				Away any `json:"away"`
			} `json:"extratime"`
			Penalty struct {
				Home any `json:"home"`
				Away any `json:"away"`
			} `json:"penalty"`
		} `json:"score"`
	} `json:"response"`
}

func GetFixtures(params map[string]string) (*FixtureResponse, error) {
	u, err := url.Parse(os.Getenv("RAPIDAPI_FOOTBALL_URL") + "/fixtures")
	if err != nil {
		return nil, err
	}
	query := u.Query()
	for k, v := range params {
		query.Set(k, v)
	}

	method := "GET"

	client := &http.Client{}

	req, err := http.NewRequest(method, "", nil)

	if err != nil {
		return nil, err
	}
	req.Header.Add("x-rapidapi-key", os.Getenv("RAPIDAPI_KEY"))
	req.Header.Add("x-rapidapi-host", os.Getenv("RAPIDAPI_HOST"))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var fixtures *FixtureResponse
	err = json.Unmarshal(body, &fixtures)
	if err != nil {
		return nil, err
	}

	return fixtures, nil
}
