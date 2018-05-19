/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package oauth2_test

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/herodot"
	hc "github.com/ory/hydra/client"
	. "github.com/ory/hydra/oauth2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2/clientcredentials"
)

func TestClientCredentials(t *testing.T) {
	router := httprouter.New()
	l := logrus.New()
	l.Level = logrus.DebugLevel
	store := NewFositeMemoryStore(hc.NewMemoryManager(hasher), time.Second)

	ts := httptest.NewServer(router)
	handler := &Handler{
		OAuth2: compose.Compose(
			fc,
			store,
			oauth2Strategy,
			nil,
			compose.OAuth2ClientCredentialsGrantFactory,
			compose.OAuth2TokenIntrospectionFactory,
		),
		//Consent:         consentStrategy,
		CookieStore:     sessions.NewCookieStore([]byte("foo-secret")),
		ForcedHTTP:      true,
		ScopeStrategy:   fosite.HierarchicScopeStrategy,
		IDTokenLifespan: time.Minute,
		H:               herodot.NewJSONWriter(l),
		L:               l,
		IssuerURL:       ts.URL,
	}

	handler.SetRoutes(router)

	require.NoError(t, store.CreateClient(&hc.Client{
		ID:            "app-client",
		Secret:        "secret",
		RedirectURIs:  []string{ts.URL + "/callback"},
		ResponseTypes: []string{"token"},
		GrantTypes:    []string{"client_credentials"},
		Scope:         "foobar",
	}))

	oauthClientConfig := &clientcredentials.Config{
		ClientID:     "app-client",
		ClientSecret: "secret",
		TokenURL:     ts.URL + "/oauth2/token",
		Scopes:       []string{"foobar"},
	}

	tok, err := oauthClientConfig.Token(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, tok.AccessToken)
}
