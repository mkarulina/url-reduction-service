package handlers

import (
	"github.com/golang/mock/gomock"
	"github.com/mkarulina/url-reduction-service/internal/mocks"
	"github.com/mkarulina/url-reduction-service/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLinkHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stg := mocks.NewMockStorage(ctrl)
	h := NewHandler(stg)

	type want struct {
		statusCode int
		body       string
		location   string
	}
	tests := []struct {
		name string
		path string
		link *storage.Link
		want want
	}{
		{
			name: "GET request with existent key",
			path: "/testKey2",
			link: &storage.Link{
				UserID:    "any",
				Key:       "any",
				Link:      "http://testhost.ru/2",
				IsDeleted: false,
			},
			want: want{
				307,
				"http://testhost.ru/2",
				"http://testhost.ru/2",
			},
		},
		{
			name: "GET request with non-existent key",
			path: "/nonexistent",
			link: nil,
			want: want{
				200,
				"Url not found",
				"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stg.EXPECT().GetLinkByKey(gomock.Any()).Return(tt.link)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()
			req.AddCookie(&http.Cookie{
				Name:  "session_token",
				Value: "testCookie",
			})

			handler := http.HandlerFunc(h.GetLinkHandler)
			handler.ServeHTTP(rec, req)
			result := rec.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.location, result.Header.Get("Location"))

			err := result.Body.Close()
			require.NoError(t, err)
		})
	}
}
