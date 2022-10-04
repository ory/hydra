package internal

import (
	"context"

	"sync"
	"testing"

	"github.com/ory/x/configx"

	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/x"
	"github.com/ory/x/contextx"
	"github.com/ory/x/sqlcon/dockertest"

	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/jwk"
	"github.com/ory/x/logrusx"
)

func resetConfig(p *config.DefaultProvider) {
	p.MustSet(context.Background(), config.KeyBCryptCost, "4")
	p.MustSet(context.Background(), config.KeySubjectIdentifierAlgorithmSalt, "00000000")
	p.MustSet(context.Background(), config.KeyGetSystemSecret, []string{"000000000000000000000000000000000000000000000000"})
	p.MustSet(context.Background(), config.KeyGetCookieSecrets, []string{"000000000000000000000000000000000000000000000000"})
	p.MustSet(context.Background(), config.KeyLogLevel, "trace")
}

func NewConfigurationWithDefaults() *config.DefaultProvider {
	p := config.MustNew(context.Background(), logrusx.New("", ""), configx.SkipValidation())
	resetConfig(p)
	p.MustSet(context.Background(), config.KeyTLSEnabled, false)
	return p
}

func NewConfigurationWithDefaultsAndHTTPS() *config.DefaultProvider {
	p := config.MustNew(context.Background(), logrusx.New("", ""), configx.SkipValidation())
	resetConfig(p)
	p.MustSet(context.Background(), config.KeyTLSEnabled, true)
	return p
}

func NewRegistryMemory(t *testing.T, c *config.DefaultProvider, ctxer contextx.Contextualizer) driver.Registry {
	return newRegistryDefault(t, "memory", c, true, ctxer)
}

func NewMockedRegistry(t *testing.T, ctxer contextx.Contextualizer) driver.Registry {
	return newRegistryDefault(t, "memory", NewConfigurationWithDefaults(), true, ctxer)
}

func NewRegistrySQLFromURL(t *testing.T, url string, migrate bool, ctxer contextx.Contextualizer) driver.Registry {
	return newRegistryDefault(t, url, NewConfigurationWithDefaults(), migrate, ctxer)
}

func newRegistryDefault(t *testing.T, url string, c *config.DefaultProvider, migrate bool, ctxer contextx.Contextualizer) driver.Registry {
	ctx := context.Background()
	c.MustSet(ctx, config.KeyLogLevel, "trace")
	c.MustSet(ctx, config.KeyDSN, url)
	c.MustSet(ctx, "dev", true)

	r, err := driver.NewRegistryFromDSN(ctx, c, logrusx.New("test_hydra", "master"), false, migrate, ctxer)
	require.NoError(t, err)

	return r
}

func CleanAndMigrate(reg driver.Registry) func(*testing.T) {
	return func(t *testing.T) {
		x.CleanSQLPop(t, reg.Persister().Connection(context.Background()))
		require.NoError(t, reg.Persister().MigrateUp(context.Background()))
		t.Log("clean and migrate done")
	}
}

func ConnectToMySQL(t *testing.T) string {
	return dockertest.RunTestMySQLWithVersion(t, "11.8")
}

func ConnectToPG(t *testing.T) string {
	return dockertest.RunTestPostgreSQLWithVersion(t, "11.8")
}

func ConnectToCRDB(t *testing.T) string {
	return dockertest.RunTestCockroachDBWithVersion(t, "v22.1.2")
}

func ConnectDatabases(t *testing.T, migrate bool, ctxer contextx.Contextualizer) (pg, mysql, crdb driver.Registry, clean func(*testing.T)) {
	var pgURL, mysqlURL, crdbURL string
	wg := sync.WaitGroup{}

	wg.Add(3)
	go func() {
		pgURL = ConnectToPG(t)
		t.Log("Pg done")
		wg.Done()
	}()
	go func() {
		mysqlURL = ConnectToMySQL(t)
		t.Log("myssql done")
		wg.Done()
	}()
	go func() {
		crdbURL = ConnectToCRDB(t)
		t.Log("crdb done")
		wg.Done()
	}()
	t.Log("beginning to wait")
	wg.Wait()
	t.Log("done waiting")

	pg = NewRegistrySQLFromURL(t, pgURL, migrate, ctxer)
	mysql = NewRegistrySQLFromURL(t, mysqlURL, migrate, ctxer)
	crdb = NewRegistrySQLFromURL(t, crdbURL, migrate, ctxer)
	dbs := []driver.Registry{pg, mysql, crdb}

	clean = func(t *testing.T) {
		wg := sync.WaitGroup{}

		wg.Add(len(dbs))
		for _, db := range dbs {
			go func(db driver.Registry) {
				defer wg.Done()
				CleanAndMigrate(db)(t)
			}(db)
		}
		wg.Wait()
	}
	clean(t)
	return
}

func MustEnsureRegistryKeys(r driver.Registry, key string) {
	if err := jwk.EnsureAsymmetricKeypairExists(context.Background(), r, "RS256", key); err != nil {
		panic(err)
	}
}
