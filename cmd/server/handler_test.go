package server

import (
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/config"
)

func TestStart(t *testing.T) {
	router := httprouter.New()
	h := &Handler{}
	h.Start(&config.Config{}, router)
}
