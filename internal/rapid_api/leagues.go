package rapidapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type FullLeaguesResponse struct {
	Get        string `json:"get"`
	Parameters struct {
		Id     string `json:"id"`
		Season string `json:"season"`
	} `json:"parameters"`
	Errors  []string `json:"errors"`
	Results int      `json:"results"`
	Paging  struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"paging"`
	Response []*LeaguesResponse `json:"response"`
}

type LeaguesResponse struct {
	League  League   `json:"league"`
	Country Country  `json:"country"`
	Seasons []Season `json:"seasons"`
}

type League struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Logo string `json:"logo"`
}

type Country struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Flag string `json:"flag"`
}

type Season struct {
	Year    int    `json:"year"`
	Start   string `json:"start"`
	End     string `json:"end"`
	Current bool   `json:"current"`
}

func GetLeagueByApiIdAndSeason(apiLeagueId, season string) (*[]*LeaguesResponse, error) {
	u, err := url.Parse(os.Getenv("RAPIDAPI_FOOTBALL_URL") + "/leagues")

	if err != nil {
		return nil, err
	}

	query := u.Query()

	query.Add("id", apiLeagueId)
	query.Add("season", season)

	u.RawQuery = query.Encode()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	method := "GET"

	req, err := http.NewRequest(method, u.String(), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("x-apisports-key", os.Getenv("RAPIDAPI_KEY"))
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

	var fullLeaguesResponse *FullLeaguesResponse

	fmt.Printf("Received the following response: %s\n", string(body))

	err = json.Unmarshal(body, &fullLeaguesResponse)

	if err != nil {
		return nil, err
	}

	// Check if there are any errors
	if len(fullLeaguesResponse.Errors) > 0 {
		return nil, ErrGeneric
	}

	return &fullLeaguesResponse.Response, nil
}

/*
	Now we gotta ask ourselves of how do we want to handle this API

	Q: Main Goal of API?
	A: Get a response from Rapid Api to get a League with a certain season.

	Q: What do we want the process to look like?
	A: createLeagueByX -> GetLeaguesResponse

	Q: How do we know that we are creating the response that we want?
	A: Maybe, before the post request, we have to call a get request to actually see if we are selecting the right league

	Q: How strict do we want the API to be?

	Q: What parts of the response do we need?
	A: We want to know
*/
