package cli

import (
	"encoding/json"
	"fmt"

	"github.com/ory/hydra/config"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/pkg"
	"github.com/spf13/cobra"
)

type JWKHandler struct {
	Config *config.Config
}

func (h *JWKHandler) newJwkManager(cmd *cobra.Command) *jwk.HTTPManager {
	dry, _ := cmd.Flags().GetBool("dry")
	term, _ := cmd.Flags().GetBool("fake-tls-termination")

	return &jwk.HTTPManager{
		Dry:                dry,
		Endpoint:           h.Config.Resolve("/keys"),
		Client:             h.Config.OAuth2Client(cmd),
		FakeTLSTermination: term,
	}
}

func newJWKHandler(c *config.Config) *JWKHandler {
	return &JWKHandler{
		Config: c,
	}
}

func (h *JWKHandler) CreateKeys(cmd *cobra.Command, args []string) {
	m := h.newJwkManager(cmd)
	if len(args) == 0 {
		fmt.Println(cmd.UsageString())
		return
	}

	alg, _ := cmd.Flags().GetString("alg")
	keys, err := m.CreateKeys(args[0], alg)
	if m.Dry {
		fmt.Printf("%s\n", err)
		return
	}
	pkg.Must(err, "Could not generate keys: %s", err)

	out, err := json.MarshalIndent(keys, "", "\t")
	pkg.Must(err, "Could not marshall keys: %s", err)

	fmt.Printf("%s\n", out)
}

func (h *JWKHandler) GetKeys(cmd *cobra.Command, args []string) {
	m := h.newJwkManager(cmd)
	if len(args) == 0 {
		fmt.Println(cmd.UsageString())
		return
	}

	keys, err := m.GetKeySet(args[0])
	if m.Dry {
		fmt.Printf("%s\n", err)
		return
	}
	pkg.Must(err, "Could not generate keys: %s", err)

	out, err := json.MarshalIndent(keys, "", "\t")
	pkg.Must(err, "Could not marshall keys: %s", err)

	fmt.Printf("%s\n", out)
}

func (h *JWKHandler) DeleteKeys(cmd *cobra.Command, args []string) {
	m := h.newJwkManager(cmd)
	if len(args) == 0 {
		fmt.Println(cmd.UsageString())
		return
	}

	err := m.DeleteKeySet(args[0])
	if m.Dry {
		fmt.Printf("%s\n", err)
		return
	}
	pkg.Must(err, "Could not generate keys: %s", err)
	fmt.Println("Key set deleted.")
}
