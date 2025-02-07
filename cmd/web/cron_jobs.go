package main

import rapidapi "github.com/jantoniogonzalez/factos/internal/rapid_api"

func updatePastLeagueResults(leagueId, last string) {
	params := make(map[string]string)
	params["league"] = leagueId
	params["last"] = last
	_, err := rapidapi.GetFixturesRapidApi(params)

	if err != nil {
		return
	}
}

func updateFutureLeagueResults(leagueId, next string) {
	params := make(map[string]string)
	params["league"] = leagueId
	params["next"] = next
	_, err := rapidapi.GetFixturesRapidApi(params)

	if err != nil {
		return
	}
}
