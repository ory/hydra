// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package dockertest

import (
	"cmp"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/pop/v6"

	"github.com/ory/dockertest/v4"
	"github.com/ory/x/stringsx"
)

func ConnectPop(t require.TestingT, url string) (c *pop.Connection) {
	require.EventuallyWithT(t, func(t *assert.CollectT) {
		var err error
		c, err = pop.NewConnection(&pop.ConnectionDetails{
			URL: url,
		})
		require.NoError(t, err)
		require.NoError(t, c.Open())
		require.NoError(t, c.RawQuery("select version()").Exec())
	}, 5*time.Minute, time.Second*5, "could not connect to database at %s", url)
	return
}

// ## PostgreSQL ##

func startPostgreSQL(t testing.TB, version string) dockertest.Resource {
	pool := dockertest.NewPoolT(t, "")
	return pool.RunT(t, "postgres",
		dockertest.WithTag(cmp.Or(version, "16")),
		dockertest.WithoutReuse(),
		dockertest.WithEnv([]string{"PGUSER=postgres", "POSTGRES_PASSWORD=secret", "POSTGRES_DB=postgres"}),
	)
}

// RunTestPostgreSQL runs a PostgreSQL database and returns the URL to it.
// If a docker container is started for the database, the container be removed
// at the end of the test.
func RunTestPostgreSQL(t testing.TB) string {
	if dsn := os.Getenv("TEST_DATABASE_POSTGRESQL"); dsn != "" {
		t.Logf("Skipping Docker setup because environment variable TEST_DATABASE_POSTGRESQL is set to: %s", dsn)
		return dsn
	}
	r := startPostgreSQL(t, "")
	return fmt.Sprintf("postgres://postgres:secret@127.0.0.1:%s/postgres?sslmode=disable", r.GetPort("5432/tcp"))
}

// RunTestPostgreSQLWithVersion connects to a PostgreSQL database.
func RunTestPostgreSQLWithVersion(t testing.TB, version string) string {
	if dsn := os.Getenv("TEST_DATABASE_POSTGRESQL"); dsn != "" {
		return dsn
	}
	r := startPostgreSQL(t, version)
	return fmt.Sprintf("postgres://postgres:secret@127.0.0.1:%s/postgres?sslmode=disable", r.GetPort("5432/tcp"))
}

// ## MySQL ##

func startMySQL(t testing.TB, version string) dockertest.Resource {
	pool := dockertest.NewPoolT(t, "")
	return pool.RunT(t, "mysql",
		dockertest.WithTag(cmp.Or(version, "8.0")),
		dockertest.WithoutReuse(),
		dockertest.WithEnv([]string{
			"MYSQL_ROOT_PASSWORD=secret",
			"MYSQL_ROOT_HOST=%",
		}),
	)
}

// RunTestMySQL runs a MySQL database and returns the URL to it.
// If a docker container is started for the database, the container be removed
// at the end of the test.
func RunTestMySQL(t testing.TB) string {
	if dsn := os.Getenv("TEST_DATABASE_MYSQL"); dsn != "" {
		t.Logf("Skipping Docker setup because environment variable TEST_DATABASE_MYSQL is set to: %s", dsn)
		return dsn
	}
	r := startMySQL(t, "")
	return fmt.Sprintf("mysql://root:secret@tcp(localhost:%s)/mysql?parseTime=true&multiStatements=true", r.GetPort("3306/tcp"))
}

// RunTestMySQLWithVersion runs a MySQL database in the specified version and returns the URL to it.
// If a docker container is started for the database, the container be removed
// at the end of the test.
func RunTestMySQLWithVersion(t testing.TB, version string) string {
	if dsn := os.Getenv("TEST_DATABASE_MYSQL"); dsn != "" {
		t.Logf("Skipping Docker setup because environment variable TEST_DATABASE_MYSQL is set to: %s", dsn)
		return dsn
	}
	r := startMySQL(t, version)
	return fmt.Sprintf("mysql://root:secret@tcp(localhost:%s)/mysql?parseTime=true&multiStatements=true", r.GetPort("3306/tcp"))
}

// ## CockroachDB

func startCockroachDB(t testing.TB, version string) dockertest.Resource {
	pool := dockertest.NewPoolT(t, "")
	return pool.RunT(t, "cockroachdb/cockroach",
		dockertest.WithTag(cmp.Or(version, "latest-v25.4")),
		dockertest.WithoutReuse(),
		dockertest.WithCmd([]string{"start-single-node", "--insecure"}),
	)
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
	r := startCockroachDB(t, version)
	return fmt.Sprintf("cockroach://root@localhost:%s/defaultdb?sslmode=disable", r.GetPort("26257/tcp"))
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

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
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
