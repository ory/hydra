package server

import (
	"testing"
	"github.com/ory-am/hydra/config"
	"github.com/julienschmidt/httprouter"
)

func TestNewHandler(t *testing.T) {
	var c = new(config.Config)
	router := httprouter.New()
	serverHandler := &Handler{}
	serverHandler.Start(c, router)
}
