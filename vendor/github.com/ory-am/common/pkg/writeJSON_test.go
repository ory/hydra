package pkg

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type rw struct {
	written []byte
	code    int
	header  http.Header
}

func (r *rw) Header() http.Header {
	return r.header
}

func (r *rw) Write(w []byte) (int, error) {
	r.written = w
	return 0, nil
}

func (r *rw) WriteHeader(c int) {
	r.code = c
}

func TestWriteJSON(t *testing.T) {
	r := &rw{header: http.Header{}}
	js := struct {
		Foo string `json:"foo"`
	}{"bar"}
	WriteJSON(r, js)
	assert.Equal(t, http.StatusOK, r.code)
	assert.Equal(t, `{"foo":"bar"}`, string(r.written))
	assert.Equal(t, r.Header().Get("Content-Type"), "application/json")

	WriteJSON(r, func() {})
	assert.Equal(t, http.StatusInternalServerError, r.code)

	WriteCreatedJSON(r, "location", js)
	assert.Equal(t, http.StatusCreated, r.code)
	assert.Equal(t, `{"foo":"bar"}`, string(r.written))
	assert.Equal(t, r.Header().Get("Content-Type"), "application/json")
	assert.Equal(t, r.Header().Get("Location"), "location")
}
