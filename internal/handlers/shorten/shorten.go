package shorten

import (
	"compress/gzip"
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/jackc/pgerrcode"
	"github.com/mkarulina/url-reduction-service/internal/helpers"
	"io"
	"log"
	"net/http"
)

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	type receivedURL struct {
		URL string `json:"url"`
	}
	type sentURL struct {
		Result string `json:"result"`
	}
	type error struct {
		Error string `json:"error"`
	}

	var reader io.Reader

	//e := encryptor.New()
	//if err := e.SetCookie(w, r); err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
	}

	reader = r.Body

	body, err := io.ReadAll(reader)
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

	link, err := helpers.ShortenLink(reqValue)
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
