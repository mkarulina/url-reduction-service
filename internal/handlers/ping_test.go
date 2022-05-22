package handlers

import (
	"github.com/mkarulina/url-reduction-service/internal/storage"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingHandler_Error(t *testing.T) {
	h := NewHandler(storage.New())

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(h.PingHandler)
	handler.ServeHTTP(rec, req)
	result := rec.Result()

	require.Equal(t, http.StatusInternalServerError, result.StatusCode)

	err := result.Body.Close()
	require.NoError(t, err)
}
