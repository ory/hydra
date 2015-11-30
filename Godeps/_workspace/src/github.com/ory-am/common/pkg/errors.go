package pkg

import (
	"github.com/go-errors/errors"
	"net/http"
)

var (
	ErrNotFound = errors.New("not found.")
)

func WriteError(w http.ResponseWriter, err error) {
	if err == ErrNotFound {
		LogError(err, http.StatusNotFound)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	LogError(err, http.StatusInternalServerError)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func HttpError(rw http.ResponseWriter, err error, code int) {
	LogError(err, code)
	http.Error(rw, err.Error(), code)
}
