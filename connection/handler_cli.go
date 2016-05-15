package connection

import (
	"fmt"

	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/pkg"
	"github.com/pborman/uuid"
	"github.com/spf13/cobra"
)

type CLIHandler struct {
	Config *config.Config
	M      *HTTPManager
}

func NewCLIHandler(c *config.Config) *CLIHandler {
	return &CLIHandler{
		Config: c,
		M:      &HTTPManager{},
	}
}

func (h *CLIHandler) CreateConnection(cmd *cobra.Command, args []string) {
	h.M.Client = h.Config.OAuth2Client(cmd)
	h.M.Endpoint = h.Config.Resolve("/connections")
	if len(args) != 3 {
		fmt.Print(cmd.UsageString())
		return
	}

	err := h.M.Create(&Connection{
		ID:            uuid.New(),
		Provider:      args[0],
		LocalSubject:  args[1],
		RemoteSubject: args[2],
	})
	pkg.Must(err, "Could not create connection: %s", err)
}

func (h *CLIHandler) DeleteConnection(cmd *cobra.Command, args []string) {
	h.M.Client = h.Config.OAuth2Client(cmd)
	h.M.Endpoint = h.Config.Resolve("/connections")
	if len(args) == 0 {
		fmt.Print(cmd.UsageString())
		return
	}

	for _, arg := range args {
		err := h.M.Delete(arg)
		pkg.Must(err, "Could not delete connection: %s", err)
		fmt.Printf("Connection %s deleted.\n", arg)
	}
}
