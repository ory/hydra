package cli

import (
	//"context"
	//"encoding/json"
	"fmt"

	"github.com/ory/hydra/config"
	//"github.com/ory/hydra/oauth2"
	//"github.com/ory/hydra/pkg"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/spf13/cobra"
)

type IntrospectionHandler struct {
	Config *config.Config
}

func newIntrospectionHandler(c *config.Config) *IntrospectionHandler {
	return &IntrospectionHandler{
		Config: c,
	}
}

func (h *IntrospectionHandler) IsAuthorized(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Print(cmd.UsageString())
		return
	}

	c := hydra.NewOAuth2ApiWithBasePath(h.Config.ClusterURL)
	c.Configuration.Username = h.Config.ClientID
	c.Configuration.Password = h.Config.ClientSecret

	//scopes, _ := cmd.Flags().GetStringSlice("scopes")
	//res, err := c.IntrospectToken(context.Background(), args[0], scopes...)
	//pkg.Must(err, "Could not validate token: %s", err)
	//
	//out, err := json.MarshalIndent(res, "", "\t")
	//pkg.Must(err, "Could not prettify token: %s", err)
	//
	//fmt.Printf("%s\n", out)
}
