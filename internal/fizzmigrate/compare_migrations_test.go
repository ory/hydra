package fizzmigrate

import (
	"context"
	"fmt"
	"github.com/gobuffalo/pop/v5"
	"github.com/jmoiron/sqlx"
	"github.com/ory/hydra/internal/fizzmigrate/client"
	"github.com/ory/hydra/internal/fizzmigrate/consent"
	"github.com/ory/hydra/internal/fizzmigrate/jwk"
	"github.com/ory/hydra/internal/fizzmigrate/oauth2"
	"github.com/ory/hydra/persistence/sql"
	"github.com/ory/hydra/x"
	"github.com/ory/x/sqlcon/dockertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"modernc.org/mathutil"
	"os/exec"
	"strings"
	"testing"
)

func connectPostgres(t *testing.T) (*pop.Connection, *sqlx.DB) {
	c := dockertest.ConnectToTestPostgreSQLPop(t)
	db, err := sqlx.Connect("postgres", c.URL())
	require.NoError(t, err)
	return c, db
}

func getContainerID(t *testing.T, containerPort string) string {
	cid, err := exec.Command("docker", "ps", "-f", fmt.Sprintf("expose=%s", containerPort), "-q").CombinedOutput()
	require.NoError(t, err)
	containerID := strings.TrimSuffix(string(cid), "\n")
	require.False(t, strings.Contains(containerID, "\n"), "there is more than one docker container running with port %s, I am confused: %s", containerPort, containerID)
	return containerID
}

func dumpArgs(t *testing.T, db string) []string {
	switch db {
	case "postgres":
		return []string{"exec", "-t", getContainerID(t, "5432"), "pg_dump", "-U", "postgres", "-s", "-T", "hydra_*_migration", "-T", "schema_migration"}
	case "mysql":
		return []string{"exec", "-t", getContainerID(t, "3306"), "/usr/bin/mysqldump", "-u", "root", "--password=secret", "mysql"}
	case "cockroach":
		return []string{"exec", "-t", getContainerID(t, "26257"), "./cockroach", "dump", "defaultdb", "--insecure", "--dump-mode=schema"}
	}
	t.Fail()
	return []string{}
}

func dump(t *testing.T, db string) string {
	dump, err := exec.Command("docker", dumpArgs(t, db)...).CombinedOutput()
	require.NoError(t, err, "%s", dump)
	return string(dump)
}

var dbConnections = map[string]func(*testing.T) (*pop.Connection, *sqlx.DB){
	"postgres": connectPostgres,
}

func migrateOldBySingleSteps(t *testing.T, m migrator, db string, stepsDone *int, maxSteps int, afterEach func(int)) {
	for ; *stepsDone < maxSteps; {
		n, err := m.CreateMaxSchemas(db, 1)
		require.NoError(t, err)
		require.Equal(t, 1, n)
		*stepsDone++
		afterEach(*stepsDone)
	}
}

func migrateOldUpSteps(t *testing.T, dbx *sqlx.DB, todo int, afterEach func(int)) {
	stepsDone := 0
	migrateOldBySingleSteps(t, client.NewMigrator(dbx), dbx.DriverName(), &stepsDone, mathutil.Min(todo, 14), afterEach)
	if stepsDone < 14+4 && todo > stepsDone {
		migrateOldBySingleSteps(t, jwk.NewMigrator(dbx), dbx.DriverName(), &stepsDone, mathutil.Min(todo-stepsDone, 14+4-stepsDone), afterEach)
	}
	if stepsDone < 14+4+14 && todo > stepsDone {
		migrateOldBySingleSteps(t, consent.NewMigrator(dbx), dbx.DriverName(), &stepsDone, mathutil.Min(todo-stepsDone, 14+4+14-stepsDone), afterEach)
	}
	if stepsDone < 14+4+14+11 && todo > stepsDone {
		migrateOldBySingleSteps(t, oauth2.NewMigrator(dbx), dbx.DriverName(), &stepsDone, mathutil.Min(todo-stepsDone, 14+4+14+11-stepsDone), afterEach)
	}
}

func TestCompareMigrations(t *testing.T) {
	for db, connect := range dbConnections {
		t.Run("db="+db, func(t *testing.T) {
			c, dbx := connect(t)
			x.CleanSQLPop(t, c)

			persister, err := sql.NewPersister(c)
			require.NoError(t, err)

			schemasOld := make([]string, 14+4+14+11)
			migrateOldUpSteps(t, dbx, 14+4+14+11, func(i int) {
				schemasOld[i] = dump(t, db)
			})

			x.CleanSQLPop(t, c)

			schemasNew := make([]string, 14+4+14+11)
			for i := 0; i < 14+4+14+11; i++ {
				n, err := persister.MigrateUpTo(context.Background(), 1)
				require.NoError(t, err)
				require.Equal(t, 1, n)
				schemasNew[i] = dump(t, db)
			}

			for i, s := range schemasOld {
				require.Equal(t, s, schemasNew[i], "%d", i)
			}
			assert.Equal(t, schemasOld, schemasNew)
			schemasOld = nil
			schemasNew = nil
		})
	}
}

func TestMixMigrations(t *testing.T) {
	for db, connect := range dbConnections {
		t.Run("db="+db, func(t *testing.T) {
			c, dbx := connect(t)
			persister, err := sql.NewPersister(c)
			require.NoError(t, err)

			schemas := make([]string, 14+4+14+11)
			for i := 0; i < 14+4+14+11; i++ {
				x.CleanSQLPop(t, c)
				migrateOldUpSteps(t, dbx, i, func(_ int) {})
				require.NoError(t, persister.MigrateUp(context.Background()))
				schemas[i] = dump(t, db)
			}
			for _, s := range schemas {
				assert.Equal(t, schemas[0], s)
			}
		})
	}
}
