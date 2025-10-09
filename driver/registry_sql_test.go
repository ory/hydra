// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/sessions"
	pkgerr "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/client"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/persistence/sql"
	"github.com/ory/x/configx"
	"github.com/ory/x/dbal"
	"github.com/ory/x/httpx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/popx"
	"github.com/ory/x/randx"
)

func init() {
	sql.SilenceMigrations = true
}

func TestGetJWKSFetcherStrategyHostEnforcement(t *testing.T) {
	t.Parallel()

	r, err := New(t.Context(), WithConfigOptions(
		configx.WithValues(map[string]any{
			config.KeyDSN:                         "memory",
			config.HSMEnabled:                     "false",
			config.KeyClientHTTPNoPrivateIPRanges: true,
		}),
		configx.WithConfigFiles("../internal/.hydra.yaml"),
	))
	require.NoError(t, err)

	_, err = r.OAuth2Config().GetJWKSFetcherStrategy(t.Context()).Resolve(t.Context(), "http://localhost:8080", true)
	require.ErrorAs(t, err, new(httpx.ErrPrivateIPAddressDisallowed))
}

func TestRegistrySQL_newKeyStrategy_handlesNetworkError(t *testing.T) {
	t.Parallel()

	// Test ensures any network specific error is logged with a
	// specific message when attempting to create a new key strategy: issue #2338

	hook := test.Hook{} // Test hook for asserting log messages

	l := logrusx.New("", "", logrusx.WithHook(&hook))
	l.Logrus().SetOutput(io.Discard)

	// Create a config and set a valid but unresolvable DSN
	c := config.MustNew(t, l,
		configx.WithConfigFiles("../internal/.hydra.yaml"),
		configx.WithValues(map[string]any{
			config.KeyDSN:     "postgres://user:password@127.0.0.1:9999/postgres",
			config.HSMEnabled: false,
		}),
	)

	r, err := newRegistryWithoutInit(c, l)
	if err != nil {
		t.Errorf("Failed to create registry: %s", err)
		return
	}

	r.initialPing = failedPing(errors.New("snizzles"))

	assert.ErrorContains(t,
		r.Init(t.Context(), true, false, nil, nil),
		"snizzles",
	)

	assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)
	assert.Contains(t, hook.LastEntry().Message, "snizzles")
}

func TestRegistrySQL_CookieStore_MaxAgeZero(t *testing.T) {
	t.Parallel()

	// Test ensures that CookieStore MaxAge option is equal to zero after initialization

	r, err := New(t.Context(), SkipNetworkInit(), DisableValidation(), WithConfigOptions(
		configx.WithValues(map[string]any{
			config.KeyDSN:             dbal.NewSQLiteInMemoryDatabase(t.Name()),
			config.KeyGetSystemSecret: []string{randx.MustString(32, randx.AlphaNum)},
		}),
	))
	require.NoError(t, err)

	s, err := r.CookieStore(t.Context())
	require.NoError(t, err)
	cs := s.(*sessions.CookieStore)

	assert.Equal(t, cs.Options.MaxAge, 0)
}

func TestRegistrySQL_HTTPClient(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	r, err := New(t.Context(), SkipNetworkInit(), DisableValidation(), WithConfigOptions(configx.WithValues(map[string]interface{}{
		config.KeyDSN:                              dbal.NewSQLiteInMemoryDatabase(t.Name()),
		config.KeyClientHTTPNoPrivateIPRanges:      true,
		config.KeyClientHTTPPrivateIPExceptionURLs: []string{ts.URL + "/exception/*"},
	})))
	require.NoError(t, err)

	t.Run("case=matches exception glob", func(t *testing.T) {
		res, err := r.HTTPClient(t.Context()).Get(ts.URL + "/exception/foo")
		require.NoError(t, err)
		assert.Equal(t, 200, res.StatusCode)
	})

	t.Run("case=does not match exception glob", func(t *testing.T) {
		_, err := r.HTTPClient(t.Context()).Get(ts.URL + "/foo")
		assert.ErrorContains(t, err, "prohibited IP address")
	})
}

func TestDefaultKeyManager_HsmDisabled(t *testing.T) {
	t.Parallel()

	r, err := New(t.Context(),
		SkipNetworkInit(),
		DisableValidation(),
		WithConfigOptions(
			configx.SkipValidation(),
			configx.WithValues(map[string]any{
				config.KeyDSN:     "memory",
				config.HSMEnabled: false,
			}),
		),
	)
	require.NoError(t, err)
	assert.IsType(t, &sql.JWKPersister{}, r.KeyManager())
}

func TestDbUnknownTableColumns(t *testing.T) {
	t.Parallel()

	dsn := dbal.NewSQLiteTestDatabase(t)
	reg, err := New(t.Context(), WithConfigOptions(configx.WithValue("dsn", dsn)), WithAutoMigrate())
	require.NoError(t, err)

	statement := `ALTER TABLE "hydra_client" ADD COLUMN "temp_column" VARCHAR(128) NOT NULL DEFAULT '';`
	require.NoError(t, reg.Persister().Connection(t.Context()).RawQuery(statement).Exec())

	cl := &client.Client{
		ID: strconv.Itoa(rand.Int()),
	}
	require.NoError(t, reg.Persister().CreateClient(t.Context(), cl))
	getClients := func(ctx context.Context, reg *RegistrySQL) ([]client.Client, error) {
		readClients := make([]client.Client, 0)
		conn := reg.Persister().Connection(ctx)
		cols := popx.DBColumns[client.Client](conn.Dialect)
		return readClients, conn.RawQuery(fmt.Sprintf(`SELECT %s, temp_column FROM "hydra_client"`, cols)).All(&readClients)
	}

	t.Run("with ignore disabled (default behavior)", func(t *testing.T) {
		_, err := getClients(t.Context(), reg)
		assert.ErrorContains(t, err, "missing destination name temp_column")
	})

	t.Run("with ignore enabled", func(t *testing.T) {
		reg, err := New(t.Context(), WithConfigOptions(
			configx.WithValue("dsn", dsn),
			configx.WithValue(config.KeyDBIgnoreUnknownTableColumns, true),
		))
		require.NoError(t, err)

		actual, err := getClients(t.Context(), reg)
		require.NoError(t, err)
		assert.Len(t, actual, 1)
	})
}

func failedPing(err error) func(context.Context, *logrusx.Logger, *sql.BasePersister) error {
	return func(context.Context, *logrusx.Logger, *sql.BasePersister) error {
		return pkgerr.WithStack(err)
	}
}
