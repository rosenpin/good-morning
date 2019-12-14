package caching

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	ImageKey = "ImageKey"
)

type Cache interface {
	Save(key string, value interface{}) error
	Load(key string) (interface{}, error)
	Age(key string) (time.Duration, error)
}

type cacheData struct {
	path     string
	creation time.Time
}

type ImageCache struct {
	cache map[string]cacheData
}

func NewImage() Cache {
	return ImageCache{cache: map[string]cacheData{}}
}

func (i ImageCache) Age(key string) (time.Duration, error) {
	if _, ok := i.cache[key]; !ok {
		return 0, fmt.Errorf("invalid key provided %v", key)
	}

	return time.Since(i.cache[key].creation), nil
}

func (i ImageCache) Save(key string, value interface{}) error {
	url, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid value provided to cache manager")
	}

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create("/tmp/image" + filepath.Ext(url))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	i.cache[ImageKey] = cacheData{path: "/tmp/image" + filepath.Ext(url), creation: time.Now()}

	return nil
}

func (i ImageCache) Load(key string) (interface{}, error) {
	if _, ok := i.cache[key]; !ok {
		return nil, fmt.Errorf("invalid key provided %v", key)
	}

	return os.Open(i.cache[key].path)
}
