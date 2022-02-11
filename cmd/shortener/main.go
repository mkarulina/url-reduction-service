package main

import (
	"github.com/asaskevich/govalidator"
	"io"
	"net/http"
	"strings"
)

type SavedUrl struct {
	Id string
	Url string
}

var urlsDB [] SavedUrl


func main() {
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			path := r.URL.Path
			reqParam := strings.Split(path, "/")[1]
			foundLink := getLinkById(reqParam)

			if foundLink != "" {
				w.WriteHeader(http.StatusTemporaryRedirect)
				w.Header().Add("location", foundLink)
				w.Write([]byte(foundLink))
			} else {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Url not found"))
			}

		case "POST":
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
			}
			generatedKey := shortenLink(reqValue)
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte ("Id для вашего url: " + generatedKey))
		}
	})
	http.ListenAndServe(":8080", nil)
}