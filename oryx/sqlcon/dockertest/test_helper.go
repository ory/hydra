// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package dockertest

import (
	"cmp"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/require"

	"github.com/ory/pop/v6"

	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/resilience"
	"github.com/ory/x/stringsx"
)

type dockerPool interface {
	Purge(r *dockertest.Resource) error
	Run(repository, tag string, env []string) (*dockertest.Resource, error)
	RunWithOptions(opts *dockertest.RunOptions, hcOpts ...func(*dc.HostConfig)) (*dockertest.Resource, error)
}

var (
	pool      dockerPool
	resources []*dockertest.Resource
	mux       sync.Mutex
)

func init() {
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		panic(err)
	}
}

// KillAllTestDatabases deletes all test databases.
func KillAllTestDatabases() {
	mux.Lock()
	defer mux.Unlock()
	for _, r := range resources {
		if err := pool.Purge(r); err != nil {
			log.Printf("Failed to purge resource: %s", err)
		}
	}

	resources = nil
}

// Register sets up OnExit.
func Register() *OnExit {
	onexit := NewOnExit()
	onexit.Add(func() {
		KillAllTestDatabases()
	})
	return onexit
}

func ConnectPop(t require.TestingT, url string) (c *pop.Connection) {
	require.NoError(t, resilience.Retry(logrusx.New("", ""), time.Second*5, time.Minute*5, func() error {
		var err error
		c, err = pop.NewConnection(&pop.ConnectionDetails{
			URL: url,
		})
		if err != nil {
			log.Printf("could not create pop connection")
			return err
		}
		if err := c.Open(); err != nil {
			// an Open error probably means we have a problem with the connections config
			log.Printf("could not open pop connection: %+v", err)
			return err
		}
		return c.RawQuery("select version()").Exec()
	}))
	return
}

// ## PostgreSQL ##

func startPostgreSQL(version string) (*dockertest.Resource, error) {
	resource, err := pool.Run("postgres", cmp.Or(version, "16"), []string{"PGUSER=postgres", "POSTGRES_PASSWORD=secret", "POSTGRES_DB=postgres"})
	if err == nil {
		mux.Lock()
		resources = append(resources, resource)
		mux.Unlock()
	}
	return resource, err
}

// RunTestPostgreSQL runs a PostgreSQL database and returns the URL to it.
// If a docker container is started for the database, the container be removed
// at the end of the test.
func RunTestPostgreSQL(t testing.TB) string {
	if dsn := os.Getenv("TEST_DATABASE_POSTGRESQL"); dsn != "" {
		t.Logf("Skipping Docker setup because environment variable TEST_DATABASE_POSTGRESQL is set to: %s", dsn)
		return dsn
	}

	u, cleanup, err := runPosgreSQLCleanup("")
	require.NoError(t, err)
	t.Cleanup(cleanup)

	return u
}

func runPosgreSQLCleanup(version string) (string, func(), error) {
	resource, err := startPostgreSQL(version)
	if err != nil {
		return "", func() {}, err
	}

	return fmt.Sprintf("postgres://postgres:secret@127.0.0.1:%s/postgres?sslmode=disable", resource.GetPort("5432/tcp")),
		func() { _ = pool.Purge(resource) }, nil
}

// RunTestPostgreSQLWithVersion connects to a PostgreSQL database .
func RunTestPostgreSQLWithVersion(t testing.TB, version string) string {
	if dsn := os.Getenv("TEST_DATABASE_POSTGRESQL"); dsn != "" {
		return dsn
	}

	resource, err := startPostgreSQL(version)
	require.NoError(t, err)
	return fmt.Sprintf("postgres://postgres:secret@127.0.0.1:%s/postgres?sslmode=disable", resource.GetPort("5432/tcp"))
}

// ## MySQL ##

func startMySQL(version string) (*dockertest.Resource, error) {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mysql",
		Tag:        cmp.Or(version, "8.0"),
		Env: []string{
			"MYSQL_ROOT_PASSWORD=secret",
			"MYSQL_ROOT_HOST=%",
		},
	})
	if err != nil {
		return nil, err
	}
	mux.Lock()
	resources = append(resources, resource)
	mux.Unlock()
	return resource, nil
}

func runMySQLCleanup(version string) (string, func(), error) {
	resource, err := startMySQL(version)
	if err != nil {
		return "", func() {}, err
	}

	return fmt.Sprintf("mysql://root:secret@tcp(localhost:%s)/mysql?parseTime=true&multiStatements=true", resource.GetPort("3306/tcp")),
		func() { _ = pool.Purge(resource) }, nil
}

// RunTestMySQL runs a MySQL database and returns the URL to it.
// If a docker container is started for the database, the container be removed
// at the end of the test.
func RunTestMySQL(t testing.TB) string {
	if dsn := os.Getenv("TEST_DATABASE_MYSQL"); dsn != "" {
		t.Logf("Skipping Docker setup because environment variable TEST_DATABASE_MYSQL is set to: %s", dsn)
		return dsn
	}

	u, cleanup, err := runMySQLCleanup("")
	require.NoError(t, err)
	t.Cleanup(cleanup)

	return u
}

// RunTestMySQLWithVersion runs a MySQL database in the specified version and returns the URL to it.
// If a docker container is started for the database, the container be removed
// at the end of the test.
func RunTestMySQLWithVersion(t testing.TB, version string) string {
	if dsn := os.Getenv("TEST_DATABASE_MYSQL"); dsn != "" {
		t.Logf("Skipping Docker setup because environment variable TEST_DATABASE_MYSQL is set to: %s", dsn)
		return dsn
	}

	u, cleanup, err := runMySQLCleanup(version)
	require.NoError(t, err)
	t.Cleanup(cleanup)

	return u
}

// ## CockroachDB

func startCockroachDB(version string) (*dockertest.Resource, error) {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "cockroachdb/cockroach",
		Tag:        cmp.Or(version, "latest-v25.4"),
		Cmd:        []string{"start-single-node", "--insecure"},
	})
	if err == nil {
		mux.Lock()
		resources = append(resources, resource)
		mux.Unlock()
	}
	return resource, err
}

func runCockroachDBWithVersionCleanup(version string) (string, func(), error) {
	resource, err := startCockroachDB(version)
	if err != nil {
		return "", func() {}, err
	}

	return fmt.Sprintf("cockroach://root@localhost:%s/defaultdb?sslmode=disable", resource.GetPort("26257/tcp")),
		func() { _ = pool.Purge(resource) },
		nil
}

// RunTestCockroachDB runs a CockroachDB database and returns the URL to it.
// If a docker container is started for the database, the container be removed
// at the end of the test.
func RunTestCockroachDB(t testing.TB) string {
	return RunTestCockroachDBWithVersion(t, "")
}

// RunTestCockroachDBWithVersion runs a CockroachDB database and returns the URL to it.
// If a docker container is started for the database, the container be removed
// at the end of the test.
func RunTestCockroachDBWithVersion(t testing.TB, version string) string {
	if dsn := os.Getenv("TEST_DATABASE_COCKROACHDB"); dsn != "" {
		t.Logf("Skipping Docker setup because environment variable TEST_DATABASE_COCKROACHDB is set to: %s", dsn)
		return dsn
	}

	u, cleanup, err := runCockroachDBWithVersionCleanup(version)
	require.NoError(t, err)
	t.Cleanup(cleanup)

	return u
}

func DumpSchema(t testing.TB, c *pop.Connection) string {
	name, database, port := c.Dialect.Name(), c.Dialect.Details().Database, c.Dialect.Details().Port
	t.Logf("Dumping schema for dialect %s, database %s on port %s", name, database, port)

	var cmd []string
	var appendToDump string
	switch dialects := stringsx.SwitchExact(name); {
	case dialects.AddCase("sqlite3"):
		return dumpSQLiteSchema(t, c)
	case dialects.AddCase("postgres"):
		cmd = []string{"pg_dump", "--username", "postgres", "--schema-only", "--dbname", database}
		// we need to set the search path because the postgres dump always unsets it
		appendToDump = "SET search_path TO public;\n"
	case dialects.AddCase("mysql"):
		cmd = []string{"mysqldump", "--user", "root", "--password=secret", "--no-data", database}
	case dialects.AddCase("cockroach"):
		cmd = []string{"cockroach", "sql", "--insecure", "--database", database, "--execute", "SHOW CREATE ALL TABLES; SHOW CREATE ALL TYPES;", "--format", "raw"}
	default:
		t.Log(dialects.ToUnknownCaseErr())
		t.FailNow()
		return ""
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	require.NoError(t, err)
	containers, err := cli.ContainerList(t.Context(), container.ListOptions{
		Filters: filters.NewArgs(filters.Arg("publish", port)),
	})
	require.NoError(t, err)
	require.Lenf(t, containers, 1, "expected exactly one %s container with port %s", name, port)

	process, err := cli.ContainerExecCreate(t.Context(), containers[0].ID, container.ExecOptions{
		Tty:          true,
		AttachStdout: true,
		Cmd:          cmd,
	})
	require.NoError(t, err)

	resp, err := cli.ContainerExecAttach(t.Context(), process.ID, container.ExecAttachOptions{
		Tty: true,
	})
	require.NoError(t, err)
	dump, err := io.ReadAll(resp.Reader)
	require.NoErrorf(t, err, "%s", dump)

	d := string(dump) + appendToDump
	d = regexp.MustCompile(`(--|#|\\|mysqldump|SHOW CREATE)[^\n]*\n`).ReplaceAllLiteralString(d, "") // comments and other non-schema lines
	d = strings.ReplaceAll(d, "\r\n", "\n")
	d = regexp.MustCompile(`\n\n+`).ReplaceAllLiteralString(d, "\n\n")
	return d
}

func dumpSQLiteSchema(t testing.TB, c *pop.Connection) string {
	var sqls []string
	require.NoError(t, c.RawQuery("SELECT sql FROM sqlite_master WHERE type IN ('table', 'index', 'trigger', 'view') AND name NOT LIKE 'sqlite_%' ORDER BY name").All(&sqls))
	return strings.Join(sqls, ";\n") + ";\n"
}
