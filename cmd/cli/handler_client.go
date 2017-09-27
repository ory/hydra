package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ory/hydra/config"
	"github.com/ory/hydra/pkg"
	hydra "github.com/ory/hydra/sdk/go/swagger"
	"github.com/spf13/cobra"
)

type ClientHandler struct {
	Config *config.Config
}

func newClientHandler(c *config.Config) *ClientHandler {
	return &ClientHandler{
		Config: c,
	}
}

func (h *ClientHandler) newClientManager(cmd *cobra.Command) *hydra.ClientsApi {
	c := hydra.NewClientsApiWithBasePath(h.Config.ClusterURL)
	c.Configuration.Transport = h.Config.OAuth2Client(cmd).Transport
	return c
}

func (h *ClientHandler) ImportClients(cmd *cobra.Command, args []string) {
	m := h.newClientManager(cmd)

	if len(args) == 0 {
		fmt.Print(cmd.UsageString())
		return
	}

	for _, path := range args {
		reader, err := os.Open(path)
		pkg.Must(err, "Could not open file %s: %s", path, err)
		var c hydra.OauthClient
		err = json.NewDecoder(reader).Decode(&c)
		pkg.Must(err, "Could not parse JSON: %s", err)

		result, _, err := m.CreateOAuthClient(c)
		pkg.Must(err, "Could not create client: %s", err)
		fmt.Printf("Imported client %s:%s from %s.\n", result.Id, result.ClientSecret, path)
	}
}

func (h *ClientHandler) CreateClient(cmd *cobra.Command, args []string) {
	var err error
	m := h.newClientManager(cmd)
	responseTypes, _ := cmd.Flags().GetStringSlice("response-types")
	grantTypes, _ := cmd.Flags().GetStringSlice("grant-types")
	allowedScopes, _ := cmd.Flags().GetStringSlice("allowed-scopes")
	callbacks, _ := cmd.Flags().GetStringSlice("callbacks")
	name, _ := cmd.Flags().GetString("name")
	secret, _ := cmd.Flags().GetString("secret")
	id, _ := cmd.Flags().GetString("id")
	public, _ := cmd.Flags().GetBool("is-public")

	if secret == "" {
		var secretb []byte
		secretb, err = pkg.GenerateSecret(26)
		pkg.Must(err, "Could not generate secret: %s", err)
		secret = string(secretb)
	} else {
		fmt.Println("You should not provide secrets using command line flags. The secret might leak to bash history and similar systems.")
	}

	cc := hydra.OauthClient{
		Id:            id,
		ClientSecret:  secret,
		ResponseTypes: responseTypes,
		Scope:         strings.Join(allowedScopes, " "),
		GrantTypes:    grantTypes,
		RedirectUris:  callbacks,
		ClientName:    name,
		Public:        public,
	}

	result, _, err := m.CreateOAuthClient(cc)
	pkg.Must(err, "Could not create client: %s", err)

	fmt.Printf("Client ID: %s\n", result.Id)
	fmt.Printf("Client Secret: %s\n", result.ClientSecret)
}

func (h *ClientHandler) DeleteClient(cmd *cobra.Command, args []string) {
	m := h.newClientManager(cmd)

	if len(args) == 0 {
		fmt.Print(cmd.UsageString())
		return
	}

	for _, c := range args {
		_, err := m.DeleteOAuthClient(c)
		pkg.Must(err, "Could not delete client: %s", err)
	}

	fmt.Println("Client(s) deleted.")
}

func (h *ClientHandler) GetClient(cmd *cobra.Command, args []string) {
	m := h.newClientManager(cmd)

	if len(args) == 0 {
		fmt.Print(cmd.UsageString())
		return
	}

	cl, _, err := m.GetOAuthClient(args[0])
	pkg.Must(err, "Could not delete client: %s", err)

	out, err := json.MarshalIndent(cl, "", "\t")
	pkg.Must(err, "Could not convert client to JSON: %s", err)

	fmt.Printf("%s\n", out)
}
