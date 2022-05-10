package allurls

import (
	"encoding/json"
	"github.com/mkarulina/url-reduction-service/internal/storage"
	"log"
	"net/http"
)

func GetAllUrlsHandler(w http.ResponseWriter, r *http.Request) {
	stg := storage.New()

	//if validCookie := helpers.VerifyCookie(r); !validCookie {
	//	w.WriteHeader(http.StatusMethodNotAllowed)
	//	return
	//}

	response, err := stg.GetAllUrls()
	if err != nil {
		log.Println("can't get urls", err)
		return
	}

	if len(response) == 0 {
		w.WriteHeader(http.StatusNoContent)
	}
	marshalResp, err := json.Marshal(response)
	if err != nil {
		log.Println("can't marshal response", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshalResp)
}
