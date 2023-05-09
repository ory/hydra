// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"flag"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"testing"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"golang.org/x/oauth2"
	"gopkg.in/square/go-jose.v2"

	hydra "github.com/ory/hydra-client-go/v2"
	hc "github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/internal"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/contextx"
	"github.com/ory/x/pointerx"
)

var (
	prof = flag.String("profile", "", "write a CPU profile to this filename")
	conc = flag.Int("conc", 100, "dispatch this many requests concurrently")
)

func BenchmarkAuthCode(b *testing.B) {
	flag.Parse()

	ctx := context.Background()

	spans := tracetest.NewSpanRecorder()
	tracer := trace.NewTracerProvider(trace.WithSpanProcessor(spans)).Tracer("")

	dsn := "postgres://postgres:secret@127.0.0.1:3445/postgres?sslmode=disable&max_conns=10&max_idle_conns=10"
	// dsn := "mysql://root:secret@tcp(localhost:3444)/mysql?max_conns=16&max_idle_conns=16"
	// dsn := "cockroach://root@localhost:3446/defaultdb?sslmode=disable&max_conns=16&max_idle_conns=16"
	reg := internal.NewRegistrySQLFromURL(b, dsn, true, new(contextx.Default)).WithTracer(tracer)
	reg.Config().MustSet(ctx, config.KeyLogLevel, "error")
	reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, "opaque")
	reg.Config().MustSet(ctx, config.KeyRefreshTokenHookURL, "")
	oauth2Keys, err := jwk.GenerateJWK(ctx, jose.ES256, x.OAuth2JWTKeyName, "sig")
	require.NoError(b, err)
	oidcKeys, err := jwk.GenerateJWK(ctx, jose.ES256, x.OpenIDConnectKeyName, "sig")
	require.NoError(b, err)
	require.NoError(b, reg.KeyManager().UpdateKeySet(ctx, x.OAuth2JWTKeyName, oauth2Keys))
	require.NoError(b, reg.KeyManager().UpdateKeySet(ctx, x.OpenIDConnectKeyName, oidcKeys))
	_, adminTS := testhelpers.NewOAuth2Server(ctx, b, reg)
	var (
		authURL  = reg.Config().OAuth2AuthURL(ctx).String()
		tokenURL = reg.Config().OAuth2TokenURL(ctx).String()
		subject  = "aeneas-rekkas"
		nonce    = uuid.New()
	)

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
				AuthURL:   authURL,
				TokenURL:  tokenURL,
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
			assert.Contains(b, rr.RequestUrl, authURL)

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

			assert.EqualValues(b, c.GetID(), pointerx.Deref(rr.Client.ClientId))
			assert.Empty(b, pointerx.Deref(rr.Client.ClientSecret))
			assert.EqualValues(b, c.GrantTypes, rr.Client.GrantTypes)
			assert.EqualValues(b, c.LogoURI, pointerx.Deref(rr.Client.LogoUri))
			assert.EqualValues(b, c.RedirectURIs, rr.Client.RedirectUris)
			assert.EqualValues(b, subject, pointerx.Deref(rr.Subject))
			assert.EqualValues(b, []string{"hydra", "offline", "openid"}, rr.RequestedScope)
			assert.EqualValues(b, r.URL.Query().Get("consent_challenge"), rr.Challenge)
			assert.Contains(b, *rr.RequestUrl, authURL)
			if checkRequestPayload != nil {
				checkRequestPayload(rr)
			}

			assert.Equal(b, map[string]interface{}{"context": "bar"}, rr.Context)
			v, _, err := adminClient.OAuth2Api.AcceptOAuth2ConsentRequest(context.Background()).
				ConsentChallenge(r.URL.Query().Get("consent_challenge")).
				AcceptOAuth2ConsentRequest(hydra.AcceptOAuth2ConsentRequest{
					GrantScope: []string{"hydra", "offline", "openid"}, Remember: pointerx.Ptr(true), RememberFor: pointerx.Ptr[int64](0),
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

	run := func(b *testing.B, strategy string) func(*testing.B) {
		reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, strategy)
		c, conf := newOAuth2Client(b, testhelpers.NewCallbackURL(b, "callback", testhelpers.HTTPServerNotImplementedHandler))
		testhelpers.NewLoginConsentUI(b, reg.Config(),
			acceptLoginHandler(b, c, subject, nil),
			acceptConsentHandler(b, c, subject, nil),
		)

		return func(b *testing.B) {
			code, _ := getAuthorizeCode(b, conf, nil, oauth2.SetAuthURLParam("nonce", nonce))
			require.NotEmpty(b, code)

			//pop.Debug = true
			_, err := conf.Exchange(context.Background(), code)
			//pop.Debug = false
			require.NoError(b, err)
		}
	}

	b.SetParallelism(*conc / runtime.GOMAXPROCS(0))

	b.Run("strategy=jwt", func(b *testing.B) {
		initialDBSpans := dbSpans(spans)
		B := run(b, "jwt")

		stop := profile(b)
		defer stop()

		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				B(b)
			}
		})

		b.ReportMetric(0, "ns/op")
		b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms/op")
		b.ReportMetric((float64(dbSpans(spans)-initialDBSpans))/float64(b.N), "queries/op")
		b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "ops/s")
	})

	b.Run("strategy=opaque", func(b *testing.B) {
		initialDBSpans := dbSpans(spans)
		B := run(b, "opaque")

		stop := profile(b)
		defer stop()

		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				B(b)
			}
		})

		b.ReportMetric(0, "ns/op")
		b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms/op")
		b.ReportMetric((float64(dbSpans(spans)-initialDBSpans))/float64(b.N), "queries/op")
		b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "ops/s")
	})

}

func profile(t testing.TB) (stop func()) {
	if *prof == "" {
		return func() {} // noop
	}
	f, err := os.Create(*prof)
	require.NoError(t, err)
	require.NoError(t, pprof.StartCPUProfile(f))
	return func() {
		pprof.StopCPUProfile()
		require.NoError(t, f.Close())
		t.Log("Wrote profile to", f.Name())
	}
}
