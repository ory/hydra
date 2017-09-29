package cli

import (
	"fmt"

	"net/http"

	"github.com/ory/hydra/config"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/spf13/cobra"
)

type JWKHandler struct {
	Config *config.Config
}

func (h *JWKHandler) newJwkManager(cmd *cobra.Command) *hydra.JsonWebKeyApi {
	c := hydra.NewJsonWebKeyApiWithBasePath(h.Config.ClusterURL)
	c.Configuration.Transport = h.Config.OAuth2Client(cmd).Transport
	if term, _ := cmd.Flags().GetBool("fake-tls-termination"); term {
		c.Configuration.DefaultHeader["X-Forwarded-Proto"] = "https"
	}

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
	keys, response, err := m.CreateJsonWebKeySet(args[0], hydra.JsonWebKeySetGeneratorRequest{Alg: alg, Kid: kid})
	checkResponse(response, err, http.StatusCreated)
	fmt.Printf("%s\n", formatResponse(keys))
}

func (h *JWKHandler) GetKeys(cmd *cobra.Command, args []string) {
	m := h.newJwkManager(cmd)
	if len(args) != 1 {
		fmt.Println(cmd.UsageString())
		return
	}

	keys, response, err := m.GetJsonWebKeySet(args[0])
	checkResponse(response, err, http.StatusOK)
	fmt.Printf("%s\n", formatResponse(keys))
}

func (h *JWKHandler) DeleteKeys(cmd *cobra.Command, args []string) {
	m := h.newJwkManager(cmd)
	if len(args) != 1 {
		fmt.Println(cmd.UsageString())
		return
	}

	response, err := m.DeleteJsonWebKeySet(args[0])
	checkResponse(response, err, http.StatusNoContent)
	fmt.Printf("Key set %s deleted.\n", args[0])
}
