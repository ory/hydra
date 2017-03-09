package integration_test

import (
	"testing"

	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/compose"
	"github.com/ory-am/fosite/handler/oauth2"
	"github.com/parnurzeal/gorequest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	goauth "golang.org/x/oauth2"
	"net/http"
)

func TestRevokeToken(t *testing.T) {
	for _, strategy := range []oauth2.AccessTokenStrategy{
		hmacStrategy,
	} {
		runRevokeTokenTest(t, strategy)
	}
}

func runRevokeTokenTest(t *testing.T, strategy oauth2.AccessTokenStrategy) {
	f := compose.Compose(new(compose.Config), fositeStore, strategy, compose.OAuth2ClientCredentialsGrantFactory, compose.OAuth2TokenIntrospectionFactory, compose.OAuth2TokenRevocationFactory)
	ts := mockServer(t, f, &fosite.DefaultSession{})
	defer ts.Close()

	oauthClient := newOAuth2AppClient(ts)
	token, err := oauthClient.Token(goauth.NoContext)
	assert.Nil(t, err)

	resp, _, errs := gorequest.New().Post(ts.URL+"/revoke").
		SetBasicAuth(oauthClient.ClientID, oauthClient.ClientSecret).
		Type("form").
		SendStruct(map[string]string{"token": "asdf"}).End()
	assert.Len(t, errs, 0)
	assert.Equal(t, 200, resp.StatusCode)

	resp, _, errs = gorequest.New().Post(ts.URL+"/revoke").
		SetBasicAuth(oauthClient.ClientID, oauthClient.ClientSecret).
		Type("form").
		SendStruct(map[string]string{"token": token.AccessToken}).End()
	assert.Len(t, errs, 0)
	assert.Equal(t, 200, resp.StatusCode)

	hres, _, errs := gorequest.New().Get(ts.URL+"/info").
		Set("Authorization", "bearer "+token.AccessToken).
		End()
	require.Len(t, errs, 0)
	assert.Equal(t, http.StatusUnauthorized, hres.StatusCode)
}
