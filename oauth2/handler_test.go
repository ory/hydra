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

package oauth2

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"encoding/json"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/storage"
	"github.com/ory/herodot"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerWellKnown(t *testing.T) {
	h := &Handler{
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

	trueConfig := WellKnown{
		Issuer:                            h.Issuer,
		AuthURL:                           h.Issuer + AuthPathT,
		TokenURL:                          h.Issuer + TokenPathT,
		JWKsURI:                           h.Issuer + JWKPathT,
		SubjectTypes:                      []string{"pairwise", "public"},
		ResponseTypes:                     []string{"code", "code id_token", "id_token", "token id_token", "token", "token id_token code"},
		ClaimsSupported:                   []string{"sub"},
		ScopesSupported:                   []string{"offline", "openid"},
		UserinfoEndpoint:                  h.Issuer + UserinfoPath,
		TokenEndpointAuthMethodsSupported: []string{"client_secret_post", "client_secret_basic"},
		IDTokenSigningAlgValuesSupported:  []string{"RS256"},
	}
	var wellKnownResp WellKnown
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

func (s *FakeConsentStrategy) ValidateConsentRequest(authorizeRequest fosite.AuthorizeRequester, token string, session *sessions.Session) (claims *Session, err error) {
	return nil, nil
}

func (s *FakeConsentStrategy) CreateConsentRequest(authorizeRequest fosite.AuthorizeRequester, redirectURL string, session *sessions.Session) (token string, err error) {
	s.RedirectURL = redirectURL
	return "token", nil
}

func TestIssuerRedirect(t *testing.T) {
	storage := storage.NewExampleStore()
	secret := []byte("my super secret password password password password")
	config := compose.Config{}
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	consentUrl, _ := url.Parse("http://consent.localhost")

	cs := &FakeConsentStrategy{}

	h := &Handler{
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
