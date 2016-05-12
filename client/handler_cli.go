package client

import (
	"github.com/ory-am/fosite"
	"github.com/spf13/cobra"
	"github.com/ory-am/hydra/config"
	"github.com/pborman/uuid"
	"github.com/ory-am/hydra/pkg"
)

type CLIHandler struct {
	Config config.Config
	M      *HTTPManager
}

func NewCLIHandler(c *config.Config) *CLIHandler {
	return &CLIHandler{
		Config: c,
		M: HTTPManager{
			Endpoint: pkg.JoinURL(c.ClusterURL, "/clients"),
		},
	}
}

func (h *CLIHandler) CreateClient(cmd *cobra.Command, args []string) {
	h.M.Client = h.Config.OAuth2Client()

	secret, err := pkg.GenerateSecret(26)
	pkg.Must(err, "Could not generate secret: %s", err)

	err = h.M.CreateClient(&fosite.DefaultClient{
		ID: uuid.New(),
		Secret: secret,
	})
	pkg.Must(err, "Could not create client: %s", err)
}

func (h *CLIHandler) DeleteClient(cmd *cobra.Command, args []string) {
	h.M.Client = h.Config.OAuth2Client()
	err := h.M.DeleteClient(args[0])
	pkg.Must(err, "Could not delete client: %s", err)
}
