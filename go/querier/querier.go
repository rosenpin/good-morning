package querier

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Querier queries a server for result and returns it
type Querier interface {
	// Query queries
	Query(url string) (interface{}, error)
}

// JSONQuerier queries JSON results
type JSONQuerier struct {
}

// Query queries the provided URL and returns its JSON result
func (q JSONQuerier) Query(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
