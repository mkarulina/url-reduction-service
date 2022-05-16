package handlers

import (
	"github.com/mkarulina/url-reduction-service/internal/storage"
	"net/http"
)

type Handler interface {
	BatchLinksHandler(w http.ResponseWriter, r *http.Request)
	GetAllUrlsHandler(w http.ResponseWriter, r *http.Request)
	GetLinkHandler(w http.ResponseWriter, r *http.Request)
	PingHandler(w http.ResponseWriter, r *http.Request)
	PostLinkHandler(w http.ResponseWriter, r *http.Request)
	ShortenHandler(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	stg storage.Storage
}

func NewHandler(s storage.Storage) Handler {
	h := &handler{
		stg: s,
	}
	return h
}
