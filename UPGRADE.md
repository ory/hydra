# Upgrading

Please refer to [CHANGELOG.md](./CHANGELOG.md) for a full list of changes.

The intent of this document is to make migration of breaking changes as easy as possible. Please note that not all
breaking changes might be included here. Refer to refer to [CHANGELOG.md](./CHANGELOG.md) for a full list of changes
before finalizing the upgrade process.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [0.10.0](#0100)
  - [Breaking Changes](#breaking-changes)
    - [Introspection now requires authorization](#introspection-now-requires-authorization)
    - [New consent flow](#new-consent-flow)
    - [Audience](#audience)
    - [Response payload changes to `/warden/token/allowed`](#response-payload-changes-to-wardentokenallowed)
    - [Go SDK](#go-sdk)
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
    - [SDK](#sdk)
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

## 0.10.0

This release has several major improvements, and some breaking changes. It focuses on cryptographic security
by leveraging best practices that emerged within the last one and a half years. Before upgrading to this version,
make a back up of the JWK table in your SQL database.

This release requires running `hydra migrate sql` before `hydra host`.

The most important breaking changes are the SDK libraries, the new consent flow, the AES-GCM improvement, and the
response payload changes to the warden.

We know that these are a lot of changes, but we highly recommend upgrading to this version. It will be the last before
releasing 1.0.0.

### Breaking Changes

#### Introspection now requires authorization

The introspection endpoint was previously accessible to anyone with valid client credentials or a valid access token.
According to spec, the introspection endpoint should be protected by additional access control mechanisms. This
version introduces new access control requirements for this endpoint.

The client id of the basic authorization / subject of the bearer token must be allowed action `introspect`
on resource `rn:hydra:oauth2:tokens`. If an access token is used for authorization, it needs to be granted the
`hydra.introspect` scope.

#### New consent flow

Previously, the consent flow looked roughly like this:

1. App asks user for Authorization by generating the authorization URL (http://hydra.mydomain.com/oauth2/auth?client_id=...).
1. Hydra asks browser of user for authentication by redirecting to the Consent App with a *consent challenge* (http://login.mydomain.com/login?challenge=xYt...).
  1. Retrieves a RSA 256 public key from Hydra.
  2. Uses said public key to verify the consent challenge.
3. User logs in and authorizes the requested scopes
4. Consent app generates the consent response
  1. Retrieves a private key from Hydra which is used to sign the consent response.
  2. Creates a response message and sign with said private key.
  3. Redirects the browser back to Hydra, appending the consent response (http://hydra.mydomain.com/oauth2/auth?client_id=...&consent=zxI...).
6. Hydra validates consent response and generates access tokens, authorize codes, refresh tokens, and id tokens.

This approach had several disadvantages:

1. Validating and generating the JSON Web Tokens (JWTs) requires libraries for each language
  1. Because libraries are required, auto generating SDKs from the swagger spec is impossible. Thus, every language
  requires a maintained SDK which significantly increases our workload.
  2. There have been at least two major bugs affecting almost all JWT libraries for any language. The spec has been criticised
  for it's mushy language.
  3. The private key used by the consent app for signing consent responses was originally thought to be stored at the consent
  app, not in Hydra. However, since Hydra offers JWK storage, it was decided to use the Hydra JWK store per default for
  retrieval of the private key to improve developer experience. However, to make really sense, the private key should have
  been stored at the consent app, not in Hydra.
2. Private/Public keypairs need to be fetched on every request or cached in a way that allows for key rotation, complicating
the consent app.
3. There is currently no good mechanism for rotating JWKs in Hydra's storage.
4. The consent challenge / response has a limited length as it's transmitted via the URL query. The length of a URL
is limited.

Due to these reasons we decided to refactor the consent flow. Instead of relying on JWTs using RSA256, a simple HTTP call
is now enough to confirm a consent request:

1. App asks user for Authorization by generating the authorization URL (http://hydra.mydomain.com/oauth2/auth?client_id=...).
1. Hydra asks browser of user for authentication by redirecting to the Consent App with a unique *consent request id* (http://login.mydomain.com/login?consent=fjad2312).
  1. Consent app makes a HTTP REST request to `http://hydra.mydomain.com/oauth2/consent/requests/fjad2312` and retrieves information on the authorization request.
3. User logs in and authorizes the requested scopes
4. Consent app accepts or denies the consent request by making a HTTP REST request to `http://hydra.mydomain.com/oauth2/consent/requests/fjad2312/accept` or `http://hydra.mydomain.com/oauth2/consent/requests/fjad2312/reject`.
5. Redirects the browser back to Hydra.
6. Hydra validates consent request by checking if it was accepted and generates access tokens, authorize codes, refresh tokens, and id tokens.

Learn more on how the new consent flow works in the guide: https://ory.gitbooks.io/hydra/content/oauth2.html#consent-flow

#### Audience

Previously, the audience terminology was used as a synonym for OAuth2 client IDs. This is no longer the case. The audience
is typically a URL identifying the endpoint(s) the token is intended for. For example, if a client requires access to
endpoint `http://mydomain.com/users`, then the audience would be `http://mydomain.com/users`.

This changes the payload of `/warden/token/allowed` and is incorporated in the new consent flow as well. Please note
that it is currently not possible to set the audience of a token. This feature is tracked with [here](https://github.com/ory/hydra/issues/687).

**IMPORTANT NOTE:** In OpenID Connect ID Tokens, the token is issued for that client. Thus, the `aud` claim must equal
to the `client_id` that initiated the request.

#### Response payload changes to `/warden/token/allowed`

Previously, the response of the warden endpoint contained shorthands like `aud`, `iss`, and so on. Those have now been changed
to their full names:

* `sub` is now named `subject`.
* `scopes` is now named `grantedScopes`.
* `iss` is now named `issuer`.
* `aud` is now named `clientId`.
* `iat` is now named `issuedAt`.
* `exp` is now named `expiresAt`.
* `ext` is now named `accessTokenExtra`.

#### Go SDK

The Go SDK was completely replaced in favor of a SDK based on `swagger-codegen`. Unfortunately this means that
any code relying on the old SDK has to be replaced. On the bright side the dependency tree is much smaller as
no direct dependencies to ORY Hydra's code base exist any more.

Read more on it here: https://ory.gitbooks.io/hydra/content/sdk/go.html

#### Health endpoints

* `GET /health` is now `GET /health/status`
* `GET /health/stats` is now `GET /health/metrics`

#### Group endpoints

`GET /warden/groups` now returns a list of groups, not just a list of strings (group ids).

#### Replacing hierarchical scope strategy with wildcard scope strategy

The previous scope matching strategy has been replaced in favor of a wildcard-based matching strategy. Previously,
`foo` matched `foo` and `foo.bar` and `foo.baz`. This is no longer the case. So `foo` matches only `foo`. Matching
subsets is possible using wildcards. `foo.*` matches `foo.bar` and `foo.baz`.

This change makes setting scopes more explicit and is more secure, as it is less likely to make mistakes.

Read more on this strategy [here](https://ory.gitbooks.io/hydra/content/oauth2.html#oauth2-scopes).

To fall back to hierarchical scope matching, set the environment variable `SCOPE_STRATEGY=DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY`.
This feature *might* be fully removed in a later version.

#### AES-GCM nonce storage

Our use of `crypto/aes`'s AES-GCM was replaced in favor of [`cryptopasta/encrypt`](https://github.com/gtank/cryptopasta/blob/master/encrypt.go).
As this includes a change of how nonces are appended to the ciphertext, ORY Hydra will be unable to decipher existing
databases.

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

The signature algorithm used to generate authorize codes, access tokens, and refresh tokens has been upgraded
from HMAC-SHA256 to HMAC-SHA512. With upgrading to alpha.9, all previously issued authorize codes, access tokens, and refresh will thus be
rendered invalid. Apart from some re-authorization procedures, which are usually automated, this should not have any
significant impact on your installation.

##### HS256 JWK Generator now uses all 256 bit

The HS256 (symmetric/shared keys) JWK Generator now uses the full 256 bit range to generate secrets instead of a
predefined rune sequence. This change only affects keys generated in the future.

##### ES512 Key generator

The JWK algorithm `ES521` was renamed to `ES512`. If you want to generate a key using this algorithm, you have to use
the update name in the future.

##### Build tags deprecated

This release removes build tags `-http`, `-automigrate`, `-without-telemetry` from the docker hub repository and replaces
it with a new and tiny (~6MB) docker image containing the binary only. Please note that this docker image does not have
a shell, which makes it harder to penetrate.

Instead of relying on tags to pass arguments, it is now possible to pass command arguments such as `docker run oryd/hydra:v0.10.0 host --dangerous-force-http`
directly.

**Version 0.10.8 reintroduces an image with a shell, appended with tag `-alpine`.**

### Important Additions

#### Prefixing Resources Names

It is now possible to alter resource name prefixes (`rn:hydra`) using the `RESOURCE_NAME_PREFIX` environment variable.

#### Refreshing OpenID Connect ID Token using `refresh_token` grant type

1. It is now possible to refresh openid connect tokens using the refresh_token grant. An ID Token is issued if the scope
`openid` was requested, and the client is allowed to receive an ID Token.

### Important Changes

#### Telemetry

To improve ORY Hydra and understand how the software is used, optional, anonymized telemetry data is shared with ORY.
A change was made to help us understand which telemetry sources belong to the same installation by hashing (SHA256)
two environment variables which make up a unique identifier. [Click here](https://ory.gitbooks.io/hydra/content/telemetry.html)
to read more about how we collect telemetry data, why we do it, and how to enable or disable it.

#### URL Encoding Root Client Credentials

This release adds the possibility to specify special characters in the `FORCE_ROOT_CLIENT_CREDENTIALS` by `www-url-decoding`
the values. If you have characters that are not url safe in your root client credentials, please use the following
form to specify them: `"FORCE_ROOT_CLIENT_CREDENTIALS=urlencode(id):urlencode(secret)"`.

## 0.9.0

This version adds performance metrics to `/health` and sends anonymous usage statistics to our servers, [click here](https://ory.gitbooks.io/hydra/content/telemetry.html) for more
details on this feature and how to disable it.

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
