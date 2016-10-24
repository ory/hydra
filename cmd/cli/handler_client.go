package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"strings"

	"github.com/ory-am/hydra/client"
	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/pkg"
	"github.com/spf13/cobra"
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

func (h *ClientHandler) ImportClients(cmd *cobra.Command, args []string) {
	h.M.Endpoint = h.Config.Resolve("/clients")
	h.M.Client = h.Config.OAuth2Client(cmd)
	if len(args) == 0 {
		fmt.Print(cmd.UsageString())
		return
	}

	for _, path := range args {
		reader, err := os.Open(path)
		pkg.Must(err, "Could not open file %s: %s", path, err)
		var c client.Client
		err = json.NewDecoder(reader).Decode(&c)
		pkg.Must(err, "Could not parse JSON: %s", err)

		err = h.M.CreateClient(&c)
		if h.M.Dry {
			fmt.Printf("%s\n", err)
			continue
		}
		pkg.Must(err, "Could not create client: %s", err)
		fmt.Printf("Imported client %s:%s from %s.\n", c.ID, c.Secret, path)
	}
}

func (h *ClientHandler) CreateClient(cmd *cobra.Command, args []string) {
	var err error

	h.M.Dry, _ = cmd.Flags().GetBool("dry")
	h.M.Endpoint = h.Config.Resolve("/clients")
	h.M.Client = h.Config.OAuth2Client(cmd)

	responseTypes, _ := cmd.Flags().GetStringSlice("response-types")
	grantTypes, _ := cmd.Flags().GetStringSlice("grant-types")
	allowedScopes, _ := cmd.Flags().GetStringSlice("allowed-scopes")
	callbacks, _ := cmd.Flags().GetStringSlice("callbacks")
	name, _ := cmd.Flags().GetString("name")
	id, _ := cmd.Flags().GetString("id")
	public, _ := cmd.Flags().GetBool("is-public")

	secret, err := pkg.GenerateSecret(26)
	pkg.Must(err, "Could not generate secret: %s", err)

	cc := &client.Client{
		ID:            id,
		Secret:        string(secret),
		ResponseTypes: responseTypes,
		Scope:         strings.Join(allowedScopes, " "),
		GrantTypes:    grantTypes,
		RedirectURIs:  callbacks,
		Name:          name,
		Public:        public,
	}
	err = h.M.CreateClient(cc)
	if h.M.Dry {
		fmt.Printf("%s\n", err)
		return
	}
	pkg.Must(err, "Could not create client: %s", err)

	fmt.Printf("Client ID: %s\n", cc.ID)
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
		if h.M.Dry {
			fmt.Printf("%s\n", err)
			continue
		}
		pkg.Must(err, "Could not delete client: %s", err)
	}

	fmt.Println("Client(s) deleted.")
}
