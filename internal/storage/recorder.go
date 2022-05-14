package storage

import (
	"encoding/json"
	"os"
)

type Recorder struct {
	file    *os.File
	encoder *json.Encoder
}

type Reader struct {
	file    *os.File
	decoder *json.Decoder
}

func NewRecorder(fileName string) (*Recorder, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &Recorder{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func NewReader(fileName string) (*Reader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &Reader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (r *Recorder) WriteLink(link *Link) error {
	return r.encoder.Encode(&link)
}

func (r *Recorder) Close() error {
	return r.file.Close()
}

func (r *Reader) ReadLink() (*Link, error) {
	link := &Link{}
	if err := r.decoder.Decode(&link); err != nil {
		return nil, err
	}
	return link, nil
}

func (r *Reader) Close() error {
	return r.file.Close()
}
