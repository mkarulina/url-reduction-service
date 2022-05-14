package handlers

import (
	"github.com/golang/mock/gomock"
	"github.com/mkarulina/url-reduction-service/config"
	"github.com/mkarulina/url-reduction-service/internal/mocks"
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
	_, err := config.LoadConfig("../../config")
	if err != nil {
		log.Fatal(err)
	}
	os.Setenv("FILE_STORAGE_PATH", "../../tmp/test_urls.log")
	defer os.Remove("../../tmp/test_urls.log")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type want struct {
		statusCode int
		body       string
		location   string
	}
	tests := []struct {
		name string
		path string
		link string
		want want
	}{
		{
			name: "GET request with existent key",
			path: "/testKey2",
			link: "http://testhost.ru/2",
			want: want{
				307,
				"http://testhost.ru/2",
				"http://testhost.ru/2",
			},
		},
		{
			name: "GET request with non-existent key",
			path: "/nonexistent",
			link: "",
			want: want{
				200,
				"Url not found",
				"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stg := mocks.NewMockStorage(ctrl)
			stg.EXPECT().GetLinkByKey(gomock.Any()).Return(tt.link)

			h := NewHandler(stg)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()
			req.AddCookie(&http.Cookie{
				Name:  "session_token",
				Value: "testCookie",
			})

			handler := http.HandlerFunc(h.GetLinkHandler)
			handler.ServeHTTP(rec, req)
			result := rec.Result()

			body, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.body, string(body))
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))

			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}
