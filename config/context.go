package config

import (
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/handler/oauth2"
	"github.com/ory-am/hydra/firewall"
	"github.com/ory-am/hydra/jwk"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/ladon"
)

type Context struct {
	Connection interface{}

	Hasher         fosite.Hasher
	Warden         firewall.Firewall
	LadonManager   ladon.Manager
	FositeStrategy oauth2.CoreStrategy
	FositeStore    pkg.FositeStorer
	KeyManager     jwk.Manager
}
