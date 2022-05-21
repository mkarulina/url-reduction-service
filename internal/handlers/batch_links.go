package handlers

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/jackc/pgerrcode"
	"io"
	"log"
	"net/http"
)

func (h *handler) BatchLinksHandler(w http.ResponseWriter, r *http.Request) {
	type receivedURL struct {
		ID          string `json:"correlation_id"`
		OriginalURL string `json:"original_url"`
	}
	type sentURL struct {
		ID       string `json:"correlation_id"`
		ShortURL string `json:"short_url"`
	}

	type request []receivedURL
	type response []sentURL

	var resp response

	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Println(err)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("can't read body", err)
		return
	}
	unmarshalBody := request{}
	if err := json.Unmarshal(body, &unmarshalBody); err != nil {
		log.Println("can't unmarshal request body", err)
		return
	}

	for _, v := range unmarshalBody {
		reqValue := v.OriginalURL
		validURL := govalidator.IsURL(reqValue)

		if !validURL {
			log.Println("invalid link: ", v)
			continue
		}
		link, err := h.stg.ShortenLink(cookie.Value, reqValue)
		if err != nil {
			if code := err.Error(); code != pgerrcode.UniqueViolation {
				log.Panic(err)
			}
		}
		resp = append(resp, sentURL{v.ID, link})
	}

	marshalResult, err := json.Marshal(resp)
	if err != nil {
		log.Println("can't marshal response", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(marshalResult)
}
