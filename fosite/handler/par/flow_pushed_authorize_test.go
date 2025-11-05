// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package par_test

import (
	"context"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/fosite/storage"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/handler/par"
)

func parseURL(uu string) *url.URL {
	u, _ := url.Parse(uu)
	return u
}

func TestAuthorizeCode_HandleAuthorizeEndpointRequest(t *testing.T) {
	requestURIPrefix := "urn:ietf:params:oauth:request_uri_diff:"
	store := storage.NewMemoryStore()
	handler := par.PushedAuthorizeHandler{
		Storage: store,
		Config: &fosite.Config{
			PushedAuthorizeContextLifespan:  30 * time.Minute,
			PushedAuthorizeRequestURIPrefix: requestURIPrefix,
			ScopeStrategy:                   fosite.HierarchicScopeStrategy,
			AudienceMatchingStrategy:        fosite.DefaultAudienceMatchingStrategy,
		},
	}
	for _, c := range []struct {
		handler     par.PushedAuthorizeHandler
		areq        *fosite.AuthorizeRequest
		description string
		expectErr   error
		expect      func(t *testing.T, areq *fosite.AuthorizeRequest, aresp *fosite.PushedAuthorizeResponse)
	}{
		{
			handler: handler,
			areq: &fosite.AuthorizeRequest{
				ResponseTypes: fosite.Arguments{""},
				Request:       *fosite.NewRequest(),
			},
			description: "should pass because not responsible for handling an empty response type",
		},
		{
			handler: handler,
			areq: &fosite.AuthorizeRequest{
				ResponseTypes: fosite.Arguments{"foo"},
				Request:       *fosite.NewRequest(),
			},
			description: "should pass because not responsible for handling an invalid response type",
		},
		{
			handler: handler,
			areq: &fosite.AuthorizeRequest{
				ResponseTypes: fosite.Arguments{"code"},
				Request: fosite.Request{
					Client: &fosite.DefaultClient{
						ResponseTypes: fosite.Arguments{"code"},
						RedirectURIs:  []string{"http://asdf.com/cb"},
					},
				},
				RedirectURI: parseURL("http://asdf.com/cb"),
			},
			description: "should fail because redirect uri is not https",
			expectErr:   fosite.ErrInvalidRequest,
		},
		{
			handler: handler,
			areq: &fosite.AuthorizeRequest{
				ResponseTypes: fosite.Arguments{"code"},
				Request: fosite.Request{
					Client: &fosite.DefaultClient{
						ResponseTypes: fosite.Arguments{"code"},
						RedirectURIs:  []string{"https://asdf.com/cb"},
						Audience:      []string{"https://www.ory.sh/api"},
					},
					RequestedAudience: []string{"https://www.ory.sh/not-api"},
				},
				RedirectURI: parseURL("https://asdf.com/cb"),
			},
			description: "should fail because audience doesn't match",
			expectErr:   fosite.ErrInvalidRequest,
		},
		{
			handler: handler,
			areq: &fosite.AuthorizeRequest{
				ResponseTypes: fosite.Arguments{"code"},
				Request: fosite.Request{
					Client: &fosite.DefaultClient{
						ResponseTypes: fosite.Arguments{"code"},
						RedirectURIs:  []string{"https://asdf.de/cb"},
						Audience:      []string{"https://www.ory.sh/api"},
					},
					RequestedAudience: []string{"https://www.ory.sh/api"},
					GrantedScope:      fosite.Arguments{"a", "b"},
					Session: &fosite.DefaultSession{
						ExpiresAt: map[fosite.TokenType]time.Time{fosite.AccessToken: time.Now().UTC().Add(time.Hour)},
					},
					RequestedAt: time.Now().UTC(),
				},
				State:       "superstate",
				RedirectURI: parseURL("https://asdf.de/cb"),
			},
			description: "should pass",
			expect: func(t *testing.T, areq *fosite.AuthorizeRequest, aresp *fosite.PushedAuthorizeResponse) {
				requestURI := aresp.RequestURI
				assert.NotEmpty(t, requestURI)
				assert.True(t, strings.HasPrefix(requestURI, requestURIPrefix), "requestURI does not match: %s", requestURI)
			},
		},
	} {
		t.Run("case="+c.description, func(t *testing.T) {
			aresp := &fosite.PushedAuthorizeResponse{}
			err := c.handler.HandlePushedAuthorizeEndpointRequest(context.Background(), c.areq, aresp)
			if c.expectErr != nil {
				require.EqualError(t, err, c.expectErr.Error())
			} else {
				require.NoError(t, err)
			}

			if c.expect != nil {
				c.expect(t, c.areq, aresp)
			}
		})
	}
}
