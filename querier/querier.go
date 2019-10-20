package querier

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"gitlab.com/rosenpin/good-morning/config"
)

const (
	request = "https://www.googleapis.com/customsearch/v1"

	// resultNumRand is the value used for selecting the random result index from the list of results
	resultNumRand = 5
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Querier queries a server for result and returns it
type Querier interface {
	// Query queris
	Query() (interface{}, error)
}

type imagesQuerier struct {
	config config.Config
}

// Query queries google images api using the configuration
func (q imagesQuerier) Query() (interface{}, error) {
	query := q.buildURL()
	fmt.Println("using query", query)
	resp, err := http.Get(query)
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

func (q imagesQuerier) buildURL() string {
	baseURL, err := url.Parse(request)
	if err != nil {
		panic("invalid base request url")
	}

	params := url.Values{}
	params.Add("key", q.config.API.APIKey)
	params.Add("cx", q.config.API.CX)
	params.Add("q", q.config.Search.BaseQuery+" "+strconv.Itoa(rand.Intn(q.config.Search.Randomness)))
	params.Add("imgSize", q.config.Image.Size)
	params.Add("start", strconv.Itoa(rand.Intn(resultNumRand)))
	params.Add("searchType", "image")
	params.Add("num", "1")
	baseURL.RawQuery = params.Encode()

	return baseURL.String()
}

// NewQuerier returns a new image querier object
func NewQuerier(config config.Config) Querier {
	return imagesQuerier{config}
}
