package handlers

import (
	"github.com/mkarulina/url-reduction-service/config"
	"sync"
)

type Container struct {
	mu   *sync.Mutex
	urls map[string]string
	file string
}

func NewContainer() Container {
	filePath := config.GetConfig("FILE_STORAGE_PATH")
	c := Container{
		mu:   new(sync.Mutex),
		urls: map[string]string{},
		file: filePath,
	}
	return c
}
