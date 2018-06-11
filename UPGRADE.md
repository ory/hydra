# Upgrading

The intent of this document is to make migration of breaking changes as easy as possible. Please note that not all
breaking changes might be included here. Please check the [CHANGELOG.md](./CHANGELOG.md) for a full list of changes
before finalizing the upgrade process.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [1.0.0-alpha.1](#100-beta1)
  - [OpenID Connect Certified](#openid-connect-certified)
  - [Breaking Changes](#breaking-changes)
    - [Introspection API](#introspection-api)
      - [Introspection is now capable of introspecting refresh tokens](#introspection-is-now-capable-of-introspecting-refresh-tokens)
    - [Access Control & Warden API](#access-control-&-warden-api)
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
  - [Breaking Changes](#breaking-changes-1)
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

## 1.0.0-beta.1

This section summarizes important changes introduced in 1.0.0. **Follow it chronologically to ensure a proper migration.**

We are very well aware that the changelist is huge and we try to prepare you as good as we can to migrate to this version.
We also understand that breaking changes are frustrating and it takes time to adopt to them. We sincerely hope that
the benefits from this version (improved consent flow, easier set up, clear boundaries & responsibilities)
outweigh the hassle of upgrading to the new version.

If you have difficulties upgrading and would like a helping hand, reach out to us at [mailto:hi@ory.sh](hi@ory.sh) and
we will help you with the upgrade process. Our services are billed by the hour and are priced fairly.

### OpenID Connect Certified

ORY Hydra is now OpenID Connect Certified! Certification spans the OAuth 2.0 Authorize Code Flow, Implicit Flow, and Hybrid Flow
as well as dynamic discovery.

The certification is one reason for the breaking changes in the consent app.

### Breaking Changes

#### Introspection API

One change has been made to the introspection API which is that key `aud` is no longer a string, but an array of strings.
As this claim has not been supported actively up until now, this will most likely not affect you at all.

##### Introspection is now capable of introspecting refresh tokens

Previously, we disabled the introspection of refresh tokens. This has now changed to comply with the OAuth 2.0 specification.
To distinguish tokens, use the `token_type` in the introspection response. It can either be `access_token` or `refresh_token`.

#### Access Control & Warden API

Internal access control, access control policies, and the Warden API have moved to a separate project called [ORY Keto](https://github.com/ory/keto).
You will be able to run a combination of ORY Hydra, [ORY Oathkeeper](https://github.com/ory/oathkeeper), and [ORY Keto](https://github.com/ory/keto) which will be backwards compatible with
ORY Hydra before the 1.0.0 release. This section explains how to upgrade and links to an example explaining the set up
of the three services.

**This means that ORY Hydra has no longer any type of internal access control. Endpoints such as `POST /clients` no longer
require access tokens to be accessed. You must secure these endpoints yourself. For more information, [click here](https://www.ory.sh/docs/guides/master/hydra/2-environment/1-securing-ory-hydra).**

[ORY Keto](https://github.com/ory/keto) handles access control using access control policies. The project currently supports the Warden API, Access Control Policy
management, and Roles (previously known as [Warden Groups](#warden-groups)). ORY Keto is independent from ORY Hydra
as it does not rely on any proprietary APIs but instead uses open standards such as OAuth 2.0 Token Introspection
and the OAuth 2.0 Client Credentials Grant to authenticate credentials. ORY Keto can be used as a standalone project,
and might even be used with other OAuth 2.0 providers, opening up tons of possible use cases and scenarios. To learn
more about the project, head over to [github.com/ory/keto](https://github.com/ory/keto).

Assuming that you have the 1.0.0 release binary of ORY Hydra and ORY Keto locally installed, you can migrate the existing
policies and Warden Groups using the migrate commands. Please back up your database before doing this:

```
$ export DATABASE_URL=<your-database-url>

# Migrate the policies and warden groups to keto
$ keto migrate hydra $DATABASE_URL

# Create other Keto database schemas
$ keto migrate sql $DATABASE_URL

# Run Hydra migrations
$ hydra migrate sql $DATABASE_URL
```

Now you can run `keto serve` and endpoints `/policies` as well as `/warden` will be available at ORY Keto's URL.

##### Running the backwards compatible set up

We have set up a docker-compose example of a set up that resembles ORY Hydra prior to this release. You can find
the source and documentation at [github.com/ory/examples](https://github.com/ory/examples).

If you find it difficult to run this set up but would like to use the old access control mechanisms, feel free
to reach out to us at [mailto:hi@ory.sh](hi@ory.sh).

###### Warden API

The Warden endpoints have moved to a new project. Thus, obviously, the URL changes too. The Warden API paths have changed
as well:

* `/warden/allowed` is now `/warden/subjects/authorize`
* `/warden/token/allowed` is now `/warden/oauth2/access-tokens/authorize`
* `/warden/oauth2/clients/authorize` is a new endpoint that lets you authorize OAuth 2.0 Clients using their ID and secret.

The backwards compatible set up properly forwards the old paths. If you use that image and you have been using
`http://my-hydra/warden/token/allowed` previously, you can still use that URL to access that functionality if the
backwards-compatible image is hosted at that location. This image does, however, currently not rewrite the request
and response payloads. If you think that's a good idea, [let us know](https://github.com/ory/examples/issues/new).

The request payload of these endpoints has changed:

* `/warden/token/allowed` - only key `scopes` was renamed to `scope` in order to have a coherent API with any OAuth 2.0
endpoints which use the `scope` for singular and plural:
  * Key `scopes` is now `scope` - a response body is `{ "token": "...", "action": "...", "resource": "...", "scope": ["scope-a", "scope-b"] }`
  instead of (previously) `{ "token": "...", "action": "...", "resource": "...", "scopes": ["scope-a", "scope-b"] }`.

All other endpoints have not experienced any request payload changes.

The response payload of these endpoints has changed:

* `/warden/token/allowed` - keys have been changed to conform to the OAuth 2.0 Introspection response payload and offer
a coherent API.
  * Key `grantedScopes` is now `scope` and is no longer an array string but rather a space-delimited string ("scope-a scope-b").
  * Key `clientId` is now `client_id`.
  * Key `issuedAt` is now `iat`.
  * Key `expires_at` is now `exp`.
  * Key `subject` is now `sub`.
  * Key `accessTokenExtra` is now `session` and might be omitted if the OAuth 2.0 Introspection Endpoint does not provide
  session data.
  * Key `aud` ("audience") has been added as a string array.
  * Key `iss` ("issuer") has been added.
  * Key `nbf` ("not before") has been added.

We are aware that these changes are rather serious, especially if you rely on the Warden API in each of your endpoints.
If you have ideas on how to improve upgrading or offer a backwards compatible API, please [open an issue](https://github.com/ory/keto/issues/new)
and let us know.

All other endpoints have not experienced any response payload changes.

###### Warden Groups

Warden Groups have been an experiment determined to simplify managing multiple subjects with the same access rights.
In ORY Keto, Warden Groups have been renamed to **Roles** and the endpoint has moved from `/warden/groups` to `/roles`.
No request or response payloads have changed, only the URL is a different one.

If you use the backwards-compatible image, you can access roles using the `/warden/groups` path as you did before.

#### jwk: Forces JWK to have a unique ID

Previously, JSON Web Keys did not have to specify a unique id. JWKs
generated by ORY Hydra typically only used `public` or `private`
as KeyID. This patch changes that and appends a unique id if no
KeyID was given. To be able to separate between public and private key
pairs in resource name, the public/private convention was kept.

This change targets specifically the OpenID Connect ID Token and HTTP
TLS keys. The ID Token key was previously "hydra.openid.id-token:public"
and "hydra.openid.id-token:private" which now changed to something like
"hydra.openid.id-token:public:9a458aa3-65a0-4982-835f-343eec45183c" and
"hydra.openid.id-token:private:fa353995-d77d-420a-b967-63bf0721271b"
with the UUID part being random for every installation.

This change will help greatly with key rotation in the future.

If you rely on these keys in your applications and if they are hardcoded in some way, you may want to use the `./well-known/openid-configuration`
or `./well-known/jwks.json` endpoints instead. Libraries, which handle these standards appropriately, exist for almost any
programming language.

These keys will be generated automatically if they do not exist yet in the database. No further steps for upgrading are
required.

#### Consent Flow

The consent flow has been refactored in order to implement session (login & consent) management in ORY Hydra and in order
to properly support OpenID Connect parameters such as `prompt`, `max_age`, and others.

First, the consent flow has been renamed to "User Login and Consent Flow". The consent app has been renamed to `User Login Provider`
and `User Consent Provider`. If you implement both features (explained in the next sections) in one program, you can call
it the `User Login and Consent Provider`.

A reference implementation of the new User Login and Consent Provider is available at
[github.com/ory/hydra-login-consent-node](https://github.com/ory/hydra-login-consent-node).

The major difference between the old and new flow is, that authentication (user login) and scope authorization (user consent)
are now two separate endpoints.

The new User Login and Consent Flow is documented in the [developer guide](https://www.ory.sh/docs/guides/latest/1-hydra/).

#### Changes to the CLI

The CLI has changed in order to improve developer experience and adopt to the changes made with this release.

##### `hydra host`

The command `hydra host` has been renamed to `hydra serve` as projects ORY Oathkeeper and ORY Keto use the `serve` terminology
as well.

Because this patch removes the internal access control, no root client and root policy will be created upon start up. Thus,
environment variable `FORCE_ROOT_CLIENT_CREDENTIALS` has been removed without replacement.

To better reflect what environment variables touch which system, ISSUER has been renamed to `OAUTH2_ISSUER_URL` and
`CONSENT_URL` has been renamed to `OAUTH2_CONSENT_URL`.

Additionally, flag `--dangerous-force-auto-logon` has been removed it has no effect any more.

##### `hydra connect`

The command `hydra connect` has been removed as it no longer serves a purpose now that the internal access control
has been removed. Every command you call now needs the environment variable `HYDRA_URL` (previously named `CLUSTER_URL`)
which should point to ORY Hydra's URL. Removing this command has an additional benefit - privileged client IDs and secrets
will no longer be stored in a plaintext file on your system if you use this command.

As access control has been removed, most commands (except `token user`, `token client`, `token revoke`, `token introspect`)
work without supplying any credentials at all. The listed exceptions support setting an OAuth 2.0 Client ID and Client Secret
using flags `--client-id` and `--client-secret` or environment variables `OAUTH2_CLIENT_ID` and `OAUTH2_CLIENT_SECRET`.

All other commands, such as `hydra clients create`, still support scenarios where you would need an OAuth2 Access Token.
In those cases, you can supply the access token using flag `--access-token` or environment variable `OAUTH2_ACCESS_TOKEN`.
Assuming that you would like to automate management in a protected scenario, you could do something like this:

```
$ token=$(hydra token client --client-id foo --client-secret bar --endpoint http://foobar)
$ hydra clients create --access-token $token ...
```

All commands now support the `--endpoint` flag which sets the `HYDRA_URL` in case you don't want to use environment variables.

##### `hydra token user`

Flags `--id` and `--secret` are now called `--client-id` and `--client-secret`.

##### `hydra token client`

Flags `--client-id` and `--client-secret` have been added.

Flag `--scopes` has been renamed to `--scope`.

##### `hydra token validate`

This command has been renamed to `hydra token introspect` to properly reflect that you are performing OAuth 2.0
Token Introspection.

Flags `--client-id` and `--client-secret` have been added.

Flag `--scopes` has been renamed to `--scope`.

##### `hydra clients create`

As OAuth 2.0 specifies that terminology `scope` does not have a plural `scopes`, we updated the places where the
incorrect `scopes` was used in order to provide a more consistent developer experience.

This command renamed flag `--allowed-scopes` to `--scope`.

##### `hydra migrate ladon`

This command is a relict of an old version of ORY Hydra which is, according to our metrics, not being used any more.

##### `hydra policies`

This command has moved to Keto. All commands work the same way, but you have to have Keto installed and replace `hydra`
with `keto`. For example `hydra policies create ...` is now `keto policies create ...`

##### `hydra groups`

This command has moved to Keto. All commands work the same way, but you have to have Keto installed and replace `hydra groups`
with `keto roles`. For example `hydra groups create ...` is now `keto roles create ...`

#### SDK

As the SDK is code-generated, and we are not specialists in every language, we have only documented changes to the Go API.
Please help improving this section by adding upgrade guides for the SDK you upgraded.

The following methods have been moved.

* The Access Control Policy SDK has moved to ORY Keto:
  * `CreatePolicy(body swagger.Policy) (*swagger.Policy, *swagger.APIResponse, error)` is now available via `github.com/ory/keto/sdk/go/keto`. The method signature has not changed, apart from types `github.com/ory/hydra/sdk/go/hydra/swagger` now being included from `github.com/ory/keto/sdk/go/keto/swagger`.
  * `DeletePolicy(id string) (*swagger.APIResponse, error)` is now available via `github.com/ory/keto/sdk/go/keto`. The method signature has not changed, apart from types `github.com/ory/hydra/sdk/go/hydra/swagger` now being included from `github.com/ory/keto/sdk/go/keto/swagger`.
  * `GetPolicy(id string) (*swagger.Policy, *swagger.APIResponse, error)` is now available via `github.com/ory/keto/sdk/go/keto`. The method signature has not changed, apart from types `github.com/ory/hydra/sdk/go/hydra/swagger` now being included from `github.com/ory/keto/sdk/go/keto/swagger`.
  * `ListPolicies(offset int64, limit int64) ([]swagger.Policy, *swagger.APIResponse, error)` is now available via `github.com/ory/keto/sdk/go/keto`. The method signature has not changed, apart from types `github.com/ory/hydra/sdk/go/hydra/swagger` now being included from `github.com/ory/keto/sdk/go/keto/swagger`.
  * `UpdatePolicy(id string, body swagger.Policy) (*swagger.Policy, *swagger.APIResponse, error)` is now available via `github.com/ory/keto/sdk/go/keto`. The method signature has not changed, apart from types `github.com/ory/hydra/sdk/go/hydra/swagger` now being included from `github.com/ory/keto/sdk/go/keto/swagger`.
* The Warden Group SDK has moved to Keto:
  - `AddMembersToGroup(id string, body swagger.GroupMembers) (*swagger.APIResponse, error)` is now `AddMembersToRole(id string, body swagger.RoleMembers) (*swagger.APIResponse, error)` and is now available via `github.com/ory/keto/sdk/go/keto`.
  - `CreateGroup(body swagger.Group) (*swagger.Group, *swagger.APIResponse, error)` is now `CreateRole(body swagger.Role) (*swagger.Role, *swagger.APIResponse, error` and is now available via `github.com/ory/keto/sdk/go/keto`.
  - `DeleteGroup(id string) (*swagger.APIResponse, error)` is now `DeleteRole(id string) (*swagger.APIResponse, error)` and is now available via `github.com/ory/keto/sdk/go/keto`.
  - `ListGroups(member string, limit, offset int64) ([]swagger.Group, *swagger.APIResponse, error)` is now `ListRoles(member string, limit int64, offset int64) ([]swagger.Role, *swagger.APIResponse, error)` and is now available via `github.com/ory/keto/sdk/go/keto`.
  - `GetGroup(id string) (*swagger.Group, *swagger.APIResponse, error)` is now `GetRole(id string) (*swagger.Role, *swagger.APIResponse, error)` and is now available via `github.com/ory/keto/sdk/go/keto`.
  - `RemoveMembersFromGroup(id string, body swagger.GroupMembers) (*swagger.APIResponse, error)` is now `RemoveMembersFromRole(id string, body swagger.RoleMembers) (*swagger.APIResponse, error)` and is now available via `github.com/ory/keto/sdk/go/keto`.
* The Warden API SDK has moved to Keto:
  - `DoesWardenAllowAccessRequest(body swagger.WardenAccessRequest) (*swagger.WardenAccessRequestResponse, *swagger.APIResponse, error)` is now `IsSubjectAuthorized(body swagger.WardenSubjectAuthorizationRequest) (*swagger.WardenSubjectAuthorizationResponse, *swagger.APIResponse, error)`. Please check out the changes to the request/response body as well.
  - `DoesWardenAllowTokenAccessRequest(body swagger.WardenTokenAccessRequest) (*swagger.WardenTokenAccessRequestResponse, *swagger.APIResponse, error)` is now `IsOAuth2AccessTokenAuthorized(body swagger.WardenOAuth2AccessTokenAuthorizationRequest) (*swagger.WardenOAuth2AccessTokenAuthorizationResponse, *swagger.APIResponse, error)`. Please check out the changes to the request/response body as well.
* The Consent API SDK has been deprecated:
  - `AcceptOAuth2ConsentRequest(id string, body swagger.ConsentRequestAcceptance) (*swagger.APIResponse, error)` has been removed without replacement.
  - `GetOAuth2ConsentRequest(id string) (*swagger.OAuth2ConsentRequest, *swagger.APIResponse, error)` has been removed without replacement.
  - `RejectOAuth2ConsentRequest(id string, body swagger.ConsentRequestRejection) (*swagger.APIResponse, error)` has been removed without replacement.
* The Login & Consent API SDK has been added:
  - `AcceptConsentRequest(challenge string, body swagger.AcceptConsentRequest) (*swagger.CompletedRequest, *swagger.APIResponse, error)`
  - `AcceptLoginRequest(challenge string, body swagger.AcceptLoginRequest) (*swagger.CompletedRequest, *swagger.APIResponse, error)`
  - `RejectConsentRequest(challenge string, body swagger.RejectRequest) (*swagger.CompletedRequest, *swagger.APIResponse, error)`
  - `RejectLoginRequest(challenge string, body swagger.RejectRequest) (*swagger.CompletedRequest, *swagger.APIResponse, error)`
  - `GetLoginRequest(challenge string) (*swagger.LoginRequest, *swagger.APIResponse, error)`
  - `GetConsentRequest(challenge string) (*swagger.ConsentRequest, *swagger.APIResponse, error)`

Additionally, the following methods have been removed as they were of very little use and also mixed the Client Credentials
flow with the Authorize Code Flow which lead to weird usage. It's much easier to configure `clientcredentials.Config` or
`oauth2.Config` yourself.

* `GetOAuth2ClientConfig() (*clientcredentials.Config)`
* `GetOAuth2Config() (*oauth2.Config)`

### Improvements

#### Health Check endpoint has moved

The health check endpoint has moved from `/health/status` to `/health/alive`. We set up a 308 redirect from `/health/status`
to `/health/alive` so this should not cause any issues.

The `/health/alive` endpoint returns `200 OK` as soon as the HTTP server is responsive. Another endpoint `/health/ready`
was added which returns `200 OK` only if the database connection is working as well.

As part of this change, function `getInstanceStatus` of the SDK is now `isInstanceAlive` and `isInstanceReady`.

#### Unknown request body payloads result in error

Previously, if you had a typo in the JSON (e.g. `client_nme` instead of `client_name`), ORY Hydra simply ignored that key.
Now, an error is thrown if unknown JSON keys are included.

#### UTC everywhere

ORY Hydra now uses UTC everywhere, reducing the possibility of errors stemming from different timezones.

#### Pagination everywhere

Each endpoint that returns a list of items now supports pagination using `limit` and `offset` query parameters.

#### Flushing old access tokens

An endpoint (`/oauth2/flush`) has been added that allows you to flush old access tokens.

#### Prometheus endpoint

An endpoint `/health/prometheus` for providing data to Prometheus has been added.

## 0.11.12

This release resolves a security issue (reported by [platform.sh](https://www.platform.sh)) related to the fosite
storage implementation in this project. Fosite used to pass all of the request body from both authorize and token
endpoints to the storage adapters. As some of these values are needed in consecutive requests, the storage adapter
of this project chose to drop all of the key/value pairs to the database in plaintext.

This implied that confidential parameters, such as the `client_secret` which can be passed in the request body since
fosite version 0.15.0, were stored as key/value pairs in plaintext in the database. While most client secrets are generated
programmatically (as opposed to set by the user) and most popular OAuth2 providers choose to store the secret in plaintext
for later retrieval, we see it as a considerable security issue nonetheless.

The issue has been resolved by sanitizing the request body and only including those values truly required by their
respective handlers. This also implies that typos (eg `client_secet`) won't "leak" to the database.

There are no special upgrade paths required for this version.

This issue does not apply to you if you do not use an SQL backend. If you do upgrade to this version, you need to run
`hydra migrate sql path://to.your/database`.

If your users use POST body client authentication, it might
be a good move to remove old data. There are multiple ways of doing that. **Back up your data before you do this**:

1. **Radical solution:** Drop all rows from tables `hydra_oauth2_refresh`, `hydra_oauth2_access`, `hydra_oauth2_oidc`,
`hydra_oauth2_code`. This implies that all your users have to re-authorize.
2. **Sensitive solution:** Replace all values in column `form_data` in tables `hydra_oauth2_refresh`, `hydra_oauth2_access` with
an empty string. This will keep all authorization sessions alive. Tables `hydra_oauth2_oidc` and `hydra_oauth2_code`
do not contain sensitive information, unless your users accidentally sent the client_secret to the `/oauth2/auth` endpoint.

We would like to thank [platform.sh](https://www.platform.sh) for sponsoring the development of a patch that resolves this
issue.

## 0.11.3

The experimental endpoint `/health/metrics` has been removed as it caused various issues such as increased memory usage,
and it was apparently not used at all.

## 0.11.0

This release has a minor breaking change in the experimental Warden Group SDK: 
`FindGroupsByMember(member string) ([]swagger.Group, *swagger.APIResponse, error)` is now
`ListGroups(member string, limit, offset int64) ([]swagger.Group, *swagger.APIResponse, error)`.
The change has to be applied in a similar fashion to other SDKs generated using swagger.

Leave the `member` parameter empty to list all groups, and add it to filter groups by member id.

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
