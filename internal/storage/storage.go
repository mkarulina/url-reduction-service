package storage

import (
	"github.com/spf13/viper"
	"io"
	"log"
	"sync"
)

type Storage interface {
	ShortenLink(userID string, link string) (string, error)
	AddLinkToDB(link *Link) error
	GetLinkByKey(linkKey string) string
	GetKeyByLink(link string) string
	GetAllUrlsByUserID(userID string) ([]ResponseLink, error)
}

type storage struct {
	mu        sync.RWMutex
	file      string
	dbAddress string
	urls      []Link
}

type Link struct {
	UserID string `json:"user_id"`
	Key    string `json:"key"`
	Link   string `json:"link"`
}

type ResponseLink struct {
	Key  string `json:"short_url"`
	Link string `json:"original_url"`
}

func New() Storage {
	s := &storage{
		mu:        sync.RWMutex{},
		file:      viper.GetString("FILE_STORAGE_PATH"),
		dbAddress: viper.GetString("DATABASE_DSN"),
		urls:      []Link{},
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
				s.urls = append(s.urls, Link{link.UserID, link.Key, link.Link})
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

		if err = AddURLToTable(&Link{link.UserID, link.Key, link.Link}); err != nil {
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
			for _, v := range s.urls {
				if v.Link == link {
					return v.Key
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
			for _, v := range s.urls {
				if v.Key == linkKey {
					return v.Link
				}
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

func (s *storage) GetAllUrlsByUserID(userID string) ([]ResponseLink, error) {
	var response []ResponseLink

	baseURL := viper.GetString("BASE_URL")

	if s.dbAddress == "" {
		if s.file == "" {
			for _, v := range s.urls {
				if v.UserID == userID {
					response = append(response, ResponseLink{Key: baseURL + "/" + v.Key, Link: v.Link})
				}
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
			if readLine.UserID == userID {
				response = append(response, ResponseLink{baseURL + "/" + readLine.Key, readLine.Link})
			}
		}
		return response, nil
	}

	resp, err := GetAllRowsByUserID(userID)
	if err != nil {
		return nil, err
	}

	for _, v := range resp {
		response = append(response, ResponseLink{baseURL + "/" + v.Key, v.Link})
	}

	return response, nil
}
