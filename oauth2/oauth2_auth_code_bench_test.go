// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/gobuffalo/pop/v6"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"golang.org/x/oauth2"

	hydra "github.com/ory/hydra-client-go/v2"
	hc "github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/internal"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/x/contextx"
	"github.com/ory/x/pointerx"
)

func BenchmarkAuthCode(b *testing.B) {
	ctx := context.Background()

	spans := tracetest.NewSpanRecorder()
	tracer := trace.NewTracerProvider(trace.WithSpanProcessor(spans)).Tracer("")

	dsn := "postgres://postgres:secret@127.0.0.1:3445/postgres?sslmode=disable"
	reg := internal.NewRegistrySQLFromURL(b, dsn, true, new(contextx.Default)).WithTracer(tracer)
	reg.Config().MustSet(ctx, config.KeyLogLevel, "error")
	reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "opaque")
	reg.Config().MustSet(ctx, config.KeyRefreshTokenHookURL, "")
	_, adminTS := testhelpers.NewOAuth2Server(ctx, b, reg)

	newOAuth2Client := func(b *testing.B, cb string) (*hc.Client, *oauth2.Config) {
		secret := uuid.New()
		c := &hc.Client{
			Secret:        secret,
			RedirectURIs:  []string{cb},
			ResponseTypes: []string{"id_token", "code", "token"},
			GrantTypes:    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
			Scope:         "hydra offline openid",
			Audience:      []string{"https://api.ory.sh/"},
		}
		require.NoError(b, reg.ClientManager().CreateClient(context.TODO(), c))
		return c, &oauth2.Config{
			ClientID:     c.GetID(),
			ClientSecret: secret,
			Endpoint: oauth2.Endpoint{
				AuthURL:   reg.Config().OAuth2AuthURL(ctx).String(),
				TokenURL:  reg.Config().OAuth2TokenURL(ctx).String(),
				AuthStyle: oauth2.AuthStyleInHeader,
			},
			Scopes: strings.Split(c.Scope, " "),
		}
	}

	adminClient := hydra.NewAPIClient(hydra.NewConfiguration())
	adminClient.GetConfig().Servers = hydra.ServerConfigurations{{URL: adminTS.URL}}

	getAuthorizeCode := func(b *testing.B, conf *oauth2.Config, c *http.Client, params ...oauth2.AuthCodeOption) (string, *http.Response) {
		if c == nil {
			c = testhelpers.NewEmptyJarClient(b)
		}

		state := uuid.New()
		resp, err := c.Get(conf.AuthCodeURL(state, params...))
		require.NoError(b, err)
		defer resp.Body.Close()

		q := resp.Request.URL.Query()
		require.EqualValues(b, state, q.Get("state"))
		return q.Get("code"), resp
	}

	acceptLoginHandler := func(b *testing.B, c *hc.Client, subject string, checkRequestPayload func(request *hydra.OAuth2LoginRequest) *hydra.AcceptOAuth2LoginRequest) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			rr, _, err := adminClient.OAuth2Api.GetOAuth2LoginRequest(context.Background()).LoginChallenge(r.URL.Query().Get("login_challenge")).Execute()
			require.NoError(b, err)

			assert.EqualValues(b, c.GetID(), pointerx.Deref(rr.Client.ClientId))
			assert.Empty(b, pointerx.Deref(rr.Client.ClientSecret))
			assert.EqualValues(b, c.GrantTypes, rr.Client.GrantTypes)
			assert.EqualValues(b, c.LogoURI, pointerx.Deref(rr.Client.LogoUri))
			assert.EqualValues(b, c.RedirectURIs, rr.Client.RedirectUris)
			assert.EqualValues(b, r.URL.Query().Get("login_challenge"), rr.Challenge)
			assert.EqualValues(b, []string{"hydra", "offline", "openid"}, rr.RequestedScope)
			assert.Contains(b, rr.RequestUrl, reg.Config().OAuth2AuthURL(ctx).String())

			acceptBody := hydra.AcceptOAuth2LoginRequest{
				Subject:  subject,
				Remember: pointerx.Ptr(!rr.Skip),
				Acr:      pointerx.Ptr("1"),
				Amr:      []string{"pwd"},
				Context:  map[string]interface{}{"context": "bar"},
			}
			if checkRequestPayload != nil {
				if b := checkRequestPayload(rr); b != nil {
					acceptBody = *b
				}
			}

			v, _, err := adminClient.OAuth2Api.AcceptOAuth2LoginRequest(context.Background()).
				LoginChallenge(r.URL.Query().Get("login_challenge")).
				AcceptOAuth2LoginRequest(acceptBody).
				Execute()
			require.NoError(b, err)
			require.NotEmpty(b, v.RedirectTo)
			http.Redirect(w, r, v.RedirectTo, http.StatusFound)
		}
	}

	acceptConsentHandler := func(b *testing.B, c *hc.Client, subject string, checkRequestPayload func(*hydra.OAuth2ConsentRequest)) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			rr, _, err := adminClient.OAuth2Api.GetOAuth2ConsentRequest(context.Background()).ConsentChallenge(r.URL.Query().Get("consent_challenge")).Execute()
			require.NoError(b, err)

			assert.EqualValues(b, c.GetID(), pointerx.StringR(rr.Client.ClientId))
			assert.Empty(b, pointerx.StringR(rr.Client.ClientSecret))
			assert.EqualValues(b, c.GrantTypes, rr.Client.GrantTypes)
			assert.EqualValues(b, c.LogoURI, pointerx.StringR(rr.Client.LogoUri))
			assert.EqualValues(b, c.RedirectURIs, rr.Client.RedirectUris)
			assert.EqualValues(b, subject, pointerx.StringR(rr.Subject))
			assert.EqualValues(b, []string{"hydra", "offline", "openid"}, rr.RequestedScope)
			assert.EqualValues(b, r.URL.Query().Get("consent_challenge"), rr.Challenge)
			assert.Contains(b, *rr.RequestUrl, reg.Config().OAuth2AuthURL(ctx).String())
			if checkRequestPayload != nil {
				checkRequestPayload(rr)
			}

			assert.Equal(b, map[string]interface{}{"context": "bar"}, rr.Context)
			v, _, err := adminClient.OAuth2Api.AcceptOAuth2ConsentRequest(context.Background()).
				ConsentChallenge(r.URL.Query().Get("consent_challenge")).
				AcceptOAuth2ConsentRequest(hydra.AcceptOAuth2ConsentRequest{
					GrantScope: []string{"hydra", "offline", "openid"}, Remember: pointerx.Bool(true), RememberFor: pointerx.Int64(0),
					GrantAccessTokenAudience: rr.RequestedAccessTokenAudience,
					Session: &hydra.AcceptOAuth2ConsentRequestSession{
						AccessToken: map[string]interface{}{"foo": "bar"},
						IdToken:     map[string]interface{}{"bar": "baz"},
					},
				}).
				Execute()
			require.NoError(b, err)
			require.NotEmpty(b, v.RedirectTo)
			http.Redirect(w, r, v.RedirectTo, http.StatusFound)
		}
	}

	subject := "aeneas-rekkas"
	nonce := uuid.New()
	c, conf := newOAuth2Client(b, testhelpers.NewCallbackURL(b, "callback", testhelpers.HTTPServerNotImplementedHandler))

	testhelpers.NewLoginConsentUI(b, reg.Config(),
		acceptLoginHandler(b, c, subject, nil),
		acceptConsentHandler(b, c, subject, nil),
	)

	run := func(b *testing.B, strategy string, exchange bool) {
		//pop.Debug = true
		code, _ := getAuthorizeCode(b, conf, nil, oauth2.SetAuthURLParam("nonce", nonce))
		//pop.Debug = false
		require.NotEmpty(b, code)

		if exchange {
			pop.Debug = true
			_, err := conf.Exchange(context.Background(), code)
			pop.Debug = false
			require.NoError(b, err)
		}
	}

	b.Run("strategy=jwt", func(b *testing.B) {
		reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "jwt")
		b.Run("case=with token exchange", func(b *testing.B) {

			initialDBSpans := dbSpans(spans)
			for i := 0; i < b.N; i++ {
				run(b, "jwt", true)
			}
			b.ReportMetric(0, "ns/op")
			b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms/op")
			b.ReportMetric((float64(dbSpans(spans)-initialDBSpans))/float64(b.N), "queries/op")
		})

		b.Run("case=without exchange", func(b *testing.B) {
			reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "jwt")
			initialDBSpans := dbSpans(spans)
			for i := 0; i < b.N; i++ {
				run(b, "jwt", false)
			}
			b.ReportMetric(0, "ns/op")
			b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms/op")
			b.ReportMetric((float64(dbSpans(spans)-initialDBSpans))/float64(b.N), "queries/op")
		})
	})

	//b.Run("strategy=opaque", func(b *testing.B) {
	//	reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "opaque")
	//	initialDBSpans := dbSpans(spans)
	//	for i := 0; i < b.N; i++ {
	//		run(b, "opaque")
	//	}
	//	b.ReportMetric(0, "ns/op")
	//	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms/op")
	//	b.ReportMetric((float64(dbSpans(spans)-initialDBSpans))/float64(b.N), "queries/op")
	//})
}
