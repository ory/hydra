package server

import (
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/hydra/config"
)

func TestStart(t *testing.T) {
	router := httprouter.New()
	h := &Handler{
		Config: &config.Config{},
	}
	h.registerRoutes(router)
}
