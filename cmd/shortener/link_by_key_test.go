package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetLinkByKey(t *testing.T) {
	tests := []struct {
		name string
		addedKey string
		searchKey string
		link string
		want string
	}{
		{
			"I can get url by existing key",
			"testKey11",
			"testKey11",
			"ya.ru/11",
			"ya.ru/11",
		},
		{
			"I can't get url by non-existent key",
			"testKey12",
			"nonExistentKey",
			"ya.ru/12",
			"",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			UrlsDB = append(UrlsDB, SavedUrl{test.addedKey, test.link})
			foundUrl := GetLinkByKey(test.searchKey)
			require.Equal(t, test.want, foundUrl, "found value does not match expected")
		})
	}
}
