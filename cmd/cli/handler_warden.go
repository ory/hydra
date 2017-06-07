package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ory/hydra/config"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	"github.com/spf13/cobra"
)

type WardenHandler struct {
	Config *config.Config
}

func newWardenHandler(c *config.Config) *WardenHandler {
	return &WardenHandler{
		Config: c,
	}
}

func (h *WardenHandler) IsAuthorized(cmd *cobra.Command, args []string) {
	m := &oauth2.HTTPIntrospector{
		Endpoint: h.Config.Resolve("/oauth2/introspect"),
		Client:   h.Config.OAuth2Client(cmd),
	}

	if len(args) != 1 {
		fmt.Print(cmd.UsageString())
		return
	}

	scopes, _ := cmd.Flags().GetStringSlice("scopes")
	res, err := m.IntrospectToken(context.Background(), args[0], scopes...)
	pkg.Must(err, "Could not validate token: %s", err)

	out, err := json.MarshalIndent(res, "", "\t")
	pkg.Must(err, "Could not prettify token: %s", err)

	fmt.Printf("%s\n", out)
}
