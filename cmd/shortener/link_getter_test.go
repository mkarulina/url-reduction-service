package main

import (
	"github.com/mkarulina/url-reduction-service/internal/storage"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestGetLinkByKey(t *testing.T) {
	c := NewContainer()
	var wg sync.WaitGroup
	wg.Add(1)
	go c.AddLinkToDB(&storage.Link{"testKey11", "http://testhost.ru/11"}, &wg)
	wg.Wait()

	tests := []struct {
		name string
		key  string
		want string
	}{
		{
			"I can get url by existing key",
			"testKey11",
			"http://testhost.ru/11",
		},
		{
			"I can't get url by non-existent key",
			"nonExistentKey",
			"",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			foundURL := c.GetLinkByKey(test.key)
			require.Equal(t, test.want, foundURL, "found value does not match expected")
		})
	}
}

func TestGetKeyByLink(t *testing.T) {
	c := NewContainer()
	var wg sync.WaitGroup
	wg.Add(1)
	go c.AddLinkToDB(&storage.Link{"testKey12", "http://testhost.ru/12"}, &wg)
	wg.Wait()

	tests := []struct {
		name string
		link string
		want string
	}{
		{
			"I can get url by existing key",
			"http://testhost.ru/12",
			"testKey12",
		},
		{
			"I can't get url by non-existent key",
			"http://nonexistentlink.ru",
			"",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			foundKey := c.GetKeyByLink(test.link)
			require.Equal(t, test.want, foundKey, "found value does not match expected")
		})
	}
}
