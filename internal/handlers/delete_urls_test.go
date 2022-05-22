package handlers

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/mkarulina/url-reduction-service/internal/mocks"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handler_DeleteUserUrls(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stg := mocks.NewMockStorage(ctrl)
	h := NewHandler(stg)

	tests := []struct {
		name       string
		body       string
		statusCode int
	}{
		{
			name:       "ok",
			body:       "[\"key1\", \"key2\"]",
			statusCode: http.StatusAccepted,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stg.EXPECT().DeleteUrls(gomock.Any()).Return(nil)

			req := httptest.NewRequest(http.MethodGet, "/api/user/urls", bytes.NewBufferString(tt.body))
			rec := httptest.NewRecorder()
			req.AddCookie(&http.Cookie{
				Name:  "session_token",
				Value: "testCookie",
			})

			handler := http.HandlerFunc(h.DeleteUserUrls)
			handler.ServeHTTP(rec, req)
			result := rec.Result()

			require.Equal(t, tt.statusCode, result.StatusCode)

			err := result.Body.Close()
			require.NoError(t, err)
		})
	}
}
