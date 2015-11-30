package pkg

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(rw http.ResponseWriter, data interface{}) {
	writeJSON(rw, data, http.StatusOK)
}

func WriteCreatedJSON(rw http.ResponseWriter, url string, data interface{}) {
	rw.Header().Add("Location", url)
	writeJSON(rw, data, http.StatusCreated)
}

func writeJSON(rw http.ResponseWriter, data interface{}, code int) {
	js, err := json.Marshal(data)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	rw.Write(js)

}
