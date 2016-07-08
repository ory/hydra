package cli

import (
	"github.com/ory-am/hydra/config"
)

type Handler struct {
	Clients     *ClientHandler
	Connections *ConnectionHandler
	Policies    *PolicyHandler
	Keys        *JWKHandler
	Warden      *WardenHandler
}

func NewHandler(c *config.Config) *Handler {
	return &Handler{
		Clients:     newClientHandler(c),
		Connections: newConnectionHandler(c),
		Policies:    newPolicyHandler(c),
		Keys:        newJWKHandler(c),
		Warden:      newWardenHandler(c),
	}
}
