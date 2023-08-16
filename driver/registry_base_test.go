// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ory/x/randx"

	"github.com/stretchr/testify/require"

	"github.com/ory/x/httpx"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/x/configx"
	"github.com/ory/x/contextx"
	"github.com/ory/x/logrusx"

	"github.com/gorilla/sessions"
)

func TestGetJWKSFetcherStrategyHostEnforcment(t *testing.T) {
	ctx := context.Background()
	l := logrusx.New("", "")
	c := config.MustNew(context.Background(), l, configx.WithConfigFiles("../internal/.hydra.yaml"))
	c.MustSet(ctx, config.KeyDSN, "memory")
	c.MustSet(ctx, config.HSMEnabled, "false")
	c.MustSet(ctx, config.KeyClientHTTPNoPrivateIPRanges, true)

	registry, err := NewRegistryWithoutInit(c, l)
	require.NoError(t, err)

	_, err = registry.GetJWKSFetcherStrategy().Resolve(ctx, "http://localhost:8080", true)
	require.ErrorAs(t, err, new(httpx.ErrPrivateIPAddressDisallowed))
}

func TestRegistryBase_newKeyStrategy_handlesNetworkError(t *testing.T) {
	// Test ensures any network specific error is logged with a
	// specific message when attempting to create a new key strategy: issue #2338

	hook := test.Hook{} // Test hook for asserting log messages
	ctx := context.Background()

	l := logrusx.New("", "", logrusx.WithHook(&hook))
	l.Logrus().SetOutput(io.Discard)
	l.Logrus().ExitFunc = func(int) {} // Override the exit func to avoid call to os.Exit

	// Create a config and set a valid but unresolvable DSN
	c := config.MustNew(context.Background(), l, configx.WithConfigFiles("../internal/.hydra.yaml"))
	c.MustSet(ctx, config.KeyDSN, "postgres://user:password@127.0.0.1:9999/postgres")
	c.MustSet(ctx, config.HSMEnabled, "false")

	registry, err := NewRegistryWithoutInit(c, l)
	if err != nil {
		t.Errorf("Failed to create registry: %s", err)
		return
	}

	r := registry.(*RegistrySQL)
	r.initialPing = failedPing(errors.New("snizzles"))

	_ = r.Init(context.Background(), true, false, &contextx.TestContextualizer{}, nil, nil)

	registryBase := RegistryBase{r: r, l: l}
	registryBase.WithConfig(c)

	assert.Equal(t, logrus.FatalLevel, hook.LastEntry().Level)
	assert.Contains(t, hook.LastEntry().Message, "snizzles")
}

func TestRegistryBase_CookieStore_MaxAgeZero(t *testing.T) {
	// Test ensures that CookieStore MaxAge option is equal to zero after initialization

	ctx := context.Background()
	r := new(RegistryBase)
	r.WithConfig(config.MustNew(context.Background(), logrusx.New("", ""), configx.WithValue(config.KeyGetSystemSecret, []string{randx.MustString(32, randx.AlphaNum)})))

	s, err := r.CookieStore(ctx)
	require.NoError(t, err)
	cs := s.(*sessions.CookieStore)

	assert.Equal(t, cs.Options.MaxAge, 0)
}

func TestRegistryBase_HTTPClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	t.Setenv("CLIENTS_HTTP_PRIVATE_IP_EXCEPTION_URLS", fmt.Sprintf("[%q]", ts.URL+"/exception/*"))

	ctx := context.Background()
	r := new(RegistryBase)
	r.WithConfig(config.MustNew(
		ctx,
		logrusx.New("", ""),
		configx.WithValues(map[string]interface{}{
			config.KeyClientHTTPNoPrivateIPRanges: true,
		}),
	))

	t.Run("case=matches exception glob", func(t *testing.T) {
		res, err := r.HTTPClient(ctx).Get(ts.URL + "/exception/foo")
		require.NoError(t, err)
		assert.Equal(t, 200, res.StatusCode)
	})

	t.Run("case=does not match exception glob", func(t *testing.T) {
		_, err := r.HTTPClient(ctx).Get(ts.URL + "/foo")
		require.Error(t, err)
	})
}
