package cli

import (
	"github.com/ory-am/hydra/client"
	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/connection"
	"github.com/ory-am/hydra/policy"
)

type Handler struct {
	Clients     *client.CLIHandler
	Connections *connection.CLIHandler
	Policies    *policy.CLIHandler
}

func NewHandler(c *config.Config) *Handler {
	return &Handler{
		Clients: client.NewCLIHandler(c),
		Connections: connection.NewCLIHandler(c),
		Policies: policy.NewCLIHandler(c),
	}
}
