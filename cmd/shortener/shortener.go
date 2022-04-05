package main

import (
	"bytes"
	"encoding/base64"
	"github.com/mkarulina/url-reduction-service/config"
	"github.com/mkarulina/url-reduction-service/internal/storage"
	"sync"
)

func (c *Container) ShortenLink(link string) string {
	var wg sync.WaitGroup
	var key string

	if existentKey := c.GetKeyByLink(link); existentKey != "" {
		key = existentKey
	}

	if key == "" {
		buf := bytes.Buffer{}

		encoder := base64.NewEncoder(base64.URLEncoding, &buf)
		encoder.Write([]byte(link))
		key = buf.String()

		wg.Add(1)
		go c.AddLinkToDB(&storage.Link{key, link}, &wg)
		wg.Wait()
	}
	shortLink := config.GetConfig("BASE_URL") + key
	return shortLink
}
