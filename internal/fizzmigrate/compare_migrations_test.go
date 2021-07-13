// +build legacy_migration_test

package fizzmigrate

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/pop/v5/logging"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/sqlcon/dockertest"

	"github.com/ory/hydra/internal/fizzmigrate/client"
	"github.com/ory/hydra/internal/fizzmigrate/consent"
	"github.com/ory/hydra/internal/fizzmigrate/jwk"
	"github.com/ory/hydra/internal/fizzmigrate/oauth2"
	"github.com/ory/hydra/persistence/sql"
	"github.com/ory/hydra/x"
)

type migrator interface {
	PlanMigration(string) ([]*migrate.PlannedMigration, error)
	CreateSchemas(string) (int, error)
	CreateMaxSchemas(string, int) (int, error)
}

func connectPostgres(t *testing.T) (*pop.Connection, *sqlx.DB) {
	c := dockertest.ConnectToTestPostgreSQLPop(t)
	db, err := sqlx.Connect("postgres", c.URL())
	require.NoError(t, err)
	return c, db
}

func connectMySQL(t *testing.T) (*pop.Connection, *sqlx.DB) {
	c := dockertest.ConnectToTestMySQLPop(t)
	u := c.URL()
	db, err := sqlx.Connect("mysql", u)
	require.NoError(t, err)
	return c, db
}

func connectCockroach(t *testing.T) (*pop.Connection, *sqlx.DB) {
	c := dockertest.ConnectToTestCockroachDBPop(t)
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

var comments = regexp.MustCompile("(--[^\n]*\n)|(?s:/\\*.+\\*/)")
var migrationTableStatements = regexp.MustCompile("[^;]*(hydra_[a-zA-Z0-9_]*_migration|schema_migration)[^;]*;")

func stripDump(d string) string {
	d = comments.ReplaceAllLiteralString(d, "")
	d = migrationTableStatements.ReplaceAllLiteralString(d, "")
	return strings.ReplaceAll(d, "\r\n", "")
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
	return stripDump(string(dump))
}

var dbConnections = map[string]func(*testing.T) (*pop.Connection, *sqlx.DB){
	"postgres":  connectPostgres,
	"mysql":     connectMySQL,
	"cockroach": connectCockroach,
}

func migrateOldBySingleSteps(t *testing.T, m migrator, db string, stepsDone *int, maxSteps int, afterEach func(int)) {
	startSteps := *stepsDone
	for ; *stepsDone < startSteps+maxSteps; *stepsDone++ {
		n, err := m.CreateMaxSchemas(db, 1)
		require.NoError(t, err)
		require.Equal(t, 1, n)
		afterEach(*stepsDone)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func migrateOldUpSteps(t *testing.T, dbx *sqlx.DB, db string, todo int, afterEach func(int)) {
	stepsDone := 0
	migrateOldBySingleSteps(t, client.NewMigrator(dbx), db, &stepsDone, min(todo, 14), afterEach)
	if todo > stepsDone {
		migrateOldBySingleSteps(t, jwk.NewMigrator(dbx), db, &stepsDone, min(todo-stepsDone, 4), afterEach)
	}
	if todo > stepsDone {
		migrateOldBySingleSteps(t, consent.NewMigrator(dbx), db, &stepsDone, min(todo-stepsDone, 14), afterEach)
	}
	if todo > stepsDone {
		migrateOldBySingleSteps(t, oauth2.NewMigrator(dbx), db, &stepsDone, min(todo-stepsDone, 11), afterEach)
	}
}

func TestCompareMigrations(t *testing.T) {
	for db, connect := range dbConnections {
		t.Run("db="+db, func(t *testing.T) {
			c, dbx := connect(t)
			x.CleanSQLPop(t, c)

			persister, err := sql.NewPersister(c)
			require.NoError(t, err)

			schemasOld := make([]string, 43)
			migrateOldUpSteps(t, dbx, db, 43, func(i int) {
				schemasOld[i] = dump(t, db)
			})

			x.CleanSQLPop(t, c)

			schemasNew := make([]string, 43)
			for i := 0; i < 43; i++ {
				n, err := persister.MigrateUpTo(context.Background(), 1)
				require.NoError(t, err)
				require.Equal(t, 1, n)
				schemasNew[i] = dump(t, db)
			}

			for i, s := range schemasOld {
				require.Equal(t, s, schemasNew[i], "%d", i)
			}
		})
	}
}

func TestMixMigrations(t *testing.T) {
	for db, connect := range dbConnections {
		t.Run("db="+db, func(t *testing.T) {
			pop.SetLogger(func(lvl logging.Level, s string, args ...interface{}) {
				t.Logf("Running SQL Query: %s", fmt.Sprintf(s, args...))
			})

			c, dbx := connect(t)
			persister, err := sql.NewPersister(c)
			require.NoError(t, err)

			schemas := make([]string, 43)
			for i := 0; i < 43; i++ {
				start := time.Now()
				x.CleanSQLPop(t, c)
				migrateOldUpSteps(t, dbx, db, i, func(_ int) {})
				require.NoError(t, persister.MigrateUp(context.Background()))
				t.Logf("migrate took: %dms", time.Now().Sub(start).Milliseconds())

				start = time.Now()
				schemas[i] = dump(t, db)
				t.Logf("dump took: %dms", time.Now().Sub(start).Milliseconds())
			}
			for _, s := range schemas {
				assert.Equal(t, schemas[0], s)
			}
		})
	}
}
