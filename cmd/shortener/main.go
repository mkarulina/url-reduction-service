package main

import (
	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

type Container struct {
	mu       *sync.Mutex
	urls map[string]string
}

func main() {
	c := NewContainer()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/{linkKey}", c.GetLinkHandler)
		r.Post("/", c.PostLinkHandler)
	})
	log.Fatal(http.ListenAndServe(":8080", r))
}

func NewContainer() Container {
	c := Container{
		urls: map[string]string{},
		mu: new(sync.Mutex),
	}
	return c
}

func (c *Container) AddLinkToDB(link string, key string, wg *sync.WaitGroup) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.urls[key] = link
	wg.Done()
}

func (c *Container) GetLinkHandler(w http.ResponseWriter, r *http.Request)  {
	linkKey := strings.Split(r.URL.Path, "/")[1]
	foundLink := c.GetLinkByKey(linkKey)

	if foundLink == "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Url not found"))
		return
	} else {
		w.Header().Set("Location", foundLink)
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte(foundLink))
		return
	}
}

func (c *Container) PostLinkHandler(w http.ResponseWriter, r *http.Request)  {
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
		generatedKey := c.ShortenLink(reqValue)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte (generatedKey))
	}
}
