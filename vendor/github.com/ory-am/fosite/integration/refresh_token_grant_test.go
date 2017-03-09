package integration_test

import (
	"testing"

	"net/http"
	"time"

	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/compose"
	hst "github.com/ory-am/fosite/handler/oauth2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestRefreshTokenFlow(t *testing.T) {
	for _, strategy := range []hst.AccessTokenStrategy{
		hmacStrategy,
	} {
		runRefreshTokenGrantTest(t, strategy)
	}
}

func runRefreshTokenGrantTest(t *testing.T, strategy interface{}) {
	f := compose.Compose(
		new(compose.Config),
		fositeStore,
		strategy,
		compose.OAuth2AuthorizeExplicitFactory,
		compose.OAuth2RefreshTokenGrantFactory,
		compose.OAuth2TokenIntrospectionFactory,
	)
	ts := mockServer(t, f, &fosite.DefaultSession{})
	defer ts.Close()

	oauthClient := newOAuth2Client(ts)
	state := "1234567890"
	fositeStore.Clients["my-client"].RedirectURIs[0] = ts.URL + "/callback"
	for k, c := range []struct {
		description string
		setup       func()
		pass        bool
	}{
		{
			description: "should fail because scope missing",
			setup:       func() {},
			pass:        false,
		},
		{
			description: "should pass",
			setup: func() {
				oauthClient.Scopes = []string{"fosite", "offline"}
			},
			pass: true,
		},
	} {
		c.setup()

		resp, err := http.Get(oauthClient.AuthCodeURL(state))
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode, "(%d) %s", k, c.description)

		if resp.StatusCode == http.StatusOK {
			token, err := oauthClient.Exchange(oauth2.NoContext, resp.Request.URL.Query().Get("code"))
			require.Nil(t, err, "(%d) %s", k, c.description)
			require.NotEmpty(t, token.AccessToken, "(%d) %s", k, c.description)

			t.Logf("Token %s\n", token)
			token.Expiry = token.Expiry.Add(-time.Hour * 24)
			t.Logf("Token %s\n", token)

			tokenSource := oauthClient.TokenSource(oauth2.NoContext, token)
			refreshed, err := tokenSource.Token()
			if c.pass {
				require.Nil(t, err, "(%d) %s: %s", k, c.description, err)
				assert.NotEqual(t, token.RefreshToken, refreshed.RefreshToken, "(%d) %s", k, c.description)
				assert.NotEqual(t, token.AccessToken, refreshed.AccessToken, "(%d) %s", k, c.description)
			} else {
				require.NotNil(t, err, "(%d) %s: %s", k, c.description, err)

			}
		}
		t.Logf("Passed test case (%d) %s", k, c.description)
	}
}
