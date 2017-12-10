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
		Expiry:      time.Now().UTC().Add(time.Hour),
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
