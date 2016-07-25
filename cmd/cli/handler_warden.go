package cli

import (
	"encoding/json"
	"fmt"

	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/hydra/warden"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

type WardenHandler struct {
	Config *config.Config
	M      *warden.HTTPWarden
}

func newWardenHandler(c *config.Config) *WardenHandler {
	return &WardenHandler{
		Config: c,
		M:      &warden.HTTPWarden{},
	}
}

func (h *WardenHandler) IsAuthorized(cmd *cobra.Command, args []string) {
	h.M.Dry, _ = cmd.Flags().GetBool("dry")
	h.M.Client = h.Config.OAuth2Client(cmd)
	h.M.Endpoint = h.Config.Resolve("/connections")

	if len(args) != 1 {
		fmt.Print(cmd.UsageString())
		return
	}

	scopes, _ := cmd.Flags().GetStringSlice("scopes")
	if len(scopes) == 0 {
		scopes = []string{"core"}
	}

	res, err := h.M.Authorized(context.Background(), args[0], scopes...)
	pkg.Must(err, "Could not validate token: %s", err)

	out, err := json.MarshalIndent(res, "", "\t")
	pkg.Must(err, "Could not marshall keys: %s", err)

	fmt.Printf("%s\n", out)
}
