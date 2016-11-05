package cli

import (
	"github.com/ory-am/hydra/config"
)

type Handler struct {
	Clients  *ClientHandler
	Policies *PolicyHandler
	Keys     *JWKHandler
	Warden   *WardenHandler
}

func NewHandler(c *config.Config) *Handler {
	return &Handler{
		Clients:  newClientHandler(c),
		Policies: newPolicyHandler(c),
		Keys:     newJWKHandler(c),
		Warden:   newWardenHandler(c),
	}
}
