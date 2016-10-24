package oauth2_test

import (
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/compose"
	"github.com/ory-am/fosite/storage"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/oauth2"
	"github.com/ory-am/hydra/pkg"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	goauth2 "golang.org/x/oauth2"
)

var (
	introspectors = make(map[string]oauth2.Introspector)
	now           = time.Now().Round(time.Second)
	tokens        = pkg.Tokens(3)
	fositeStore   = storage.NewExampleStore()
)

func init() {
	introspectors = make(map[string]oauth2.Introspector)
	now = time.Now().Round(time.Second)
	tokens = pkg.Tokens(3)
	fositeStore = storage.NewExampleStore()
	r := httprouter.New()
	serv := &oauth2.Handler{
		OAuth2: compose.Compose(
			fc,
			fositeStore,
			&compose.CommonStrategy{
				CoreStrategy:               compose.NewOAuth2HMACStrategy(fc, []byte("1234567890123456789012345678901234567890")),
				OpenIDConnectTokenStrategy: compose.NewOpenIDConnectStrategy(pkg.MustRSAKey()),
			},
			compose.OAuth2AuthorizeExplicitFactory,
			compose.OAuth2TokenIntrospectionFactory,
		),
		H: &herodot.JSON{},
	}
	serv.SetRoutes(r)
	ts = httptest.NewServer(r)

	ar := fosite.NewAccessRequest(oauth2.NewSession("alice"))
	ar.GrantedScopes = fosite.Arguments{"core"}
	ar.RequestedAt = now
	ar.Client = &fosite.DefaultClient{ID: "siri"}
	ar.Session.SetExpiresAt(fosite.AccessToken, now.Add(time.Hour))
	ar.Session.(*oauth2.Session).Extra = map[string]interface{}{"foo": "bar"}
	fositeStore.CreateAccessTokenSession(nil, tokens[0][0], ar)

	ar2 := fosite.NewAccessRequest(oauth2.NewSession("siri"))
	ar2.GrantedScopes = fosite.Arguments{"core"}
	ar2.RequestedAt = now
	ar2.Session.(*oauth2.Session).Extra = map[string]interface{}{"foo": "bar"}
	ar2.Session.SetExpiresAt(fosite.AccessToken, now.Add(time.Hour))
	ar2.Client = &fosite.DefaultClient{ID: "siri"}
	fositeStore.CreateAccessTokenSession(nil, tokens[1][0], ar2)

	ar3 := fosite.NewAccessRequest(oauth2.NewSession("siri"))
	ar3.GrantedScopes = fosite.Arguments{"core"}
	ar3.RequestedAt = now
	ar3.Session.(*oauth2.Session).Extra = map[string]interface{}{"foo": "bar"}
	ar3.Client = &fosite.DefaultClient{ID: "doesnt-exist"}
	ar3.Session.SetExpiresAt(fosite.AccessToken, now.Add(-time.Hour))
	fositeStore.CreateAccessTokenSession(nil, tokens[2][0], ar3)

	conf := &goauth2.Config{
		Scopes:   []string{},
		Endpoint: goauth2.Endpoint{},
	}

	ep, err := url.Parse(ts.URL)
	if err != nil {
		logrus.Fatalf("%s", err)
	}
	introspectors["http"] = &oauth2.HTTPIntrospector{
		Endpoint: ep,
		Client: conf.Client(goauth2.NoContext, &goauth2.Token{
			AccessToken: tokens[1][1],
			Expiry:      now.Add(time.Hour),
			TokenType:   "bearer",
		}),
	}
}

func TestIntrospect(t *testing.T) {
	for k, w := range introspectors {
		for _, c := range []struct {
			token     string
			expectErr bool
			assert    func(*oauth2.Introspection)
		}{
			{
				token:     "invalid",
				expectErr: true,
			},
			{
				token:     tokens[2][1],
				expectErr: true,
			},
			{
				token:     tokens[1][1],
				expectErr: true,
			},
			{
				token:     tokens[0][1],
				expectErr: false,
			},
			{
				token:     tokens[0][1],
				expectErr: false,
				assert: func(c *oauth2.Introspection) {
					assert.Equal(t, "alice", c.Subject)
					//assert.Equal(t, "tests", c.Issuer)
					assert.Equal(t, now.Add(time.Hour).Unix(), c.ExpiresAt, "expires at")
					assert.Equal(t, now.Unix(), c.IssuedAt, "issued at")
					assert.Equal(t, map[string]interface{}{"foo": "bar"}, c.Extra)
				},
			},
		} {
			t.Run(fmt.Sprintf("case=%s", k), func(t *testing.T) {
				ctx, err := w.IntrospectToken(context.Background(), c.token)
				pkg.AssertError(t, c.expectErr, err)
				if err == nil && c.assert != nil {
					c.assert(ctx)
				}
			})
		}
	}
}
