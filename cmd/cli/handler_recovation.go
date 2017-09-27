package cli

import (
	"fmt"
	"net/http"

	"github.com/ory/hydra/config"
	"github.com/ory/hydra/pkg"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/spf13/cobra"
)

type RevocationHandler struct {
	Config *config.Config
	M      *hydra.OAuth2Api
}

func newRevocationHandler(c *config.Config) *RevocationHandler {
	handler := hydra.NewOAuth2ApiWithBasePath(c.ClusterURL)
	handler.Configuration.Username = c.ClientID
	handler.Configuration.Password = c.ClientSecret
	return &RevocationHandler{M: handler}
}

func (h *RevocationHandler) RevokeToken(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Print(cmd.UsageString())
		return
	}

	token := args[0]
	response, err := h.M.RevokeOAuth2Token(args[0])
	pkg.Must(err, "Could not revoke token: %s", err)
	if response.StatusCode != http.StatusOK {
		fmt.Printf("Could not revoke token becase status code %d was received, expected %d.", response.StatusCode, http.StatusOK)
		return
	}

	fmt.Printf("Revoked token %s", token)
}
