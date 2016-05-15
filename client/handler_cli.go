package client

import (
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/pkg"
	"github.com/pborman/uuid"
	"github.com/spf13/cobra"
	"fmt"
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

func (h *CLIHandler) CreateClient(cmd *cobra.Command, args []string) {
	var err error

	h.M.Endpoint = h.Config.Resolve("/clients")
	h.M.Client = h.Config.OAuth2Client(cmd)

	responseTypes, _ := cmd.Flags().GetStringSlice("response-types")
	grantTypes, _ := cmd.Flags().GetStringSlice("grant-types")
	allowedScopes, _ := cmd.Flags().GetStringSlice("allowed-scopes")
	callbacks, _ := cmd.Flags().GetStringSlice("callbacks")
	name, _ := cmd.Flags().GetStringSlice("name")
	id, _ := cmd.Flags().GetStringSlice("id")
	secret, _ := cmd.Flags().GetStringSlice("secret")
	if secret == "" {
		secret, err = pkg.GenerateSecret(26)
		pkg.Must(err, "Could not generate secret: %s", err)
	}

	client := &fosite.DefaultClient{
		ID:     uuid.New(),
		Secret: secret,
		ResponseTypes: responseTypes,
		GrantedScopes: allowedScopes,
		GrantTypes: grantTypes,
		RedirectURIs: callbacks,
		Name: name,
	}
	err = h.M.CreateClient(client)
	pkg.Must(err, "Could not create client: %s", err)

	fmt.Printf("Client ID: %s\n", client.ID)
	fmt.Printf("Client Secret: %s\n", secret)
}

func (h *CLIHandler) DeleteClient(cmd *cobra.Command, args []string) {
	h.M.Endpoint = h.Config.Resolve("/clients")
	h.M.Client = h.Config.OAuth2Client(cmd)
	if len(args) == 0 {
		fmt.Print(cmd.UsageString())
		return
	}

	for _, c := range args{
		err := h.M.DeleteClient(c)
		pkg.Must(err, "Could not delete client: %s", err)
	}

	fmt.Println("Client(s) deleted.")
}
