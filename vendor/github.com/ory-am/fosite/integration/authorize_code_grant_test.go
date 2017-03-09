package integration_test

import (
	"testing"

	"net/http"

	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/compose"
	"github.com/ory-am/fosite/handler/oauth2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	goauth "golang.org/x/oauth2"
)

func TestAuthorizeCodeFlow(t *testing.T) {
	for _, strategy := range []oauth2.AccessTokenStrategy{
		hmacStrategy,
	} {
		runAuthorizeCodeGrantTest(t, strategy)
	}
}

func runAuthorizeCodeGrantTest(t *testing.T, strategy interface{}) {
	f := compose.Compose(new(compose.Config), fositeStore, strategy, compose.OAuth2AuthorizeExplicitFactory, compose.OAuth2TokenIntrospectionFactory)
	ts := mockServer(t, f, &fosite.DefaultSession{})
	defer ts.Close()

	oauthClient := newOAuth2Client(ts)
	fositeStore.Clients["my-client"].RedirectURIs[0] = ts.URL + "/callback"

	var state string
	for k, c := range []struct {
		description    string
		setup          func()
		authStatusCode int
	}{
		{
			description: "should pass",
			setup: func() {
				oauthClient = newOAuth2Client(ts)
				state = "12345678901234567890"
			},
			authStatusCode: http.StatusOK,
		},
	} {
		c.setup()

		resp, err := http.Get(oauthClient.AuthCodeURL(state))
		require.Nil(t, err)
		require.Equal(t, c.authStatusCode, resp.StatusCode, "(%d) %s", k, c.description)

		if resp.StatusCode == http.StatusOK {
			token, err := oauthClient.Exchange(goauth.NoContext, resp.Request.URL.Query().Get("code"))
			require.Nil(t, err, "(%d) %s", k, c.description)
			require.NotEmpty(t, token.AccessToken, "(%d) %s", k, c.description)

			httpClient := oauthClient.Client(goauth.NoContext, token)
			resp, err := httpClient.Get(ts.URL + "/info")
			require.Nil(t, err, "(%d) %s", k, c.description)
			assert.Equal(t, http.StatusNoContent, resp.StatusCode, "(%d) %s", k, c.description)
		}
		t.Logf("Passed test case (%d) %s", k, c.description)
	}
}
