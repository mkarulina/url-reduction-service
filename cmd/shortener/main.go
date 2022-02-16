package main

import (
	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"io"
	"log"
	"net/http"
	"strings"
)

type SavedUrl struct {
	Key, Url string
}

var UrlsDB []SavedUrl

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/{linkKey}", GetLinkHandler)
		r.Post("/", PostLinkHandler)
	})
	log.Fatal(http.ListenAndServe(":8080", r))
}

func GetLinkHandler(w http.ResponseWriter, r *http.Request)  {
	linkKey := strings.Split(r.URL.Path, "/")[1]
	foundLink := GetLinkByKey(linkKey)

	if foundLink == "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Url not found"))
		return
	} else {
		w.Header().Set("Location", foundLink)
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte(foundLink))
	}
}

func PostLinkHandler(w http.ResponseWriter, r *http.Request)  {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	reqValue := string(body)
	validURL := govalidator.IsURL(reqValue)
	if !validURL {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Проверьте формат url в теле запроса"))
		return
	} else {
		generatedKey := ShortenLink(reqValue)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte (generatedKey))
	}
}
