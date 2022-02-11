package main

import (
	"strconv"
)

func shortenLink(link string) string {
	var key string

	for i := 0; i < len(urlsDB); i++ {
		if urlsDB[i].Url == link {
			key = urlsDB[i].Id
			return key
		}
	}
	if key == "" {
		newKey := strconv.Itoa(len(urlsDB))
		urlsDB = append(urlsDB, SavedUrl{newKey, link})
		key = newKey
	}
	return key
}
