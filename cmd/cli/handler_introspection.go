package cli

import (
	//"context"
	//"encoding/json"
	"fmt"

	"github.com/ory/hydra/config"
	//"github.com/ory/hydra/oauth2"

	"net/http"
	"strings"

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
	c.Configuration.Transport = h.Config.OAuth2Client(cmd).Transport

	if term, _ := cmd.Flags().GetBool("fake-tls-termination"); term {
		c.Configuration.DefaultHeader["X-Forwarded-Proto"] = "https"
	}

	scopes, _ := cmd.Flags().GetStringSlice("scopes")
	result, response, err := c.IntrospectOAuth2Token(args[0], strings.Join(scopes, " "))
	checkResponse(response, err, http.StatusOK)
	fmt.Printf("%s\n", formatResponse(result))
}
