package fosite

import (
	"encoding/json"
	"net/http"
)

func (c *Fosite) WriteAccessResponse(rw http.ResponseWriter, requester AccessRequester, responder AccessResponder) {
	js, err := json.Marshal(responder.ToMap())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json;charset=UTF-8")
	rw.Header().Set("Cache-Control", "no-store")
	rw.Header().Set("Pragma", "no-cache")

	rw.WriteHeader(http.StatusOK)
	rw.Write(js)
}
