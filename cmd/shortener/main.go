package main

import (
	"compress/gzip"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mkarulina/url-reduction-service/config"
	"github.com/mkarulina/url-reduction-service/internal/handlers"
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

type gzipReader struct {
	http.ResponseWriter
	Reader io.Reader
}

func main() {
	_, err := config.LoadConfig("config")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	c := handlers.NewContainer()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/{linkKey}", c.GetLinkHandler)
		r.Post("/", c.PostLinkHandler)
		r.Post("/api/shorten", c.ShortenHandler)
	})

	address := viper.GetString("SERVER_ADDRESS")
	if err := http.ListenAndServe(address, gzipHandler(r)); err != nil {
		log.Fatal(err)
	}
}

func (gw gzipWriter) Write(b []byte) (int, error) {
	return gw.Writer.Write(b)
}

func (gr gzipReader) Read(b []byte) (int, error) {
	return gr.Reader.Read(b)
}

func gzipHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get(`Content-Encoding`) == `gzip` {

			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				r.Header.Del("Accept-Encoding")
			}
			gzr, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer gzr.Close()

			_, err = io.ReadAll(gzr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			handler.ServeHTTP(gzipReader{ResponseWriter: w, Reader: gzr}, r)

		} else if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {

			gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			defer gz.Close()

			w.Header().Set("Content-Encoding", "gzip")
			handler.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)

		} else {
			handler.ServeHTTP(w, r)
			return
		}

	})

}
