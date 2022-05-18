package internal

import (
	"context"
	"strings"
	"sync"
	"testing"

	"github.com/ory/x/configx"
	"github.com/ory/x/networkx"

	"github.com/stretchr/testify/require"

	"github.com/ory/hydra/x"
	"github.com/ory/x/sqlcon/dockertest"

	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/config"
	"github.com/ory/hydra/jwk"
	"github.com/ory/x/logrusx"
)

func resetConfig(p *config.Provider) {
	p.MustSet(config.KeyBCryptCost, "4")
	p.MustSet(config.KeySubjectIdentifierAlgorithmSalt, "00000000")
	p.MustSet(config.KeyGetSystemSecret, []string{"000000000000000000000000000000000000000000000000"})
	p.MustSet(config.KeyGetCookieSecrets, []string{"000000000000000000000000000000000000000000000000"})
	p.MustSet(config.KeyLogLevel, "trace")
}

func NewConfigurationWithDefaults() *config.Provider {
	p := config.MustNew(context.Background(), logrusx.New("", ""), configx.SkipValidation())
	resetConfig(p)
	p.MustSet("dangerous-force-http", true)
	return p
}

func NewConfigurationWithDefaultsAndHTTPS() *config.Provider {
	p := config.MustNew(context.Background(), logrusx.New("", ""), configx.SkipValidation())
	resetConfig(p)
	p.MustSet("dangerous-force-http", false)
	return p
}

func NewRegistryMemory(t *testing.T, c *config.Provider) driver.Registry {
	return newRegistryDefault(t, "memory", c, true)
}

func NewMockedRegistry(t *testing.T) driver.Registry {
	return newRegistryDefault(t, "memory", NewConfigurationWithDefaults(), true)
}

func NewRegistrySQLFromURL(t *testing.T, url string, migrate bool) driver.Registry {
	return newRegistryDefault(t, url, NewConfigurationWithDefaults(), migrate)
}

func newRegistryDefault(t *testing.T, url string, c *config.Provider, migrate bool) driver.Registry {
	c.MustSet(config.KeyLogLevel, "trace")
	c.MustSet(config.KeyDSN, url)

	r, err := driver.NewRegistryFromDSN(context.Background(), c, logrusx.New("test_hydra", "master"), false, migrate)
	require.NoError(t, err)

	kg := map[string]jwk.KeyGenerator{
		"RS256": new(veryInsecureRS256Generator),
		"ES256": &jwk.ECDSA256Generator{},
		"ES512": &jwk.ECDSA512Generator{},
		"EdDSA": &jwk.EdDSAGenerator{},
		"HS256": &jwk.HS256Generator{},
		"HS512": &jwk.HS512Generator{},
	}

	r = r.WithKeyGenerators(kg)

	return r
}

func CleanAndMigrate(reg driver.Registry) func(*testing.T) {
	return func(t *testing.T) {
		net := &networkx.Network{}
		recreateNetwork := reg.Persister().Connection(context.Background()).First(net) == nil
		x.CleanSQLPop(t, reg.Persister().Connection(context.Background()))
		require.NoError(t, reg.Persister().MigrateUp(context.Background()))
		if recreateNetwork {
			require.NoError(t, reg.Persister().Connection(context.Background()).RawQuery("DELETE FROM networks").Exec())
			require.NoError(t, reg.Persister().Connection(context.Background()).Create(net))
		}
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

func ConnectDatabases(t *testing.T, migrate bool) (pg, mysql, crdb driver.Registry, clean func(*testing.T)) {
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

	pg = NewRegistrySQLFromURL(t, pgURL, migrate)
	mysql = NewRegistrySQLFromURL(t, mysqlURL, migrate)
	crdb = NewRegistrySQLFromURL(t, crdbURL, migrate)
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
