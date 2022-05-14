package handlers

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/jackc/pgerrcode"
	"io"
	"log"
	"net/http"
)

func (h *handler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	type receivedURL struct {
		URL string `json:"url"`
	}
	type sentURL struct {
		Result string `json:"result"`
	}
	type error struct {
		Error string `json:"error"`
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Println(err)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("can't read body", err)
		return
	}
	unmarshalBody := receivedURL{}
	if err := json.Unmarshal(body, &unmarshalBody); err != nil {
		log.Println("can't unmarshal request body", err)
		return
	}
	reqValue := unmarshalBody.URL

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
	}

	link, err := h.stg.ShortenLink(cookie.Value, reqValue)
	if err != nil {
		if code := err.Error(); code == pgerrcode.UniqueViolation {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
	}

	generatedLink := sentURL{link}
	marshalResult, err := json.Marshal(generatedLink)
	if err != nil {
		log.Println("can't marshal response", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(marshalResult)
}
