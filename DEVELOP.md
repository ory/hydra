# Development

This document explains how to develop Ory Hydra, run tests, and work with the
tooling around it.

## Upgrading and changelog

New releases might introduce breaking changes. To help you identify and
incorporate those changes, we document these changes in
[CHANGELOG.md](./CHANGELOG.md).

## Command line documentation

To see available commands and flags, run:

```shell
hydra -h
# or
hydra help
```

## Contribution guidelines

We encourage all contributions. Before opening a pull request, read the
[contribution guidelines](./CONTRIBUTING.md).

## Prerequisites

You need Go 1.13+ with `GO111MODULE=on` and, for the test suites:

- Docker and Docker Compose
- Makefile
- Node.js and npm

It is possible to develop Ory Hydra on Windows, but please be aware that all
guides assume a Unix shell like bash or zsh.

## Formatting code

Format all code using:

```shell
make format
```

The continuous integration pipeline checks code formatting.

## Running tests

There are three types of tests:

- Short tests that do not require a SQL database
- Regular tests that require PostgreSQL, MySQL, and CockroachDB
- End to end tests that use real databases and a test browser

### Short tests

Short tests run fairly quickly and use SQLite in-memory.

All tests run against a sqlite in-memory database, thus it is required to use
the `-tags sqlite` build tag.

Run all short tests:

```shell
go test -v -failfast -short -tags sqlite ./...
```

Run short tests in a specific module:

```shell
go test -v -failfast -short -tags sqlite ./client
```

Run a specific test:

```shell
go test -v -failfast -short -tags sqlite -run ^TestName$ ./...
```

### Regular tests

Regular tests require a database setup.

The test suite can use [ory/dockertest](https://github.com/ory/dockertest) to
work with docker directly, but we encourage using the Makefile instead. Using
dockertest can bloat the number of Docker Images on your system and are quite
slow.

Run the full test suite:

```shell
make test
```

> Note: `make test` recreates the databases every time. This can be annoying if
> you are trying to fix something very specific and need the database tests all
> the time.

If you want to reuse databases across test runs, initialize them once:

```shell
make test-resetdb
export TEST_DATABASE_MYSQL='mysql://root:secret@(127.0.0.1:3444)/mysql?parseTime=true&multiStatements=true'
export TEST_DATABASE_POSTGRESQL='postgres://postgres:secret@127.0.0.1:3445/postgres?sslmode=disable'
export TEST_DATABASE_COCKROACHDB='cockroach://root@127.0.0.1:3446/defaultdb?sslmode=disable'
```

Then you can run Go tests directly as often as needed:

```shell
go test -p 1 ./...

# or in a module:
cd client
go test .
```

### End-to-end tests

The E2E tests use [Cypress](https://www.cypress.io) to run full browser tests.

Run e2e tests:

```
make e2e
```

The runner will not show the Browser window, as it runs in the CI Mode
(background). That makes debugging these type of tests very difficult, but
thankfully you can run the e2e test in the browser which helps with debugging:

```shell
./test/e2e/circle-ci.bash memory --watch

# Or for the JSON Web Token Access Token strategy:
# ./test/e2e/circle-ci.bash memory-jwt --watch
```

Or if you would like to test one of the databases:

```shell
make test-resetdb
export TEST_DATABASE_MYSQL='mysql://root:secret@(127.0.0.1:3444)/mysql?parseTime=true&multiStatements=true'
export TEST_DATABASE_POSTGRESQL='postgres://postgres:secret@127.0.0.1:3445/postgres?sslmode=disable'
export TEST_DATABASE_COCKROACHDB='cockroach://root@127.0.0.1:3446/defaultdb?sslmode=disable'

# You can test against each individual database:
./test/e2e/circle-ci.bash postgres --watch
./test/e2e/circle-ci.bash memory --watch
./test/e2e/circle-ci.bash mysql --watch
# ...
```

Once you run the script, a Cypress window will appear. Hit the button "Run all
Specs"!

The code for these tests is located in
[./cypress/integration](./cypress/integration) and
[./cypress/support](./cypress/support) and
[./cypress/helpers](./cypress/helpers). The website you're seeing is located in
[./test/e2e/oauth2-client](./test/e2e/oauth2-client).

#### OpenID Connect conformity tests

To run Ory Hydra against the OpenID Connect conformity suite, run:

```shell
./test/conformity/start.sh --build
```

and then in a separate shell:

```shell
./test/conformity/test.sh
```

Running these tests will take a significant amount of time which is why they are
not part of the CI pipeline.

## Build Docker image

To build a development Docker Image:

```shell
make docker
```

> [!WARNING] If you already have a production image (e.g. `oryd/hydra:v2.2.0`)
> pulled, the above `make docker` command will replace it with a local build of
> the image that is more equivalent to the `-distroless` variant on Docker Hub.
>
> You can pull the production image any time using `docker pull`

## Run the Docker Compose quickstarts

If you wish to check your code changes against any of the docker-compose
quickstart files, run:

```shell
docker compose -f quickstart.yml up --build
```

## Add a new migration

1. `mkdir persistence/sql/src/YYYYMMDD000001_migration_name/`
2. Put the migration files into this directory, following the standard naming
   conventions. If you wish to execute different parts of a migration in
   separate transactions, add split marks (lines with the text `--split`) where
   desired. Why this might be necessary is explained in
   https://github.com/gobuffalo/fizz/issues/104.
3. Run `make persistence/sql/migrations/<migration_id>` to generate migration
   fragments.
4. If an update causes the migration to have fewer fragments than the number
   already generated, run
   `make persistence/sql/migrations/<migration_id>-clean`. This is equivalent to
   a `rm` command with the right parameters, but comes with better tab
   completion.
5. Before committing generated migration fragments, run the above clean command
   and generate a fresh copy of migration fragments to make sure the `sql/src`
   and `sql/migrations` directories are consistent.
