// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fositex

import (
	"context"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/fosite"
	"github.com/ory/hydra/v2/persistence"
	"github.com/ory/hydra/v2/spec"
	"github.com/ory/x/configx"
	"github.com/ory/x/contextx"
	"github.com/ory/x/httpx"
	"github.com/ory/x/logrusx"
)

type stubConfigDeps struct {
	conf *config.DefaultProvider
}

var _ configDependencies = (*stubConfigDeps)(nil)

func (s *stubConfigDeps) Config() *config.DefaultProvider  { return s.conf }
func (s *stubConfigDeps) Persister() persistence.Persister { return nil }
func (s *stubConfigDeps) HTTPClient(context.Context, ...httpx.ResilientOptions) *retryablehttp.Client {
	return nil
}
func (s *stubConfigDeps) ClientHasher() fosite.Hasher     { return nil }
func (s *stubConfigDeps) ExtraFositeFactories() []Factory { return nil }

func newTestConfig(t *testing.T, opts ...configx.OptionModifier) *config.DefaultProvider {
	t.Helper()

	defaults := []configx.OptionModifier{
		configx.SkipValidation(),
		configx.WithValues(map[string]any{
			config.KeyBCryptCost:                     4,
			config.KeySubjectIdentifierAlgorithmSalt: "00000000",
			config.KeyGetSystemSecret:                []string{"000000000000000000000000000000000000000000000000"},
			config.KeyGetCookieSecrets:               []string{"000000000000000000000000000000000000000000000000"},
			config.KeyLogLevel:                       "trace",
			config.KeyDevelopmentMode:                true,
			"serve.public.host":                      "localhost",
		}),
	}

	all := append(defaults, opts...)
	p, err := configx.New(t.Context(), spec.ConfigValidationSchema, all...)
	require.NoError(t, err)
	return config.NewCustom(logrusx.New("", ""), p, contextx.NewTestConfigProvider(spec.ConfigValidationSchema, all...))
}

func TestGetTokenURLs(t *testing.T) {
	ctx := context.Background()

	t.Run("case=custom domain issuer accepted as valid JWT bearer audience", func(t *testing.T) {
		conf := newTestConfig(t,
			configx.WithValue(config.KeyIssuerURL, "https://auth.example.com"),
			configx.WithValue(config.KeyPublicURL, "https://hydra.example"),
		)

		c := NewConfig(&stubConfigDeps{conf: conf})
		tokenURLs := c.GetTokenURLs(ctx)

		assert.Contains(t, tokenURLs, "https://auth.example.com/oauth2/token",
			"issuer-derived token URL must be accepted as a valid audience")
		assert.Contains(t, tokenURLs, "https://hydra.example/oauth2/token",
			"internal token URL must still be accepted")
	})

	t.Run("case=deduplicates when issuer equals public URL", func(t *testing.T) {
		conf := newTestConfig(t,
			configx.WithValue(config.KeyIssuerURL, "https://hydra.example"),
			configx.WithValue(config.KeyPublicURL, "https://hydra.example"),
		)

		c := NewConfig(&stubConfigDeps{conf: conf})
		urls := c.GetTokenURLs(ctx)

		require.Len(t, urls, 1)
		assert.Equal(t, "https://hydra.example/oauth2/token", urls[0])
	})

}
