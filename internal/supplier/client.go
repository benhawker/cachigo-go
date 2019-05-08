package supplier

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
    "sync"
)

// Cli ...
type Cli interface {
	MakeRequest(url string) (Response, error)
    MakeRequestXX(name, url string, wg *sync.WaitGroup, responseCh chan ResponseWithError) ()
}

type ResponseWithError struct {
    Name string
    Body map[string]float64
    Error error
}

// Response defines the standard response structure from suppliers
type Response map[string]float64

// Client ...
type Client struct{}

// NewClient ...
func NewClient() (Cli, error) {
	return Client{}, nil
}

func (c Client) MakeRequestXX(name, url string, wg *sync.WaitGroup, responseCh chan ResponseWithError) {
    wg.Add(1)

    timeout := time.Duration(15 * time.Second)
    httpClient := http.Client{Timeout: timeout}

    resp, err := httpClient.Get(url)
    if err != nil {
        // body <- b
        // return nil, err
        responseCh <- ResponseWithError{name, map[string]float64{}, err}
    }

    if resp.StatusCode != 200 {
        // return nil, fmt.Errorf("the request to %s returned HTTP status code: %s", url, resp.Status)
        responseCh <- ResponseWithError{name, map[string]float64{}, fmt.Errorf("the request to %s returned HTTP status code: %s", url, resp.Status)}
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        // return nil, err
        responseCh <- ResponseWithError{name, map[string]float64{}, err}
    }

    var response map[string]float64
    err = json.Unmarshal(body, &response)
    if err != nil {
        // return Response{}, err
        responseCh <- ResponseWithError{name, map[string]float64{}, err}
    }

    responseCh <- ResponseWithError{name, map[string]float64{}, nil}
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
