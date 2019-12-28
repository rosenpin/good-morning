package url

import (
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

// Creator creates URLs using pre defined heuristics and the provided params
type Creator interface {
	Create(params interface{}) (string, error)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	request = "https://www.googleapis.com/customsearch/v1"

	imageSearchType = "image"
)

// GoogleImagesCreator creates URLs for querying Google images
type GoogleImagesCreator struct {
}

// GoogleImagesParams is the struct providing the params for the GoogleImagesCreator
type GoogleImagesParams struct {
	APIKey     string
	APICX      string
	BaseQuery  string
	ImgSize    string
	Randomness int
	ResultNum  int
	Features   []string
}

// Create creates a Google Images url using a set of URL parameters
func (g GoogleImagesCreator) Create(rawParams interface{}) (string, error) {
	gparams, ok := rawParams.(GoogleImagesParams)
	if !ok {
		return "", fmt.Errorf("GoogleImagesCreator must take a GoogleImagesParams params")
	}

	baseURL, err := url.Parse(request)
	if err != nil {
		panic("invalid base request url")
	}

	params := &url.Values{}
	baseURL.RawQuery = g.addParams(params, gparams)

	return baseURL.String(), nil

}

func (g GoogleImagesCreator) addParams(params *url.Values, gparams GoogleImagesParams) string {
	params.Add("key", gparams.APIKey)
	params.Add("cx", gparams.APICX)
	params.Add("q", gparams.BaseQuery+" "+gparams.Features[rand.Intn(len(gparams.Features))])
	params.Add("imgSize", gparams.ImgSize)
	params.Add("start", strconv.Itoa(rand.Intn(gparams.Randomness)))
	params.Add("searchType", imageSearchType)
	params.Add("num", strconv.Itoa(gparams.ResultNum))

	return params.Encode()
}
