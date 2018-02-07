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

package oauth2_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/storage"
	"github.com/ory/herodot"
	"github.com/ory/hydra/client"
	c2 "github.com/ory/hydra/compose"
	"github.com/ory/hydra/oauth2"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/ory/ladon"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var lifespan = time.Hour
var flushRequests = []*fosite.Request{
	{
		ID:            "flush-1",
		RequestedAt:   time.Now().Round(time.Second),
		Client:        &client.Client{ID: "foobar"},
		Scopes:        fosite.Arguments{"fa", "ba"},
		GrantedScopes: fosite.Arguments{"fa", "ba"},
		Form:          url.Values{"foo": []string{"bar", "baz"}},
		Session:       &fosite.DefaultSession{Subject: "bar"},
	},
	{
		ID:            "flush-2",
		RequestedAt:   time.Now().Round(time.Second).Add(-(lifespan + time.Minute)),
		Client:        &client.Client{ID: "foobar"},
		Scopes:        fosite.Arguments{"fa", "ba"},
		GrantedScopes: fosite.Arguments{"fa", "ba"},
		Form:          url.Values{"foo": []string{"bar", "baz"}},
		Session:       &fosite.DefaultSession{Subject: "bar"},
	},
	{
		ID:            "flush-3",
		RequestedAt:   time.Now().Round(time.Second).Add(-(lifespan + time.Hour)),
		Client:        &client.Client{ID: "foobar"},
		Scopes:        fosite.Arguments{"fa", "ba"},
		GrantedScopes: fosite.Arguments{"fa", "ba"},
		Form:          url.Values{"foo": []string{"bar", "baz"}},
		Session:       &fosite.DefaultSession{Subject: "bar"},
	},
}

func TestHandlerFlushHandler(t *testing.T) {
	localWarden, httpClient := c2.NewMockFirewall(
		"tests",
		"alice",
		fosite.Arguments{
			"hydra.oauth2.flush",
		}, &ladon.DefaultPolicy{
			ID:        "1",
			Subjects:  []string{"<.*>"},
			Resources: []string{"rn:hydra:oauth2:tokens"},
			Actions:   []string{"flush"},
			Effect:    ladon.AllowAccess,
		},
	)

	store := oauth2.NewFositeMemoryStore(nil, lifespan)
	h := &oauth2.Handler{
		H:             herodot.NewJSONWriter(nil),
		W:             localWarden,
		ScopeStrategy: fosite.HierarchicScopeStrategy,
		Issuer:        "http://hydra.localhost",
		Storage:       store,
	}

	for _, r := range flushRequests {
		require.NoError(t, store.CreateAccessTokenSession(nil, r.ID, r))
	}

	r := httprouter.New()
	h.SetRoutes(r)
	ts := httptest.NewServer(r)
	c := hydra.NewOAuth2ApiWithBasePath(ts.URL)
	c.Configuration.Transport = httpClient.Transport

	ds := new(fosite.DefaultSession)
	ctx := context.Background()

	resp, err := c.FlushInactiveOAuth2Tokens(hydra.FlushInactiveOAuth2TokensRequest{NotAfter: time.Now().Add(-time.Hour * 24)})
	require.NoError(t, err)
	assert.EqualValues(t, http.StatusNoContent, resp.StatusCode)

	_, err = store.GetAccessTokenSession(ctx, "flush-1", ds)
	require.NoError(t, err)
	_, err = store.GetAccessTokenSession(ctx, "flush-2", ds)
	require.NoError(t, err)
	_, err = store.GetAccessTokenSession(ctx, "flush-3", ds)
	require.NoError(t, err)

	resp, err = c.FlushInactiveOAuth2Tokens(hydra.FlushInactiveOAuth2TokensRequest{NotAfter: time.Now().Add(-(lifespan + time.Hour/2))})
	require.NoError(t, err)
	assert.EqualValues(t, http.StatusNoContent, resp.StatusCode)

	_, err = store.GetAccessTokenSession(ctx, "flush-1", ds)
	require.NoError(t, err)
	_, err = store.GetAccessTokenSession(ctx, "flush-2", ds)
	require.NoError(t, err)
	_, err = store.GetAccessTokenSession(ctx, "flush-3", ds)
	require.Error(t, err)

	resp, err = c.FlushInactiveOAuth2Tokens(hydra.FlushInactiveOAuth2TokensRequest{NotAfter: time.Now()})
	require.NoError(t, err)
	assert.EqualValues(t, http.StatusNoContent, resp.StatusCode)

	_, err = store.GetAccessTokenSession(ctx, "flush-1", ds)
	require.NoError(t, err)
	_, err = store.GetAccessTokenSession(ctx, "flush-2", ds)
	require.Error(t, err)
	_, err = store.GetAccessTokenSession(ctx, "flush-3", ds)
	require.Error(t, err)
}

func TestHandlerWellKnown(t *testing.T) {
	h := &oauth2.Handler{
		H:             herodot.NewJSONWriter(nil),
		ScopeStrategy: fosite.HierarchicScopeStrategy,
		Issuer:        "http://hydra.localhost",
	}

	AuthPathT := "/oauth2/auth"
	TokenPathT := "/oauth2/token"
	JWKPathT := "/.well-known/jwks.json"

	r := httprouter.New()
	h.SetRoutes(r)
	ts := httptest.NewServer(r)

	res, err := http.Get(ts.URL + "/.well-known/openid-configuration")
	require.NoError(t, err)
	defer res.Body.Close()

	trueConfig := oauth2.WellKnown{
		Issuer:                            h.Issuer,
		AuthURL:                           h.Issuer + AuthPathT,
		TokenURL:                          h.Issuer + TokenPathT,
		JWKsURI:                           h.Issuer + JWKPathT,
		SubjectTypes:                      []string{"pairwise", "public"},
		ResponseTypes:                     []string{"code", "code id_token", "id_token", "token id_token", "token", "token id_token code"},
		ClaimsSupported:                   []string{"sub"},
		ScopesSupported:                   []string{"offline", "openid"},
		UserinfoEndpoint:                  h.Issuer + oauth2.UserinfoPath,
		TokenEndpointAuthMethodsSupported: []string{"client_secret_post", "client_secret_basic"},
		IDTokenSigningAlgValuesSupported:  []string{"RS256"},
	}
	var wellKnownResp oauth2.WellKnown
	err = json.NewDecoder(res.Body).Decode(&wellKnownResp)
	require.NoError(t, err, "problem decoding wellknown json response: %+v", err)
	assert.EqualValues(t, trueConfig, wellKnownResp)

	h.ScopesSupported = "foo,bar"
	h.ClaimsSupported = "baz,oof"
	h.UserinfoEndpoint = "bar"

	res, err = http.Get(ts.URL + "/.well-known/openid-configuration")
	require.NoError(t, err)
	defer res.Body.Close()
	require.NoError(t, json.NewDecoder(res.Body).Decode(&wellKnownResp))

	assert.EqualValues(t, wellKnownResp.ClaimsSupported, []string{"sub", "baz", "oof"})
	assert.EqualValues(t, wellKnownResp.ScopesSupported, []string{"offline", "openid", "foo", "bar"})
	assert.Equal(t, wellKnownResp.UserinfoEndpoint, "bar")
}

type FakeConsentStrategy struct {
	RedirectURL string
}

func (s *FakeConsentStrategy) ValidateConsentRequest(authorizeRequest fosite.AuthorizeRequester, token string, session *sessions.Session) (claims *oauth2.Session, err error) {
	return oauth2.NewSession("consent-user"), nil
}

func (s *FakeConsentStrategy) CreateConsentRequest(authorizeRequest fosite.AuthorizeRequester, redirectURL string, session *sessions.Session) (token string, err error) {
	s.RedirectURL = redirectURL
	return "token", nil
}

func (s *FakeConsentStrategy) HandleConsentRequest(authorizeRequest fosite.AuthorizeRequester, session *sessions.Session) (claims *oauth2.Session, err error) {
	return nil, oauth2.ErrRequiresAuthentication
}

func TestIssuerRedirect(t *testing.T) {
	storage := storage.NewExampleStore()
	secret := []byte("my super secret password password password password")
	config := compose.Config{}
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	consentUrl, _ := url.Parse("http://consent.localhost")

	cs := &FakeConsentStrategy{}

	h := &oauth2.Handler{
		H:             herodot.NewJSONWriter(nil),
		Issuer:        "http://127.0.0.1/some/proxied/path",
		OAuth2:        compose.ComposeAllEnabled(&config, storage, secret, privateKey),
		ConsentURL:    *consentUrl,
		ScopeStrategy: fosite.WildcardScopeStrategy,
		CookieStore:   sessions.NewCookieStore([]byte("my super secret password")),
		Consent:       cs,
		L:             logrus.New(),
	}

	r := httprouter.New()
	h.SetRoutes(r)
	ts := httptest.NewServer(r)

	authUrl, _ := url.Parse(ts.URL)
	v := url.Values{}
	v.Set("response_type", "code")
	v.Set("client_id", "my-client")
	v.Set("redirect_uri", "http://localhost:3846/callback")
	v.Set("scope", "openid")
	v.Set("state", "my super secret state")
	authUrl.Path = "/oauth2/auth"
	authUrl.RawQuery = v.Encode()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, _ := client.Get(authUrl.String())

	authRedirect, _ := url.Parse(cs.RedirectURL)
	assert.Equal(t, "/some/proxied/path/oauth2/auth", authRedirect.Path, "The redirect URL sent in the challenge includes the full issuer path")
	assert.Equal(t, authUrl.Query(), authRedirect.Query(), "The auth redirect should have the same parameters with the addition of challenge")

	defer res.Body.Close()
}
