// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package testhelpers

import (
	"context"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/v2/driver"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/jwk"
	"github.com/ory/hydra/v2/persistence/sql"
	"github.com/ory/pop/v6"
	"github.com/ory/x/configx"
	"github.com/ory/x/contextx"
	"github.com/ory/x/dbal"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/sqlcon/dockertest"
)

var defaultConfig = map[string]any{
	config.KeyBCryptCost:                     4,
	config.KeySubjectIdentifierAlgorithmSalt: "00000000",
	config.KeyGetSystemSecret:                []string{"000000000000000000000000000000000000000000000000"},
	config.KeyGetCookieSecrets:               []string{"000000000000000000000000000000000000000000000000"},
	config.KeyLogLevel:                       "trace",
}

func NewConfigurationWithDefaults() *config.DefaultProvider {
	return config.MustNew(context.Background(), logrusx.New("", ""),
		configx.SkipValidation(),
		configx.WithValues(defaultConfig),
		configx.WithValue(config.KeyTLSEnabled, false),
		configx.WithValue("log.leak_sensitive_values", true),
	)
}

func NewConfigurationWithDefaultsAndHTTPS() *config.DefaultProvider {
	return config.MustNew(context.Background(), logrusx.New("", ""),
		configx.SkipValidation(),
		configx.WithValues(defaultConfig),
		configx.WithValue(config.KeyTLSEnabled, true),
	)
}

func NewRegistryMemory(t testing.TB, c *config.DefaultProvider, ctxer contextx.Contextualizer) driver.Registry {
	return registryFactory(t, dbal.NewSQLiteTestDatabase(t), c, true, ctxer)
}

func NewMockedRegistry(t testing.TB, ctxer contextx.Contextualizer) driver.Registry {
	return registryFactory(t, dbal.NewSQLiteTestDatabase(t), NewConfigurationWithDefaults(), true, ctxer)
}

func NewRegistrySQLFromURL(t testing.TB, url string, migrate bool, ctxer contextx.Contextualizer) driver.Registry {
	return registryFactory(t, url, NewConfigurationWithDefaults(), migrate, ctxer)
}

func registryFactory(t testing.TB, url string, c *config.DefaultProvider, migrate bool, ctxer contextx.Contextualizer) driver.Registry {
	return RegistryFactory(t, url, c, !migrate, migrate, ctxer)
}

func RegistryFactory(t testing.TB, url string, c *config.DefaultProvider, networkInit, migrate bool, ctxer contextx.Contextualizer) driver.Registry {
	ctx := t.Context()
	sql.SilenceMigrations = true
	c.MustSet(ctx, config.KeyLogLevel, "trace")
	c.MustSet(ctx, config.KeyDSN, url)
	c.MustSet(ctx, "dev", true)
	l := logrusx.New("test_hydra", "master", logrusx.WithConfigurator(c.Source(ctx)))

	r, err := driver.NewRegistryFromDSN(ctx, c, l, networkInit, migrate, ctxer)
	require.NoError(t, err)

	return r
}

func ConnectToMySQL(t testing.TB) string { return dockertest.RunTestMySQLWithVersion(t, "8.0") }
func ConnectToPG(t testing.TB) string    { return dockertest.RunTestPostgreSQLWithVersion(t, "16") }
func ConnectToCRDB(t testing.TB) string {
	return dockertest.RunTestCockroachDBWithVersion(t, "latest-v24.1")
}

func ConnectDatabasesURLs(t *testing.T) (pgURL, mysqlURL, crdbURL string) {
	wg := sync.WaitGroup{}

	wg.Add(3)
	go func() {
		pgURL = ConnectToPG(t)
		t.Log("Pg done")

		require.EventuallyWithT(t, func(t *assert.CollectT) {
			c, err := pop.NewConnection(&pop.ConnectionDetails{URL: pgURL})
			require.NoError(t, err)
			require.NoError(t, c.Open())
			dbName := "testdb" + strings.ReplaceAll(uuid.Must(uuid.NewV4()).String(), "-", "")
			require.NoError(t, c.RawQuery("CREATE DATABASE "+dbName).Exec())
			pgURL = regexp.MustCompile(`/[a-z0-9]+\?`).ReplaceAllString(pgURL, "/"+dbName+"?")
		}, 20*time.Second, 100*time.Millisecond)

		wg.Done()
	}()
	go func() {
		mysqlURL = ConnectToMySQL(t)
		t.Log("myssql done")

		require.EventuallyWithT(t, func(t *assert.CollectT) {
			c, err := pop.NewConnection(&pop.ConnectionDetails{URL: mysqlURL})
			require.NoError(t, err)
			require.NoError(t, c.Open())
			dbName := "testdb" + strings.ReplaceAll(uuid.Must(uuid.NewV4()).String(), "-", "")
			require.NoError(t, c.RawQuery("CREATE DATABASE "+dbName).Exec())
			mysqlURL = regexp.MustCompile(`/[a-z0-9]+\?`).ReplaceAllString(mysqlURL, "/"+dbName+"?")
		}, 20*time.Second, 100*time.Millisecond)

		wg.Done()
	}()
	go func() {
		crdbURL = ConnectToCRDB(t)
		t.Log("crdb done")

		require.EventuallyWithT(t, func(t *assert.CollectT) {
			c, err := pop.NewConnection(&pop.ConnectionDetails{URL: crdbURL})
			require.NoError(t, err)
			require.NoError(t, c.Open())
			dbName := "testdb" + strings.ReplaceAll(uuid.Must(uuid.NewV4()).String(), "-", "")
			require.NoError(t, c.RawQuery("CREATE DATABASE "+dbName).Exec())
			crdbURL = regexp.MustCompile(`/[a-z0-9]+\?`).ReplaceAllString(crdbURL, "/"+dbName+"?")
		}, 20*time.Second, 100*time.Millisecond)

		wg.Done()
	}()
	t.Log("beginning to wait")
	wg.Wait()
	t.Log("done waiting")

	return
}

func ConnectDatabases(t *testing.T, migrate bool) map[string]driver.Registry {
	regs := make(map[string]driver.Registry)
	regs["memory"] = NewRegistryMemory(t, NewConfigurationWithDefaults(), &contextx.Default{})
	if !testing.Short() {
		pg, mysql, crdb := ConnectDatabasesURLs(t)
		regs["postgres"] = NewRegistrySQLFromURL(t, pg, migrate, &contextx.Default{})
		regs["mysql"] = NewRegistrySQLFromURL(t, mysql, migrate, &contextx.Default{})
		regs["cockroach"] = NewRegistrySQLFromURL(t, crdb, migrate, &contextx.Default{})
	}
	return regs
}

func MustEnsureRegistryKeys(ctx context.Context, r driver.Registry, key string) {
	if err := jwk.EnsureAsymmetricKeypairExists(ctx, r, string(jose.ES256), key); err != nil {
		panic(err)
	}
}
