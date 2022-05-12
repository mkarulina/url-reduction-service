package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
)

type Encryptor interface {
	SetCookie(w http.ResponseWriter, r *http.Request) error
	VerifyCookie(cookieValue string) bool
	EncryptData(data []byte) ([]byte, error)
	DecryptData(data []byte) ([]byte, error)
	generateRandom(size int) ([]byte, error)
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

func (e *encryptor) SetCookie(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Println("can't get cookie", err)
	}

	if cookie != nil && len(cookie.Value) >= 16 {
		if valid := e.VerifyCookie(cookie.Value); valid {
			return nil
		}
	}

	random, err := e.generateRandom(16)
	if err != nil {
		return err
	}

	newUser := hex.EncodeToString(random)

	token, err := e.EncryptData([]byte(newUser))
	if err != nil {
		return err
	}

	newCookie := &http.Cookie{
		Name:  "session_token",
		Value: string(token),
	}

	http.SetCookie(w, newCookie)
	return nil
}

func (e *encryptor) VerifyCookie(cookieValue string) bool {
	isValid := false

	decrCookie, err := e.DecryptData([]byte(cookieValue))
	if err != nil {
		return false
	}
	if decrCookie != nil {
		isValid = true
	}
	return isValid
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

	nonce, err := e.generateRandom(aesgcm.NonceSize())
	if err != nil {
		return nil, err
	}

	encrypted := aesgcm.Seal(nil, nonce, data, nil)
	encrypted = append(encrypted, nonce...)

	return encrypted, nil
}

func (e *encryptor) DecryptData(data []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	nonce := data[len(data)-aesgcm.NonceSize():]
	enc := data[:len(data)-aesgcm.NonceSize()]

	decrypted, err := aesgcm.Open(nil, nonce, enc, nil)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

func (e *encryptor) generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
