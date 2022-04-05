package main

import (
	"bytes"
	"encoding/json"
	"github.com/mkarulina/url-reduction-service/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestPostLinkHandler(t *testing.T) {
	c := NewContainer()

	type want struct {
		statusCode int
		key        string
	}
	tests := []struct {
		name string
		url  string
		body string
		want want
	}{
		{
			"POST request with valid link",
			"/",
			"http://testhost.ru/1",
			want{
				201,
				`^[0-9,a-zA-Z]*$`,
			},
		},
		{
			"POST request with not valid link",
			"/",
			"test",
			want{
				400,
				`^Проверьте формат url в теле запроса$`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, test.url, bytes.NewBufferString(test.body))
			rec := httptest.NewRecorder()
			handler := http.HandlerFunc(c.PostLinkHandler)
			handler.ServeHTTP(rec, req)
			result := rec.Result()

			body, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.want.statusCode, result.StatusCode)
			assert.Regexp(t, test.want.key, string(body))
		})
	}
}

func TestGetLinkHandler(t *testing.T) {
	c := NewContainer()
	var wg sync.WaitGroup
	wg.Add(1)
	go c.AddLinkToDB(&storage.Link{"testKey2", "http://testhost.ru/2"}, &wg)
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

func TestShortenHandler(t *testing.T) {
	c := NewContainer()

	type requestBody struct {
		Url string `json:"url"`
	}
	type responseBody struct {
		Result string `json:"result"`
		Error  string `json:"error"`
	}
	type want struct {
		statusCode int
		key        string
		wantError  bool
	}
	tests := []struct {
		name string
		url  string
		body string
		want want
	}{
		{
			"POST request to /api/shorten with valid link",
			"/api/shorten",
			"http://testhost.ru/3",
			want{
				201,
				`^[0-9,a-zA-Z]*$`,
				false,
			},
		},
		{
			"POST request /api/shorten with not valid link",
			"/api/shorten",
			"testhost",
			want{
				400,
				`^Проверьте формат url в теле запроса$`,
				true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			marshalBody, err := json.Marshal(requestBody{test.body})
			if err != nil {
				log.Println("can't marshal url", err)
				return
			}

			req := httptest.NewRequest(http.MethodPost, test.url, bytes.NewReader(marshalBody))
			rec := httptest.NewRecorder()
			handler := http.HandlerFunc(c.ShortenHandler)
			handler.ServeHTTP(rec, req)
			result := rec.Result()

			body, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			unmarshalBody := responseBody{}
			err = json.Unmarshal(body, &unmarshalBody)
			require.NoError(t, err)

			err = result.Body.Close()
			require.NoError(t, err)

			if test.want.wantError {
				assert.Regexp(t, test.want.key, unmarshalBody.Error)
			} else {
				assert.Regexp(t, test.want.key, unmarshalBody.Result)
			}
			assert.Equal(t, test.want.statusCode, result.StatusCode)
		})
	}
}
