package getlink

import (
	"github.com/mkarulina/url-reduction-service/internal/storage"
	"net/http"
	"strings"
)

func GetLinkHandler(w http.ResponseWriter, r *http.Request) {
	stg := storage.New()
	//e := encryptor.New()
	//
	//cookie, err := r.Cookie("session_token")
	//if err != nil {
	//	log.Println("can't get cookie", err)
	//}
	//if validCookie := e.VerifyCookie(cookie.Value); !validCookie {
	//	w.WriteHeader(http.StatusMethodNotAllowed)
	//	return
	//}

	linkKey := strings.Split(r.URL.Path, "/")[1]
	foundLink := stg.GetLinkByKey(linkKey)

	if foundLink == "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Url not found"))
		return
	}

	w.Header().Set("Location", foundLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte(foundLink))
}
