package postlink

import (
	"compress/gzip"
	"github.com/asaskevich/govalidator"
	"github.com/jackc/pgerrcode"
	"github.com/mkarulina/url-reduction-service/internal/helpers"
	"io"
	"log"
	"net/http"
)

func PostLinkHandler(w http.ResponseWriter, r *http.Request) {
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

	body, err := io.ReadAll(reader)
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
		generatedLink, err := helpers.ShortenLink(reqValue)
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
