package handlers

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/mkarulina/url-reduction-service/internal/mocks"
	"github.com/mkarulina/url-reduction-service/internal/storage"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllUrlsHandler(t *testing.T) {
	urls := []storage.ResponseLink{
		{Key: "testKey11", Link: "http://testhost.ru/11"},
		{Key: "testKey12", Link: "http://testhost.ru/12"},
		{Key: "testKey13", Link: "http://testhost.ru/13"},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stg := mocks.NewMockStorage(ctrl)
	stg.EXPECT().GetAllUrlsByUserID(gomock.Any()).Return(urls, nil)

	h := NewHandler(stg)

	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "testCookie",
	})

	handler := http.HandlerFunc(h.GetAllUrlsHandler)

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
