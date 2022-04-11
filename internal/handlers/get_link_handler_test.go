package handlers

import (
	"github.com/mkarulina/url-reduction-service/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestGetLinkHandler(t *testing.T) {
	c := NewContainer()
	var wg sync.WaitGroup
	wg.Add(1)
	go c.AddLinkToDB(&storage.Link{Key: "testKey2", Link: "http://testhost.ru/2"}, &wg)
	wg.Wait()

	type want struct {
		statusCode int
		body       string
		location   string
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			"GET request with existent key",
			"/testKey2",
			want{
				307,
				"http://testhost.ru/2",
				"http://testhost.ru/2",
			},
		},
		{
			"GET request with non-existent key",
			"/nonexistent",
			want{
				200,
				"Url not found",
				"",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.url, nil)
			rec := httptest.NewRecorder()
			handler := http.HandlerFunc(c.GetLinkHandler)
			handler.ServeHTTP(rec, req)
			result := rec.Result()

			body, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.want.statusCode, result.StatusCode)
			assert.Equal(t, test.want.body, string(body))
			assert.Equal(t, test.want.location, result.Header.Get("Location"))
		})
	}
}

func TestGetLinkByKey(t *testing.T) {
	c := NewContainer()
	var wg sync.WaitGroup
	wg.Add(1)
	go c.AddLinkToDB(&storage.Link{Key: "testKey11", Link: "http://testhost.ru/11"}, &wg)
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
	go c.AddLinkToDB(&storage.Link{Key: "testKey12", Link: "http://testhost.ru/12"}, &wg)
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
