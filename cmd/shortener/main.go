package main

import (
	"compress/gzip"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mkarulina/url-reduction-service/config"
	"github.com/mkarulina/url-reduction-service/internal/handlers"
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
	c := handlers.NewContainer()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/{linkKey}", c.GetLinkHandler)
		r.Post("/", c.PostLinkHandler)
		r.Post("/api/shorten", c.ShortenHandler)
	})

	config.SetFlags()

	port := config.GetConfig("SERVER_ADDRESS")
	if err := http.ListenAndServe(port, gzipHandler(r)); err != nil {
		log.Fatal(err)
	}
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			handler.ServeHTTP(w, r)
			return
		} else if !strings.Contains(r.Header.Get("Content-Type"), "application/json") && !strings.Contains(r.Header.Get("Content-Type"), "text/plain") {
			handler.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		handler.ServeHTTP(gzipWriter{w, gz}, r)
	})
}
