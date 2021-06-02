---
id: dependencies-environment
title: Database Setup and Configuration
---

ORY Hydra is built cloud native and implements
[12factor](https://www.12factor.net/) principles. The Docker Image is 5 MB light
and versioned with
[verbose upgrade instructions](https://github.com/ory/hydra/blob/master/UPGRADE.md)
and
[detailed changelogs](https://github.com/ory/hydra/blob/master/CHANGELOG.md).
Auto-scaling, migrations, health checks, it all works with zero additional work
required. It is possible to run ORY Hydra on any platform, including but not
limited to OSX, Linux, Windows, ARM, FreeBSD and more.

ORY Hydra has two operational modes:

- In-memory: This mode does not work with more than one instance ("cluster") and
  any state is lost after restarting the instance.
- SQL: This mode works with more than one instance and state is not lost after
  restarts.

No further dependencies are required for a production-ready instance.

## SQL

The SQL adapter supports two DBMS: PostgreSQL 9.6+ and MySQL 5.7+. Please note
that older MySQL versions have issues with ORY Hydra's database schema. For more
information [go here](https://github.com/ory/hydra/issues/377).

If you do run the SQL adapter, you must first create the database schema. The
`hydra serve` command does not do this automatically, instead you must run
`hydra migrate sql` to create the schemas. The `hydra migrate sql` command also
runs database migrations in case of an upgrade. Please follow the
[upgrade instructions](https://github.com/ory/hydra/blob/master/UPGRADE.md) to
see when you need to run this command. Always create a backup before running
`hydra migrate sql`!

Running SQL migrations in Docker is very easy, check out the
[docker-compose](https://github.com/ory/hydra/blob/master/quickstart-postgres.yml)
example to see how we did it!

### Configuration

Both MySQL and PostgreSQL adapters support the following settings. You can
modify these settings by appending query parameters to your DSN
(`postgres://user:pw@host:port/database?setting1=foo&setting2=bar`):

- `max_conns` sets the maximum number of open connections to the database.
  Defaults to the number of CPUs. Example
  `postgres://user:pw@host:port/database?max_conns=10`.
- `max_idle_conns` sets the maximum number of connections in the idle connection
  pool. Defaults to the number of CPUs. Example
  `postgres://user:pw@host:port/database?max_idle_conns=5`.
- `max_conn_lifetime` sets the maximum amount of time (`ms`, `s`, `m`, `h`) a
  connection may be reused. Defaults to 0. Example
  `postgres://user:pw@host:port/database?max_conn_lifetime=10s`.
- `max_conn_idle_time` sets the maximum amount of time (`ms`, `s`, `m`, `h`) a
  connection can be kept alive. Defaults to 0. Example
  `postgres://user:pw@host:port/database?max_conn_idle_time=10s`.

#### MySQL

> DSN Layout:
> `mysql://user:pw@tcp(host:port)/database?someSetting=value&foo=bar`.

On top of the settings above, MySQL supports additional settings:

- `sql_notes`, if set to `false`, ignores MySQL notices. If left empty or set to
  `true`, they will be treated as warnings. Example
  `mysql://user:pw@tcp(host:port)/database?sql_notes=false`.
- `sql_mode` sets the server-side strict mode. Read more about possible values
  [here](https://dev.mysql.com/doc/refman/5.7/en/sql-mode.html).

#### PostgreSQL

> DSN Layout: `postgres://user:pw@host:port/database?someSetting=value&foo=bar`.

On top of the settings above, PostgreSQL supports additional settings:

- `sslmode` sets whether or not to use SSL (default is require, this is not the
  default for libpq). Valid values for sslmode are: _ `disable` - No SSL _
  `require` - Always SSL (skip verification) _ `verify-ca` - Always SSL (verify
  that the certificate presented by the server was signed by a trusted CA) _
  `verify-full` - Always SSL (verify that the certification presented by the
  server was signed by a trusted CA and the server host name matches the one in
  the certificate)
- `fallback_application_name` - An application_name to fall back to if one isn't
  provided.
- `connect_timeout` - Maximum wait for connection, in seconds. Zero or not
  specified means wait indefinitely.
- `sslcert` - Cert file location. The file must contain PEM encoded data.
- `sslkey` - Key file location. The file must contain PEM encoded data.
- `sslrootcert` - The location of the root certificate file. The file must
  contain PEM encoded data.

## Plugin

It is possible to implement your own DBAL using Go Plugins.

> We strongly discourage you from implementing your own DBAL. Special knowledge
> is required and internal interfaces will break without notice and migration
> guides. This can lead to serious security issues and vulnerabilities. USE AT
> YOUR OWN RISK.

Your plugin must implement interface `github.com/ory/hydra/driver.Registry`. You
can load the plugin as follows:

```yaml
dsn: plugin:///path/to/file.so
```

We strongly discourage you from using this feature.
