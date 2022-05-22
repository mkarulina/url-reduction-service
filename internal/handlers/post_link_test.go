package handlers

import (
	"bytes"
	"github.com/mkarulina/url-reduction-service/internal/storage"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostLinkHandler(t *testing.T) {
	baseURL := viper.GetString("BASE_URL")
	h := NewHandler(storage.New())

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
				baseURL + `\/[0-9,a-zA-Z]*$`,
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
			req.AddCookie(&http.Cookie{
				Name:  "session_token",
				Value: "testCookie",
			})

			handler := http.HandlerFunc(h.PostLinkHandler)
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
