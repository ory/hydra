package cli

import (
	"encoding/json"
	"fmt"

	"github.com/ory/hydra/config"
	"github.com/ory/hydra/pkg"
	hydra "github.com/ory/hydra/sdk/go/swagger"
	"github.com/spf13/cobra"
)

type JWKHandler struct {
	Config *config.Config
}

func (h *JWKHandler) newJwkManager(cmd *cobra.Command) *hydra.JsonWebKeyApi {
	c := hydra.NewJsonWebKeyApiWithBasePath(h.Config.ClusterURL)
	c.Configuration.Transport = h.Config.OAuth2Client(cmd).Transport
	return c
}

func newJWKHandler(c *config.Config) *JWKHandler {
	return &JWKHandler{Config: c}
}

func (h *JWKHandler) CreateKeys(cmd *cobra.Command, args []string) {
	m := h.newJwkManager(cmd)
	if len(args) < 1 || len(args) > 2 {
		fmt.Println(cmd.UsageString())
		return
	}

	kid := ""
	if len(args) == 2 {
		kid = args[1]
	}

	alg, _ := cmd.Flags().GetString("alg")
	keys, _, err := m.CreateJsonWebKeySet(args[0], hydra.CreateJsonWebKeySetPayload{Alg: alg, Kid: kid})
	pkg.Must(err, "Could not generate keys: %s", err)

	out, err := json.MarshalIndent(keys, "", "\t")
	pkg.Must(err, "Could not marshall keys: %s", err)

	fmt.Printf("%s\n", out)
}

func (h *JWKHandler) GetKeys(cmd *cobra.Command, args []string) {
	m := h.newJwkManager(cmd)
	if len(args) != 1 {
		fmt.Println(cmd.UsageString())
		return
	}

	keys, _, err := m.GetJsonWebKeySet(args[0])
	pkg.Must(err, "Could not generate keys: %s", err)

	out, err := json.MarshalIndent(keys, "", "\t")
	pkg.Must(err, "Could not marshall keys: %s", err)

	fmt.Printf("%s\n", out)
}

func (h *JWKHandler) DeleteKeys(cmd *cobra.Command, args []string) {
	m := h.newJwkManager(cmd)
	if len(args) != 1 {
		fmt.Println(cmd.UsageString())
		return
	}

	_, err := m.DeleteJsonWebKeySet(args[0])
	pkg.Must(err, "Could not generate keys: %s", err)
	fmt.Printf("Key set %s deleted.\n", args[0])
}
