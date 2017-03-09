package fosite

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Fosite) WriteAccessError(rw http.ResponseWriter, _ AccessRequester, err error) {
	writeJsonError(rw, err)
}

func writeJsonError(rw http.ResponseWriter, err error) {
	rw.Header().Set("Content-Type", "application/json;charset=UTF-8")

	rfcerr := ErrorToRFC6749Error(err)
	js, err := json.Marshal(rfcerr)
	if err != nil {
		http.Error(rw, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(rfcerr.StatusCode)
	rw.Write(js)
}
