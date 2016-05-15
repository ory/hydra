package policy

import (
	"github.com/ory-am/hydra/config"
	"github.com/spf13/cobra"
	"github.com/ory-am/hydra/pkg"
	"fmt"
	"os"
	"github.com/square/go-jose/json"
	"github.com/ory-am/ladon"
)

type CLIHandler struct {
	Config *config.Config
	M      *HTTPManager
}

func NewCLIHandler(c *config.Config) *CLIHandler {
	return &CLIHandler{
		Config: c,
		M: &HTTPManager{},
	}
}

func (h *CLIHandler) CreatePolicy(cmd *cobra.Command, args []string) {
	h.M.Endpoint = h.Config.Resolve("/policies")
	if len(args) == 0 {
		fmt.Print(cmd.UsageString())
		return
	}

	for _, path := range args {
		reader, err := os.Open(path)
		pkg.Must(err, "Could not open file %s: %s", path, err)
		var policy ladon.DefaultPolicy
		err = json.NewDecoder(reader).Decode(&policy)
		pkg.Must(err, "Could not parse JSON: %s", err)
		err = h.M.Create(&policy)
		pkg.Must(err, "Could not create policy: %s", err)
	}
}

func (h *CLIHandler) DeletePolicy(cmd *cobra.Command, args []string) {
	h.M.Endpoint = h.Config.Resolve("/policies")
	if len(args) == 0 {
		fmt.Print(cmd.UsageString())
		return
	}

	h.M.Client = h.Config.OAuth2Client(cmd)
	for _, arg := range args {
		err := h.M.Delete(arg)
		pkg.Must(err, "Could not delete policy: %s", err)
	}
}
