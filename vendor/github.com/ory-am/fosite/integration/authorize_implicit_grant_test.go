package integration_test

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/compose"
	"github.com/ory-am/fosite/handler/oauth2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	goauth "golang.org/x/oauth2"
)

func TestAuthorizeImplicitFlow(t *testing.T) {
	for _, strategy := range []oauth2.AccessTokenStrategy{
		hmacStrategy,
	} {
		runTestAuthorizeImplicitGrant(t, strategy)
	}
}

func runTestAuthorizeImplicitGrant(t *testing.T, strategy interface{}) {
	f := compose.Compose(new(compose.Config), fositeStore, strategy, compose.OAuth2AuthorizeImplicitFactory, compose.OAuth2TokenIntrospectionFactory)
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
				state = "12345678901234567890"
			},
			authStatusCode: http.StatusOK,
		},
	} {
		c.setup()

		var callbackURL *url.URL
		authURL := strings.Replace(oauthClient.AuthCodeURL(state), "response_type=code", "response_type=token", -1)
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				callbackURL = req.URL
				return errors.New("Dont follow redirects")
			},
		}
		resp, err := client.Get(authURL)
		require.NotNil(t, err)

		if resp.StatusCode == http.StatusOK {
			fragment, err := url.ParseQuery(callbackURL.Fragment)
			require.Nil(t, err)
			expires, err := strconv.Atoi(fragment.Get("expires_in"))
			require.Nil(t, err)
			token := &goauth.Token{
				AccessToken:  fragment.Get("access_token"),
				TokenType:    fragment.Get("token_type"),
				RefreshToken: fragment.Get("refresh_token"),
				Expiry:       time.Now().Add(time.Duration(expires) * time.Second),
			}

			httpClient := oauthClient.Client(goauth.NoContext, token)
			resp, err := httpClient.Get(ts.URL + "/info")
			require.Nil(t, err, "(%d) %s", k, c.description)
			assert.Equal(t, http.StatusNoContent, resp.StatusCode, "(%d) %s", k, c.description)
		}
		t.Logf("Passed test case (%d) %s", k, c.description)
	}
}
