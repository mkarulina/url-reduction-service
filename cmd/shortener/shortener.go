package main

import (
	"bytes"
	"encoding/base64"
)

func ShortenLink(link string) string {
	var key string

	for i := 0; i < len(UrlsDB); i++ {
		if UrlsDB[i].Url == link {
			key = UrlsDB[i].Key
			return key
		}
	}
	if key == "" {
		buf := bytes.Buffer{}
		encoder := base64.NewEncoder(base64.URLEncoding, &buf)
		encoder.Write([]byte(link))
		UrlsDB = append(UrlsDB, SavedUrl{buf.String(), link})
		key = buf.String()
	}
	return key
}