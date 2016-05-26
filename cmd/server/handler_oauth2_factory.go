package server

import (
	"time"

	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/handler/core"
	oc "github.com/ory-am/fosite/handler/core/client"
	"github.com/ory-am/fosite/handler/core/explicit"
	"github.com/ory-am/fosite/handler/core/implicit"
	"github.com/ory-am/fosite/handler/core/refresh"
	"github.com/ory-am/fosite/handler/oidc"
	oe "github.com/ory-am/fosite/handler/oidc/explicit"
	"github.com/ory-am/fosite/handler/oidc/hybrid"
	oi "github.com/ory-am/fosite/handler/oidc/implicit"
	os "github.com/ory-am/fosite/handler/oidc/strategy"
	"github.com/ory-am/fosite/hash"
	"github.com/ory-am/fosite/token/jwt"
	"github.com/ory-am/hydra/client"
	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/internal"
	"github.com/ory-am/hydra/jwk"
	"github.com/ory-am/hydra/oauth2"
	"github.com/ory-am/hydra/pkg"
)

func injectFositeStore(c *config.Config, clients client.Manager) {
	var ctx = c.Context()
	var store pkg.FositeStorer

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
		break
	default:
		panic("Unknown connection type.")
	}

	ctx.FositeStore = store
}

func newOAuth2Handler(c *config.Config, router *httprouter.Router, km jwk.Manager) *oauth2.Handler {
	var ctx = c.Context()
	var store = ctx.FositeStore

	keys, err := km.GetKey(oauth2.OpenIDConnectKeyName, "private")
	if errors.Is(err, pkg.ErrNotFound) {
		logrus.Warnln("Could not find OpenID Connect singing keys. Generating a new keypair...")
		k, err := new(jwk.RS256Generator).Generate("")
		pkg.Must(err, "Could not generate signing key for OpenID Connect")
		km.AddKeySet(oauth2.OpenIDConnectKeyName, k)
		keys, err = km.GetKey(oauth2.OpenIDConnectKeyName, "private")
		pkg.Must(err, "Could not fetch signing key for OpenID Connect")
		logrus.Warnln("Keypair generated.")
		logrus.Warnln("WARNING: Automated key creation causes low entropy. Replace the keys as soon as possible.")
	} else {
		pkg.Must(err, "Could not fetch signing key for OpenID Connect")
	}

	rsaKey := jwk.MustRSAPrivate(jwk.First(keys.Keys))

	idStrategy := &os.DefaultStrategy{
		RS256JWTStrategy: &jwt.RS256JWTStrategy{
			PrivateKey: rsaKey,
		},
	}

	oauth2HandleHelper := &core.HandleHelper{
		AccessTokenStrategy: ctx.FositeStrategy,
		AccessTokenStorage:  store,
		AccessTokenLifespan: c.GetAccessTokenLifespan(),
	}

	oidcHelper := &oidc.IDTokenHandleHelper{IDTokenStrategy: idStrategy}

	explicitHandler := &explicit.AuthorizeExplicitGrantTypeHandler{
		AccessTokenStrategy:       ctx.FositeStrategy,
		RefreshTokenStrategy:      ctx.FositeStrategy,
		AuthorizeCodeStrategy:     ctx.FositeStrategy,
		AuthorizeCodeGrantStorage: store,
		AuthCodeLifespan:          time.Hour,
		AccessTokenLifespan:       c.GetAccessTokenLifespan(),
	}

	// The OpenID Connect Authorize Code Flow.
	oidcExplicit := &oe.OpenIDConnectExplicitHandler{
		OpenIDConnectRequestStorage: store,
		IDTokenHandleHelper:         oidcHelper,
	}

	implicitHandler := &implicit.AuthorizeImplicitGrantTypeHandler{
		AccessTokenStrategy: ctx.FositeStrategy,
		AccessTokenStorage:  store,
		AccessTokenLifespan: c.GetAccessTokenLifespan(),
	}

	consentURL, err := url.Parse(c.ConsentURL)
	pkg.Must(err, "Could not parse consent url.")

	handler := &oauth2.Handler{
		OAuth2: &fosite.Fosite{
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
					AccessTokenStrategy:      ctx.FositeStrategy,
					RefreshTokenStrategy:     ctx.FositeStrategy,
					RefreshTokenGrantStorage: store,
					AccessTokenLifespan:      c.GetAccessTokenLifespan(),
				},
				&oc.ClientCredentialsGrantHandler{
					HandleHelper: oauth2HandleHelper,
				},
			},
			AuthorizedRequestValidators: fosite.AuthorizedRequestValidators{
				&core.CoreValidator{
					AccessTokenStrategy: ctx.FositeStrategy,
					AccessTokenStorage:  store,
				},
			},
			Hasher: &hash.BCrypt{},
		},
		Consent: &oauth2.DefaultConsentStrategy{
			Issuer:     c.Issuer,
			KeyManager: km,
		},
		ConsentURL: *consentURL,
	}

	handler.SetRoutes(router)
	return handler
}
