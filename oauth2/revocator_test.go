package oauth2_test

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/storage"
	"github.com/ory/herodot"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createAccessTokenSession(subject, client string, token string, expiresAt time.Time, fs *storage.MemoryStore, scopes fosite.Arguments) {
	ar := fosite.NewAccessRequest(oauth2.NewSession(subject))
	ar.GrantedScopes = fosite.Arguments{"core"}
	if scopes != nil {
		ar.GrantedScopes = scopes
	}
	ar.RequestedAt = time.Now().Round(time.Minute)
	ar.Client = &fosite.DefaultClient{ID: client}
	ar.Session.SetExpiresAt(fosite.AccessToken, expiresAt)
	ar.Session.(*oauth2.Session).Extra = map[string]interface{}{"foo": "bar"}
	fs.CreateAccessTokenSession(nil, token, ar)
}

func TestRevoke(t *testing.T) {
	var (
		tokens = pkg.Tokens(3)
		store  = storage.NewExampleStore()
		now    = time.Now().Round(time.Second)
	)

	handler := &oauth2.Handler{
		OAuth2: compose.Compose(
			fc,
			store,
			&compose.CommonStrategy{
				CoreStrategy:               compose.NewOAuth2HMACStrategy(fc, []byte("1234567890123456789012345678901234567890")),
				OpenIDConnectTokenStrategy: compose.NewOpenIDConnectStrategy(pkg.MustRSAKey()),
			},
			nil,
			compose.OAuth2TokenIntrospectionFactory,
			compose.OAuth2TokenRevocationFactory,
		),
		H: herodot.NewJSONWriter(nil),
	}

	router := httprouter.New()
	handler.SetRoutes(router)
	server := httptest.NewServer(router)

	createAccessTokenSession("alice", "siri", tokens[0][0], now.Add(time.Hour), store, nil)
	createAccessTokenSession("siri", "siri", tokens[1][0], now.Add(time.Hour), store, nil)
	createAccessTokenSession("siri", "doesnt-exist", tokens[2][0], now.Add(-time.Hour), store, nil)

	client := hydra.NewOAuth2ApiWithBasePath(server.URL)
	client.Configuration.Username = "my-client"
	client.Configuration.Password = "foobar"

	for k, c := range []struct {
		token  string
		assert func(*testing.T)
	}{
		{
			token: "invalid",
		},
		{
			token: tokens[0][1],
			assert: func(t *testing.T) {
				assert.Len(t, store.AccessTokens, 2)
			},
		},
		{
			token: tokens[0][1],
		},
		{
			token: tokens[2][1],
			assert: func(t *testing.T) {
				assert.Len(t, store.AccessTokens, 1)
			},
		},
		{
			token: tokens[1][1],
			assert: func(t *testing.T) {
				assert.Len(t, store.AccessTokens, 0)
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			response, err := client.RevokeOAuth2Token(c.token)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, response.StatusCode)

			if c.assert != nil {
				c.assert(t)
			}
		})
	}
}
