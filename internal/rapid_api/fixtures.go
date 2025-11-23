package rapidapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/jantoniogonzalez/factos/internal/models"
)

func GetFixturesRapidApi(params map[string]string) (*models.FixtureResponse, error) {
	u, err := url.Parse(os.Getenv("RAPIDAPI_FOOTBALL_URL") + "/fixtures")
	if err != nil {
		return nil, err
	}
	query := u.Query()
	for k, v := range params {
		query.Set(k, v)
	}
	u.RawQuery = query.Encode()
	fmt.Printf("Url we are looking for is %v\n", u.String())

	method := "GET"

	client := &http.Client{}

	req, err := http.NewRequest(method, u.String(), nil)

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

	fmt.Printf("The returned body is: %v\n", string(body))

	var fixtures *models.FixtureResponse
	err = json.Unmarshal(body, &fixtures)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Response in json: %v\n", fixtures)
	return fixtures, nil
}
