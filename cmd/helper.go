package cmd

import (
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/fosite"
	"github.com/ory-am/ladon/memory"
	"github.com/ory-am/ladon"
	"github.com/ory-am/fosite/handler/core/strategy"
	"github.com/ory-am/fosite/token/hmac"
	"github.com/ory-am/hydra/client"
	"github.com/ory-am/fosite/handler/core"
	"github.com/ory-am/fosite/handler/oidc"
)

type fositeStore interface {
	fosite.Storage
	core.AccessTokenStorage
	core.RefreshTokenStorage
	core.AuthorizeCodeStorage
	oidc.OpenIDConnectRequestStorage
}

func newFositeStore(c config) fositeStore {
	if c.BackendURL == "" {
	}

	return pkg.FositeStore()
}

func newLadonStore(c config) ladon.Manager {
	if c.BackendURL == "" {
	}

	return &memory.Manager{
		Policies: make(map[string]ladon.Policy),
	}

}

func newClientStore(c config) client.Manager {
	if c.BackendURL == "" {
	}


	return &client.MemoryManager{
		Clients: make(map[string]*fosite.DefaultClient),
	}
}

func newHmacStrategy(c config) *strategy.HMACSHAStrategy {
	if c.BackendURL == "" {
	}

	return &strategy.HMACSHAStrategy{
		Enigma: &hmac.HMACStrategy{
			GlobalSecret: c.SystemSecret,
		},
	}

}