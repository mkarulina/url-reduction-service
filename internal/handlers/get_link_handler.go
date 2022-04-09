package handlers

import (
	"github.com/mkarulina/url-reduction-service/internal/storage"
	"log"
	"net/http"
	"strings"
)

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

func (c *Container) GetLinkByKey(linkKey string) string {
	var foundLink string

	if c.file == "" {
		if val, found := c.urls[linkKey]; found {
			foundLink = val
		}
	} else {
		reader, err := storage.NewReader(c.file)
		if err != nil {
			log.Fatal(err)
		}
		defer reader.Close()
		for {
			readLine, err := reader.ReadLink()
			if readLine == nil {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			if readLine.Key == linkKey {
				foundLink = readLine.Link
				break
			}
		}
	}

	return foundLink
}

func (c *Container) GetKeyByLink(link string) string {
	var foundKey string

	if c.file == "" {
		for key, value := range c.urls {
			if value == link {
				foundKey = key
				return foundKey
			}
		}
	} else {
		reader, err := storage.NewReader(c.file)
		if err != nil {
			log.Fatal(err)
		}
		defer reader.Close()

		for {
			readLine, err := reader.ReadLink()
			if readLine == nil {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			if readLine.Link == link {
				foundKey = readLine.Key
				break
			}
		}
	}

	return foundKey
}
