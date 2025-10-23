// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oauth2_test

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	goauth2 "golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	hc "github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/internal/testhelpers"
	"github.com/ory/x/configx"
	"github.com/ory/x/otelx"
)

func BenchmarkClientCredentials(b *testing.B) {
	ctx := context.Background()

	spans := tracetest.NewSpanRecorder()
	tracer := trace.NewTracerProvider(trace.WithSpanProcessor(spans)).Tracer("")

	dsn := "postgres://postgres:secret@127.0.0.1:3445/postgres?sslmode=disable"
	reg, err := testhelpers.NewRegistrySQLFromURL(b.Context(), dsn, true,
		driver.WithTracerWrapper(func(t *otelx.Tracer) *otelx.Tracer { return new(otelx.Tracer).WithOTLP(tracer) }),
		driver.WithConfigOptions(configx.WithValue(config.KeyAccessTokenStrategy, "opaque")))
	require.NoError(b, err)
	public, admin := testhelpers.NewOAuth2Server(ctx, b, reg)

	newCustomClient := func(b *testing.B, c *hc.Client) (*hc.Client, clientcredentials.Config) {
		unhashedSecret := c.Secret
		require.NoError(b, reg.ClientManager().CreateClient(ctx, c))
		return c, clientcredentials.Config{
			ClientID:       c.GetID(),
			ClientSecret:   unhashedSecret,
			TokenURL:       reg.Config().OAuth2TokenURL(ctx).String(),
			Scopes:         strings.Split(c.Scope, " "),
			EndpointParams: url.Values{"audience": c.Audience},
		}
	}

	newClient := func(b *testing.B) (*hc.Client, clientcredentials.Config) {
		return newCustomClient(b, &hc.Client{
			Secret:        uuid.Must(uuid.NewV4()).String(),
			RedirectURIs:  []string{public.URL + "/callback"},
			ResponseTypes: []string{"token"},
			GrantTypes:    []string{"client_credentials"},
			Scope:         "foobar",
			Audience:      []string{"https://api.ory.sh/"},
		})
	}

	getToken := func(t *testing.B, conf clientcredentials.Config) (*goauth2.Token, error) {
		conf.AuthStyle = goauth2.AuthStyleInHeader
		return conf.Token(context.Background())
	}

	encodeOr := func(b *testing.B, val interface{}, or string) string {
		out, err := json.Marshal(val)
		require.NoError(b, err)
		if string(out) == "null" {
			return or
		}

		return string(out)
	}

	inspectToken := func(b *testing.B, token *goauth2.Token, cl *hc.Client, conf clientcredentials.Config, strategy string, expectedExp time.Time, checkExtraClaims bool) {
		introspection := testhelpers.IntrospectToken(b, token.AccessToken, admin)

		check := func(res gjson.Result) {
			assert.EqualValues(b, cl.GetID(), res.Get("client_id").String(), "%s", res.Raw)
			assert.EqualValues(b, cl.GetID(), res.Get("sub").String(), "%s", res.Raw)
			assert.EqualValues(b, reg.Config().IssuerURL(ctx).String(), res.Get("iss").String(), "%s", res.Raw)

			assert.EqualValues(b, res.Get("nbf").Int(), res.Get("iat").Int(), "%s", res.Raw)
			assert.WithinDuration(b, expectedExp, time.Unix(res.Get("exp").Int(), 0), time.Second)

			assert.EqualValues(b, encodeOr(b, conf.EndpointParams["audience"], "[]"), res.Get("aud").Raw, "%s", res.Raw)

			if checkExtraClaims {
				require.True(b, res.Get("ext.hooked").Bool())
			}
		}

		check(introspection)
		assert.True(b, introspection.Get("active").Bool())
		assert.EqualValues(b, "access_token", introspection.Get("token_use").String())
		assert.EqualValues(b, "Bearer", introspection.Get("token_type").String())
		assert.EqualValues(b, strings.Join(conf.Scopes, " "), introspection.Get("scope").String(), "%s", introspection.Raw)

		if strategy != "jwt" {
			return
		}

		jwtClaims := gjson.ParseBytes(testhelpers.InsecureDecodeJWT(b, token.AccessToken))
		assert.NotEmpty(b, jwtClaims.Get("jti").String())
		assert.EqualValues(b, encodeOr(b, conf.Scopes, "[]"), jwtClaims.Get("scp").Raw, "%s", introspection.Raw)
		check(jwtClaims)
	}

	getAndInspectToken := func(b *testing.B, cl *hc.Client, conf clientcredentials.Config, strategy string, expectedExp time.Time, checkExtraClaims bool) {
		token, err := getToken(b, conf)
		require.NoError(b, err)
		inspectToken(b, token, cl, conf, strategy, expectedExp, checkExtraClaims)
	}

	run := func(strategy string) func(b *testing.B) {
		return func(t *testing.B) {
			reg.Config().MustSet(ctx, config.KeyAccessTokenStrategy, strategy)

			cl, conf := newClient(b)
			getAndInspectToken(b, cl, conf, strategy, time.Now().Add(reg.Config().GetAccessTokenLifespan(ctx)), false)
		}
	}

	b.Run("strategy=jwt", func(b *testing.B) {
		initialDBSpans := dbSpans(spans)
		for i := 0; i < b.N; i++ {
			run("jwt")(b)
		}
		b.ReportMetric(0, "ns/op")
		b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms/op")
		b.ReportMetric((float64(dbSpans(spans)-initialDBSpans))/float64(b.N), "queries/op")
	})

	b.Run("strategy=opaque", func(b *testing.B) {
		initialDBSpans := dbSpans(spans)
		for i := 0; i < b.N; i++ {
			run("opaque")(b)
		}
		b.ReportMetric(0, "ns/op")
		b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms/op")
		b.ReportMetric((float64(dbSpans(spans)-initialDBSpans))/float64(b.N), "queries/op")
	})
}

func dbSpans(spans *tracetest.SpanRecorder) (count int) {
	for _, s := range spans.Started() {
		if strings.HasPrefix(s.Name(), "sql-") {
			count++
		}
	}
	return
}
