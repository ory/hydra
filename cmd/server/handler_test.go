package server

import (
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/hydra/config"
)

func TestStart(t *testing.T) {
	router := httprouter.New()
	h := &Handler{
		Config: &config.Config{
			DatabaseURL: "memory",
		},
	}
	h.registerRoutes(router)
}
