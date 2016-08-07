package warden_test

import (
	"log"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	foauth2 "github.com/ory-am/fosite/handler/oauth2"
	"github.com/ory-am/hydra/firewall"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/oauth2"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/hydra/warden"
	"github.com/ory-am/ladon"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	coauth2 "golang.org/x/oauth2"
)

var ts *httptest.Server

var wardens = map[string]firewall.Firewall{}

var ladonWarden = pkg.LadonWarden(map[string]ladon.Policy{
	"1": &ladon.DefaultPolicy{
		ID:        "1",
		Subjects:  []string{"alice"},
		Resources: []string{"matrix"},
		Actions:   []string{"create"},
		Effect:    ladon.AllowAccess,
	},
	"2": &ladon.DefaultPolicy{
		ID:        "2",
		Subjects:  []string{"siri"},
		Resources: []string{"<.*>"},
		Actions:   []string{},
		Effect:    ladon.AllowAccess,
	},
})

var fositeStore = pkg.FositeStore()

var now = time.Now().Round(time.Second)

var tokens = pkg.Tokens(3)

func init() {
	wardens["local"] = &warden.LocalWarden{
		Warden: ladonWarden,
		OAuth2: &fosite.Fosite{
			Store: fositeStore,
			TokenValidators: fosite.TokenValidators{
				&foauth2.CoreValidator{
					CoreStrategy: pkg.HMACStrategy,
					CoreStorage:  fositeStore,
				},
			},
		},
		Issuer:              "tests",
		AccessTokenLifespan: time.Hour,
	}

	r := httprouter.New()
	serv := &warden.WardenHandler{
		H:      &herodot.JSON{},
		Warden: wardens["local"],
	}
	serv.SetRoutes(r)
	ts = httptest.NewServer(r)

	url, err := url.Parse(ts.URL + warden.TokenAllowedHandlerPath)
	if err != nil {
		log.Fatalf("%s", err)
	}

	ar := fosite.NewAccessRequest(oauth2.NewSession("alice"))
	ar.GrantedScopes = fosite.Arguments{"core"}
	ar.RequestedAt = now
	fositeStore.CreateAccessTokenSession(nil, tokens[0][0], ar)

	ar2 := fosite.NewAccessRequest(oauth2.NewSession("siri"))
	ar2.GrantedScopes = fosite.Arguments{"core"}
	ar2.RequestedAt = now
	fositeStore.CreateAccessTokenSession(nil, tokens[1][0], ar2)

	ar3 := fosite.NewAccessRequest(oauth2.NewSession("siri"))
	ar3.GrantedScopes = fosite.Arguments{"core"}
	ar3.RequestedAt = now
	ar3.Session.(*oauth2.Session).AccessTokenExpiry = time.Now().Add(-time.Hour)
	fositeStore.CreateAccessTokenSession(nil, tokens[2][0], ar3)

	conf := &coauth2.Config{
		Scopes:   []string{},
		Endpoint: coauth2.Endpoint{},
	}
	wardens["http"] = &warden.HTTPWarden{
		Endpoint: url,
		Client: conf.Client(coauth2.NoContext, &coauth2.Token{
			AccessToken: tokens[1][1],
			Expiry:      time.Now().Add(time.Hour),
			TokenType:   "bearer",
		}),
	}
}

func TestActionAllowed(t *testing.T) {
	for n, w := range wardens {
		for k, c := range []struct {
			token     string
			req       *ladon.Request
			scopes    []string
			expectErr bool
			assert    func(*firewall.Context)
		}{
			{
				token:     "invalid",
				req:       &ladon.Request{},
				scopes:    []string{},
				expectErr: true,
			},
			{
				token: "invalid",
				req: &ladon.Request{
					Subject: "mallet",
				},
				scopes:    []string{},
				expectErr: true,
			},
			{
				token: tokens[0][1],
				req: &ladon.Request{
					Subject: "mallet",
				},
				scopes:    []string{"core"},
				expectErr: true,
			},
			{
				token: tokens[0][1],
				req: &ladon.Request{
					Subject: "alice",
				},
				scopes:    []string{"foo"},
				expectErr: true,
			},
			{
				token: tokens[0][1],
				req: &ladon.Request{
					Subject:  "alice",
					Resource: "matrix",
					Action:   "create",
					Context:  ladon.Context{},
				},
				scopes:    []string{"foo"},
				expectErr: true,
			},
			{
				token: tokens[0][1],
				req: &ladon.Request{
					Subject:  "alice",
					Resource: "matrix",
					Action:   "delete",
					Context:  ladon.Context{},
				},
				scopes:    []string{"core"},
				expectErr: true,
			},
			{
				token: tokens[0][1],
				req: &ladon.Request{
					Subject:  "alice",
					Resource: "matrix",
					Action:   "create",
					Context:  ladon.Context{},
				},
				scopes:    []string{"core"},
				expectErr: false,
				assert: func(c *firewall.Context) {
					assert.Equal(t, "alice", c.Subject)
					assert.Equal(t, "tests", c.Issuer)
					assert.Equal(t, now.Add(time.Hour), c.ExpiresAt)
					assert.Equal(t, now, c.IssuedAt)
				},
			},
		} {
			ctx, err := w.TokenAllowed(context.Background(), c.token, c.req, c.scopes...)
			pkg.AssertError(t, c.expectErr, err, "ActionAllowed case", n, k)
			if err == nil && c.assert != nil {
				c.assert(ctx)
			}
		}
	}
}

func TestAuthorized(t *testing.T) {
	for n, w := range wardens {
		for k, c := range []struct {
			token     string
			scopes    []string
			expectErr bool
			assert    func(*firewall.Context)
		}{
			{
				token:     "invalid",
				expectErr: true,
			},
			{
				token:     "invalid",
				expectErr: true,
			},
			{
				token:     tokens[0][1],
				scopes:    []string{"foo"},
				expectErr: true,
			},
			{
				token:     tokens[1][1],
				scopes:    []string{"core"},
				expectErr: false,
				assert: func(c *firewall.Context) {
					assert.Equal(t, "siri", c.Subject)
					assert.Equal(t, "tests", c.Issuer)
					assert.Equal(t, now.Add(time.Hour), c.ExpiresAt, "expires at", n)
					assert.Equal(t, now, c.IssuedAt, "issued at", n)
				},
			},
			{
				token:     tokens[2][1],
				scopes:    []string{"core"},
				expectErr: true,
			},
		} {
			ctx, err := w.InspectToken(context.Background(), c.token, c.scopes...)
			pkg.AssertError(t, c.expectErr, err, "ActionAllowed case", n, k)
			if err == nil && c.assert != nil {
				c.assert(ctx)
			}
		}
	}
}
