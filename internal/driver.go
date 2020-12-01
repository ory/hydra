package internal

import (
	"context"
	"strings"
	"sync"
	"testing"

	"github.com/spf13/pflag"

	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/x"
	"github.com/ory/x/sqlcon/dockertest"

	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/jwk"
	"github.com/ory/x/logrusx"
)

func resetConfig(p *config.ViperProvider) {
	p.Set(config.ViperKeyBCryptCost, "4")
	p.Set(config.ViperKeySubjectIdentifierAlgorithmSalt, "00000000")
	p.Set(config.ViperKeyGetSystemSecret, []string{"000000000000000000000000000000000000000000000000"})
	p.Set(config.ViperKeyGetCookieSecrets, []string{"000000000000000000000000000000000000000000000000"})
	p.Set(config.ViperKeyLogLevel, "trace")
}

func NewConfigurationWithDefaults() *config.ViperProvider {
	p := config.MustNew(pflag.NewFlagSet("config", pflag.ContinueOnError), logrusx.New("", ""))
	resetConfig(p)
	p.Set("dangerous-force-http", true)
	return p
}

func NewConfigurationWithDefaultsAndHTTPS() *config.ViperProvider {
	p := config.MustNew(pflag.NewFlagSet("config", pflag.ContinueOnError), logrusx.New("", ""))
	resetConfig(p)
	p.Set("dangerous-force-http", false)
	return p
}

func NewRegistryMemory(t *testing.T, c *config.ViperProvider) driver.Registry {
	return newRegistryDefault(t, "memory", c)
}

func NewMockedRegistry(t *testing.T) driver.Registry {
	return newRegistryDefault(t, "memory", NewConfigurationWithDefaults())
}

func NewRegistrySQLFromURL(t *testing.T, url string) driver.Registry {
	return newRegistryDefault(t, url, NewConfigurationWithDefaults())
}

func newRegistryDefault(t *testing.T, url string, c *config.ViperProvider) driver.Registry {
	c.Set(config.ViperKeyLogLevel, "trace")
	c.Set(config.ViperKeyDSN, url)

	r, err := driver.NewRegistryFromDSN(c, logrusx.New("test_hydra", "master"))
	require.NoError(t, err)

	require.NoError(t, r.Init())
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
	c := dockertest.ConnectToTestMySQLPop(t)
	url := c.URL()
	if !strings.HasPrefix(url, "mysql://") {
		url = "mysql://" + url
	}
	require.NoError(t, c.Close())
	return url
}

func ConnectToPG(t *testing.T) string {
	c := dockertest.ConnectToTestPostgreSQLPop(t)
	require.NoError(t, c.Close())
	return c.URL()
}

func ConnectToCRDB(t *testing.T) string {
	c := dockertest.ConnectToTestCockroachDBPop(t)
	url := c.URL()
	if !strings.HasPrefix(url, "cockroach") {
		url = "cockroach://" + strings.Split(url, "://")[1]
	}
	require.NoError(t, c.Close())
	return url
}

func ConnectDatabases(t *testing.T) (pg, mysql, crdb driver.Registry, clean func(*testing.T)) {
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

	pg = NewRegistrySQLFromURL(t, pgURL)
	mysql = NewRegistrySQLFromURL(t, mysqlURL)
	crdb = NewRegistrySQLFromURL(t, crdbURL)
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
	if err := jwk.EnsureAsymmetricKeypairExists(context.Background(), r, new(veryInsecureRS256Generator), key); err != nil {
		panic(err)
	}
}
