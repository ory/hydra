package warden_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"log"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/ory/hydra/firewall"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/ory/hydra/warden"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	coauth2 "golang.org/x/oauth2"
)

func TestWardenSDK(t *testing.T) {
	router := httprouter.New()
	handler := &warden.WardenHandler{
		H:      herodot.NewJSONWriter(nil),
		Warden: wardens["local"],
	}
	handler.SetRoutes(router)
	server := httptest.NewServer(router)

	conf := &coauth2.Config{
		Scopes:   []string{},
		Endpoint: coauth2.Endpoint{},
	}

	client := hydra.NewWardenApiWithBasePath(server.URL)
	client.Configuration.Transport = conf.Client(coauth2.NoContext, &coauth2.Token{
		AccessToken: tokens[1][1],
		Expiry:      time.Now().Add(time.Hour),
		TokenType:   "bearer",
	}).Transport

	t.Run("DoesWardenAllowAccessRequest", func(t *testing.T) {
		for k, c := range accessRequestTestCases {
			t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
				result, response, err := client.DoesWardenAllowAccessRequest(hydra.WardenAccessRequest{
					Action:   c.req.Action,
					Resource: c.req.Resource,
					Subject:  c.req.Subject,
					Context:  c.req.Context,
				})

				require.NoError(t, err)
				require.Equal(t, http.StatusOK, response.StatusCode)
				assert.Equal(t, !c.expectErr, result.Allowed)
			})
		}
	})

	t.Run("DoesWardenAllowAccessRequest", func(t *testing.T) {
		for k, c := range accessRequestTokenTestCases {
			t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
				result, response, err := client.DoesWardenAllowTokenAccessRequest(hydra.WardenTokenAccessRequest{
					Action:   c.req.Action,
					Resource: c.req.Resource,
					Token:    c.token,
					Scopes:   c.scopes,
					Context:  c.req.Context,
				})

				require.NoError(t, err)
				require.Equal(t, http.StatusOK, response.StatusCode)
				assert.Equal(t, !c.expectErr, result.Allowed)

				if err == nil && c.assert != nil {
					c.assert(t, &firewall.Context{
						Subject:       result.Subject,
						GrantedScopes: result.GrantedScopes,
						Issuer:        result.Issuer,
						ClientID:      result.ClientId,
						Extra:         result.AccessTokenExtra,
						ExpiresAt:     mustParseTime(result.ExpiresAt),
						IssuedAt:      mustParseTime(result.IssuedAt),
					})
				}
			})
		}
	})
}

func mustParseTime(t string) time.Time {
	result, err := time.Parse(time.RFC3339Nano, t)
	if err != nil {
		log.Fatalf("Could not parse date time %s because %s", t, err)
		return time.Time{}
	}
	return result
}
