package storage

import (
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"sync"
)

type Storage interface {
	AddLinkToDB(link *Link) error
	GetLinkByKey(linkKey string) string
	GetKeyByLink(link string) string
	GetAllUrls() ([]Link, error)
}

type storage struct {
	mu        sync.Mutex
	urls      map[string]string
	file      string
	dbAddress string
}

type Link struct {
	Key  string `json:"short_url"`
	Link string `json:"original_url"`
}

func New() Storage {
	s := &storage{
		mu:        sync.Mutex{},
		urls:      map[string]string{},
		file:      viper.GetString("FILE_STORAGE_PATH"),
		dbAddress: viper.GetString("DATABASE_DSN"),
	}
	return s
}

func (s *storage) AddLinkToDB(link *Link) error {
	var wg sync.WaitGroup
	var err error

	doIncrement := func() error {
		s.mu.Lock()
		defer s.mu.Unlock()

		if s.dbAddress == "" {
			if s.file == "" {
				s.urls[link.Key] = link.Link
				wg.Done()
				return nil
			}

			recorder, err := NewRecorder(s.file)
			if err != nil {
				return err
			}
			defer recorder.Close()

			if err = recorder.WriteLink(link); err != nil {
				return err
			}
			wg.Done()
			return nil
		}

		if err = AddURLToTable(&Link{Key: link.Key, Link: link.Link}); err != nil {
			wg.Done()
			return err
		}

		wg.Done()
		return nil
	}

	wg.Add(1)
	go doIncrement()
	wg.Wait()

	if err != nil {
		return err
	}
	return nil
}

func (s *storage) GetKeyByLink(link string) string {
	var foundKey string

	if s.dbAddress == "" {
		if s.file == "" {
			for key, value := range s.urls {
				fmt.Println("v:", value)
				if value == link {
					foundKey = key
					break
				}
			}
			return foundKey
		}

		reader, err := NewReader(s.file)
		if err != nil {
			log.Panic(err)
		}
		defer reader.Close()

		for {
			readLine, err := reader.ReadLink()
			if readLine == nil {
				break
			}
			if err != nil {
				log.Panic(err)
			}

			if readLine.Link == link {
				foundKey = readLine.Key
				break
			}
		}
		return foundKey
	}

	sqlResp, err := FindValueInDB(link)
	if err != nil {
		log.Panic(err)
	}
	foundKey = sqlResp.Key
	return foundKey
}

func (s *storage) GetLinkByKey(linkKey string) string {
	var foundLink string

	if s.dbAddress == "" {
		if s.file == "" {
			if val, found := s.urls[linkKey]; found {
				foundLink = val
			}
			return foundLink
		}

		reader, err := NewReader(s.file)
		if err != nil {
			log.Fatal(err)
		}
		defer reader.Close()

		for {
			readLine, err := reader.ReadLink()
			if readLine == nil {
				break
			}
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			if readLine.Key == linkKey {
				foundLink = readLine.Link
				break
			}
		}
		return foundLink
	}

	sqlResp, err := FindValueInDB(linkKey)
	if err != nil {
		log.Panic(err)
	}
	foundLink = sqlResp.Link
	return foundLink
}

func (s *storage) GetAllUrls() ([]Link, error) {
	var response []Link

	if s.dbAddress == "" {
		if s.file == "" {
			for key, value := range s.urls {
				response = append(response, Link{key, value})
			}
			return response, nil
		}

		reader, err := NewReader(s.file)
		if err != nil {
			log.Fatal(err)
		}
		defer reader.Close()

		for {
			readLine, err := reader.ReadLink()
			if readLine == nil {
				break
			}
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			response = append(response, Link{readLine.Key, readLine.Link})
		}
		return response, nil
	}

	response, err := GetAllRows()
	if err != nil {
		return nil, err
	}

	return response, nil
}
