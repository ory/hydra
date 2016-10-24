package oauth2_test

import (
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/fosite/compose"
	"github.com/ory-am/fosite/storage"
	"github.com/ory-am/hydra/herodot"
	"github.com/ory-am/hydra/oauth2"
	"github.com/ory-am/hydra/pkg"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	revocators           = make(map[string]oauth2.Revocator)
	nowRecovator         = time.Now().Round(time.Second)
	tokensRecovator      = pkg.Tokens(3)
	fositeStoreRecovator = storage.NewExampleStore()
)

func init() {

	r := httprouter.New()
	serv := &oauth2.Handler{
		OAuth2: compose.Compose(
			fc,
			fositeStoreRecovator,
			&compose.CommonStrategy{
				CoreStrategy:               compose.NewOAuth2HMACStrategy(fc, []byte("1234567890123456789012345678901234567890")),
				OpenIDConnectTokenStrategy: compose.NewOpenIDConnectStrategy(pkg.MustRSAKey()),
			},
			compose.OAuth2TokenIntrospectionFactory,
			compose.OAuth2TokenRevocationFactory,
		),
		H: &herodot.JSON{},
	}
	serv.SetRoutes(r)
	ts = httptest.NewServer(r)

	ar := fosite.NewAccessRequest(oauth2.NewSession("alice"))
	ar.GrantedScopes = fosite.Arguments{"core"}
	ar.RequestedAt = nowRecovator
	ar.Client = &fosite.DefaultClient{ID: "siri"}
	ar.Session.SetExpiresAt(fosite.AccessToken, nowRecovator.Add(time.Hour))
	ar.Session.(*oauth2.Session).Extra = map[string]interface{}{"foo": "bar"}
	fositeStoreRecovator.CreateAccessTokenSession(nil, tokensRecovator[0][0], ar)

	ar2 := fosite.NewAccessRequest(oauth2.NewSession("siri"))
	ar2.GrantedScopes = fosite.Arguments{"core"}
	ar2.RequestedAt = nowRecovator
	ar2.Session.(*oauth2.Session).Extra = map[string]interface{}{"foo": "bar"}
	ar2.Session.SetExpiresAt(fosite.AccessToken, nowRecovator.Add(time.Hour))
	ar2.Client = &fosite.DefaultClient{ID: "siri"}
	fositeStoreRecovator.CreateAccessTokenSession(nil, tokensRecovator[1][0], ar2)

	ar3 := fosite.NewAccessRequest(oauth2.NewSession("siri"))
	ar3.GrantedScopes = fosite.Arguments{"core"}
	ar3.RequestedAt = nowRecovator
	ar3.Session.(*oauth2.Session).Extra = map[string]interface{}{"foo": "bar"}
	ar3.Client = &fosite.DefaultClient{ID: "doesnt-exist"}
	ar3.Session.SetExpiresAt(fosite.AccessToken, nowRecovator.Add(-time.Hour))
	fositeStoreRecovator.CreateAccessTokenSession(nil, tokensRecovator[2][0], ar3)

	ep, err := url.Parse(ts.URL)
	if err != nil {
		logrus.Fatalf("%s", err)
	}
	revocators["http"] = &oauth2.HTTPRecovator{
		Endpoint: ep,
		Config: &clientcredentials.Config{
			ClientID:     "my-client",
			ClientSecret: "foobar",
		},
	}
}

func TestRevoke(t *testing.T) {
	for k, w := range revocators {
		for _, c := range []struct {
			token     string
			expectErr bool
		}{
			{
				token:     "invalid",
				expectErr: false,
			},
			{
				token:     tokensRecovator[0][1],
				expectErr: false,
			},
			{
				token:     tokensRecovator[0][1],
				expectErr: false,
			},
			{
				token:     tokensRecovator[2][1],
				expectErr: false,
			},
			{
				token:     tokensRecovator[1][1],
				expectErr: false,
			},
		} {
			t.Run(fmt.Sprintf("case=%s", k), func(t *testing.T) {
				err := w.RevokeToken(context.Background(), c.token)
				pkg.AssertError(t, c.expectErr, err)
			})
		}
	}
}
