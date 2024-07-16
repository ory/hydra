// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"context"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/persistence/sql"
	"github.com/ory/x/configx"
	"github.com/ory/x/contextx"
	"github.com/ory/x/errorsx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/sqlcon/dockertest"

	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/ory/x/randx"

	"github.com/ory/x/httpx"

	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
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

func TestRegistrySQL_newKeyStrategy_handlesNetworkError(t *testing.T) {
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

	assert.Equal(t, logrus.FatalLevel, hook.LastEntry().Level)
	assert.Contains(t, hook.LastEntry().Message, "snizzles")
}

func TestRegistrySQL_CookieStore_MaxAgeZero(t *testing.T) {
	// Test ensures that CookieStore MaxAge option is equal to zero after initialization

	ctx := context.Background()
	r := new(RegistrySQL)
	r.WithConfig(config.MustNew(context.Background(), logrusx.New("", ""), configx.WithValue(config.KeyGetSystemSecret, []string{randx.MustString(32, randx.AlphaNum)})))

	s, err := r.CookieStore(ctx)
	require.NoError(t, err)
	cs := s.(*sessions.CookieStore)

	assert.Equal(t, cs.Options.MaxAge, 0)
}

func TestRegistrySQL_HTTPClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	t.Setenv("CLIENTS_HTTP_PRIVATE_IP_EXCEPTION_URLS", fmt.Sprintf("[%q]", ts.URL+"/exception/*"))

	ctx := context.Background()
	r := new(RegistrySQL)
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

func TestDefaultKeyManager_HsmDisabled(t *testing.T) {
	l := logrusx.New("", "")
	c := config.MustNew(context.Background(), l, configx.SkipValidation())
	c.MustSet(context.Background(), config.KeyDSN, "postgres://user:password@127.0.0.1:9999/postgres")
	c.MustSet(context.Background(), config.HSMEnabled, "false")
	reg, err := NewRegistryWithoutInit(c, l)
	r := reg.(*RegistrySQL)
	r.initialPing = sussessfulPing()
	if err := r.Init(context.Background(), true, false, &contextx.Default{}, nil, nil); err != nil {
		t.Fatalf("unable to init registry: %s", err)
	}
	assert.NoError(t, err)
	assert.IsType(t, &sql.Persister{}, reg.KeyManager())
	assert.IsType(t, &sql.Persister{}, reg.SoftwareKeyManager())
}

func TestDbUnknownTableColumns(t *testing.T) {
	ctx := context.Background()
	l := logrusx.New("", "")
	c := config.MustNew(ctx, l, configx.SkipValidation())
	postgresDsn := dockertest.RunTestPostgreSQL(t)
	c.MustSet(ctx, config.KeyDSN, postgresDsn)
	reg, err := NewRegistryFromDSN(ctx, c, l, false, true, &contextx.Default{})
	require.NoError(t, err)

	statement := "ALTER TABLE \"hydra_client\" ADD COLUMN \"temp_column\" VARCHAR(128) NOT NULL DEFAULT '';"
	require.NoError(t, reg.Persister().Connection(ctx).RawQuery(statement).Exec())

	cl := &client.Client{
		ID: strconv.Itoa(rand.Int()),
	}
	require.NoError(t, reg.Persister().CreateClient(ctx, cl))
	getClients := func(reg Registry) ([]client.Client, error) {
		readClients := make([]client.Client, 0)
		return readClients, reg.Persister().Connection(ctx).RawQuery("SELECT * FROM \"hydra_client\"").All(&readClients)
	}

	t.Run("with ignore disabled (default behavior)", func(t *testing.T) {
		_, err := getClients(reg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing destination name temp_column")
	})

	t.Run("with ignore enabled", func(t *testing.T) {
		c.MustSet(ctx, config.KeyDBIgnoreUnknownTableColumns, true)
		reg, err := NewRegistryFromDSN(ctx, c, l, false, true, &contextx.Default{})
		require.NoError(t, err)

		actual, err := getClients(reg)
		require.NoError(t, err)
		assert.Len(t, actual, 1)
	})
}

func sussessfulPing() func(r *RegistrySQL) error {
	return func(r *RegistrySQL) error {
		// fake that ping is successful
		return nil
	}
}

func failedPing(err error) func(r *RegistrySQL) error {
	return func(r *RegistrySQL) error {
		r.Logger().Fatalf(err.Error())
		return errorsx.WithStack(err)
	}
}
