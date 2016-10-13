package oauth2_test

import (
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	foauth2 "github.com/ory-am/fosite/handler/oauth2"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/oauth2"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/hydra/warden"
	"github.com/ory-am/ladon"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	goauth2 "golang.org/x/oauth2"
	"fmt"
)

var (
	introspectors = make(map[string]oauth2.Introspector)
	now           = time.Now().Round(time.Second)
	tokens        = pkg.Tokens(3)
	fositeStore   = pkg.FositeStore()
)

var ladonWarden = pkg.LadonWarden(map[string]ladon.Policy{
	"1": &ladon.DefaultPolicy{
		ID:        "1",
		Subjects:  []string{"alice"},
		Resources: []string{"matrix", "rn:hydra:token<.*>"},
		Actions:   []string{"create", "decide"},
		Effect:    ladon.AllowAccess,
	},
	"2": &ladon.DefaultPolicy{
		ID:        "2",
		Subjects:  []string{"siri"},
		Resources: []string{"<.*>"},
		Actions:   []string{"decide"},
		Effect:    ladon.AllowAccess,
	},
})

var localWarden = &warden.LocalWarden{
	Warden: ladonWarden,
	OAuth2: &fosite.Fosite{
		Store: fositeStore,
		TokenValidators: fosite.TokenValidators{
			0: &foauth2.CoreValidator{
				CoreStrategy:  pkg.HMACStrategy,
				CoreStorage:   fositeStore,
				ScopeStrategy: fosite.HierarchicScopeStrategy,
			},
		},
		ScopeStrategy: fosite.HierarchicScopeStrategy,
	},
	Issuer:              "tests",
	AccessTokenLifespan: time.Hour,
}

func init() {
	introspectors["local"] = &oauth2.LocalIntrospector{
		OAuth2:              localWarden.OAuth2,
		Issuer:              "tests",
		AccessTokenLifespan: time.Hour,
	}

	r := httprouter.New()
	serv := &oauth2.Handler{
		Firewall:     localWarden,
		H:            &herodot.JSON{},
		Introspector: introspectors["local"],
	}
	serv.SetRoutes(r)
	ts = httptest.NewServer(r)

	ar := fosite.NewAccessRequest(oauth2.NewSession("alice"))
	ar.GrantedScopes = fosite.Arguments{"core"}
	ar.RequestedAt = now
	ar.Client = &fosite.DefaultClient{ID: "siri"}
	ar.Session.(*oauth2.Session).Extra = map[string]interface{}{"foo": "bar"}
	fositeStore.CreateAccessTokenSession(nil, tokens[0][0], ar)

	ar2 := fosite.NewAccessRequest(oauth2.NewSession("siri"))
	ar2.GrantedScopes = fosite.Arguments{"core"}
	ar2.RequestedAt = now
	ar2.Session.(*oauth2.Session).Extra = map[string]interface{}{"foo": "bar"}
	ar2.Client = &fosite.DefaultClient{ID: "siri"}
	fositeStore.CreateAccessTokenSession(nil, tokens[1][0], ar2)

	ar3 := fosite.NewAccessRequest(oauth2.NewSession("siri"))
	ar3.GrantedScopes = fosite.Arguments{"core"}
	ar3.RequestedAt = now
	ar2.Session.(*oauth2.Session).Extra = map[string]interface{}{"foo": "bar"}
	ar3.Client = &fosite.DefaultClient{ID: "doesnt-exist"}
	ar3.Session.(*oauth2.Session).AccessTokenExpiry = time.Now().Add(-time.Hour)
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
			Expiry:      time.Now().Add(time.Hour),
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
				expectErr: false,
			},
			{
				token:     tokens[0][1],
				expectErr: false,
				assert: func(c *oauth2.Introspection) {
					assert.Equal(t, "alice", c.Subject)
					assert.Equal(t, "tests", c.Issuer)
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
