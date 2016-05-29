package cli

import (
	"github.com/ory-am/hydra/config"
)

type Handler struct {
	Clients     *ClientHandler
	Connections *ConnectionHandler
	Policies    *PolicyHandler
	Keys        *JWKHandler
}

func NewHandler(c *config.Config) *Handler {
	return &Handler{
		Clients:     newClientHandler(c),
		Connections: newConnectionHandler(c),
		Policies:    newPolicHandler(c),
		Keys:        newJWKHandler(c),
	}
}
