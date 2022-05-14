package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
)

type Encryptor interface {
	EncryptData(data []byte) ([]byte, error)
	DecryptData(data string) ([]byte, error)
	GenerateRandom(size int) ([]byte, error)
}

type encryptor struct {
	key []byte
}

func New() Encryptor {
	e := &encryptor{
		key: []byte("passphrasewhichneedstobe32bytes!"),
	}
	return e
}

func (e *encryptor) EncryptData(data []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	nonce, err := e.GenerateRandom(aesgcm.NonceSize())
	if err != nil {
		return nil, err
	}

	encrypted := aesgcm.Seal(nil, nonce, data, nil)
	encrypted = append(encrypted, nonce...)

	return encrypted, nil
}

func (e *encryptor) DecryptData(data string) ([]byte, error) {
	aesblock, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	encData, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}

	nonce := encData[len(encData)-aesgcm.NonceSize():]
	enc := encData[:len(encData)-aesgcm.NonceSize()]

	decrypted, err := aesgcm.Open(nil, nonce, enc, nil)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

func (e *encryptor) GenerateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
