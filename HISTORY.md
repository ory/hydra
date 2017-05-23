# History

This list makes you aware of (breaking) changes. For patch notes, please check the [releases tab](https://github.com/ory/hydra/releases).

## 0.8.0

This PR improves some performance bottlenecks, offers more control over Hydra, moves to Go 1.8,
and moves the REST documentation to swagger.

**Before applying this update, please make a back up of your database. Do not upgrade directly from versions
below 0.7.0**.

To upgrade the database schemas, please run the following commands in exactly this order

```sh
$ hydra help migrate sql
$ hydra help migrate ladon
```

```sh
$ hydra migrate sql mysql://...
$ hydra migrate ladon 0.6.0 mysql://...
```

### Breaking changes

#### Ladon updated to 0.6.0

Ladon was greatly improved with version 0.6.0, resolving various performance bottlenecks. Please read more on this
release [here](https://github.com/ory/ladon/blob/master/HISTORY.md#060).

#### Redis and RethinkDB deprecated

Redis and RethinkDB are removed from the repository now and no longer supported, see
[this issue](https://github.com/ory/hydra/issues/425).

#### Moved to ory namespace

To reflect the GitHub organization rename, Hydra was moved from `https://github.com/ory-am/hydra` to
`https://github.com/ory/hydra`.

#### SDK

The method `FindPoliciesForSubject` of the policy SDK was removed. Instead, `List` was added. The HTTP endpoint `GET /policies`
no longer allows to query by subject.

#### JWK

To generate JWKs previously the payload at `POST /keys` was `{ "alg": "...", "id": "some-id" }`. `id` was changed to
`kid` so this is now `{ "alg": "...", "kid": "some-id" }`.

#### Migrations are no longer automatically applied

SQL Migrations are no longer automatically applied. Instead you need to run `hydra migrate sql` after upgrading
to a Hydra version that includes a breaking schema change.

### Changes

#### Log format: json

Set the log format to json using `export LOG_FORMAT=json`

#### SQL Connection Control

You can configure SQL connection limits by appending parameters `max_conns`, `max_idle_conns`, or `max_conn_lifetime`
to the DSN: `postgres://foo:bar@host:port/database?max_conns=12`.

#### REST API Docs are now generated from source code

... and are swagger 2.0 spec.

#### Documentation on scopes

Documentation on scopes (e.g. offline) was added.

#### New response writer library

Hydra now uses `github.com/ory/herodot` for writing REST responses. This increases compatibility with other libraries
and resolves a few other issues.

#### Graceful http handling

Hydra is now capable of gracefully handling SIGINT.

#### Best practice HTTP server config

Hydra now implements best practices for running HTTP servers that are exposed to the public internet.
