// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package integration_test

import (
	"net/http"
	"testing"

	"github.com/parnurzeal/gorequest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	goauth "golang.org/x/oauth2"

	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/fosite/compose"
	"github.com/ory/hydra/v2/fosite/handler/oauth2"
)

func TestRevokeToken(t *testing.T) {
	for _, strategy := range []oauth2.AccessTokenStrategy{
		hmacStrategy,
	} {
		runRevokeTokenTest(t, strategy)
	}
}

func runRevokeTokenTest(t *testing.T, strategy oauth2.AccessTokenStrategy) {
	f := compose.Compose(new(fosite.Config), fositeStore, strategy, compose.OAuth2ClientCredentialsGrantFactory, compose.OAuth2TokenIntrospectionFactory, compose.OAuth2TokenRevocationFactory)
	ts := mockServer(t, f, &fosite.DefaultSession{})
	defer ts.Close()

	oauthClient := newOAuth2AppClient(ts)
	token, err := oauthClient.Token(goauth.NoContext)
	require.NoError(t, err)

	resp, _, errs := gorequest.New().Post(ts.URL+"/revoke").
		SetBasicAuth(oauthClient.ClientID, oauthClient.ClientSecret).
		Type("form").
		SendStruct(map[string]string{"token": "asdf"}).End()
	require.Len(t, errs, 0)
	assert.Equal(t, 200, resp.StatusCode)

	resp, _, errs = gorequest.New().Post(ts.URL+"/revoke").
		SetBasicAuth(oauthClient.ClientID, oauthClient.ClientSecret).
		Type("form").
		SendStruct(map[string]string{"token": token.AccessToken}).End()
	require.Len(t, errs, 0)
	assert.Equal(t, 200, resp.StatusCode)

	hres, _, errs := gorequest.New().Get(ts.URL+"/info").
		Set("Authorization", "bearer "+token.AccessToken).
		End()
	require.Len(t, errs, 0)
	assert.Equal(t, http.StatusUnauthorized, hres.StatusCode)
}
