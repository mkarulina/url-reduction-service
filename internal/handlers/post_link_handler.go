package handlers

import (
	"bytes"
	"encoding/base64"
	"github.com/asaskevich/govalidator"
	"github.com/mkarulina/url-reduction-service/config"
	"github.com/mkarulina/url-reduction-service/internal/storage"
	"io"
	"log"
	"net/http"
	"sync"
)

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

func (c *Container) ShortenLink(link string) string {
	var wg sync.WaitGroup
	var key string

	if existentKey := c.GetKeyByLink(link); existentKey != "" {
		key = existentKey
	}

	if key == "" {
		buf := bytes.Buffer{}

		encoder := base64.NewEncoder(base64.URLEncoding, &buf)
		encoder.Write([]byte(link))
		key = buf.String()

		wg.Add(1)
		go c.AddLinkToDB(&storage.Link{key, link}, &wg)
		wg.Wait()
	}
	shortLink := config.GetConfig("BASE_URL") + key
	return shortLink
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
