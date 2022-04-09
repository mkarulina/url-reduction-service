package handlers

import (
	"bytes"
	"github.com/mkarulina/url-reduction-service/cmd/shortener"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostLinkHandler(t *testing.T) {
	c := main.NewContainer()

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

func TestShortenLink(t *testing.T) {
	c := main.NewContainer()

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
