package ping

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestPingHandler_Error(t *testing.T) {
	os.Setenv("DATABASE_DSN", "")

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(PingHandler)
	handler.ServeHTTP(rec, req)
	result := rec.Result()

	require.Equal(t, http.StatusInternalServerError, result.StatusCode)
}
