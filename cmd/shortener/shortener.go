package main

import (
	"bytes"
	"encoding/base64"
	"sync"
)

func (c *Container) ShortenLink(link string) string {
	var key string

	for k, v := range c.urls {
		if v == link {
			key = k
			return key
		}
	}

	if key == "" {
		var wg sync.WaitGroup
		buf := bytes.Buffer{}

		encoder := base64.NewEncoder(base64.URLEncoding, &buf)
		encoder.Write([]byte(link))
		key = buf.String()

		wg.Add(1)
		go c.AddLinkToDB(link, buf.String(), &wg)
		wg.Wait()
	}
	return key
}
