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
	"github.com/ory/hydra/v2/spec"
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
	config.KeyDevelopmentMode:                true,
	"serve.public.host":                      "localhost",
}

func NewConfigurationWithDefaults(t testing.TB, opts ...configx.OptionModifier) *config.DefaultProvider {
	allOpts := append([]configx.OptionModifier{
		configx.SkipValidation(),
		configx.WithValues(defaultConfig),
		configx.WithValue("log.leak_sensitive_values", true),
	}, opts...)
	p, err := configx.New(t.Context(), spec.ConfigValidationSchema, allOpts...)
	require.NoError(t, err)
	return config.NewCustom(logrusx.New("", ""), p, contextx.NewTestConfigProvider(spec.ConfigValidationSchema, allOpts...))
}

func NewRegistryMemory(t testing.TB, configOpts ...configx.OptionModifier) driver.Registry {
	return NewRegistrySQLFromURL(t, dbal.NewSQLiteTestDatabase(t), true, configOpts...)
}

func NewRegistrySQLFromURL(t testing.TB, dsn string, migrate bool, configOpts ...configx.OptionModifier) driver.Registry {
	return registryFactory(t,
		NewConfigurationWithDefaults(t, append(
			[]configx.OptionModifier{configx.WithValue(config.KeyDSN, dsn)},
			configOpts...,
		)...), migrate)
}

func registryFactory(t testing.TB, c *config.DefaultProvider, migrate bool) driver.Registry {
	ctx := t.Context()
	sql.SilenceMigrations = true
	l := logrusx.New("test_hydra", "master", logrusx.WithConfigurator(c.Source(ctx)))

	r, err := driver.NewRegistryFromDSN(ctx, c, l, !migrate, migrate, &contextx.Default{})
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
	regs["memory"] = NewRegistryMemory(t)
	if !testing.Short() {
		pg, mysql, crdb := ConnectDatabasesURLs(t)
		regs["postgres"] = NewRegistrySQLFromURL(t, pg, migrate)
		regs["mysql"] = NewRegistrySQLFromURL(t, mysql, migrate)
		regs["cockroach"] = NewRegistrySQLFromURL(t, crdb, migrate)
	}
	return regs
}

func MustEnsureRegistryKeys(ctx context.Context, r driver.Registry, key string) {
	if err := jwk.EnsureAsymmetricKeypairExists(ctx, r, string(jose.ES256), key); err != nil {
		panic(err)
	}
}
