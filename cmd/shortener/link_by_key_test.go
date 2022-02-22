package main

import (
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestGetLinkByKey(t *testing.T) {
	c := NewContainer()
	var wg sync.WaitGroup
	wg.Add(1)
	go c.AddLinkToDB("http://ya.ru/11", "testKey11", &wg)
	wg.Wait()

	tests := []struct {
		name string
		key string
		want string
	}{
		{
			"I can get url by existing key",
			"testKey11",
			"http://ya.ru/11",
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
