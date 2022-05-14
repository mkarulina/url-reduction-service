package handlers

import (
	"net/http"
	"strings"
)

func (h *handler) GetLinkHandler(w http.ResponseWriter, r *http.Request) {
	linkKey := strings.Split(r.URL.Path, "/")[1]
	foundLink := h.stg.GetLinkByKey(linkKey)

	if foundLink == "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Url not found"))
		return
	}

	w.Header().Set("Location", foundLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte(foundLink))
}
