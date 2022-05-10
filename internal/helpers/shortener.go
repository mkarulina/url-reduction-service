package helpers

import (
	"bytes"
	"encoding/base64"
	"github.com/mkarulina/url-reduction-service/internal/storage"
	"github.com/spf13/viper"
)

func ShortenLink(link string) (string, error) {
	var insertErr error

	stg := storage.New()

	buf := bytes.Buffer{}

	encoder := base64.NewEncoder(base64.URLEncoding, &buf)
	encoder.Write([]byte(link))

	err := stg.AddLinkToDB(&storage.Link{Key: buf.String(), Link: link})
	if err != nil {
		insertErr = err
	}

	shortLink := viper.GetString("BASE_URL") + "/" + buf.String()
	return shortLink, insertErr
}
