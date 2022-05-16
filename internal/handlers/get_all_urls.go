package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *handler) GetAllUrlsHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Println(err)
	}

	response, err := h.stg.GetAllUrlsByUserID(cookie.Value)
	if err != nil {
		log.Println("can't get urls", err)
		return
	}

	if len(response) == 0 {
		w.WriteHeader(http.StatusNoContent)
	}

	marshalResp, err := json.Marshal(response)
	if err != nil {
		log.Println("can't marshal response", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshalResp)
}
