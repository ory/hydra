package cli

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/oauth2"
	"github.com/ory-am/hydra/pkg"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
)

type RevocationHandler struct {
	Config *config.Config
	M      *oauth2.HTTPRecovator
}

func newRevocationHandler(c *config.Config) *RevocationHandler {
	return &RevocationHandler{
		Config: c,
		M:      &oauth2.HTTPRecovator{},
	}
}

func (h *RevocationHandler) RevokeToken(cmd *cobra.Command, args []string) {
	if ok, _ := cmd.Flags().GetBool("skip-tls-verify"); ok {
		fmt.Println("Warning: Skipping TLS Certificate Verification.")
		h.M.Client = &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}
	}

	h.M.Endpoint = h.Config.Resolve("/oauth2/revoke")
	h.M.Config = &clientcredentials.Config{
		ClientID:     h.Config.ClientID,
		ClientSecret: h.Config.ClientSecret,
	}

	if len(args) != 1 {
		fmt.Print(cmd.UsageString())
		return
	}

	token := args[0]
	err := h.M.RevokeToken(context.Background(), args[0])
	pkg.Must(err, "Could not revoke token: %s", err)
	fmt.Printf("Revoked token %s", token)
}
