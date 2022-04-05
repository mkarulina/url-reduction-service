package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestShortenLink(t *testing.T) {
	c := NewContainer()

	tests := []struct {
		name string
		link string
	}{
		{
			"I can generate key for url without protocol",
			"testhost.ru/21",
		},
		{
			"I can generate key for url with http",
			"http://testhost.ru/22",
		},
		{
			"I can generate key for url with https",
			"httpы://testhost.ru/23",
		},
		{
			"I can generate key for url with params",
			"https://www.avito.ru/murino/kvartiry/prodam/novostroyka-ASgBAQICAUSSA8YQAUDmBxSOUg?cd=1&f=ASgBAQECAUSSA8YQA0DkBxT~UeYHFI5SyggU_lgCRYQJE3siZnJvbSI6MjAsInRvIjo1MH3GmgweeyJmcm9tIjo1MDAwMDAwLCJ0byI6MTUwMDAwMDB9",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			key := c.ShortenLink(test.link)
			require.NotEmptyf(t, key, "key not generated", test.link)
		})
	}
}
