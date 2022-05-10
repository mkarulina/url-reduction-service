package helpers

import (
	"github.com/mkarulina/url-reduction-service/config"
	"github.com/mkarulina/url-reduction-service/internal/storage"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

func TestShortenLink(t *testing.T) {
	_, err := config.LoadConfig("../../config")
	if err != nil {
		log.Fatal(err)
	}
	os.Setenv("FILE_STORAGE_PATH", "../../tmp/test_urls.log")
	defer os.Remove("../../tmp/test_urls.log")

	storage.New()

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
			"https://testhost.ru/23",
		},
		{
			"I can generate key for url with params",
			"https://www.avito.ru/murino/kvartiry/prodam/novostroyka-ASgBAQICAUSSA8YQAUDmBxSOUg?cd=1&f=ASgBAQECAUSSA8YQA0DkBxT~UeYHFI5SyggU_lgCRYQJE3siZnJvbSI6MjAsInRvIjo1MH3GmgweeyJmcm9tIjo1MDAwMDAwLCJ0byI6MTUwMDAwMDB9",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			key, _ := ShortenLink(test.link)
			require.NotEmptyf(t, key, "key not generated", test.link)
		})
	}
}
