package batch

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

func BatchLinksHandler(w http.ResponseWriter, r *http.Request) {
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

	var reader io.Reader
	var resp response

	//if err := helpers.SetCookie(w, r); err != nil {
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
		err = gz.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		reader = r.Body
	}

	body, err := io.ReadAll(reader)
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
		} else {
			link, err := helpers.ShortenLink(reqValue)
			if err != nil {
				if code := err.Error(); code != pgerrcode.UniqueViolation {
					log.Panic(err)
				}
			}
			resp = append(resp, sentURL{v.ID, link})
		}
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
