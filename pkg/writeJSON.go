package pkg

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(rw http.ResponseWriter, data interface{}) {
	js, err := json.Marshal(data)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(js)
}
