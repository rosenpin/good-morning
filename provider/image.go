package provider

import (
	"fmt"
	"io"
	"time"

	"gitlab.com/rosenpin/good-morning/caching"
	"gitlab.com/rosenpin/good-morning/config"
	"gitlab.com/rosenpin/good-morning/querier"
	"gitlab.com/rosenpin/good-morning/result"
	"gitlab.com/rosenpin/good-morning/url"
)

// ImageProvider provides images from the web using the provided querier
// It searches for images using the configuration and creates the search URL using the URL creator
type ImageProvider struct {
	querier    querier.Querier
	urlCreator url.Creator
	parser     result.Parser
	cache      caching.Cache
	config     config.Config
}

// Provide provides the image
func (p ImageProvider) Provide() (io.ReadCloser, error) {
	if p.isCacheValid() {
		return p.loadFromCache()
	}

	fmt.Println("cache invalid, reloading..")

	link, err := p.getImageURL()
	if err != nil {
		return nil, err
	}

	p.cache.Save(caching.ImageKey, link)

	return p.loadFromCache()
}

func (p ImageProvider) isCacheValid() bool {
	cacheAge, err := p.cache.Age(caching.ImageKey)
	if err == nil && cacheAge < time.Hour*time.Duration(p.config.Image.LifeSpan) {
		return true
	}

	return false
}

func (p ImageProvider) loadFromCache() (io.ReadCloser, error) {
	cache, err := p.cache.Load(caching.ImageKey)
	if err != nil {
		return nil, err
	}

	r, ok := cache.(io.ReadCloser)
	if !ok {
		return nil, fmt.Errorf("caching return invalid type")
	}
	return r, nil
}

func (p ImageProvider) getImageURL() (string, error) {
	query, err := p.urlCreator.Create(configToParams(p.config))
	if err != nil {
		return "", err
	}

	fmt.Println("sending request:", query)
	result, err := p.querier.Query(query)
	if err != nil {
		return "", err
	}

	rawParsed, err := p.parser.Parse(result)
	if err != nil {
		return "", err
	}

	link, ok := rawParsed.(string)
	if !ok {
		return "", fmt.Errorf("unexpected result returned from parser, %T:%v", link, link)
	}

	return link, nil
}

func configToParams(c config.Config) url.GoogleImagesParams {
	p := url.GoogleImagesParams{}

	p.APIKey = c.API.APIKey
	p.APICX = c.API.CX
	p.BaseQuery = c.Search.BaseQuery
	p.ImgSize = c.Image.Size
	p.Randomness = c.Search.Randomness
	p.ResultNum = 1

	return p
}

// NewImageProvider creates a new image provider object
func NewImageProvider(querier querier.Querier, urlCreator url.Creator, parser result.Parser, config config.Config, cache caching.Cache) ImageProvider {
	return ImageProvider{querier: querier, urlCreator: urlCreator, parser: parser, config: config, cache: cache}
}
