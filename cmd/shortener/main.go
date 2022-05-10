package main

import (
	"compress/gzip"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mkarulina/url-reduction-service/config"
	"github.com/mkarulina/url-reduction-service/internal/handlers/all_urls"
	"github.com/mkarulina/url-reduction-service/internal/handlers/batch"
	"github.com/mkarulina/url-reduction-service/internal/handlers/get_link"
	"github.com/mkarulina/url-reduction-service/internal/handlers/ping"
	"github.com/mkarulina/url-reduction-service/internal/handlers/post_link"
	"github.com/mkarulina/url-reduction-service/internal/handlers/shorten"
	"github.com/mkarulina/url-reduction-service/internal/storage"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func main() {
	_, err := config.LoadConfig("config")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	storage.New()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/ping", ping.PingHandler)
		r.Get("/{linkKey}", getlink.GetLinkHandler)
		r.Get("/api/user/urls", allurls.GetAllUrlsHandler)
		r.Post("/", postlink.PostLinkHandler)
		r.Post("/api/shorten", shorten.ShortenHandler)
		r.Post("/api/shorten/batch", batch.BatchLinksHandler)
	})

	address := viper.GetString("SERVER_ADDRESS")
	handler := gzipHandle(r)

	if err := http.ListenAndServe(address, handler); err != nil {
		log.Fatal(err)
	}
}

func (gw gzipWriter) Write(b []byte) (int, error) {
	return gw.Writer.Write(b)
}

func gzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})

}
