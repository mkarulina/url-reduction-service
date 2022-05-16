package main

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_gzipHandle(t *testing.T) {
	type request struct {
		method string
		header http.Header
		target string
	}
	tests := []struct {
		name            string
		req             request
		next            http.Handler
		wantHeaderValue string
	}{
		{
			name: "request with gzip header",
			req: request{
				method: http.MethodPost,
				header: http.Header{
					"Accept-Encoding": []string{"gzip"},
				},
				target: "http://localhost:8080",
			},
			next:            http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			wantHeaderValue: "gzip",
		},
		{
			name: "request without gzip header",
			req: request{
				method: http.MethodPost,
				header: http.Header{},
				target: "http://localhost:8080",
			},
			next:            http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			wantHeaderValue: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.req.method, tt.req.target, nil)
			req.Header = tt.req.header
			res := httptest.NewRecorder()

			gzh := gzipHandle(tt.next)
			gzh.ServeHTTP(res, req)

			require.Equal(t, tt.wantHeaderValue, res.Header().Get("Content-Encoding"))
		})
	}
}
