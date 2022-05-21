package handlers

import (
	"encoding/json"
	"github.com/mkarulina/url-reduction-service/internal/storage"
	"io"
	"log"
	"net/http"
)

func (h *handler) DeleteUserUrls(w http.ResponseWriter, r *http.Request) {
	var keys []string

	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Println(err)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("can't read body", err)
		return
	}

	if err := json.Unmarshal(body, &keys); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request body"))
	}

	chIn := genChan(cookie.Value, keys)

	if err = h.stg.DeleteUrls(chIn); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func genChan(cookieValue string, keys []string) chan storage.UserKeys {
	inputCh := make(chan storage.UserKeys)

	go func() {
		inputCh <- storage.UserKeys{
			Cookie: cookieValue,
			Keys:   keys,
		}

		close(inputCh)
	}()
	return inputCh
}
