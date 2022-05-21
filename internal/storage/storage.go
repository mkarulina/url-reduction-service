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
	GetLinkByKey(linkKey string) *Link
	GetKeyByLink(link string) string
	GetAllUrlsByUserID(userID string) ([]ResponseLink, error)
	DeleteUrls(inputChs ...chan UserKeys) error
}

type storage struct {
	mu        sync.RWMutex
	file      string
	dbAddress string
	urls      []Link
}

type Link struct {
	UserID    string `json:"user_id"`
	Key       string `json:"key"`
	Link      string `json:"link"`
	IsDeleted bool   `json:"is_deleted"`
}

type ResponseLink struct {
	Key  string `json:"short_url"`
	Link string `json:"original_url"`
}

type UserKeys struct {
	Cookie string
	Keys   []string
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
	var errToReturn error

	doIncrement := func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		defer wg.Done()

		if s.dbAddress == "" {
			if s.file == "" {
				s.urls = append(s.urls, Link{
					UserID:    link.UserID,
					Key:       link.Key,
					Link:      link.Link,
					IsDeleted: link.IsDeleted,
				})
				return
			}

			recorder, err := NewRecorder(s.file)
			if err != nil {
				errToReturn = err
				return
			}
			defer recorder.Close()

			if err = recorder.WriteLink(link); err != nil {
				errToReturn = err
				return
			}
			return
		}

		if err := AddURLToTable(&Link{
			UserID: link.UserID,
			Key:    link.Key,
			Link:   link.Link,
		}); err != nil {
			errToReturn = err
			return
		}
	}

	wg.Add(1)
	go doIncrement()
	wg.Wait()

	return errToReturn
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
	if sqlResp.IsDeleted == true {

	}

	foundKey = sqlResp.Key
	return foundKey
}

func (s *storage) GetLinkByKey(linkKey string) *Link {
	if s.dbAddress == "" {
		if s.file == "" {
			for _, v := range s.urls {
				if v.Key == linkKey {
					return &v
				}
			}
			return nil
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
				return readLine
				break
			}
		}
	}

	sqlResp, err := FindValueInDB(linkKey)
	if err != nil {
		log.Panic(err)
	}
	return &sqlResp
}

func (s *storage) GetAllUrlsByUserID(userID string) ([]ResponseLink, error) {
	var response []ResponseLink

	baseURL := viper.GetString("BASE_URL")

	if s.dbAddress == "" {
		if s.file == "" {
			for _, v := range s.urls {
				if v.UserID == userID && v.IsDeleted != true {
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
			if readLine.UserID == userID && readLine.IsDeleted != true {
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

func (s *storage) DeleteUrls(inputChs ...chan UserKeys) error {
	var wg sync.WaitGroup
	var errToReturn error

	outCh := fanIn(inputChs...)

	doIncrement := func(ch UserKeys) {
		s.mu.Lock()
		defer s.mu.Unlock()
		defer wg.Done()

		if s.dbAddress == "" {
			if s.file == "" {
				for _, key := range ch.Keys {
					for i, v := range s.urls {
						if v.Key == key && v.UserID == ch.Cookie {
							s.urls[i].IsDeleted = true
						}
					}
				}
			}

			reader, err := NewReader(s.file)
			if err != nil {
				errToReturn = err
			}
			defer reader.Close()

			recorder, err := NewRecorder(s.file)
			if err != nil {
				errToReturn = err
			}
			defer recorder.Close()

			for _, key := range ch.Keys {
				readLine, err := reader.ReadLink()
				if readLine == nil {
					break
				}
				if err != nil && err != io.EOF {
					log.Fatal(err)
				}

				if readLine.Key == key && readLine.UserID == ch.Cookie {
					readLine.IsDeleted = true
					if err = recorder.WriteLink(readLine); err != nil {
						errToReturn = err
					}
				}
			}
		}

		if err := SetIsDeletedFlag(ch.Cookie, ch.Keys); err != nil {
			errToReturn = err
		}
	}

	for ch := range outCh {
		wg.Add(1)
		go doIncrement(ch)
		wg.Wait()
	}

	return errToReturn
}

func fanIn(inputChs ...chan UserKeys) chan UserKeys {
	outCh := make(chan UserKeys)

	go func() {
		wg := &sync.WaitGroup{}

		for _, ch := range inputChs {
			wg.Add(1)

			go func(ch chan UserKeys) {
				defer wg.Done()
				for i := range ch {
					outCh <- i
				}
			}(ch)
		}

		wg.Wait()
		close(outCh)
	}()

	return outCh
}
