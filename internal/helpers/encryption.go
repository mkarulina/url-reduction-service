package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
)

var key = []byte("passphrasewhichneedstobe32bytes!")

func SetCookie(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Println("can't get cookie", err)
	}

	if cookie != nil {
		if valid := VerifyCookie(r); valid {
			return nil
		}
	}

	random, err := generateRandom(16)
	if err != nil {
		return err
	}

	newUser := hex.EncodeToString(random)

	token, err := EncryptData([]byte(newUser), key)
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

func VerifyCookie(r *http.Request) bool {
	isValid := false
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Println("can't get cookie", err)
		return false
	}

	decrCookie, err := DecryptData([]byte(cookie.Value), key)
	if err != nil {
		return false
	}
	if decrCookie != nil {
		isValid = true
	}

	return isValid
}

func EncryptData(data []byte, key []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	nonce, err := generateRandom(aesgcm.NonceSize())
	if err != nil {
		return nil, err
	}

	encrypted := aesgcm.Seal(nonce, nonce, data, nil)

	return encrypted, nil
}

func DecryptData(data []byte, key []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	nonce := data[:aesgcm.NonceSize()]

	decrypted, err := aesgcm.Open(nonce, nonce, data, nil)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
