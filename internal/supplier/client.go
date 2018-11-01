package supplier

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Cli ...
type Cli interface {
	MakeRequest(url string) (Response, error)
}

// Response defines the standard response structure from suppliers
type Response map[string]float64

// Client ...
type Client struct{}

// NewClient ...
func NewClient() (Cli, error) {
	return Client{}, nil
}

// MakeRequest handles the request to suppliers
func (c Client) MakeRequest(url string) (Response, error) {
	timeout := time.Duration(15 * time.Second)
	httpClient := http.Client{Timeout: timeout}

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("the request to %s returned HTTP status code: %s", url, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return Response{}, err
	}
	return response, nil
}
