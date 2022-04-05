package main

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mkarulina/url-reduction-service/config"
	"github.com/mkarulina/url-reduction-service/internal/storage"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

type Container struct {
	mu   *sync.Mutex
	urls map[string]string
	file string
}

func main() {
	c := NewContainer()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/{linkKey}", c.GetLinkHandler)
		r.Post("/", c.PostLinkHandler)
		r.Post("/api/shorten", c.ShortenHandler)
	})

	config.SetFlags()

	fmt.Println("address:", config.GetConfig("SERVER_ADDRESS"))
	fmt.Println("base url:", config.GetConfig("BASE_URL"))
	fmt.Println("file:", config.GetConfig("FILE_STORAGE_PATH"))

	port := config.GetConfig("SERVER_ADDRESS")
	log.Fatal(http.ListenAndServe(port, r))
}

func NewContainer() Container {
	filePath := config.GetConfig("FILE_STORAGE_PATH")
	c := Container{
		mu:   new(sync.Mutex),
		urls: map[string]string{},
		file: filePath,
	}
	return c
}

func (c *Container) GetLinkHandler(w http.ResponseWriter, r *http.Request) {
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

func (c *Container) PostLinkHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("can't read body", err)
		return
	}
	reqValue := string(body)
	validURL := govalidator.IsURL(reqValue)
	if !validURL {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Проверьте формат url в теле запроса"))
		return
	} else {
		generatedLink := c.ShortenLink(reqValue)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(generatedLink))
	}
}

func (c *Container) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	type receivedUrl struct {
		Url string `json:"url"`
	}
	type sentUrl struct {
		Result string `json:"result"`
	}
	type error struct {
		Error string `json:"error"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("can't read body", err)
		return
	}
	unmarshalBody := receivedUrl{}
	if err := json.Unmarshal(body, &unmarshalBody); err != nil {
		log.Println("can't unmarshal request body", err)
		return
	}
	reqValue := unmarshalBody.Url

	validURL := govalidator.IsURL(reqValue)
	if !validURL {
		errText, err := json.Marshal(error{"Проверьте формат url в теле запроса"})
		if err != nil {
			log.Println("can't marshal response", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errText)
		return
	} else {
		generatedLink := sentUrl{c.ShortenLink(reqValue)}
		marshalResult, err := json.Marshal(generatedLink)
		if err != nil {
			log.Println("can't marshal response", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(marshalResult)
	}
}

func (c *Container) AddLinkToDB(link *storage.Link, wg *sync.WaitGroup) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.file == "" {
		c.urls[link.Key] = link.Link
	} else {
		recorder, err := storage.NewRecorder(c.file)
		if err != nil {
			log.Fatal(err)
		}
		defer recorder.Close()

		if err := recorder.WriteLink(link); err != nil {
			log.Fatal(err)
		}
	}
	wg.Done()
}
