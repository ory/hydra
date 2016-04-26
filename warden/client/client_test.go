package client_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/handler/core"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/oauth2"
	"github.com/ory-am/hydra/pkg"
	"github.com/ory-am/hydra/server"
	"github.com/ory-am/hydra/warden"
	"github.com/ory-am/hydra/warden/client"
	"github.com/ory-am/ladon"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	coauth2 "golang.org/x/oauth2"
)

var ts *httptest.Server

var wardens = map[string]warden.Warden{}

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
		Actions: []string{
			"an:hydra:warden:allowed",
			"an:hydra:warden:authorized",
		},
		Effect: ladon.AllowAccess,
	},
})

var fositeStore = pkg.FositeStore()

var tokens = pkg.Tokens(2)

func init() {
	wardens["local"] = &client.LocalWarden{
		Warden: ladonWarden,
		TokenValidator: &core.CoreValidator{
			AccessTokenStrategy: pkg.HMACStrategy,
			AccessTokenStorage:  fositeStore,
		},
		Issuer: "tests",
	}

	r := httprouter.New()
	serv := &server.Warden{
		Ladon:  ladonWarden,
		H:      &herodot.JSON{},
		Warden: wardens["local"],
	}
	serv.SetRoutes(r)
	ts = httptest.NewServer(r)

	url, err := url.Parse(ts.URL + server.AllowedHandlerPath)
	if err != nil {
		log.Fatalf("%s", err)
	}

	ar := fosite.NewAccessRequest(&oauth2.Session{Subject: "alice"})
	ar.GrantedScopes = fosite.Arguments{"core"}
	fositeStore.CreateAccessTokenSession(nil, tokens[0][0], ar)

	ar = fosite.NewAccessRequest(&oauth2.Session{Subject: "siri"})
	ar.GrantedScopes = fosite.Arguments{"hydra.warden"}
	fositeStore.CreateAccessTokenSession(nil, tokens[1][0], ar)

	conf := &coauth2.Config{
		Scopes:   []string{},
		Endpoint: coauth2.Endpoint{},
	}
	wardens["http"] = &client.HTTPWarden{
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
			assert    func(*warden.Context)
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
				assert: func(c *warden.Context) {
					assert.Equal(t, "alice", c.Subject)
					assert.Equal(t, "tests", c.Issuer)

				},
			},
		} {
			ctx, err := w.ActionAllowed(context.Background(), c.token, c.req, c.scopes...)
			pkg.AssertError(t, c.expectErr, err, n, "ActionAllowed", k)
			if err == nil && c.assert != nil {
				c.assert(ctx)
			}

			httpreq := &http.Request{Header: http.Header{}}
			httpreq.Header.Set("Authorization", "bearer "+c.token)
			ctx, err = w.HTTPActionAllowed(context.Background(), httpreq, c.req, c.scopes...)
			pkg.AssertError(t, c.expectErr, err, n, "HTTPActionAllowed", k)
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
			assert    func(*warden.Context)
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
				token:     tokens[0][1],
				scopes:    []string{"core"},
				expectErr: false,
				assert: func(c *warden.Context) {
					assert.Equal(t, "alice", c.Subject)
					assert.Equal(t, "tests", c.Issuer)

				},
			},
		} {
			ctx, err := w.Authorized(context.Background(), c.token, c.scopes...)
			pkg.AssertError(t, c.expectErr, err, n, "ActionAllowed", k)
			if err == nil && c.assert != nil {
				c.assert(ctx)
			}

			httpreq := &http.Request{Header: http.Header{}}
			httpreq.Header.Set("Authorization", "bearer "+c.token)
			ctx, err = w.HTTPAuthorized(context.Background(), httpreq, c.scopes...)
			pkg.AssertError(t, c.expectErr, err, n, "HTTPActionAllowed", k)
			if err == nil && c.assert != nil {
				c.assert(ctx)
			}
		}
	}
}
