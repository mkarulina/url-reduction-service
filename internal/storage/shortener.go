package storage

import (
	"bytes"
	"encoding/base64"
	"github.com/spf13/viper"
)

func (s *storage) ShortenLink(userID string, link string) (string, error) {
	var insertErr error

	buf := bytes.Buffer{}

	encoder := base64.NewEncoder(base64.URLEncoding, &buf)
	encoder.Write([]byte(link))
	key := buf.String()

	err := s.AddLinkToDB(&Link{UserID: userID, Key: key, Link: link})
	if err != nil {
		insertErr = err
	}

	shortLink := viper.GetString("BASE_URL") + "/" + key
	return shortLink, insertErr
}
