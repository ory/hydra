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
	"golang.org/x/sync/errgroup"

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
	"github.com/ory/x/servicelocatorx"
	"github.com/ory/x/sqlcon/dockertest"
)

var ConfigDefaults = []configx.OptionModifier{
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
	configx.WithValue("log.leak_sensitive_values", true),
}

func NewConfigurationWithDefaults(t testing.TB, opts ...configx.OptionModifier) *config.DefaultProvider {
	p, err := configx.New(t.Context(), spec.ConfigValidationSchema, append(ConfigDefaults, opts...)...)
	require.NoError(t, err)
	return config.NewCustom(logrusx.New("", ""), p, contextx.NewTestConfigProvider(spec.ConfigValidationSchema, append(ConfigDefaults, opts...)...))
}

func NewRegistryMemory(t testing.TB, opts ...driver.OptionsModifier) *driver.RegistrySQL {
	reg, err := NewRegistrySQLFromURL(t.Context(), dbal.NewSQLiteTestDatabase(t), true, opts...)
	require.NoError(t, err)
	return reg
}

func NewRegistrySQLFromURL(ctx context.Context, dsn string, migrate bool, opts ...driver.OptionsModifier) (*driver.RegistrySQL, error) {
	sql.SilenceMigrations = true

	configOpts := append(ConfigDefaults, configx.WithValue(config.KeyDSN, dsn))
	regOpts := append([]driver.OptionsModifier{
		driver.WithConfigOptions(configOpts...),
		driver.WithServiceLocatorOptions(servicelocatorx.WithContextualizer(contextx.NewTestConfigProvider(spec.ConfigValidationSchema, configOpts...))),
	}, opts...)
	if migrate {
		regOpts = append(regOpts, driver.WithAutoMigrate())
	}

	return driver.New(ctx, regOpts...)
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

func ConnectDatabases(t *testing.T, migrate bool, opts ...driver.OptionsModifier) map[string]*driver.RegistrySQL {
	regs := make(map[string]*driver.RegistrySQL)
	regs["memory"] = NewRegistryMemory(t, opts...)
	if !testing.Short() {
		pg, mysql, crdb := ConnectDatabasesURLs(t)
		eg, ctx := errgroup.WithContext(t.Context())
		eg.Go(func() (err error) {
			regs["postgres"], err = NewRegistrySQLFromURL(ctx, pg, migrate, opts...)
			return
		})
		eg.Go(func() (err error) {
			regs["mysql"], err = NewRegistrySQLFromURL(ctx, mysql, migrate, opts...)
			return
		})
		eg.Go(func() (err error) {
			regs["cockroach"], err = NewRegistrySQLFromURL(ctx, crdb, migrate, opts...)
			return
		})
		require.NoError(t, eg.Wait())
	}
	return regs
}

func MustEnsureRegistryKeys(t testing.TB, r *driver.RegistrySQL, key string) {
	jwk.EnsureAsymmetricKeypairExists(t, r, string(jose.ES256), key)
}
