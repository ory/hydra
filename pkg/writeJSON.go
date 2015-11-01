package pkg

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(rw http.ResponseWriter, connection interface{}) {
	js, err := json.Marshal(connection)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(js)
}
