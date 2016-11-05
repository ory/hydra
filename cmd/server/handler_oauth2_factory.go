package server

import (
	"fmt"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/compose"
	"github.com/ory-am/hydra/client"
	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/jwk"
	"github.com/ory-am/hydra/oauth2"
	"github.com/ory-am/hydra/pkg"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	r "gopkg.in/dancannon/gorethink.v2"
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
		m := &oauth2.FositeSQLStore{
			DB:      con.GetDatabase(),
			Manager: clients,
		}
		if err := m.CreateSchemas(); err != nil {
			logrus.Fatalf("Could not create oauth2 schema: %s", err)
		}
		store = m
		break
	case *config.RethinkDBConnection:
		con.CreateTableIfNotExists("hydra_oauth2_authorize_code")
		con.CreateTableIfNotExists("hydra_oauth2_id_sessions")
		con.CreateTableIfNotExists("hydra_oauth2_access_token")
		con.CreateTableIfNotExists("hydra_oauth2_implicit")
		con.CreateTableIfNotExists("hydra_oauth2_refresh_token")
		m := &oauth2.FositeRehinkDBStore{
			Session:             con.GetSession(),
			Manager:             clients,
			AuthorizeCodesTable: r.Table("hydra_oauth2_authorize_code"),
			IDSessionsTable:     r.Table("hydra_oauth2_id_sessions"),
			AccessTokensTable:   r.Table("hydra_oauth2_access_token"),
			RefreshTokensTable:  r.Table("hydra_oauth2_refresh_token"),
			AuthorizeCodes:      make(oauth2.RDBItems),
			IDSessions:          make(oauth2.RDBItems),
			AccessTokens:        make(oauth2.RDBItems),
			RefreshTokens:       make(oauth2.RDBItems),
		}
		if err := m.ColdStart(); err != nil {
			logrus.Fatalf("Could not fetch initial state: %s", err)
		}
		m.Watch(context.Background())
		store = m
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
		logrus.Warnln("Could not find OpenID Connect signing keys. Generating a new keypair...")
		keys, err = new(jwk.RS256Generator).Generate("")

		pkg.Must(err, "Could not generate signing key for OpenID Connect")
		km.AddKeySet(oauth2.OpenIDConnectKeyName, keys)
		logrus.Infoln("Keypair generated.")
		logrus.Warnln("WARNING: Automated key creation causes low entropy. Replace the keys as soon as possible.")
	} else {
		pkg.Must(err, "Could not fetch signing key for OpenID Connect")
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
		ConsentURL: *consentURL,
		H:          &herodot.JSON{},
	}

	handler.SetRoutes(router)
	return handler
}
