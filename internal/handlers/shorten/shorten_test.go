package shorten

import (
	"bytes"
	"encoding/json"
	"github.com/mkarulina/url-reduction-service/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShortenHandler(t *testing.T) {
	_, err := config.LoadConfig("../../../config")
	if err != nil {
		log.Fatal(err)
	}

	type requestBody struct {
		URL string `json:"url"`
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
				`(.*\/)?[0-9,a-zA-Z]*$`,
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
			handler := http.HandlerFunc(ShortenHandler)
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
