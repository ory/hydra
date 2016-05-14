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
	h.M.Endpoint = h.Config.Resolve("/clients")
	h.M.Client = h.Config.OAuth2Client()

	secret, err := pkg.GenerateSecret(26)
	pkg.Must(err, "Could not generate secret: %s", err)

	client := &fosite.DefaultClient{
		ID:     uuid.New(),
		Secret: secret,
	}
	err = h.M.CreateClient(client)
	pkg.Must(err, "Could not create client: %s", err)

	fmt.Printf("Client ID: %s\n", client.ID)
	fmt.Printf("Client Secret: %s\n", secret)
}

func (h *CLIHandler) DeleteClient(cmd *cobra.Command, args []string) {
	h.M.Endpoint = h.Config.Resolve("/clients")
	h.M.Client = h.Config.OAuth2Client()
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
