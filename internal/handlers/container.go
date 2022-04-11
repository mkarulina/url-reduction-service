package handlers

import (
	"github.com/spf13/viper"
	"sync"
)

type Container struct {
	mu   *sync.Mutex
	urls map[string]string
	file string
}

func NewContainer() Container {
	filePath := viper.GetString("FILE_STORAGE_PATH")
	c := Container{
		mu:   new(sync.Mutex),
		urls: map[string]string{},
		file: filePath,
	}
	return c
}
