package cli

import (
	"github.com/ory-am/hydra/client"
	"github.com/ory-am/hydra/config"
)

type Handler struct {
	Clients *client.CLIHandler
}

func NewHandler(c *config.Config) *Handler {
	return &Handler{
		Clients: client.NewCLIHandler(c),
	}
}