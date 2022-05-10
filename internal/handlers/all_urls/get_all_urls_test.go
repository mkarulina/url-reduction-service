package allurls

import (
	"encoding/json"
	"github.com/mkarulina/url-reduction-service/config"
	"github.com/mkarulina/url-reduction-service/internal/storage"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetAllUrlsHandler(t *testing.T) {
	_, err := config.LoadConfig("../../../config")
	if err != nil {
		log.Fatal(err)
	}
	os.Setenv("FILE_STORAGE_PATH", "../../../tmp/test_urls.log")
	defer os.Remove("../../../tmp/test_urls.log")

	stg := storage.New()

	urls := []storage.Link{
		{Key: "testKey11", Link: "http://testhost.ru/11"},
		{Key: "testKey12", Link: "http://testhost.ru/12"},
		{Key: "testKey13", Link: "http://testhost.ru/13"},
	}

	for i := 0; i < len(urls); i++ {
		err = stg.AddLinkToDB(&storage.Link{Key: urls[i].Key, Link: urls[i].Link})
		if err != nil {
			log.Println("can't add link to db", err)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(GetAllUrlsHandler)

	handler.ServeHTTP(rec, req)
	result := rec.Result()

	wantResp, err := json.Marshal(urls)
	if err != nil {
		log.Println("can't marshal response", err)
		return
	}

	body, err := io.ReadAll(result.Body)
	require.NoError(t, err)
	require.Equal(t, wantResp, body)
	err = result.Body.Close()
	require.NoError(t, err)
}
