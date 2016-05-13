package oauth2

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/internal"
	"github.com/ory-am/fosite/handler/core"
	"time"
	"github.com/ory-am/fosite/handler/oidc/hybrid"
	"github.com/ory-am/fosite/handler/core/refresh"
	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/client"
	"github.com/ory-am/fosite/handler/oidc"
	"github.com/ory-am/fosite/handler/core/explicit"
	"github.com/ory-am/fosite/handler/core/implicit"
	oe "github.com/ory-am/fosite/handler/oidc/explicit"
	oi "github.com/ory-am/fosite/handler/oidc/implicit"
	"github.com/ory-am/fosite/handler/core/strategy"
	oidcstrategy "github.com/ory-am/fosite/handler/oidc/strategy"
	"github.com/ory-am/fosite/token/hmac"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/hydra/jwk"
	"github.com/ory-am/fosite/token/jwt"
)

type fositeStorer interface {
	core.AccessTokenStorage
	fosite.Storage
	explicit.AuthorizeCodeGrantStorage
	refresh.RefreshTokenGrantStorage
	implicit.ImplicitGrantStorage
	oidc.OpenIDConnectRequestStorage
}

func NewHandler(c *config.Config, router *httprouter.Router, clients client.Manager, keys jwk.Manager) *Handler {
	ctx := c.Context()

	var fos *fosite.Fosite
	var store fositeStorer
	h := &Handler{}
	h.SetRoutes(router)

	switch ctx.Connection.(type) {
	case *config.MemoryConnection:
		store = &internal.FositeMemoryStore{
			Manager:        clients,
			AuthorizeCodes: make(map[string]fosite.Requester),
			IDSessions:     make(map[string]fosite.Requester),
			AccessTokens:   make(map[string]fosite.Requester),
			Implicit:       make(map[string]fosite.Requester),
			RefreshTokens:  make(map[string]fosite.Requester),
		}
		break;
	default:
		panic("Unknown connection type.")
	}

	h.OAuth2 = fos
	strategy := &strategy.HMACSHAStrategy{
		Enigma: &hmac.HMACStrategy{
			GlobalSecret: c.GetSystemSecret(),
		},
	}

	key, err := keys.GetKey("openid-connect", "public")
	pkg.Must(err, "Could not fetch signing key for OpenID Connect")
	idStrategy :=  &oidcstrategy.DefaultStrategy{
		RS256JWTStrategy: &jwt.RS256JWTStrategy{
			PrivateKey: key.Key,
			PublicKey:  key.Public,
		},
	}

	oauth2HandleHelper := &core.HandleHelper{
		AccessTokenStrategy: strategy,
		AccessTokenStorage:  store,
		AccessTokenLifespan: c.GetAccessTokenLifespan(),
	}

	oidcHelper := &oidc.IDTokenHandleHelper{IDTokenStrategy: idStrategy}

	explicitHandler := &explicit.AuthorizeExplicitGrantTypeHandler{
		AccessTokenStrategy:       strategy,
		RefreshTokenStrategy:      strategy,
		AuthorizeCodeStrategy:     strategy,
		AuthorizeCodeGrantStorage: store,
		AuthCodeLifespan:          time.Hour,
		AccessTokenLifespan:       c.GetAccessTokenLifespan(),
	}
	fos.AuthorizeEndpointHandlers.Append(explicitHandler)
	fos.TokenEndpointHandlers.Append(explicitHandler)

	// The OpenID Connect Authorize Code Flow.
	oidcExplicit := &oe.OpenIDConnectExplicitHandler{
		OpenIDConnectRequestStorage: store,
		IDTokenHandleHelper:         oidcHelper,
	}
	fos.AuthorizeEndpointHandlers.Append(oidcExplicit)
	fos.TokenEndpointHandlers.Append(oidcExplicit)

	implicitHandler := &implicit.AuthorizeImplicitGrantTypeHandler{
		AccessTokenStrategy: strategy,
		AccessTokenStorage:  store,
		AccessTokenLifespan: c.GetAccessTokenLifespan(),
	}

	return &fosite.Fosite{
		Store:          store,
		MandatoryScope: "core",
		AuthorizeEndpointHandlers: fosite.AuthorizeEndpointHandlers{
			explicitHandler,
			implicitHandler,
			oidcExplicit,
			&oi.OpenIDConnectImplicitHandler{
				IDTokenHandleHelper:               oidcHelper,
				AuthorizeImplicitGrantTypeHandler: implicitHandler,
			},
			&hybrid.OpenIDConnectHybridHandler{
				IDTokenHandleHelper:               oidcHelper,
				AuthorizeExplicitGrantTypeHandler: explicitHandler,
				AuthorizeImplicitGrantTypeHandler: implicitHandler,
			},
		},
		TokenEndpointHandlers: fosite.TokenEndpointHandlers{
			explicitHandler,
			oidcExplicit,
			&refresh.RefreshTokenGrantHandler{
				AccessTokenStrategy:      strategy,
				RefreshTokenStrategy:     strategy,
				RefreshTokenGrantStorage: store,
				AccessTokenLifespan:      c.GetAccessTokenLifespan(),
			},
			&client.ClientCredentialsGrantHandler{
				HandleHelper: oauth2HandleHelper,
			},
		},
		AuthorizedRequestValidators: fosite.AuthorizedRequestValidators{&core.CoreValidator{
			AccessTokenStrategy: strategy,
			AccessTokenStorage:  store,
		}},
		Hasher: hasher,
	}


	return h
}
