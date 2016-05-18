package cli

import (
	"fmt"

	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/pkg"
	"github.com/spf13/cobra"
	"github.com/ory-am/hydra/client"
)

type ClientHandler struct {
	Config *config.Config
	M      *client.HTTPManager
}

func newClientHandler(c *config.Config) *ClientHandler {
	return &ClientHandler{
		Config: c,
		M:      &client.HTTPManager{},
	}
}

func (h *ClientHandler) CreateClient(cmd *cobra.Command, args []string) {
	var err error

	h.M.Endpoint = h.Config.Resolve("/clients")
	h.M.Client = h.Config.OAuth2Client(cmd)

	responseTypes, _ := cmd.Flags().GetStringSlice("response-types")
	grantTypes, _ := cmd.Flags().GetStringSlice("grant-types")
	allowedScopes, _ := cmd.Flags().GetStringSlice("allowed-scopes")
	callbacks, _ := cmd.Flags().GetStringSlice("callbacks")
	name, _ := cmd.Flags().GetString("name")
	id, _ := cmd.Flags().GetString("id")
	secretFlag, _ := cmd.Flags().GetString("secret")

	secret := []byte(secretFlag)
	if secretFlag == "" {
		secret, err = pkg.GenerateSecret(26)
		pkg.Must(err, "Could not generate secret: %s", err)
	}

	client := &fosite.DefaultClient{
		ID:            id,
		Secret:        secret,
		ResponseTypes: responseTypes,
		GrantedScopes: allowedScopes,
		GrantTypes:    grantTypes,
		RedirectURIs:  callbacks,
		Name:          name,
	}
	err = h.M.CreateClient(client)
	pkg.Must(err, "Could not create client: %s", err)

	fmt.Printf("Client ID: %s\n", client.ID)
	fmt.Printf("Client Secret: %s\n", secret)
}

func (h *ClientHandler) DeleteClient(cmd *cobra.Command, args []string) {
	h.M.Endpoint = h.Config.Resolve("/clients")
	h.M.Client = h.Config.OAuth2Client(cmd)
	if len(args) == 0 {
		fmt.Print(cmd.UsageString())
		return
	}

	for _, c := range args {
		err := h.M.DeleteClient(c)
		pkg.Must(err, "Could not delete client: %s", err)
	}

	fmt.Println("Client(s) deleted.")
}
