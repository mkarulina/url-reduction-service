package main

import (
	"compress/gzip"
	"encoding/hex"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mkarulina/url-reduction-service/config"
	"github.com/mkarulina/url-reduction-service/internal/handlers"
	"github.com/mkarulina/url-reduction-service/internal/helpers/encryptor"
	"github.com/mkarulina/url-reduction-service/internal/storage"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
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

	h := handlers.NewHandler(storage.New())

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(gzipHandle)
	r.Use(cookieHandle)

	r.Route("/", func(r chi.Router) {
		r.Get("/ping", h.PingHandler)
		r.Get("/{linkKey}", h.GetLinkHandler)
		r.Get("/api/user/urls", h.GetAllUrlsHandler)
		r.Post("/", h.PostLinkHandler)
		r.Post("/api/shorten", h.ShortenHandler)
		r.Post("/api/shorten/batch", h.BatchLinksHandler)
	})

	address := viper.GetString("SERVER_ADDRESS")

	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatal(err)
	}
}

func (gw gzipWriter) Write(b []byte) (int, error) {
	return gw.Writer.Write(b)
}

func gzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.Header.Get(`Content-Encoding`), `gzip`) {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = gz
			next.ServeHTTP(w, r)
			return
		}

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

func cookieHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e := encryptor.New()

		cookie, err := r.Cookie("session_token")
		if err != nil {
			log.Println(err)
		}

		if cookie != nil && len(cookie.Value) >= 16 {
			decrCookie, err := e.DecryptData(cookie.Value)
			if err != nil {
				log.Println(err)
			}
			if decrCookie != nil {
				next.ServeHTTP(w, r)
				return
			}
		}

		random, err := e.GenerateRandom(16)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newUser := hex.EncodeToString(random)

		t, err := e.EncryptData([]byte(newUser))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token := hex.EncodeToString(t)

		newCookie := &http.Cookie{
			Name:    "session_token",
			Value:   token,
			Expires: time.Now().Add(3 * time.Hour),
			Secure:  false,
		}

		http.SetCookie(w, newCookie)
		r.AddCookie(newCookie)
		next.ServeHTTP(w, r)
	})
}
