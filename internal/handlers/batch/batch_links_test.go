package batch

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBatchLinksHandler(t *testing.T) {
	type receivedURL struct {
		ID          string `json:"correlation_id"`
		OriginalURL string `json:"original_url"`
	}

	reqBody := []receivedURL{
		{ID: "111", OriginalURL: "ya111.ru"},
		{ID: "112", OriginalURL: "ya112.ru"},
		{ID: "113", OriginalURL: "ya113.ru"},
	}

	marshalReqBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Println("can't marshal request", err)
		return
	}

	r := bytes.NewReader(marshalReqBody)

	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", r)
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(BatchLinksHandler)

	handler.ServeHTTP(rec, req)
	result := rec.Result()

	require.Equal(t, http.StatusCreated, result.StatusCode)

	body, err := io.ReadAll(result.Body)
	require.NoError(t, err)
	require.NotEmpty(t, body)

	err = result.Body.Close()
	require.NoError(t, err)
}
