# History

This list makes you aware of (breaking) changes. For patch notes, please check the [releases tab](https://github.com/ory/hydra/releases).

## 0.10.0-alpha1

**Warning: This version introduces breaking changes and is not suited for production use yet.**

Version 0.10.0 is a preview tag of the 1.0.0 release. It contains multiple breaking changes.

This release requires running `hydra migrate sql` before `hydra host`.

Please also note that the new scope strategy might render your administrative client incapable of performing requests.
Set the environment variable `SCOPE_STRATEGY=DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY` to temporarily use the previous
scope strategy and migrate the scopes manually. You may append `.*` to all scopes. For example, `hydra` is now `hydra hydra.*`

## New consent flow

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

This approach has several disadvantages:

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

## Audience

Previously, the audience terminology was used as a synonym for OAuth2 client IDs. This is no longer the case. The audience
is typically a URL identifying the endpoint(s) the token is intended for. For example, if a client requires access to
endpoint `http://mydomain.com/users`, then the audience would be `http://mydomain.com/users`.

The audience feature is currently not supported in Hydra, only the terminology changed. Fields named `audience` are thus
renamed to `clientId` (where previously named `audience`) and `cid` (where previously named `aud`).

**IMPORTANT NOTE:** This does **not** apply to OpenID Connect ID tokens. There, the `aud` claim **MUST** match the `client_id`.
This discrepancy between OpenID Connect and OAuth 2.0 is what caused the confusion with the OAuth 2.0 audience terminology.

## Response payload changes to `/warden/token/allowed`

Previously, the response of the warden endpoint contained shorthands like `aud`, `iss`, and so on. Those have now been changed
to their full names. For example, `iss` is now `issuer`. Additionally, `aud` is now named `clientId`.

## Go SDK

The Go SDK was completely replaced in favor of a SDK based on `swagger-codegen`. Read more on it here: https://ory.gitbooks.io/hydra/content/sdk/go.html

## Health endpoints

* `GET /health` is now `GET /health/status`
* `GET /health/stats` is now `GET /health/metrics`

## Group endpoints

* `GET /warden/groups` now returns a list of groups, not just a group id

## Refreshing OpenID Connect ID Token using `refresh_token` grant type

1. It is now possible to refresh openid connect tokens using the refresh_token grant. An ID Token is issued if the scope
`openid` was requested, and the client is allowed to receive an ID Token.

## Replacing hierarchical scope strategy with wildcard scope strategy

The previous scope matching strategy has been replaced in favor of a wildcard-based matching strategy. Read more
on this strategy [here](https://ory.gitbooks.io/hydra/content/oauth2.html#oauth2-scopes).

To fall back to hierarchical scope matching, set the environment variable `SCOPE_STRATEGY=DEPRECATED_HIERARCHICAL_SCOPE_STRATEGY`.
This feature *might* be fully removed in the final 1.0.0 version.

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
