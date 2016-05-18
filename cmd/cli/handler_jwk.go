package cli

import (
	"encoding/json"
	"fmt"

	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/jwk"
	"github.com/ory-am/hydra/pkg"
	"github.com/spf13/cobra"
)

type CLIHandler struct {
	Config *config.Config
	M      *jwk.HTTPManager
}

func newJWKCLIHandler(c *config.Config) *CLIHandler {
	return &CLIHandler{
		Config: c,
		M:      &jwk.HTTPManager{},
	}
}

func (h *CLIHandler) CreateKeys(cmd *cobra.Command, args []string) {
	h.M.Endpoint = h.Config.Resolve("/keys")
	h.M.Client = h.Config.OAuth2Client(cmd)
	if len(args) == 0 {
		fmt.Println(cmd.UsageString())
		return
	}

	alg, _ := cmd.Flags().GetString("alg")
	keys, err := h.M.CreateKeys(args[0], alg)
	pkg.Must(err, "Could not generate keys: %s", err)

	out, err := json.MarshalIndent(keys, "", "\t")
	pkg.Must(err, "Could not marshall keys: %s", err)

	fmt.Printf("%s\n", out)
}

func (h *CLIHandler) GetKeys(cmd *cobra.Command, args []string) {
	h.M.Endpoint = h.Config.Resolve("/keys")
	h.M.Client = h.Config.OAuth2Client(cmd)
	if len(args) == 0 {
		fmt.Println(cmd.UsageString())
		return
	}

	keys, err := h.M.GetKeySet(args[0])
	pkg.Must(err, "Could not generate keys: %s", err)

	out, err := json.MarshalIndent(keys, "", "\t")
	pkg.Must(err, "Could not marshall keys: %s", err)

	fmt.Printf("%s\n", out)
}

func (h *CLIHandler) DeleteKeys(cmd *cobra.Command, args []string) {
	h.M.Endpoint = h.Config.Resolve("/keys")
	h.M.Client = h.Config.OAuth2Client(cmd)
	if len(args) == 0 {
		fmt.Println(cmd.UsageString())
		return
	}

	err := h.M.DeleteKeySet(args[0])
	pkg.Must(err, "Could not generate keys: %s", err)
	fmt.Println("Key set deleted.")
}
