package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)


func TestPostLinkHandler(t *testing.T) {
	type want struct {
		statusCode int
		key string
	}
	tests := []struct {
		name string
		url string
		body string
		want want
	}{
		{
			"POST request with valid link",
			"/",
			"testlink.com/1",
			want {
				201,
				`^[0-9,a-zA-Z]*$`,
			},
		},
		{
			"POST request with not valid link",
			"/",
			"test",
			want {
				400,
				`^Проверьте формат url в теле запроса$`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, test.url, bytes.NewBufferString(test.body))
			rec := httptest.NewRecorder()
			handler := http.HandlerFunc(PostLinkHandler)
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
	UrlsDB = append(UrlsDB, SavedUrl{"testKey", "google.com"})
	t.Log(UrlsDB)

	type want struct {
		statusCode int
		body string
	}
	tests := []struct {
		name string
		url string
		want want
	}{
		{
			"GET request with existent key",
			"/testKey",
			want {
				307,
				"google.com",
			},
		},
		{
			"GET request with non-existent key",
			"/nonexistent",
			want {
				200,
				"Url not found",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.url, nil)
			rec := httptest.NewRecorder()
			handler := http.HandlerFunc(GetLinkHandler)
			handler.ServeHTTP(rec, req)
			result := rec.Result()

			body, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.want.statusCode, result.StatusCode)
			assert.Equal(t, test.want.body, string(body))

		})
	}
}