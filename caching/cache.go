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
	Age(key string) (time.Time, error)
}

type ImageCache struct {
	cache map[string]string
}

func NewImage() Cache {
	return ImageCache{cache: map[string]string{}}
}

func (i ImageCache) Age(key string) (time.Time, error) {
	panic("not implemented")
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

	i.cache[ImageKey] = "/tmp/image" + filepath.Ext(url)

	return nil
}

func (i ImageCache) Load(key string) (interface{}, error) {
	return os.Open(i.cache[key])
}
