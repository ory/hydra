// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package dockertest

import (
	"context"
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
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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

// Parallel runs tasks in parallel.
func Parallel(fs []func()) {
	wg := sync.WaitGroup{}

	wg.Add(len(fs))
	for _, f := range fs {
		go func(ff func()) {
			defer wg.Done()
			ff()
		}(f)
	}

	wg.Wait()
}

func connect(dialect, driver, dsn string) (db *sqlx.DB, err error) {
	if scheme := strings.Split(dsn, "://")[0]; scheme == "mysql" {
		dsn = strings.Replace(dsn, "mysql://", "", -1)
	} else if scheme == "cockroach" {
		dsn = strings.Replace(dsn, "cockroach://", "postgres://", 1)
	}
	err = resilience.Retry(
		logrusx.New("", ""),
		time.Second*5,
		time.Minute*5,
		func() (err error) {
			db, err = sqlx.Open(dialect, dsn)
			if err != nil {
				log.Printf("Connecting to database %s failed: %s", driver, err)
				return err
			}

			if err := db.Ping(); err != nil {
				log.Printf("Pinging database %s failed: %s", driver, err)
				return err
			}

			return nil
		},
	)
	if err != nil {
		return nil, errors.Errorf("Unable to connect to %s (%s): %s", driver, dsn, err)
	}
	log.Printf("Connected to database %s", driver)
	return db, nil
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
	resource, err := pool.Run("postgres", stringsx.Coalesce(version, "16"), []string{"PGUSER=postgres", "POSTGRES_PASSWORD=secret", "POSTGRES_DB=postgres"})
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

// RunPostgreSQL runs a PostgreSQL database and returns the URL to it.
func RunPostgreSQL() (string, error) {
	dsn, _, err := runPosgreSQLCleanup("")
	return dsn, err
}

func runPosgreSQLCleanup(version string) (string, func(), error) {
	resource, err := startPostgreSQL(version)
	if err != nil {
		return "", func() {}, err
	}

	return fmt.Sprintf("postgres://postgres:secret@127.0.0.1:%s/postgres?sslmode=disable", resource.GetPort("5432/tcp")),
		func() { _ = pool.Purge(resource) }, nil
}

// ConnectToTestPostgreSQL connects to a PostgreSQL database.
func ConnectToTestPostgreSQL() (*sqlx.DB, error) {
	if dsn := os.Getenv("TEST_DATABASE_POSTGRESQL"); dsn != "" {
		return connect("pgx", "postgres", dsn)
	}

	resource, err := startPostgreSQL("")
	if err != nil {
		return nil, errors.Wrap(err, "Could not start resource")
	}

	db := bootstrap("postgres://postgres:secret@localhost:%s/postgres?sslmode=disable", "5432/tcp", "pgx", pool, resource)
	return db, nil
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

// ConnectToTestPostgreSQLPop connects to a test PostgreSQL database.
// If a docker container is started for the database, the container be removed
// at the end of the test.
func ConnectToTestPostgreSQLPop(t testing.TB) *pop.Connection {
	url := RunTestPostgreSQL(t)
	return ConnectPop(t, url)
}

// ## MySQL ##

func startMySQL(version string) (*dockertest.Resource, error) {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mysql",
		Tag:        stringsx.Coalesce(version, "8.0"),
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

// RunMySQL runs a RunMySQL database and returns the URL to it.
func RunMySQL() (string, error) {
	dsn, _, err := runMySQLCleanup("")
	return dsn, err
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

// ConnectToTestMySQL connects to a MySQL database.
func ConnectToTestMySQL() (*sqlx.DB, error) {
	if dsn := os.Getenv("TEST_DATABASE_MYSQL"); dsn != "" {
		log.Println("Found mysql test database config, skipping dockertest...")
		return connect("mysql", "mysql", dsn)
	}

	resource, err := startMySQL("")
	if err != nil {
		return nil, errors.Wrap(err, "Could not start resource")
	}

	db := bootstrap("root:secret@(localhost:%s)/mysql?parseTime=true", "3306/tcp", "mysql", pool, resource)
	return db, nil
}

func ConnectToTestMySQLPop(t testing.TB) *pop.Connection {
	url := RunTestMySQL(t)
	return ConnectPop(t, url)
}

// ## CockroachDB

func startCockroachDB(version string) (*dockertest.Resource, error) {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "cockroachdb/cockroach",
		Tag:        stringsx.Coalesce(version, "latest-v25.3"),
		Cmd:        []string{"start-single-node", "--insecure"},
	})
	if err == nil {
		mux.Lock()
		resources = append(resources, resource)
		mux.Unlock()
	}
	return resource, err
}

// RunCockroachDB runs a CockroachDB database and returns the URL to it.
func RunCockroachDB() (string, error) {
	return RunCockroachDBWithVersion("")
}

// RunCockroachDBWithVersion runs a CockroachDB database with the specified version and returns the URL to it.
func RunCockroachDBWithVersion(version string) (string, error) {
	resource, err := startCockroachDB(version)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("cockroach://root@localhost:%s/defaultdb?sslmode=disable", resource.GetPort("26257/tcp")), nil
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

// ConnectToTestCockroachDB connects to a CockroachDB database.
func ConnectToTestCockroachDB() (*sqlx.DB, error) {
	if dsn := os.Getenv("TEST_DATABASE_COCKROACHDB"); dsn != "" {
		log.Println("Found cockroachdb test database config, skipping dockertest...")
		return connect("pgx", "cockroach", dsn)
	}

	resource, err := startCockroachDB("")
	if err != nil {
		return nil, errors.Wrap(err, "Could not start resource")
	}

	db := bootstrap("postgres://root@localhost:%s/defaultdb?sslmode=disable", "26257/tcp", "pgx", pool, resource)
	return db, nil
}

// ConnectToTestCockroachDBPop connects to a test CockroachDB database.
// If a docker container is started for the database, the container be removed
// at the end of the test.
func ConnectToTestCockroachDBPop(t testing.TB) *pop.Connection {
	url := RunTestCockroachDB(t)
	return ConnectPop(t, url)
}

func bootstrap(u, port, d string, pool dockerPool, resource *dockertest.Resource) (db *sqlx.DB) {
	if err := resilience.Retry(logrusx.New("", ""), time.Second*5, time.Minute*5, func() error {
		var err error
		db, err = sqlx.Open(d, fmt.Sprintf(u, resource.GetPort(port)))
		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		if pErr := pool.Purge(resource); pErr != nil {
			log.Fatalf("Could not connect to docker and unable to remove image: %s - %s", err, pErr)
		}
		log.Fatalf("Could not connect to docker: %s", err)
	}
	return
}

var comments = regexp.MustCompile("(--[^\n]*\n)|(?s:/\\*.+\\*/)")

func StripDump(d string) string {
	d = comments.ReplaceAllLiteralString(d, "")
	d = strings.TrimPrefix(d, "Command \"dump\" is deprecated, cockroach dump will be removed in a subsequent release.\r\nFor details, see: https://github.com/cockroachdb/cockroach/issues/54040\r\n")
	d = strings.ReplaceAll(d, "\r\n", "")
	d = strings.ReplaceAll(d, "\t", " ")
	d = strings.ReplaceAll(d, "\n", " ")
	return d
}

func DumpSchema(ctx context.Context, t *testing.T, db string) string {
	var containerPort string
	var cmd []string

	switch c := stringsx.SwitchExact(db); {
	case c.AddCase("postgres"):
		containerPort = "5432"
		cmd = []string{"pg_dump", "-U", "postgres", "-s", "-T", "hydra_*_migration", "-T", "schema_migration"}
	case c.AddCase("mysql"):
		containerPort = "3306"
		cmd = []string{"/usr/bin/mysqldump", "-u", "root", "--password=secret", "mysql"}
	case c.AddCase("cockroach"):
		containerPort = "26257"
		cmd = []string{"./cockroach", "dump", "defaultdb", "--insecure", "--dump-mode=schema"}
	default:
		t.Log(c.ToUnknownCaseErr())
		t.FailNow()
		return ""
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	require.NoError(t, err)
	containers, err := cli.ContainerList(ctx, container.ListOptions{
		Filters: filters.NewArgs(filters.Arg("expose", containerPort)),
	})
	require.NoError(t, err)

	if len(containers) != 1 {
		t.Logf("Ambiguous amount of %s containers: %d", db, len(containers))
		t.FailNow()
	}

	process, err := cli.ContainerExecCreate(ctx, containers[0].ID, container.ExecOptions{
		Tty:          true,
		AttachStdout: true,
		Cmd:          cmd,
	})
	require.NoError(t, err)

	resp, err := cli.ContainerExecAttach(ctx, process.ID, container.ExecAttachOptions{
		Tty: true,
	})
	require.NoError(t, err)
	dump, err := io.ReadAll(resp.Reader)
	require.NoError(t, err, "%s", dump)

	return StripDump(string(dump))
}
