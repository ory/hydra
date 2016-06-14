package cli

import (
	"fmt"

	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/connection"
	"github.com/ory-am/hydra/pkg"
	"github.com/pborman/uuid"
	"github.com/spf13/cobra"
)

type ConnectionHandler struct {
	Config *config.Config
	M      *connection.HTTPManager
}

func newConnectionHandler(c *config.Config) *ConnectionHandler {
	return &ConnectionHandler{
		Config: c,
		M:      &connection.HTTPManager{},
	}
}

func (h *ConnectionHandler) CreateConnection(cmd *cobra.Command, args []string) {
	h.M.Dry = *h.Config.Dry
	h.M.Client = h.Config.OAuth2Client(cmd)
	h.M.Endpoint = h.Config.Resolve("/connections")
	if len(args) != 3 {
		fmt.Print(cmd.UsageString())
		return
	}

	err := h.M.Create(&connection.Connection{
		ID:            uuid.New(),
		Provider:      args[0],
		LocalSubject:  args[1],
		RemoteSubject: args[2],
	})
	if h.M.Dry {
		fmt.Printf("%s\n", err)
		return
	}
	pkg.Must(err, "Could not create connection: %s", err)
}

func (h *ConnectionHandler) DeleteConnection(cmd *cobra.Command, args []string) {
	h.M.Dry = *h.Config.Dry
	h.M.Client = h.Config.OAuth2Client(cmd)
	h.M.Endpoint = h.Config.Resolve("/connections")
	if len(args) == 0 {
		fmt.Print(cmd.UsageString())
		return
	}

	for _, arg := range args {
		err := h.M.Delete(arg)
		if h.M.Dry {
			fmt.Printf("%s\n", err)
			continue
		}
		pkg.Must(err, "Could not delete connection: %s", err)
		fmt.Printf("Connection %s deleted.\n", arg)
	}
}
