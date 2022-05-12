package storage

import (
	"github.com/mkarulina/url-reduction-service/config"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

func Test_storage_GetAllUrls(t *testing.T) {
	_, err := config.LoadConfig("../../config")
	if err != nil {
		log.Fatal(err)
	}
	os.Setenv("FILE_STORAGE_PATH", "../../tmp/test_urls.log")
	defer os.Remove("../../tmp/test_urls.log")

	stg := New()

	urls := []Link{
		{Key: "testKey11", Link: "http://testhost.ru/11"},
		{Key: "testKey12", Link: "http://testhost.ru/12"},
		{Key: "testKey13", Link: "http://testhost.ru/13"},
	}

	for i := 0; i < len(urls); i++ {
		err = stg.AddLinkToDB(&Link{Key: urls[i].Key, Link: urls[i].Link})
		if err != nil {
			log.Println("can't add link to db", err)
		}
	}

	got, err := stg.GetAllUrls()
	require.NoError(t, err)
	require.Equal(t, len(got), len(urls))

	for i := 0; i < len(urls); i++ {
		require.Contains(t, got, urls[i])
	}
}

func Test_storage_GetKeyByLink(t *testing.T) {
	_, err := config.LoadConfig("../../config")
	if err != nil {
		log.Fatal(err)
	}
	os.Setenv("FILE_STORAGE_PATH", "../../tmp/test_urls.log")
	defer os.Remove("../../tmp/test_urls.log")

	stg := New()

	tests := []struct {
		name string
		link string
		key  string
	}{
		{
			name: "ok",
			link: "http://ya.ru/2",
			key:  "key2",
		},
		{
			name: "not existent link",
			link: "http://notExistent.ru",
			key:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err = stg.AddLinkToDB(&Link{Key: tt.key, Link: tt.link})
			if err != nil {
				log.Println("can't add link to db", err)
			}

			if got := stg.GetKeyByLink(tt.link); got != tt.key {
				t.Errorf("GetKeyByLink() = %v, want %v", got, tt.key)
			}
		})
	}
}

func Test_storage_GetLinkByKey(t *testing.T) {
	_, err := config.LoadConfig("../../config")
	if err != nil {
		log.Fatal(err)
	}
	os.Setenv("FILE_STORAGE_PATH", "../../tmp/test_urls.log")
	defer os.Remove("../../tmp/test_urls.log")

	stg := New()

	tests := []struct {
		name string
		key  string
		link string
	}{
		{
			name: "ok",
			key:  "http://ya.ru/2",
			link: "key2",
		},
		{
			name: "not existent link",
			key:  "http://notExistent.ru",
			link: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err = stg.AddLinkToDB(&Link{Key: tt.key, Link: tt.link})
			if err != nil {
				log.Println("can't add link to db", err)
			}

			if got := stg.GetLinkByKey(tt.key); got != tt.link {
				t.Errorf("GetLinkByKey() = %v, want %v", got, tt.link)
			}
		})
	}
}
