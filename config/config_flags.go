package config

import (
	"flag"
	"os"
)

type Config struct {
	serverAddress string `yaml:"SERVER_ADDRESS"`
	baseURL       string `yaml:"BASE_URL"`
	filePath      string `yaml:"FILE_STORAGE_PATH"`
}

var Flags map[string]string

func SetFlags() {
	conf := &Config{}

	os.Args[0] = "shortener"
	flag.StringVar(&conf.serverAddress, "a", ":8080", "port to listen on")
	flag.StringVar(&conf.baseURL, "b", "http://localhost:8080/", "base short url")
	flag.StringVar(&conf.filePath, "f", "../urls.log", "file path for url saving")

	flag.Parse()

	Flags = make(map[string]string)
	Flags["a"] = "SERVER_ADDRESS"
	Flags["b"] = "BASE_URL"
	Flags["f"] = "FILE_STORAGE_PATH"
}
