package server

import (
	"fmt"
	"net/url"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/herodot"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/config"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	"github.com/pkg/errors"
	"os"
)

func injectFositeStore(c *config.Config, clients client.Manager) {
	var ctx = c.Context()
	var store pkg.FositeStorer

	switch con := ctx.Connection.(type) {
	case *config.MemoryConnection:
		store = &oauth2.FositeMemoryStore{
			Manager:        clients,
			AuthorizeCodes: make(map[string]fosite.Requester),
			IDSessions:     make(map[string]fosite.Requester),
			AccessTokens:   make(map[string]fosite.Requester),
			RefreshTokens:  make(map[string]fosite.Requester),
		}
		break
	case *config.SQLConnection:
		store = &oauth2.FositeSQLStore{
			DB:      con.GetDatabase(),
			Manager: clients,
			L:       c.GetLogger(),
		}
		break
	default:
		panic("Unknown connection type.")
	}

	ctx.FositeStore = store
}

func newOAuth2Provider(c *config.Config, km jwk.Manager) fosite.OAuth2Provider {
	var ctx = c.Context()
	var store = ctx.FositeStore

	createRS256KeysIfNotExist(c, oauth2.OpenIDConnectKeyName, "private", "sig")
	keys, err := km.GetKey(oauth2.OpenIDConnectKeyName, "private")
	if errors.Cause(err) == pkg.ErrNotFound {
		c.GetLogger().Warnln("Could not find OpenID Connect signing keys. Generating a new keypair...")
		keys, err = new(jwk.RS256Generator).Generate("")

		pkg.Must(err, "Could not generate signing key for OpenID Connect")
		km.AddKeySet(oauth2.OpenIDConnectKeyName, keys)
		c.GetLogger().Infoln("Keypair generated.")
		c.GetLogger().Warnln("WARNING: Automated key creation causes low entropy. Replace the keys as soon as possible.")
	} else if err != nil {
		fmt.Fprintf(os.Stderr, `Could not fetch signing key for OpenID Connect - did you forget to run "hydra migrate sql" or forget to set the SYSTEM_SECRET? Got error: %s` + "\n", err.Error())
		os.Exit(1)
	}

	rsaKey := jwk.MustRSAPrivate(jwk.First(keys.Keys))
	fc := &compose.Config{
		AccessTokenLifespan:   c.GetAccessTokenLifespan(),
		AuthorizeCodeLifespan: c.GetAuthCodeLifespan(),
		IDTokenLifespan:       c.GetIDTokenLifespan(),
		HashCost:              c.BCryptWorkFactor,
	}
	return compose.Compose(
		fc,
		store,
		&compose.CommonStrategy{
			CoreStrategy:               compose.NewOAuth2HMACStrategy(fc, c.GetSystemSecret()),
			OpenIDConnectTokenStrategy: compose.NewOpenIDConnectStrategy(rsaKey),
		},
		compose.OAuth2AuthorizeExplicitFactory,
		compose.OAuth2AuthorizeImplicitFactory,
		compose.OAuth2ClientCredentialsGrantFactory,
		compose.OAuth2RefreshTokenGrantFactory,
		compose.OpenIDConnectExplicitFactory,
		compose.OpenIDConnectHybridFactory,
		compose.OpenIDConnectImplicitFactory,
		compose.OAuth2TokenRevocationFactory,
		compose.OAuth2TokenIntrospectionFactory,
	)
}

func newOAuth2Handler(c *config.Config, router *httprouter.Router, km jwk.Manager, o fosite.OAuth2Provider) *oauth2.Handler {
	if c.ConsentURL == "" {
		proto := "https"
		if c.ForceHTTP {
			proto = "http"
		}
		host := "localhost"
		if c.BindHost != "" {
			host = c.BindHost
		}
		c.ConsentURL = fmt.Sprintf("%s://%s:%d/oauth2/consent", proto, host, c.BindPort)
	}

	consentURL, err := url.Parse(c.ConsentURL)
	pkg.Must(err, "Could not parse consent url %s.", c.ConsentURL)

	handler := &oauth2.Handler{
		ForcedHTTP: c.ForceHTTP,
		OAuth2:     o,
		Consent: &oauth2.DefaultConsentStrategy{
			Issuer:                   c.Issuer,
			KeyManager:               km,
			DefaultChallengeLifespan: c.GetChallengeTokenLifespan(),
			DefaultIDTokenLifespan:   c.GetIDTokenLifespan(),
		},
		ConsentURL:          *consentURL,
		H:                   herodot.NewJSONWriter(c.GetLogger()),
		AccessTokenLifespan: c.GetAccessTokenLifespan(),
		CookieStore:         sessions.NewCookieStore(c.GetCookieSecret()),
		Issuer:              c.Issuer,
	}

	handler.SetRoutes(router)
	return handler
}
