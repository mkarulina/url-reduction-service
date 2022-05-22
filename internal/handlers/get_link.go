package handlers

import (
	"net/http"
	"strings"
)

func (h *handler) GetLinkHandler(w http.ResponseWriter, r *http.Request) {
	linkKey := strings.Split(r.URL.Path, "/")[1]
	foundLink := h.stg.GetLinkByKey(linkKey)

	if foundLink == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Url not found"))
		return
	}

	if foundLink.IsDeleted {
		w.WriteHeader(http.StatusGone)
		w.Write([]byte("Url has been deleted"))
		return
	}

	w.Header().Set("Location", foundLink.Link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
