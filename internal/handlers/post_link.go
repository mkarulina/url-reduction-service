package handlers

import (
	"github.com/asaskevich/govalidator"
	"github.com/jackc/pgerrcode"
	"io"
	"log"
	"net/http"
)

func (h *handler) PostLinkHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("can't read body", err)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie == nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	reqValue := string(body)
	validURL := govalidator.IsURL(reqValue)

	if !validURL {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Проверьте формат url в теле запроса"))
		return
	} else {
		generatedLink, err := h.stg.ShortenLink(cookie.Value, reqValue)
		if err != nil {
			if code := err.Error(); code == pgerrcode.UniqueViolation {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(generatedLink))
				return
			}
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(generatedLink))
	}
}
