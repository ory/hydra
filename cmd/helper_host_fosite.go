package cmd

import (
	"time"

	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/handler/core"
	"github.com/ory-am/fosite/handler/core/client"
	"github.com/ory-am/fosite/handler/core/explicit"
	"github.com/ory-am/fosite/handler/core/implicit"
	"github.com/ory-am/fosite/handler/core/refresh"
	"github.com/ory-am/fosite/handler/oidc"
	oe "github.com/ory-am/fosite/handler/oidc/explicit"
	"github.com/ory-am/fosite/handler/oidc/hybrid"
	oi "github.com/ory-am/fosite/handler/oidc/implicit"
	"github.com/ory-am/fosite/hash"
	"github.com/ory-am/hydra/key"
	"github.com/ory-am/hydra/oauth2"
)

func newFosite(c *configuration, strategy core.CoreStrategy, idStrategy oidc.OpenIDConnectTokenStrategy, fositeStore fositeStorer, hasher hash.Hasher) fosite.OAuth2Provider {
	oauth2HandleHelper := &core.HandleHelper{
		AccessTokenStrategy: strategy,
		AccessTokenStorage:  fositeStore,
		AccessTokenLifespan: c.GetAccessTokenLifespan(),
	}

	oidcHelper := &oidc.IDTokenHandleHelper{IDTokenStrategy: idStrategy}

	explicitHandler := &explicit.AuthorizeExplicitGrantTypeHandler{
		AccessTokenStrategy:       strategy,
		RefreshTokenStrategy:      strategy,
		AuthorizeCodeStrategy:     strategy,
		AuthorizeCodeGrantStorage: fositeStore,
		AuthCodeLifespan:          time.Hour,
		AccessTokenLifespan:       c.GetAccessTokenLifespan(),
	}

	// The OpenID Connect Authorize Code Flow.
	oidcExplicit := &oe.OpenIDConnectExplicitHandler{
		OpenIDConnectRequestStorage: fositeStore,
		IDTokenHandleHelper:         oidcHelper,
	}

	implicitHandler := &implicit.AuthorizeImplicitGrantTypeHandler{
		AccessTokenStrategy: strategy,
		AccessTokenStorage:  fositeStore,
		AccessTokenLifespan: c.GetAccessTokenLifespan(),
	}

	return &fosite.Fosite{
		Store:          fositeStore,
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
				RefreshTokenGrantStorage: fositeStore,
				AccessTokenLifespan:      c.GetAccessTokenLifespan(),
			},
			&client.ClientCredentialsGrantHandler{
				HandleHelper: oauth2HandleHelper,
			},
		},
		AuthorizedRequestValidators: fosite.AuthorizedRequestValidators{&core.CoreValidator{
			AccessTokenStrategy: strategy,
			AccessTokenStorage:  fositeStore,
		}},
		Hasher: hasher,
	}
}

func newOAuth2Handler(c *configuration, fosite fosite.OAuth2Provider, keyManager key.Manager) *oauth2.Handler {
	return &oauth2.Handler{
		OAuth2: fosite,
		Consent: &oauth2.DefaultConsentStrategy{
			Issuer:     c.Issuer,
			KeyManager: keyManager,
		},
	}
}
