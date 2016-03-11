package pkg

import (
	"encoding/json"
	"net/http"
)

func WriteIndentJSON(rw http.ResponseWriter, data interface{}) {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(rw, js, http.StatusOK)
}

func WriteJSON(rw http.ResponseWriter, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(rw, js, http.StatusOK)
}

func WriteCreatedJSON(rw http.ResponseWriter, url string, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Location", url)
	writeJSON(rw, js, http.StatusCreated)
}

func writeJSON(rw http.ResponseWriter, js []byte, code int) {

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	rw.Write(js)
}
