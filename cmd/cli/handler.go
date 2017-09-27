package cli

import (
	"github.com/ory/hydra/config"
)

type Handler struct {
	Clients    *ClientHandler
	Policies   *PolicyHandler
	Keys       *JWKHandler
	Warden     *IntrospectionHandler
	Revocation *RevocationHandler
	Groups     *GroupHandler
	Migration  *MigrateHandler
}

func NewHandler(c *config.Config) *Handler {
	return &Handler{
		Clients:    newClientHandler(c),
		Policies:   newPolicyHandler(c),
		Keys:       newJWKHandler(c),
		Warden:     newIntrospectionHandler(c),
		Revocation: newRevocationHandler(c),
		Groups:     newGroupHandler(c),
		Migration:  newMigrateHandler(c),
	}
}
