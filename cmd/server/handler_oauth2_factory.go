// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"fmt"
	"net/url"

	"os"

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
	"github.com/ory/hydra/warden"
	"github.com/pkg/errors"
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
	case *config.PluginConnection:
		var err error
		if store, err = con.NewOAuth2Manager(clients); err != nil {
			c.GetLogger().Fatalf("Could not load client manager plugin %s", err)
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
		fmt.Fprintf(os.Stderr, `Could not fetch signing key for OpenID Connect - did you forget to run "hydra migrate sql" or forget to set the SYSTEM_SECRET? Got error: %s`+"\n", err.Error())
		os.Exit(1)
	}

	rsaKey := jwk.MustRSAPrivate(jwk.First(keys.Keys))
	fc := &compose.Config{
		AccessTokenLifespan:   c.GetAccessTokenLifespan(),
		AuthorizeCodeLifespan: c.GetAuthCodeLifespan(),
		IDTokenLifespan:       c.GetIDTokenLifespan(),
		HashCost:              c.BCryptWorkFactor,
		ScopeStrategy:         c.GetScopeStrategy(),
	}
	return compose.Compose(
		fc,
		store,
		&compose.CommonStrategy{
			CoreStrategy:               compose.NewOAuth2HMACStrategy(fc, c.GetSystemSecret()),
			OpenIDConnectTokenStrategy: compose.NewOpenIDConnectStrategy(rsaKey),
		},
		nil,
		compose.OAuth2AuthorizeExplicitFactory,
		compose.OAuth2AuthorizeImplicitFactory,
		compose.OAuth2ClientCredentialsGrantFactory,
		compose.OAuth2RefreshTokenGrantFactory,
		compose.OpenIDConnectExplicitFactory,
		compose.OpenIDConnectHybridFactory,
		compose.OpenIDConnectImplicitFactory,
		compose.OAuth2TokenRevocationFactory,
		warden.OAuth2TokenIntrospectionFactory,
	)
}

func newOAuth2Handler(c *config.Config, router *httprouter.Router, cm oauth2.ConsentRequestManager, o fosite.OAuth2Provider) *oauth2.Handler {
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
		ScopesSupported:  c.OpenIDDiscoveryScopesSupported,
		UserinfoEndpoint: c.OpenIDDiscoveryUserinfoEndpoint,
		ClaimsSupported:  c.OpenIDDiscoveryClaimsSupported,
		ForcedHTTP:       c.ForceHTTP,
		OAuth2:           o,
		ScopeStrategy:    c.GetScopeStrategy(),
		Consent: &oauth2.DefaultConsentStrategy{
			Issuer:                   c.Issuer,
			ConsentManager:           c.Context().ConsentManager,
			DefaultChallengeLifespan: c.GetChallengeTokenLifespan(),
			DefaultIDTokenLifespan:   c.GetIDTokenLifespan(),
		},
		ConsentURL:          *consentURL,
		H:                   herodot.NewJSONWriter(c.GetLogger()),
		AccessTokenLifespan: c.GetAccessTokenLifespan(),
		CookieStore:         sessions.NewCookieStore(c.GetCookieSecret()),
		Issuer:              c.Issuer,
		L:                   c.GetLogger(),
		W:                   c.Context().Warden,
		ResourcePrefix:      c.AccessControlResourcePrefix,
	}

	handler.SetRoutes(router)
	return handler
}
