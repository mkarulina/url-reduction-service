package getlink

import (
	"github.com/mkarulina/url-reduction-service/config"
	"github.com/mkarulina/url-reduction-service/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetLinkHandler(t *testing.T) {
	_, err := config.LoadConfig("../../../config")
	if err != nil {
		log.Fatal(err)
	}
	os.Setenv("FILE_STORAGE_PATH", "../../../tmp/test_urls.log")
	defer os.Remove("../../../tmp/test_urls.log")

	stg := storage.New()

	err = stg.AddLinkToDB(&storage.Link{Key: "testKey2", Link: "http://testhost.ru/2"})
	if err != nil {
		log.Println("can't add link to db", err)
	}

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

			handler := http.HandlerFunc(GetLinkHandler)
			handler.ServeHTTP(rec, req)
			result := rec.Result()

			body, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			assert.Equal(t, test.want.statusCode, result.StatusCode)
			assert.Equal(t, test.want.body, string(body))
			assert.Equal(t, test.want.location, result.Header.Get("Location"))

			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}
