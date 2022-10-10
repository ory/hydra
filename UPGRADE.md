# Upgrading

The intent of this document is to make migration of breaking changes as easy as
possible. Please note that not all breaking changes might be included here.
Please check the [CHANGELOG.md](./CHANGELOG.md) for a full list of changes
before finalizing the upgrade process.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Hassle-free upgrades](#hassle-free-upgrades)
- [1.4](#14)
- [1.3.0](#130)
- [1.2.0](#120)
- [1.1.0](#110)
- [1.0.9](#109)
  - [Schema Changes](#schema-changes)
- [1.0.0-rc.10](#100-rc10)
  - [OpenID Connect Front-/Backchannel Logout 1.0](#openid-connect-front-backchannel-logout-10)
  - [Schema Changes](#schema-changes-1)
  - [SQL Migrations now require user-input or `--yes` flag](#sql-migrations-now-require-user-input-or---yes-flag)
  - [Login and Consent Management](#login-and-consent-management)
- [1.0.0-rc.9](#100-rc9)
  - [Go SDK](#go-sdk)
  - [Accepting Login and Consent Requests](#accepting-login-and-consent-requests)
- [1.0.0-rc.7](#100-rc7)
  - [Configuration changes](#configuration-changes)
  - [System secret rotation](#system-secret-rotation)
  - [Database Plugins](#database-plugins)
- [1.0.0-rc.4](#100-rc4)
- [1.0.0-rc.1](#100-rc1)
  - [Schema Changes](#schema-changes-2)
    - [Foreign Keys](#foreign-keys)
      - [Removing inconsistent oauth2 data](#removing-inconsistent-oauth2-data)
      - [Removing inconsistent login & consent data](#removing-inconsistent-login--consent-data)
    - [Indices](#indices)
  - [Non-breaking Changes](#non-breaking-changes)
    - [Access Token Audience](#access-token-audience)
    - [Refresh Grant](#refresh-grant)
    - [Customise login and consent flow timeout](#customise-login-and-consent-flow-timeout)
  - [Breaking Changes](#breaking-changes)
    - [Refresh Token Expiry](#refresh-token-expiry)
    - [Swagger & SDK Restructuring](#swagger--sdk-restructuring)
      - [Go](#go)
      - [Others](#others)
    - [JSON Web Token formatted Access Token data](#json-web-token-formatted-access-token-data)
    - [CLI Changes](#cli-changes)
    - [API Changes](#api-changes)
- [1.0.0-beta.9](#100-beta9)
  - [CORS is disabled by default](#cors-is-disabled-by-default)
- [1.0.0-beta.8](#100-beta8)
  - [Schema Changes](#schema-changes-3)
  - [Split of Public and Administrative Endpoints](#split-of-public-and-administrative-endpoints)
  - [Golang SDK `Configuration.EndpointURL` is now `Configuration.AdminURL`](#golang-sdk-configurationendpointurl-is-now-configurationadminurl)
  - [`hydra serve` is now `hydra serve all`](#hydra-serve-is-now-hydra-serve-all)
  - [Environment variable `HYDRA_URL` now is `HYDRA_ADMIN_URL` for admin commands](#environment-variable-hydra_url-now-is-hydra_admin_url-for-admin-commands)
  - [OAuth 2.0 Token Introspection](#oauth-20-token-introspection)
  - [OAuth 2.0 Client flag `public` has been removed](#oauth-20-client-flag-public-has-been-removed)
- [1.0.0-beta.7](#100-beta7)
  - [Regenerated OpenID Connect ID Token cryptographic keys](#regenerated-openid-connect-id-token-cryptographic-keys)
- [1.0.0-beta.5](#100-beta5)
  - [OAuth 2.0 Client Response Type changes](#oauth-20-client-response-type-changes)
  - [Schema Changes](#schema-changes-4)
  - [HTTP Error Payload](#http-error-payload)
  - [OAuth 2.0 Clients must specify correct `token_endpoint_auth_method`](#oauth-20-clients-must-specify-correct-token_endpoint_auth_method)
  - [OAuth 2.0 Client field `id` is now `client_id`](#oauth-20-client-field-id-is-now-client_id)
- [1.0.0-beta.1](#100-beta1)
  - [Upgrading from versions v0.9.x](#upgrading-from-versions-v09x)
  - [OpenID Connect Certified](#openid-connect-certified)
  - [Breaking Changes](#breaking-changes-1)
    - [Introspection API](#introspection-api)
      - [Introspection is now capable of introspecting refresh tokens](#introspection-is-now-capable-of-introspecting-refresh-tokens)
    - [Access Control & Warden API](#access-control--warden-api)
      - [Running the backwards compatible set up](#running-the-backwards-compatible-set-up)
        - [Warden API](#warden-api)
        - [Warden Groups](#warden-groups)
    - [jwk: Forces JWK to have a unique ID](#jwk-forces-jwk-to-have-a-unique-id)
    - [Consent Flow](#consent-flow)
    - [Changes to the CLI](#changes-to-the-cli)
      - [`hydra host`](#hydra-host)
      - [`hydra connect`](#hydra-connect)
      - [`hydra token user`](#hydra-token-user)
      - [`hydra token client`](#hydra-token-client)
      - [`hydra token validate`](#hydra-token-validate)
      - [`hydra clients create`](#hydra-clients-create)
      - [`hydra migrate ladon`](#hydra-migrate-ladon)
      - [`hydra policies`](#hydra-policies)
      - [`hydra groups`](#hydra-groups)
    - [SDK](#sdk)
  - [Improvements](#improvements)
    - [Health Check endpoint has moved](#health-check-endpoint-has-moved)
    - [Unknown request body payloads result in error](#unknown-request-body-payloads-result-in-error)
    - [UTC everywhere](#utc-everywhere)
    - [Pagination everywhere](#pagination-everywhere)
    - [Flushing old access tokens](#flushing-old-access-tokens)
    - [Prometheus endpoint](#prometheus-endpoint)
- [0.11.12](#01112)
- [0.11.3](#0113)
- [0.11.0](#0110)
- [0.10.0](#0100)
  - [Breaking Changes](#breaking-changes-2)
    - [Introspection now requires authorization](#introspection-now-requires-authorization)
    - [New consent flow](#new-consent-flow)
    - [Audience](#audience)
    - [Response payload changes to `/warden/token/allowed`](#response-payload-changes-to-wardentokenallowed)
    - [Go SDK](#go-sdk-1)
    - [Health endpoints](#health-endpoints)
    - [Group endpoints](#group-endpoints)
    - [Replacing hierarchical scope strategy with wildcard scope strategy](#replacing-hierarchical-scope-strategy-with-wildcard-scope-strategy)
    - [AES-GCM nonce storage](#aes-gcm-nonce-storage)
    - [Minor Breaking Changes](#minor-breaking-changes)
      - [Token signature algorithm changed from HMAC-SHA256 to HMAC-SHA512](#token-signature-algorithm-changed-from-hmac-sha256-to-hmac-sha512)
      - [HS256 JWK Generator now uses all 256 bit](#hs256-jwk-generator-now-uses-all-256-bit)
      - [ES512 Key generator](#es512-key-generator)
      - [Build tags deprecated](#build-tags-deprecated)
  - [Important Additions](#important-additions)
    - [Prefixing Resources Names](#prefixing-resources-names)
    - [Refreshing OpenID Connect ID Token using `refresh_token` grant type](#refreshing-openid-connect-id-token-using-refresh_token-grant-type)
  - [Important Changes](#important-changes)
    - [Telemetry](#telemetry)
    - [URL Encoding Root Client Credentials](#url-encoding-root-client-credentials)
- [0.9.0](#090)
- [0.8.0](#080)
  - [Breaking changes](#breaking-changes)
    - [Ladon updated to 0.6.0](#ladon-updated-to-060)
    - [Redis and RethinkDB deprecated](#redis-and-rethinkdb-deprecated)
    - [Moved to ory namespace](#moved-to-ory-namespace)
    - [SDK](#sdk-1)
    - [JWK](#jwk)
    - [Migrations are no longer automatically applied](#migrations-are-no-longer-automatically-applied)
  - [Changes](#changes)
    - [Log format: json](#log-format-json)
    - [SQL Connection Control](#sql-connection-control)
    - [REST API Docs are now generated from source code](#rest-api-docs-are-now-generated-from-source-code)
    - [Documentation on scopes](#documentation-on-scopes)
    - [New response writer library](#new-response-writer-library)
    - [Graceful http handling](#graceful-http-handling)
    - [Best practice HTTP server config](#best-practice-http-server-config)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Hassle-free upgrades

Do you want the latest features and patches without work and hassle? Are you
looking for a reliable, scalable, and secure deployment with zero effort? We can
run it for you! If you're interested, [contact us now](mailto:office@ory.sh)!

## 1.5.0

Migrations are now handled with https://github.com/gobuffalo/fizz. Please run
`hydra migrate sql` when upgrading to this version.

For a full list of changes please check:
https://github.com/ory/hydra/compare/v1.4...v1.5.0

## 1.4

Please run `hydra migrate sql` when upgrading to this version. For more
information, check
https://github.com/ory/hydra/commit/700d17d3b7d507de1b1d459a7261d6fb2571ebe3.

For a full list of changes please check:
https://github.com/ory/hydra/compare/v1.3.2...v1.4.1

## 1.3.0

Please run `hydra migrate sql` when upgrading to this version. For more
information, check
https://github.com/ory/hydra/commit/d9308fa0dba26019a59e4d97e85b036133ad8362.

## 1.2.0

This release focuses on a rework of the SDK pipeline. First of all, we have
introduced new SDKs for all popular programming languages and published them on
their respective package repositories:

- [Python](https://pypi.org/project/ory-hydra-client/)
- [PHP](https://packagist.org/packages/ory/hydra-client)
- [Go](https://github.com/ory/hydra-client-go)
- [NodeJS](https://www.npmjs.com/package/@oryd/hydra-client) (with TypeScript)
- [Java](https://search.maven.org/artifact/sh.ory.hydra/hydra-client)
- [Ruby](https://rubygems.org/gems/ory-hydra-client)

The SDKs hosted in this repository (under ./sdk/...) have been completely
removed. Please use only the SDKs from the above sources from now on as it will
also remove several issues that were caused by the previous SDK pipeline.

Unfortunately, there were breaking changes introduced by the new SDK generation:

- Several structs and fields have been renamed in the Go SDK. However, nothing
  else changed so upgrading should be a matter of half an hour if you made
  extensive use of the SDK, or several minutes if just one or two methods are
  being used.
- All other SDKs changed to `openapi-generator`, which is a better maintained
  generator that creates better code than the one previously used. This
  manifests in TypeScript definitions for the NodeJS SDK and several other
  goodies. We do not have a proper migration path for those, unfortunately.

If you have issues with upgrading the SDK, please let us know in an issue on
this repository!

## 1.1.0

Several indices have been added to the SQL Migrations. There are no backwards
incompatible changes in this release but we advise to do a test-run of the SQL
Migrations before applying them, as they might lock some tables which may cause
downtimes.

> Make a backup of your database before applying this change.

After applying these SQL Migrations, several queries and endpoints will be much
faster than before.

## 1.0.9

### Schema Changes

A minor Schema change was introduced to the OAuth 2.0 Clients table. It is now
possible to store arbitrary metadata for a client.

> Make a backup of your database before applying this change.

## 1.0.0-rc.10

### OpenID Connect Front-/Backchannel Logout 1.0

This patch implements OpenID Connect Front-/Backchannel Logout 1.0
([read docs](https://www.ory.sh/docs/hydra/oauth2#logout)). Therefore, endpoint
`/oauth2/auth/sessions/login/revoke` has been deprecated.

### Schema Changes

Please read all paragraphs of this section with the utmost care, before
executing `hydra migrate sql`. Do not take this change lightly and create a
backup of the database before you begin. To be sure, copy the database and do a
dry-run locally.

> Be aware that running these migrations might take some time when using large
> databases. Do a dry-run before hammering your production database.

### SQL Migrations now require user-input or `--yes` flag

`hydra migrate sql` now shows an execution plan and asks for confirmation before
executing the migrations. To run migrations without user interaction, add flag
`--yes`.

### Login and Consent Management

Orthogonal to the changes when accepting and rejection consent and login
requests, the following endpoints have been updated as well:

- `DELETE /oauth2/auth/sessions/login/:subject` ->
  `DELETE /oauth2/auth/sessions/login?subject={subject}`
- `GET /oauth2/auth/sessions/consent/:subject` ->
  `GET /oauth2/auth/sessions/login?subject={subject}`
- `DELETE /oauth2/auth/sessions/consent/:subject` ->
  `DELETE /oauth2/auth/sessions/login?subject={subject}`
- `DELETE /oauth2/auth/sessions/consent/:subject/:client` ->
  `DELETE /oauth2/auth/sessions/login?subject={subject}&client={client}`

While this does not include a security warning, this patch allows developers to
use slashes in dots in their subject/user IDs.

## 1.0.0-rc.9

### Go SDK

The Go SDK is now being generated using `go-swagger`. The SDK generated using
`swagger-codegen` is no longer supported. The old Go SDK is still available but
moved to a new path. To use it, change:

```
- import "github.com/ory/hydra/sdk/go/hydra"
- import "github.com/ory/hydra/sdk/go/hydra/swagger"

+ import hydra "github.com/ory/hydra-legacy-sdk"
+ import "github.com/ory/hydra-legacy-sdk/swagger"
```

### Accepting Login and Consent Requests

Previously, login and consent requests were accepted/rejected by doing one of:

```
GET /oauth2/auth/requests/login/{challenge}
PUT /oauth2/auth/requests/login/{challenge}/accept
PUT /oauth2/auth/requests/login/{challenge}/reject

GET /oauth2/auth/requests/consent/{challenge}
PUT /oauth2/auth/requests/consent/{challenge}/accept
PUT /oauth2/auth/requests/consent/{challenge}/reject
```

We observed login/consent apps that did not properly sanitize the `{challenge}`
parameter, making it possible to escape the path by using `..` in the challenge
parameter (e.g. `http://my-login-app/login?challenge=../../whatever`) causing
the login/consent app to execute a request it is not supposed to be making (e.g.
`/oauth2/auth/requests/login/../../whatever/accept`).

From now on, the challenge has to be sent using a query parameter instead:

```
GET /oauth2/auth/requests/login?challenge={challenge}
PUT /oauth2/auth/requests/login/accept?challenge={challenge}
PUT /oauth2/auth/requests/login/reject?challenge={challenge}

GET /oauth2/auth/requests/consent?challenge={challenge}
PUT /oauth2/auth/requests/consent/accept?challenge={challenge}
PUT /oauth2/auth/requests/consent/reject?challenge={challenge}
```

Implementers will still need to make sure that `challenge` is properly (query)
scaped, but it's generally easier to secure than a path parameter.

We've decided to make this a hard breaking change in order to force everybody to
check if their application is vulnerable to this issue and to upgrade their
code. The required code change is minimal but the resulting security
improvements are potentially large.

## 1.0.0-rc.7

### Configuration changes

This patch introduces changes to the way configuration works in ORY Hydra. It
allows ORY Hydra to be configured from a variety of sources including
environment variables and a configuration file. In the future, ORY Hydra might
be configurable using etcd or consul. The changes allow ORY Hydra to reload
configuration without restarting in the future.

An overview of configuration settings can be found
[here](https://github.com/ory/hydra/blob/master/docs/config.yaml).

All changes are backwards compatible except for the way key rotation works (see
next section) and the way DBAL plugins are loaded (see section after next).

### System secret rotation

Rotating system secrets was fairly cumbersome in the past and required a restart
of ORY Hydra. This changed. The system secret is now an array where the first
element is used for encryption and all elements can be used for decryption.

For more information on this topic, click
[here](https://www.ory.sh/docs/hydra/advanced#rotation-of-hmac-token-signing-and-database-and-cookie-encryption-keys).

To make this change work, environment variable `ROTATED_SYSTEM_SECRET` has been
removed and can no longer be used. Command `hydra migrate secret` has also been
removed without replacement as it is no longer required for rotating secrets.

### Database Plugins

Environment variable `DATABASE_PLUGIN` has been replaced by `dsn`. To load a
plugin, set `dsn: plugin:///path/to/plugin.so`.

Please note that internals have changed radically with this patch and that there
is some refactoring effort required to make plugins work with the most recent
version.

## 1.0.0-rc.4

This patch requires you to run SQL migrations. No other important changes have
been made.

## 1.0.0-rc.1

This release ships with major scalability and reliability improvements and
resolves several bugs.

### Schema Changes

Please read all paragraphs of this section with the utmost care, before
executing `hydra migrate sql`. Do not take this change lightly and create a
backup of the database before you begin. To be sure, copy the database and do a
dry-run locally.

> Be aware that running these migrations might take some time when using large
> databases. Do a dry-run before hammering your production database.

#### Foreign Keys

In order to keep data consistent across tables, several foreign key constraints
have been added between consent, oauth2, client tables. If you are running a
large database take enough time to run this migration - it might take a while
depending on the amount of data and the database version and driver. Before
executing this migration, you should _manually_ check and remove inconsistent
data.

##### Removing inconsistent oauth2 data

This migration automatically removes inconsistent OAuth 2.0 and OpenID Connect
data. Possible impacts are:

1. Existing authorize codes, access, refresh tokens might be invalidated (all
   flows, including PKCE and OpenID Connect)

As OAuth 2.0 clients are generally capable of handling re-authorization, this
should not have a serious impact. Removing this data increases security through
strong consistency. The following data-altering statements will be executed:

```sql
-- First we need to delete all rows that point to a non-existing oauth2 client.
DELETE FROM hydra_oauth2_access WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_access.client_id = hydra_client.id);
DELETE FROM hydra_oauth2_refresh WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_refresh.client_id = hydra_client.id);
DELETE FROM hydra_oauth2_code WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_code.client_id = hydra_client.id);
DELETE FROM hydra_oauth2_oidc WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_oidc.client_id = hydra_client.id);
DELETE FROM hydra_oauth2_pkce WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_pkce.client_id = hydra_client.id);

-- request_id is a 40 varchar in the referenced table which is why we are resizing
-- 1. We must remove request_ids longer than 40 chars. This should never happen as we've never issued them longer than this
DELETE FROM hydra_oauth2_access WHERE LENGTH(request_id) > 40;
DELETE FROM hydra_oauth2_refresh WHERE LENGTH(request_id) > 40;
DELETE FROM hydra_oauth2_code WHERE LENGTH(request_id) > 40;
DELETE FROM hydra_oauth2_oidc WHERE LENGTH(request_id) > 40;
DELETE FROM hydra_oauth2_pkce WHERE LENGTH(request_id) > 40;

-- 2. Next we're actually resizing
ALTER TABLE hydra_oauth2_access ALTER COLUMN request_id TYPE varchar(40);
ALTER TABLE hydra_oauth2_refresh ALTER COLUMN request_id TYPE varchar(40);
ALTER TABLE hydra_oauth2_code ALTER COLUMN request_id TYPE varchar(40);
ALTER TABLE hydra_oauth2_oidc ALTER COLUMN request_id TYPE varchar(40);
ALTER TABLE hydra_oauth2_pkce ALTER COLUMN request_id TYPE varchar(40);

-- In preparation for creating the client_id index and foreign key, we must set it to varchar(255) which is also
-- the length of hydra_client.id
DELETE FROM hydra_oauth2_access WHERE LENGTH(client_id) > 255;
DELETE FROM hydra_oauth2_refresh WHERE LENGTH(client_id) > 255;
DELETE FROM hydra_oauth2_code WHERE LENGTH(client_id) > 255;
DELETE FROM hydra_oauth2_oidc WHERE LENGTH(client_id) > 255;
DELETE FROM hydra_oauth2_pkce WHERE LENGTH(client_id) > 255;
ALTER TABLE hydra_oauth2_access ALTER COLUMN client_id TYPE varchar(255);
ALTER TABLE hydra_oauth2_refresh ALTER COLUMN client_id TYPE varchar(255);
ALTER TABLE hydra_oauth2_code ALTER COLUMN client_id TYPE varchar(255);
ALTER TABLE hydra_oauth2_oidc ALTER COLUMN client_id TYPE varchar(255);
ALTER TABLE hydra_oauth2_pkce ALTER COLUMN client_id TYPE varchar(255);
```

##### Removing inconsistent login & consent data

This migration automatically removes inconsistent login & consent data. Possible
impacts are:

1. Users that set `remember` to true during login have to re-authenticate.
2. Users that set `remember` to true during consent have to re-authorize
   requested OAuth 2.0 Scope.
3. Data associated with OAuth 2.0 Clients that have been removed will be
   deleted.

That is achieved by running the following queries. Make sure you understand what
these queries do and what impact they may have on your system before executing
`hydra migrate sql`:

```sql
-- This can be null when no previous login session exists, so let's remove default
ALTER TABLE hydra_oauth2_authentication_request ALTER COLUMN login_session_id DROP DEFAULT;

-- This can be null when no previous login session exists or if that session has been removed, so let's remove default
ALTER TABLE hydra_oauth2_consent_request ALTER COLUMN login_session_id DROP DEFAULT;

-- This can be null when the login_challenge was deleted (should not delete the consent itself)
ALTER TABLE hydra_oauth2_consent_request ALTER COLUMN login_challenge DROP DEFAULT;

-- Consent requests that point to an empty or invalid login request should set their login_challenge to NULL
UPDATE hydra_oauth2_consent_request SET login_challenge = NULL WHERE NOT EXISTS (
  SELECT 1 FROM hydra_oauth2_authentication_request WHERE hydra_oauth2_consent_request.login_challenge = hydra_oauth2_authentication_request.challenge
);

-- Consent requests that point to an empty or invalid login session should set their login_session_id to NULL
UPDATE hydra_oauth2_consent_request SET login_session_id = NULL WHERE NOT EXISTS (
  SELECT 1 FROM hydra_oauth2_authentication_session WHERE hydra_oauth2_consent_request.login_session_id = hydra_oauth2_authentication_session.id
);

-- Login requests that point to a login session that no longer exists (or was never set in the first place) should set that to NULL
UPDATE hydra_oauth2_authentication_request SET login_session_id = NULL WHERE NOT EXISTS (
  SELECT 1 FROM hydra_oauth2_authentication_session WHERE hydra_oauth2_authentication_request.login_session_id = hydra_oauth2_authentication_session.id
);

-- Login, consent, obfuscated sessions that point to a client which no longer exists must be deleted
DELETE FROM hydra_oauth2_authentication_request WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_authentication_request.client_id = hydra_client.id);
DELETE FROM hydra_oauth2_consent_request WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_consent_request.client_id = hydra_client.id);
DELETE FROM hydra_oauth2_obfuscated_authentication_session WHERE NOT EXISTS (SELECT 1 FROM hydra_client WHERE hydra_oauth2_obfuscated_authentication_session.client_id = hydra_client.id);

-- Handled login and consent requests which point to a consent/login request that no longer exists must be deleted
DELETE FROM hydra_oauth2_consent_request_handled WHERE NOT EXISTS (SELECT 1 FROM hydra_oauth2_consent_request WHERE hydra_oauth2_consent_request_handled.challenge = hydra_oauth2_consent_request.challenge);
DELETE FROM hydra_oauth2_authentication_request_handled WHERE NOT EXISTS (SELECT 1 FROM hydra_oauth2_consent_request WHERE hydra_oauth2_authentication_request_handled.challenge = hydra_oauth2_consent_request.challenge);
```

Be aware that some queries might cascade and remove other data to. One such
example is checking `hydra_oauth2_consent_request` for rows that have no
associated `login_challenge`. If such a row is removed, the associated
`hydra_oauth2_consent_request_handled` is removed as well.

#### Indices

Several indices have been added which should resolve table locking when
searching in large data sets.

### Non-breaking Changes

#### Access Token Audience

This patch adds the access token audience feature. For more information on this,
head over to [the docs](https://www.ory.sh/docs/hydra/advanced).

#### Refresh Grant

Previously, the refresh grant did not check whether a client's allowed scope or
audience changed. This has now been added. If an OAuth 2.0 Client performs the
refresh flow but the requested token includes a scope which has not been
whitelisted at the client, the flow will fail and no refresh token will be
granted.

#### Customise login and consent flow timeout

You can now set the login and consent flow timeout using environment variable
`LOGIN_CONSENT_REQUEST_LIFESPAN`.

### Breaking Changes

#### Refresh Token Expiry

All refresh tokens issued with this release will expire after 30 days of
non-use. This behaviour can be modified using the `REFRESH_TOKEN_LIFESPAN`
environment variable. By setting `REFRESH_TOKEN_LIFESPAN=-1`, refresh tokens are
set to never expire, which is the previous behaviour.

Tokens issued before this change will still be valid forever.

We discourage setting `REFRESH_TOKEN_LIFESPAN=-1` as it might clog the database
with tokens that will never be used again. In high-scale systems,
`REFRESH_TOKEN_LIFESPAN` should be set to something like 15 or 30 days.

#### Swagger & SDK Restructuring

To better represent the public and admin endpoint, previous swagger tags (like
oAuth2, jwks, ...) have been deprecated in favor of tags `public` and `admin`.
This has different impacts for the different code-generated client libraries.

##### Go

If you use the `hydra.SDK` interface only and the `hydra.NewSDK()` factory,
everything will work as before. If you rely on e.g. `hydra.Ne.OAuth2Api.)`, you
will be affected by this change.

##### Others

All method signatures stayed the same, but the factory names for instantiating
the SDK client have changed. For example, `hydra.Ne.OAuth2Api.)` is now
`hydra.NewAdminApi()` and `hydra.NewPublicApi()` - depending on which endpoints
you need to interact with.

#### JSON Web Token formatted Access Token data

Previously, extra fields coming from `session.access_token` where directly
embedded in the OAuth 2.0 Access Token when the JSON Web Token strategy was
used. However, the token introspection response returned the extra data as a
field `ext: {...}`.

In order to have a streamlined experience, session data is from now on stored in
a field `ext: {...}` for Access Tokens formatted as JSON Web Tokens.

This change does not impact the opaque strategy, which is the default one.

#### CLI Changes

Flags `https-tls-key-path` and `https-tls-cert-path` have been removed from the
`hydra serve *` commands. Use environment variables `HTTPS_TLS_CERT_PATH` and
`HTTPS_TLS_KEY_PATH` instead.

#### API Changes

Endpoint `/health/status`, which redirected to `/health/alive` was deprecated
and has been removed.

## 1.0.0-beta.9

### CORS is disabled by default

A new environment variable `CORS_ENABLED` was introduced. It sets whether CORS
is enabled ("true") or not ("false")". Default is disabled.

## 1.0.0-beta.8

### Schema Changes

This patch introduces some minor database schema changes. Before you apply it,
you must run `hydra migrate sql` against your database.

### Split of Public and Administrative Endpoints

Previously, all endpoints were exposed at one port. Since access control was
removed with version 1.0.0, administrative endpoints (JWKs management, OAuth 2.0
Client Management, Login & Consent Management) were exposed and had to be
secured with sophisticated set ups using, for example, an API gateway to control
which endpoints can be accessed by whom.

This version introduces a new port (default `:4445`, configurable using
environment variables `ADMIN_PORT` and `ADMIN_POST`) which is serves all
administrative APIs:

- All `/clients` endpoints.
- All `/jwks` endpoints.
- All `/health`, `/metrics`, `/version` endpoints.
- All `/oauth2/auth/requests` endpoints.
- Endpoint `/oauth2/introspect`.
- Endpoint `/oauth2/flush`.

The second port exposes API endpoints generally available to the public (default
`:4444`, configurable using environment variables `PUBLIC_PORT` and
`PUBLIC_HOST`):

- `./well-known/jwks.json`
- `./well-known/openid-configuration`
- `/oauth2/auth`
- `/oauth2/token`
- `/oauth2/revoke`
- `/oauth2/fallbacks/consent`
- `/oauth2/fallbacks/error`
- `/userinfo`

The simplest way to starting both ports is to run `hydra serve`. This will start
a process which listens on both ports and exposes their respective features. All
settings (cors, database, tls, ...) will be shared by both listeners.

To configure each listener differently - for example setting CORS for public but
not privileged APIs - you can run `hydra serve public` and `hydra serve admin`
with different settings. Be aware that this will not work with `DATABASE=memory`
and that both services must use the same secrets.

### Golang SDK `Configuration.EndpointURL` is now `Configuration.AdminURL`

To reflect the changes made in this patch, the SDK's configuration struct has
been updated. Additionally, `Configuration.PublicURL` has been added in case you
need to perform OAuth2 flows with the SDK before accessing the admin endpoints.

### `hydra serve` is now `hydra serve all`

To reflect the changes of public and administrative ports, command `hydra serve`
is now `hydra serve all`.

### Environment variable `HYDRA_URL` now is `HYDRA_ADMIN_URL` for admin commands

CLI Commands like `hydra clients ...`, `hydra keys ...`, `hydra token flush`,
`hydra token introspect` no longer use environment variable `HYDRA_URL` as
default for `--endpoint` but instead `HYDRA_ADMIN_URL`.

### OAuth 2.0 Token Introspection

Previously, OAuth 2.0 Token Introspection was protected with HTTP Basic
Authorization (a valid OAuth 2.0 Client with Client ID and Client Secret was
needed) or HTTP Bearer Authorization (a valid OAuth 2.0 Access Token was
needed).

As OAuth 2.0 Token Introspection is generally an internal-facing endpoint used
by resource servers to validate OAuth 2.0 Access Tokens, this endpoint has moved
to the privileged port. The specification does not implore which authorization
scheme must be used - it only shows that HTTP Basic/Bearer Authorization may be
used. By exposing this endpoint to the privileged port a strong authorization
scheme is implemented and no further authorization is needed. Thus, access
control was stripped from this endpoint, making integration with other API
gateways easier.

You may still choose to export this endpoint to the public internet and
implement any access control mechanism you find appropriate.

### OAuth 2.0 Client flag `public` has been removed

Previously, OAuth 2.0 Clients had a flag called `public`. If set to true, the
OAuth 2.0 Client was able to exchange authorize codes for access tokens without
a password. This is useful in scenarios where the device can not keep a secret
(browser app, mobile app).

Since OpenID Connect Dynamic Discovery was added, this flag collided with the
`token_endpoint_auth_method`. If `token_endpoint_auth_method` is set to `none`,
then that is equal to setting `public` to `true`. To remove this ambiguity the
`public` flag was removed.

If you wish to create a client that runs on an untrusted device (browser app,
mobile app), simply set `"token_endpoint_auth_method": "none"` in the JSON
request.

If you are using the ORY Hydra CLI, you can use
`--token-endpoint-auth-method none` to achieve what `--is-public` did
previously.

The SQL migrations will automatically migrate clients that have `public` set to
`true` by setting `token_endpoint_auth_method` to `none`.

## 1.0.0-beta.7

### Regenerated OpenID Connect ID Token cryptographic keys

This patch resolves an issue which caused the migration to fail from beta.4 to
beta.5 / beta.6. The reason being that the keys stored in the data store had
mismatching `kid` values if generated by <= beta.5. This patch runs a SQL
migration script which removes the old key and then, after booting up ORY Hydra,
regenerates it.

To apply this change, please run you must run `hydra migrate sql` against your
database.

## 1.0.0-beta.5

This patch implements the OpenID Connect Dynamic Client registration
specification and thus now supports client authentication via JSON Web Tokens
signed with RSA public/private keypairs, alongside HTTP Basic Authorization and
sending the client's ID and secret in the POST body.

For more information on this, please refer to the
[specification](http://openid.net/specs/openid-connect-core-1_0.html#ClientAuthentication).

### OAuth 2.0 Client Response Type changes

Previously, when response types such as `code token id_token` were requested
(OpenID Connect Hybrid Flow) it was enough for the client to have
`response_types=["code", "token", "id_token"]`. This is however incompatible
with the OpenID Connect Dynamic Client Registration 1.0 spec which dictates that
the `response_types` have to match exactly.

Assuming you are requesting `&response_types=code+token+id_token`, your client
should have `response_types=["code token id_token"]`, if other response types
are required (e.g. `&response_types=code`, `&response_types=token`) they too
must be included: `response_types=["code", "token", "code token id_token"]`.

This will only affect you if you have clients requesting OpenID Connect Hybrid
flows where more than one response_type is requested.

### Schema Changes

This patch introduces some minor database schema changes. Before you apply it,
you must run `hydra migrate sql` against your database.

### HTTP Error Payload

Previously, errors have been returned as nested objects:

```json
{
  "error": {
    "error": "invalid_request"
    // ...
  }
}
```

while other endpoints, specifically those under OAuth 2.0 / OpenID Connect
returned them without nesting:

```json
{
  "error": "invalid_request"
  // ...
}
```

This patch updates all error responses and formats them coherently as across all
APIs:

```json
{
  "error": "invalid_request"
  // ...
}
```

### OAuth 2.0 Clients must specify correct `token_endpoint_auth_method`

With support for the OpenID Connect Dynamic Discovery specification, a new field
has been added to the OAuth 2.0 Client's metadata which is
`token_endpoint_auth_method`. The `token_endpoint_auth_method` specifies which
authentication methods the client can use at the token, introspection, and
revocation endpoint.

The default value for this method is `client_secret_basic` which uses the Basic
HTTP Authorization scheme. If your client uses the POST body to perform
authentication, this value must be changed to `client_secret_post`

### OAuth 2.0 Client field `id` is now `client_id`

The
[OpenID Connect Dynamic Client Registration 1.0](https://openid.net/specs/openid-connect-registration-1_0.html)
spec formulates that the client's ID should be sent as field `client_id`. Until
now, the id was sent as field `id`. This release changes that. For example, what
was previously

```
$ curl http://hydra/clients/my-client

{
    "id": "my-client",
    // ...
}
```

is now

```
$ curl http://hydra/clients/my-client

{
    "client_id": "my-client",
    // ...
}
```

For now, the `id` field will still be returned, but is marked deprecated and
will be removed in future releases.

## 1.0.0-beta.1

This section summarizes important changes introduced in 1.0.0. **Follow it
chronologically to ensure a proper migration.**

We are very well aware that the changelist is huge and we try to prepare you as
good as we can to migrate to this version. We also understand that breaking
changes are frustrating and it takes time to adopt to them. We sincerely hope
that the benefits from this version (improved consent flow, easier set up, clear
boundaries & responsibilities) outweigh the hassle of upgrading to the new
version.

If you have difficulties upgrading and would like a helping hand, reach out to
us at [mailto:hi@ory.sh](hi@ory.sh) and we will help you with the upgrade
process. Our services are billed by the hour and are priced fairly.

### Upgrading from versions v0.9.x

This is a (potentially incomplete) summary of what needs to be done to upgrade
from versions of the 0.9.x branch, which are still common. As always, try this
with a staging environment first and create back ups. Never run this directly in
production unless you are 100% sure everything works:

1. Get the latest ORY Hydra binary or docker image from the 1.0.0 branch.
2. Run `$ export DATABASE_URL=<your-database-url>`.
3. If you want to keep using Access Control Policies and the Warden API, you
   must install ORY Keto using the binaries or the docker image:
   1. `$ keto migrate hydra $DATABASE_URL`.
   2. `$ keto migrate sql $DATABASE_URL`.
   3. Read how to update the JWK storage:
      [AES-GCM nonce storage](#aes-gcm-nonce-storage). If you only use
      auto-generated keys and have never used POST or PUT on the `/keys` API,
      you can probably just execute a `DELETE FROM hydra_jwk` to just remove all
      the auto-generated keys. When starting the ORY Hydra 1.0.0 CLI the
      required keys will be re-generated automatically.
   4. `$ hydra migrate sql $DATABASE_URL`.
4. If you don't use Access Control Policies nor the Warden API, you can skip ORY
   Keto:
   1. Read how to update the JWK storage:
      [AES-GCM nonce storage](#aes-gcm-nonce-storage) and **read point 3.3 of
      this list**.
   2. `$ hydra migrate sql $DATABASE_URL`.
5. `$ export SCOPE_STRATEGY=DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY` - this will
   set the
   [scope strategy](#replacing-hierarchical-scope-strategy-with-wildcard-scope-strategy)
   to the old scope strategy used in version 0.9.x. If you set this, you don't
   need to update the scopes your OAuth 2.0 Clients are allowed to request.
6. `$ hydra help serve`.
7. `$ hydra serve --your-flags`.

It's still a good idea to read through the changes in 0.10.0, for example:
[Response payload changes to `/warden/token/allowed`](#response-payload-changes-to-wardentokenallowed).
You can, however, skip the [New consent flow](#new-consent-flow) subsection in
the 0.10.0 section. All required changes are explained in detail in this
release's consent flow description.

### OpenID Connect Certified

ORY Hydra is now OpenID Connect Certified! Certification spans the OAuth 2.0
Authorize Code Flow, Implicit Flow, and Hybrid Flow as well as dynamic
discovery.

The certification is one reason for the breaking changes in the consent app.

### Breaking Changes

#### Introspection API

One change has been made to the introspection API which is that key `aud` is no
longer a string, but an array of strings. As this claim has not been supported
actively up until now, this will most likely not affect you at all.

##### Introspection is now capable of introspecting refresh tokens

Previously, we disabled the introspection of refresh tokens. This has now
changed to comply with the OAuth 2.0 specification. To distinguish tokens, use
the `token_type` in the introspection response. It can either be `access_token`
or `refresh_token`.

#### Access Control & Warden API

Internal access control, access control policies, and the Warden API have moved
to a separate project called [ORY Keto](https://github.com/ory/keto). You will
be able to run a combination of ORY Hydra,
[ORY Oathkeeper](https://github.com/ory/oathkeeper), and
[ORY Keto](https://github.com/ory/keto) which will be backwards compatible with
ORY Hydra before the 1.0.0 release. This section explains how to upgrade and
links to an example explaining the set up of the three services.

**This means that ORY Hydra has no longer any type of internal access control.
Endpoints such as `POST /clients` no longer require access tokens to be
accessed. You must secure these endpoints yourself. For more information,
[click here](https://www.ory.sh/docs/hydra/production).**

[ORY Keto](https://github.com/ory/keto) handles access control using access
control policies. The project currently supports the Warden API, Access Control
Policy management, and Roles (previously known as
[Warden Groups](#warden-groups)). ORY Keto is independent from ORY Hydra as it
does not rely on any proprietary APIs but instead uses open standards such as
OAuth 2.0 Token Introspection and the OAuth 2.0 Client Credentials Grant to
authenticate credentials. ORY Keto can be used as a standalone project, and
might even be used with other OAuth 2.0 providers, opening up tons of possible
use cases and scenarios. To learn more about the project, head over to
[github.com/ory/keto](https://github.com/ory/keto).

Assuming that you have the 1.0.0 release binary of ORY Hydra and ORY Keto
locally installed, you can migrate the existing policies and Warden Groups using
the migrate commands. Please back up your database before doing this:

```
$ export DATABASE_URL=<your-database-url>

# Migrate the policies and warden groups to keto
$ keto migrate hydra $DATABASE_URL

# Create other Keto database schemas
$ keto migrate sql $DATABASE_URL

# Run Hydra migrations
$ hydra migrate sql $DATABASE_URL
```

Now you can run `keto serve` and endpoints `/policies` as well as `/warden` will
be available at ORY Keto's URL.

##### Running the backwards compatible set up

We have set up a docker-compose example of a set up that resembles ORY Hydra
prior to this release. You can find the source and documentation at
[github.com/ory/examples](https://github.com/ory/examples).

If you find it difficult to run this set up but would like to use the old access
control mechanisms, feel free to reach out to us at
[mailto:hi@ory.sh](hi@ory.sh).

###### Warden API

The Warden endpoints have moved to a new project. Thus, obviously, the URL
changes too. The Warden API paths have changed as well:

- `/warden/allowed` is now `/warden/subjects/authorize`
- `/warden/token/allowed` is now `/warden/oauth2/access-tokens/authorize`
- `/warden/oauth2/clients/authorize` is a new endpoint that lets you authorize
  OAuth 2.0 Clients using their ID and secret.

The backwards compatible set up properly forwards the old paths. If you use that
image and you have been using `http://my-hydra/warden/token/allowed` previously,
you can still use that URL to access that functionality if the
backwards-compatible image is hosted at that location. This image does, however,
currently not rewrite the request and response payloads. If you think that's a
good idea, [let us know](https://github.com/ory/examples/issues/new).

The request payload of these endpoints has changed:

- `/warden/token/allowed` - only key `scopes` was renamed to `scope` in order to
  have a coherent API with any OAuth 2.0 endpoints which use the `scope` for
  singular and plural:
  - Key `scopes` is now `scope` - a response body is
    `{ "token": "...", "action": "...", "resource": "...", "scope": ["scope-a", "scope-b"] }`
    instead of (previously)
    `{ "token": "...", "action": "...", "resource": "...", "scopes": ["scope-a", "scope-b"] }`.

All other endpoints have not experienced any request payload changes.

The response payload of these endpoints has changed:

- `/warden/token/allowed` - keys have been changed to conform to the OAuth 2.0
  Introspection response payload and offer a coherent API.
  - Key `grantedScopes` is now `scope` and is no longer an array string but
    rather a space-delimited string ("scope-a scope-b").
  - Key `clientId` is now `client_id`.
  - Key `issuedAt` is now `iat`.
  - Key `expires_at` is now `exp`.
  - Key `subject` is now `sub`.
  - Key `accessTokenExtra` is now `session` and might be omitted if the OAuth
    2.0 Introspection Endpoint does not provide session data.
  - Key `aud` ("audience") has been added as a string array.
  - Key `iss` ("issuer") has been added.
  - Key `nbf` ("not before") has been added.

We are aware that these changes are rather serious, especially if you rely on
the Warden API in each of your endpoints. If you have ideas on how to improve
upgrading or offer a backwards compatible API, please
[open an issue](https://github.com/ory/keto/issues/new) and let us know.

All other endpoints have not experienced any response payload changes.

###### Warden Groups

Warden Groups have been an experiment determined to simplify managing multiple
subjects with the same access rights. In ORY Keto, Warden Groups have been
renamed to **Roles** and the endpoint has moved from `/warden/groups` to
`/roles`. No request or response payloads have changed, only the URL is a
different one.

If you use the backwards-compatible image, you can access roles using the
`/warden/groups` path as you did before.

#### jwk: Forces JWK to have a unique ID

Previously, JSON Web Keys did not have to specify a unique id. JWKs generated by
ORY Hydra typically only used `public` or `private` as KeyID. This patch changes
that and appends a unique id if no KeyID was given. To be able to separate
between public and private key pairs in resource name, the public/private
convention was kept.

This change targets specifically the OpenID Connect ID Token and HTTP TLS keys.
The ID Token key was previously "hydra.openid.id-token:public" and
"hydra.openid.id-token:private" which now changed to something like
"hydra.openid.id-token:public:9a458aa3-65a0-4982-835f-343eec45183c" and
"hydra.openid.id-token:private:fa353995-d77d-420a-b967-63bf0721271b" with the
UUID part being random for every installation.

This change will help greatly with key rotation in the future.

If you rely on these keys in your applications and if they are hardcoded in some
way, you may want to use the `/.well-known/openid-configuration` or
`/.well-known/jwks.json` endpoints instead. Libraries, which handle these
standards appropriately, exist for almost any programming language.

These keys will be generated automatically if they do not exist yet in the
database. No further steps for upgrading are required.

#### Consent Flow

The consent flow has been refactored in order to implement session (login &
consent) management in ORY Hydra and in order to properly support OpenID Connect
parameters such as `prompt`, `max_age`, and others.

First, the consent flow has been renamed to "User Login and Consent Flow". The
consent app has been renamed to `User Login Provider` and
`User Consent Provider`. If you implement both features (explained in the next
sections) in one program, you can call it the `User Login and Consent Provider`.

A reference implementation of the new User Login and Consent Provider is
available at
[github.com/ory/hydra-login-consent-node](https://github.com/ory/hydra-login-consent-node).

The major difference between the old and new flow is, that authentication (user
login) and scope authorization (user consent) are now two separate endpoints.

The new User Login and Consent Flow is documented in the
[developer guide](https://www.ory.sh/docs/hydra/).

#### Changes to the CLI

The CLI has changed in order to improve developer experience and adopt to the
changes made with this release.

##### `hydra host`

The command `hydra host` has been renamed to `hydra serve` as projects ORY
Oathkeeper and ORY Keto use the `serve` terminology as well.

Because this patch removes the internal access control, no root client and root
policy will be created upon start up. Thus, environment variable
`FORCE_ROOT_CLIENT_CREDENTIALS` has been removed without replacement.

To better reflect what environment variables touch which system, ISSUER has been
renamed to `OAUTH2_ISSUER_URL` and `CONSENT_URL` has been renamed to
`OAUTH2_CONSENT_URL`.

Additionally, flag `--dangerous-force-auto-logon` has been removed it has no
effect any more.

##### `hydra connect`

The command `hydra connect` has been removed as it no longer serves a purpose
now that the internal access control has been removed. Every command you call
now needs the environment variable `HYDRA_URL` (previously named `CLUSTER_URL`)
which should point to ORY Hydra's URL. Removing this command has an additional
benefit - privileged client IDs and secrets will no longer be stored in a
plaintext file on your system if you use this command.

As access control has been removed, most commands (except `token user`,
`token client`, `token revoke`, `token introspect`) work without supplying any
credentials at all. The listed exceptions support setting an OAuth 2.0 Client ID
and Client Secret using flags `--client-id` and `--client-secret` or environment
variables `OAUTH2_CLIENT_ID` and `OAUTH2_CLIENT_SECRET`.

All other commands, such as `hydra clients create`, still support scenarios
where you would need an OAuth2 Access Token. In those cases, you can supply the
access token using flag `--access-token` or environment variable
`OAUTH2_ACCESS_TOKEN`. Assuming that you would like to automate management in a
protected scenario, you could do something like this:

```
$ token=$(hydra token client --client-id foo --client-secret bar --endpoint http://foobar)
$ hydra clients create --access-token $token ...
```

All commands now support the `--endpoint` flag which sets the `HYDRA_URL` in
case you don't want to use environment variables.

##### `hydra token user`

Flags `--id` and `--secret` are now called `--client-id` and `--client-secret`.

##### `hydra token client`

Flags `--client-id` and `--client-secret` have been added.

Flag `--scopes` has been renamed to `--scope`.

##### `hydra token validate`

This command has been renamed to `hydra token introspect` to properly reflect
that you are performing OAuth 2.0 Token Introspection.

Flags `--client-id` and `--client-secret` have been added.

Flag `--scopes` has been renamed to `--scope`.

##### `hydra clients create`

As OAuth 2.0 specifies that terminology `scope` does not have a plural `scopes`,
we updated the places where the incorrect `scopes` was used in order to provide
a more consistent developer experience.

This command renamed flag `--allowed-scopes` to `--scope`.

##### `hydra migrate ladon`

This command is a relict of an old version of ORY Hydra which is, according to
our metrics, not being used any more.

##### `hydra policies`

This command has moved to Keto. All commands work the same way, but you have to
have Keto installed and replace `hydra` with `keto`. For example
`hydra policies create ...` is now `keto policies create ...`

##### `hydra groups`

This command has moved to Keto. All commands work the same way, but you have to
have Keto installed and replace `hydra groups` with `keto roles`. For example
`hydra groups create ...` is now `keto roles create ...`

#### SDK

As the SDK is code-generated, and we are not specialists in every language, we
have only documented changes to the Go API. Please help improving this section
by adding upgrade guides for the SDK you upgraded.

The following methods have been moved.

- The Access Control Policy SDK has moved to ORY Keto:
  - `CreatePolicy(body swagger.Policy) (*swagger.Policy, *swagger.APIResponse, error)`
    is now available via `github.com/ory/keto/sdk/go/keto`. The method signature
    has not changed, apart from types
    `github.com/ory/hydra/sdk/go/hydra/swagger` now being included from
    `github.com/ory/keto/sdk/go/keto/swagger`.
  - `DeletePolicy(id string) (*swagger.APIResponse, error)` is now available via
    `github.com/ory/keto/sdk/go/keto`. The method signature has not changed,
    apart from types `github.com/ory/hydra/sdk/go/hydra/swagger` now being
    included from `github.com/ory/keto/sdk/go/keto/swagger`.
  - `GetPolicy(id string) (*swagger.Policy, *swagger.APIResponse, error)` is now
    available via `github.com/ory/keto/sdk/go/keto`. The method signature has
    not changed, apart from types `github.com/ory/hydra/sdk/go/hydra/swagger`
    now being included from `github.com/ory/keto/sdk/go/keto/swagger`.
  - `ListPolicies(offset int64, limit int64) ([]swagger.Policy, *swagger.APIResponse, error)`
    is now available via `github.com/ory/keto/sdk/go/keto`. The method signature
    has not changed, apart from types
    `github.com/ory/hydra/sdk/go/hydra/swagger` now being included from
    `github.com/ory/keto/sdk/go/keto/swagger`.
  - `UpdatePolicy(id string, body swagger.Policy) (*swagger.Policy, *swagger.APIResponse, error)`
    is now available via `github.com/ory/keto/sdk/go/keto`. The method signature
    has not changed, apart from types
    `github.com/ory/hydra/sdk/go/hydra/swagger` now being included from
    `github.com/ory/keto/sdk/go/keto/swagger`.
- The Warden Group SDK has moved to Keto:
  - `AddMembersToGroup(id string, body swagger.GroupMembers) (*swagger.APIResponse, error)`
    is now
    `AddMembersToRole(id string, body swagger.RoleMembers) (*swagger.APIResponse, error)`
    and is now available via `github.com/ory/keto/sdk/go/keto`.
  - `CreateGroup(body swagger.Group) (*swagger.Group, *swagger.APIResponse, error)`
    is now
    `CreateRole(body swagger.Role) (*swagger.Role, *swagger.APIResponse, error`
    and is now available via `github.com/ory/keto/sdk/go/keto`.
  - `DeleteGroup(id string) (*swagger.APIResponse, error)` is now
    `DeleteRole(id string) (*swagger.APIResponse, error)` and is now available
    via `github.com/ory/keto/sdk/go/keto`.
  - `ListGroups(member string, limit, offset int64) ([]swagger.Group, *swagger.APIResponse, error)`
    is now
    `ListRoles(member string, limit int64, offset int64) ([]swagger.Role, *swagger.APIResponse, error)`
    and is now available via `github.com/ory/keto/sdk/go/keto`.
  - `GetGroup(id string) (*swagger.Group, *swagger.APIResponse, error)` is now
    `GetRole(id string) (*swagger.Role, *swagger.APIResponse, error)` and is now
    available via `github.com/ory/keto/sdk/go/keto`.
  - `RemoveMembersFromGroup(id string, body swagger.GroupMembers) (*swagger.APIResponse, error)`
    is now
    `RemoveMembersFromRole(id string, body swagger.RoleMembers) (*swagger.APIResponse, error)`
    and is now available via `github.com/ory/keto/sdk/go/keto`.
- The Warden API SDK has moved to Keto:
  - `DoesWardenAllowAccessRequest(body swagger.WardenAccessRequest) (*swagger.WardenAccessRequestResponse, *swagger.APIResponse, error)`
    is now
    `IsSubjectAuthorized(body swagger.WardenSubjectAuthorizationRequest) (*swagger.WardenSubjectAuthorizationResponse, *swagger.APIResponse, error)`.
    Please check out the changes to the request/response body as well.
  - `DoesWardenAllowTokenAccessRequest(body swagger.WardenTokenAccessRequest) (*swagger.WardenTokenAccessRequestResponse, *swagger.APIResponse, error)`
    is now
    `IsOAuth2AccessTokenAuthorized(body swagger.WardenOAuth2AccessTokenAuthorizationRequest) (*swagger.WardenOAuth2AccessTokenAuthorizationResponse, *swagger.APIResponse, error)`.
    Please check out the changes to the request/response body as well.
- The Consent API SDK has been deprecated:
  - `AcceptOAuth2ConsentRequest(id string, body swagger.ConsentRequestAcceptance) (*swagger.APIResponse, error)`
    has been removed without replacement.
  - `GetOAuth2ConsentRequest(id string) (*swagger.OAuth2ConsentRequest, *swagger.APIResponse, error)`
    has been removed without replacement.
  - `RejectOAuth2ConsentRequest(id string, body swagger.ConsentRequestRejection) (*swagger.APIResponse, error)`
    has been removed without replacement.
- The Login & Consent API SDK has been added:
  - `AcceptConsentRequest(challenge string, body swagger.AcceptConsentRequest) (*swagger.CompletedRequest, *swagger.APIResponse, error)`
  - `AcceptLoginRequest(challenge string, body swagger.AcceptLoginRequest) (*swagger.CompletedRequest, *swagger.APIResponse, error)`
  - `RejectConsentRequest(challenge string, body swagger.RejectRequest) (*swagger.CompletedRequest, *swagger.APIResponse, error)`
  - `RejectLoginRequest(challenge string, body swagger.RejectRequest) (*swagger.CompletedRequest, *swagger.APIResponse, error)`
  - `GetLoginRequest(challenge string) (*swagger.LoginRequest, *swagger.APIResponse, error)`
  - `GetConsentRequest(challenge string) (*swagger.ConsentRequest, *swagger.APIResponse, error)`

Additionally, the following methods have been removed as they were of very
little use and also mixed the Client Credentials flow with the Authorize Code
Flow which lead to weird usage. It's much easier to configure
`clientcredentials.Config` or `oauth2.Config` yourself.

- `GetOAuth2ClientConfig() (*clientcredentials.Config)`
- `GetOAuth2Config() (*oauth2.Config)`

### Improvements

#### Health Check endpoint has moved

The health check endpoint has moved from `/health/status` to `/health/alive`. We
set up a 308 redirect from `/health/status` to `/health/alive` so this should
not cause any issues.

The `/health/alive` endpoint returns `200 OK` as soon as the HTTP server is
responsive. Another endpoint `/health/ready` was added which returns `200 OK`
only if the database connection is working as well.

As part of this change, function `getInstanceStatus` of the SDK is now
`isInstanceAlive` and `isInstanceReady`.

#### Unknown request body payloads result in error

Previously, if you had a typo in the JSON (e.g. `client_nme` instead of
`client_name`), ORY Hydra simply ignored that key. Now, an error is thrown if
unknown JSON keys are included.

#### UTC everywhere

ORY Hydra now uses UTC everywhere, reducing the possibility of errors stemming
from different timezones.

#### Pagination everywhere

Each endpoint that returns a list of items now supports pagination using `limit`
and `offset` query parameters.

#### Flushing old access tokens

An endpoint (`/oauth2/flush`) has been added that allows you to flush old access
tokens.

#### Prometheus endpoint

An endpoint `/health/prometheus` for providing data to Prometheus has been
added.

## 0.11.12

This release resolves a security issue (reported by
[platform.sh](https://www.platform.sh)) related to the fosite storage
implementation in this project. Fosite used to pass all of the request body from
both authorize and token endpoints to the storage adapters. As some of these
values are needed in consecutive requests, the storage adapter of this project
chose to drop all of the key/value pairs to the database in plaintext.

This implied that confidential parameters, such as the `client_secret` which can
be passed in the request body since fosite version 0.15.0, were stored as
key/value pairs in plaintext in the database. While most client secrets are
generated programmatically (as opposed to set by the user) and most popular
OAuth2 providers choose to store the secret in plaintext for later retrieval, we
see it as a considerable security issue nonetheless.

The issue has been resolved by sanitizing the request body and only including
those values truly required by their respective handlers. This also implies that
typos (eg `client_secet`) won't "leak" to the database.

There are no special upgrade paths required for this version.

This issue does not apply to you if you do not use an SQL backend. If you do
upgrade to this version, you need to run
`hydra migrate sql path://to.your/database`.

If your users use POST body client authentication, it might be a good move to
remove old data. There are multiple ways of doing that. **Back up your data
before you do this**:

1. **Radical solution:** Drop all rows from tables `hydra_oauth2_refresh`,
   `hydra_oauth2_access`, `hydra_oauth2_oidc`, `hydra_oauth2_code`. This implies
   that all your users have to re-authorize.
2. **Sensitive solution:** Replace all values in column `form_data` in tables
   `hydra_oauth2_refresh`, `hydra_oauth2_access` with an empty string. This will
   keep all authorization sessions alive. Tables `hydra_oauth2_oidc` and
   `hydra_oauth2_code` do not contain sensitive information, unless your users
   accidentally sent the client_secret to the `/oauth2/auth` endpoint.

We would like to thank [platform.sh](https://www.platform.sh) for sponsoring the
development of a patch that resolves this issue.

## 0.11.3

The experimental endpoint `/health/metrics` has been removed as it caused
various issues such as increased memory usage, and it was apparently not used at
all.

## 0.11.0

This release has a minor breaking change in the experimental Warden Group SDK:
`FindGroupsByMember(member string) ([]swagger.Group, *swagger.APIResponse, error)`
is now
`ListGroups(member string, limit, offset int64) ([]swagger.Group, *swagger.APIResponse, error)`.
The change has to be applied in a similar fashion to other SDKs generated using
swagger.

Leave the `member` parameter empty to list all groups, and add it to filter
groups by member id.

## 0.10.0

This release has several major improvements, and some breaking changes. It
focuses on cryptographic security by leveraging best practices that emerged
within the last one and a half years. Before upgrading to this version, make a
back up of the JWK table in your SQL database.

This release requires running `hydra migrate sql` before `hydra host`.

The most important breaking changes are the SDK libraries, the new consent flow,
the AES-GCM improvement, and the response payload changes to the warden.

We know that these are a lot of changes, but we highly recommend upgrading to
this version. It will be the last before releasing 1.0.0.

### Breaking Changes

#### Introspection now requires authorization

The introspection endpoint was previously accessible to anyone with valid client
credentials or a valid access token. According to spec, the introspection
endpoint should be protected by additional access control mechanisms. This
version introduces new access control requirements for this endpoint.

The client id of the basic authorization / subject of the bearer token must be
allowed action `introspect` on resource `rn:hydra:oauth2:tokens`. If an access
token is used for authorization, it needs to be granted the `hydra.introspect`
scope.

#### New consent flow

Previously, the consent flow looked roughly like this:

1. App asks user for Authorization by generating the authorization URL
   (http://hydra.mydomain.com/oauth2/auth?client_id=...).
1. Hydra asks browser of user for authentication by redirecting to the Consent
   App with a _consent challenge_
   (http://login.mydomain.com/login?challenge=xYt...).
1. Retrieves a RSA 256 public key from Hydra.
1. Uses said public key to verify the consent challenge.
1. User logs in and authorizes the requested scopes
1. Consent app generates the consent response
1. Retrieves a private key from Hydra which is used to sign the consent
   response.
1. Creates a response message and sign with said private key.
1. Redirects the browser back to Hydra, appending the consent response
   (http://hydra.mydomain.com/oauth2/auth?client_id=...&consent=zxI...).
1. Hydra validates consent response and generates access tokens, authorize
   codes, refresh tokens, and id tokens.

This approach had several disadvantages:

1. Validating and generating the JSON Web Tokens (JWTs) requires libraries for
   each language
1. Because libraries are required, auto generating SDKs from the swagger spec is
   impossible. Thus, every language requires a maintained SDK which
   significantly increases our workload.
1. There have been at least two major bugs affecting almost all JWT libraries
   for any language. The spec has been criticised for it's mushy language.
1. The private key used by the consent app for signing consent responses was
   originally thought to be stored at the consent app, not in Hydra. However,
   since Hydra offers JWK storage, it was decided to use the Hydra JWK store per
   default for retrieval of the private key to improve developer experience.
   However, to make really sense, the private key should have been stored at the
   consent app, not in Hydra.
1. Private/Public keypairs need to be fetched on every request or cached in a
   way that allows for key rotation, complicating the consent app.
1. There is currently no good mechanism for rotating JWKs in Hydra's storage.
1. The consent challenge / response has a limited length as it's transmitted via
   the URL query. The length of a URL is limited.

Due to these reasons we decided to refactor the consent flow. Instead of relying
on JWTs using RSA256, a simple HTTP call is now enough to confirm a consent
request:

1. App asks user for Authorization by generating the authorization URL
   (http://hydra.mydomain.com/oauth2/auth?client_id=...).
1. Hydra asks browser of user for authentication by redirecting to the Consent
   App with a unique _consent request id_
   (http://login.mydomain.com/login?consent=fjad2312).
1. Consent app makes a HTTP REST request to
   `http://hydra.mydomain.com/oauth2/consent/requests/fjad2312` and retrieves
   information on the authorization request.
1. User logs in and authorizes the requested scopes
1. Consent app accepts or denies the consent request by making a HTTP REST
   request to
   `http://hydra.mydomain.com/oauth2/consent/requests/fjad2312/accept` or
   `http://hydra.mydomain.com/oauth2/consent/requests/fjad2312/reject`.
1. Redirects the browser back to Hydra.
1. Hydra validates consent request by checking if it was accepted and generates
   access tokens, authorize codes, refresh tokens, and id tokens.

Learn more on how the new consent flow works in the guide:
https://ory.gitbooks.io/hydra/content/oauth2.html#consent-flow

#### Audience

Previously, the audience terminology was used as a synonym for OAuth2 client
IDs. This is no longer the case. The audience is typically a URL identifying the
endpoint(s) the token is intended for. For example, if a client requires access
to endpoint `http://mydomain.com/users`, then the audience would be
`http://mydomain.com/users`.

This changes the payload of `/warden/token/allowed` and is incorporated in the
new consent flow as well. Please note that it is currently not possible to set
the audience of a token. This feature is tracked with
[here](https://github.com/ory/hydra/issues/687).

**IMPORTANT NOTE:** In OpenID Connect ID Tokens, the token is issued for that
client. Thus, the `aud` claim must equal to the `client_id` that initiated the
request.

#### Response payload changes to `/warden/token/allowed`

Previously, the response of the warden endpoint contained shorthands like `aud`,
`iss`, and so on. Those have now been changed to their full names:

- `sub` is now named `subject`.
- `scopes` is now named `grantedScopes`.
- `iss` is now named `issuer`.
- `aud` is now named `clientId`.
- `iat` is now named `issuedAt`.
- `exp` is now named `expiresAt`.
- `ext` is now named `accessTokenExtra`.

#### Go SDK

The Go SDK was completely replaced in favor of a SDK based on `swagger-codegen`.
Unfortunately this means that any code relying on the old SDK has to be
replaced. On the bright side the dependency tree is much smaller as no direct
dependencies to ORY Hydra's code base exist any more.

Read more on it here: https://ory.gitbooks.io/hydra/content/sdk/go.html

#### Health endpoints

- `GET /health` is now `GET /health/status`
- `GET /health/stats` is now `GET /health/metrics`

#### Group endpoints

`GET /warden/groups` now returns a list of groups, not just a list of strings
(group ids).

#### Replacing hierarchical scope strategy with wildcard scope strategy

The previous scope matching strategy has been replaced in favor of a
wildcard-based matching strategy. Previously, `foo` matched `foo` and `foo.bar`
and `foo.baz`. This is no longer the case. So `foo` matches only `foo`. Matching
subsets is possible using wildcards. `foo.*` matches `foo.bar` and `foo.baz`.

This change makes setting scopes more explicit and is more secure, as it is less
likely to make mistakes.

Read more on this strategy
[here](https://www.ory.sh/docs/hydra/oauth2#oauth-20-scope).

To fall back to hierarchical scope matching, set the environment variable
`SCOPE_STRATEGY=DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY`. This feature _might_ be
fully removed in a later version.

#### AES-GCM nonce storage

Our use of `crypto/aes`'s AES-GCM was replaced in favor of
[`cryptopasta/encrypt`](https://github.com/gtank/cryptopasta/blob/master/encrypt.go).
As this includes a change of how nonces are appended to the ciphertext, ORY
Hydra will be unable to decipher existing databases.

There are two paths to migrate this change:

1. If you have not added any keys to the JWK store:
   1. Stop all Hydra instances.
   2. Drop all rows from the `hydra_jwk` table.
   3. Start **one** Hydra instance and wait for it to boot.
   4. Restart all remaining Hydra instances.
2. If you added keys to the JWK store:
   1. If you can afford to re-generate those keys:
      1. Write down all key ids you generated.
      2. Stop all Hydra instances.
      3. Drop all rows from the `hydra_jwk` table.
      4. Start **one** Hydra instance and wait for it to boot.
      5. Restart all remaining Hydra instances.
      6. Regenerate the keys and use the key ids you wrote down.
   2. If you can not afford to re-generate the keys:
      1. Export said keys using the REST API.
      2. Stop all Hydra instances.
      3. Drop all rows from the `hydra_jwk` table.
      4. Start **one** Hydra instance and wait for it to boot.
      5. Restart all remaining Hydra instances.
      6. Import said keys using the REST API.

#### Minor Breaking Changes

##### Token signature algorithm changed from HMAC-SHA256 to HMAC-SHA512

The signature algorithm used to generate authorize codes, access tokens, and
refresh tokens has been upgraded from HMAC-SHA256 to HMAC-SHA512. With upgrading
to alpha.9, all previously issued authorize codes, access tokens, and refresh
will thus be rendered invalid. Apart from some re-authorization procedures,
which are usually automated, this should not have any significant impact on your
installation.

##### HS256 JWK Generator now uses all 256 bit

The HS256 (symmetric/shared keys) JWK Generator now uses the full 256 bit range
to generate secrets instead of a predefined rune sequence. This change only
affects keys generated in the future.

##### ES512 Key generator

The JWK algorithm `ES521` was renamed to `ES512`. If you want to generate a key
using this algorithm, you have to use the update name in the future.

##### Build tags deprecated

This release removes build tags `-http`, `-automigrate`, `-without-telemetry`
from the docker hub repository and replaces it with a new and tiny (~6MB) docker
image containing the binary only. Please note that this docker image does not
have a shell, which makes it harder to penetrate.

Instead of relying on tags to pass arguments, it is now possible to pass command
arguments such as `docker run oryd/hydra:v0.10.0 host --dangerous-force-http`
directly.

**Version 0.10.8 reintroduces an image with a shell, appended with tag
`-alpine`.**

### Important Additions

#### Prefixing Resources Names

It is now possible to alter resource name prefixes (`rn:hydra`) using the
`RESOURCE_NAME_PREFIX` environment variable.

#### Refreshing OpenID Connect ID Token using `refresh_token` grant type

1. It is now possible to refresh openid connect tokens using the refresh_token
   grant. An ID Token is issued if the scope `openid` was requested, and the
   client is allowed to receive an ID Token.

### Important Changes

#### Telemetry

To improve ORY Hydra and understand how the software is used, optional,
anonymized telemetry data is shared with ORY. A change was made to help us
understand which telemetry sources belong to the same installation by hashing
(SHA256) two environment variables which make up a unique identifier.
[Click here](https://www.ory.sh/docs/ecosystem/sqa) to read more about how we
collect telemetry data, why we do it, and how to enable or disable it.

#### URL Encoding Root Client Credentials

This release adds the possibility to specify special characters in the
`FORCE_ROOT_CLIENT_CREDENTIALS` by `www-url-decoding` the values. If you have
characters that are not url safe in your root client credentials, please use the
following form to specify them:
`"FORCE_ROOT_CLIENT_CREDENTIALS=urlencode(id):urlencode(secret)"`.

## 0.9.0

This version adds performance metrics to `/health` and sends anonymous usage
statistics to our servers, [click here](https://www.ory.sh/docs/ecosystem/sqa)
for more details on this feature and how to disable it.

## 0.8.0

This PR improves some performance bottlenecks, offers more control over Hydra,
moves to Go 1.8, and moves the REST documentation to swagger.

**Before applying this update, please make a back up of your database. Do not
upgrade directly from versions below 0.7.0**.

To upgrade the database schemas, please run the following commands in exactly
this order

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

Ladon was greatly improved with version 0.6.0, resolving various performance
bottlenecks. Please read more on this release
[here](https://github.com/ory/ladon/blob/master/HISTORY.md#060).

#### Redis and RethinkDB deprecated

Redis and RethinkDB are removed from the repository now and no longer supported,
see [this issue](https://github.com/ory/hydra/issues/425).

#### Moved to ory namespace

To reflect the GitHub organization rename, Hydra was moved from
`https://github.com/ory-am/hydra` to `https://github.com/ory/hydra`.

#### SDK

The method `FindPoliciesForSubject` of the policy SDK was removed. Instead,
`List` was added. The HTTP endpoint `GET /policies` no longer allows to query by
subject.

#### JWK

To generate JWKs previously the payload at `POST /keys` was
`{ "alg": "...", "id": "some-id" }`. `id` was changed to `kid` so this is now
`{ "alg": "...", "kid": "some-id" }`.

#### Migrations are no longer automatically applied

SQL Migrations are no longer automatically applied. Instead you need to run
`hydra migrate sql` after upgrading to a Hydra version that includes a breaking
schema change.

### Changes

#### Log format: json

Set the log format to json using `export LOG_FORMAT=json`

#### SQL Connection Control

You can configure SQL connection limits by appending parameters `max_conns`,
`max_idle_conns`, or `max_conn_lifetime` to the DSN:
`postgres://foo:bar@host:port/database?max_conns=12`.

#### REST API Docs are now generated from source code

... and are swagger 2.0 spec.

#### Documentation on scopes

Documentation on scopes (e.g. offline) was added.

#### New response writer library

Hydra now uses `github.com/ory/herodot` for writing REST responses. This
increases compatibility with other libraries and resolves a few other issues.

#### Graceful http handling

Hydra is now capable of gracefully handling SIGINT.

#### Best practice HTTP server config

Hydra now implements best practices for running HTTP servers that are exposed to
the public internet.
