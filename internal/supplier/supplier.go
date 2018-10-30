package supplier

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Response map[string]float64

func MakeRequest(url string) (Response, error) {
	resp, _ := http.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)

	var response Response
	err := json.Unmarshal(body, &response)
	if err != nil {
		return Response{}, err
	}
	return response, nil
}
